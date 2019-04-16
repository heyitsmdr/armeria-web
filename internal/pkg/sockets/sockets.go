package sockets

import (
	"armeria/internal/pkg/game/schema"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs upgrades the connection to a WebSocket
func ServeWs(gs schema.IGameState, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[sockets] ServeWs: %s", err)
	}


	p := gs.PlayerManager().NewPlayer(conn)
	p.SetupPumps()

	p.ClientActions().ShowIntroText()
}
