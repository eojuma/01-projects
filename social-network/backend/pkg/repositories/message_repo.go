package repositories

import (
	"database/sql"
	"social-network/backend/pkg/models"
	"time"

	"github.com/gofrs/uuid"
)

var DB *sql.DB

func SetDatabase(db *sql.DB) {
	DB = db
}

func generateUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func now() int64 {
	return time.Now().Unix()
}

func CreateConversation(user1ID, user2ID string) (string, error) {
	convID := generateUUID()
	tx, err := DB.Begin()
	if err != nil {
		return "", err
	}

	_, err = tx.Exec(`INSERT INTO conversations (id, created_at) VALUES (?, ?)`, convID, now())
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec(`INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`, convID, user1ID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	_, err = tx.Exec(`INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?)`, convID, user2ID)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}
	return convID, nil
}

func GetOrCreateConversation(user1ID, user2ID string) (string, error) {
	var convID string
	query := `
        SELECT c.id FROM conversations c
        JOIN conversation_participants p1 ON c.id = p1.conversation_id
        JOIN conversation_participants p2 ON c.id = p2.conversation_id
        WHERE p1.user_id = ? AND p2.user_id = ?
    `
	err := DB.QueryRow(query, user1ID, user2ID).Scan(&convID)
	if err == sql.ErrNoRows {
		return CreateConversation(user1ID, user2ID)
	}
	return convID, err
}

func SaveMessage(convID, senderID, content string) (*models.Message, error) {
	msg := &models.Message{
		ID:             generateUUID(),
		ConversationID: convID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      now(),
		IsRead:         false,
	}
	_, err := DB.Exec(`
        INSERT INTO messages (id, conversation_id, sender_id, content, created_at, is_read)
        VALUES (?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ConversationID, msg.SenderID, msg.Content, msg.CreatedAt, msg.IsRead,
	)
	return msg, err
}

func GetMessages(convID string, limit, offset int) ([]models.Message, error) {
	rows, err := DB.Query(`
        SELECT id, conversation_id, sender_id, content, created_at, is_read
        FROM messages
        WHERE conversation_id = ?
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?`,
		convID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.SenderID, &m.Content, &m.CreatedAt, &m.IsRead); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
