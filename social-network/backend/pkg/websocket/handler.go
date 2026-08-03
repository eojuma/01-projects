package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var hub = NewHub()

func init() {
	go hub.Run()
}

func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade", http.StatusInternalServerError)
		return
	}

	client := &Client{
		ID:     userID,
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		IsOpen: true,
	}

	hub.Register <- client

	go client.WritePump()
	client.ReadPump()
}

var sessionStore interface{}

func SetSessionStore(store interface{}) {
	sessionStore = store
}
