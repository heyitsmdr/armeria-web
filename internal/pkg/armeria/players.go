package armeria

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type PlayerManager struct {
	gameState *GameState
	players   map[*Player]bool
	mux       sync.Mutex
}

// Init creates a new player Manager instance
func NewPlayerManager(state *GameState) *PlayerManager {
	return &PlayerManager{
		gameState: state,
		players:   make(map[*Player]bool),
	}
}

// NewPlayer creates a new Player instance and returns it
func (m *PlayerManager) NewPlayer(conn *websocket.Conn) *Player {
	p := &Player{
		gameState:        m.gameState,
		socket:           conn,
		pumpsInitialized: false,
		sendData:         make(chan *OutgoingDataStructure, 256),
	}

	p.clientActions = NewClientActions(p)

	m.players[p] = true

	log.Printf("[players] new player connected from %s (%d total players)", conn.RemoteAddr().String(), len(m.players))

	return p
}

// DisconnectPlayer will gracefully remove the player from the game and terminate the socket connection
func (m *PlayerManager) DisconnectPlayer(p *Player) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if !m.players[p] {
		return
	}

	if p.character != nil {
		// Notify character of logout
		p.character.LoggedOut(m.gameState)
		// Unset player from character
		p.character.SetPlayer(nil)
		// Log
		log.Printf("[players] character logged out: %s", p.character.GetName())
		// Unset character from player
		p.AttachCharacter(nil)
	}

	// Fatal if data should of been sent but wasn't
	if len(p.sendData) > 0 {
		log.Fatal("[players] player disconnected with unsent data")
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
