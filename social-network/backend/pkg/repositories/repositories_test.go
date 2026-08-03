package repositories

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"social-network/backend/pkg/models"

	_ "github.com/mattn/go-sqlite3" // Ensure sqlite3 driver is available
)

// Helper function to initialize an in-memory SQLite database and create all tables
func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create all schemas used by the various repositories
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			first_name TEXT DEFAULT '',
			last_name TEXT DEFAULT '',
			bio TEXT DEFAULT '',
			avatar_url TEXT DEFAULT '',
			nickname TEXT DEFAULT '',
			date_of_birth TEXT DEFAULT '',
			is_private INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS followers (
			follower_id TEXT NOT NULL,
			following_id TEXT NOT NULL,
			created_at INTEGER,
			PRIMARY KEY (follower_id, following_id)
		);`,
		`CREATE TABLE IF NOT EXISTS follow_requests (
			id TEXT PRIMARY KEY,
			requester_id TEXT NOT NULL,
			target_id TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			created_at INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			type TEXT,
			source_id TEXT,
			ref_id TEXT,
			content TEXT,
			is_read INTEGER DEFAULT 0,
			created_at INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS groups (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			creator_id TEXT NOT NULL,
			created_at INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS group_members (
			group_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			joined_at INTEGER,
			PRIMARY KEY (group_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS group_messages (
			id TEXT PRIMARY KEY,
			group_id TEXT NOT NULL,
			sender_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS conversations (
			id TEXT PRIMARY KEY,
			created_at INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS conversation_participants (
			conversation_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			PRIMARY KEY (conversation_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			conversation_id TEXT NOT NULL,
			sender_id TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at INTEGER,
			is_read INTEGER DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS posts (
			id TEXT PRIMARY KEY,
			author_id TEXT NOT NULL,
			content TEXT,
			privacy TEXT NOT NULL,
			created_at INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS post_likes (
			post_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			PRIMARY KEY (post_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			expires_at TIMESTAMP NOT NULL
		);`,
	}

	for _, query := range schemas {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("failed to build schema: %v", err)
		}
	}

	return db
}

// Helper to seed a basic user into the database with default string values to bypass NULL scan exceptions
func seedUser(t *testing.T, db *sql.DB, username string, isPrivate bool) string {
	privateVal := 0
	if isPrivate {
		privateVal = 1
	}
	res, err := db.Exec(`
		INSERT INTO users (email, username, password_hash, first_name, last_name, bio, avatar_url, nickname, date_of_birth, is_private)
		VALUES (?, ?, 'hash', 'First', 'Last', '', '', '', '', ?)`,
		username+"@example.com", username, privateVal,
	)
	if err != nil {
		t.Fatalf("failed to seed user %s: %v", username, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("failed to fetch last insert id: %v", err)
	}
	return fmt.Sprintf("%d", id)
}

// ==========================================
// 1. TESTS FOR USERS REPOSITORY (users.go)
// ==========================================

func TestCreateUserAndGet(t *testing.T) {
	db := newTestDB(t)

	user := &models.User{
		Email:        "test@test.com",
		Username:     "testuser",
		PasswordHash: "secret",
		FirstName:    "John",
		LastName:     "Doe",
		IsPrivate:    false,
	}

	err := CreateUser(db, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if user.ID == "" {
		t.Error("expected user ID to be set after creation")
	}

	fetched, err := GetUserByEmail(db, "test@test.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}

	if fetched.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, fetched.Username)
	}
}

// ==========================================
// 2. TESTS FOR FOLLOW REPOSITORY (follow.go)
// ==========================================

func TestHandleFollowRequest_SelfFollow(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	alice := seedUser(t, db, "alice", false)

	_, _, err := repo.HandleFollowRequest("req1", "notif1", alice, alice)
	if err == nil {
		t.Fatal("expected error when following oneself, got nil")
	}
}

func TestHandleFollowRequest_PublicTarget(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)
	alice := seedUser(t, db, "alice", false)
	bob := seedUser(t, db, "bob", false)

	result, _, err := repo.HandleFollowRequest("req1", "notif1", bob, alice)
	if err != nil {
		t.Fatalf("HandleFollowRequest failed: %v", err)
	}

	if !result.Following || result.Status != "following" {
		t.Errorf("expected status 'following', got: %+v", result)
	}
}

// ==========================================
// 3. TESTS FOR SESSIONS REPOSITORY (session_repo.go)
// ==========================================

func TestSessionLifecycle(t *testing.T) {
	db := newTestDB(t)
	repo, err := NewSQLiteSessionRepo(db)
	if err != nil {
		t.Fatalf("NewSQLiteSessionRepo failed: %v", err)
	}

	sessionID := "sess-123"
	userID := "1"
	expiry := time.Now().Add(24 * time.Hour)

	err = repo.Create(sessionID, userID, expiry)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	foundUID, err := repo.LookupUserID(sessionID)
	if err != nil {
		t.Fatalf("failed to lookup session: %v", err)
	}

	if foundUID != userID {
		t.Errorf("expected userID %s, got %s", userID, foundUID)
	}

	err = repo.Delete(sessionID)
	if err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}
}

// ==========================================
// 4. TESTS FOR MESSAGES REPOSITORY (message_repo.go)
// ==========================================

func TestConversationsAndMessages(t *testing.T) {
	db := newTestDB(t)
	SetDatabase(db)

	u1 := seedUser(t, db, "user1", false)
	u2 := seedUser(t, db, "user2", false)

	convID, err := GetOrCreateConversation(u1, u2)
	if err != nil {
		t.Fatalf("failed to get/create conversation: %v", err)
	}

	msg, err := SaveMessage(convID, u1, "Hello user 2!")
	if err != nil {
		t.Fatalf("failed to save message: %v", err)
	}

	if msg.Content != "Hello user 2!" {
		t.Errorf("unexpected content: %s", msg.Content)
	}

	msgs, err := GetMessages(convID, 10, 0)
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}

	if len(msgs) == 0 {
		t.Fatal("expected messages to be found")
	}
}

// ==========================================
// 5. TESTS FOR GROUPS REPOSITORY (group_repo.go)
// ==========================================

func TestGroupCreationAndManagement(t *testing.T) {
	db := newTestDB(t)
	u1 := seedUser(t, db, "creator", false)

	req := models.CreateGroupRequest{
		Title:       "Golang Experts",
		Description: "A group for Go developers",
	}

	g, err := CreateGroup(db, req, u1)
	if err != nil {
		t.Fatalf("failed to create group: %v", err)
	}

	if g.Title != "Golang Experts" {
		t.Errorf("expected title 'Golang Experts', got '%s'", g.Title)
	}
}

// ==========================================
// 6. TESTS FOR GROUP MESSAGES (group_message_repo.go)
// ==========================================

func TestGroupMessages(t *testing.T) {
	db := newTestDB(t)
	DB = db

	u1 := seedUser(t, db, "member1", false)
	groupID := "grp-abc"

	_, err := db.Exec(`INSERT INTO groups (id, title, creator_id, created_at) VALUES (?, 'Test Group', ?, ?)`, groupID, u1, time.Now().Unix())
	if err != nil {
		t.Fatalf("setup group failed: %v", err)
	}

	msg, err := SaveGroupMessage(groupID, u1, "Hello Group!")
	if err != nil {
		t.Fatalf("SaveGroupMessage failed: %v", err)
	}

	if msg.Content != "Hello Group!" {
		t.Errorf("expected 'Hello Group!', got '%s'", msg.Content)
	}
}

// ==========================================
// 7. TESTS FOR PROFILE REPOSITORY (profile_repo.go)
// ==========================================

func TestUpdateUserProfile(t *testing.T) {
	db := newTestDB(t)
	u1Str := seedUser(t, db, "profileuser", false)
	var u1 int
	fmt.Sscanf(u1Str, "%d", &u1)

	newBio := "Software Engineer & Architect"
	newNickname := "GoGuru"

	params := UpdateUserProfileParams{
		Bio:      &newBio,
		Nickname: &newNickname,
	}

	updatedUser, err := UpdateUserProfile(db, u1, params)
	if err != nil {
		t.Fatalf("UpdateUserProfile failed: %v", err)
	}

	if updatedUser.Bio != newBio || updatedUser.Nickname != newNickname {
		t.Errorf("profile parameters did not update properly: %+v", updatedUser)
	}
}