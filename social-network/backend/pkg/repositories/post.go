package repositories

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"social-network/backend/pkg/models"
)

func (r *Repository) CreatePost(input models.CreatePostInput) (*models.Post, error) {
	if input.Content == "" && len(input.ImagePaths) == 0 {
		return nil, errors.New("content or image is required")
	}
	if input.Privacy != "public" && input.Privacy != "followers" && input.Privacy != "selected" {
		return nil, errors.New("invalid privacy value")
	}
	if input.Privacy == "selected" && len(input.AllowedUsers) == 0 {
		return nil, errors.New("allowed_users is required for selected privacy")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now().Unix()

	_, err = tx.Exec(
		`INSERT INTO posts (id, author_id, content, privacy, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		input.ID, input.AuthorID, input.Content, input.Privacy, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert post: %w", err)
	}

	if len(input.ImagePaths) > 0 {
		for idx, imgPath := range input.ImagePaths {
			imageID := fmt.Sprintf("%s_img_%d", input.ID, idx)
			_, err = tx.Exec(
				`INSERT INTO post_images (id, post_id, image_path, created_at)
				 VALUES (?, ?, ?, ?)`,
				imageID, input.ID, imgPath, now,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to insert post image: %w", err)
			}
		}
	}

	if input.Privacy == "selected" {
		for _, userID := range input.AllowedUsers {
			_, err = tx.Exec(
				`INSERT INTO post_permissions (post_id, user_id) VALUES (?, ?)`,
				input.ID, userID,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to insert post permission for user %s: %w", userID, err)
			}
		}
	}

	var authorUsername string
	err = tx.QueryRow(`SELECT username FROM users WHERE id = ?`, input.AuthorID).Scan(&authorUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch author username: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Post{
		ID:             input.ID,
		AuthorID:       input.AuthorID,
		AuthorUsername: authorUsername,
		Content:        input.Content,
		Privacy:        input.Privacy,
		CreatedAt:      now,
	}, nil
}

func (r *Repository) GetFeed(sessionUserID, cursor string, limit int) (*models.FeedResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	args := []interface{}{
		sessionUserID, sessionUserID, sessionUserID, sessionUserID,
		sessionUserID, sessionUserID, sessionUserID, sessionUserID,
		sessionUserID,
	}

	query := `
		SELECT p.id, p.author_id, p.content, p.privacy, p.created_at,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count,
		       u.username,
		       (SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS likes_count,
		       EXISTS(SELECT 1 FROM post_likes WHERE post_id = p.id AND user_id = ?) AS liked,
		       CASE
		         WHEN p.author_id = ? THEN 'self'
		         WHEN EXISTS(SELECT 1 FROM followers WHERE follower_id = ? AND following_id = p.author_id) THEN 'following'
		         WHEN EXISTS(SELECT 1 FROM follow_requests WHERE requester_id = ? AND target_id = p.author_id AND status = 'pending') THEN 'pending'
		         ELSE 'none'
		       END AS follow_status,
		       u.is_private AS author_is_private
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE (
			u.is_private = 0
			OR (u.is_private = 1 AND EXISTS (
				SELECT 1 FROM followers WHERE follower_id = ? AND following_id = u.id
			))
			OR u.id = ?
		)
		AND (
			p.privacy = 'public'
			OR p.author_id = ?
			OR (p.privacy = 'followers' AND p.author_id IN (
				SELECT following_id FROM followers WHERE follower_id = ?
			))
			OR (p.privacy = 'selected' AND p.id IN (
				SELECT post_id FROM post_permissions WHERE user_id = ?
			))
		)`

	if cursor != "" {
		decoded, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		parts := strings.SplitN(string(decoded), "_", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid cursor")
		}
		cursorCreatedAt, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		cursorID := parts[1]
		query += ` AND (p.created_at < ? OR (p.created_at = ? AND p.id < ?))`
		args = append(args, cursorCreatedAt, cursorCreatedAt, cursorID)
	}

	query += ` ORDER BY p.created_at DESC, p.id DESC LIMIT ?`
	args = append(args, limit+1)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query feed: %w", err)
	}
	defer rows.Close()

	posts := []models.FeedPost{}
	for rows.Next() {
		var p models.FeedPost
		if err := rows.Scan(
			&p.ID, &p.AuthorID, &p.Content, &p.Privacy, &p.CreatedAt,
			&p.CommentCount, &p.AuthorUsername, &p.LikesCount, &p.Liked,
			&p.FollowStatus, &p.AuthorIsPrivate,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	nextCursor := ""
	if len(posts) > limit {
		last := posts[limit-1]
		nextCursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d_%s", last.CreatedAt, last.ID)))
		posts = posts[:limit]
	}

	for i := range posts {
		imgRows, err := r.DB.Query(
			`SELECT image_path FROM post_images WHERE post_id = ?`, posts[i].ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query post images: %w", err)
		}
		var images []string
		for imgRows.Next() {
			var path string
			if err := imgRows.Scan(&path); err != nil {
				imgRows.Close()
				return nil, err
			}
			images = append(images, path)
		}
		imgRows.Close()
		posts[i].ImageURLs = images
	}

	return &models.FeedResponse{
		Posts:      posts,
		NextCursor: nextCursor,
	}, nil
}

func (r *Repository) CanViewPost(sessionUserID, postID string) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(*) FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = ?
		AND (
			u.is_private = 0
			OR (u.is_private = 1 AND EXISTS (
				SELECT 1 FROM followers WHERE follower_id = ? AND following_id = u.id
			))
			OR u.id = ?
		)
		AND (
			p.privacy = 'public'
			OR p.author_id = ?
			OR (p.privacy = 'followers' AND p.author_id IN (
				SELECT following_id FROM followers WHERE follower_id = ?
			))
			OR (p.privacy = 'selected' AND p.id IN (
				SELECT post_id FROM post_permissions WHERE user_id = ?
			))
		)`,
		postID, sessionUserID, sessionUserID, sessionUserID, sessionUserID, sessionUserID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check post visibility: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) CreateComment(sessionUserID string, input models.CreateCommentInput) (*models.Comment, error) {
	if input.Content == "" && input.ImagePath == "" {
		return nil, errors.New("content or image is required")
	}

	var postExists int
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM posts WHERE id = ?`, input.PostID).Scan(&postExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check post: %w", err)
	}
	if postExists == 0 {
		return nil, errors.New("post not found")
	}

	canView, err := r.CanViewPost(sessionUserID, input.PostID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, errors.New("forbidden")
	}

	now := time.Now().Unix()

	_, err = r.DB.Exec(
		`INSERT INTO comments (id, post_id, author_id, content, image_path, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		input.ID, input.PostID, input.AuthorID, input.Content, nullableString(input.ImagePath), now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert comment: %w", err)
	}

	var author models.UserSummary
	err = r.DB.QueryRow(
		`SELECT id, username, COALESCE(avatar_url, '') FROM users WHERE id = ?`,
		input.AuthorID,
	).Scan(&author.ID, &author.Username, &author.AvatarURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch author: %w", err)
	}

	return &models.Comment{
		ID:        input.ID,
		PostID:    input.PostID,
		AuthorID:  input.AuthorID,
		Content:   input.Content,
		ImagePath: input.ImagePath,
		CreatedAt: now,
		Author:    author,
	}, nil
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func (r *Repository) GetUserPosts(sessionUserID, targetUserID, cursor string, limit int) (*models.FeedResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	var isPrivate int
	err := r.DB.QueryRow(`SELECT is_private FROM users WHERE id = ?`, targetUserID).Scan(&isPrivate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if isPrivate == 1 && sessionUserID != targetUserID {
		var isFollower int
		err := r.DB.QueryRow(
			`SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ?`,
			sessionUserID, targetUserID,
		).Scan(&isFollower)
		if err != nil {
			return nil, err
		}
		if isFollower == 0 {
			return &models.FeedResponse{Posts: []models.FeedPost{}}, nil
		}
	}

	args := []interface{}{sessionUserID, targetUserID, sessionUserID, sessionUserID, sessionUserID}

	query := `
		SELECT p.id, p.author_id, p.content, p.privacy, p.created_at,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count,
		       u.username,
		       (SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS likes_count,
		       EXISTS(SELECT 1 FROM post_likes WHERE post_id = p.id AND user_id = ?) AS liked
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.author_id = ?
		AND (
			p.privacy = 'public'
			OR p.author_id = ?
			OR (p.privacy = 'followers' AND p.author_id IN (
				SELECT following_id FROM followers WHERE follower_id = ?
			))
			OR (p.privacy = 'selected' AND p.id IN (
				SELECT post_id FROM post_permissions WHERE user_id = ?
			))
		)`

	if cursor != "" {
		decoded, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		parts := strings.SplitN(string(decoded), "_", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid cursor")
		}
		cursorCreatedAt, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		cursorID := parts[1]
		query += ` AND (p.created_at < ? OR (p.created_at = ? AND p.id < ?))`
		args = append(args, cursorCreatedAt, cursorCreatedAt, cursorID)
	}

	query += ` ORDER BY p.created_at DESC, p.id DESC LIMIT ?`
	args = append(args, limit+1)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query user posts: %w", err)
	}
	defer rows.Close()

	posts := []models.FeedPost{}
	for rows.Next() {
		var p models.FeedPost
		if err := rows.Scan(
			&p.ID, &p.AuthorID, &p.Content, &p.Privacy, &p.CreatedAt,
			&p.CommentCount, &p.AuthorUsername, &p.LikesCount, &p.Liked,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	nextCursor := ""
	if len(posts) > limit {
		last := posts[limit-1]
		nextCursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d_%s", last.CreatedAt, last.ID)))
		posts = posts[:limit]
	}

	for i := range posts {
		imgRows, err := r.DB.Query(
			`SELECT image_path FROM post_images WHERE post_id = ?`, posts[i].ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query post images: %w", err)
		}
		var images []string
		for imgRows.Next() {
			var path string
			if err := imgRows.Scan(&path); err != nil {
				imgRows.Close()
				return nil, err
			}
			images = append(images, path)
		}
		imgRows.Close()
		posts[i].ImageURLs = images
	}

	return &models.FeedResponse{
		Posts:      posts,
		NextCursor: nextCursor,
	}, nil
}

func (r *Repository) GetComments(sessionUserID string, postID string) ([]models.Comment, error) {
	var postExists int
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM posts WHERE id = ?`, postID).Scan(&postExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check post: %w", err)
	}
	if postExists == 0 {
		return nil, errors.New("post not found")
	}

	canView, err := r.CanViewPost(sessionUserID, postID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, errors.New("forbidden")
	}

	rows, err := r.DB.Query(
		`SELECT c.id, c.post_id, c.author_id, c.content, COALESCE(c.image_path, ''), c.created_at,
		        u.id, u.username, COALESCE(u.avatar_url, '')
		 FROM comments c
		 JOIN users u ON c.author_id = u.id
		 WHERE c.post_id = ?
		 ORDER BY c.created_at ASC`,
		postID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(
			&comment.ID, &comment.PostID, &comment.AuthorID, &comment.Content, &comment.ImagePath, &comment.CreatedAt,
			&comment.Author.ID, &comment.Author.Username, &comment.Author.AvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if comments == nil {
		comments = []models.Comment{}
	}

	return comments, nil
}

func (r *Repository) ToggleLike(sessionUserID, postID string) (bool, int, error) {
	var postExists int
	err := r.DB.QueryRow(`SELECT COUNT(*) FROM posts WHERE id = ?`, postID).Scan(&postExists)
	if err != nil {
		return false, 0, fmt.Errorf("failed to check post: %w", err)
	}
	if postExists == 0 {
		return false, 0, errors.New("post not found")
	}

	canView, err := r.CanViewPost(sessionUserID, postID)
	if err != nil {
		return false, 0, err
	}
	if !canView {
		return false, 0, errors.New("forbidden")
	}

	var alreadyLiked bool
	err = r.DB.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = ? AND user_id = ?)`,
		postID, sessionUserID,
	).Scan(&alreadyLiked)
	if err != nil {
		return false, 0, fmt.Errorf("failed to check like: %w", err)
	}

	if alreadyLiked {
		_, err = r.DB.Exec(`DELETE FROM post_likes WHERE post_id = ? AND user_id = ?`, postID, sessionUserID)
		if err != nil {
			return false, 0, fmt.Errorf("failed to delete like: %w", err)
		}
	} else {
		_, err = r.DB.Exec(
			`INSERT INTO post_likes (post_id, user_id, created_at) VALUES (?, ?, ?)`,
			postID, sessionUserID, time.Now().Unix(),
		)
		if err != nil {
			return false, 0, fmt.Errorf("failed to insert like: %w", err)
		}
	}

	var count int
	err = r.DB.QueryRow(`SELECT COUNT(*) FROM post_likes WHERE post_id = ?`, postID).Scan(&count)
	if err != nil {
		return false, 0, fmt.Errorf("failed to query count: %w", err)
	}

	return !alreadyLiked, count, nil
}
