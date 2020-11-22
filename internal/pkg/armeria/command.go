package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Command struct {
	Parent       *Command                `json:"-"`
	Name         string                  `json:"name"`
	AltNames     []string                `json:"altNames"`
	Help         string                  `json:"help"`
	Hidden       bool                    `json:"-"`
	Alias        string                  `json:"alias"`
	Permissions  *CommandPermissions     `json:"permissions"`
	AllowedRoles []int                   `json:"-"`
	Arguments    []*CommandArgument      `json:"args"`
	Subcommands  []*Command              `json:"subCommands"`
	Handler      func(r *CommandContext) `json:"-"`
}

type CommandArgument struct {
	Name             string
	IncludeRemaining bool
	Optional         bool
	NoLog            bool
	Help             string
}

type CommandPermissions struct {
	RequireNoCharacter bool
	RequireCharacter   bool
	RequirePermission  string
}

type CommandContext struct {
	Command         *Command
	Player          *Player
	PlayerInitiated bool
	Character       *Character
	Args            map[string]string
	HandlerStart    time.Time
}

// CheckPermissions returns whether or not a parent can see/use the command.
func (cmd *Command) CheckPermissions(p *Player) bool {
	if cmd.Permissions == nil {
		return true
	}

	if cmd.Permissions.RequireNoCharacter {
		if p.Character() != nil {
			return false
		}
	}

	if cmd.Permissions.RequireCharacter {
		if p.Character() == nil {
			return false
		}
	}

	if len(cmd.Permissions.RequirePermission) > 0 {
		if !p.Character().HasPermission(cmd.Permissions.RequirePermission) {
			return false
		}
	}
	return true
}

// ShowSubcommandHelp returns the list of sub-commands that the parent has access to as a string.
func (cmd *Command) ShowSubcommandHelp(p *Player, commandsEntered []string) string {
	if len(cmd.Subcommands) == 0 {
		return "There are no sub-commands available."
	}

	output := []string{
		cmd.Help,
		fmt.Sprintf("%s /%s &lt;sub-command&gt;\n",
			TextStyle("Syntax:", WithBold()),
			strings.Join(commandsEntered, " "),
		),
		TextStyle("Sub-commands:", WithBold()),
	}

	var rows []string
	for _, scmd := range cmd.Subcommands {
		if cmd.CheckPermissions(p) {
			rows = append(rows, TableRow(
				TableCell{content: TextStyle(scmd.Name, WithBold())},
				TableCell{content: scmd.Help},
			))
		}
	}

	output = append(output, TextTable(rows...))

	return strings.Join(output, "\n")
}

// ShowArgumentHelp returns help for command arguments.
func (cmd *Command) ShowArgumentHelp(commandsEntered []string) string {
	if len(cmd.Arguments) == 0 {
		return "There are no command arguments."
	}

	var argumentStrings []string
	var argumentRows []string
	for _, arg := range cmd.Arguments {
		if arg.Optional {
			argumentStrings = append(argumentStrings, fmt.Sprintf("[%s]", arg.Name))
		} else {
			argumentStrings = append(argumentStrings, fmt.Sprintf("&lt;%s&gt;", arg.Name))
		}
		argumentRows = append(argumentRows, TableRow(
			TableCell{content: TextStyle(arg.Name, WithBold())},
			TableCell{content: TextStyle(misc.BoolToWords(arg.Optional, "Optional", "Required"), WithItalics())},
			TableCell{content: arg.Help},
		))
	}

	output := []string{
		cmd.Help,
		fmt.Sprintf(
			"%s /%s %s\n\n%s\n%s",
			TextStyle("Syntax:", WithBold()),
			strings.Join(commandsEntered, " "),
			strings.Join(argumentStrings, " "),
			TextStyle("Arguments:", WithBold()),
			TextTable(argumentRows...),
		),
	}

	return strings.Join(output, "\n")
}

// ArgumentByName returns a CommandArgument that matches the argument's name.
func (cmd *Command) ArgumentByName(name string) *CommandArgument {
	for _, a := range cmd.Arguments {
		if strings.ToLower(a.Name) == strings.ToLower(name) {
			return a
		}
	}

	return nil
}

// LogCtx logs a parent using a command.
func (cmd *Command) LogCtx(ctx *CommandContext) {
	handlerDuration := time.Since(ctx.HandlerStart)

	var args []string
	for k, v := range ctx.Args {
		a := cmd.ArgumentByName(k)
		if a == nil || !a.NoLog {
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}
	}

	c := "Anonymous"
	if ctx.Character != nil {
		c = ctx.Character.Name()
	}

	if ctx.Command.Parent != nil {
		Armeria.log.Info("character executed command",
			zap.String("character", c),
			zap.String("command", ctx.Command.Parent.Name),
			zap.String("sub-command", ctx.Command.Name),
			zap.Strings("arguments", args),
			zap.Duration("duration", handlerDuration),
		)
	} else {
		Armeria.log.Info("character executed command",
			zap.String("character", c),
			zap.String("command", ctx.Command.Name),
			zap.Strings("arguments", args),
			zap.Duration("duration", handlerDuration),
		)
	}
}
