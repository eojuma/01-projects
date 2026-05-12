package services
 
import (
	"testing"
 
	"forum/internals/models"
)
 
func TestAddLike_ValidLike_NoError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	err := AddLike(db, models.Like{
		LikeID: "like-1",
		Counts: 1,
		UserID: user.UserID,
		PostID: postID,
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
 
func TestHasUserLiked_ReturnsFalse_WhenNoLike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	liked, err := HasUserLiked(db, "user-1", "post-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if liked {
		t.Error("expected false, got true")
	}
}
 
func TestHasUserLiked_ReturnsTrue_AfterLike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	AddLike(db, models.Like{LikeID: "like-1", Counts: 1, UserID: user.UserID, PostID: postID})
 
	liked, err := HasUserLiked(db, user.UserID, postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !liked {
		t.Error("expected true after adding like")
	}
}
 
func TestRemoveLike_RemovesLike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	AddLike(db, models.Like{LikeID: "like-1", Counts: 1, UserID: user.UserID, PostID: postID})
 
	err := RemoveLike(db, user.UserID, postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
 
	liked, _ := HasUserLiked(db, user.UserID, postID)
	if liked {
		t.Error("expected false after removing like")
	}
}
 
func TestGetMyLikes_ReturnsPostsLikedByUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	AddLike(db, models.Like{LikeID: "like-1", Counts: 1, UserID: user.UserID, PostID: postID})
 
	posts, err := GetMyLikes(db, user.UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(posts))
	}
}
 
func TestGetMyPosts_ReturnsPostsByUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	seedPost(t, db, user.UserID)
 
	posts, err := GetMyPosts(db, user.UserID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("expected 1 post, got %d", len(posts))
	}
}