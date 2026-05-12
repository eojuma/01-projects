package services

import (
	"database/sql"
	"errors"
	"strings"

	"forum/internals/models"
)

/*
func CreatePost(db *sql.DB, post models.Post) error {

	query := `
	INSERT INTO posts (user_id,title,contents,category)
	VALUES (?,?,?,?)

	`

	_, err := db.Exec(query, post.UserID, post.Title, post.Content, post.Category)

	if err != nil {
		return err
	}
	return nil
}

func GetPosts(db *sql.DB) ([]models.Post, error) {
	query := `
	SELECT P.id,p.user_id,p.title,p.contents,p.category,p.created_at
	FROM posts p
	JOIN users u ON p.user_id = u.id
	ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.PostID, &p.UserID, &p.Title, &p.Content, &p.Category, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
*/

// Post start
func AddPost(db *sql.DB, post models.Post) error {
	// prepare add post sql statement with placeholders
	addPost, err := db.Prepare("INSERT INTO POSTS(postId, userid, title, content, category, createdAt, image) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return errors.New("Error adding post: " + err.Error())
	}

	// close after user
	defer addPost.Close()

	// execute sql query
	_, err = addPost.Exec(post.PostID, post.UserID, post.Title, post.Content, post.Category, post.CreatedAt, post.Image)
	if err != nil {
		return errors.New("Error creating table to add post: " + err.Error())
	}
	return nil
}

func ViewPost(db *sql.DB) ([]models.Post, error) {
	// prepare sql query
	viewPost, err := db.Prepare("SELECT posts.*,users.username  FROM POSTS JOIN users ON posts.userId=users.userId ORDER BY createdAt DESC")
	if err != nil {
		return nil, err
	}
	// close on complete
	defer viewPost.Close()

	// execute query all posts
	posts, err := viewPost.Query()
	if err != nil {
		return nil, err
	}

	// slice to hold all posts
	var Allposts []models.Post

	// Read results
	for posts.Next() {
		// initialize post struct
		var p models.Post
		// write values to struct
		err := posts.Scan(&p.PostID, &p.UserID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.Image, &p.Username)
		if err != nil {
			return nil, err
		}
		// Capitalize first letter and set Initials
		if len(p.Username) > 0 {
			p.Username = strings.ToUpper(string(p.Username[0])) + p.Username[1:]
			p.Initials = string(p.Username[0])
		}
		// append post to Allpost Slice
		Allposts = append(Allposts, p)
	}

	for i, p := range Allposts {
		db.QueryRow(`SELECT COUNT(*) FROM likes WHERE postId = ?`, p.PostID).Scan(&Allposts[i].Likes)
		db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE postId = ?`, p.PostID).Scan(&Allposts[i].Dislikes)
	}
	return Allposts, nil
}

func EditPost(db *sql.DB, post models.Post) error {
	// prepare edit query
	editPost, err := db.Prepare("UPDATE POSTS SET postId=?, userid=?, title=?, content=? , category=?, createdAt=?, image=? WHERE postid=?")
	if err != nil {
		return err
	}
	// close statement after use
	defer editPost.Close()

	// pass post data to update and execute update
	_, err = editPost.Exec(post.PostID, post.UserID, post.Title, post.Content, post.Category, post.CreatedAt, post.Image)
	if err != nil {
		return err
	}

	return nil
}

func GetPostById(db *sql.DB, postID string) (models.Post, error) {
	// prepare query to get post by id
	getPost, err := db.Prepare("SELECT p.*, u.username FROM posts p JOIN users u ON p.userId = u.userId WHERE p.postId = ?")
	if err != nil {
		return models.Post{}, err
	}
	defer getPost.Close()

	// execute query with postID as parameter
	row := getPost.QueryRow(postID)

	var p models.Post

	// scan post row including username
	err = row.Scan(&p.PostID, &p.UserID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.Image, &p.Username)
	if err != nil {
		return models.Post{}, err
	}

	if len(p.Username) > 0 {
		p.Username = strings.ToUpper(string(p.Username[0])) + p.Username[1:]
		p.Initials = string(p.Username[0])
	}

	// populate like and dislike counts after confirming post exists
	err = db.QueryRow(`SELECT COUNT(*) FROM likes WHERE postId = ?`, postID).Scan(&p.Likes)
	if err != nil {
		return models.Post{}, err
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE postId = ?`, postID).Scan(&p.Dislikes)
	if err != nil {
		return models.Post{}, err
	}

	return p, nil
}

func SearchPosts(db *sql.DB, query string) ([]models.Post, error) {
	q := "%" + query + "%"

	rows, err := db.Query(`
        SELECT p.postId, p.userId, u.username, p.title, p.content, p.createdAt
    FROM posts p
    JOIN users u ON p.userId = u.userId
    WHERE p.title LIKE ?
       OR p.content LIKE ?
       OR u.username LIKE ?
    ORDER BY p.createdAt DESC
    `, q, q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.PostID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		if len(p.Username) > 0 {
			p.Username = strings.ToUpper(string(p.Username[0])) + p.Username[1:]
			p.Initials = string(p.Username[0])
		}
		// populate like and dislike counts
		db.QueryRow(`SELECT COUNT(*) FROM likes WHERE postId = ?`, p.PostID).Scan(&p.Likes)
		db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE postId = ?`, p.PostID).Scan(&p.Dislikes)

		posts = append(posts, p)
	}

	return posts, nil
}

// Post End
