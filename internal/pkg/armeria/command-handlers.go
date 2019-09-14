package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"

	lua "github.com/yuin/gopher-lua"
)

func handleLoginCommand(ctx *CommandContext) {
	var c *Character

	if len(ctx.Args) == 1 {
		// token auth
		sections := strings.Split(ctx.Args["token"], ":")

		if len(sections) != 2 {
			ctx.Player.client.ShowText("Token formatted incorrectly.")
			return
		}

		c = Armeria.characterManager.CharacterByName(sections[0])
		if c == nil {
			ctx.Player.client.ShowText("Character not found.")
			return
		}

		if c.PasswordHash() != sections[1] {
			ctx.Player.client.ShowColorizedText("Invalid token for that character.", ColorError)
			return
		}
	} else {
		// basic auth
		c = Armeria.characterManager.CharacterByName(ctx.Args["character"])
		if c == nil {
			ctx.Player.client.ShowText("Character not found.")
			return
		}

		if !c.CheckPassword(ctx.Args["password"]) {
			ctx.Player.client.ShowColorizedText("Password incorrect for that character.", ColorError)
			return
		}
	}

	if c.Player() != nil {
		ctx.Player.client.ShowColorizedText("This character is already logged in.", ColorError)
		return
	}

	if c.Room() == nil {
		ctx.Player.client.ShowColorizedText("This character logged out of a room which no longer exists.", ColorError)
		return
	}

	ctx.Player.AttachCharacter(c)
	c.SetPlayer(ctx.Player)

	ctx.Player.client.ShowColorizedText(fmt.Sprintf("You've entered Armeria as %s!", c.FormattedName()), ColorSuccess)

	c.LoggedIn()
}

func handleCreateCommand(ctx *CommandContext) {
	ctx.Player.client.ShowText("You cannot manually create new characters yet. Stay tuned.")
}

