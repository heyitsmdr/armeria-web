package players

import (
	"github.com/gorilla/websocket"
	"log"
)

// Manager is the global manager instance for Player objects
var Manager *manager

type manager struct {
	players map[*Player]bool
}

// Init creates a new player Manager instance
func Init() {
	Manager = &manager{
		players: make(map[*Player]bool),
	}
}

// NewPlayer creates a new Player instance and returns it
func (m *manager) NewPlayer(conn *websocket.Conn) *Player {
	p := &Player{
		socket:           conn,
		pumpsInitialized: false,
		sendData:         make(chan []byte, 256),
	}

	log.Printf("[players] new player connected from %s", conn.RemoteAddr().String())

	m.players[p] = true

	return p
}
