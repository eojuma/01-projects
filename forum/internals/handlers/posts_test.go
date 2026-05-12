package handlers
 
import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)
 
// -----------------------------------------------------------------------
// CreatePostHandler tests
// -----------------------------------------------------------------------
 
func TestCreatePostHandler_GET_Unauthenticated_RedirectsToLogin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodGet, "/post/create", nil)
	rr := httptest.NewRecorder()
 
	CreatePostHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rr.Code)
	}
	if rr.Header().Get("Location") != "/login" {
		t.Errorf("expected redirect to /login, got %s", rr.Header().Get("Location"))
	}
}
 
func TestCreatePostHandler_POST_ValidPost_Redirects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user := seedUser(t, db)
	cookie := sessionCookie(t, db, user.UserID)
 
	form := url.Values{}
	form.Set("title", "Test Title")
	form.Set("content", "Test content body")
	form.Add("category", "Tech")
 
	req := httptest.NewRequest(http.MethodPost, "/post/create", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(cookie)
	rr := httptest.NewRecorder()
 
	CreatePostHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rr.Code)
	}
}
 
func TestCreatePostHandler_InvalidMethod_ReturnsMethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodPut, "/post/create", nil)
	rr := httptest.NewRecorder()
 
	CreatePostHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
 
// -----------------------------------------------------------------------
// ViewPostHandler tests
// -----------------------------------------------------------------------
 
func TestViewPostHandler_InvalidMethod_ReturnsMethodNotAllowed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodPost, "/post/view", nil)
	rr := httptest.NewRecorder()
 
	ViewPostHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}
 
func TestViewPostHandler_MissingPostID_ReturnsBadRequest(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodGet, "/post/view", nil)
	rr := httptest.NewRecorder()
 
	ViewPostHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
 
// -----------------------------------------------------------------------
// FilterByCategory tests
// -----------------------------------------------------------------------
 
func TestFilterByCategory_Unauthenticated_DoesNotError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodGet, "/post/category?q=Tech", nil)
	rr := httptest.NewRecorder()
 
	FilterByCategory(db).ServeHTTP(rr, req)
 
	if rr.Code == http.StatusUnauthorized {
		t.Errorf("category filter should be accessible without login, got 401")
	}
}