func handleLookCommand(ctx *CommandContext) {
	r := ctx.Character.Room()
	at := ctx.Args["at"]

	// Look _at_ something?
	if len(at) > 0 {
		var o interface{}
		var rt RegistryType
		var oc *ObjectContainer
		var searchInv bool
		if len(at) > 4 && at[0:4] == "inv:" {
			at = at[4:]
			oc = ctx.Character.Inventory()
			searchInv = true
		} else {
			oc = ctx.Character.Room().Here()
			searchInv = false
		}

		o, _, rt = oc.Get(at)
		if rt == RegistryTypeUnknown {
			o, _, rt = oc.GetByName(at)
			if rt == RegistryTypeUnknown {
				ctx.Player.client.ShowColorizedText("You don't see anything by that name.", ColorError)
				return
			}
		}

		co := o.(ContainerObject)
		var lookResult string
		if rt == RegistryTypeItemInstance {
			lookResult = co.Attribute(AttributeDescription)
		} else if rt == RegistryTypeCharacter {
			lookResult = fmt.Sprintf("There is nothing special about %s.", co.(*Character).Pronoun(PronounObjective))
		}

		if len(lookResult) == 0 {
			lookResult = "There is nothing special about it."
		}

		ctx.Player.client.ShowText(
			fmt.Sprintf("You take a look at %s.\n%s", TextStyle(co.FormattedName(), TextStyleBold), lookResult),
		)

		if ctx.PlayerInitiated {
			for _, c := range r.Here().Characters(true, ctx.Character) {
				if searchInv {
					c.Player().client.ShowText(
						fmt.Sprintf("%s is taking a look at something within %s inventory.",
							ctx.Character.FormattedName(),
							ctx.Character.Pronoun(PronounPossessiveAdjective),
						),
					)
				} else {
					c.Player().client.ShowText(
						fmt.Sprintf("%s is taking a look at %s.",
							ctx.Character.FormattedName(),
							TextStyle(co.FormattedName(), TextStyleBold),
						),
					)
				}
			}
		}
		return
	}

	var objNames []string
	for _, o := range r.Here().All() {
		obj := o.(ContainerObject)

		if obj.Type() == ContainerObjectTypeCharacter && obj.(*Character).Player() == nil {
			continue
		}

		if obj.ID() != ctx.Character.ID() {
			objNames = append(objNames, obj.FormattedName())
		}
	}

	var withYou string
	if len(objNames) > 0 {
		withYou = fmt.Sprintf("\nHere with you: %s.", strings.Join(objNames, ", "))
	}

	ar := r.AdjacentRooms()
	var validDirs []string
	if ar.North != nil {
		validDirs = append(validDirs, TextStyle("north", TextStyleBold))
	}
	if ar.South != nil {
		validDirs = append(validDirs, TextStyle("south", TextStyleBold))
	}
	if ar.East != nil {
		validDirs = append(validDirs, TextStyle("east", TextStyleBold))
	}
	if ar.West != nil {
		validDirs = append(validDirs, TextStyle("west", TextStyleBold))
	}
	if ar.Up != nil {
		validDirs = append(validDirs, TextStyle("up", TextStyleBold))
	}
	if ar.Down != nil {
		validDirs = append(validDirs, TextStyle("down", TextStyleBold))
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

	ctx.Player.client.ShowText(
		ctx.Character.Colorize(r.Attribute("title"), ColorRoomTitle) + "\n" +
			r.Attribute("description") +
			ctx.Character.Colorize(validDirString, ColorRoomDirs) +
			withYou,
	)

	if ctx.PlayerInitiated {
		for _, c := range r.Here().Characters(true, ctx.Character) {
			c.Player().client.ShowText(
				fmt.Sprintf("%s takes a look around.", ctx.Character.FormattedName()),
			)
		}
	}
}

func handleGlanceCommand(ctx *CommandContext) {
	r := ctx.Character.Room()

	var objNames []string
	for _, o := range r.Here().All() {
		obj := o.(ContainerObject)

		if obj.Type() == ContainerObjectTypeCharacter && obj.(*Character).Player() == nil {
			continue
		}

		if obj.ID() != ctx.Character.ID() {
			objNames = append(objNames, obj.FormattedName())
		}
	}

	var withYou string
	if len(objNames) > 0 {
		withYou = fmt.Sprintf("\nHere with you: %s.", strings.Join(objNames, ", "))
	}

	ar := r.AdjacentRooms()
	var validDirs []string
	if ar.North != nil {
		validDirs = append(validDirs, TextStyle("north", TextStyleBold))
	}
	if ar.South != nil {
		validDirs = append(validDirs, TextStyle("south", TextStyleBold))
	}
	if ar.East != nil {
		validDirs = append(validDirs, TextStyle("east", TextStyleBold))
	}
	if ar.West != nil {
		validDirs = append(validDirs, TextStyle("west", TextStyleBold))
	}
	if ar.Up != nil {
		validDirs = append(validDirs, TextStyle("up", TextStyleBold))
	}
	if ar.Down != nil {
		validDirs = append(validDirs, TextStyle("down", TextStyleBold))
	}
	var validDirString string
	for i, d := range validDirs {
		if i == 0 {
			validDirString = fmt.Sprintf("\n[ Exits: %s", d)
			if i == len(validDirs)-1 {
				validDirString = validDirString + " ]"
			}
		} else if i == len(validDirs)-1 {
			validDirString = fmt.Sprintf("%s and %s ]", validDirString, d)
		} else {
			validDirString = fmt.Sprintf("%s, %s", validDirString, d)
		}
	}

	ctx.Player.client.ShowText(
		ctx.Character.Colorize(r.Attribute("title"), ColorRoomTitle) +
			ctx.Character.Colorize(validDirString, ColorRoomDirs) +
			withYou,
	)

	if ctx.PlayerInitiated {
		for _, c := range r.Here().Characters(true, ctx.Character) {
			c.Player().client.ShowText(
				fmt.Sprintf("%s glances around.", ctx.Character.FormattedName()),
			)
		}
	}
}

func handleSayCommand(ctx *CommandContext) {
	if len(ctx.Args) == 0 {
		ctx.Player.client.ShowText("Say what?")
		return
	}

	var moveOverride = []string{"n", "s", "e", "w", "u", "d"}
	for _, mo := range moveOverride {
		if ctx.Args["text"] == mo {
			Armeria.commandManager.ProcessCommand(ctx.Player, "move "+mo, true)
			return
		}
	}

	normalizedText, textType := TextPunctuation(ctx.Args["text"])

	var verbs []string
	switch textType {
	case TextQuestion:
		verbs = []string{"ask", "asks"}
	case TextExclaim:
		verbs = []string{"exclaim", "exclaims"}
	default:
		verbs = []string{"say", "says"}
	}

	ctx.Player.client.ShowText(
		ctx.Player.Character().Colorize(fmt.Sprintf("You %s, \"%s\"", verbs[0], normalizedText), ColorSay),
	)

	room := ctx.Character.Room()
	for _, c := range room.Here().Characters(true, ctx.Character) {
		c.Player().client.ShowText(
			c.Player().Character().Colorize(
				fmt.Sprintf("%s %s, \"%s\"", ctx.Character.FormattedName(), verbs[1], normalizedText),
				ColorSay,
			),
		)
	}

	for _, mi := range room.Here().Mobs() {
		go CallMobFunc(
			ctx.Character,
			mi,
			"character_said",
			lua.LString(ctx.Character.Name()),
			lua.LString(ctx.Args["text"]),
		)
	}
}

func handleMoveCommand(ctx *CommandContext) {
	d := ctx.Args["direction"]
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
		ctx.Player.client.ShowColorizedText("That's not a valid direction to move in.", ColorError)
		return
	}

	newRoom := Armeria.worldManager.RoomInDirection(ctx.Character.Room(), d)

	areaLink := ctx.Character.Room().Attribute(d)
	if areaLink != "" {
		na := strings.Split(areaLink, ",")
		if len(na) != 4 {
			ctx.Player.client.ShowColorizedText("You cannot go that way.", ColorError)
			Armeria.channels[ChannelBuilders].Broadcast(
				nil,
				fmt.Sprintf("ERROR: Invalid area link at %s (%s) due to using an invalid format.", ctx.Character.Room().LocationString(), d),
			)
			return
		}
		ta := Armeria.worldManager.AreaByName(na[0])
		if ta == nil {
			ctx.Player.client.ShowColorizedText("You cannot go that way.", ColorError)
			Armeria.channels[ChannelBuilders].Broadcast(
				nil,
				fmt.Sprintf("ERROR: Invalid area link at %s (%s) due to the area not existing.", ctx.Character.Room().LocationString(), d),
			)
			return
		}
		x, xerr := strconv.Atoi(na[1])
		y, yerr := strconv.Atoi(na[2])
		z, zerr := strconv.Atoi(na[3])
		if xerr != nil || yerr != nil || zerr != nil {
			ctx.Player.client.ShowColorizedText("You cannot go that way.", ColorError)
			Armeria.channels[ChannelBuilders].Broadcast(
				nil,
				fmt.Sprintf("ERROR: Invalid area link at %s (%s) due to the coords not existing within the target area.", ctx.Character.Room().LocationString(), d),
			)
			return
		}
		newRoom = ta.RoomAt(NewCoords(x, y, z, 0))
	}

	moveAllowed, moveError := ctx.Character.MoveAllowed(newRoom)
	if !moveAllowed {
		ctx.Player.client.ShowColorizedText(moveError, ColorError)
		return
	}

	oldRoomUUID := ctx.Character.Room().ParentArea.Id()
	ctx.Character.Move(
		newRoom,
		ctx.Character.Colorize(fmt.Sprintf("You walk to %s.", walkDir), ColorMovement),
		ctx.Character.Colorize(fmt.Sprintf("%s walks to %s.", ctx.Character.FormattedName(), walkDir), ColorMovement),
		ctx.Character.Colorize(fmt.Sprintf("%s walked in from %s.", ctx.Character.FormattedName(), arriveDir), ColorMovement),
	)

	if newRoom.ParentArea.Id() != oldRoomUUID {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("You've just entered %s.", TextStyle(newRoom.ParentArea.Name(), TextStyleBold)),
			ColorMovementAlt,
		)
	}

	if ctx.Character.Setting(SettingBrief) == "true" {
		Armeria.commandManager.ProcessCommand(ctx.Player, "glance", false)
	} else {
		Armeria.commandManager.ProcessCommand(ctx.Player, "look", false)
	}
}

func handleRoomEditCommand(ctx *CommandContext) {
	t := ctx.Args["target"]
	a := ctx.Character.Room().ParentArea
	tr := ctx.Character.Room()

	args := strings.Split(t, ",")

	if args[0] == "" {
		ctx.Player.client.ShowObjectEditor(tr.EditorData())
		return
	}

	if len(args) != 3 {
		ctx.Player.client.ShowColorizedText("Incorrect format for the target room. Use [x],[y],[z].", ColorError)
		return
	}

	x, xerr := strconv.Atoi(args[0])
	y, yerr := strconv.Atoi(args[1])
	z, zerr := strconv.Atoi(args[2])
	if xerr != nil || yerr != nil || zerr != nil {
		ctx.Player.client.ShowColorizedText("The x, y, and z coordinates must be valid numbers.", ColorError)
		return
	}

	tr = a.RoomAt(NewCoords(x, y, z, 0))
	if tr == nil {
		ctx.Player.client.ShowColorizedText("The specified room does not exist.", ColorError)
		return
	}

	ctx.Player.client.ShowObjectEditor(tr.EditorData())
}

