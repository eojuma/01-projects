package services

import (
	"database/sql"
	"net/http"

	"forum/internals/models"
)

func GetAuthenticatedUser(db *sql.DB, r *http.Request) (models.User, error) {
	cookie, err := r.Cookie("user_session")
	if err != nil {
		return models.User{}, err
	}

	session, err := GetSession(db, cookie.Value)
	if err != nil {
		return models.User{}, err
	}

	if IsSessionExpired(session.ExpiresAt) {
		DeleteSession(db, session.SessionID)
		return models.User{}, sql.ErrNoRows
	}

	user, err := GetUserByID(db, session.UserID)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
