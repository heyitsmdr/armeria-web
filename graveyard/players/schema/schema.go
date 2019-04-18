package schema

import (
	"github.com/gorilla/websocket"
)

type IPlayerManager interface {
	NewPlayer(conn *websocket.Conn) IPlayer
	DisconnectPlayer(p IPlayer)
}

type IPlayer interface {
	SetupPumps()
	CallClientAction(actionName string, payload interface{})
	ClientActions() IClientAction
}

type IClientAction interface {
	ShowText(text string)
	ShowIntroText()
}
