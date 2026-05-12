package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Initialize tables
	err = InitTables(db)
	if err != nil {
		return db, err
	}

	return db, nil
}

func InitTables(db *sql.DB) error {
	// Sql Scripts to create tables
	createUsers := `CREATE TABLE IF NOT EXISTS users(
		userId TEXT PRIMARY KEY NOT NULL, email TEXT NOT NULL UNIQUE, 
		username TEXT NOT NULL, password TEXT NOT NULL, profile TEXT
	)`

	createComments := `CREATE TABLE IF NOT EXISTS comments(
		commentId TEXT PRIMARY KEY NOT NULL,
		userId TEXT NOT NULL,
		postId TEXT NOT NULL,
		content TEXT NOT NULL,
		createdAt TEXT NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE,
		FOREIGN KEY (postId) REFERENCES posts(postId) ON DELETE CASCADE
	)`

	createPosts := `CREATE TABLE IF NOT EXISTS posts(
		postId TEXT PRIMARY KEY NOT NULL,
		userId TEXT NOT NULL,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		category TEXT NOT NULL,
		createdAt TEXT NOT NULL,
		image TEXT,
		FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE 	
	)`

	createLikes := `CREATE TABLE IF NOT EXISTS likes (
		likeId TEXT PRIMARY KEY NOT NULL,
		counts INTEGER NOT NULL,
		userId TEXT NOT NULL,
		postId TEXT NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE,
		FOREIGN KEY (postId) REFERENCES posts(postId) ON DELETE CASCADE
	)`

	createDislikes := `CREATE TABLE IF NOT EXISTS dislikes (
		dislikeId TEXT PRIMARY KEY NOT NULL,
		counts INTEGER NOT NULL,
		userId TEXT NOT NULL,
		postId TEXT NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE,
		FOREIGN KEY (postId) REFERENCES posts(postId) ON DELETE CASCADE
)`

	createCommentLikes := `CREATE TABLE IF NOT EXISTS comment_likes (
		likeId TEXT PRIMARY KEY NOT NULL,
		counts INTEGER NOT NULL,
		userId TEXT NOT NULL,
		commentId TEXT NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE,
		FOREIGN KEY (commentId) REFERENCES comments(commentId) ON DELETE CASCADE
	)`

	createCommentDislikes := `CREATE TABLE IF NOT EXISTS comment_dislikes (
		dislikeId TEXT PRIMARY KEY NOT NULL,
		counts INTEGER NOT NULL,
		userId TEXT NOT NULL,
		commentId TEXT NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE,
		FOREIGN KEY (commentId) REFERENCES comments(commentId) ON DELETE CASCADE
	)`

	createSessions := `CREATE TABLE IF NOT EXISTS sessions (
		sessionId TEXT PRIMARY KEY NOT NULL,
		userId TEXT NOT NULL,
		expiresAt TEXT NOT NULL,
		FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE
	)`

	// Prepare statements for execution
	users, err := db.Prepare(createUsers)
	if err != nil {
		return err
	}

	comments, err := db.Prepare(createComments)
	if err != nil {
		return err
	}
	posts, err := db.Prepare(createPosts)
	if err != nil {
		return err
	}

	likes, err := db.Prepare(createLikes)
	if err != nil {
		return nil
	}

	dislikes, err := db.Prepare(createDislikes)
	if err != nil {
		return nil
	}

	commentLikes, err := db.Prepare(createCommentLikes)
	if err != nil {
		return nil
	}

	commentDislikes, err := db.Prepare(createCommentDislikes)
	if err != nil {
		return nil
	}

	sessions, err := db.Prepare(createSessions)
	if err != nil {
		return err
	}

	// Close the statements after completion
	defer users.Close()
	defer comments.Close()
	defer posts.Close()
	defer likes.Close()
	defer dislikes.Close()
	defer commentLikes.Close()
	defer commentDislikes.Close()
	defer sessions.Close()

	// Execute queries
	_, err = users.Exec()
	if err != nil {
		return err
	}

	_, err = comments.Exec()
	if err != nil {
		return err
	}

	_, err = posts.Exec()
	if err != nil {
		return err
	}

	_, err = likes.Exec()
	if err != nil {
		return nil
	}

	_, err = dislikes.Exec()
	if err != nil {
		return nil
	}

	_, err = commentLikes.Exec()
	if err != nil {
		return nil
	}

	_, err = commentDislikes.Exec()
	if err != nil {
		return nil
	}

	_, err = sessions.Exec()
	if err != nil {
		return err
	}

	// clear sessions
	//_, err = db.Exec("DELETE FROM sqlite_sequence WHERE table='sessions'")
	//if err != nil {
	//	fmt.Println(err)
	//}

	return nil
}
