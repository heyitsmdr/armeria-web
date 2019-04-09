package players

import (
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	socket           *websocket.Conn
	pumpsInitialized bool
	sendData         chan []byte
}

func (p *Player) readPump() {

}

func (p *Player) writePump() {
	for {
		select {
		case message, ok := <-p.sendData:
			if !ok {
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
