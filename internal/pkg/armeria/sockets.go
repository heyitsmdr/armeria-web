package armeria

import (
	"net/http"

	"go.uber.org/zap"

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
func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Armeria.log.Error("error upgrading socket connection",
			zap.Error(err),
		)
		return
	}

	p := Armeria.playerManager.NewPlayer(conn)
	p.SetupPumps()
	p.ShowConnectionText()
}
