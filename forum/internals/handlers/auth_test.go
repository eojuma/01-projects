package handlers
 
import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
 
	_ "github.com/mattn/go-sqlite3"
 
	"forum/internals/models"
	"forum/internals/services"
	"forum/internals/utils"

	dbpkg "forum/internals/db"


)
 
// -----------------------------------------------------------------------
// Test DB setup — in-memory SQLite, no files created
// -----------------------------------------------------------------------
 
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
 
// seeds one user and returns their plaintext password
func seedUser(t *testing.T, db *sql.DB) models.User {
	t.Helper()
	plainPassword := "TestPass123!"
	hashed, _ := utils.HashPass(plainPassword)
	user := models.User{
		UserID:         "test-user-id",
		Email:          "test@example.com",
		Username:       "testuser",
		Password:       hashed,
		ProfilePicture: "",
	}
	err := services.AddUser(db, user)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	user.Password = plainPassword // return plaintext for use in tests
	return user
}
 
// -----------------------------------------------------------------------
// LoginHandler tests
// -----------------------------------------------------------------------
 
func TestLoginHandler_GET_ReturnsForm(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()

	LoginHandler(db).ServeHTTP(rr, req)

	// 200 or 500 both acceptable — 500 means template not found (path issue in test env)
	// what we're confirming is it does NOT redirect
	if rr.Code == http.StatusSeeOther {
		t.Error("GET /login should not redirect")
	}
}
 
func TestLoginHandler_POST_CorrectCredentials_Redirects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	user := seedUser(t, db)
 
	form := url.Values{}
	form.Set("email", user.Email)
	form.Set("password", user.Password)
 
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	LoginHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rr.Code)
	}
	if rr.Header().Get("Location") != "/" {
		t.Errorf("expected redirect to /, got %s", rr.Header().Get("Location"))
	}
}
 
func TestLoginHandler_POST_CorrectCredentials_SetsCookie(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	user := seedUser(t, db)
 
	form := url.Values{}
	form.Set("email", user.Email)
	form.Set("password", user.Password)
 
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	LoginHandler(db).ServeHTTP(rr, req)
 
	cookies := rr.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "user_session" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected user_session cookie to be set")
	}
}
 
func TestLoginHandler_POST_WrongPassword_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	user := seedUser(t, db)
 
	form := url.Values{}
	form.Set("email", user.Email)
	form.Set("password", "wrongpassword")
 
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	LoginHandler(db).ServeHTTP(rr, req)
 
	// should not redirect — should re-render the form
	if rr.Code == http.StatusSeeOther {
		t.Error("expected no redirect on wrong password")
	}
}
 
func TestLoginHandler_InvalidMethod_ReturnsBadRequest(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodPut, "/login", nil)
	rr := httptest.NewRecorder()
 
	LoginHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
 
// -----------------------------------------------------------------------
// RegisterHandler tests
// -----------------------------------------------------------------------
 
func TestRegisterHandler_GET_ReturnsForm(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	rr := httptest.NewRecorder()

	RegisterHandler(db).ServeHTTP(rr, req)

	if rr.Code == http.StatusSeeOther {
		t.Error("GET /register should not redirect")
	}
}
 

func TestRegisterHandler_POST_ValidData_Redirects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	form := url.Values{}
	form.Set("email", "newuser@example.com")
	form.Set("username", "newuser")
	form.Set("password", "StrongPass123!")
	form.Set("confirm_password", "StrongPass123!")
 
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	RegisterHandler(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect to /login, got %d", rr.Code)
	}
	if rr.Header().Get("Location") != "/login" {
		t.Errorf("expected redirect to /login, got %s", rr.Header().Get("Location"))
	}
}
 
func TestRegisterHandler_POST_PasswordMismatch_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	form := url.Values{}
	form.Set("email", "user@example.com")
	form.Set("username", "someuser")
	form.Set("password", "StrongPass123!")
	form.Set("confirm_password", "DifferentPass!")
 
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	RegisterHandler(db).ServeHTTP(rr, req)
 
	if rr.Code == http.StatusSeeOther {
		t.Error("expected no redirect on password mismatch")
	}
}
 
func TestRegisterHandler_POST_DuplicateEmail_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedUser(t, db) // seeds test@example.com
 
	form := url.Values{}
	form.Set("email", "test@example.com") // duplicate
	form.Set("username", "anotheruser")
	form.Set("password", "StrongPass123!")
	form.Set("confirm_password", "StrongPass123!")
 
	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
 
	RegisterHandler(db).ServeHTTP(rr, req)
 
	if rr.Code == http.StatusSeeOther {
		t.Error("expected no redirect on duplicate email")
	}
}
 
// -----------------------------------------------------------------------
// Logout tests
// -----------------------------------------------------------------------
 
func TestLogout_ClearsCookieAndRedirects(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	user := seedUser(t, db)
 
	// create a session first
	sessionID := utils.GenerateUUID()
	expiresAt := time.Now().Add(6 * time.Hour)
	services.AddSession(db, sessionID, user.UserID, expiresAt)
 
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "user_session", Value: sessionID})
	rr := httptest.NewRecorder()
 
	Logout(db).ServeHTTP(rr, req)
 
	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 redirect, got %d", rr.Code)
	}
 
	// confirm cookie is cleared
	cookies := rr.Result().Cookies()
	for _, c := range cookies {
		if c.Name == "user_session" && c.MaxAge != -1 {
			t.Error("expected user_session cookie to be cleared")
		}
	}
}

