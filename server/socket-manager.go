package arcadia

import "github.com/gorilla/websocket"

// SocketManager oversees all sockets connected to the game
type SocketManager struct {
}

var upgrader = websocket.Upgrader{}

// NewSocketManager creates a new SocketManager instance
func NewSocketManager() *SocketManager {

	return &SocketManager{}
}
