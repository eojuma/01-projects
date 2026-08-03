package models

type Conversation struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
}

type Message struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	CreatedAt      int64  `json:"created_at"`
	IsRead         bool   `json:"is_read"`
}

type GroupMessage struct {
	ID        string `json:"id"`
	GroupID   string `json:"group_id"`
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type WebSocketMessage struct {
	Type    string      `json:"type"`
	To      string      `json:"to"`
	Payload interface{} `json:"payload"`
}
