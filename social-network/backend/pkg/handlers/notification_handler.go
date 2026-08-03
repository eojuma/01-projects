package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"social-network/backend/pkg/websocket"
)

type NotificationResponse struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
	Actionable bool   `json:"actionable"`
	Read       bool   `json:"read"`
	RefID      string `json:"ref_id,omitempty"`
	GroupID    string `json:"group_id,omitempty"`
}

func GetNotificationsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := authUserIDFromContext(r)
		if !ok || userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query(`
			SELECT id, type, content, is_read, ref_id, created_at
			FROM notifications
			WHERE user_id = ?
			ORDER BY created_at DESC`,
			userID,
		)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notifs []NotificationResponse
		for rows.Next() {
			var id, notifType, content, refID sql.NullString
			var isRead int
			var createdAt int64

			err := rows.Scan(&id, &notifType, &content, &isRead, &refID, &createdAt)
			if err != nil {
				continue
			}

			n := NotificationResponse{
				ID:         id.String,
				Type:       notifType.String,
				Message:    content.String,
				Timestamp:  time.Unix(createdAt, 0).Format(time.RFC3339),
				Read:       isRead == 1,
				RefID:      refID.String,
				Actionable: notifType.String == "follow_request" || notifType.String == "group_invite" || notifType.String == "group_join_request",
			}

			// Look up GroupID using RefID for group-related notifications
			if n.Type == "group_invite" && n.RefID != "" {
				_ = db.QueryRow("SELECT group_id FROM group_invites WHERE id = ?", n.RefID).Scan(&n.GroupID)
			} else if n.Type == "group_join_request" && n.RefID != "" {
				_ = db.QueryRow("SELECT group_id FROM group_join_requests WHERE id = ?", n.RefID).Scan(&n.GroupID)
			} else if n.Type == "event" && n.RefID != "" {
				_ = db.QueryRow("SELECT group_id FROM events WHERE id = ?", n.RefID).Scan(&n.GroupID)
			}

			notifs = append(notifs, n)
		}

		if notifs == nil {
			notifs = []NotificationResponse{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notifs)
	}
}

func MarkNotificationsReadHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, ok := authUserIDFromContext(r)
		if !ok || userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		notifID := r.URL.Query().Get("id")
		if notifID == "" {
			var req struct {
				ID string `json:"id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&req)
			notifID = req.ID
		}

		if notifID != "" {
			_, err := db.Exec("UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?", notifID, userID)
			if err != nil {
				http.Error(w, "database error", http.StatusInternalServerError)
				return
			}
		} else {
			_, err := db.Exec("UPDATE notifications SET is_read = 1 WHERE user_id = ?", userID)
			if err != nil {
				http.Error(w, "database error", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func PushNotification(db *sql.DB, notifID, userID string) {
	var id, notifType, content, refID sql.NullString
	var isRead int
	var createdAt int64

	err := db.QueryRow(`
		SELECT id, type, content, is_read, ref_id, created_at
		FROM notifications
		WHERE id = ? AND user_id = ?`,
		notifID, userID,
	).Scan(&id, &notifType, &content, &isRead, &refID, &createdAt)
	if err != nil {
		return
	}

	n := NotificationResponse{
		ID:         id.String,
		Type:       notifType.String,
		Message:    content.String,
		Timestamp:  time.Unix(createdAt, 0).Format(time.RFC3339),
		Read:       isRead == 1,
		RefID:      refID.String,
		Actionable: notifType.String == "follow_request" || notifType.String == "group_invite" || notifType.String == "group_join_request",
	}

	// Look up GroupID using RefID for group-related notifications
	if n.Type == "group_invite" && n.RefID != "" {
		_ = db.QueryRow("SELECT group_id FROM group_invites WHERE id = ?", n.RefID).Scan(&n.GroupID)
	} else if n.Type == "group_join_request" && n.RefID != "" {
		_ = db.QueryRow("SELECT group_id FROM group_join_requests WHERE id = ?", n.RefID).Scan(&n.GroupID)
	} else if n.Type == "event" && n.RefID != "" {
		_ = db.QueryRow("SELECT group_id FROM events WHERE id = ?", n.RefID).Scan(&n.GroupID)
	}

	if websocket.ActiveHub != nil {
		websocket.ActiveHub.SendToUser(userID, n)
	}
}

