package armeria

import (
	"fmt"
	"log"
	"strings"
)

// Manager is the global manager instance for Command objects
type CommandManager struct {
	gameState *GameState
	commands  []*Command
}

type Command struct {
	Name       string
	SyntaxHelp string
	Handler    func(r *CommandRequest)
}

type CommandRequest struct {
	Command *Command
	Player  *Player
	Args    []string
}

func NewCommandManager(state *GameState) *CommandManager {
	return &CommandManager{
		gameState: state,
		commands:  []*Command{},
	}
}

func (m *CommandManager) RegisterCommand(c *Command) {
	m.commands = append(m.commands, c)
	log.Printf("[commands] command registered: %s", fmt.Sprintf("/%s", c.Name))
}

func (m *CommandManager) ProcessCommand(p *Player, cmd string) {
	sections := strings.Fields(cmd)
	if len(sections) == 0 {
		return
	}

	// Get command name and trim the first character
	commandName := sections[0][1:]

	for _, c := range m.commands {
		if c.Name == commandName {
			c.Handler(&CommandRequest{
				Command: c,
				Player:  p,
				Args:    sections[1:],
			})
			return
		}
	}

	p.clientActions.ShowText("That's an invalid command.")
}
