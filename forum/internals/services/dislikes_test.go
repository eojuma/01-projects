package services
 
import (
	"testing"
 
	"forum/internals/models"
)
 
func TestAddDislike_ValidDislike_NoError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	err := AddDislike(db, models.Dislike{
		DislikeID: "dislike-1",
		Counts:    1,
		UserID:    user.UserID,
		PostID:    postID,
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
 
func TestHasUserDisliked_ReturnsFalse_WhenNoDislike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	disliked, err := HasUserDisliked(db, "user-1", "post-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if disliked {
		t.Error("expected false, got true")
	}
}
 
func TestHasUserDisliked_ReturnsTrue_AfterDislike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	AddDislike(db, models.Dislike{DislikeID: "dislike-1", Counts: 1, UserID: user.UserID, PostID: postID})
 
	disliked, err := HasUserDisliked(db, user.UserID, postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !disliked {
		t.Error("expected true after adding dislike")
	}
}
 
func TestRemoveDislike_RemovesDislike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	AddDislike(db, models.Dislike{DislikeID: "dislike-1", Counts: 1, UserID: user.UserID, PostID: postID})
 
	err := RemoveDislike(db, user.UserID, postID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
 
	disliked, _ := HasUserDisliked(db, user.UserID, postID)
	if disliked {
		t.Error("expected false after removing dislike")
	}
}