package services
 
import (
	"testing"
	"time"
 
	"forum/internals/models"
	dbpkg "forum/internals/db"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)
 
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	err = dbpkg.InitTables(db)
	if err != nil {
		t.Fatalf("failed to init tables: %v", err)
	}
	return db
}
 
func seedUser(t *testing.T, db *sql.DB) models.User {
	t.Helper()
	user := models.User{
		UserID:   "user-1",
		Email:    "test@example.com",
		Username: "testuser",
		Password: "hashedpassword",
	}
	err := AddUser(db, user)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	return user
}
 
func seedPost(t *testing.T, db *sql.DB, userID string) string {
	t.Helper()
	postID := "post-1"
	err := AddPost(db, models.Post{
		PostID:    postID,
		UserID:    userID,
		Title:     "Test Post",
		Content:   "Test content",
		Category:  "Tech",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		t.Fatalf("failed to seed post: %v", err)
	}
	return postID
}
 
// -----------------------------------------------------------------------
// AddComment tests
// -----------------------------------------------------------------------
 
func TestAddComment_ValidComment_NoError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	err := AddComment(db, models.Comment{
		CommentID: "comment-1",
		UserID:    user.UserID,
		PostID:    postID,
		Content:   "Great post!",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
 
// -----------------------------------------------------------------------
// GetCommentsByPost tests
// -----------------------------------------------------------------------
 
func TestGetCommentsByPost_ReturnsComments(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	AddComment(db, models.Comment{
		CommentID: "comment-1",
		UserID:    user.UserID,
		PostID:    postID,
		Content:   "First comment",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	})
 
	comments, err := GetCommentsByPost(db, postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}
}
 
func TestGetCommentsByPost_NoComments_ReturnsEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	comments, err := GetCommentsByPost(db, postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 0 {
		t.Errorf("expected 0 comments, got %d", len(comments))
	}
}
 
// -----------------------------------------------------------------------
// Comment like tests
// -----------------------------------------------------------------------
 
func TestHasUserLikedComment_ReturnsFalse_WhenNoLike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	liked, err := HasUserLikedComment(db, "user-1", "comment-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if liked {
		t.Error("expected false, got true")
	}
}
 
func TestAddCommentLike_AndHasUserLiked_ReturnsTrue(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
	AddComment(db, models.Comment{CommentID: "comment-1", UserID: user.UserID, PostID: postID, Content: "hi", CreatedAt: time.Now().Format("2006-01-02 15:04:05")})
 
	err := AddCommentLike(db, models.CommentLike{LikeID: "like-1", Counts: 1, UserID: user.UserID, CommentID: "comment-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
 
	liked, _ := HasUserLikedComment(db, user.UserID, "comment-1")
	if !liked {
		t.Error("expected true after adding like")
	}
}
 
func TestRemoveCommentLike_RemovesLike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
	AddComment(db, models.Comment{CommentID: "comment-1", UserID: user.UserID, PostID: postID, Content: "hi", CreatedAt: time.Now().Format("2006-01-02 15:04:05")})
	AddCommentLike(db, models.CommentLike{LikeID: "like-1", Counts: 1, UserID: user.UserID, CommentID: "comment-1"})
 
	err := RemoveCommentLike(db, user.UserID, "comment-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
 
	liked, _ := HasUserLikedComment(db, user.UserID, "comment-1")
	if liked {
		t.Error("expected false after removing like")
	}
}
 
// -----------------------------------------------------------------------
// Comment dislike tests
// -----------------------------------------------------------------------
 
func TestHasUserDislikedComment_ReturnsFalse_WhenNoDislike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	disliked, err := HasUserDislikedComment(db, "user-1", "comment-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if disliked {
		t.Error("expected false, got true")
	}
}
 
func TestRemoveCommentDislike_RemovesDislike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
	AddComment(db, models.Comment{CommentID: "comment-1", UserID: user.UserID, PostID: postID, Content: "hi", CreatedAt: time.Now().Format("2006-01-02 15:04:05")})
	AddCommentDislike(db, models.CommentDislike{DislikeID: "dislike-1", Counts: 1, UserID: user.UserID, CommentID: "comment-1"})
 
	err := RemoveCommentDislike(db, user.UserID, "comment-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
 
	disliked, _ := HasUserDislikedComment(db, user.UserID, "comment-1")
	if disliked {
		t.Error("expected false after removing dislike")
	}
}