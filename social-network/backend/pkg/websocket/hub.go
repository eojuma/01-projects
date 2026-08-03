package websocket

import (
	"encoding/json"
	"log"
	"social-network/backend/pkg/repositories"
)

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Handlers   map[string]func(*Client, WebSocketMessage)
}

var ActiveHub *Hub

func NewHub() *Hub {
	h := &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Handlers:   make(map[string]func(*Client, WebSocketMessage)),
	}

	h.RegisterHandler("chat", h.handlePrivateChat)
	h.RegisterHandler("group_chat", h.handleGroupChat)
	h.RegisterHandler("typing", h.handleTyping)

	ActiveHub = h
	return h
}

func (h *Hub) RegisterHandler(t string, fn func(*Client, WebSocketMessage)) {
	h.Handlers[t] = fn
}

func (h *Hub) HandleMessage(c *Client, msg WebSocketMessage) {
	if fn, ok := h.Handlers[msg.Type]; ok {
		fn(c, msg)
	} else {
		log.Printf("unknown message type: %s", msg.Type)
	}
}

func (h *Hub) handlePrivateChat(c *Client, msg WebSocketMessage) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}
	content, _ := payload["content"].(string)
	if content == "" {
		return
	}

	// Verify follow relationship in at least one direction
	var followCount int
	err := repositories.DB.QueryRow(`
		SELECT COUNT(*) FROM followers 
		WHERE (follower_id = ? AND following_id = ?) 
		   OR (follower_id = ? AND following_id = ?)`,
		c.ID, msg.To, msg.To, c.ID,
	).Scan(&followCount)
	if err != nil || followCount == 0 {
		log.Printf("unauthorized chat message from %s to %s", c.ID, msg.To)
		return
	}

	convID, err := repositories.GetOrCreateConversation(c.ID, msg.To)
	if err != nil {
		log.Printf("failed to get conversation: %v", err)
		return
	}

	savedMsg, err := repositories.SaveMessage(convID, c.ID, content)
	if err != nil {
		log.Printf("failed to save message: %v", err)
		return
	}

	resp := map[string]interface{}{
		"type": "chat",
		"payload": map[string]interface{}{
			"id":              savedMsg.ID,
			"from":            c.ID,
			"content":         savedMsg.Content,
			"timestamp":       savedMsg.CreatedAt,
			"conversation_id": convID,
		},
	}
	data, _ := json.Marshal(resp)

	if recipient, ok := h.Clients[msg.To]; ok {
		recipient.Send <- data
	}

	c.Send <- data
}

func (h *Hub) handleGroupChat(c *Client, msg WebSocketMessage) {
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		return
	}
	content, _ := payload["content"].(string)
	if content == "" {
		return
	}

	groupID := msg.To
	isMember, err := repositories.IsUserInGroup(c.ID, groupID)
	if err != nil || !isMember {
		return
	}

	savedMsg, err := repositories.SaveGroupMessage(groupID, c.ID, content)
	if err != nil {
		log.Printf("failed to save group message: %v", err)
		return
	}

	var senderUsername string
	_ = repositories.DB.QueryRow("SELECT username FROM users WHERE id = ?", c.ID).Scan(&senderUsername)

	members, err := repositories.GetGroupMembers(groupID)
	if err != nil {
		return
	}

	resp := map[string]interface{}{
		"type": "group_chat",
		"payload": map[string]interface{}{
			"id":        savedMsg.ID,
			"group_id":  groupID,
			"from":      c.ID,
			"username":  senderUsername,
			"content":   savedMsg.Content,
			"timestamp": savedMsg.CreatedAt,
		},
	}
	data, _ := json.Marshal(resp)

	for _, memberID := range members {
		if memberID == c.ID {
			continue
		}
		if client, ok := h.Clients[memberID]; ok {
			client.Send <- data
		}
	}
}

func (h *Hub) handleTyping(c *Client, msg WebSocketMessage) {
	resp := map[string]interface{}{
		"type": "typing",
		"payload": map[string]interface{}{
			"from": c.ID,
			"to":   msg.To,
		},
	}
	data, _ := json.Marshal(resp)
	if recipient, ok := h.Clients[msg.To]; ok {
		recipient.Send <- data
	}
}

func (h *Hub) SendToUser(userID string, payload interface{}) {
	if client, ok := h.Clients[userID]; ok {
		msg := map[string]interface{}{
			"type":    "notification",
			"payload": payload,
		}
		data, _ := json.Marshal(msg)
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(h.Clients, client.ID)
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Clients[c.ID] = c
			log.Printf("client %s connected", c.ID)

		case c := <-h.Unregister:
			if _, ok := h.Clients[c.ID]; ok {
				delete(h.Clients, c.ID)
				close(c.Send)
				log.Printf("client %s disconnected", c.ID)
			}

		case msg := <-h.Broadcast:
			for _, c := range h.Clients {
				select {
				case c.Send <- msg:
				default:
					close(c.Send)
					delete(h.Clients, c.ID)
				}
			}
		}
	}
}
