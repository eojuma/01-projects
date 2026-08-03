package models

import "time"

// Session represents a user session
type Session struct {
	ID        string    `db:"id"`
	UserID    int       `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}