func handleRoomSetCommand(ctx *CommandContext) {
	attr := strings.ToLower(ctx.Args["property"])
	if !misc.Contains(ValidRoomAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid room attribute.", ColorError)
		return
	}
	ta := ctx.Args["target"]
	tr := ctx.Character.Room()

	if ta != "." {
		ts := strings.Split(ta, ",")

		if len(ts) != 3 {
			ctx.Player.client.ShowColorizedText(
				"Incorrect format for the target room. Use [x],[y],[z] (or \".\" to reference the current room).",
				ColorError,
			)
			return
		}

		x, xerr := strconv.Atoi(ts[0])
		y, yerr := strconv.Atoi(ts[1])
		z, zerr := strconv.Atoi(ts[2])
		if xerr != nil || yerr != nil || zerr != nil {
			ctx.Player.client.ShowColorizedText("The x, y, and z coordinates must be valid numbers.", ColorError)
			return
		}
		tr = ctx.Character.Room().ParentArea.RoomAt(NewCoords(x, y, z, 0))
	}

	if tr != nil {
		tr.SetAttribute(attr, ctx.Args["value"])
	} else {
		ctx.Player.client.ShowColorizedText("The specified room does not exist.", ColorError)
		return
	}

	ctx.Player.client.SyncMap()

	for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
		c.Player().client.ShowText(
			fmt.Sprintf("%s modified the room.", ctx.Character.FormattedName()),
		)
	}
	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the room (%s).", TextStyle(attr, TextStyleBold), ta),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(tr.EditorData())
	}
}

func handleRoomCreateCommand(ctx *CommandContext) {
	d := ctx.Args["direction"]

	o := misc.DirectionOffsets(d)
	if o == nil {
		ctx.Player.client.ShowColorizedText("That's not a valid direction to create a room in.", ColorError)
		return
	}

	co := ctx.Character.Room().Coords
	x := co.X() + o["x"]
	y := co.Y() + o["y"]
	z := co.Z() + o["z"]

	c := NewCoords(x, y, z, 0)
	if ctx.Character.Room().ParentArea.RoomAt(c) != nil {
		ctx.Player.client.ShowColorizedText("There's already a room in that direction.", ColorError)
		return
	}

	_ = Armeria.worldManager.CreateRoom(ctx.Character.Room().ParentArea, c)

	for _, c := range ctx.Character.Room().ParentArea.Characters() {
		c.Player().client.SyncMap()
	}

	ctx.Player.client.ShowText("A new room has been created.")
}

func handleRoomDestroyCommand(ctx *CommandContext) {
	d := ctx.Args["direction"]

	o := misc.DirectionOffsets(d)
	if o == nil {
		ctx.Player.client.ShowColorizedText("That's not a valid direction to destroy a room in.", ColorError)
		return
	}

	co := ctx.Character.Room().Coords
	x := co.X() + o["x"]
	y := co.Y() + o["y"]
	z := co.Z() + o["z"]

	c := NewCoords(x, y, z, 0)
	r := ctx.Character.Room().ParentArea.RoomAt(c)
	if r == nil {
		ctx.Player.client.ShowColorizedText("There's no room in that direction.", ColorError)
		return
	}

	if r.Here().Count() > 0 {
		ctx.Player.client.ShowColorizedText("There is something in the room you're attempting to destroy.", ColorError)
		return
	}

	r.ParentArea.RemoveRoom(r)

	for _, c := range ctx.Character.Room().ParentArea.Characters() {
		c.Player().client.SyncMap()
	}

	ctx.Player.client.ShowColorizedText("The room has been destroyed.", ColorSuccess)
}

func handleSaveCommand(ctx *CommandContext) {
	Armeria.Save()
	ctx.Player.client.ShowText("The game data has been saved to disk.")
}

func handleRefreshCommand(ctx *CommandContext) {
	ctx.Player.client.SyncMap()
	ctx.Player.client.SyncRoomObjects()
	ctx.Player.client.SyncRoomTitle()
	ctx.Player.client.SyncInventory()
	ctx.Player.client.SyncPermissions()
	ctx.Player.client.SyncPlayerInfo()
	ctx.Player.client.ShowText("Client data has been refreshed.")
}

func handleWhisperCommand(ctx *CommandContext) {
	t := ctx.Args["target"]
	m := ctx.Args["message"]

	c := Armeria.characterManager.CharacterByName(t)
	if c == nil {
		ctx.Player.client.ShowColorizedText("That's not a valid character name.", ColorError)
		return
	} else if c.Player() == nil {
		ctx.Player.client.ShowColorizedText("That character is not online.", ColorError)
		return
	}

	c.SetTempAttribute(TempAttributeReplyTo, ctx.Character.Name())

	normalizedText, _ := TextPunctuation(m)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You whisper to %s, \"%s\"", c.FormattedNameWithTitle(), normalizedText),
		ColorWhisper,
	)

	c.Player().client.ShowColorizedText(
		fmt.Sprintf("%s whispers to you from %s, \"%s\"",
			ctx.Character.FormattedNameWithTitle(),
			c.Room().ParentArea.Name(),
			normalizedText,
		),
		ColorWhisper,
	)
}

func handleReplyCommand(ctx *CommandContext) {
	m := ctx.Args["message"]

	rt := ctx.Character.TempAttribute(TempAttributeReplyTo)
	if len(rt) == 0 {
		ctx.Player.client.ShowColorizedText("No one has sent you a whisper yet.", ColorError)
		return
	}

	Armeria.commandManager.ProcessCommand(ctx.Player, fmt.Sprintf("whisper %s %s", rt, m), false)
}

func handleWhoCommand(ctx *CommandContext) {
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

	ctx.Player.client.ShowText(
		fmt.Sprintf(
			"There %s %d %s playing right now:\n%s",
			verb,
			len(chars),
			noun,
			strings.Join(fn, ", ")+".",
		),
	)
}

