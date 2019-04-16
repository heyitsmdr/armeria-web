package characters

import (
	schemaCommands "armeria/internal/pkg/commands/schema"
	schemaPlayers "armeria/internal/pkg/players/schema"
)

func (m *manager) registerCommands() {
	m.gameState.CommandManager().RegisterCommand(schemaCommands.Command{
		Name: "login",
		Handler: login,
	})
}

func login(p schemaPlayers.IPlayer) {
	p.ClientActions().ShowText("Trying to login? Okay!")
}