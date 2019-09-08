package armeria

import (
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
)

type Player struct {
	sync.RWMutex
	client           ClientActions
	socket           *websocket.Conn
	pumpsInitialized bool
	sendData         chan *OutgoingDataStructure
	character        *Character
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
	defer Armeria.playerManager.DisconnectPlayer(p)

	// Set max size of a single message to 512KB
	p.socket.SetReadLimit(512000)

	for {
		messageRead := &IncomingDataStructure{}
		err := p.socket.ReadJSON(messageRead)
		if err != nil {
			Armeria.log.Debug("socket read error",
				zap.Error(err),
			)
			break
		}

		switch messageRead.Type {
		case "command":
			cmd := messageRead.Payload.(string)
			Armeria.commandManager.ProcessCommand(p, cmd[1:], true)
		case "objectEditorOpen":
			open := messageRead.Payload.(bool)
			if open {
				p.Character().SetTempAttribute("editorOpen", "true")
			} else {
				p.Character().SetTempAttribute("editorOpen", "false")
			}
		case "objectPictureUpload":
			StoreObjectPicture(p, messageRead.Payload.(map[string]interface{}))
		default:
			p.client.ShowText("Your client sent invalid data.")
		}
	}
}

func (p *Player) writePump() {
	defer Armeria.playerManager.DisconnectPlayer(p)

	for {
		select {
		case message, channelOpen := <-p.sendData:
			// Has the sendData chan been closed?
			if !channelOpen {
				return
			}

			err := p.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				Armeria.log.Error("error setting write deadline",
					zap.Error(err),
				)
				return
			}

			if err := p.socket.WriteJSON(message); err != nil {
				Armeria.log.Error("error writing to socket",
					zap.Error(err),
				)
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
			Armeria.log.Error("error flushing writes",
				zap.Error(err),
			)
		}
	}
}

// SetupPumps will create two go routines for reading and writing from the socket
func (p *Player) SetupPumps() {
	if p.pumpsInitialized {
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

// Connected is called when the parent successfully connects to the game (pre-login).
func (p *Player) Connected() {

}

func (p *Player) AttachCharacter(c *Character) {
	p.Lock()
	defer p.Unlock()

	p.character = c
}

func (p *Player) Character() *Character {
	p.RLock()
	defer p.RUnlock()

	return p.character
}

func (p *Player) PlayerInfoJSON() string {
	pi := map[string]string{
		"uuid": p.Character().ID(),
	}

	piJSON, err := json.Marshal(pi)
	if err != nil {
		Armeria.log.Fatal("failed to marshal player info data",
			zap.String("character", p.Character().UUID),
			zap.Error(err),
		)
	}

	return string(piJSON)
}
