package middleware

import (
	"database/sql"
	"net/http"
	"forum/internals/services"
)

func Authorize(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := services.GetAuthenticatedUser(db, r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
