package handlers
import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"forum/internals/models"
	"forum/internals/services"
	"forum/internals/utils"
)

// listen for both get and post requests
func CommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Get authenticated user
			user, err := services.GetAuthenticatedUser(db, r)
			if err != nil {
				http.Error(w, "You must be logged in to comment", http.StatusUnauthorized)
				return
			}

			commentText := r.FormValue("comment")
			postID := r.FormValue("post-id")
			if postID == "" {
				// fmt.Fprintln(w, "Post ID is required: ", http.StatusBadRequest)
				http.Error(w, "Post ID is required", http.StatusBadRequest)
				return
			}
			if strings.TrimSpace(commentText) == "" {
				// fmt.Fprintln(w, "Comment content can't be empty: ", http.StatusBadRequest)
				http.Error(w, "Comment content can't be empty", http.StatusBadRequest)
				return
			} else {
				comment := models.Comment{
					CommentID: utils.GenerateUUID(),
					UserID:    user.UserID,
					PostID:    postID,
					Content:   strings.TrimSpace(commentText),
					CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
				}
				err := services.AddComment(db, comment)
				if err != nil {
					// fmt.Fprintln(w, "Error adding comment: ", err.Error(), http.StatusInternalServerError)
					http.Error(w, "Error adding comment: "+err.Error(), http.StatusInternalServerError)
					return
				}
				// redirect back to the post view so the user sees their comment
				http.Redirect(w, r, "/post/view?id="+postID, http.StatusSeeOther)
			}

		}
	}
}

// CommentLikeHandler handles liking a comment
func CommentLikeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
		commentID := r.FormValue("comment_id")
		userID := r.FormValue("user_id")
		if commentID == "" || userID == "" {
			http.Error(w, "missing parameters", http.StatusBadRequest)
			return
		}

		// If user had disliked the comment, remove the dislike
		hasDisliked, err := services.HasUserDislikedComment(db, userID, commentID)
		if err != nil {
			http.Error(w, "Failed to check dislike status", http.StatusInternalServerError)
			return
		}
		if hasDisliked {
			if err := services.RemoveCommentDislike(db, userID, commentID); err != nil {
				http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
				return
			}
		}

		// If already liked, remove like (toggle)
		hasLiked, err := services.HasUserLikedComment(db, userID, commentID)
		if err != nil {
			http.Error(w, "Failed to check like status", http.StatusInternalServerError)
			return
		}
		if hasLiked {
			if err := services.RemoveCommentLike(db, userID, commentID); err != nil {
				http.Error(w, "Failed to remove like", http.StatusInternalServerError)
				return
			}
		} else {
			like := models.CommentLike{
				LikeID:    utils.GenerateUUID(),
				Counts:    1,
				UserID:    userID,
				CommentID: commentID,
			}
			if err := services.AddCommentLike(db, like); err != nil {
				http.Error(w, "Failed to add like", http.StatusInternalServerError)
				return
			}
		}

		// redirect back to the post view (comment form uses post-id hidden field)
		postID := r.FormValue("post_id")
		if postID == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/post/view?id="+postID, http.StatusSeeOther)
	}
}

// CommentDislikeHandler handles disliking a comment
func CommentDislikeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RenderError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
		commentID := r.FormValue("comment_id")
		userID := r.FormValue("user_id")
		if commentID == "" || userID == "" {
			http.Error(w, "missing parameters", http.StatusBadRequest)
			return
		}

		// If user had liked the comment, remove the like
		hasLiked, err := services.HasUserLikedComment(db, userID, commentID)
		if err != nil {
			http.Error(w, "Failed to check like status", http.StatusInternalServerError)
			return
		}
		if hasLiked {
			if err := services.RemoveCommentLike(db, userID, commentID); err != nil {
				http.Error(w, "Failed to remove like", http.StatusInternalServerError)
				return
			}
		}

		// Toggle dislike
		hasDisliked, err := services.HasUserDislikedComment(db, userID, commentID)
		if err != nil {
			http.Error(w, "Failed to check dislike status", http.StatusInternalServerError)
			return
		}
		if hasDisliked {
			if err := services.RemoveCommentDislike(db, userID, commentID); err != nil {
				http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
				return
			}
		} else {
			dis := models.CommentDislike{
				DislikeID: utils.GenerateUUID(),
				Counts:    1,
				UserID:    userID,
				CommentID: commentID,
			}
			if err := services.AddCommentDislike(db, dis); err != nil {
				http.Error(w, "Failed to add dislike", http.StatusInternalServerError)
				return
			}
		}

		postID := r.FormValue("post_id")
		if postID == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/post/view?id="+postID, http.StatusSeeOther)
	}
}
