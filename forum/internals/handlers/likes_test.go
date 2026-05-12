package handlers
 
import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
 
	"forum/internals/models"
	"forum/internals/services"
	"forum/internals/utils"
)
 
// -----------------------------------------------------------------------
// LikeHandler tests
// -----------------------------------------------------------------------
 
func TestLikeHandler_ValidLike_Redirects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	form := url.Values{}
	form.Set("post_id", postID)
	form.Set("user_id", user.UserID)
 
	req := httptest.NewRequest(http.MethodPost, "/vote/like", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	LikeHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rr.Code)
	}
}
 
func TestLikeHandler_ToggleLike_RemovesLike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	// like once
	services.AddLike(db, models.Like{
		LikeID: utils.GenerateUUID(),
		PostID: postID,
		UserID: user.UserID,
		Counts: 1,
	})
 
	// like again — should toggle off
	form := url.Values{}
	form.Set("post_id", postID)
	form.Set("user_id", user.UserID)
 
	req := httptest.NewRequest(http.MethodPost, "/vote/like", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	LikeHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 after toggle, got %d", rr.Code)
	}
 
	hasLiked, _ := services.HasUserLiked(db, user.UserID, postID)
	if hasLiked {
		t.Error("expected like to be removed after toggle")
	}
}
 
func TestLikeHandler_RemovesDislike_BeforeAddsLike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	// dislike first
	services.AddDislike(db, models.Dislike{
		DislikeID: utils.GenerateUUID(),
		PostID:    postID,
		UserID:    user.UserID,
		Counts:    1,
	})
 
	// now like
	form := url.Values{}
	form.Set("post_id", postID)
	form.Set("user_id", user.UserID)
 
	req := httptest.NewRequest(http.MethodPost, "/vote/like", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	LikeHandler(db).ServeHTTP(rr, req)
 
	hasDisliked, _ := services.HasUserDisliked(db, user.UserID, postID)
	if hasDisliked {
		t.Error("expected dislike to be removed after liking")
	}
}
 
// -----------------------------------------------------------------------
// DislikeHandler tests
// -----------------------------------------------------------------------
 
func TestDislikeHandler_ValidDislike_Redirects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	form := url.Values{}
	form.Set("post_id", postID)
	form.Set("user_id", user.UserID)
 
	req := httptest.NewRequest(http.MethodPost, "/vote/dislike", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	DislikeHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303, got %d", rr.Code)
	}
}
 
func TestDislikeHandler_RemovesLike_BeforeAddsDislike(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	postID := seedPost(t, db, user.UserID)
 
	// like first
	services.AddLike(db, models.Like{
		LikeID: utils.GenerateUUID(),
		PostID: postID,
		UserID: user.UserID,
		Counts: 1,
	})
 
	// now dislike
	form := url.Values{}
	form.Set("post_id", postID)
	form.Set("user_id", user.UserID)
 
	req := httptest.NewRequest(http.MethodPost, "/vote/dislike", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	DislikeHandler(db).ServeHTTP(rr, req)
 
	hasLiked, _ := services.HasUserLiked(db, user.UserID, postID)
	if hasLiked {
		t.Error("expected like to be removed after disliking")
	}
}
 
// -----------------------------------------------------------------------
// MyLikesHandler tests
// -----------------------------------------------------------------------
 
func TestMyLikesHandler_Unauthenticated_RedirectsToLogin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodGet, "/my-likes", nil)
	rr := httptest.NewRecorder()
 
	MyLikesHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect to login, got %d", rr.Code)
	}
	if rr.Header().Get("Location") != "/login" {
		t.Errorf("expected redirect to /login, got %s", rr.Header().Get("Location"))
	}
}
 
func TestMyLikesHandler_InvalidMethod_ReturnsMethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodPost, "/my-likes", nil)
	rr := httptest.NewRecorder()
 
	MyLikesHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
 
// -----------------------------------------------------------------------
// MyPostsHandler tests
// -----------------------------------------------------------------------
 
func TestMyPostsHandler_Unauthenticated_RedirectsToLogin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodGet, "/my-posts", nil)
	rr := httptest.NewRecorder()
 
	MyPostsHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect to login, got %d", rr.Code)
	}
	if rr.Header().Get("Location") != "/login" {
		t.Errorf("expected redirect to /login, got %s", rr.Header().Get("Location"))
	}
}
 
func TestMyPostsHandler_InvalidMethod_ReturnsMethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodPost, "/my-posts", nil)
	rr := httptest.NewRecorder()
 
	MyPostsHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}