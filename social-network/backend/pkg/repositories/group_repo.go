package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"social-network/backend/pkg/models"

	"github.com/gofrs/uuid"
)

func CreateGroup(db *sql.DB, req models.CreateGroupRequest, creatorID string) (*models.Group, error) {
	groupID := uuid.Must(uuid.NewV4()).String()

	group := &models.Group{
		ID:          groupID,
		Title:       req.Title,
		Description: req.Description,
		CreatorID:   creatorID,
		CreatedAt:   time.Now().Unix(),
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO groups (id, title, description, creator_id, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		group.ID, group.Title, group.Description, group.CreatorID, group.CreatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("could not insert group: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO group_members (group_id, user_id, joined_at)
		VALUES (?, ?, ?)`,
		group.ID, group.CreatorID, group.CreatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("could not add creator as member: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}
	return group, nil
}

func GetAllGroups(db *sql.DB, userID string, page, limit int, search string) ([]models.GroupWithCount, error) {
	offset := (page - 1) * limit
	searchParam := "%" + search + "%"

	rows, err := db.Query(`
		SELECT g.id, g.title, g.description, g.creator_id, g.created_at,
			   COUNT(DISTINCT gm.user_id) as member_count,
			   EXISTS(SELECT 1 FROM group_members WHERE group_id = g.id AND user_id = ?) as is_member,
			   EXISTS(SELECT 1 FROM group_invites WHERE group_id = g.id AND invitee_id = ? AND status = 'pending') as is_invited,
			   (SELECT id FROM group_invites WHERE group_id = g.id AND invitee_id = ? AND status = 'pending' LIMIT 1) as invite_id,
			   EXISTS(SELECT 1 FROM group_join_requests WHERE group_id = g.id AND user_id = ? AND status = 'pending') as is_requested
		FROM groups g
		LEFT JOIN group_members gm ON g.id = gm.group_id
		WHERE g.title LIKE ?
		GROUP BY g.id
		ORDER BY g.created_at DESC
		LIMIT ? OFFSET ?`,
		userID, userID, userID, userID, searchParam, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("could not query groups: %w", err)
	}
	defer rows.Close()

	var groups []models.GroupWithCount
	for rows.Next() {
		var g models.GroupWithCount
		var inviteID sql.NullString
		err := rows.Scan(&g.ID, &g.Title, &g.Description, &g.CreatorID, &g.CreatedAt, &g.MemberCount, &g.IsMember, &g.IsInvited, &inviteID, &g.IsRequested)
		if err != nil {
			return nil, fmt.Errorf("could not scan group: %w", err)
		}
		if inviteID.Valid {
			g.InviteID = inviteID.String
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func GetGroupByID(db *sql.DB, groupID string) (*models.Group, error) {
	var gr models.Group
	err := db.QueryRow(`
		SELECT id, title, description, creator_id, created_at
		FROM groups WHERE id = ?`, groupID,
	).Scan(&gr.ID, &gr.Title, &gr.Description, &gr.CreatorID, &gr.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not get group: %w", err)
	}
	return &gr, nil
}

func GetPendingJoinRequests(db *sql.DB, groupID string) ([]models.GroupJoinRequestWithUser, error) {
	rows, err := db.Query(`
		SELECT r.id, r.group_id, r.user_id, r.status, r.created_at,
		       u.username, COALESCE(u.first_name, ''), COALESCE(u.last_name, '')
		FROM group_join_requests r
		JOIN users u ON r.user_id = u.id
		WHERE r.group_id = ? AND r.status = 'pending'`,
		groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not query join requests: %w", err)
	}
	defer rows.Close()

	var reqs []models.GroupJoinRequestWithUser
	for rows.Next() {
		var r models.GroupJoinRequestWithUser
		err := rows.Scan(&r.ID, &r.GroupID, &r.UserID, &r.Status, &r.CreatedAt, &r.Username, &r.FirstName, &r.LastName)
		if err != nil {
			return nil, fmt.Errorf("could not scan join request: %w", err)
		}
		reqs = append(reqs, r)
	}

	if reqs == nil {
		reqs = []models.GroupJoinRequestWithUser{}
	}
	return reqs, nil
}

func IsMember(db *sql.DB, groupID, userID string) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM group_members
		WHERE group_id = ? AND user_id = ?`,
		groupID, userID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("could not check membership: %w", err)
	}
	return count > 0, nil
}

func IsGroupCreator(db *sql.DB, groupID, userID string) (bool, error) {
	var creatorID string
	err := db.QueryRow(`
		SELECT creator_id FROM groups WHERE id = ?`, groupID,
	).Scan(&creatorID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("could not check creator: %w", err)
	}
	return creatorID == userID, nil
}

func AddMember(db *sql.DB, groupID, userID string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE INTO group_members (group_id, user_id, joined_at)
		VALUES (?, ?, ?)`,
		groupID, userID, time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("could not add member: %w", err)
	}
	return nil
}

func CreateInvite(db *sql.DB, groupID, inviterID, inviteeID string) (*models.GroupInvite, error) {
	member, err := IsMember(db, groupID, inviterID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if !member {
		return nil, fmt.Errorf("only group members can send invites")
	}

	alreadyMember, err := IsMember(db, groupID, inviteeID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if alreadyMember {
		return nil, fmt.Errorf("user is already a member")
	}

	inviteID := uuid.Must(uuid.NewV4()).String()

	invite := &models.GroupInvite{
		ID:        inviteID,
		GroupID:   groupID,
		InviterID: inviterID,
		InviteeID: inviteeID,
		Status:    "pending",
		CreatedAt: time.Now().Unix(),
	}

	_, err = db.Exec(`
		INSERT INTO group_invites (id, group_id, inviter_id, invitee_id, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		invite.ID, invite.GroupID, invite.InviterID, invite.InviteeID, invite.Status, invite.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create invite: %w", err)
	}

	var groupTitle string
	err = db.QueryRow(`SELECT title FROM groups WHERE id = ?`, groupID).Scan(&groupTitle)
	if err == nil {
		notifID := uuid.Must(uuid.NewV4()).String()
		_, errNotif := db.Exec(`
			INSERT INTO notifications (id, user_id, type, content, is_read, ref_id, created_at)
			VALUES (?, ?, 'group_invite', ?, 0, ?, ?)`,
			notifID, inviteeID, fmt.Sprintf("You have been invited to join group: %s", groupTitle), invite.ID, time.Now().Unix(),
		)
		if errNotif == nil && OnNotificationCreated != nil {
			OnNotificationCreated(notifID, inviteeID)
		}
	}

	return invite, nil
}

func UpdateInviteStatus(db *sql.DB, inviteID, userID, status string) error {
	var inviteeID string
	var groupID string
	err := db.QueryRow(`
		SELECT invitee_id, group_id FROM group_invites WHERE id = ?`, inviteID,
	).Scan(&inviteeID, &groupID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("invite not found")
	}
	if err != nil {
		return fmt.Errorf("could not get invite: %w", err)
	}

	if inviteeID != userID {
		return fmt.Errorf("unauthorized")
	}

	_, err = db.Exec(`
		UPDATE group_invites SET status = ? WHERE id = ?`,
		status, inviteID,
	)
	if err != nil {
		return fmt.Errorf("could not update invite: %w", err)
	}

	// Delete the original group_invite notification for the invitee so it doesn't reappear
	_, _ = db.Exec(
		`DELETE FROM notifications WHERE type = 'group_invite' AND ref_id = ?`,
		inviteID,
	)

	if status == "accepted" {
		if err := AddMember(db, groupID, userID); err != nil {
			return fmt.Errorf("could not add member: %w", err)
		}
	}

	// Notify the inviter — but only once per invite
	var inviterID string
	_ = db.QueryRow("SELECT inviter_id FROM group_invites WHERE id = ?", inviteID).Scan(&inviterID)

	var inviteeUsername string
	_ = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&inviteeUsername)

	var groupTitle string
	_ = db.QueryRow("SELECT title FROM groups WHERE id = ?", groupID).Scan(&groupTitle)

	// Guard: don't insert a duplicate notification for the same invite
	var existing int
	_ = db.QueryRow(
		`SELECT COUNT(*) FROM notifications WHERE type = 'group_invite_accepted' AND ref_id = ?`,
		inviteID,
	).Scan(&existing)

	if existing == 0 {
		notifID := uuid.Must(uuid.NewV4()).String()
		var content string
		if status == "accepted" {
			content = fmt.Sprintf("@%s accepted your invite to join group '%s'.", inviteeUsername, groupTitle)
		} else {
			content = fmt.Sprintf("@%s declined your invite to join group '%s'.", inviteeUsername, groupTitle)
		}

		_, errInsert := db.Exec(`
			INSERT INTO notifications (id, user_id, type, content, is_read, ref_id, created_at)
			VALUES (?, ?, 'group_invite_accepted', ?, 0, ?, ?)`,
			notifID, inviterID, content, inviteID, time.Now().Unix(),
		)
		if errInsert == nil && OnNotificationCreated != nil {
			OnNotificationCreated(notifID, inviterID)
		}
	}

	return nil
}

func CreateJoinRequest(db *sql.DB, groupID, userID string) (*models.GroupJoinRequest, error) {
	already, err := IsMember(db, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if already {
		return nil, fmt.Errorf("already a member")
	}

	reqID := uuid.Must(uuid.NewV4()).String()

	req := &models.GroupJoinRequest{
		ID:        reqID,
		GroupID:   groupID,
		UserID:    userID,
		Status:    "pending",
		CreatedAt: time.Now().Unix(),
	}

	_, err = db.Exec(`
		INSERT INTO group_join_requests (id, group_id, user_id, status, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		req.ID, req.GroupID, req.UserID, req.Status, req.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create join request: %w", err)
	}

	var groupTitle, creatorID string
	err = db.QueryRow(`SELECT title, creator_id FROM groups WHERE id = ?`, groupID).Scan(&groupTitle, &creatorID)
	if err == nil {
		var requesterUsername string
		err = db.QueryRow(`SELECT username FROM users WHERE id = ?`, userID).Scan(&requesterUsername)
		if err == nil {
			notifID := uuid.Must(uuid.NewV4()).String()
			_, errNotif := db.Exec(`
				INSERT INTO notifications (id, user_id, type, content, is_read, ref_id, created_at)
				VALUES (?, ?, 'group_join_request', ?, 0, ?, ?)`,
				notifID, creatorID, fmt.Sprintf("User @%s requested to join group: %s", requesterUsername, groupTitle), req.ID, time.Now().Unix(),
			)
			if errNotif == nil && OnNotificationCreated != nil {
				OnNotificationCreated(notifID, creatorID)
			}
		}
	}

	return req, nil
}

func UpdateRequestStatus(db *sql.DB, requestID, creatorID, status string) error {
	var userID string
	var groupID string
	err := db.QueryRow(`
		SELECT user_id, group_id FROM group_join_requests WHERE id = ?`, requestID,
	).Scan(&userID, &groupID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("request not found")
	}
	if err != nil {
		return fmt.Errorf("could not get request: %w", err)
	}

	isCreator, err := IsGroupCreator(db, groupID, creatorID)
	if err != nil {
		return fmt.Errorf("could not check creator: %w", err)
	}
	if !isCreator {
		return fmt.Errorf("only the group creator can approve or decline requests")
	}

	_, err = db.Exec(`
		UPDATE group_join_requests SET status = ? WHERE id = ?`,
		status, requestID,
	)
	if err != nil {
		return fmt.Errorf("could not update request: %w", err)
	}

	// Delete the original group_join_request notification for the creator so it doesn't reappear
	_, _ = db.Exec(
		`DELETE FROM notifications WHERE type = 'group_join_request' AND ref_id = ?`,
		requestID,
	)

	if status == "approved" {
		if err := AddMember(db, groupID, userID); err != nil {
			return fmt.Errorf("could not add member: %w", err)
		}
	}

	// Notify the requester — but only once
	var groupTitle string
	_ = db.QueryRow("SELECT title FROM groups WHERE id = ?", groupID).Scan(&groupTitle)

	var existing int
	_ = db.QueryRow(
		`SELECT COUNT(*) FROM notifications WHERE type = 'group_join_response' AND ref_id = ?`,
		requestID,
	).Scan(&existing)

	if existing == 0 {
		notifID := uuid.Must(uuid.NewV4()).String()
		var content string
		if status == "approved" {
			content = fmt.Sprintf("Your request to join group '%s' was approved.", groupTitle)
		} else {
			content = fmt.Sprintf("Your request to join group '%s' was declined.", groupTitle)
		}

		_, errInsert := db.Exec(`
			INSERT INTO notifications (id, user_id, type, content, is_read, ref_id, created_at)
			VALUES (?, ?, 'group_join_response', ?, 0, ?, ?)`,
			notifID, userID, content, requestID, time.Now().Unix(),
		)
		if errInsert == nil && OnNotificationCreated != nil {
			OnNotificationCreated(notifID, userID)
		}
	}

	return nil
}

func CreateGroupPost(db *sql.DB, groupID, authorID, content string) (*models.GroupPost, error) {
	member, err := IsMember(db, groupID, authorID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if !member {
		return nil, fmt.Errorf("only group members can create posts")
	}

	postID := uuid.Must(uuid.NewV4()).String()

	post := &models.GroupPost{
		ID:        postID,
		GroupID:   groupID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	_, err = db.Exec(`
		INSERT INTO group_posts (id, group_id, author_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		post.ID, post.GroupID, post.AuthorID, post.Content, post.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create group post: %w", err)
	}

	return post, nil
}

func GetGroupPosts(db *sql.DB, groupID, requestingUserID string, page, limit int) ([]models.GroupPostWithMeta, error) {
	member, err := IsMember(db, groupID, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if !member {
		return nil, fmt.Errorf("only group members can view posts")
	}

	offset := (page - 1) * limit

	rows, err := db.Query(`
		SELECT 
			gp.id,
			gp.group_id,
			gp.author_id,
			u.username,
			gp.content,
			gp.created_at,
			COUNT(gc.id) as comment_count
		FROM group_posts gp
		JOIN users u ON gp.author_id = u.id
		LEFT JOIN group_comments gc ON gp.id = gc.post_id
		WHERE gp.group_id = ?
		GROUP BY gp.id
		ORDER BY gp.created_at DESC
		LIMIT ? OFFSET ?`,
		groupID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("could not query group posts: %w", err)
	}
	defer rows.Close()

	var posts []models.GroupPostWithMeta
	for rows.Next() {
		var p models.GroupPostWithMeta
		err := rows.Scan(
			&p.ID, &p.GroupID, &p.AuthorID, &p.AuthorUsername,
			&p.Content, &p.CreatedAt, &p.CommentCount,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func CreateGroupComment(db *sql.DB, groupID, postID, authorID, content string) (*models.GroupComment, error) {
	member, err := IsMember(db, groupID, authorID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if !member {
		return nil, fmt.Errorf("only group members can comment")
	}

	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM group_posts
		WHERE id = ? AND group_id = ?`,
		postID, groupID,
	).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("could not verify post: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("post not found in this group")
	}

	commentID := uuid.Must(uuid.NewV4()).String()

	comment := &models.GroupComment{
		ID:        commentID,
		PostID:    postID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now().Unix(),
	}

	_, err = db.Exec(`
		INSERT INTO group_comments (id, post_id, author_id, content, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		comment.ID, comment.PostID, comment.AuthorID, comment.Content, comment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create comment: %w", err)
	}

	return comment, nil
}

func GetGroupComments(db *sql.DB, groupID, postID, requestingUserID string) ([]models.GroupCommentWithMeta, error) {
	member, err := IsMember(db, groupID, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if !member {
		return nil, fmt.Errorf("only group members can view comments")
	}

	var count int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM group_posts
		WHERE id = ? AND group_id = ?`,
		postID, groupID,
	).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("could not verify post: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("post not found in this group")
	}

	rows, err := db.Query(`
		SELECT 
			gc.id, gc.post_id, gc.author_id,
			u.username, gc.content, gc.created_at
		FROM group_comments gc
		JOIN users u ON gc.author_id = u.id
		WHERE gc.post_id = ?
		ORDER BY gc.created_at ASC`,
		postID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not query comments: %w", err)
	}
	defer rows.Close()

	var comments []models.GroupCommentWithMeta
	for rows.Next() {
		var c models.GroupCommentWithMeta
		err := rows.Scan(
			&c.ID, &c.PostID, &c.AuthorID,
			&c.AuthorUsername, &c.Content, &c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("could not scan comment: %w", err)
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func CreateEvent(db *sql.DB, groupID, creatorID string, req models.CreateEventRequest) (*models.Event, error) {
	member, err := IsMember(db, groupID, creatorID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if !member {
		return nil, fmt.Errorf("only group members can create an event")
	}

	eventID := uuid.Must(uuid.NewV4()).String()

	event := &models.Event{
		ID:          eventID,
		GroupID:     groupID,
		CreatorID:   creatorID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   req.EventDate,
		CreatedAt:   time.Now().Unix(),
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("could not begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the event
	_, err = tx.Exec(
		`INSERT INTO events(id, group_id, creator_id, title, description, event_date, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.GroupID, event.CreatorID,
		event.Title, event.Description, event.EventDate, event.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create event: %w", err)
	}

	// Get group title for notification
	var groupTitle string
	err = tx.QueryRow(`SELECT title FROM groups WHERE id = ?`, groupID).Scan(&groupTitle)
	if err != nil {
		// If we can't get the group title, still return the event
		// but log the error
		log.Printf("Warning: Failed to get group title for notification: %v", err)
	} else {
		// Get all group members except the creator
		rows, err := tx.Query(`SELECT user_id FROM group_members WHERE group_id = ? AND user_id != ?`, groupID, creatorID)
		if err != nil {
			log.Printf("Warning: Failed to query group members for notifications: %v", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var memberID string
				if err := rows.Scan(&memberID); err != nil {
					log.Printf("Warning: Failed to scan member ID: %v", err)
					continue
				}

				notifID := uuid.Must(uuid.NewV4()).String()

				_, errNotif := tx.Exec(`
                    INSERT INTO notifications (id, user_id, type, content, is_read, ref_id, created_at)
                    VALUES (?, ?, 'event', ?, 0, ?, ?)`,
					notifID, memberID,
					fmt.Sprintf("New event '%s' created in group: %s", req.Title, groupTitle),
					event.ID, time.Now().Unix(),
				)
				if errNotif != nil {
					log.Printf("Warning: Failed to insert notification for user %s: %v", memberID, errNotif)
				}
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("could not commit transaction: %w", err)
	}

	// NOW send notifications after commit - this ensures the notification is in the database
	// Get all members again and send notifications
	rows, err := db.Query(`SELECT user_id FROM group_members WHERE group_id = ? AND user_id != ?`, groupID, creatorID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var memberID string
			if err := rows.Scan(&memberID); err == nil {
				// Get the notification ID for this member
				var notifID string
				err := db.QueryRow(`
                    SELECT id FROM notifications 
                    WHERE ref_id = ? AND user_id = ? AND type = 'event' 
                    ORDER BY created_at DESC LIMIT 1`,
					event.ID, memberID,
				).Scan(&notifID)
				if err == nil && OnNotificationCreated != nil {
					OnNotificationCreated(notifID, memberID)
				}
			}
		}
	}

	return event, nil
}

func GetGroupEvents(db *sql.DB, groupID, userID string) ([]models.EventWithResponses, error) {
	member, err := IsMember(db, groupID, userID)
	if err != nil {
		return nil, fmt.Errorf("could not check membership: %w", err)
	}
	if !member {
		return nil, fmt.Errorf("only group members can view events")
	}

	rows, err := db.Query(`
		SELECT e.id, e.group_id, e.creator_id, e.title, e.description, e.event_date, e.created_at,
		       COALESCE((SELECT COUNT(*) FROM event_responses WHERE event_id = e.id AND response = 'going'), 0) as going,
		       COALESCE((SELECT COUNT(*) FROM event_responses WHERE event_id = e.id AND response = 'not_going'), 0) as not_going,
		       COALESCE((SELECT response FROM event_responses WHERE event_id = e.id AND user_id = ?), '') as user_response
		FROM events e
		WHERE e.group_id = ?
		ORDER BY e.event_date ASC`,
		userID, groupID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not query events: %w", err)
	}
	defer rows.Close()

	var events []models.EventWithResponses
	for rows.Next() {
		var e models.EventWithResponses
		err := rows.Scan(&e.ID, &e.GroupID, &e.CreatorID, &e.Title, &e.Description, &e.EventDate, &e.CreatedAt, &e.Going, &e.NotGoing, &e.UserResponse)
		if err != nil {
			return nil, fmt.Errorf("could not scan event: %w", err)
		}
		events = append(events, e)
	}
	return events, nil
}

func RespondToEvent(db *sql.DB, eventID, userID, response string) (*models.EventResponse, error) {
	var groupID string
	err := db.QueryRow("SELECT group_id FROM events WHERE id = ?", eventID).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, err
	}

	member, err := IsMember(db, groupID, userID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, fmt.Errorf("only group members can respond to events")
	}

	respID := uuid.Must(uuid.NewV4()).String()

	resp := &models.EventResponse{
		ID:        respID,
		EventID:   eventID,
		UserID:    userID,
		Response:  response,
		CreatedAt: time.Now().Unix(),
	}

	_, err = db.Exec(`
		INSERT INTO event_responses (id, event_id, user_id, response, created_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(event_id, user_id) DO UPDATE SET response = ?, created_at = ?`,
		resp.ID, resp.EventID, resp.UserID, resp.Response, resp.CreatedAt,
		resp.Response, resp.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("could not save event response: %w", err)
	}

	return resp, nil
}

func UpdateGroup(db *sql.DB, groupID, requestingUserID, title, description string) (*models.Group, error) {
	isCreator, err := IsGroupCreator(db, groupID, requestingUserID)
	if err != nil {
		return nil, fmt.Errorf("could not verify creator: %w", err)
	}
	if !isCreator {
		return nil, fmt.Errorf("only the group creator can edit the group")
	}

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	_, err = db.Exec(`UPDATE groups SET title = ?, description = ? WHERE id = ?`, title, description, groupID)
	if err != nil {
		return nil, fmt.Errorf("could not update group: %w", err)
	}

	return GetGroupByID(db, groupID)
}