func handleCharacterEditCommand(ctx *CommandContext) {
	char := ctx.Args["character"]
	var c *Character
	if len(char) == 0 {
		c = ctx.Character
	} else {
		c = Armeria.characterManager.CharacterByName(strings.ToLower(char))
		if c == nil {
			ctx.Player.client.ShowColorizedText("That character doesn't exist.", ColorError)
			return
		}
	}

	ctx.Player.client.ShowObjectEditor(c.EditorData())
}

func handleCharacterListCommand(ctx *CommandContext) {
	f := ctx.Args["filter"]

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
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("There are no characters matching \"%s\".", f),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(
		fmt.Sprintf("There are %s characters%s: %s.", TextStyle(len(chars), TextStyleBold), matchingText, strings.Join(chars, ", ")),
	)
}

func handleCharacterSetCommand(ctx *CommandContext) {
	char := strings.ToLower(ctx.Args["character"])
	attr := strings.ToLower(ctx.Args["property"])
	val := ctx.Args["value"]

	c := Armeria.characterManager.CharacterByName(char)
	if c == nil {
		ctx.Player.client.ShowColorizedText("That character doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(ValidCharacterAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid character attribute.", ColorError)
		return
	}

	if len(val) > 0 {
		valid, why := ValidateCharacterAttribute(attr, val)
		if !valid {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", why), ColorError)
			return
		}
	}

	_ = c.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the character %s.", TextStyle(attr, TextStyleBold), c.FormattedName()),
		ColorSuccess,
	)

	if c.Name() != ctx.Character.Name() && c.Player() != nil {
		c.Player().client.ShowText(
			fmt.Sprintf("Your character was modified by %s.", ctx.Character.FormattedName()),
		)

	}

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(c.EditorData())
	}
}

func handleMobListCommand(ctx *CommandContext) {
	f := ctx.Args["filter"]

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
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("There are no mobs matching \"%s\".", f),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(
		fmt.Sprintf("There are %s mobs%s: %s.", TextStyle(len(mobs), TextStyleBold), matchingText, strings.Join(mobs, ", ")),
	)
}

func handleMobCreateCommand(ctx *CommandContext) {
	n := ctx.Args["name"]

	if Armeria.mobManager.MobByName(n) != nil {
		ctx.Player.client.ShowColorizedText("A mob already exists with that name.", ColorError)
		return
	}

	m := Armeria.mobManager.CreateMob(n)
	Armeria.mobManager.AddMob(m)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("A mob named %s has been created.", TextStyle(n, TextStyleBold)),
		ColorSuccess,
	)
}

func handleMobEditCommand(ctx *CommandContext) {
	mname := ctx.Args["mob"]

	m := Armeria.mobManager.MobByName(mname)
	if m == nil {
		ctx.Player.client.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	ctx.Player.client.ShowObjectEditor(m.EditorData())
}

func handleMobInstanceEditCommand(ctx *CommandContext) {
	o, rt := Armeria.registry.Get(ctx.Args["uuid"])
	if rt == RegistryTypeUnknown {
		ctx.Player.client.ShowColorizedText("That uuid doesn't exist.", ColorError)
		return
	} else if rt != RegistryTypeMobInstance {
		ctx.Player.client.ShowColorizedText("That uuid is not a mob.", ColorError)
		return
	}

	mi := o.(*MobInstance)

	ctx.Player.client.ShowObjectEditor(mi.EditorData())
}

func handleMobSetCommand(ctx *CommandContext) {
	mob := strings.ToLower(ctx.Args["mob"])
	attr := strings.ToLower(ctx.Args["property"])
	val := ctx.Args["value"]

	m := Armeria.mobManager.MobByName(mob)
	if m == nil {
		ctx.Player.client.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(ValidMobAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid mob attribute.", ColorError)
		return
	}

	if len(val) > 0 {
		valid, why := ValidateMobAttribute(attr, val)
		if !valid {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", why), ColorError)
			return
		}
	}

	m.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the mob %s.",
			TextStyle(attr, TextStyleBold),
			TextStyle(m.UnsafeName, TextStyleBold),
		),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(m.EditorData())
	}
}

func handleMobInstanceSetCommand(ctx *CommandContext) {
	o, rt := Armeria.registry.Get(ctx.Args["uuid"])
	if rt == RegistryTypeUnknown {
		ctx.Player.client.ShowColorizedText("That uuid doesn't exist.", ColorError)
		return
	} else if rt != RegistryTypeMobInstance {
		ctx.Player.client.ShowColorizedText("That uuid is not a mob.", ColorError)
		return
	}

	mi := o.(*MobInstance)
	attr := strings.ToLower(ctx.Args["property"])
	val := strings.ToLower(ctx.Args["value"])

	if !misc.Contains(ValidMobAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid mob attribute.", ColorError)
		return
	} else if !misc.Contains(ValidMobInstanceAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not an attribute you can set on a mob instance.", ColorError)
		return
	}

	if len(val) > 0 {
		valid, why := ValidateMobAttribute(attr, val)
		if !valid {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", why), ColorError)
			return
		}
	}

	_ = mi.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the mob instace %s (%s).",
			TextStyle(attr, TextStyleBold),
			TextStyle(mi.Name(), TextStyleBold),
			mi.ID(),
		),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(mi.EditorData())
	}
}

