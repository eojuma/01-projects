package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestHashAndComparePassword(t *testing.T) {
	password := "mysecurepassword"

	// Test hashing
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}

	// Test empty password hashing
	_, err = HashPassword("")
	if err == nil {
		t.Error("expected error when hashing empty password, got nil")
	}

	// Test correct comparison
	err = ComparePasswordHash(hash, password)
	if err != nil {
		t.Errorf("ComparePasswordHash failed with correct password: %v", err)
	}

	// Test incorrect comparison
	err = ComparePasswordHash(hash, "wrongpassword")
	if err == nil {
		t.Error("expected error when comparing incorrect password, got nil")
	}
}

func TestGenerateSessionID(t *testing.T) {
	sessionID, err := GenerateSessionID()
	if err != nil {
		t.Fatalf("failed to generate session ID: %v", err)
	}

	// SessionIDLength is 32 bytes, hex encoded it should be 64 characters
	expectedLength := SessionIDLength * 2
	if len(sessionID) != expectedLength {
		t.Errorf("expected session ID length to be %d, got %d", expectedLength, len(sessionID))
	}
}

func TestGetUserIDFromSession(t *testing.T) {
	// Setup in-memory sqlite db
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	defer testDB.Close()

	// Create sessions table
	_, err = testDB.Exec(`
		CREATE TABLE sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			expires_at INTEGER NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create sessions table: %v", err)
	}

	// Set db using SetDB
	SetDB(testDB)
	defer SetDB(nil)

	// Seed session
	sessionID := "testsessionid12345"
	userID := "user-99"
	expiresAt := time.Now().Add(1 * time.Hour).Unix()

	_, err = testDB.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", sessionID, userID, expiresAt)
	if err != nil {
		t.Fatalf("failed to seed session: %v", err)
	}

	// Test with valid cookie
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: sessionID,
	})

	gotUserID, err := GetUserIDFromSession(req)
	if err != nil {
		t.Fatalf("GetUserIDFromSession failed with valid cookie: %v", err)
	}
	if gotUserID != userID {
		t.Errorf("expected user ID %q, got %q", userID, gotUserID)
	}

	// Test with no cookie
	reqNoCookie := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err = GetUserIDFromSession(reqNoCookie)
	if err == nil {
		t.Error("expected error when no cookie is present, got nil")
	}

	// Test with expired session
	expiredSessionID := "expiredsessionid"
	expiredExpiresAt := time.Now().Add(-1 * time.Hour).Unix()
	_, err = testDB.Exec("INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)", expiredSessionID, userID, expiredExpiresAt)
	if err != nil {
		t.Fatalf("failed to seed expired session: %v", err)
	}

	reqExpired := httptest.NewRequest(http.MethodGet, "/", nil)
	reqExpired.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: expiredSessionID,
	})
	_, err = GetUserIDFromSession(reqExpired)
	if err == nil {
		t.Error("expected error when session is expired, got nil")
	}
}
