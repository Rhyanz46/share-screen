# Makefile for Share Screen project

.PHONY: help build run test clean docker-build docker-run docker-stop certs dev prod

# Default target
help: ## Show this help message
	@echo "Share Screen - Makefile Commands"
	@echo "================================="
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
build: ## Build the Go application
	@echo "Building share-screen..."
	@go build -o bin/share-screen main.go
	@echo "Build complete: bin/share-screen"

run: ## Run the application locally
	@echo "Starting share-screen..."
	@go run main.go

dev: ## Run in development mode with live reload (requires air)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Install air for live reload: go install github.com/cosmtrek/air@latest"; \
		make run; \
	fi

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -w .

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: fmt vet ## Run formatting and vet checks

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf certs/
	@echo "Clean complete"

# Certificate management
certs: ## Generate self-signed certificates for HTTPS
	@echo "Generating certificates..."
	@./scripts/generate-certs.sh

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t share-screen:latest .
	@echo "Docker image built: share-screen:latest"

docker-run: ## Run with Docker Compose (HTTP)
	@echo "Starting share-screen with Docker Compose (HTTP)..."
	@docker-compose up -d
	@echo "Application started at http://localhost:8080"

docker-run-https: certs ## Run with Docker Compose (HTTPS)
	@echo "Starting share-screen with Docker Compose (HTTPS)..."
	@docker-compose --profile https up -d share-screen-https
	@echo "Application started at https://localhost:8443"

docker-stop: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped"

docker-logs: ## Show Docker logs
	@docker-compose logs -f

docker-clean: docker-stop ## Clean Docker resources
	@echo "Cleaning Docker resources..."
	@docker-compose down -v --remove-orphans
	@docker image rm share-screen:latest 2>/dev/null || true
	@echo "Docker cleanup complete"

# Production deployment
prod-http: ## Run in production HTTP mode
	@echo "Starting in production HTTP mode..."
	@PORT=8080 ENABLE_HTTPS=false ./bin/share-screen

prod-https: certs ## Run in production HTTPS mode
	@echo "Starting in production HTTPS mode..."
	@PORT=8443 ENABLE_HTTPS=true ./bin/share-screen

# Quick start commands
start: docker-run ## Quick start with HTTP
start-https: docker-run-https ## Quick start with HTTPS

# Health check
health: ## Check if the service is running
	@echo "Checking service health..."
	@curl -f http://localhost:8080/ > /dev/null 2>&1 && echo "✅ HTTP service is healthy" || echo "❌ HTTP service is not responding"
	@curl -f -k https://localhost:8443/ > /dev/null 2>&1 && echo "✅ HTTPS service is healthy" || echo "❌ HTTPS service is not responding"

# Development setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	@go mod tidy
	@go mod download
	@mkdir -p bin/
	@mkdir -p logs/
	@cp .env.example .env 2>/dev/null || true
	@echo "Development environment setup complete"
	@echo ""
	@echo "Next steps:"
	@echo "1. Edit .env file for your configuration"
	@echo "2. Run 'make certs' to generate certificates for HTTPS"
	@echo "3. Run 'make start' for HTTP or 'make start-https' for HTTPS"

# Deployment commands
setup-deployment: ## Setup GitHub Actions deployment configuration
	@echo "Setting up deployment configuration..."
	@chmod +x scripts/setup-deployment.sh
	@./scripts/setup-deployment.sh

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

security-scan: ## Run security scan (requires gosec)
	@echo "Running security scan..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "Install gosec: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

ci-test: test-coverage security-scan ## Run all CI tests locally

deploy-check: ## Check deployment readiness
	@echo "Checking deployment readiness..."
	@echo "✅ Checking project structure..."
	@test -f main.go || (echo "❌ main.go not found" && exit 1)
	@test -d .github/workflows || (echo "❌ .github/workflows not found" && exit 1)
	@echo "✅ Checking scripts..."
	@test -x scripts/generate-certs.sh || (echo "❌ scripts/generate-certs.sh not executable" && exit 1)
	@test -x scripts/setup-deployment.sh || (echo "❌ scripts/setup-deployment.sh not executable" && exit 1)
	@echo "✅ Running tests..."
	@go test ./... > /dev/null || (echo "❌ Tests failing" && exit 1)
	@echo "✅ All checks passed! Ready for deployment."

# Quick deployment test
deploy-test: build ## Test deployment locally
	@echo "Testing deployment locally..."
	@./bin/share-screen --help || (echo "❌ Binary not working" && exit 1)
	@echo "✅ Local deployment test passed"