package services

import (
	"database/sql"
	"errors"
	"strings"

	"forum/internals/models"
)

// Comment start
func AddComment(db *sql.DB, comment models.Comment) error {
	addComment, err := db.Prepare("INSERT INTO COMMENTS(commentId,userId,postId,content,createdAt) VALUES(?,?,?,?,?)")
	if err != nil {
		return errors.New("Error adding comment: " + err.Error())
	}
	defer addComment.Close()
	_, err = addComment.Exec(comment.CommentID, comment.UserID, comment.PostID, comment.Content, comment.CreatedAt)
	if err != nil {
		return errors.New("Error creating table to add comment: " + err.Error())
	}
	// fmt.Println(result.LastInsertId())
	return nil
}

func ViewComment(db *sql.DB) ([]models.Post, error) {
	viewComment, err := db.Prepare(`SELECT 
		p.postId,
		p.userId,
		p.title,
		p.content,
		p.category,
		p.createdAt,
		p.image,
		u.username,
	COUNT(c.commentId) AS comment_count
	FROM posts p
	LEFT JOIN comments c ON p.postId = c.postId
	JOIN users u ON p.userId = u.userId
	GROUP BY 
		p.postId, p.userId, p.title, p.content, 
		p.category, p.createdAt, p.image, u.username
	ORDER BY comment_count DESC;`)
	if err != nil {

		// fmt.Println("Error editing comment: ", err)
		return nil, err
	}
	defer viewComment.Close()
	comments, err := viewComment.Query()
	if err != nil {
		// fmt.Println("Error retrieving users: ", err)
		return nil, err
	}

	posts := []models.Post{}
	for comments.Next() {
		var p models.Post
		var commentCount int

		err := comments.Scan(&p.PostID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.Category,
			&p.CreatedAt,
			&p.Image,
			&p.Username,
			&commentCount)
		if err != nil {
			return nil, err
		}
		// capitalize username for display
		if p.Username != "" {
			p.Username = strings.ToUpper(string(p.Username[0])) + p.Username[1:]
			p.Initials = string(p.Username[0])
		}
		// populate likes/dislikes so template shows them
		db.QueryRow(`SELECT COUNT(*) FROM likes WHERE postId = ?`, p.PostID).Scan(&p.Likes)
		db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE postId = ?`, p.PostID).Scan(&p.Dislikes)

		posts = append(posts, p)
	}
	return posts, nil
}

// GetCommentsByPost returns all comments for a given post including commenter username
func GetCommentsByPost(db *sql.DB, postID string) ([]models.Comment, error) {
	query, err := db.Prepare(`SELECT c.commentId, c.userId, c.postId, c.content, c.createdAt, u.username
		FROM comments c
		JOIN users u ON c.userId = u.userId
		WHERE c.postId = ?
		ORDER BY c.createdAt ASC`)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	rows, err := query.Query(postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.CommentID, &c.UserID, &c.PostID, &c.Content, &c.CreatedAt, &c.Username)
		if err != nil {
			return nil, err
		}
		// populate like/dislike counts for the comment
		db.QueryRow(`SELECT COUNT(*) FROM comment_likes WHERE commentId = ?`, c.CommentID).Scan(&c.Likes)
		db.QueryRow(`SELECT COUNT(*) FROM comment_dislikes WHERE commentId = ?`, c.CommentID).Scan(&c.Dislikes)

		comments = append(comments, c)
	}
	return comments, nil
}

func EditComment(db *sql.DB, comment models.Comment) {
	editComment, err := db.Prepare("UPDATE COMMENT SET commentId=?, userId=?, postId=?, content=?, createdAt=? WHERE commentId=?")
	if err != nil {
		// fmt.Println("Error editing comment: ", err)
		return
	}
	defer editComment.Close()
	_, err = editComment.Exec(comment.CommentID, comment.UserID, comment.PostID, comment.Content, comment.CreatedAt)
	if err != nil {
		// fmt.Println("Error editing comment: ", err)
	}
	// fmt.Println(result.LastInsertId())
}

// Comment likes/dislikes services
func AddCommentLike(db *sql.DB, like models.CommentLike) error {
	stmt, err := db.Prepare("INSERT INTO comment_likes (likeId, counts, userId, commentId) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(like.LikeID, like.Counts, like.UserID, like.CommentID)
	return err
}

func RemoveCommentLike(db *sql.DB, userID string, commentID string) error {
	_, err := db.Exec("DELETE FROM comment_likes WHERE userId = ? AND commentId = ?", userID, commentID)
	return err
}

func HasUserLikedComment(db *sql.DB, userID string, commentID string) (bool, error) {
	var id string
	err := db.QueryRow("SELECT likeId FROM comment_likes WHERE userId = ? AND commentId = ?", userID, commentID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetCommentLikeCount(db *sql.DB, commentID string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM comment_likes WHERE commentId = ?", commentID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AddCommentDislike(db *sql.DB, dis models.CommentDislike) error {
	stmt, err := db.Prepare("INSERT INTO comment_dislikes (dislikeId, counts, userId, commentId) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(dis.DislikeID, dis.Counts, dis.UserID, dis.CommentID)
	return err
}

func RemoveCommentDislike(db *sql.DB, userID string, commentID string) error {
	_, err := db.Exec("DELETE FROM comment_dislikes WHERE userId = ? AND commentId = ?", userID, commentID)
	return err
}

func HasUserDislikedComment(db *sql.DB, userID string, commentID string) (bool, error) {
	var id string
	err := db.QueryRow("SELECT dislikeId FROM comment_dislikes WHERE userId = ? AND commentId = ?", userID, commentID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetCommentDislikeCount(db *sql.DB, commentID string) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM comment_dislikes WHERE commentId = ?", commentID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Comment End
