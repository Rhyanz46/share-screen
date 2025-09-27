package config

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	Port        string
	STUNServer  string
	TokenExpiry time.Duration
	EnableHTTPS bool
	CertFile    string
	KeyFile     string
}

// LoadConfig loads configuration from environment variables and command line flags
func LoadConfig() *Config {
	// Load .env file first
	loadEnv()

	// Define flags
	port := flag.String("port", "8080", "Server port")
	stunServer := flag.String("stun", "stun:stun.l.google.com:19302", "STUN server URL")
	tokenExpiry := flag.Duration("token-expiry", 30*time.Minute, "Token expiry duration")
	enableHTTPS := flag.Bool("https", false, "Enable HTTPS")
	certFile := flag.String("cert", "certs/server.crt", "Path to TLS certificate file")
	keyFile := flag.String("key", "certs/server.key", "Path to TLS private key file")
	flag.Parse()

	// Override with environment variables
	if envPort := os.Getenv("PORT"); envPort != "" {
		*port = envPort
	}
	if envStun := os.Getenv("STUN_SERVER"); envStun != "" {
		*stunServer = envStun
	}
	if envExpiry := os.Getenv("TOKEN_EXPIRY"); envExpiry != "" {
		if duration, err := time.ParseDuration(envExpiry); err == nil {
			*tokenExpiry = duration
		}
	}
	if envHTTPS := os.Getenv("ENABLE_HTTPS"); envHTTPS != "" {
		*enableHTTPS = envHTTPS == "true"
	}
	if envCert := os.Getenv("TLS_CERT_FILE"); envCert != "" {
		*certFile = envCert
	}
	if envKey := os.Getenv("TLS_KEY_FILE"); envKey != "" {
		*keyFile = envKey
	}

	return &Config{
		Port:        *port,
		STUNServer:  *stunServer,
		TokenExpiry: *tokenExpiry,
		EnableHTTPS: *enableHTTPS,
		CertFile:    *certFile,
		KeyFile:     *keyFile,
	}
}

// loadEnv loads environment variables from .env file
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // .env file not found, continue with defaults
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if err := os.Setenv(key, value); err != nil {
				log.Printf("Error setting env var %s: %v", key, err)
			}
		}
	}
}
