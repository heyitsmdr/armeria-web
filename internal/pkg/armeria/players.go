package armeria

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerManager struct {
	gameState *GameState
	players   map[*Player]bool
	mux       sync.Mutex
}

type Player struct {
	gameState        *GameState
	clientActions    *ClientActions
	socket           *websocket.Conn
	pumpsInitialized bool
	sendData         chan *OutgoingDataStructure
}

type IncomingDataStructure struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type OutgoingDataStructure struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"data"`
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

func (p *Player) readPump() {
	defer p.gameState.playerManager.DisconnectPlayer(p)

	p.socket.SetReadLimit(512)

	for {
		messageRead := &IncomingDataStructure{}
		err := p.socket.ReadJSON(messageRead)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[players] error reading from socket: %s", err)
			}
			break
		}

		switch messageRead.Type {
		case "command":
			p.gameState.commandManager.ProcessCommand(p, messageRead.Payload.(string))
		default:
			p.clientActions.ShowText("Your client sent invalid data.")
		}
	}
}

func (p *Player) writePump() {
	defer p.gameState.playerManager.DisconnectPlayer(p)

	for {
		select {
		case message, channelOpen := <-p.sendData:
			// Has the sendData chan been closed?
			if !channelOpen {
				return
			}

			err := p.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				log.Printf("[players] error setting write deadline: %s", err)
				return
			}

			if err := p.socket.WriteJSON(message); err != nil {
				log.Printf("[players] error writing to socket: %s", err)
				return
			}
		}
	}
}

// SetupPumps will create two go routines for reading and writing from the socket
func (p *Player) SetupPumps() {
	if p.pumpsInitialized {
		log.Printf("[players] call to SetupPumps failed (pumps already set up)")
		return
	}

	go p.readPump()
	go p.writePump()
	p.pumpsInitialized = true
}

// CallClientAction sends a socket event to call a Vuex action on the webapp
func (p *Player) CallClientAction(actionName string, payload interface{}) {
	p.sendData <- &OutgoingDataStructure{Action: actionName, Payload: payload}
}

func (p *Player) ShowConnectionText() {
	lines := []string{
		"Welcome to the world of Armeria!\n",
		"[ARMERIA ASCII ART HERE]\n",
		"If you have a character, you can <b>/login</b>. Otherwise, use <b>/create</b>.",
	}

	p.clientActions.ShowText(strings.Join(lines, "<br>"))
}
