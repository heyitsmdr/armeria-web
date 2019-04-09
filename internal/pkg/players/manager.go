package players

import "github.com/gorilla/websocket"

var Manager *manager

type manager struct {
	players []*Player
}

func Init() {
	Manager = &manager{}
}

// NewPlayer creates a new Player instance
func (m *manager) NewPlayer(conn *websocket.Conn) *Player {
	p := &Player{
		socket: conn,
	}
	m.players
	return p
}