func handleMobSpawnCommand(ctx *CommandContext) {
	m := Armeria.mobManager.MobByName(ctx.Args["mob"])
	if m == nil {
		ctx.Player.client.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	mi := m.CreateInstance()
	_ = ctx.Character.Room().Here().Add(mi.ID())

	for _, c := range ctx.Character.Room().Here().Characters(true) {
		c.Player().client.ShowText(
			fmt.Sprintf("With a flash of light, a %s appeared out of nowhere!", mi.FormattedName()),
		)
		c.Player().client.SyncRoomObjects()
	}
}

func handleMobInstancesCommand(ctx *CommandContext) {
	m := Armeria.mobManager.MobByName(ctx.Args["mob"])
	if m == nil {
		ctx.Player.client.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	rows := []string{TableRow(
		TableCell{content: "Mob", header: true},
		TableCell{content: "UUID", header: true},
		TableCell{content: "Location", header: true},
	)}

	for _, mi := range m.Instances() {
		rows = append(rows, TableRow(
			TableCell{content: mi.FormattedName()},
			TableCell{content: mi.ID()},
			TableCell{content: fmt.Sprintf("%s (%s)", mi.Room().LocationString(), mi.Room().Attribute("title"))},
		))
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleWipeCommand(ctx *CommandContext) {
	for _, o := range ctx.Character.Room().Here().All() {
		obj := o.(ContainerObject)
		switch obj.Type() {
		case ContainerObjectTypeMob:
			m := Armeria.mobManager.MobByName(obj.Name())
			ctx.Character.Room().Here().Remove(obj.ID())
			if m != nil {
				m.DeleteInstance(obj.(*MobInstance))
			}
		case ContainerObjectTypeItem:
			i := Armeria.itemManager.ItemByName(obj.Name())
			ctx.Character.Room().Here().Remove(obj.ID())
			if i != nil {
				i.DeleteInstance(obj.(*ItemInstance))
			}
		}
	}

	for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
		c.Player().client.ShowText(
			fmt.Sprintf("%s wiped the room.", ctx.Character.FormattedName()),
		)
		c.Player().client.SyncRoomObjects()
	}

	ctx.Player.client.ShowColorizedText("You wiped the room.", ColorSuccess)
	ctx.Player.client.SyncRoomObjects()
}

func handleItemCreateCommand(ctx *CommandContext) {
	n := ctx.Args["name"]

	if Armeria.itemManager.ItemByName(n) != nil {
		ctx.Player.client.ShowColorizedText("An item already exists with that name.", ColorError)
		return
	}

	i := Armeria.itemManager.CreateItem(n)
	Armeria.itemManager.AddItem(i)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("An item named %s has been created.", TextStyle(n, TextStyleBold)),
		ColorSuccess,
	)
}

func handleItemListCommand(ctx *CommandContext) {
	f := ctx.Args["filter"]

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
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("There are no items matching \"%s\".", f),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(
		fmt.Sprintf("There are %s items%s: %s.", TextStyle(len(items), TextStyleBold), matchingText, strings.Join(items, ", ")),
	)
}

func handleItemSpawnCommand(ctx *CommandContext) {
	i := Armeria.itemManager.ItemByName(ctx.Args["item"])
	if i == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist.", ColorError)
		return
	}

	ii := i.CreateInstance()
	_ = ctx.Character.Room().Here().Add(ii.ID())

	for _, c := range ctx.Character.Room().Here().Characters(true) {
		c.Player().client.ShowText(
			fmt.Sprintf("With a flash of light, a %s appeared out of nowhere!", ii.FormattedName()),
		)
		c.Player().client.SyncRoomObjects()
	}
}

func handleItemEditCommand(ctx *CommandContext) {
	i := Armeria.itemManager.ItemByName(ctx.Args["item"])
	if i == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist.", ColorError)
		return
	}

	ctx.Player.client.ShowObjectEditor(i.EditorData())
}

func handleItemInstanceEditCommand(ctx *CommandContext) {
	o, rt := Armeria.registry.Get(ctx.Args["uuid"])
	if rt == RegistryTypeUnknown {
		ctx.Player.client.ShowColorizedText("That uuid doesn't exist.", ColorError)
		return
	} else if rt != RegistryTypeItemInstance {
		ctx.Player.client.ShowColorizedText("That uuid is not an item.", ColorError)
		return
	}

	ii := o.(*ItemInstance)

	ctx.Player.client.ShowObjectEditor(ii.EditorData())
}

func handleItemSetCommand(ctx *CommandContext) {
	item := strings.ToLower(ctx.Args["item"])
	attr := strings.ToLower(ctx.Args["property"])
	val := ctx.Args["value"]

	i := Armeria.itemManager.ItemByName(item)
	if i == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(ValidItemAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid item attribute.", ColorError)
		return
	}

	if len(val) > 0 {
		valid, why := ValidateItemAttribute(attr, val)
		if !valid {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", why), ColorError)
			return
		}
	}

	i.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the item %s.",
			TextStyle(attr, TextStyleBold),
			TextStyle(i.Name(), TextStyleBold),
		),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(i.EditorData())
	}
}

func handleItemInstanceSetCommand(ctx *CommandContext) {
	o, rt := Armeria.registry.Get(ctx.Args["uuid"])
	if rt == RegistryTypeUnknown {
		ctx.Player.client.ShowColorizedText("That uuid doesn't exist.", ColorError)
		return
	} else if rt != RegistryTypeItemInstance {
		ctx.Player.client.ShowColorizedText("That uuid is not an item.", ColorError)
		return
	}

	ii := o.(*ItemInstance)
	attr := strings.ToLower(ctx.Args["property"])
	val := strings.ToLower(ctx.Args["value"])

	if !misc.Contains(ValidItemAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid item attribute.", ColorError)
		return
	} else if !misc.Contains(ValidItemInstanceAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not an attribute you can set on an item instance.", ColorError)
		return
	}

	if len(val) > 0 {
		valid, why := ValidateItemAttribute(attr, val)
		if !valid {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", why), ColorError)
			return
		}
	}

	_ = ii.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the item instace %s (%s).",
			TextStyle(attr, TextStyleBold),
			TextStyle(ii.Name(), TextStyleBold),
			ii.ID(),
		),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(ii.EditorData())
	}
}

func handleItemInstancesCommand(ctx *CommandContext) {
	i := Armeria.itemManager.ItemByName(ctx.Args["item"])
	if i == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist.", ColorError)
		return
	}

	rows := []string{TableRow(
		TableCell{content: "Item", header: true},
		TableCell{content: "UUID", header: true},
		TableCell{content: "Location", header: true},
	)}

	for _, ii := range i.Instances() {
		ctr := Armeria.registry.GetObjectContainer(ii.ID())
		if ctr.ParentType() == ContainerParentTypeRoom {
			rows = append(rows, TableRow(
				TableCell{content: ii.FormattedName()},
				TableCell{content: ii.ID()},
				TableCell{content: fmt.Sprintf("%s (%s)", ii.Room().LocationString(), ii.Room().Attribute("title"))},
			))
		} else if ctr.ParentType() == ContainerParentTypeCharacter {
			rows = append(rows, TableRow(
				TableCell{content: ii.FormattedName()},
				TableCell{content: ii.ID()},
				TableCell{
					content: fmt.Sprintf(
						"%s (slot %d)",
						ii.Character().FormattedName(),
						ii.Character().Inventory().Slot(ii.ID()),
					),
					styling: "",
				},
			))
		}
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleGhostCommand(ctx *CommandContext) {
	if len(ctx.Character.TempAttribute(TempAttributeGhost)) > 0 {
		ctx.Character.SetTempAttribute(TempAttributeGhost, "")
		ctx.Player.client.ShowColorizedText("You are no longer ghostly.", ColorSuccess)
	} else {
		ctx.Character.SetTempAttribute(TempAttributeGhost, "1")
		ctx.Player.client.ShowColorizedText("You are now ghostly.", ColorSuccess)
	}
}

func handleAreaCreateCommand(ctx *CommandContext) {
	n := ctx.Args["name"]

	if Armeria.worldManager.AreaByName(n) != nil {
		ctx.Player.client.ShowColorizedText("An area by that name already exists.", ColorError)
		return
	}

	a := Armeria.worldManager.CreateArea(n)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("An area named %s has been created!", TextStyle(a.Name(), TextStyleBold)),
		ColorSuccess,
	)
}

