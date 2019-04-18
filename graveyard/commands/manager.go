package commands

import (
	"armeria/internal/pkg/commands/schema"
	schemaGame "armeria/internal/pkg/game/schema"
	schemaPlayers "armeria/internal/pkg/players/schema"
	"fmt"
	"log"
	"strings"
)

// Manager is the global manager instance for Command objects
type manager struct {
	gameState schemaGame.IGameState
	commands  []schema.Command
}

func Init(gs schemaGame.IGameState) *manager {
	return &manager{
		gameState: gs,
		commands:  []schema.Command{},
	}
}

func (m *manager) RegisterCommand(c schema.Command) {
	m.commands = append(m.commands, c)
	log.Printf("[commands] command registered: %s", fmt.Sprintf("/%s", c.Name))
}

func (m *manager) ProcessCommand(p schemaPlayers.IPlayer, cmd string) {
	sections := strings.Fields(cmd)
	if len(sections) == 0 {
		return
	}

	// Get command name and trim the first character
	commandName := sections[0][1:]

	for _, c := range m.commands {
		if c.Name == commandName {
			c.Handler(&schema.CommandRequest{
				Command: &c,
				Player:  p,
				Args:    sections[1:],
			})
			return
		}
	}

	p.ClientActions().ShowText("That's an invalid command.")
}
