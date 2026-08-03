package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"social-network/backend/pkg/models"
	"social-network/backend/pkg/repositories"

	"github.com/gofrs/uuid"
)

type PostHandler struct {
	Repo      *repositories.Repository
	UploadDir string
}

func NewPostHandler(r *repositories.Repository, uploadDir string) *PostHandler {
	os.MkdirAll(uploadDir, 0755)
	return &PostHandler{Repo: r, UploadDir: uploadDir}
}

// POST /api/posts
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	authorID, ok := r.Context().Value("userID").(string)
	if !ok || authorID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Strictly restricts the entire network request payload to 10MB maximum
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))

	privacy := r.FormValue("privacy")
	if privacy != "public" && privacy != "followers" && privacy != "selected" {
		http.Error(w, "privacy must be public, followers, or selected", http.StatusBadRequest)
		return
	}

	allowedUsers := r.Form["allowed_users[]"]
	if privacy == "selected" && len(allowedUsers) == 0 {
		http.Error(w, "allowed_users[] is required for selected privacy", http.StatusBadRequest)
		return
	}

	// Handle images (allow multiple uploads)
	imagePaths, err := SaveUploadedImages(r, "image", h.UploadDir)
	if err != nil {
		if errors.Is(err, ErrUnsupportedMediaType) {
			http.Error(w, "image must be JPEG, PNG, or GIF", http.StatusUnsupportedMediaType)
			return
		}
		http.Error(w, "failed to process image", http.StatusInternalServerError)
		return
	}

	// Require either content or image
	if content == "" && len(imagePaths) == 0 {
		http.Error(w, "content or image is required", http.StatusBadRequest)
		return
	}

	post, err := h.Repo.CreatePost(models.CreatePostInput{
		ID:           uuid.Must(uuid.NewV4()).String(),
		AuthorID:     authorID,
		Content:      content,
		Privacy:      privacy,
		ImagePaths:   imagePaths,
		AllowedUsers: allowedUsers,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch any images associated with this post so frontend can display them immediately
	var imageURLs []string
	rows, qerr := h.Repo.DB.Query(`SELECT image_path FROM post_images WHERE post_id = ?`, post.ID)
	if qerr == nil {
		defer rows.Close()
		for rows.Next() {
			var p string
			if err := rows.Scan(&p); err == nil {
				if strings.HasPrefix(p, "/") {
					imageURLs = append(imageURLs, p)
				} else {
					imageURLs = append(imageURLs, "/"+p)
				}
			}
		}
	}

	resp := map[string]interface{}{
		"id":              post.ID,
		"author_id":       post.AuthorID,
		"author_username": post.AuthorUsername,
		"content":         post.Content,
		"privacy":         post.Privacy,
		"created_at":      post.CreatedAt,
		"image_urls":      imageURLs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
func (h *PostHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	sessionUserID, ok := r.Context().Value("userID").(string)
	if !ok || sessionUserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	feed, err := h.Repo.GetFeed(sessionUserID, cursor, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}

// POST /api/posts/{id}/comments
func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	authorID, ok := r.Context().Value("userID").(string)
	if !ok || authorID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postID := r.PathValue("id")
	if postID == "" {
		http.Error(w, "missing post id", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))

	// Handle optional image — same validation as posts
	imagePath, err := SaveUploadedImage(r, "image", h.UploadDir)
	if err != nil {
		if errors.Is(err, ErrUnsupportedMediaType) {
			http.Error(w, "image must be JPEG, PNG, or GIF", http.StatusUnsupportedMediaType)
			return
		}
		http.Error(w, "failed to process image", http.StatusInternalServerError)
		return
	}

	if content == "" && imagePath == "" {
		http.Error(w, "content or image is required", http.StatusBadRequest)
		return
	}

	comment, err := h.Repo.CreateComment(authorID, models.CreateCommentInput{
		ID:        uuid.Must(uuid.NewV4()).String(),
		PostID:    postID,
		AuthorID:  authorID,
		Content:   content,
		ImagePath: imagePath,
	})
	if err != nil {
		switch err.Error() {
		case "post not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "forbidden":
			http.Error(w, "you do not have permission to view this post", http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GET /api/users/{id}/posts?cursor=...
func (h *PostHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	sessionUserID, ok := r.Context().Value("userID").(string)
	if !ok || sessionUserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	targetUserID := r.PathValue("id")
	if targetUserID == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	cursor := r.URL.Query().Get("cursor")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	posts, err := h.Repo.GetUserPosts(sessionUserID, targetUserID, cursor, limit)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// GET /api/posts/{id}/comments
func (h *PostHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	sessionUserID, ok := r.Context().Value("userID").(string)
	if !ok || sessionUserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postID := r.PathValue("id")
	if postID == "" {
		http.Error(w, "missing post id", http.StatusBadRequest)
		return
	}

	comments, err := h.Repo.GetComments(sessionUserID, postID)
	if err != nil {
		switch err.Error() {
		case "post not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "forbidden":
			http.Error(w, "you do not have permission to view comments for this post", http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// POST /api/posts/{id}/like
func (h *PostHandler) ToggleLike(w http.ResponseWriter, r *http.Request) {
	sessionUserID, ok := r.Context().Value("userID").(string)
	if !ok || sessionUserID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	postID := r.PathValue("id")
	if postID == "" {
		http.Error(w, "missing post id", http.StatusBadRequest)
		return
	}

	liked, count, err := h.Repo.ToggleLike(sessionUserID, postID)
	if err != nil {
		switch err.Error() {
		case "post not found":
			http.Error(w, err.Error(), http.StatusNotFound)
		case "forbidden":
			http.Error(w, "you do not have permission to like this post", http.StatusForbidden)
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"liked":       liked,
		"likes_count": count,
	})
}
