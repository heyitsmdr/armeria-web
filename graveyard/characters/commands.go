package characters

import (
	schemaCommands "armeria/internal/pkg/commands/schema"
	"fmt"
)

func (m *manager) registerCommands() {
	m.gameState.CommandManager().RegisterCommand(schemaCommands.Command{
		Name:       "login",
		SyntaxHelp: "/login [character] [password]",
		Handler:    login,
	})
}

func login(r *schemaCommands.CommandRequest) {
	if len(r.Args) != 2 {
		r.Player.ClientActions().ShowText(fmt.Sprintf("[b]Syntax:[/b] %s", r.Command.SyntaxHelp))
		return
	}

	//character := r.Args[0]
	//password := r.Args[1]

	r.Player.ClientActions().ShowText("Trying to login!")
}