func handleAreaListCommand(ctx *CommandContext) {
	f := ctx.Args["filter"]

	var areas []string
	for _, a := range Armeria.worldManager.Areas() {
		if len(f) == 0 || strings.Contains(strings.ToLower(a.Name()), strings.ToLower(f)) {
			areas = append(areas, a.Name())
		}
	}

	var matchingText string
	if len(f) > 0 {
		matchingText = " matching \"" + f + "\""
	}

	if len(areas) == 0 {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("There are no areas matching \"%s\".", f),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(
		fmt.Sprintf("There are %s areas%s: %s.", TextStyle(len(areas), TextStyleBold), matchingText, strings.Join(areas, ", ")),
	)
}

func handleAreaEditCommand(ctx *CommandContext) {
	area := ctx.Args["area"]
	var a *Area
	if len(area) == 0 {
		a = ctx.Character.Room().ParentArea
	} else {
		a = Armeria.worldManager.AreaByName(area)
		if a == nil {
			ctx.Player.client.ShowColorizedText("That area doesn't exist.", ColorError)
			return
		}
	}

	ctx.Player.client.ShowObjectEditor(a.EditorData())
}

func handlePasswordCommand(ctx *CommandContext) {
	pw := ctx.Args["password"]
	ctx.Character.SetPassword(pw)
	ctx.Player.client.ShowColorizedText("Your character password has been set.", ColorSuccess)
}

func handleTeleportCommand(ctx *CommandContext) {
	t := ctx.Args["destination"]

	var destination *Room
	var moveMsg string
	if t[0:1] == "@" {
		cn := t[1:]
		c := Armeria.characterManager.CharacterByName(cn)
		if c == nil {
			ctx.Player.client.ShowColorizedText("There is no character by that name.", ColorError)
			return
		} else if c.Player() == nil {
			ctx.Player.client.ShowColorizedText("That character is not online.", ColorError)
			return
		}

		destination = c.Room()
		moveMsg = fmt.Sprintf("You teleported to %s.", c.FormattedName())
	} else {
		loc := strings.Split(t, ",")
		if len(loc) != 4 {
			ctx.Player.client.ShowColorizedText("Incorrect format for teleport. Use [area],[x],[y],[z].", ColorError)
			return
		}

		a := Armeria.worldManager.AreaByName(loc[0])
		if a == nil {
			ctx.Player.client.ShowColorizedText("That is not a valid area.", ColorError)
			return
		}

		var x, y, z int
		x, xerr := strconv.Atoi(loc[1])
		y, yerr := strconv.Atoi(loc[2])
		z, zerr := strconv.Atoi(loc[3])
		if xerr != nil || yerr != nil || zerr != nil {
			ctx.Player.client.ShowColorizedText("The x, y, and z coordinate must be a valid number.", ColorError)
			return
		}

		destination = a.RoomAt(NewCoords(x, y, z, 0))
		moveMsg = fmt.Sprintf("You teleported to %s at %d, %d, %d.", TextStyle(a.Name(), TextStyleBold), x, y, z)
	}

	if destination == nil {
		ctx.Player.client.ShowColorizedText("You cannot teleport there!", ColorError)
		return
	}

	ctx.Character.Move(
		destination,
		ctx.Character.Colorize(moveMsg, ColorMovement),
		ctx.Character.Colorize(fmt.Sprintf("%s teleported away!", ctx.Character.FormattedName()), ColorMovement),
		ctx.Character.Colorize(fmt.Sprintf("%s teleported here!", ctx.Character.FormattedName()), ColorMovement),
	)

	Armeria.commandManager.ProcessCommand(ctx.Player, "look", false)
}

func handleCommandsCommand(ctx *CommandContext) {
	var valid []*Command
	var largest int
	for _, cmd := range Armeria.commandManager.Commands() {
		if cmd.CheckPermissions(ctx.Player) && len(cmd.Alias) == 0 && !cmd.Hidden {
			valid = append(valid, cmd)

			if len(cmd.Name) > largest {
				largest = len(cmd.Name)
			}
		}
	}

	rows := []string{TableRow(
		TableCell{content: "Command", header: true},
		TableCell{content: "Description", header: true},
	)}

	for _, cmd := range valid {
		rows = append(rows, TableRow(
			TableCell{content: TextStyle("/"+cmd.Name, TextStyleBold), styling: "padding:0px 2px"},
			TableCell{content: cmd.Help},
		))
	}

	ctx.Player.client.ShowColorizedText(TextTable(rows...), ColorCmdHelp)
}

func handleClipboardCopyCommand(ctx *CommandContext) {
	t := strings.ToLower(ctx.Args["type"])
	n := ctx.Args["name"]
	a := ctx.Args["attributes"]

	attrs := strings.Split(a, " ")

	switch t {
	case "room":
		// validate room (based on "name")
		var r *Room
		if n == "." || n == "here" {
			r = ctx.Character.Room()
		} else {
			ctx.Player.client.ShowColorizedText("That room is not valid.", ColorError)
			return
		}
		// validate room attributes; allow * for everything
		if a == "*" {
			attrs = ValidRoomAttributes()
		}
		for _, attr := range attrs {
			if !misc.Contains(ValidRoomAttributes(), attr) {
				ctx.Player.client.ShowColorizedText(
					fmt.Sprintf("Invalid room attribute: %s.", attr),
					ColorError,
				)
				return
			}
		}
		ctx.Character.SetTempAttribute("clipboard_type", t)
		ctx.Character.SetTempAttribute("clipboard_attributes", a)
		for _, attr := range attrs {
			ctx.Character.SetTempAttribute("clipboard_attribute_"+attr, r.Attribute(attr))
		}
		ctx.Player.client.ShowColorizedText("Room attributes copied to clipboard.", ColorSuccess)
	default:
		ctx.Player.client.ShowColorizedText("That object type is not supported by the clipboard.", ColorError)
		return
	}
}

func handleClipboardPasteCommand(ctx *CommandContext) {
	n := ctx.Args["name"]

	cbt := ctx.Character.TempAttribute("clipboard_type")
	if len(cbt) == 0 {
		ctx.Player.client.ShowColorizedText("You don't have anything on your clipboard.", ColorError)
		return
	}

	cba := strings.Split(ctx.Character.TempAttribute("clipboard_attributes"), " ")

	switch cbt {
	case "room":
		// validate room (based on "name")
		var r *Room
		if n == "." || n == "here" {
			r = ctx.Character.Room()
		} else {
			ctx.Player.client.ShowColorizedText("That room is not valid.", ColorError)
			return
		}
		// paste room attributes
		for _, attr := range cba {
			attrValue := ctx.Character.TempAttribute("clipboard_attribute_" + attr)
			r.SetAttribute(attr, attrValue)
		}
		ctx.Player.client.ShowColorizedText("Room attributes on the clipboard have been applied.", ColorSuccess)
	default:
		ctx.Player.client.ShowColorizedText("That object type cannot be pasted anywhere.", ColorError)
		return
	}
}

func handleClipboardClearCommand(ctx *CommandContext) {
	cba := strings.Split(ctx.Character.TempAttribute("clipboard_attributes"), " ")
	for _, attr := range cba {
		ctx.Character.SetTempAttribute("clipboard_attribute_"+attr, "")
	}

	ctx.Character.SetTempAttribute("clipboard_type", "")
	ctx.Character.SetTempAttribute("clipboard_attributes", "")

	ctx.Player.client.ShowColorizedText("Your clipboard has been cleared.", ColorSuccess)
}

func handleGetCommand(ctx *CommandContext) {
	istring := ctx.Args["item"]

	roomObjects := ctx.Character.Room().Here()
	o, _, rt := roomObjects.GetByName(istring)
	if o == nil {
		ctx.Player.client.ShowColorizedText("There is nothing here by that name.", ColorError)
		return
	} else if rt != RegistryTypeItemInstance {
		ctx.Player.client.ShowColorizedText("You cannot pick that up.", ColorError)
		return
	}

	item := o.(*ItemInstance)

	roomObjects.Remove(item.ID())
	err := ctx.Character.Inventory().Add(item.ID())
	if err == ErrContainerNoRoom {
		_ = roomObjects.Add(item.ID())
		ctx.Player.client.ShowColorizedText("You have no room in your inventory.", ColorError)
		return
	} else if err == ErrContainerDuplicate {
		_ = roomObjects.Add(item.ID())
		ctx.Player.client.ShowColorizedText("You already have that item instance in your inventory.", ColorError)
		return
	}

	ctx.Player.client.SyncRoomObjects()
	ctx.Player.client.SyncInventory()
	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You picked up a %s.", item.FormattedName()),
		ColorSuccess,
	)

	for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
		c.Player().client.SyncRoomObjects()
		c.Player().client.ShowColorizedText(
			fmt.Sprintf("%s picked up a %s.", ctx.Character.FormattedName(), item.FormattedName()),
			ColorSuccess,
		)
	}
}

