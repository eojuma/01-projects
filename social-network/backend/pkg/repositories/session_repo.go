package repositories

import (
	"database/sql"
	"strconv"
	"time"
)

type SessionRepository interface {
	Create(sessionID, userID string, expiresAt time.Time) error
	LookupUserID(sessionID string) (string, error)
	Delete(sessionID string) error
	DeleteExpired() error
}

type SQLiteSessionRepo struct {
	DB *sql.DB
}

func NewSQLiteSessionRepo(db *sql.DB) (*SQLiteSessionRepo, error) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS sessions (
        id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        expires_at INTEGER NOT NULL
    );`)
	if err != nil {
		return nil, err
	}
	return &SQLiteSessionRepo{DB: db}, nil
}

func (r *SQLiteSessionRepo) Create(sessionID, userID string, expiresAt time.Time) error {
	_, err := r.DB.Exec(`INSERT INTO sessions(id, user_id, expires_at) VALUES (?, ?, ?)`, sessionID, userID, expiresAt.Unix())
	return err
}

func (r *SQLiteSessionRepo) LookupUserID(sessionID string) (string, error) {
	var userID int64
	var expires time.Time
	err := r.DB.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE id = ?`, sessionID).Scan(&userID, &expires)
	if err != nil {
		return "", err
	}
	if expires.Before(time.Now()) {
		// expired — remove and report not found
		_ = r.Delete(sessionID)
		return "", sql.ErrNoRows
	}
	return strconv.FormatInt(userID, 10), nil
}

func (r *SQLiteSessionRepo) Delete(sessionID string) error {
	_, err := r.DB.Exec(`DELETE FROM sessions WHERE id = ?`, sessionID)
	return err
}

// DeleteExpired removes expired sessions from the store.
func (r *SQLiteSessionRepo) DeleteExpired() error {
	_, err := r.DB.Exec(`DELETE FROM sessions WHERE expires_at < ?`, time.Now().Unix())
	return err
}
