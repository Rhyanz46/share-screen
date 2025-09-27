// Mac → iPhone Screen Share (no‑login, LAN only)
// ------------------------------------------------
// Minimal Go + WebRTC signaling over HTTP (no WebSockets, no auth/login).
// Use on the same LAN. One sender (Mac) mirrors screen to one viewer (iPhone Safari).
//
// How to run:
// 1) `go run main.go`
// 2) On your Mac: open http://localhost:8080/sender and click "Start Share".
//    The page will show a Viewer URL (with a one-time token).
// 3) On your iPhone: open the Viewer URL in Safari. Boom — mirrored.
//
// Notes:
// - Uses `getDisplayMedia` (you choose which screen/window to share).
// - Codec is whatever Safari negotiates (H.264/VP8). No audio, just video.
// - LAN only by default; NAT traversal via Google STUN for convenience.
// - Single viewer at a time per token. No persistence, no login, no tracking.
// - This is intentionally bare-bones; tweak constraints or add PIN if needed.
//
// Tested on: macOS (Chrome/Safari) as sender, iOS Safari as viewer.
// ------------------------------------------------

package main

import (
	"log"
	"net/http"
	"time"

	"share-screen/pkg/infrastructure/template"
	"share-screen/pkg/infrastructure/config"
	"share-screen/pkg/infrastructure/network"
	"share-screen/pkg/infrastructure/repository"
	httphandlers "share-screen/pkg/presentation/http"
	"share-screen/pkg/usecase/usecases"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize dependencies following Clean Architecture
	dependencies := initializeDependencies(cfg)

	// Setup routes
	setupRoutes(dependencies.staticHandlers, dependencies.apiHandlers)

	// Start background services
	startBackgroundServices(dependencies.sessionRepo, cfg.TokenExpiry)

	// Start server
	startServer(cfg)
}

// Dependencies holds all application dependencies
type Dependencies struct {
	sessionRepo       *repository.MemorySessionRepository
	networkService    *network.NetworkService
	templateService   *template.TemplateService
	sessionUseCase    *usecases.SessionUseCase
	serverInfoUseCase *usecases.ServerInfoUseCase
	staticHandlers    *httphandlers.StaticHandlers
	apiHandlers       *httphandlers.APIHandlers
}

// initializeDependencies sets up dependency injection following Clean Architecture
func initializeDependencies(cfg *config.Config) *Dependencies {
	// Infrastructure Layer
	sessionRepo := repository.NewMemorySessionRepository().(*repository.MemorySessionRepository)
	networkService := network.NewNetworkService().(*network.NetworkService)

	templateService, err := template.NewTemplateService("web/templates", cfg.STUNServer)
	if err != nil {
		log.Fatalf("Failed to initialize template service: %v", err)
	}

	// Use Case Layer
	sessionUseCase := usecases.NewSessionUseCase(sessionRepo, cfg.TokenExpiry)
	serverInfoUseCase := usecases.NewServerInfoUseCase(networkService, cfg.STUNServer, "1.0.0")

	// Presentation Layer
	staticHandlers := httphandlers.NewStaticHandlers(templateService)
	apiHandlers := httphandlers.NewAPIHandlers(sessionUseCase, serverInfoUseCase)

	return &Dependencies{
		sessionRepo:       sessionRepo,
		networkService:    networkService,
		templateService:   templateService,
		sessionUseCase:    sessionUseCase,
		serverInfoUseCase: serverInfoUseCase,
		staticHandlers:    staticHandlers,
		apiHandlers:       apiHandlers,
	}
}

// startBackgroundServices starts background processes like garbage collection
func startBackgroundServices(sessionRepo *repository.MemorySessionRepository, tokenExpiry time.Duration) {
	// Start garbage collection for expired sessions
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		log.Printf("🗑️  Token garbage collector started (cleanup every 1 min, expiry: %v)", tokenExpiry)

		for range ticker.C {
			sessionRepo.CleanupExpiredSessions()
		}
	}()
}

// setupRoutes configures all HTTP routes
func setupRoutes(static *httphandlers.StaticHandlers, api *httphandlers.APIHandlers) {
	// Static pages
	http.HandleFunc("/", static.ServeIndex)
	http.HandleFunc("/sender", static.ServeSender)
	http.HandleFunc("/viewer", static.ServeViewer)

	// Static assets (CSS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	// Dynamic JavaScript (with template rendering)
	http.HandleFunc("/static/js/sender.js", static.ServeSenderJS)
	http.HandleFunc("/static/js/viewer.js", static.ServeViewerJS)

	// API endpoints
	http.HandleFunc("/api/new", api.HandleNewToken)
	http.HandleFunc("/api/offer", api.HandleOffer)
	http.HandleFunc("/api/answer", api.HandleAnswer)
	http.HandleFunc("/api/info", api.HandleInfo)
}

// startServer starts the HTTP or HTTPS server based on configuration
func startServer(cfg *config.Config) {
	addr := ":" + cfg.Port
	protocol := "HTTP"
	if cfg.EnableHTTPS {
		protocol = "HTTPS"
	}

	log.Printf("%s Server listening on %s", protocol, addr)
	log.Printf("STUN Server: %s", cfg.STUNServer)
	log.Printf("Token Expiry: %s", cfg.TokenExpiry)

	var err error
	if cfg.EnableHTTPS {
		log.Printf("TLS Certificate: %s", cfg.CertFile)
		log.Printf("TLS Private Key: %s", cfg.KeyFile)
		err = http.ListenAndServeTLS(addr, cfg.CertFile, cfg.KeyFile, nil)
	} else {
		log.Printf("⚠️  Running in HTTP mode - consider enabling HTTPS for production")
		err = http.ListenAndServe(addr, nil)
	}

	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
