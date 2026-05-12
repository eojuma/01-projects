package services
 
import (
	"testing"
 
	"forum/internals/models"
)
 
func TestAddUser_ValidUser_NoError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	err := AddUser(db, models.User{
		UserID: "u1", Email: "a@b.com",
		Username: "user", Password: "pass",
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
 
func TestAddUser_DuplicateEmail_ReturnsError(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	AddUser(db, models.User{UserID: "u1", Email: "a@b.com", Username: "user", Password: "pass"})
	err := AddUser(db, models.User{UserID: "u2", Email: "a@b.com", Username: "user2", Password: "pass"})
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
}
 
func TestGetUser_ReturnsUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	AddUser(db, models.User{UserID: "u1", Email: "a@b.com", Username: "user", Password: "pass"})
 
	user, err := GetUser(db, "a@b.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "a@b.com" {
		t.Errorf("expected email a@b.com, got %s", user.Email)
	}
}
 
func TestGetUser_NonexistentEmail_ReturnsEmptyUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
 
	user, err := GetUser(db, "ghost@nowhere.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.UserID != "" {
		t.Errorf("expected empty user, got userID: %s", user.UserID)
	}
}