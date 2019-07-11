package armeria

import (
	"sync"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

type PlayerManager struct {
	sync.RWMutex
	players map[*Player]bool
}

// Init creates a new player Manager instance
func NewPlayerManager() *PlayerManager {
	return &PlayerManager{
		players: make(map[*Player]bool),
	}
}

// NewPlayer creates a new Player instance, adds it to memory, and returns Player.
func (m *PlayerManager) NewPlayer(conn *websocket.Conn) *Player {
	m.Lock()
	defer m.Unlock()

	p := &Player{
		socket:           conn,
		pumpsInitialized: false,
		sendData:         make(chan *OutgoingDataStructure, 256),
	}

	p.clientActions = NewClientActions(p)

	m.players[p] = true

	Armeria.log.Info("player connected",
		zap.String("ip", conn.RemoteAddr().String()),
		zap.Int("players", len(m.players)),
	)

	return p
}

// DisconnectPlayer will gracefully remove the player from the game and terminate the socket connection
func (m *PlayerManager) DisconnectPlayer(p *Player) {
	m.Lock()
	defer m.Unlock()

	if !m.players[p] {
		return
	}

	if p.character != nil {
		// Notify character of logout
		p.character.LoggedOut()
		// Unset player from character
		p.character.SetPlayer(nil)
		// Unset character from player
		p.AttachCharacter(nil)
	}

	// Fatal if data should of been sent but wasn't
	if len(p.sendData) > 0 {
		Armeria.log.Error("player disconnected with unsent data",
			zap.Int("dataSize", len(p.sendData)),
		)
	}

	// Close the socket connection
	err := p.socket.Close()
	if err != nil {
		Armeria.log.Error("error closing socket",
			zap.Error(err),
		)
	}

	// Close the player's write channel
	close(p.sendData)

	// Remove the player from the manager
	delete(m.players, p)

	Armeria.log.Info("player disconnected",
		zap.Int("players", len(m.players)),
	)
}
