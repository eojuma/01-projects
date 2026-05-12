package services
 
import (
	"testing"
	"time"
 
	"forum/internals/models"
)
 
func TestAddPost_ValidPost_NoError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	user := seedUser(t, db)
 
	err := AddPost(db, models.Post{
		PostID: "p1", UserID: user.UserID, Title: "Title",
		Content: "Content", Category: "Tech",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
 
func TestViewPost_ReturnsPosts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	user := seedUser(t, db)
	seedPost(t, db, user.UserID)
 
	posts, err := ViewPost(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(posts))
	}
}
 
func TestGetPostById_ReturnsCorrectPost(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	post, err := GetPostById(db, postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.PostID != postID {
		t.Errorf("expected postID %s, got %s", postID, post.PostID)
	}
}
 
func TestGetPostById_InvalidID_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	_, err := GetPostById(db, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent post, got nil")
	}
}
 
func TestSearchPosts_ReturnsMatchingPosts(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	user := seedUser(t, db)
	seedPost(t, db, user.UserID)
 
	posts, err := SearchPosts(db, "Test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) == 0 {
		t.Error("expected at least 1 result for 'Test'")
	}
}