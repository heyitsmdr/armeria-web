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

// login ethryx xyrhte89
func (m *CommandManager) FindCommand(p *Player, searchWithin []*Command, cmd string, alreadyProcessed []string) (*Command, map[string]string, string) {
	sections := strings.Fields(cmd)
	cmdName := sections[0]

	for _, cmd := range searchWithin {
		if strings.ToLower(cmd.Name) == strings.ToLower(cmdName) {
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
					} else {
						commandArgs[arg.Name] = sections[pos+1]
					}
				}
			}

			return cmd, commandArgs, ""
		}
	}

	return nil, nil, "That's an invalid command."
}

func (m *CommandManager) ProcessCommand(p *Player, command string) {
	sections := strings.Fields(command)
	if len(sections) == 0 {
		return
	}

	cmd, cmdArgs, errorMsg := m.FindCommand(p, m.commands, strings.Join(sections, " "), []string{})

	if cmd == nil {
		p.clientActions.ShowText(errorMsg)
		return
	}

	ctx := &CommandContext{
		GameState: m.gameState,
		Command:   cmd,
		Player:    p,
		Args:      cmdArgs,
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
