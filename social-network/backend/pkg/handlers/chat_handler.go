package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"social-network/backend/pkg/repositories"
)

type ConversationResponse struct {
	ID            string `json:"id"`
	OtherUserID   string `json:"other_user_id"`
	OtherUsername string `json:"other_username"`
	AvatarURL     string `json:"avatar_url"`
	LastMessage   string `json:"last_message"`
	LastMessageAt int64  `json:"last_message_at"`
	UnreadCount   int    `json:"unread_count"`
}

type MessageResponse struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	CreatedAt      int64  `json:"created_at"`
	IsRead         bool   `json:"is_read"`
}

type GroupMessageResponse struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	SenderID  string `json:"sender_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

func GetConversationsHandler(db *sql.DB) http.HandlerFunc {
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
			SELECT c.id, u.id, u.username, COALESCE(u.avatar_url, ''),
			       COALESCE((SELECT content FROM messages WHERE conversation_id = c.id ORDER BY created_at DESC LIMIT 1), '') as last_message,
			       COALESCE((SELECT created_at FROM messages WHERE conversation_id = c.id ORDER BY created_at DESC LIMIT 1), 0) as last_message_time,
			       COALESCE((SELECT COUNT(*) FROM messages WHERE conversation_id = c.id AND sender_id != ? AND is_read = 0), 0) as unread_count
			FROM conversation_participants cp
			JOIN conversations c ON c.id = cp.conversation_id
			JOIN conversation_participants cp2 ON cp2.conversation_id = c.id AND cp2.user_id != cp.user_id
			JOIN users u ON u.id = cp2.user_id
			WHERE cp.user_id = ?
			ORDER BY last_message_time DESC, c.created_at DESC`,
			userID, userID,
		)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var convs []ConversationResponse
		for rows.Next() {
			var c ConversationResponse
			err := rows.Scan(&c.ID, &c.OtherUserID, &c.OtherUsername, &c.AvatarURL, &c.LastMessage, &c.LastMessageAt, &c.UnreadCount)
			if err != nil {
				continue
			}
			convs = append(convs, c)
		}

		if convs == nil {
			convs = []ConversationResponse{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(convs)
	}
}

func GetMessagesHandler(db *sql.DB) http.HandlerFunc {
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

		convID := r.PathValue("id")
		if convID == "" {
			http.Error(w, "missing conversation id", http.StatusBadRequest)
			return
		}

		// Verify membership in conversation
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM conversation_participants WHERE conversation_id = ? AND user_id = ?", convID, userID).Scan(&count)
		if err != nil || count == 0 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		_, _ = db.Exec("UPDATE messages SET is_read = 1 WHERE conversation_id = ? AND sender_id != ?", convID, userID)

		rows, err := db.Query(`
			SELECT id, conversation_id, sender_id, content, created_at, is_read
			FROM messages
			WHERE conversation_id = ?
			ORDER BY created_at ASC`,
			convID,
		)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var msgs []MessageResponse
		for rows.Next() {
			var m MessageResponse
			var isRead int
			err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.CreatedAt, &isRead)
			if err != nil {
				continue
			}
			m.IsRead = isRead == 1
			msgs = append(msgs, m)
		}

		if msgs == nil {
			msgs = []MessageResponse{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msgs)
	}
}

func GetOrCreateConversationMessagesHandler(db *sql.DB) http.HandlerFunc {
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

		otherUserID := r.PathValue("userID")
		if otherUserID == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}

		// Verify there is a follow relationship in at least one direction
		var followCount int
		err := db.QueryRow(`
			SELECT COUNT(*) FROM followers 
			WHERE (follower_id = ? AND following_id = ?) 
			   OR (follower_id = ? AND following_id = ?)`,
			userID, otherUserID, otherUserID, userID,
		).Scan(&followCount)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		if followCount == 0 {
			http.Error(w, "cannot chat with non-followers", http.StatusForbidden)
			return
		}

		convID, err := repositories.GetOrCreateConversation(userID, otherUserID)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		_, _ = db.Exec("UPDATE messages SET is_read = 1 WHERE conversation_id = ? AND sender_id != ?", convID, userID)

		rows, err := db.Query(`
			SELECT id, conversation_id, sender_id, content, created_at, is_read
			FROM messages
			WHERE conversation_id = ?
			ORDER BY created_at ASC`,
			convID,
		)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var msgs []MessageResponse
		for rows.Next() {
			var m MessageResponse
			var isRead int
			err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.CreatedAt, &isRead)
			if err != nil {
				continue
			}
			m.IsRead = isRead == 1
			msgs = append(msgs, m)
		}

		if msgs == nil {
			msgs = []MessageResponse{}
		}

		resp := map[string]interface{}{
			"conversation_id": convID,
			"messages":        msgs,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func ChatUploadHandler(db *sql.DB) http.HandlerFunc {
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

		path, err := SaveUploadedImage(r, "image", "uploads/chat")
		if err != nil {
			if errors.Is(err, ErrUnsupportedMediaType) {
				http.Error(w, "invalid file type", http.StatusUnsupportedMediaType)
				return
			}
			http.Error(w, "failed to upload image", http.StatusInternalServerError)
			return
		}

		if path == "" {
			http.Error(w, "no file uploaded", http.StatusBadRequest)
			return
		}

		urlPath := "/" + filepath.ToSlash(path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"url": urlPath,
		})
	}
}

func GetGroupMessagesHandler(db *sql.DB) http.HandlerFunc {
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

		groupID := r.PathValue("groupID")
		if groupID == "" {
			http.Error(w, "missing group id", http.StatusBadRequest)
			return
		}

		// Verify membership
		isMember, err := repositories.IsUserInGroup(userID, groupID)
		if err != nil || !isMember {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		msgs, err := repositories.GetGroupMessages(groupID, 100, 0)
		if err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}

		type GroupMessageWithSender struct {
			ID        string `json:"id"`
			GroupID   string `json:"group_id"`
			SenderID  string `json:"sender_id"`
			Username  string `json:"username"`
			Content   string `json:"content"`
			CreatedAt int64  `json:"created_at"`
		}

		var resp []GroupMessageWithSender
		// reverse order so chronological
		for i := len(msgs) - 1; i >= 0; i-- {
			m := msgs[i]
			var username string
			_ = db.QueryRow("SELECT username FROM users WHERE id = ?", m.SenderID).Scan(&username)
			resp = append(resp, GroupMessageWithSender{
				ID:        m.ID,
				GroupID:   m.GroupID,
				SenderID:  m.SenderID,
				Username:  username,
				Content:   m.Content,
				CreatedAt: m.CreatedAt,
			})
		}

		if resp == nil {
			resp = []GroupMessageWithSender{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
