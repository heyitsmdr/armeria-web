package sockets

import (
	"github.com/googollee/go-socket.io"
	"log"
)

// Server is the socket.io server instance
var Server *socketio.Server

// Init will initialize the socket.io server
func Init() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal("[sockets] socketio.NewServer: ", err)
	}

	Server = server

	Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		return nil
	})
}
