package armeria

import (
	"fmt"
	"strings"
)

// Channel describes a particular talking channel.
type Channel struct {
	Name              string
	Description       string
	SlashCommand      string
	Color             int
	RequirePermission string
}

// NewChannels returns a map containing an instance of each talking channel.
func NewChannels() map[string]*Channel {
	return map[string]*Channel{
		"general": {
			Name:         "General",
			Description:  "Chat about anything game-related.",
			Color:        ColorChannelGeneral,
			SlashCommand: "/general",
		},
		"core": {
			Name:              "Core",
			Description:       "Chat with other Armeria Core developers.",
			SlashCommand:      "/core",
			Color:             ColorChannelCore,
			RequirePermission: "CAN_SYSOP",
		},
	}
}

// ChannelByName returns the matching Channel.
func ChannelByName(name string) *Channel {
	for _, c := range Armeria.channels {
		if strings.ToLower(c.Name) == strings.ToLower(name) {
			return c
		}
	}

	return nil
}

// HasPermission returns a bool indicating whether the character can participate in the channel.
func (c *Channel) HasPermission(char *Character) bool {
	if len(c.RequirePermission) > 0 {
		return char.HasPermission(c.RequirePermission)
	}
	return true
}

// Broadcast sends a message to all logged-in players that have joined the channel. You can pass
// nil as the Character if this is coming from a system rather than a particular character.
func (c *Channel) Broadcast(from *Character, text string) {
	var msgToOthers string
	var msgToFrom string
	var verbs []string

	if from != nil {
		normalizedText, textType := TextPunctuation(text)
		switch textType {
		case TextQuestion:
			verbs = []string{"ask", "asks"}
		case TextExclaim:
			verbs = []string{"exclaim", "exclaims"}
		default:
			verbs = []string{"say", "says"}
		}
		msgToOthers = fmt.Sprintf("[[b]%s[/b]] %s %s, \"%s\"", c.Name, from.FormattedNameWithTitle(), verbs[1], normalizedText)
		msgToFrom = fmt.Sprintf("[[b]%s[/b]] You %s, \"%s\"", c.Name, verbs[0], normalizedText)
	} else {
		msgToOthers = fmt.Sprintf("[[b]%s[/b]] %s", c.Name, text)
	}

	for _, char := range Armeria.characterManager.OnlineCharacters() {
		if char.InChannel(c) {
			if from == nil || from.ID() != char.ID() {
				char.Player().client.ShowColorizedText(
					msgToOthers,
					c.Color,
				)
			}
		}
	}

	if from != nil {
		from.Player().client.ShowColorizedText(
			msgToFrom,
			c.Color,
		)
	}
}
