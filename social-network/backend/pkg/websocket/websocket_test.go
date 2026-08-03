package websocket

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"social-network/backend/pkg/repositories"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

var originalHub *Hub

func init() {
	originalHub = ActiveHub
}

// setupTestDB initializes an in-memory SQLite database with required schemas.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	schemas := []string{
		`CREATE TABLE users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL
		);`,
		`CREATE TABLE followers (
			follower_id TEXT,
			following_id TEXT,
			PRIMARY KEY (follower_id, following_id)
		);`,
		`CREATE TABLE conversations (
			id TEXT PRIMARY KEY,
			created_at INTEGER NOT NULL
		);`,
		`CREATE TABLE conversation_participants (
			conversation_id TEXT,
			user_id TEXT,
			PRIMARY KEY (conversation_id, user_id)
		);`,
		`CREATE TABLE messages (
			id TEXT PRIMARY KEY,
			conversation_id TEXT,
			sender_id TEXT,
			content TEXT,
			created_at INTEGER,
			is_read BOOLEAN
		);`,
		`CREATE TABLE group_members (
			group_id TEXT,
			user_id TEXT,
			PRIMARY KEY (group_id, user_id)
		);`,
		`CREATE TABLE group_messages (
			id TEXT PRIMARY KEY,
			group_id TEXT,
			sender_id TEXT,
			content TEXT,
			created_at INTEGER
		);`,
	}

	for _, s := range schemas {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("failed to execute schema: %v", err)
		}
	}

	repositories.SetDatabase(db)
	return db
}

func clearActiveHub() {
	if originalHub != nil {
		ActiveHub = originalHub
	}
	if ActiveHub != nil {
		for k := range ActiveHub.Clients {
			delete(ActiveHub.Clients, k)
		}
	}
}

func waitClientsRegistered(t *testing.T, userIDs ...string) {
	for _, id := range userIDs {
		registered := false
		for i := 0; i < 100; i++ {
			if ActiveHub != nil {
				if _, ok := ActiveHub.Clients[id]; ok {
					registered = true
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		if !registered {
			t.Fatalf("timeout waiting for client %s to register", id)
		}
	}
}

func TestNewHub(t *testing.T) {
	defer func() {
		ActiveHub = originalHub
	}()

	h := NewHub()
	if h == nil {
		t.Fatal("expected non-nil Hub")
	}
	if h.Clients == nil || h.Register == nil || h.Unregister == nil || h.Broadcast == nil || h.Handlers == nil {
		t.Fatal("expected Hub fields to be initialized")
	}
	if _, ok := h.Handlers["chat"]; !ok {
		t.Error("expected chat handler to be registered")
	}
	if _, ok := h.Handlers["group_chat"]; !ok {
		t.Error("expected group_chat handler to be registered")
	}
	if _, ok := h.Handlers["typing"]; !ok {
		t.Error("expected typing handler to be registered")
	}
}

func TestWebSocketPrivateChat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	clearActiveHub()

	// Seed users
	_, err := db.Exec("INSERT INTO users (id, username) VALUES ('userA', 'Alice'), ('userB', 'Bob')")
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	// Seed follow relationships (Alice follows Bob, Bob follows Alice)
	_, err = db.Exec("INSERT INTO followers (follower_id, following_id) VALUES ('userA', 'userB'), ('userB', 'userA')")
	if err != nil {
		t.Fatalf("failed to seed followers: %v", err)
	}

	// Serve handler
	testWSHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		ctx := context.WithValue(r.Context(), "userID", userID)
		ServeWebSocket(w, r.WithContext(ctx))
	})

	server := httptest.NewServer(testWSHandler)
	defer server.Close()

	// Dial clients
	urlA := strings.Replace(server.URL, "http", "ws", 1) + "?userId=userA"
	urlB := strings.Replace(server.URL, "http", "ws", 1) + "?userId=userB"

	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	connA, _, err := dialer.Dial(urlA, nil)
	if err != nil {
		t.Fatalf("failed to connect client A: %v", err)
	}
	defer connA.Close()

	connB, _, err := dialer.Dial(urlB, nil)
	if err != nil {
		t.Fatalf("failed to connect client B: %v", err)
	}
	defer connB.Close()

	// Wait for registration dynamically
	waitClientsRegistered(t, "userA", "userB")

	// Send message from A to B
	msg := WebSocketMessage{
		Type: "chat",
		To:   "userB",
		Payload: map[string]interface{}{
			"content": "Hello, Bob!",
		},
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	if err := connA.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		t.Fatalf("failed to write message from A: %v", err)
	}

	// Assert client A receives confirmation/echo
	_, respA, err := connA.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read response client A: %v", err)
	}

	var parsedRespA map[string]interface{}
	if err := json.Unmarshal(respA, &parsedRespA); err != nil {
		t.Fatalf("failed to unmarshal response A: %v", err)
	}

	if parsedRespA["type"] != "chat" {
		t.Errorf("expected type chat, got %v", parsedRespA["type"])
	}

	// Assert client B receives the message
	_, respB, err := connB.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message client B: %v", err)
	}

	var parsedRespB map[string]interface{}
	if err := json.Unmarshal(respB, &parsedRespB); err != nil {
		t.Fatalf("failed to unmarshal response B: %v", err)
	}

	payloadB, ok := parsedRespB["payload"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected payload map, got %v", parsedRespB["payload"])
	}

	if payloadB["content"] != "Hello, Bob!" {
		t.Errorf("expected content 'Hello, Bob!', got %v", payloadB["content"])
	}

	if payloadB["from"] != "userA" {
		t.Errorf("expected from 'userA', got %v", payloadB["from"])
	}

	// Assert saved in database
	var dbCount int
	err = db.QueryRow("SELECT COUNT(*) FROM messages WHERE sender_id = 'userA' AND content = 'Hello, Bob!'").Scan(&dbCount)
	if err != nil {
		t.Fatalf("failed to query database for saved message: %v", err)
	}
	if dbCount != 1 {
		t.Errorf("expected message to be saved in DB, got count %d", dbCount)
	}
}