func handleDropCommand(ctx *CommandContext) {
	istring := ctx.Args["item"]

	var item *ItemInstance

	// check using uuid first, followed by item name
	i, _, _ := ctx.Character.Inventory().Get(istring)
	if i != nil {
		item = i.(*ItemInstance)
	} else {
		i, _, _ = ctx.Character.Inventory().GetByName(istring)
		if i != nil {
			item = i.(*ItemInstance)
		} else {
			ctx.Player.client.ShowColorizedText("You don't have that item in your inventory.", ColorError)
			return
		}
	}

	ctx.Character.Inventory().Remove(item.ID())
	_ = ctx.Character.Room().Here().Add(item.ID())

	ctx.Player.client.SyncRoomObjects()
	ctx.Player.client.SyncInventory()
	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You dropped a %s.", item.FormattedName()),
		ColorSuccess,
	)

	for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
		c.Player().client.SyncRoomObjects()
		c.Player().client.ShowColorizedText(
			fmt.Sprintf("%s dropped a %s.", ctx.Character.FormattedName(), item.FormattedName()),
			ColorSuccess,
		)
	}
}

func handleSwapCommand(ctx *CommandContext) {
	source := ctx.Args["source"]
	destination := ctx.Args["destination"]

	snum, err := strconv.Atoi(source)
	if err != nil {
		ctx.Player.client.ShowColorizedText("You must enter a valid slot number as the source slot.", ColorError)
		return
	}

	dnum, err := strconv.Atoi(destination)
	if err != nil {
		ctx.Player.client.ShowColorizedText("You must enter a valid slot number as the destination slot.", ColorError)
		return
	}

	if snum >= ctx.Character.Inventory().MaxSize() {
		ctx.Player.client.ShowColorizedText("Source slot does not exist in your inventory.", ColorError)
		return
	}

	if dnum >= ctx.Character.Inventory().MaxSize() {
		ctx.Player.client.ShowColorizedText("Destination slot does not exist in your inventory.", ColorError)
		return
	}

	sitem, _, _ := ctx.Character.Inventory().AtSlot(snum)
	ditem, _, _ := ctx.Character.Inventory().AtSlot(dnum)
	if sitem != nil {
		ctx.Character.Inventory().SetSlot(sitem.(*ItemInstance).ID(), dnum)
	}
	if ditem != nil {
		ctx.Character.Inventory().SetSlot(ditem.(*ItemInstance).ID(), snum)
	}

	ctx.Player.client.SyncInventory()
}

func handleAutoLoginCommand(ctx *CommandContext) {
	ctx.Player.client.ToggleAutologin()
}

func handleChannelListCommand(ctx *CommandContext) {
	rows := []string{TableRow(
		TableCell{content: "Channel", header: true},
		TableCell{content: "Description", header: true},
		TableCell{content: "Joined", header: true},
	)}

	for _, c := range Armeria.channels {
		if c.HasPermission(ctx.Character) {
			if ctx.Character.InChannel(c) {
				rows = append(rows, TableRow(
					TableCell{content: TextStyle(c.Name, TextStyleBold)},
					TableCell{content: c.Description},
					TableCell{content: "Yes"},
				))
			} else {
				rows = append(rows, TableRow(
					TableCell{content: TextStyle(c.Name, TextStyleBold)},
					TableCell{content: c.Description},
					TableCell{content: ""},
				))
			}
		}
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleChannelJoinCommand(ctx *CommandContext) {
	channelName := ctx.Args["channel"]

	ch := ChannelByName(channelName)
	if ch == nil {
		ctx.Player.client.ShowColorizedText("You must enter a valid channel name to join.", ColorError)
		return
	}

	if ctx.Character.InChannel(ch) {
		ctx.Player.client.ShowColorizedText("You are already in that channel!", ColorError)
		return
	}

	ctx.Character.JoinChannel(ch)
	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You joined the %s channel. You can use %s to communicate.",
			TextStyle(ch.Name, TextStyleBold),
			TextStyle(ch.SlashCommand, TextStyleBold),
		),
		ColorSuccess,
	)
}

