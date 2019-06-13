package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func handleLoginCommand(r *CommandContext) {
	if len(r.Args) != 2 {
		return
	}

	c := Armeria.characterManager.CharacterByName(r.Args["character"])
	if c == nil {
		r.Player.clientActions.ShowText("Character not found.")
		return
	}

	if c.Password() != r.Args["password"] {
		r.Player.clientActions.ShowColorizedText("Password incorrect for that character.", ColorError)
		return
	}

	if c.Player() != nil {
		r.Player.clientActions.ShowColorizedText("This character is already logged in.", ColorError)
		return
	}

	if c.Room() == nil {
		r.Player.clientActions.ShowColorizedText("This character logged out of a room which no longer exists.", ColorError)
		return
	}

	r.Player.AttachCharacter(c)
	c.SetPlayer(r.Player)

	r.Player.clientActions.ShowColorizedText(fmt.Sprintf("You've entered Armeria as %s!", c.FormattedName()), ColorSuccess)

	c.LoggedIn()
}

func handleLookCommand(r *CommandContext) {
	rm := Armeria.worldManager.RoomFromLocation(r.Character.Location())

	var objNames []string
	for _, o := range rm.Objects() {
		if o.Type() != ObjectTypeCharacter || o.Name() != r.Character.Name() {
			objNames = append(objNames, o.FormattedName())
		}
	}

	var withYou string
	if len(objNames) > 0 {
		withYou = fmt.Sprintf("\nHere with you: %s.", strings.Join(objNames, ", "))
	}

	ar := r.Character.Area().AdjacentRooms(rm)
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
		r.Character.Colorize(rm.Attribute("title"), ColorRoomTitle) + "\n" +
			rm.Attribute("description") +
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

	room := Armeria.worldManager.RoomFromLocation(r.Player.GetCharacter().Location())
	otherChars := room.Characters(r.Player.GetCharacter())
	for _, c := range otherChars {
		c.Player().clientActions.ShowText(
			c.Player().GetCharacter().Colorize(
				fmt.Sprintf("%s %s, \"%s\".", r.Player.GetCharacter().FormattedName(), verbs[1], r.Args["text"]),
				ColorSay,
			),
		)
	}

	for _, o := range room.Objects() {
		if o.Type() == ObjectTypeMob {
			go CallMobFunc(
				r.Character,
				o.(*MobInstance),
				"character_said",
				lua.LString(r.Args["text"]),
			)
		}
	}
}

func handleMoveCommand(r *CommandContext) {
	loc := r.Character.Location()
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
		r.Character.Colorize(fmt.Sprintf("%s walks to %s.", r.Character.FormattedName(), walkDir), ColorMovement),
		r.Character.Colorize(fmt.Sprintf("%s walked in from %s.", r.Character.FormattedName(), arriveDir), ColorMovement),
	)

	Armeria.commandManager.ProcessCommand(r.Player, "look")
}

func handleRoomEditCommand(r *CommandContext) {
	r.Player.clientActions.ShowObjectEditor(r.Character.Room().EditorData())
}

func handleRoomSetCommand(r *CommandContext) {
	attr := strings.ToLower(r.Args["property"])

	if !misc.Contains(ValidRoomAttributes(), attr) {
		r.Player.clientActions.ShowColorizedText("That's not a valid room attribute.", ColorError)
		return
	}

	r.Character.Room().SetAttribute(attr, r.Args["value"])

	for _, c := range r.Character.Room().Characters(r.Character) {
		c.Player().clientActions.ShowText(
			fmt.Sprintf("%s modified the room.", r.Character.FormattedName()),
		)
	}

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("You modified the [b]%s[/b] property of the room.", attr),
		ColorSuccess,
	)

	editorOpen := r.Character.GetTempAttribute("editorOpen")
	if editorOpen == "true" {
		r.Player.clientActions.ShowObjectEditor(r.Character.Room().EditorData())
	}
}

