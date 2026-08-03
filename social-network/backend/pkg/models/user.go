package models

import "time"

// User represents a user in the system
type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Bio          string    `db:"bio"`
	AvatarURL    string    `db:"avatar_url"`
	Nickname     string    `db:"nickname"`
	DateOfBirth  string    `db:"date_of_birth"`
	IsPrivate    bool      `db:"is_private"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// UserResponse is the user data returned in API responses (excludes password)
type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	Nickname    string    `json:"nickname"`
	DateOfBirth string    `json:"date_of_birth"`
	IsPrivate   bool      `json:"is_private"`
	CreatedAt   time.Time `json:"created_at"`
}