func handleChannelSayCommand(ctx *CommandContext) {
	channelName := ctx.Args["channel"]
	sayText := ctx.Args["text"]

	ch := ChannelByName(channelName)
	if ch == nil {
		ctx.Player.client.ShowColorizedText("You must enter a valid channel name to talk to.", ColorError)
		return
	}

	if !ctx.Character.InChannel(ch) {
		ctx.Player.client.ShowColorizedText("You are not in that channel.", ColorError)
		return
	}

	ch.Broadcast(ctx.Character, sayText)
}

func handleSettingsCommand(ctx *CommandContext) {
	setting := strings.ToLower(ctx.Args["name"])

	if !misc.Contains(ValidSettings(), setting) {
		if setting != "" {
			ctx.Player.client.ShowColorizedText(
				fmt.Sprintf(
					"%s is not a valid setting name. Please check available settings below:",
					TextStyle(setting, TextStyleBold),
				),
				ColorError,
			)
		}

		rows := []string{TableRow(
			TableCell{content: "Name", header: true},
			TableCell{content: "Description", header: true},
			TableCell{content: "Default", header: true},
			TableCell{content: "Current", header: true},
		)}

		valid := ValidSettings()

		for _, s := range valid {
			rows = append(rows, TableRow(
				TableCell{content: s},
				TableCell{content: SettingDesc(s), styling: "padding:0px 2px"},
				TableCell{content: SettingDefault(s), styling: "padding:0px 2px"},
				TableCell{content: ctx.Character.Setting(s), styling: "padding:0px 2px"},
			))
		}
		ctx.Player.client.ShowText(TextTable(rows...))
		return
	}

	var newValue string
	if ctx.Character.Setting(setting) == "true" || ctx.Character.Setting(setting) == "false" {
		newValue = misc.ToggleStringBool(ctx.Character.Setting(setting))
	} else {
		newValue = ctx.Args["value"]
	}

	_ = ctx.Character.SetSetting(setting, newValue)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("Setting %s has been set to '%s'.", TextStyle(setting, TextStyleBold), newValue),
		ColorSuccess,
	)
}

func handleBugCommand(ctx *CommandContext) {
	bug := ctx.Args["bug"]

	issueBody := fmt.Sprintf(
		"This was reported in-game by **%s**.\n\n"+
			"**Location:** %s\n\n%s",
		ctx.Character.Name(),
		ctx.Character.Room().LocationString(),
		bug,
	)

	issue, err := Armeria.github.CreateIssue(ctx.Character.Name(), issueBody, bug)
	if err != nil {
		Armeria.log.Error(
			"error submitting bug to github repo",
			zap.Error(err),
		)
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf(
				"There was an error submitting the bug report to GitHub.\nAs an alternative, you can manually "+
					"submit a bug report here: %s",
				TextStyle("https://github.com/heyitsmdr/armeria/issues/new", TextStyleLink),
			),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(
		fmt.Sprintf("Thank you for your submission! You can view/track it here: %s",
			TextStyle(issue.GetHTMLURL(), TextStyleLink),
		),
	)
}

func handleGiveCommand(ctx *CommandContext) {
	target := ctx.Args["target"]
	item := ctx.Args["item"]

	ctr := ctx.Character.Room().Here()
	tobj, _, trt := ctr.GetByName(target)
	if trt == RegistryTypeUnknown {
		ctx.Player.client.ShowColorizedText(CommonTargetNotFoundHere, ColorError)
		return
	}

	iobj, _, irt := ctx.Character.Inventory().GetByName(item)
	if irt == RegistryTypeUnknown {
		ctx.Player.client.ShowColorizedText(CommonItemNotFoundOnCharacter, ColorError)
		return
	}

	if trt != RegistryTypeCharacter && trt != RegistryTypeMobInstance {
		ctx.Player.client.ShowColorizedText("You can only give things to other characters or mobs!", ColorError)
		return
	} else if tobj.(ContainerObject).ID() == ctx.Character.ID() {
		ctx.Player.client.ShowColorizedText("You cannot give things to yourself.", ColorError)
		return
	}

	// set the target object container
	var toc *ObjectContainer
	if trt == RegistryTypeCharacter {
		toc = tobj.(*Character).Inventory()
	} else {
		toc = tobj.(*MobInstance).Inventory()
	}

	if toc.MaxSize() > 0 && toc.Count() >= toc.MaxSize() {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf(
				"%s does not have enough room to hold that!",
				tobj.(ContainerObject).FormattedName(),
			),
			ColorError,
		)
		return
	}

	ii := iobj.(*ItemInstance)
	tco := tobj.(ContainerObject)

	// remove item from source
	ctx.Character.Inventory().Remove(ii.ID())

	// add item to target
	_ = toc.Add(ii.ID())

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf(
			"You gave %s a %s.",
			tco.FormattedName(),
			ii.FormattedName(),
		),
		ColorSuccess,
	)
	ctx.Player.client.SyncInventory()

	if trt == RegistryTypeCharacter {
		tobj.(*Character).Player().client.ShowText(
			fmt.Sprintf(
				"%s gave you a %s.",
				ctx.Character.FormattedName(),
				ii.FormattedName(),
			),
		)
		tobj.(*Character).Player().client.SyncInventory()
	} else if trt == RegistryTypeMobInstance {
		go CallMobFunc(
			ctx.Character,
			tobj.(*MobInstance),
			"received_item",
			lua.LString(ctx.Character.ID()),
			lua.LString(ii.ID()),
		)
	}

	roomExceptions := []*Character{ctx.Character}
	if trt == RegistryTypeCharacter {
		roomExceptions = append(roomExceptions, tobj.(*Character))
	}
	for _, c := range ctx.Character.Room().Here().Characters(true, roomExceptions...) {
		c.Player().client.ShowText(
			fmt.Sprintf(
				"%s gave %s something.",
				ctx.Character.FormattedName(),
				tco.FormattedName(),
			),
		)
	}
}
