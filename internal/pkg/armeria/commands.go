package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"log"
	"strings"
)

// Manager is the global manager instance for Command objects
type CommandManager struct {
	gameState *GameState
	commands  []*Command
}

// NewCommandManager will return a new instance of the command manager.
func NewCommandManager(state *GameState) *CommandManager {
	return &CommandManager{
		gameState: state,
		commands:  []*Command{},
	}
}

// RegisterCommand will register a Command with the command manager with the arguments
// parsed out.
func (m *CommandManager) RegisterCommand(c *Command) {
	m.commands = append(m.commands, c)
	log.Printf("[commands] command registered: %s", fmt.Sprintf("/%s", c.Name))
}

// FindCommand will return a matched registered Command.
func (m *CommandManager) FindCommand(p *Player, searchWithin []*Command, cmd string, alreadyProcessed []string) (*Command, map[string]string, string) {
	sections := strings.Fields(cmd)
	cmdName := strings.ToLower(sections[0])

	for _, cmd := range searchWithin {
		if strings.ToLower(cmd.Name) == cmdName || misc.Contains(cmd.AltNames, cmdName) {
			// Handle permissions
			if !cmd.CheckPermissions(p) {
				return nil, nil, "You cannot use that command right now."
			}

			// Handle sub-commands
			if cmd.Subcommands != nil {
				processedCommands := append(alreadyProcessed, cmdName)
				if len(sections) == 1 {
					return nil, nil, cmd.ShowSubcommandHelp(p, processedCommands)

				}
				return m.FindCommand(p, cmd.Subcommands, strings.Join(sections[1:], " "), processedCommands)
			}

			// Go through arguments
			commandArgs := make(map[string]string)
			if cmd.Arguments != nil {
				for pos, arg := range cmd.Arguments {
					if !arg.Optional && len(sections) < (pos+2) {
						return nil, nil, cmd.ShowArgumentHelp(p, append(alreadyProcessed, cmdName))
					}
					if arg.IncludeRemaining {
						commandArgs[arg.Name] = strings.Join(sections[pos+1:], " ")
					} else if len(sections) >= pos+2 {
						commandArgs[arg.Name] = sections[pos+1]
					} else {
						commandArgs[arg.Name] = ""
					}
				}
			}

			return cmd, commandArgs, ""
		}
	}

	return nil, nil, "That's an invalid command."
}

// ProcessCommand will evaluate and process a command sent by the player either
// manually or programmatically.
func (m *CommandManager) ProcessCommand(p *Player, command string) {
	sections := strings.Fields(command)
	if len(sections) == 0 {
		return
	}

	cmd, cmdArgs, errorMsg := m.FindCommand(p, m.commands, strings.Join(sections, " "), []string{})

	if cmd == nil {
		p.clientActions.ShowColorizedText(errorMsg, ColorError)
		return
	}

	ctx := &CommandContext{
		Command: cmd,
		Player:  p,
		Args:    cmdArgs,
	}

	if p.GetCharacter() != nil {
		ctx.Character = p.GetCharacter()
	}

	if len(cmd.Alias) > 0 {
		m.ProcessCommand(p, cmd.Alias)
		return
	}

	cmd.Handler(ctx)
}
