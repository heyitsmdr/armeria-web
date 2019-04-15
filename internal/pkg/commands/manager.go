package commands

import (
	"armeria/internal/pkg/players"
	"fmt"
	"log"
	"strings"
	"sync"
)

// Manager is the global manager instance for Command objects
type Manager struct {
	commands []*Command
	mux sync.Mutex
}

type Command struct {
	Name    string
	Handler func(p *players.Player)
}

func Init() *Manager {
	return &Manager{
		commands: []*Command{},
	}
}

func (m *Manager) RegisterCommand(c *Command) {
	m.commands = append(m.commands, c)
	log.Printf("[commands] command registered: %s", fmt.Sprintf("/%s", c.Name))
}

func (m *Manager) ProcessCommand(p *players.Player, cmd string) {
	sections := strings.Fields(cmd)
	if len(sections) == 0 {
		return
	}

	commandName := sections[0]

	for _, c := range m.commands {
		if c.Name == commandName {
			c.Handler(p)
			break
		}
	}

	p.ClientAction.ShowText("That's an invalid command.")
}