func TestWebSocketGroupChat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	clearActiveHub()

	// Seed users
	_, err := db.Exec(`
		INSERT INTO users (id, username) VALUES 
		('userA', 'Alice'), 
		('userB', 'Bob'), 
		('userC', 'Charlie')`)
	if err != nil {
		t.Fatalf("failed to seed users: %v", err)
	}

	// Seed group membership
	_, err = db.Exec(`
		INSERT INTO group_members (group_id, user_id) VALUES 
		('group1', 'userA'), 
		('group1', 'userB'), 
		('group1', 'userC')`)
	if err != nil {
		t.Fatalf("failed to seed group members: %v", err)
	}

	// Serve handler
	testWSHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		ctx := context.WithValue(r.Context(), "userID", userID)
		ServeWebSocket(w, r.WithContext(ctx))
	})

	server := httptest.NewServer(testWSHandler)
	defer server.Close()

	// Dial clients
	urlA := strings.Replace(server.URL, "http", "ws", 1) + "?userId=userA"
	urlB := strings.Replace(server.URL, "http", "ws", 1) + "?userId=userB"
	urlC := strings.Replace(server.URL, "http", "ws", 1) + "?userId=userC"

	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	connA, _, err := dialer.Dial(urlA, nil)
	if err != nil {
		t.Fatalf("failed to connect client A: %v", err)
	}
	defer connA.Close()

	connB, _, err := dialer.Dial(urlB, nil)
	if err != nil {
		t.Fatalf("failed to connect client B: %v", err)
	}
	defer connB.Close()

	connC, _, err := dialer.Dial(urlC, nil)
	if err != nil {
		t.Fatalf("failed to connect client C: %v", err)
	}
	defer connC.Close()

	// Wait for registration dynamically
	waitClientsRegistered(t, "userA", "userB", "userC")

	// Send message from A to Group1
	msg := WebSocketMessage{
		Type: "group_chat",
		To:   "group1",
		Payload: map[string]interface{}{
			"content": "Hello, Group!",
		},
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	if err := connA.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		t.Fatalf("failed to write message from A: %v", err)
	}

	// Assert client B and client C receive it
	checkGroupResponse := func(conn *websocket.Conn) {
		_, resp, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("failed to read group message: %v", err)
		}

		var parsed map[string]interface{}
		if err := json.Unmarshal(resp, &parsed); err != nil {
			t.Fatalf("failed to unmarshal group response: %v", err)
		}

		if parsed["type"] != "group_chat" {
			t.Errorf("expected type group_chat, got %v", parsed["type"])
		}

		payload, ok := parsed["payload"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected payload map, got %v", parsed["payload"])
		}

		if payload["content"] != "Hello, Group!" {
			t.Errorf("expected content 'Hello, Group!', got %v", payload["content"])
		}

		if payload["group_id"] != "group1" {
			t.Errorf("expected group_id 'group1', got %v", payload["group_id"])
		}
	}

	checkGroupResponse(connB)
	checkGroupResponse(connC)

	// Assert saved in database
	var dbCount int
	err = db.QueryRow("SELECT COUNT(*) FROM group_messages WHERE group_id = 'group1' AND sender_id = 'userA' AND content = 'Hello, Group!'").Scan(&dbCount)
	if err != nil {
		t.Fatalf("failed to query database for saved group message: %v", err)
	}
	if dbCount != 1 {
		t.Errorf("expected group message to be saved in DB, got count %d", dbCount)
	}
}