func handleRoomCreateCommand(r *CommandContext) {
	d := r.Args["direction"]

	o := misc.DirectionOffsets(d)
	if o == nil {
		r.Player.clientActions.ShowColorizedText("That's not a valid direction to create a room in.", ColorError)
		return
	}

	loc := r.Character.Location()
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

	if Armeria.worldManager.RoomFromLocation(newLoc) != nil {
		r.Player.clientActions.ShowColorizedText("There's already a room in that direction.", ColorError)
		return
	}

	r.Character.Area().AddRoom(&Room{
		UnafeCoords: newLoc.Coords,
	})

	for _, c := range r.Character.Area().Characters(nil) {
		c.Player().clientActions.RenderMap()
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

	loc := r.Character.Location()
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

	rm := Armeria.worldManager.RoomFromLocation(l)
	if rm == nil {
		r.Player.clientActions.ShowColorizedText("There's no room in that direction.", ColorError)
		return
	}

	if len(rm.Characters(nil)) > 0 {
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
	r.Player.clientActions.SyncRoomTitle()
	r.Player.clientActions.ShowText("Client data has been refreshed.")
}

func handleWhisperCommand(r *CommandContext) {
	t := r.Args["target"]
	m := r.Args["message"]

	c := Armeria.characterManager.CharacterByName(t)
	if c == nil {
		r.Player.clientActions.ShowColorizedText("That's not a valid character name.", ColorError)
		return
	} else if c.Player() == nil {
		r.Player.clientActions.ShowColorizedText("That character is not online.", ColorError)
		return
	}

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("You whisper to %s, \"%s\".", c.FormattedName(), m),
		ColorWhisper,
	)

	c.Player().clientActions.ShowColorizedText(
		fmt.Sprintf("%s whispers to you, \"%s\".", r.Character.FormattedName(), m),
		ColorWhisper,
	)
}

func handleWhoCommand(r *CommandContext) {
	chars := Armeria.characterManager.OnlineCharacters()

	noun := "characters"
	verb := "are"
	if len(chars) < 2 {
		noun = "character"
		verb = "is"
	}

	var fn []string
	for _, c := range chars {
		fn = append(fn, c.FormattedNameWithTitle())
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
		c = Armeria.characterManager.CharacterByName(strings.ToLower(char))
		if c == nil {
			r.Player.clientActions.ShowColorizedText("That character doesn't exist.", ColorError)
			return
		}
	}

	r.Player.clientActions.ShowObjectEditor(c.GetEditorData())
}

func handleCharacterListCommand(r *CommandContext) {
	f := r.Args["filter"]

	var chars []string
	for _, c := range Armeria.characterManager.Characters() {
		if len(f) == 0 || strings.Contains(strings.ToLower(c.Name()), strings.ToLower(f)) {
			chars = append(chars, c.Name())
		}
	}

	var matchingText string
	if len(f) > 0 {
		matchingText = " matching \"" + f + "\""
	}

	if len(chars) == 0 {
		r.Player.clientActions.ShowColorizedText(
			fmt.Sprintf("There are no characters matching \"%s\".", f),
			ColorError,
		)
		return
	}

	r.Player.clientActions.ShowText(
		fmt.Sprintf("There are [b]%d[/b] characters%s: %s.", len(chars), matchingText, strings.Join(chars, ", ")),
	)
}

func handleCharacterSetCommand(r *CommandContext) {
	char := strings.ToLower(r.Args["character"])
	attr := strings.ToLower(r.Args["property"])
	val := r.Args["value"]

	c := Armeria.characterManager.CharacterByName(char)
	if c == nil {
		r.Player.clientActions.ShowColorizedText("That character doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(ValidCharacterAttributes(), attr) {
		r.Player.clientActions.ShowColorizedText("That's not a valid character attribute.", ColorError)
		return
	}

	c.SetAttribute(attr, val)

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("You modified the [b]%s[/b] property of the character %s.", attr, c.FormattedName()),
		ColorSuccess,
	)

	if c.Name() != r.Character.Name() && c.Player() != nil {
		c.Player().clientActions.ShowText(
			fmt.Sprintf("Your character was modified by %s.", r.Character.FormattedName()),
		)

	}

	editorOpen := r.Character.GetTempAttribute("editorOpen")
	if editorOpen == "true" {
		r.Player.clientActions.ShowObjectEditor(c.GetEditorData())
	}
}

func handleMobListCommand(r *CommandContext) {
	f := r.Args["filter"]

	var mobs []string
	for _, m := range Armeria.mobManager.Mobs() {
		if len(f) == 0 || strings.Contains(strings.ToLower(m.Name()), strings.ToLower(f)) {
			mobs = append(mobs, m.Name())
		}
	}

	var matchingText string
	if len(f) > 0 {
		matchingText = " matching \"" + f + "\""
	}

	if len(mobs) == 0 {
		r.Player.clientActions.ShowColorizedText(
			fmt.Sprintf("There are no mobs matching \"%s\".", f),
			ColorError,
		)
		return
	}

	r.Player.clientActions.ShowText(
		fmt.Sprintf("There are [b]%d[/b] mobs%s: %s.", len(mobs), matchingText, strings.Join(mobs, ", ")),
	)
}

func handleMobCreateCommand(r *CommandContext) {
	n := r.Args["name"]

	if Armeria.mobManager.MobByName(n) != nil {
		r.Player.clientActions.ShowColorizedText("A mob already exists with that name.", ColorError)
		return
	}

	m := Armeria.mobManager.CreateMob(n)
	Armeria.mobManager.AddMob(m)

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("A mob named [b]%s[/b] has been created.", n),
		ColorSuccess,
	)
}

