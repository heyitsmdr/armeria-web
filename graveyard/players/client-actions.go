package players

import (
	schemaCommands "armeria/internal/pkg/commands/schema"
	"fmt"
	"strings"
)

type clientAction struct {
	player *player
}

func newClientAction(p *player) clientAction {
	return clientAction{
		player: p,
	}
}

// ShowText displays raw text on the player's main text window
func (ca clientAction) ShowText(text string) {
	ca.player.CallClientAction("showText", text)
}

// ShowIntroText displays the text the player will see when first connecting
func (ca clientAction) ShowIntroText() {
	lines := []string{
		"Welcome to the world of Armeria!\n",
		"[ARMERIA ASCII ART HERE]\n",
		"If you have a character, you can <b>/login</b>. Otherwise, use <b>/create</b>.",
	}

	ca.player.CallClientAction("showText", strings.Join(lines, "<br>"))
}

func (ca clientAction) ShowCommandSyntax(command *schemaCommands.Command) {
	syntax := fmt.Sprintf("Syntax: %s", command.SyntaxHelp)
	ca.player.CallClientAction("showText", syntax)
}
