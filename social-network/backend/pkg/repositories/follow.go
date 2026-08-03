package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"social-network/backend/pkg/models"

	"github.com/gofrs/uuid"
)

var OnNotificationCreated func(notificationID, targetUserID string)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) HandleFollowRequest(requestID, notificationID, requesterID, targetID string) (models.FollowResult, *models.NotificationData, error) {
	var result models.FollowResult

	if requesterID == targetID {
		return result, nil, errors.New("you cannot follow yourself")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return result, nil, err
	}
	defer tx.Rollback()

	// Check if target is private
	var isPrivate int
	err = tx.QueryRow("SELECT is_private FROM users WHERE id = ?", targetID).Scan(&isPrivate)
	if err != nil {
		return result, nil, fmt.Errorf("failed to check target user privacy: %w", err)
	}

	now := time.Now().Unix()

	var requesterUsername string
	_ = tx.QueryRow("SELECT username FROM users WHERE id = ?", requesterID).Scan(&requesterUsername)

	if isPrivate == 0 {
		// Public profile: bypass request, follow immediately
		queryFollower := `
            INSERT INTO followers (follower_id, following_id, created_at)
            VALUES (?, ?, ?)
            ON CONFLICT(follower_id, following_id) DO NOTHING`

		_, err = tx.Exec(queryFollower, requesterID, targetID, now)
		if err != nil {
			return result, nil, fmt.Errorf("failed to insert follower: %w", err)
		}

		result.Following = true
		result.Status = "following"

		content := fmt.Sprintf("@%s started following you.", requesterUsername)
		queryNotification := `
			INSERT INTO notifications (id, user_id, type, content, is_read, ref_id, created_at)
			VALUES (?, ?, 'follow_accepted', ?, 0, '', ?)`

		_, err = tx.Exec(queryNotification, notificationID, targetID, content, now)
		if err != nil {
			return result, nil, fmt.Errorf("failed to save follow notification: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return result, nil, err
		}

		if OnNotificationCreated != nil {
			OnNotificationCreated(notificationID, targetID)
		}

		notifPayload := &models.NotificationData{
			NotificationID: notificationID,
			SenderID:       requesterID,
			Type:           "follow_accepted",
		}
		return result, notifPayload, nil
	}

	queryRequest := `
        INSERT INTO follow_requests (id, requester_id, target_id, status, created_at)
        VALUES (?, ?, ?, 'pending', ?)
        ON CONFLICT(requester_id, target_id) DO UPDATE SET status='pending', created_at=?`

	_, err = tx.Exec(queryRequest, requestID, requesterID, targetID, now, now)
	if err != nil {
		return result, nil, fmt.Errorf("failed to store follow request: %w", err)
	}

	content := fmt.Sprintf("User @%s requested to follow you.", requesterUsername)
	queryNotification := `
        INSERT INTO notifications (id, user_id, type, content, is_read, ref_id, created_at)
        VALUES (?, ?, 'follow_request', ?, 0, ?, ?)`

	_, err = tx.Exec(queryNotification, notificationID, targetID, content, requestID, now)
	if err != nil {
		return result, nil, fmt.Errorf("failed to save notification: %w", err)
	}

	result.Following = false
	result.Status = "pending"

	if err := tx.Commit(); err != nil {
		return result, nil, err
	}

	if OnNotificationCreated != nil {
		OnNotificationCreated(notificationID, targetID)
	}

	notifPayload := &models.NotificationData{
		NotificationID: notificationID,
		SenderID:       requesterID,
		Type:           "follow_request",
	}

	return result, notifPayload, nil
}
func (r *Repository) RespondToFollowRequest(requestID, targetID, action string) (models.FollowResult, error) {
	var result models.FollowResult

	if action != "accept" && action != "decline" {
		return result, errors.New("invalid action: use 'accept' or 'decline'")
	}

	var requesterID string
	var currentStatus string
	err := r.DB.QueryRow(
		`SELECT requester_id, status FROM follow_requests 
         WHERE id = ? AND target_id = ?`,
		requestID, targetID,
	).Scan(&requesterID, &currentStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("follow request not found for ID: %s", requestID)
		}
		return result, fmt.Errorf("database error fetching request: %w", err)
	}

	if currentStatus != "pending" {
		return result, fmt.Errorf("follow request already %s", currentStatus)
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return result, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().Unix()
	var notifID string

	switch action {
	case "accept":
		var alreadyFollowing int
		err = tx.QueryRow(
			`SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ?`,
			requesterID, targetID,
		).Scan(&alreadyFollowing)
		if err != nil {
			return result, fmt.Errorf("failed to check existing follow: %w", err)
		}

		if alreadyFollowing == 0 {
			_, err = tx.Exec(
				`INSERT INTO followers (follower_id, following_id, created_at)
                 VALUES (?, ?, ?)`,
				requesterID, targetID, now,
			)
			if err != nil {
				return result, fmt.Errorf("failed to insert follower: %w", err)
			}
		}

		result.Following = true
		result.Status = "accepted"

		var targetUsername string
		_ = tx.QueryRow("SELECT username FROM users WHERE id = ?", targetID).Scan(&targetUsername)
		notifID = uuid.Must(uuid.NewV4()).String()
		content := fmt.Sprintf("@%s accepted your follow request.", targetUsername)
		_, err = tx.Exec(
			`INSERT INTO notifications (id, user_id, type, content, is_read, ref_id, created_at)
             VALUES (?, ?, 'follow_accepted', ?, 0, '', ?)`,
			notifID, requesterID, content, now,
		)
		if err != nil {
			return result, fmt.Errorf("failed to save accepted notification: %w", err)
		}

	case "decline":
		result.Following = false
		result.Status = "declined"
	}

	// Delete the follow request notification
	_, err = tx.Exec(
		`DELETE FROM notifications WHERE type = 'follow_request' AND ref_id = ?`,
		requestID,
	)
	if err != nil {
		return result, fmt.Errorf("failed to delete follow request notification: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM follow_requests WHERE id = ?`, requestID)
	if err != nil {
		return result, fmt.Errorf("failed to delete follow request: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return result, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if action == "accept" && notifID != "" && OnNotificationCreated != nil {
		OnNotificationCreated(notifID, requesterID)
	}

	return result, nil
}

func (r *Repository) Unfollow(followerID, followingID string) error {
	if followerID == followingID {
		return errors.New("you cannot unfollow yourself")
	}

	result, err := r.DB.Exec(
		`DELETE FROM followers WHERE follower_id = ? AND following_id = ?`,
		followerID, followingID,
	)
	if err != nil {
		return fmt.Errorf("failed to unfollow: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("not found")
	}

	return nil
}

// GetFollowers returns paginated list of users who follow targetID.
// If target profile is private, only visible to followers of that user.
func (r *Repository) GetFollowers(targetID, sessionUserID string, page, perPage int) (*models.PaginatedUsers, error) {
	var isPrivate int
	err := r.DB.QueryRow(`SELECT is_private FROM users WHERE id = ?`, targetID).Scan(&isPrivate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if isPrivate == 1 && sessionUserID != targetID {
		var count int
		err := r.DB.QueryRow(
			`SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ?`,
			sessionUserID, targetID,
		).Scan(&count)
		if err != nil || count == 0 {
			return nil, errors.New("profile is private")
		}
	}

	offset := (page - 1) * perPage

	var total int
	err = r.DB.QueryRow(
		`SELECT COUNT(*) FROM followers WHERE following_id = ?`, targetID,
	).Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(`
        SELECT u.id, u.username, COALESCE(u.avatar_url, ''),
               COUNT(f2.follower_id) > 0 AS is_following
        FROM followers f
        JOIN users u ON u.id = f.follower_id
        LEFT JOIN followers f2 ON f2.follower_id = ? AND f2.following_id = u.id
        WHERE f.following_id = ?
        GROUP BY u.id
        LIMIT ? OFFSET ?`,
		sessionUserID, targetID, perPage, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query followers: %w", err)
	}
	defer rows.Close()

	users := []models.UserSummary{}
	for rows.Next() {
		var u models.UserSummary
		if err := rows.Scan(&u.ID, &u.Username, &u.AvatarURL, &u.IsFollowing); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return &models.PaginatedUsers{
		Users:   users,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}

// GetFollowing returns paginated list of users that targetID follows.
func (r *Repository) GetFollowing(targetID, sessionUserID string, page, perPage int) (*models.PaginatedUsers, error) {
	// Check privacy
	var isPrivate int
	err := r.DB.QueryRow(`SELECT is_private FROM users WHERE id = ?`, targetID).Scan(&isPrivate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if isPrivate == 1 && sessionUserID != targetID {
		var count int
		err := r.DB.QueryRow(
			`SELECT COUNT(*) FROM followers WHERE follower_id = ? AND following_id = ?`,
			sessionUserID, targetID,
		).Scan(&count)
		if err != nil || count == 0 {
			return nil, errors.New("profile is private")
		}
	}

	offset := (page - 1) * perPage

	var total int
	err = r.DB.QueryRow(
		`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, targetID,
	).Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(`
        SELECT u.id, u.username, COALESCE(u.avatar_url, ''),
               COUNT(f2.follower_id) > 0 AS is_following
        FROM followers f
        JOIN users u ON u.id = f.following_id
        LEFT JOIN followers f2 ON f2.follower_id = ? AND f2.following_id = u.id
        WHERE f.follower_id = ?
        GROUP BY u.id
        LIMIT ? OFFSET ?`,
		sessionUserID, targetID, perPage, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query following: %w", err)
	}
	defer rows.Close()

	users := []models.UserSummary{}
	for rows.Next() {
		var u models.UserSummary
		if err := rows.Scan(&u.ID, &u.Username, &u.AvatarURL, &u.IsFollowing); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return &models.PaginatedUsers{
		Users:   users,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	}, nil
}
