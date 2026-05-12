package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"
	"time"
	"unicode"

	"forum/internals/models"
	"forum/internals/services"
	"forum/internals/utils"
)

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			email := r.FormValue("email")
			password := r.FormValue("password")
			// get user from db
			user, err := services.GetUser(db, email, "")
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "We couldn't access your account data right now.")
				return
			}
			// compare password and hash
			if !utils.ComparePass(password, user.Password) {
				error := "Incorrect username or password"
				renderTemplate(w, http.StatusUnauthorized, "login", map[string]any{"Error": error})
				return
			}

			// Create session
			sessionID := utils.GenerateUUID()
			expiresAt := time.Now().Add(6 * time.Hour)
			err = services.AddSession(db, sessionID, user.UserID, expiresAt)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "Failed to create session")
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "user_session",
				Value:    sessionID,
				Path:     "/",
				Expires:  expiresAt,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else if r.Method == "GET" {
			renderTemplate(w, http.StatusOK, "login", models.ResponseData{IsLoggedIn: false})
		} else {
			RenderError(w, http.StatusBadRequest, "Bad request")
		}
	}
}

func RegisterHandler(database *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			email := r.FormValue("email")
			username := r.FormValue("username")
			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirm_password")

			// Validate input format
			if err := utils.ValidateEmail(email); err != nil {
				renderTemplate(w, http.StatusBadRequest, "register", map[string]any{"Error": err.Error()})
				return
			}
			if err := utils.ValidateUsername(username); err != nil {
				renderTemplate(w, http.StatusBadRequest, "register", map[string]any{"Error": err.Error()})
				return
			}
			if err := utils.ValidatePassword(password); err != nil {
				renderTemplate(w, http.StatusBadRequest, "register", map[string]any{"Error": err.Error()})
				return
			}

			if password != confirmPassword {
				renderTemplate(w, http.StatusBadRequest, "register", map[string]any{"Error": "Passwords do not match"})
				return
			}

			hashedPass, _ := utils.HashPass(password)

			newuser := models.User{
				Username:       username,
				Password:       hashedPass,
				Email:          strings.ToLower(email),
				UserID:         utils.GenerateUUID(),
				ProfilePicture: "",
			}

			err := services.AddUser(database, newuser)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "Failed to create your account. Please try again.")
				return
			}
			// Registration successful (temp credentials)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		} else if r.Method == "GET" {
			renderTemplate(w, http.StatusOK, "register", models.ResponseData{IsLoggedIn: false})
			return
		} else {
			RenderError(w, http.StatusMethodNotAllowed, "Invalid Request Method")
			return
		}
	}
}

func Logout(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_session")
		if err == nil && cookie.Value != "" {
			services.DeleteSession(db, cookie.Value)
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "user_session",
			Value:   "",
			Path:    "/",
			MaxAge:  -1,
			Expires: time.Unix(0, 0),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Helper function to render templates to avoid code duplication
func renderTemplate(w http.ResponseWriter, code int, tmpl string, data interface{}) {
	// You must parse the base and the specific page together
	files := []string{"web/pages/base.html", "web/pages/" + tmpl + ".html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// We execute "base" because that's the wrapper defined in base.html
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, err.Error(), code)
	}
}

func HomeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			RenderError(w, http.StatusNotFound, "The page you are looking for does not exist.")
			return
		}

		user, err := services.GetAuthenticatedUser(db, r)
		isLoggedIn := err == nil

		initialPosts, err := services.ViewPost(db) // pass db, not nil
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Database connection failed.")
			return
		}

		data := models.ResponseData{
			IsLoggedIn: isLoggedIn,
			Posts:      initialPosts,
			User:       user,
		}

		if isLoggedIn {
			initials := ""
			// Get first Letter and change it to caps
			for _, v := range user.Username {
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

func RenderError(w http.ResponseWriter, status int, message string) {
	// If no message is provided, use the official HTTP text
	if message == "" {
		message = http.StatusText(status)
	}

	data := models.ResponseData{
		IsLoggedIn: false,
		StatusCode: status,
		Error:      message,
	}

	// Set the actual HTTP header status
	w.WriteHeader(status)
	// Render the new error page
	renderTemplate(w, status, "error", data)
}