func TestWebSocketTyping(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	clearActiveHub()

	testWSHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		ctx := context.WithValue(r.Context(), "userID", userID)
		ServeWebSocket(w, r.WithContext(ctx))
	})

	server := httptest.NewServer(testWSHandler)
	defer server.Close()

	urlA := strings.Replace(server.URL, "http", "ws", 1) + "?userId=userA"
	urlB := strings.Replace(server.URL, "http", "ws", 1) + "?userId=userB"

	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	connA, _, err := dialer.Dial(urlA, nil)
	if err != nil {
		t.Fatalf("failed to connect client A: %v", err)
	}
	defer connA.Close()

	connB, _, err := dialer.Dial(urlB, nil)
	if err != nil {
		t.Fatalf("failed to connect client B: %v", err)
	}
	defer connB.Close()

	// Wait for registration dynamically
	waitClientsRegistered(t, "userA", "userB")

	// Send typing event from A to B
	msg := WebSocketMessage{
		Type: "typing",
		To:   "userB",
		Payload: map[string]interface{}{},
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal typing message: %v", err)
	}

	if err := connA.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		t.Fatalf("failed to write typing message from A: %v", err)
	}

	// Assert client B receives the typing event
	_, respB, err := connB.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read typing event client B: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(respB, &parsed); err != nil {
		t.Fatalf("failed to unmarshal typing response: %v", err)
	}

	if parsed["type"] != "typing" {
		t.Errorf("expected type typing, got %v", parsed["type"])
	}

	payload, ok := parsed["payload"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected payload map, got %v", parsed["payload"])
	}

	if payload["from"] != "userA" {
		t.Errorf("expected from 'userA', got %v", payload["from"])
	}
	if payload["to"] != "userB" {
		t.Errorf("expected to 'userB', got %v", payload["to"])
	}
}

func TestSendToUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	clearActiveHub()

	testWSHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userId")
		ctx := context.WithValue(r.Context(), "userID", userID)
		ServeWebSocket(w, r.WithContext(ctx))
	})

	server := httptest.NewServer(testWSHandler)
	defer server.Close()

	urlA := strings.Replace(server.URL, "http", "ws", 1) + "?userId=userA"

	dialer := websocket.Dialer{HandshakeTimeout: 5 * time.Second}
	connA, _, err := dialer.Dial(urlA, nil)
	if err != nil {
		t.Fatalf("failed to connect client A: %v", err)
	}
	defer connA.Close()

	// Wait for registration dynamically
	waitClientsRegistered(t, "userA")

	// Send to user A from hub
	ActiveHub.SendToUser("userA", map[string]string{"message": "Hello Alice!"})

	// Assert client A receives it
	_, respA, err := connA.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(respA, &parsed); err != nil {
		t.Fatalf("failed to unmarshal notification: %v", err)
	}

	if parsed["type"] != "notification" {
		t.Errorf("expected type notification, got %v", parsed["type"])
	}

	payload, ok := parsed["payload"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected payload map, got %v", parsed["payload"])
	}

	if payload["message"] != "Hello Alice!" {
		t.Errorf("expected message 'Hello Alice!', got %v", payload["message"])
	}
}
