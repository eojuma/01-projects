package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const SessionIDLength = 32

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

func ComparePasswordHash(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}
	return nil
}

func GenerateSessionID() (string, error) {
	b := make([]byte, SessionIDLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate session id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func GetUserIDFromSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	var userID string
	var expiresAt int64
	err = db.QueryRow(`
        SELECT user_id, expires_at FROM sessions 
        WHERE id = ? AND expires_at > ?`,
		cookie.Value, time.Now().Unix(),
	).Scan(&userID, &expiresAt)
	if err != nil {
		return "", err
	}
	return userID, nil
}