func handleMobEditCommand(r *CommandContext) {
	mname := r.Args["mob"]

	m := Armeria.mobManager.MobByName(mname)
	if m == nil {
		r.Player.clientActions.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	r.Player.clientActions.ShowObjectEditor(m.EditorData())
}

func handleMobSetCommand(r *CommandContext) {
	mob := strings.ToLower(r.Args["mob"])
	attr := strings.ToLower(r.Args["property"])
	val := strings.ToLower(r.Args["value"])

	m := Armeria.mobManager.MobByName(mob)
	if m == nil {
		r.Player.clientActions.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(ValidMobAttributes(), attr) {
		r.Player.clientActions.ShowColorizedText("That's not a valid mob attribute.", ColorError)
		return
	}

	valid, why := ValidateMobAttribute(attr, val)
	if !valid {
		r.Player.clientActions.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", why), ColorError)
		return
	}

	m.SetAttribute(attr, val)

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("You modified the [b]%s[/b] property of the mob [b]%s[/b].", attr, m.UnsafeName),
		ColorSuccess,
	)

	editorOpen := r.Character.GetTempAttribute("editorOpen")
	if editorOpen == "true" {
		r.Player.clientActions.ShowObjectEditor(m.EditorData())
	}
}

func handleMobSpawnCommand(r *CommandContext) {
	mob := strings.ToLower(r.Args["mob"])

	m := Armeria.mobManager.MobByName(mob)
	if m == nil {
		r.Player.clientActions.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	l := r.Character.Location()
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
	r.Character.Room().AddObjectToRoom(mi)

	for _, c := range r.Character.Room().Characters(nil) {
		c.Player().clientActions.ShowText(
			fmt.Sprintf("With a flash of light, a %s appeared out of nowhere!", mi.FormattedName()),
		)
		c.Player().clientActions.SyncRoomObjects()
	}
}

func handleMobInstancesCommand(r *CommandContext) {
	mob := strings.ToLower(r.Args["mob"])

	m := Armeria.mobManager.MobByName(mob)
	if m == nil {
		r.Player.clientActions.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	var mobLocations []string
	for i, mi := range m.Instances() {
		mobLocations = append(
			mobLocations,
			fmt.Sprintf(
				"  %d) %s (%s) is currently at %s,%d,%d,%d (%s).",
				i+1,
				mi.FormattedName(),
				mi.Id(),
				mi.Location().AreaName,
				mi.Location().Coords.X,
				mi.Location().Coords.Y,
				mi.Location().Coords.Z,
				mi.Room().Attribute("title"),
			),
		)
	}

	r.Player.clientActions.ShowText(
		fmt.Sprintf(
			"Instances of %s:\n%s",
			m.Name(),
			strings.Join(mobLocations, "\n"),
		),
	)
}

func handleWipeCommand(r *CommandContext) {
	for _, o := range r.Character.Room().Objects() {
		switch o.Type() {
		case ObjectTypeMob:
			m := Armeria.mobManager.MobByName(o.Name())
			s := r.Character.Room().RemoveObjectFromRoom(o)
			if m != nil && s {
				m.DeleteInstance(o.(*MobInstance))
			}
		}
	}

	for _, c := range r.Character.Room().Characters(r.Character) {
		c.Player().clientActions.ShowText(
			fmt.Sprintf("%s wiped the room.", r.Character.FormattedName()),
		)
		c.Player().clientActions.SyncRoomObjects()
	}

	r.Player.clientActions.ShowColorizedText("You wiped the room.", ColorSuccess)
	r.Player.clientActions.SyncRoomObjects()
}

func handleItemCreateCommand(r *CommandContext) {
	n := r.Args["name"]

	if Armeria.itemManager.ItemByName(n) != nil {
		r.Player.clientActions.ShowColorizedText("An item already exists with that name.", ColorError)
		return
	}

	i := &Item{
		UnsafeName: n,
	}

	Armeria.itemManager.AddItem(i)

	r.Player.clientActions.ShowColorizedText(
		fmt.Sprintf("An item named [b]%s[/b] has been created.", n),
		ColorSuccess,
	)
}

func handleItemListCommand(r *CommandContext) {
	f := r.Args["filter"]

	var items []string
	for _, i := range Armeria.itemManager.Items() {
		if len(f) == 0 || strings.Contains(strings.ToLower(i.Name()), strings.ToLower(f)) {
			items = append(items, i.Name())
		}
	}

	var matchingText string
	if len(f) > 0 {
		matchingText = " matching \"" + f + "\""
	}

	if len(items) == 0 {
		r.Player.clientActions.ShowColorizedText(
			fmt.Sprintf("There are no items matching \"%s\".", f),
			ColorError,
		)
		return
	}

	r.Player.clientActions.ShowText(
		fmt.Sprintf("There are [b]%d[/b] items%s: %s.", len(items), matchingText, strings.Join(items, ", ")),
	)
}
