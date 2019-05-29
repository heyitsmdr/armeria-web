package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"strings"
)

func handleLoginCommand(r *CommandContext) {
	if len(r.Args) != 2 {
		return
	}

	c := Armeria.characterManager.GetCharacterByName(r.Args["character"])
	if c == nil {
		r.Player.clientActions.ShowText("Character not found.")
		return
	}

	if c.GetPassword() != r.Args["password"] {
		r.Player.clientActions.ShowColorizedText("Password incorrect for that character.", ColorError)
		return
	}

	if c.GetPlayer() != nil {
		r.Player.clientActions.ShowColorizedText("This character is already logged in.", ColorError)
		return
	}

	if c.GetRoom() == nil {
		r.Player.clientActions.ShowColorizedText("This character logged out of a room which no longer exists.", ColorError)
		return
	}

	r.Player.AttachCharacter(c)
	c.SetPlayer(r.Player)

	r.Player.clientActions.ShowColorizedText(fmt.Sprintf("You've entered Armeria as %s!", c.GetFName()), ColorSuccess)

	c.LoggedIn()
}

func handleLookCommand(r *CommandContext) {
	rm := Armeria.worldManager.GetRoomFromLocation(r.Character.GetLocation())

	var objNames []string
	for _, o := range rm.GetObjects() {
		if o.GetType() != ObjectTypeCharacter || o.GetName() != r.Character.GetName() {
			objNames = append(objNames, o.GetFName())
		}
	}

	var withYou string
	if len(objNames) > 0 {
		withYou = fmt.Sprintf("\nHere with you: %s.", strings.Join(objNames, ", "))
	}

	ar := r.Character.GetArea().GetAdjacentRooms(rm)
	var validDirs []string
	if ar.North != nil {
		validDirs = append(validDirs, "[b]north[/b]")
	}
	if ar.South != nil {
		validDirs = append(validDirs, "[b]south[/b]")
	}
	if ar.East != nil {
		validDirs = append(validDirs, "[b]east[/b]")
	}
	if ar.West != nil {
		validDirs = append(validDirs, "[b]west[/b]")
	}
	if ar.Up != nil {
		validDirs = append(validDirs, "[b]up[/b]")
	}
	if ar.Down != nil {
		validDirs = append(validDirs, "[b]down[/b]")
	}
	var validDirString string
	for i, d := range validDirs {
		if i == 0 {
			validDirString = fmt.Sprintf("\nYou can walk %s", d)
			if i == len(validDirs)-1 {
				validDirString = validDirString + "."
			}
		} else if i == len(validDirs)-1 {
			validDirString = fmt.Sprintf("%s and %s.", validDirString, d)
		} else {
			validDirString = fmt.Sprintf("%s, %s", validDirString, d)
		}
	}

	r.Player.clientActions.ShowText(
		r.Character.Colorize(rm.GetAttribute("title"), ColorRoomTitle) + "\n" +
			rm.GetAttribute("description") +
			r.Character.Colorize(validDirString, ColorRoomDirs) +
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
			Armeria.commandManager.ProcessCommand(r.Player, "move "+mo)
			return
		}
	}

	verbs := []string{"say", "says"}
	lastChar := r.Args["text"][len(r.Args["text"])-1:]
	if lastChar == "?" {
		verbs = []string{"ask", "asks"}
	} else if lastChar == "!" {
		verbs = []string{"exclaim", "exclaims"}
	}

	r.Player.clientActions.ShowText(
		r.Player.GetCharacter().Colorize(fmt.Sprintf("You %s, \"%s\".", verbs[0], r.Args["text"]), ColorSay),
	)

	room := Armeria.worldManager.GetRoomFromLocation(r.Player.GetCharacter().GetLocation())
	otherChars := room.GetCharacters(r.Player.GetCharacter())
	for _, c := range otherChars {
		c.GetPlayer().clientActions.ShowText(
			c.GetPlayer().GetCharacter().Colorize(
				fmt.Sprintf("%s %s, \"%s\".", r.Player.GetCharacter().GetFName(), verbs[1], r.Args["text"]),
				ColorSay,
			),
		)
	}
}

func handleMoveCommand(r *CommandContext) {
	loc := r.Character.GetLocation()
	d := r.Args["direction"]

	walkDir := ""
	arriveDir := ""
	switch strings.ToLower(d) {
	case "north", "n":
		d = "north"
		walkDir = "the north"
		arriveDir = "the south"
	case "south", "s":
		d = "south"
		walkDir = "the south"
		arriveDir = "the north"
	case "east", "e":
		d = "east"
		walkDir = "the east"
		arriveDir = "the west"
	case "west", "w":
		d = "west"
		walkDir = "west"
		arriveDir = "the east"
	case "up", "u":
		d = "up"
		walkDir = "up"
		arriveDir = "below"
	case "down", "d":
		d = "down"
		walkDir = "down"
		arriveDir = "above"
	default:
		r.Player.clientActions.ShowText("That's not a valid direction to move in.")
		return
	}

	o := misc.DirectionOffsets(d)
	x := loc.Coords.X + o["x"]
	y := loc.Coords.Y + o["y"]
	z := loc.Coords.Z + o["z"]

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
		r.Character.Colorize(fmt.Sprintf("You walk to %s.", walkDir), ColorMovement),
		r.Character.Colorize(fmt.Sprintf("%s walks to %s.", r.Character.GetFName(), walkDir), ColorMovement),
		r.Character.Colorize(fmt.Sprintf("%s walked in from %s.", r.Character.GetFName(), arriveDir), ColorMovement),
	)

	Armeria.commandManager.ProcessCommand(r.Player, "look")
}

