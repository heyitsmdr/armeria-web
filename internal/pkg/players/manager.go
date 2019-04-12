package players

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// Manager is the global manager instance for Player objects
var Manager *manager

type manager struct {
	players map[*Player]bool
	mux sync.Mutex
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
		sendData:         make(chan socketDataStructure, 256),
	}

	log.Printf("[players] new player connected from %s", conn.RemoteAddr().String())

	m.players[p] = true

	return p
}

// DisconnectPlayer will gracefully remove the player from the game and terminate the socket connection
func (m *manager) DisconnectPlayer(p *Player) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.players[p] {
		return
	}

	err := p.socket.Close()
	if err != nil {
		log.Printf("[players] error closing socket in DisconnectPlayer: %s", err)
	}

	delete(m.players, p)

	log.Printf("[players] player disconnected")
}