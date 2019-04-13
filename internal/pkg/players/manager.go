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
		sendData:         make(chan *outgoingDataStructure, 256),
	}

	p.ClientAction = newClientAction(p)

	m.players[p] = true

	log.Printf("[players] new player connected from %s (%d total players)", conn.RemoteAddr().String(), len(m.players))

	return p
}

// DisconnectPlayer will gracefully remove the player from the game and terminate the socket connection
func (m *manager) DisconnectPlayer(p *Player) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.players[p] {
		return
	}

	// Close the socket connection
	err := p.socket.Close()
	if err != nil {
		log.Printf("[players] error closing socket in DisconnectPlayer: %s", err)
	}

	// Close the player's write channel
	close(p.sendData)

	// Remove the player from the manager
	delete(m.players, p)

	log.Printf("[players] player disconnected (%d total players)", len(m.players))
}