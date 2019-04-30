package armeria

import (
	"fmt"
	"strconv"
	"strings"
)

type Command struct {
	Name         string
	Help         string
	Alias        string
	Permissions  *CommandPermissions
	AllowedRoles []int
	Arguments    []*CommandArgument
	Subcommands  []*Command
	Handler      func(r *CommandContext)
}

type CommandArgument struct {
	Position         int
	Name             string
	IncludeRemaining bool
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
	Args      map[string]string
}

// CheckPermissions returns whether or not a player can use the command
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

// GetSubcommands returns the list of sub-commands that the player has access to as a string
func (cmd *Command) GetSubcommands(p *Player) string {
	if len(cmd.Subcommands) == 0 {
		return ""
	}

	var allowedSubCommands []*Command
	var longestCommandSize int
	for _, scmd := range cmd.Subcommands {
		if cmd.CheckPermissions(p) {
			if len(scmd.Name) > longestCommandSize {
				longestCommandSize = len(scmd.Name)
			}
			allowedSubCommands = append(allowedSubCommands, scmd)
		}
	}

	// NOTE: the "7" is being added to the padding to compensate for [b] and [/b]
	var output []string
	for _, scmd := range allowedSubCommands {
		output = append(output, fmt.Sprintf(
			"  %-"+strconv.Itoa(longestCommandSize+7+1)+"v %s",
			"[b]"+scmd.Name+"[/b]",
			scmd.Help,
		))
	}

	return strings.Join(output, "\n")
}
