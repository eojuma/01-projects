package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"unicode"

	"forum/internals/models"
	"forum/internals/services"
)

func SearchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		posts, err := services.SearchPosts(db, query)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Error searching posts")
			return
		}

		user, err := services.GetAuthenticatedUser(db, r)
		isLoggedIn := err == nil

		data := models.ResponseData{
			IsLoggedIn: isLoggedIn,
			Posts:      posts,
			User:       user,
		}

		if isLoggedIn {
			initials := ""
			for _, v := range data.User.Username {
				if unicode.IsLetter(v) {
					initials += strings.ToUpper(string(v))
					break
				}
			}
			data.Initials = initials
		}

		renderTemplate(w, http.StatusOK, "home", data)
	}
}
