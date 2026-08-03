package sqlite

import (
	"testing"
)

func TestSqliteDB(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite DB: %v", err)
	}
	defer db.Close()

	if db.GetConn() == nil {
		t.Error("expected non-nil database connection")
	}

	// Test Exec
	_, err = db.Exec("CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("failed to execute CREATE TABLE query: %v", err)
	}

	_, err = db.Exec("INSERT INTO test (name) VALUES (?)", "foo")
	if err != nil {
		t.Fatalf("failed to insert data: %v", err)
	}

	// Test QueryRow
	var name string
	err = db.QueryRow("SELECT name FROM test WHERE id = ?", 1).Scan(&name)
	if err != nil {
		t.Fatalf("failed to query row: %v", err)
	}
	if name != "foo" {
		t.Errorf("expected name 'foo', got %q", name)
	}

	// Test Query
	rows, err := db.Query("SELECT name FROM test")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatalf("failed to scan row: %v", err)
		}
		names = append(names, n)
	}
	if len(names) != 1 || names[0] != "foo" {
		t.Errorf("expected names ['foo'], got %v", names)
	}

	// Test BeginTx
	tx, err := db.BeginTx()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}
	_, err = tx.Exec("INSERT INTO test (name) VALUES (?)", "bar")
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to execute in transaction: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("failed to commit transaction: %v", err)
	}

	// Check if bar was inserted
	err = db.QueryRow("SELECT name FROM test WHERE id = ?", 2).Scan(&name)
	if err != nil {
		t.Fatalf("failed to query row after transaction: %v", err)
	}
	if name != "bar" {
		t.Errorf("expected name 'bar', got %q", name)
	}
}
