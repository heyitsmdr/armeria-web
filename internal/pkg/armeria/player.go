package armeria

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Player struct {
	gameState        *GameState
	clientActions    *ClientActions
	socket           *websocket.Conn
	pumpsInitialized bool
	sendData         chan *OutgoingDataStructure
	character        *Character
	mux              sync.Mutex
}

type IncomingDataStructure struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type OutgoingDataStructure struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"data"`
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
			cmd := messageRead.Payload.(string)
			p.gameState.commandManager.ProcessCommand(p, cmd[1:])
		case "objectEditorOpen":
			p.GetCharacter().SetTempAttribute("editorOpen", "true")
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

func (p *Player) FlushWrites() {
	for len(p.sendData) > 0 {
		data := <-p.sendData
		err := p.socket.WriteJSON(data)
		if err != nil {
			log.Printf("[player] error flushing writes: %s", err)
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
		"If you have a character, you can <b>/login</b>. Otherwise, use <b>/create</b>.",
	}

	p.clientActions.ShowRawText(strings.Join(lines, "\n"))
}

func (p *Player) AttachCharacter(c *Character) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.character = c
}

func (p *Player) GetCharacter() *Character {
	p.mux.Lock()
	defer p.mux.Unlock()
	return p.character
}
