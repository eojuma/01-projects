package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"unicode"

	"forum/internals/models"
	"forum/internals/services"
	"forum/internals/utils"
)
 
func LikeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postID := r.FormValue("post_id")
			userID := r.FormValue("user_id")
 
			// check if user had already disliked the post, if so remove the dislike before adding the like
			hasDisliked, err := services.HasUserDisliked(db, userID, postID)
			if err != nil {
				http.Error(w, "Failed to check dislike status", http.StatusInternalServerError)
				return
			}
			if hasDisliked {
				err := services.RemoveDislike(db, userID, postID)
				if err != nil {
					http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
					return
				}
			}
 
			// toggle — if already liked, remove and return
			hasLiked, err := services.HasUserLiked(db, userID, postID)
			if err != nil {
				http.Error(w, "Failed to check like status", http.StatusInternalServerError)
				return
			}
			if hasLiked {
				err := services.RemoveLike(db, userID, postID)
				if err != nil {
					http.Error(w, "Failed to remove like", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
 
			like := models.Like{
				LikeID: utils.GenerateUUID(),
				PostID: postID,
				UserID: userID,
				Counts: 1,
			}
 
			err = services.AddLike(db, like)
			if err != nil {
				http.Error(w, "Failed to add like", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
 
func DislikeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postID := r.FormValue("post_id")
			userID := r.FormValue("user_id")
 
			// check if user had already liked the post, if so remove the like before adding the dislike
			hasLiked, err := services.HasUserLiked(db, userID, postID)
			if err != nil {
				http.Error(w, "Failed to check like status", http.StatusInternalServerError)
				return
			}
			if hasLiked {
				err := services.RemoveLike(db, userID, postID)
				if err != nil {
					http.Error(w, "Failed to remove like", http.StatusInternalServerError)
					return
				}
			}
 
			// toggle — if already disliked, remove and return
			hasDisliked, err := services.HasUserDisliked(db, userID, postID)
			if err != nil {
				http.Error(w, "Failed to check dislike status", http.StatusInternalServerError)
				return
			}
			if hasDisliked {
				err := services.RemoveDislike(db, userID, postID)
				if err != nil {
					http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
 
			dislike := models.Dislike{
				DislikeID: utils.GenerateUUID(),
				PostID:    postID,
				UserID:    userID,
				Counts:    1,
			}
 
			err = services.AddDislike(db, dislike)
			if err != nil {
				http.Error(w, "Failed to add dislike", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}
 
func MyLikesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			user, err := services.GetAuthenticatedUser(db, r)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
 
			posts, err := services.GetMyLikes(db, user.UserID)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, err.Error())
				return
			}
 
			initials := ""
			for _, v := range user.Username {
				if unicode.IsLetter(v) {
					initials += strings.ToUpper(string(v))
					break
				}
			}
 
			data := models.ResponseData{
				IsLoggedIn: true,
				Posts:      posts,
				User:       user,
				Initials:   initials,
			}
			renderTemplate(w, http.StatusOK, "home", data)
		} else {
			RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
	}
}
 
func MyPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			user, err := services.GetAuthenticatedUser(db, r)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
 
			posts, err := services.GetMyPosts(db, user.UserID)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, err.Error())
				return
			}
 
			initials := ""
			for _, v := range user.Username {
				if unicode.IsLetter(v) {
					initials += strings.ToUpper(string(v))
					break
				}
			}
 
			data := models.ResponseData{
				IsLoggedIn: true,
				Posts:      posts,
				User:       user,
				Initials:   initials,
			}
			renderTemplate(w, http.StatusOK, "home", data)
		} else {
			RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
	}
}