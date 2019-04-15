package characters

import (
	"armeria/internal/pkg/commands"
	"armeria/internal/pkg/players"
)

func registerCommands() {
	commands.Manager.RegisterCommand(&commands.Command{
		Name: "login",
		Handler: login,
	})
}

func login(p *players.Player) {
	p.ClientAction.ShowText("Trying to login? Okay!")
}