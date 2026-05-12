package handlers

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"forum/internals/models"
	"forum/internals/services"
	"forum/internals/utils"
)

// Initialize struct and assign it posts
var data struct {
	Posts []models.Post
}

func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. If the user is just trying to VIEW the page (clicking the link)
		if r.Method == http.MethodGet {
			cookie, err := r.Cookie("user_session")
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			// Replace "post.html" with the correct path to your template
			// Ensure you are using your template execution logic here
			renderTemplate(w, http.StatusOK, "post", nil)
			return
		}

		// 2. If the user is SUBMITTING the form
		if r.Method == http.MethodPost {
			user, err := services.GetAuthenticatedUser(db, r)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			// Parse multipart form (supports optional file upload)
			_ = r.ParseMultipartForm(10 << 20) // 10 MB

			// Collect all selected categories
			categories := r.Form["category"]
			categoryString := strings.Join(categories, ", ")

			newPost := models.Post{
				PostID:    utils.GenerateUUID(),
				UserID:    user.UserID,
				Title:     r.FormValue("title"),
				Content:   r.FormValue("content"),
				Category:  categoryString,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"), // Format time as string
				Image:     "",
			}

			// Handle optional file upload named "image"
			file, header, err := r.FormFile("image")
			if err == nil {
				defer file.Close()
				// ensure uploads dir exists
				uploadDir := "web/static/uploads"
				_ = os.MkdirAll(uploadDir, 0755)

				ext := filepath.Ext(header.Filename)
				if ext == "" {
					ext = ".img"
				}
				filename := newPost.PostID + ext
				dstPath := filepath.Join(uploadDir, filename)
				dst, err := os.Create(dstPath)
				if err == nil {
					defer dst.Close()
					_, _ = io.Copy(dst, file)
					// store public path for template usage
					newPost.Image = "/static/uploads/" + filename
				}
			}

			// Attempt writeDB function to write to table posts. Err will be nil if successful
			err = services.AddPost(db, newPost)
			if err != nil {

				http.Error(w, "Could not save post: "+err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Catch-all for other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func AllPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if data.Posts == nil {
			// Get all posts from DB
			posts, err := services.ViewPost(db)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
			data.Posts = posts
		}
		tmpl, err := template.ParseFiles("web/pages/base.html", "web/pages/home.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		// Execute template and pass in data
		tmpl.Execute(w, data.Posts)
	}
}

func FilterByCategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := services.ViewPost(db)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Database error")
			return
		}

		filtered := []models.Post{}
		category := r.URL.Query().Get("q")

		for _, v := range posts {
			if strings.Contains(strings.ToLower(v.Category), strings.ToLower(category)) {
				filtered = append(filtered, v)
			}
		}

		user, err := services.GetAuthenticatedUser(db, r)
		isLoggedIn := err == nil
		
		data := models.ResponseData{
			IsLoggedIn: isLoggedIn,
			Posts:      filtered,
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

func SortByLikes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := services.ViewByLikes(db)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Database error")
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

func SortByComments(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := services.ViewComment(db)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, "Database error")
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

// ViewPostHandler shows a single post and its comments
func ViewPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}

		postID := r.URL.Query().Get("id")
		if postID == "" {
			RenderError(w, http.StatusBadRequest, "post id is required")
			return
		}

		// get post
		post, err := services.GetPostById(db, postID)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// get comments for post
		comments, err := services.GetCommentsByPost(db, postID)
		if err != nil {
			RenderError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// prepare response data
		user, err := services.GetAuthenticatedUser(db, r)
		isLoggedIn := err == nil

		data := models.ResponseData{
			IsLoggedIn: isLoggedIn,
			Posts:      []models.Post{post},
			Comments:   comments,
			User:       user,
		}

		if isLoggedIn {
			// initials
			initials := ""
			for _, v := range data.User.Username {
				if unicode.IsLetter(v) {
					initials += strings.ToUpper(string(v))
					break
				}
			}
			data.Initials = initials
		}

		renderTemplate(w, http.StatusOK, "post_view", data)
	}
}
