package schema

import schemaPlayers "armeria/internal/pkg/players/schema"

type Command struct {
	Name    string
	Handler func(p schemaPlayers.IPlayer)
}

type ICommandManager interface {
	RegisterCommand(c Command)
	ProcessCommand(p schemaPlayers.IPlayer, cmd string)
}