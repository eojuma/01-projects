package services

import (
	"database/sql"
	"fmt"
	"strings"

	"forum/internals/models"
)

func AddLike(db *sql.DB, like models.Like) error {
	likee, err := db.Prepare("INSERT INTO likes (likeId, counts, userId, postId) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer likee.Close()

	_, err = likee.Exec(like.LikeID, like.Counts, like.UserID, like.PostID)
	if err != nil {
		return err
	}

	return nil
}

func RemoveLike(db *sql.DB, userID string, postID string) error {
	_, err := db.Exec("DELETE FROM likes WHERE userId = ? AND postId = ?", userID, postID)
	if err != nil {
		return err
	}

	return nil
}

func HasUserLiked(db *sql.DB, userID string, postID string) (bool, error) {
	var likeID string
	err := db.QueryRow("SELECT likeId FROM likes WHERE userId = ? AND postId = ?", userID, postID).Scan(&likeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetLikeCount(db *sql.DB, postID string) (int, error) {
	var count int
	err := db.QueryRow("SELECT SUM(*) FROM likes WHERE postId = ?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func ViewByLikes(db *sql.DB) ([]models.Post, error) {
	viewLikes, err := db.Prepare(`SELECT 
		p.postId,
		p.userId,
		p.title,
		p.content,
		p.category,
		p.createdAt,
		p.image,
		u.username,
	COUNT(l.likeId) AS like_count
	FROM posts p
	LEFT JOIN likes l ON p.postId = l.postId
	JOIN users u ON p.userId = u.userId
	GROUP BY 
		p.postId, p.userId, p.title, p.content, 
		p.category, p.createdAt, p.image, u.username
	ORDER BY like_count DESC;`)
	if err != nil {
		fmt.Println("Error viewing by likes: ", err)
		return nil, err
	}
	defer viewLikes.Close()
	likes, err := viewLikes.Query()
	if err != nil {
		fmt.Println("Error retrieving users: ", err)
		return nil, err
	}

	posts := []models.Post{}
	for likes.Next() {
		var p models.Post
		var likeCount int

		err := likes.Scan(&p.PostID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.Category,
			&p.CreatedAt,
			&p.Image,
			&p.Username,
			&likeCount)
		if err != nil {
			return nil, err
		}
		p.Likes = likeCount
		// capitalize username first letter for display consistency
		if p.Username != "" {
			p.Username = strings.ToUpper(string(p.Username[0])) + p.Username[1:]
			p.Initials = string(p.Username[0])
		}
		// populate dislikes as well
		db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE postId = ?`, p.PostID).Scan(&p.Dislikes)
		posts = append(posts, p)

	}
	return posts, nil
}

func GetMyLikes(db *sql.DB, userId string) ([]models.Post, error) {
	query, err := db.Prepare(`SELECT p.postId, p.userId, p.title, p.content, p.category, p.createdAt, p.image, u.username
	FROM posts p
	JOIN likes l ON p.postId = l.postId
	JOIN users u ON p.userId = u.userId
	WHERE l.userId = ?`)
	if err != nil {
		return nil, err
	}
	rows, err := query.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Posts []models.Post

	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.PostID, &p.UserID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.Image, &p.Username)
		if err != nil {
			return nil, err
		}
		if len(p.Username) > 0 {
			p.Username = strings.ToUpper(string(p.Username[0])) + p.Username[1:]
			p.Initials = string(p.Username[0])
		}
		// populate likes/dislikes counts for this post
		db.QueryRow(`SELECT COUNT(*) FROM likes WHERE postId = ?`, p.PostID).Scan(&p.Likes)
		db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE postId = ?`, p.PostID).Scan(&p.Dislikes)
		Posts = append(Posts, p)
	}

	return Posts, nil
}

func GetMyPosts(db *sql.DB, userId string) ([]models.Post, error) {
	query, err := db.Prepare(`SELECT posts.*,users.username FROM posts 
	 JOIN users on posts.userId=users.userId WHERE posts.userId=?`)
	if err != nil {
		fmt.Println("Error on My post: ", err)
		return nil, err
	}
	rows, err := query.Query(userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Posts []models.Post

	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.PostID, &p.UserID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.Image, &p.Username)
		if err != nil {
			println(err.Error())
			return nil, err
		}
		if len(p.Username) > 0 {
			p.Username = strings.ToUpper(string(p.Username[0])) + p.Username[1:]
			p.Initials = string(p.Username[0])
		}
		// populate likes/dislikes counts for this post
		db.QueryRow(`SELECT COUNT(*) FROM likes WHERE postId = ?`, p.PostID).Scan(&p.Likes)
		db.QueryRow(`SELECT COUNT(*) FROM dislikes WHERE postId = ?`, p.PostID).Scan(&p.Dislikes)
		Posts = append(Posts, p)
	}
	fmt.Println(Posts)

	return Posts, nil
}
