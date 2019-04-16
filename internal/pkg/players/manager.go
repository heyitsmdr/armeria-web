package players

import (
	schemaGame "armeria/internal/pkg/game/schema"
	"armeria/internal/pkg/players/schema"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type manager struct {
	gameState schemaGame.IGameState
	players map[*player]bool
	mux sync.Mutex
}

// Init creates a new player Manager instance
func Init(gs schemaGame.IGameState) *manager {
	return &manager{
		gameState: gs,
		players: make(map[*player]bool),
	}
}

// NewPlayer creates a new Player instance and returns it
func (m *manager) NewPlayer(conn *websocket.Conn) schema.IPlayer {
	p := &player{
		gameState: 		  m.gameState,
		socket:           conn,
		pumpsInitialized: false,
		sendData:         make(chan *outgoingDataStructure, 256),
	}

	p.clientAction = newClientAction(p)

	m.players[p] = true

	log.Printf("[players] new player connected from %s (%d total players)", conn.RemoteAddr().String(), len(m.players))

	return p
}

// DisconnectPlayer will gracefully remove the player from the game and terminate the socket connection
func (m *manager) DisconnectPlayer(pp schema.IPlayer) {
	m.mux.Lock()
	defer m.mux.Unlock()

	p := pp.(*player)

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