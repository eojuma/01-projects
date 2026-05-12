package services

import (
	"database/sql"
	"forum/internals/models"
	"time"
)

func AddSession(db *sql.DB, sessionID string, userID string, expiresAt time.Time) error {
	// Enforce single active session per user: remove any existing sessions
	_, err := db.Exec("DELETE FROM sessions WHERE userId = ?", userID)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO sessions (sessionId, userId, expiresAt) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt.Format(time.RFC3339))
	return err
}

func GetSession(db *sql.DB, sessionID string) (models.Session, error) {
	var s models.Session
	err := db.QueryRow("SELECT sessionId, userId, expiresAt FROM sessions WHERE sessionId = ?", sessionID).
		Scan(&s.SessionID, &s.UserID, &s.ExpiresAt)
	return s, err
}

func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE sessionId = ?", sessionID)
	return err
}

func IsSessionExpired(expiresAt string) bool {
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return true
	}
	return time.Now().After(t)
}

func GetUserByID(db *sql.DB, userID string) (models.User, error) {
	var u models.User
	err := db.QueryRow("SELECT userId, email, username, password, profile FROM users WHERE userId = ?", userID).
		Scan(&u.UserID, &u.Email, &u.Username, &u.Password, &u.ProfilePicture)
	return u, err
}
