package armeria

import (
	"fmt"
	"log"
	"strings"
)

func RegisterGameCommands(state *GameState) {
	commands := []*Command{
		{
			Name: "login",
			Permissions: &CommandPermissions{
				RequireNoCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Position: 0,
					Name:     "character",
				},
				{
					Position: 1,
					Name:     "password",
				},
			},
			Handler: handleLoginCommand,
		},
		{
			Name: "look",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleLookCommand,
		},
		{
			Name: "say",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Position:         0,
					Name:             "text",
					IncludeRemaining: true,
				},
			},
			Handler: handleSayCommand,
		},
		{
			Name: "move",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Arguments: []*CommandArgument{
				{
					Position: 0,
					Name:     "direction",
				},
			},
			Handler: handleMoveCommand,
		},
		{Name: "north", Alias: "move north"},
		{Name: "south", Alias: "move south"},
		{Name: "east", Alias: "move east"},
		{Name: "west", Alias: "move west"},
		{Name: "up", Alias: "move up"},
		{Name: "down", Alias: "move down"},
		{
			Name: "room",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Subcommands: []*Command{
				{
					Name: "set",
					Help: "allows you to set a room attribute",
					Arguments: []*CommandArgument{
						{
							Position: 0,
							Name:     "property",
						},
						{
							Position:         1,
							Name:             "value",
							IncludeRemaining: true,
						},
					},
					Handler: handleRoomSetCommand,
				},
			},
		},
		{
			Name: "save",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleSaveCommand,
		},
	}

	for _, cmd := range commands {
		state.commandManager.RegisterCommand(cmd)
	}
}

func handleLoginCommand(r *CommandContext) {
	if len(r.Args) != 2 {
		return
	}

	c, err := r.GameState.characterManager.GetCharacterByName(r.Args["character"])
	if err != nil {
		r.Player.clientActions.ShowText("Character not found.")
		return
	}

	if c.GetPassword() != r.Args["password"] {
		r.Player.clientActions.ShowText("Password incorrect for that character.")
		return
	}

	if c.GetPlayer() != nil {
		r.Player.clientActions.ShowText("That character is already logged in.")
		return
	}

	r.Player.AttachCharacter(c)
	c.SetPlayer(r.Player)

	r.Player.clientActions.ShowText(fmt.Sprintf("You've successfully logged in to %s!", c.GetName()))

	c.LoggedIn()

	log.Printf("[characters] character logged in: %s", c.GetName())
}

func handleLookCommand(r *CommandContext) {
	room := r.GameState.worldManager.GetRoomFromLocation(r.Player.GetCharacter().GetLocation())

	var objNames []string
	for _, o := range room.GetObjects() {
		if o.GetType() != ObjectTypeCharacter || o.GetName() != r.Character.GetName() {
			objNames = append(objNames, o.GetFName())
		}
	}

	withYou := ""
	if len(objNames) > 0 {
		withYou = fmt.Sprintf("\nHere with you: %s.", strings.Join(objNames, ", "))
	}

	r.Player.clientActions.ShowText(
		"\n" +
			r.Player.GetCharacter().Colorize(room.GetTitle(), ColorRoomTitle) + "\n" +
			room.GetDescription() +
			withYou,
	)
}

func handleSayCommand(r *CommandContext) {
	if len(r.Args) == 0 {
		r.Player.clientActions.ShowText("Say what?")
		return
	}

	var moveOverride = []string{"n", "s", "e", "w", "u", "d"}
	for _, mo := range moveOverride {
		if r.Args["text"] == mo {
			r.GameState.commandManager.ProcessCommand(r.Player, "move "+mo)
			return
		}
	}

	r.Player.clientActions.ShowText(
		r.Player.GetCharacter().Colorize(fmt.Sprintf("You say, \"%s\".", r.Args["text"]), ColorSay),
	)

	room := r.GameState.worldManager.GetRoomFromLocation(r.Player.GetCharacter().GetLocation())
	otherChars := room.GetCharacters(r.Player.GetCharacter())
	for _, c := range otherChars {
		c.GetPlayer().clientActions.ShowText(
			c.GetPlayer().GetCharacter().Colorize(
				fmt.Sprintf("%s says, \"%s\".", r.Player.GetCharacter().GetFName(), r.Args["text"]),
				ColorSay,
			),
		)
	}
}

func handleMoveCommand(r *CommandContext) {
	if len(r.Args) != 1 {
		return
	}

	loc := r.Character.GetLocation()

	x := loc.Coords.X
	y := loc.Coords.Y
	z := loc.Coords.Z

	walkDir := ""
	arriveDir := ""
	switch strings.ToLower(r.Args["direction"]) {
	case "north", "n":
		y = y + 1
		walkDir = "the north"
		arriveDir = "the south"
	case "south", "s":
		y = y - 1
		walkDir = "the south"
		arriveDir = "the north"
	case "east", "e":
		x = x + 1
		walkDir = "the east"
		arriveDir = "the west"
	case "west", "w":
		x = x - 1
		walkDir = "west"
		arriveDir = "the east"
	case "up", "u":
		z = z + 1
		walkDir = "up"
		arriveDir = "below"
	case "down", "d":
		z = z - 1
		walkDir = "down"
		arriveDir = "above"
	default:
		r.Player.clientActions.ShowText("That's not a valid direction to move in.")
		return
	}

	newLocation := &Location{
		AreaName: loc.AreaName,
		Coords: &Coords{
			X: x,
			Y: y,
			Z: z,
			I: loc.Coords.I,
		},
	}

	moveAllowed, moveError := r.Character.MoveAllowed(newLocation)
	if !moveAllowed {
		r.Player.clientActions.ShowText(moveError)
		return
	}

	r.Character.Move(
		newLocation,
		r.Character.Colorize(fmt.Sprintf("\nYou walk to %s.", walkDir), ColorMovement),
		r.Character.Colorize(fmt.Sprintf("%s walks to %s.", r.Character.GetFName(), walkDir), ColorMovement),
		r.Character.Colorize(fmt.Sprintf("%s walked in from %s.", r.Character.GetFName(), arriveDir), ColorMovement),
	)

	r.GameState.commandManager.ProcessCommand(r.Player, "look")
}

func handleRoomSetCommand(r *CommandContext) {
	switch strings.ToLower(r.Args["property"]) {
	case "title":
		r.Character.GetRoom().SetTitle(r.Args["value"])

	default:
		r.Player.clientActions.ShowText("Invalid room property.")
		return
	}

	for _, c := range r.Character.GetRoom().GetCharacters(r.Character) {
		c.GetPlayer().clientActions.ShowText(
			fmt.Sprintf("%s modified the room.", r.Character.GetFName()),
		)
	}

	r.Player.clientActions.ShowText("You modified the room.")
}

func handleSaveCommand(r *CommandContext) {
	r.GameState.characterManager.SaveCharacters()
	r.GameState.worldManager.SaveWorld()
	r.Player.clientActions.ShowText("The game data has been saved to disk.")
}
