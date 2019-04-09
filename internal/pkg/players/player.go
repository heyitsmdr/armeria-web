package players

import "github.com/gorilla/websocket"

type Player struct {
	socket *websocket.Conn
}
