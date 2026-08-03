package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"social-network/backend/pkg/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// CreateUser inserts a new user and fills the generated ID.
func CreateUser(db *sql.DB, user *models.User) error {
	result, err := db.Exec(
		`INSERT INTO users (email, username, password_hash, first_name, last_name, bio, avatar_url, nickname, date_of_birth, is_private, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Bio,
		user.AvatarURL,
		user.Nickname,
		user.DateOfBirth,
		user.IsPrivate,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	// models.User.ID is a string, but sqlite LastInsertId is numeric.
	user.ID = fmt.Sprint(id)
	return nil
}

// GetUserByEmail returns a user by email.
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	return getUserByField(db, "email", email)
}

// GetUserByUsername returns a user by username.
func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	return getUserByField(db, "username", username)
}

// GetUserByID returns a user by ID.
func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(
		`SELECT id, email, username, password_hash, first_name, last_name, bio, avatar_url, nickname, date_of_birth, is_private, created_at, updated_at
		 FROM users WHERE id = ?`,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Bio,
		&user.AvatarURL,
		&user.Nickname,
		&user.DateOfBirth,
		&user.IsPrivate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByIDString returns a user by ID (string version)
func GetUserByIDString(db *sql.DB, id string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(
		`SELECT id, email, username, password_hash, first_name, last_name, bio, avatar_url, nickname, date_of_birth, is_private, created_at, updated_at
		 FROM users WHERE id = ?`,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Bio,
		&user.AvatarURL,
		&user.Nickname,
		&user.DateOfBirth,
		&user.IsPrivate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func getUserByField(db *sql.DB, field string, value interface{}) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, username, password_hash, first_name, last_name, bio, avatar_url, nickname, date_of_birth, is_private, created_at, updated_at FROM users WHERE ` + field + ` = ?`
	err := db.QueryRow(query, value).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Bio,
		&user.AvatarURL,
		&user.Nickname,
		&user.DateOfBirth,
		&user.IsPrivate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SearchUsers returns users matching query, excluding the current logged-in user.
func SearchUsers(db *sql.DB, currentUserID string, query string) ([]models.UserResponse, error) {
	searchParam := "%" + query + "%"
	rows, err := db.Query(
		`SELECT id, email, username, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(bio, ''), COALESCE(avatar_url, ''), COALESCE(nickname, ''), is_private
		 FROM users
		 WHERE id != ? AND (username LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR nickname LIKE ?)
		 LIMIT 50`,
		currentUserID, searchParam, searchParam, searchParam, searchParam,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL, &user.Nickname, &user.IsPrivate,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if users == nil {
		users = []models.UserResponse{}
	}
	return users, nil
}

// SearchChatUsers returns users matching query that are either followed by or following current user.
func SearchChatUsers(db *sql.DB, currentUserID string, query string) ([]models.UserResponse, error) {
	searchParam := "%" + query + "%"
	rows, err := db.Query(
		`SELECT id, email, username, COALESCE(first_name, ''), COALESCE(last_name, ''), COALESCE(bio, ''), COALESCE(avatar_url, ''), COALESCE(nickname, ''), is_private
		 FROM users
		 WHERE id != ? AND (username LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR nickname LIKE ?)
		   AND (
		       id IN (SELECT following_id FROM followers WHERE follower_id = ?)
		       OR id IN (SELECT follower_id FROM followers WHERE following_id = ?)
		   )
		 LIMIT 50`,
		currentUserID, searchParam, searchParam, searchParam, searchParam, currentUserID, currentUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(
			&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.Bio, &user.AvatarURL, &user.Nickname, &user.IsPrivate,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if users == nil {
		users = []models.UserResponse{}
	}
	return users, nil
}
