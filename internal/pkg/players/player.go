package players

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type socketDataStructure struct {
	messageType string
	messageData string
}

// Player is a connected player with a valid socket connection
type Player struct {
	socket           *websocket.Conn
	pumpsInitialized bool
	sendData         chan socketDataStructure
}

func (p *Player) readPump() {
	defer Manager.DisconnectPlayer(p)
	p.socket.SetReadLimit(512)
	for {
		_, message, err := p.socket.ReadJSON(socketDataStructure)
	}
}

func (p *Player) writePump() {
	defer Manager.DisconnectPlayer(p)

	for {
		select {
		case message, channelOpen := <-p.sendData:
			err := p.socket.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err != nil {
				log.Printf("[players] error setting write deadline: %s", err)
				return
			}

			// Has the sendData chan been closed?
			if !channelOpen {
				err := p.socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("[players] error writing close message to socket: %s", err)
				}
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
