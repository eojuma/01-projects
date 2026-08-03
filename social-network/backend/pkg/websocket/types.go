package websocket

type WebSocketMessage struct {
    Type    string      `json:"type"`
    To      string      `json:"to"`
    Payload interface{} `json:"payload"`
}