package schema

import (
	schemaPlayers "armeria/internal/pkg/players/schema"
)

type Command struct {
	Name       string
	SyntaxHelp string
	Handler    func(r *CommandRequest)
}

type CommandRequest struct {
	Command *Command
	Player  schemaPlayers.IPlayer
	Args    []string
}

type ICommandManager interface {
	RegisterCommand(cmd Command)
	ProcessCommand(player schemaPlayers.IPlayer, cmd string)
}
