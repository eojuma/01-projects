package config

import (
	"os"
	"time"
)

// Config holds application configuration
type Config struct {
	// Server
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration

	// Database
	DatabasePath string

	// Session
	SessionTimeout time.Duration
	CookieName     string
	CookieDomain   string
	// TLS
	TLSCertFile string
	TLSKeyFile  string

	// Environment
	Environment string
}

// New creates a new config with defaults and overrides from env vars
func New() *Config {
	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		Host:            getEnv("HOST", "0.0.0.0"),
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		DatabasePath:    getEnv("DATABASE_PATH", "social_network.db"),
		SessionTimeout:  24 * time.Hour,
		CookieName:      getEnv("COOKIE_NAME", "session_id"),
		CookieDomain:    getEnv("COOKIE_DOMAIN", ""),
		TLSCertFile:      getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:       getEnv("TLS_KEY_FILE", ""),
		Environment:     getEnv("ENVIRONMENT", "development"),
	}

	return cfg
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
