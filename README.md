# Share Screen 🖥️ ➡️ 📱

A minimal, secure screen sharing application for Mac to iPhone using WebRTC. Now **production-ready** with HTTPS support and Docker containerization.

## ✨ Features

- **Zero-login screen sharing** from Mac to iPhone
- **WebRTC-based** real-time video streaming
- **HTTPS support** with self-signed certificates
- **Docker containerization** for easy deployment
- **LAN-optimized** for local network usage
- **One-time tokens** for secure sessions
- **Automatic cleanup** of expired sessions
- **Cross-platform** browser support

## 🚀 Quick Start

### Using Docker (Recommended)

1. **Clone and setup:**
   ```bash
   git clone <repository-url>
   cd share-screen
   make setup
   ```

2. **Start with HTTP:**
   ```bash
   make start
   ```
   Open http://localhost:8080

3. **Start with HTTPS:**
   ```bash
   make start-https
   ```
   Open https://localhost:8443

### Using Go directly

1. **Install and run:**
   ```bash
   go mod tidy
   go run main.go
   ```

2. **Or build and run:**
   ```bash
   make build
   ./bin/share-screen
   ```

## 🔒 HTTPS Configuration

### Development (Self-signed certificates)

Generate certificates automatically:
```bash
make certs
```

Or manually:
```bash
./scripts/generate-certs.sh
```

### Production (Trusted certificates)

1. **Obtain certificates** from a trusted CA (Let's Encrypt, etc.)

2. **Update environment:**
   ```bash
   export ENABLE_HTTPS=true
   export TLS_CERT_FILE=/path/to/your/certificate.crt
   export TLS_KEY_FILE=/path/to/your/private.key
   export PORT=8443
   ```

3. **Run with HTTPS:**
   ```bash
   make prod-https
   ```

## 🐳 Docker Deployment

### HTTP Mode
```bash
docker-compose up -d
```

### HTTPS Mode
```bash
docker-compose --profile https up -d share-screen-https
```

### Environment Variables
```bash
# Copy and modify
cp .env.example .env
```

Key variables:
- `ENABLE_HTTPS=true/false`
- `PORT=8080` (HTTP) or `8443` (HTTPS)
- `TLS_CERT_FILE=/path/to/cert.crt`
- `TLS_KEY_FILE=/path/to/private.key`
- `STUN_SERVER=stun:stun.l.google.com:19302`
- `TOKEN_EXPIRY=30m`

## 📖 Usage

1. **Start the server** (any method above)

2. **On Mac:**
   - Open sender URL in browser
   - Click "Start Share"
   - Choose screen/window to share
   - Copy the viewer URL

3. **On iPhone:**
   - Open the viewer URL in Safari
   - Video will start automatically

## 🔧 Development

### Prerequisites
- Go 1.23.3+
- Docker & Docker Compose
- OpenSSL (for certificate generation)

### Commands
```bash
make help           # Show all available commands
make setup          # Setup development environment
make build          # Build the application
make run            # Run locally
make test           # Run tests
make lint           # Format and vet code
make certs          # Generate certificates
make docker-build   # Build Docker image
make health         # Check service health
```

### Project Structure
```
share-screen/
├── main.go                 # Main application
├── Dockerfile             # Docker configuration
├── docker-compose.yml     # Docker Compose setup
├── Makefile              # Build automation
├── .env.example          # Environment template
├── scripts/
│   └── generate-certs.sh # Certificate generation
└── certs/                # Generated certificates
    ├── server.crt
    └── server.key
```

## 🌐 Network Requirements

- **Same WiFi network** for Mac and iPhone
- **Firewall rules** allowing chosen port (8080/8443)
- **STUN server access** for NAT traversal (configurable)

## 🔐 Security Features

- **Token-based authentication** (12-character random tokens)
- **Session expiration** (configurable, default 30 minutes)
- **HTTPS support** with TLS 1.2+
- **No persistent storage** of sessions
- **Rate limiting ready** (can be added)
- **Security headers** (can be enhanced)

## 📊 Production Considerations

### Already Implemented ✅
- ✅ HTTPS support with TLS certificates
- ✅ Docker containerization
- ✅ Environment-based configuration
- ✅ Graceful error handling
- ✅ Security headers basic implementation
- ✅ Input validation for tokens
- ✅ Automatic session cleanup

### Recommended Additions 📋
- [ ] Rate limiting middleware
- [ ] Structured logging (logrus/zap)
- [ ] Metrics collection (Prometheus)
- [ ] Health check endpoints
- [ ] Load balancer support
- [ ] Database persistence (if needed)
- [ ] Authentication system (if required)
- [ ] CORS configuration
- [ ] Request size limits

## 🚀 Deployment Options

### 1. Docker Compose (Simple)
```bash
docker-compose up -d
```

### 2. Docker Swarm (Scalable)
```bash
docker stack deploy -c docker-compose.yml share-screen
```

### 3. Kubernetes (Enterprise)
Create manifests based on the Docker configuration.

### 4. Direct Binary (Lightweight)
```bash
make build
PORT=8443 ENABLE_HTTPS=true ./bin/share-screen
```

## 🐛 Troubleshooting

### Certificate Issues
```bash
# Regenerate certificates
make clean && make certs

# Check certificate validity
openssl x509 -in certs/server.crt -text -noout
```

### Docker Issues
```bash
# Clean rebuild
make docker-clean && make docker-build

# Check logs
make docker-logs
```

### WebRTC Issues
- Ensure both devices on same network
- Check STUN server connectivity
- Verify browser compatibility (Chrome/Safari recommended)

## 📜 License

This project is provided as-is for educational and development purposes.

## 🤝 Contributing

1. Fork the repository
2. Create feature branch
3. Make changes with tests
4. Submit pull request

---

**Ready for production!** 🎉

Start with `make setup && make start-https` for the full HTTPS experience.