func handleRoomEditCommand(r *CommandContext) {
	r.Player.clientActions.ShowObjectEditor(r.Character.GetRoom().GetEditorData())
}

func handleRoomSetCommand(r *CommandContext) {
	attr := strings.ToLower(r.Args["property"])

	if !misc.Contains(GetValidRoomAttributes(), attr) {
		r.Player.clientActions.ShowColorizedText("That's not a valid room attribute.", ColorError)
		return
	}

	r.Character.GetRoom().SetAttribute(attr, r.Args["value"])

	for _, c := range r.Character.GetRoom().GetCharacters(r.Character) {
		c.GetPlayer().clientActions.ShowText(
			fmt.Sprintf("%s modified the room.", r.Character.GetFName()),
		)
	}

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("You modified the [b]%s[/b] property of the room.", attr),
		ColorSuccess,
	)

	editorOpen := r.Character.GetTempAttribute("editorOpen")
	if editorOpen == "true" {
		r.Player.clientActions.ShowObjectEditor(r.Character.GetRoom().GetEditorData())
	}
}

func handleRoomCreateCommand(r *CommandContext) {
	d := r.Args["direction"]

	o := misc.DirectionOffsets(d)
	if o == nil {
		r.Player.clientActions.ShowColorizedText("That's not a valid direction to create a room in.", ColorError)
		return
	}

	loc := r.Character.GetLocation()
	x := loc.Coords.X + o["x"]
	y := loc.Coords.Y + o["y"]
	z := loc.Coords.Z + o["z"]

	newLoc := &Location{
		AreaName: loc.AreaName,
		Coords: &Coords{
			X: x,
			Y: y,
			Z: z,
		},
	}

	if Armeria.worldManager.GetRoomFromLocation(newLoc) != nil {
		r.Player.clientActions.ShowColorizedText("There's already a room in that direction.", ColorError)
		return
	}

	r.Character.GetArea().AddRoom(&Room{
		Coords: newLoc.Coords,
	})

	for _, c := range r.Character.GetArea().GetCharacters(nil) {
		c.GetPlayer().clientActions.RenderMap()
	}

	r.Player.clientActions.ShowText("A new room has been created.")
}

func handleRoomDestroyCommand(r *CommandContext) {
	d := r.Args["direction"]

	o := misc.DirectionOffsets(d)
	if o == nil {
		r.Player.clientActions.ShowColorizedText("That's not a valid direction to destroy a room in.", ColorError)
		return
	}

	loc := r.Character.GetLocation()
	x := loc.Coords.X + o["x"]
	y := loc.Coords.Y + o["y"]
	z := loc.Coords.Z + o["z"]

	l := &Location{
		AreaName: loc.AreaName,
		Coords: &Coords{
			X: x,
			Y: y,
			Z: z,
		},
	}

	rm := Armeria.worldManager.GetRoomFromLocation(l)
	if rm == nil {
		r.Player.clientActions.ShowColorizedText("There's no room in that direction.", ColorError)
		return
	}

	if len(rm.GetCharacters(nil)) > 0 {
		r.Player.clientActions.ShowColorizedText("There are characters in the room you're attempting to destroy.", ColorError)
		return
	}

	r.Player.clientActions.ShowText("Success.")
}

func handleSaveCommand(r *CommandContext) {
	Armeria.Save()
	r.Player.clientActions.ShowText("The game data has been saved to disk.")
}

func handleReloadCommand(r *CommandContext) {
	if r.Args["component"] != "server" && r.Args["component"] != "client" && r.Args["component"] != "both" {
		r.Player.clientActions.ShowText("You can reload the following components: server, client, or both.")
		return
	}

	if !Armeria.production {
		r.Player.clientActions.ShowColorizedText("You can only reload in production!", ColorError)
		return
	}

	Armeria.Reload(r.Player, r.Args["component"])
}

func handleRefreshCommand(r *CommandContext) {
	r.Player.clientActions.RenderMap()
	r.Player.clientActions.SyncRoomObjects()
	r.Player.clientActions.ShowText("Client data has been refreshed.")
}

