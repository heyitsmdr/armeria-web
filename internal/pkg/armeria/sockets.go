package armeria

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWs upgrades the connection to a WebSocket
func ServeWs(state *GameState, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[sockets] ServeWs: %s", err)
	}

	p := state.PlayerManager().NewPlayer(conn)
	p.SetupPumps()
	p.ShowConnectionText()
}
