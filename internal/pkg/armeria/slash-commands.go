package armeria

import (
	"fmt"
	"log"
	"strings"
)

func RegisterGameCommands(state *GameState) {
	commands := []*Command{
		{
			Name:       "login",
			SyntaxHelp: "/login [character] [password]",
			Permissions: &CommandPermissions{
				RequireNoCharacter: true,
			},
			Handler: handleLoginCommand,
		},
		{
			Name:       "look",
			SyntaxHelp: "/look",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleLookCommand,
		},
		{
			Name:       "say",
			SyntaxHelp: "/say [text]",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleSayCommand,
		},
		{
			Name:       "move",
			SyntaxHelp: "/move [dir]",
			Permissions: &CommandPermissions{
				RequireCharacter: true,
			},
			Handler: handleMoveCommand,
		},
	}

	for _, cmd := range commands {
		state.commandManager.RegisterCommand(cmd)
	}
}

func handleLoginCommand(r *CommandContext) {
	if len(r.Args) != 2 {
		r.Player.clientActions.ShowText(fmt.Sprintf("[b]Syntax:[/b] %s", r.Command.SyntaxHelp))
		return
	}

	character := r.Args[0]
	password := r.Args[1]

	c, err := r.GameState.characterManager.GetCharacterByName(character)
	if err != nil {
		r.Player.clientActions.ShowText("Character not found.")
		return
	}

	if c.GetPassword() != password {
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

	c.LoggedIn(r.GameState)

	log.Printf("[characters] character logged in: %s", c.GetName())
}

func handleLookCommand(r *CommandContext) {
	room := r.GameState.worldManager.GetRoomFromLocation(r.Player.GetCharacter().GetLocation())

	r.Player.clientActions.ShowText(
		"\nYou take a look around..\n" +
			r.Player.GetCharacter().Colorize(room.GetTitle(), COLOR_ROOM_TITLE) + "\n" +
			room.GetDescription(),
	)
}

func handleSayCommand(r *CommandContext) {
	if len(r.Args) == 0 {
		r.Player.clientActions.ShowText("Say what?")
		return
	}

	sayText := strings.Join(r.Args, " ")

	r.Player.clientActions.ShowText(
		r.Player.GetCharacter().Colorize(fmt.Sprintf("You say, \"%s\".", sayText), COLOR_SAY),
	)

	room := r.GameState.worldManager.GetRoomFromLocation(r.Player.GetCharacter().GetLocation())
	otherChars := room.GetCharacters(r.Player.GetCharacter())
	for _, c := range otherChars {
		c.GetPlayer().clientActions.ShowText(
			c.GetPlayer().GetCharacter().Colorize(
				fmt.Sprintf("%s says, \"%s\".", r.Player.GetCharacter().GetName(), sayText),
				COLOR_SAY,
			),
		)
	}
}

func handleMoveCommand(r *CommandContext) {
	if len(r.Args) != 1 {
		r.Player.clientActions.ShowText(fmt.Sprintf("[b]Syntax:[/b] %s", r.Command.SyntaxHelp))
		return
	}

	loc := r.Player.GetCharacter().GetLocation()
	area := r.GameState.worldManager.GetAreaFromLocation(loc)

	x := loc.Coords.X
	y := loc.Coords.Y
	z := loc.Coords.Z

	moveDir := r.Args[0]
	switch strings.ToLower(moveDir) {
	case "north", "n":
		y = y + 1
	case "south", "s":
		y = y - 1
	case "east", "e":
		x = x + 1
	case "west", "w":
		x = x - 1
	case "up", "u":
		z = z + 1
	case "down", "d":
		z = z - 1
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

	moveAllowed, moveError := r.Player.GetCharacter().MoveAllowed(r.GameState, newLocation)
	if !moveAllowed {
		r.Player.clientActions.ShowText(moveError)
		return
	}

}
