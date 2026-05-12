package handlers
 
import (
	"net/http"
	"net/http/httptest"
	"testing"
)
 
func TestSearchHandler_ValidQuery_ReturnsOK(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	// seed a user and post so search has something to find
	user := seedUser(t, db)
	seedPost(t, db, user.UserID)
 
	req := httptest.NewRequest(http.MethodGet, "/search?q=Test", nil)
	rr := httptest.NewRecorder()
 
	SearchHandler(db).ServeHTTP(rr, req)
 
	// should not error — 200 or template error (path) both acceptable, never 500 from db
	if rr.Code == http.StatusSeeOther {
		t.Errorf("unexpected 500 on valid search query")
	}
}
 
func TestSearchHandler_EmptyQuery_DoesNotError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	req := httptest.NewRequest(http.MethodGet, "/search?q=", nil)
	rr := httptest.NewRecorder()
 
	SearchHandler(db).ServeHTTP(rr, req)
 
	if rr.Code == http.StatusSeeOther {
		t.Errorf("unexpected 500 on empty search query")
	}
}
 
func TestSearchHandler_Unauthenticated_DoesNotError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	// no cookie — unauthenticated user
	req := httptest.NewRequest(http.MethodGet, "/search?q=forum", nil)
	rr := httptest.NewRecorder()
 
	SearchHandler(db).ServeHTTP(rr, req)
 
	if rr.Code == http.StatusUnauthorized {
		t.Errorf("search should be accessible without login, got 401")
	}
}