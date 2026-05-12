package db


import (
	"testing"
)

func TestInitDB_ValidPath(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if db == nil {
		t.Fatal("expected db to not be nil")
	}
	defer db.Close()
}

func TestInitTables_AllTablesCreated(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()
 
	expectedTables := []string{
		"users",
		"posts",
		"comments",
		"likes",
		"dislikes",
		"comment_likes",
		"comment_dislikes",
		"sessions",
	}
 
	for _, table := range expectedTables {
		row := db.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?",
			table,
		)
		var name string
		if err := row.Scan(&name); err != nil {
			t.Errorf("table %q was not created: %v", table, err)
		}
	}
}

func TestInitTables_Idempotent(t *testing.T) {
	// calling InitTables twice should not error (IF NOT EXISTS)
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()
 
	err = InitTables(db)
	if err != nil {
		t.Fatalf("second call to InitTables failed: %v", err)
	}
}

func TestInitTables_ForeignKeyEnforcement(t *testing.T) {
	db, err := InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}
	defer db.Close()
 
	// enable FK enforcement (SQLite has it off by default)
	db.Exec("PRAGMA foreign_keys = ON")
 
	// inserting a post with a non-existent userId should fail
	_, err = db.Exec(`
		INSERT INTO posts (postId, userId, title, content, category, createdAt)
		VALUES ('p1', 'nonexistent-user', 'Title', 'Content', 'Tech', '2024-01-01')
	`)
	if err == nil {
		t.Fatal("expected foreign key violation, got nil")
	}
}