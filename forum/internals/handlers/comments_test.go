package handlers
 
import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
 
	"forum/internals/models"
	"forum/internals/services"
	"forum/internals/utils"
)
 
// -----------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------
 
func seedPost(t *testing.T, db *sql.DB, userID string) string {
	t.Helper()
	postID := utils.GenerateUUID()
	post := models.Post{
		PostID:    postID,
		UserID:    userID,
		Title:     "Test Post",
		Content:   "Test content",
		Category:  "Tech",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Image:     "",
	}
	err := services.AddPost(db, post)
	if err != nil {
		t.Fatalf("failed to seed post: %v", err)
	}
	return postID
}
 
func seedComment(t *testing.T, db *sql.DB, userID, postID string) string {
	t.Helper()
	commentID := utils.GenerateUUID()
	comment := models.Comment{
		CommentID: commentID,
		UserID:    userID,
		PostID:    postID,
		Content:   "Test comment",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := services.AddComment(db, comment)
	if err != nil {
		t.Fatalf("failed to seed comment: %v", err)
	}
	return commentID
}
 
func sessionCookie(t *testing.T, db *sql.DB, userID string) *http.Cookie {
	t.Helper()
	sessionID := utils.GenerateUUID()
	expiresAt := time.Now().Add(6 * time.Hour)
	err := services.AddSession(db, sessionID, userID, expiresAt)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	return &http.Cookie{Name: "user_session", Value: sessionID}
}
 
// -----------------------------------------------------------------------
// CommentHandler tests
// -----------------------------------------------------------------------
 
func TestCommentHandler_POST_Unauthenticated_ReturnsUnauthorized(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	form := url.Values{}
	form.Set("comment", "Some comment")
	form.Set("post-id", "some-post-id")
 
	req := httptest.NewRequest(http.MethodPost, "/post-comment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	CommentHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", rr.Code)
	}
}
 
func TestCommentHandler_POST_EmptyComment_ReturnsBadRequest(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
	cookie := sessionCookie(t, db, user.UserID)
 
	form := url.Values{}
	form.Set("comment", "   ")
	form.Set("post-id", postID)
 
	req := httptest.NewRequest(http.MethodPost, "/post-comment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rr := httptest.NewRecorder()
 
	CommentHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
 
func TestCommentHandler_POST_ValidComment_Redirects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
	cookie := sessionCookie(t, db, user.UserID)
 
	form := url.Values{}
	form.Set("comment", "This is a test comment")
	form.Set("post-id", postID)
 
	req := httptest.NewRequest(http.MethodPost, "/post-comment", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rr := httptest.NewRecorder()
 
	CommentHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rr.Code)
	}
}
 
// -----------------------------------------------------------------------
// CommentLikeHandler tests
// -----------------------------------------------------------------------
 
func TestCommentLikeHandler_InvalidMethod_ReturnsMethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodGet, "/comment/like", nil)
	rr := httptest.NewRecorder()
 
	CommentLikeHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
 
func TestCommentLikeHandler_POST_ValidLike_Redirects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
	commentID := seedComment(t, db, user.UserID, postID)
 
	form := url.Values{}
	form.Set("comment_id", commentID)
	form.Set("user_id", user.UserID)
	form.Set("post_id", postID)
 
	req := httptest.NewRequest(http.MethodPost, "/comment/like", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	CommentLikeHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rr.Code)
	}
}
 
// -----------------------------------------------------------------------
// CommentDislikeHandler tests
// -----------------------------------------------------------------------
 
func TestCommentDislikeHandler_InvalidMethod_ReturnsMethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodGet, "/comment/dislike", nil)
	rr := httptest.NewRecorder()
 
	CommentDislikeHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
 
func TestCommentDislikeHandler_POST_ToggleDislike_RemovesDislike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
	commentID := seedComment(t, db, user.UserID, postID)
 
	// dislike once
	services.AddCommentDislike(db, models.CommentDislike{
		DislikeID: utils.GenerateUUID(),
		Counts:    1,
		UserID:    user.UserID,
		CommentID: commentID,
	})
 
	// dislike again — should toggle off
	form := url.Values{}
	form.Set("comment_id", commentID)
	form.Set("user_id", user.UserID)
	form.Set("post_id", postID)
 
	req := httptest.NewRequest(http.MethodPost, "/comment/dislike", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	CommentDislikeHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 after toggle, got %d", rr.Code)
	}
 
	hasDisliked, _ := services.HasUserDislikedComment(db, user.UserID, commentID)
	if hasDisliked {
		t.Error("expected dislike to be removed after toggle")
	}
}