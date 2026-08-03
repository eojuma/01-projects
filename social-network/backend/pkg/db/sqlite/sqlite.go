package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)
// Setting up SQLite database connection and initialization.
// DB wraps the database connection in a reusable struct
type DB struct {
	conn *sql.DB
}

// New creates a new SQLite database connection
func New(dbPath string) (*DB, error) {
	if dbPath == "" {
		dbPath = "social_network.db"
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)

	return &DB{conn: conn}, nil
}

// GetConn returns the underlying database connection
func (db *DB) GetConn() *sql.DB {
	return db.conn
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Exec executes a query that doesn't return rows
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}

// Query executes a query that returns rows
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

// QueryRow executes a query that returns a single row
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRow(query, args...)
}

// BeginTx starts a new transaction
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.conn.Begin()
}
