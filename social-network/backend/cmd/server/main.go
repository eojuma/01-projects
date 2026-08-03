package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"social-network/backend/pkg/auth"
	"social-network/backend/pkg/config"
	"social-network/backend/pkg/db/migrations"
	"social-network/backend/pkg/handlers"
	"social-network/backend/pkg/middleware"
	"social-network/backend/pkg/repositories"
	"social-network/backend/pkg/routes"
	"social-network/backend/pkg/websocket"
)


type wsSessionStoreAdapter struct {
	repo repositories.SessionRepository
}

func (a *wsSessionStoreAdapter) LookupUserID(sessionID string) (string, error) {
	return a.repo.LookupUserID(sessionID)
}

// CORS middleware to allow frontend requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from frontend
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loadEnv loads environment variables from a .env file if it exists.
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // Ignore error if file doesn't exist
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Remove quotes if present
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}

func main() {
	loadEnv()
	cfg := config.New()

	db, err := sql.Open("sqlite3", "file:social_network.db?_foreign_keys=1")
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	defer db.Close()

	migrationsPath := filepath.Join("pkg", "db", "migrations", "sqlite")
	if err := migrations.Migrate(db, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	auth.SetDB(db)
	repositories.SetDatabase(db)
	repositories.OnNotificationCreated = func(notificationID, targetUserID string) {
		handlers.PushNotification(db, notificationID, targetUserID)
	}

	hub := websocket.NewHub()
	go hub.Run()

	sessRepo, err := repositories.NewSQLiteSessionRepo(db)
	if err != nil {
		log.Fatal(err)
	}
	websocket.SetSessionStore(&wsSessionStoreAdapter{repo: sessRepo})

	// Setup auth handler
	authHandler := &handlers.AuthHandler{
		DB:             db,
		SessionRepo:    sessRepo,
		CookieName:     cfg.CookieName,
		Timeout:        cfg.SessionTimeout,
		CookieDomain:   cfg.CookieDomain,
		CookieSecure:   cfg.Environment != "development",
		CookieSameSite: http.SameSiteLaxMode,
	}

	// Setup group handler
	groupHandler := &handlers.GroupHandler{DB: db}

	// Setup repositories and handlers
	mainRepo := repositories.NewRepository(db)
	followHandler := handlers.NewFollowHandler(mainRepo, hub)
	postHandler := handlers.NewPostHandler(mainRepo, "uploads/posts")
	profileHandler := handlers.GetProfileHandler(db)

	// Setup mux and routes
	mux := http.NewServeMux()
	routes.SetupRoutes(mux, db, authHandler, groupHandler, profileHandler, followHandler, postHandler)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// Apply rate limiting and session middleware
	rlWrapped := middleware.RateLimitMiddleware(mux, 120, time.Minute)
	sessionWrapped := middleware.SessionMiddleware(rlWrapped, sessRepo, cfg.CookieName)

	// Wrap with CORS last (outermost)
	finalHandler := corsMiddleware(sessionWrapped)

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			_ = sessRepo.DeleteExpired()
		}
	}()

	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
		log.Printf("Server listening at https://%s (TLS)", addr)
		if err := http.ListenAndServeTLS(addr, cfg.TLSCertFile, cfg.TLSKeyFile, finalHandler); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
		return
	}

	log.Printf("Server listening at http://%s\n", addr)
	if err := http.ListenAndServe(addr, finalHandler); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
