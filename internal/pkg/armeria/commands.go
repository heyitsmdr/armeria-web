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
	Name         string
	Permissions  *CommandPermissions
	SyntaxHelp   string
	AllowedRoles []int
	Handler      func(r *CommandContext)
}

type CommandPermissions struct {
	RequireNoCharacter bool
	RequireCharacter   bool
}

type CommandContext struct {
	GameState *GameState
	Command   *Command
	Player    *Player
	Character *Character
	Args      []string
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

	for _, cmd := range m.commands {
		if strings.ToLower(cmd.Name) == strings.ToLower(commandName) {
			// Handle permissions
			if !cmd.CheckPermissions(p) {
				p.clientActions.ShowText("You cannot use that command right now.")
				return
			}

			ctx := &CommandContext{
				GameState: m.gameState,
				Command:   cmd,
				Player:    p,
				Args:      sections[1:],
			}

			if p.GetCharacter() != nil {
				ctx.Character = p.GetCharacter()
			}

			cmd.Handler(ctx)
			return
		}
	}

	p.clientActions.ShowText("That's an invalid command.")
}

func (cmd *Command) CheckPermissions(p *Player) bool {
	if cmd.Permissions == nil {
		return true
	}

	if cmd.Permissions.RequireNoCharacter {
		if p.GetCharacter() != nil {
			return false
		}
	}

	if cmd.Permissions.RequireCharacter {
		if p.GetCharacter() == nil {
			return false
		}
	}

	return true
}