func handleWhisperCommand(r *CommandContext) {
	t := r.Args["target"]
	m := r.Args["message"]

	c := Armeria.characterManager.GetCharacterByName(t)
	if c == nil {
		r.Player.clientActions.ShowColorizedText("That's not a valid character name.", ColorError)
		return
	} else if c.GetPlayer() == nil {
		r.Player.clientActions.ShowColorizedText("That character is not online.", ColorError)
		return
	}

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("You whisper to %s, \"%s\".", c.GetFName(), m),
		ColorWhisper,
	)

	c.GetPlayer().clientActions.ShowColorizedText(
		fmt.Sprintf("%s whispers to you, \"%s\".", r.Character.GetFName(), m),
		ColorWhisper,
	)
}

func handleWhoCommand(r *CommandContext) {
	chars := Armeria.characterManager.GetCharacters()

	noun := "characters"
	verb := "are"
	if len(chars) < 2 {
		noun = "character"
		verb = "is"
	}

	var fn []string
	for _, c := range chars {
		fn = append(fn, c.GetFNameWithTitle())
	}

	r.Player.clientActions.ShowText(
		fmt.Sprintf(
			"There %s %d %s playing right now:\n%s",
			verb,
			len(chars),
			noun,
			strings.Join(fn, ", ")+".",
		),
	)
}

func handleCharacterEditCommand(r *CommandContext) {
	char := r.Args["character"]
	var c *Character
	if len(char) == 0 {
		c = r.Character
	} else {
		c = Armeria.characterManager.GetCharacterByName(strings.ToLower(char))
		if c == nil {
			r.Player.clientActions.ShowColorizedText("That character doesn't exist.", ColorError)
			return
		}
	}

	r.Player.clientActions.ShowObjectEditor(c.GetEditorData())
}

func handleCharacterSetCommand(r *CommandContext) {
	char := strings.ToLower(r.Args["character"])
	attr := strings.ToLower(r.Args["property"])
	val := r.Args["value"]

	c := Armeria.characterManager.GetCharacterByName(char)
	if c == nil {
		r.Player.clientActions.ShowColorizedText("That character doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(GetValidCharacterAttributes(), attr) {
		r.Player.clientActions.ShowColorizedText("That's not a valid character attribute.", ColorError)
		return
	}

	c.SetAttribute(attr, val)

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("You modified the [b]%s[/b] property of the character %s.", attr, c.GetFName()),
		ColorSuccess,
	)

	if c.GetName() != r.Character.GetName() && c.GetPlayer() != nil {
		c.GetPlayer().clientActions.ShowText(
			fmt.Sprintf("Your character was modified by %s.", r.Character.GetFName()),
		)

	}

	editorOpen := r.Character.GetTempAttribute("editorOpen")
	if editorOpen == "true" {
		r.Player.clientActions.ShowObjectEditor(c.GetEditorData())
	}
}

func handleMobCreateCommand(r *CommandContext) {
	n := r.Args["name"]

	if Armeria.mobManager.GetMobByName(n) != nil {
		r.Player.clientActions.ShowColorizedText("A mob already exists with that name.", ColorError)
		return
	}

	m := &Mob{
		Name: n,
	}

	Armeria.mobManager.CreateMob(m)

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("A mob named [b]%s[/b] has been created.", n),
		ColorSuccess,
	)
}

func handleMobEditCommand(r *CommandContext) {
	mname := r.Args["mob"]

	m := Armeria.mobManager.GetMobByName(mname)
	if m == nil {
		r.Player.clientActions.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	r.Player.clientActions.ShowObjectEditor(m.GetEditorData())
}

func handleMobSetCommand(r *CommandContext) {
	mob := strings.ToLower(r.Args["mob"])
	attr := strings.ToLower(r.Args["property"])
	val := strings.ToLower(r.Args["value"])

	m := Armeria.mobManager.GetMobByName(mob)
	if m == nil {
		r.Player.clientActions.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(GetValidMobAttributes(), attr) {
		r.Player.clientActions.ShowColorizedText("That's not a valid mob attribute.", ColorError)
		return
	}

	m.SetAttribute(attr, val)

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("You modified the [b]%s[/b] property of the mob [b]%s[/b].", attr, m.Name),
		ColorSuccess,
	)

	editorOpen := r.Character.GetTempAttribute("editorOpen")
	if editorOpen == "true" {
		r.Player.clientActions.ShowObjectEditor(m.GetEditorData())
	}
}

func handleMobSpawnCommand(r *CommandContext) {
	mob := strings.ToLower(r.Args["mob"])

	m := Armeria.mobManager.GetMobByName(mob)
	if m == nil {
		r.Player.clientActions.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	l := r.Character.GetLocation()
	loc := &Location{
		AreaName: l.AreaName,
		Coords: &Coords{
			X: l.Coords.X,
			Y: l.Coords.Y,
			Z: l.Coords.Z,
			I: l.Coords.I,
		},
	}

	mi := m.CreateInstance(loc)
	r.Character.GetRoom().AddObjectToRoom(mi)

	for _, c := range r.Character.GetRoom().GetCharacters(nil) {
		c.GetPlayer().clientActions.ShowText(
			fmt.Sprintf("With a flash of light, a %s appeared out of nowhere!", mi.GetFName()),
		)
		c.GetPlayer().clientActions.SyncRoomObjects()
	}
}
