package players

import (
	"armeria/internal/pkg/game"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type incomingDataStructure struct {
	Type string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type outgoingDataStructure struct {
	Action   string 	 `json:"action"`
	Payload  interface{} `json:"data"`
}

// Player is a connected player with a valid socket connection
type Player struct {
	ClientAction	 *clientAction
	socket           *websocket.Conn
	pumpsInitialized bool
	sendData         chan *outgoingDataStructure
}

func (p *Player) readPump() {
	defer game.GameState.PlayerManager.DisconnectPlayer(p)

	p.socket.SetReadLimit(512)

	for {
		messageRead := &incomingDataStructure{}
		err := p.socket.ReadJSON(messageRead)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[players] error reading from socket: %s", err)
			}
			break
		}

		switch messageRead.Type {
		case "command":
			//commands.Manager.ProcessCommand(p, messageRead.Payload.(string))
		default:
			p.ClientAction.ShowText("Your client sent invalid data.")
		}
	}
}

func (p *Player) writePump() {
	defer game.GameState.PlayerManager.DisconnectPlayer(p)

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
	p.sendData <- &outgoingDataStructure{ Action: actionName, Payload: payload }
}