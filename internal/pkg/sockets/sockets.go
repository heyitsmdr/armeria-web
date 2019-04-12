package sockets

import (
	"armeria/internal/pkg/players"
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

// Init will initialize the socket.io server
func Init() {

}

// ServeWs upgrades the connection to a WebSocket
func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[sockets] ServeWs: %s", err)
	}

	p := players.Manager.NewPlayer(conn)
	p.SetupPumps()

}
