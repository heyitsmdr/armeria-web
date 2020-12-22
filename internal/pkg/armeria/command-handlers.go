package armeria

import (
	"armeria/internal/pkg/misc"
	"armeria/internal/pkg/sfx"
	"armeria/internal/pkg/validate"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/muesli/reflow/wordwrap"
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
			fmt.Sprintf("You take a look at %s.\n%s", TextStyle(co.FormattedName(), WithBold()), lookResult),
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
							TextStyle(co.FormattedName(), WithBold()),
						),
					)
				}
			}
		}
		return
	}

	ar := r.AdjacentRooms()
	var validDirs []string
	if ar.North != nil {
		validDirs = append(validDirs, TextStyle("north", WithBold()))
	}
	if ar.South != nil {
		validDirs = append(validDirs, TextStyle("south", WithBold()))
	}
	if ar.East != nil {
		validDirs = append(validDirs, TextStyle("east", WithBold()))
	}
	if ar.West != nil {
		validDirs = append(validDirs, TextStyle("west", WithBold()))
	}
	if ar.Up != nil {
		validDirs = append(validDirs, TextStyle("up", WithBold()))
	}
	if ar.Down != nil {
		validDirs = append(validDirs, TextStyle("down", WithBold()))
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

	wrapDescAt, err := strconv.Atoi(ctx.Character.Setting(SettingWrap))
	if err != nil {
		log.Fatalf("error converting wrap setting to int: %s", err)
	}

	ctx.Player.client.ShowText(
		TextStyle(r.Attribute(AttributeTitle), WithBold(), WithSize(14), WithUserColor(ctx.Character, ColorRoomTitle)) + "\n" +
			wordwrap.String(r.Attribute(AttributeDescription), wrapDescAt) +
			TextStyle(validDirString, WithUserColor(ctx.Character, ColorRoomDirs)),
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
		validDirs = append(validDirs, TextStyle("north", WithBold()))
	}
	if ar.South != nil {
		validDirs = append(validDirs, TextStyle("south", WithBold()))
	}
	if ar.East != nil {
		validDirs = append(validDirs, TextStyle("east", WithBold()))
	}
	if ar.West != nil {
		validDirs = append(validDirs, TextStyle("west", WithBold()))
	}
	if ar.Up != nil {
		validDirs = append(validDirs, TextStyle("up", WithBold()))
	}
	if ar.Down != nil {
		validDirs = append(validDirs, TextStyle("down", WithBold()))
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
			lua.LString(ctx.Args["text"]),
		)
	}
}

func handleMoveCommand(ctx *CommandContext) {
	d := ctx.Args["direction"]

	normDir := misc.NormalizeDirection(d)
	if len(normDir) == 0 {
		ctx.Player.client.ShowColorizedText("That's not a valid direction to move in.", ColorError)
		return
	}

	newRoom := ctx.Character.Room().ConnectedRoom(normDir)
	if newRoom == nil {
		currentRoomAttr := ctx.Character.Room().Attribute(normDir)
		if len(currentRoomAttr) > 0 && currentRoomAttr[0:1] == "!" {
			if len(currentRoomAttr) > 1 {
				ctx.Player.client.ShowColorizedText(currentRoomAttr[1:], ColorError)
			} else {
				ctx.Player.client.ShowColorizedText(CommonInvalidDirection, ColorError)
			}
		} else {
			ctx.Player.client.ShowColorizedText(CommonInvalidDirection, ColorError)
		}
		return
	}

	moveAllowed, moveError := ctx.Character.MoveAllowed(newRoom)
	if !moveAllowed {
		ctx.Player.client.ShowColorizedText(moveError, ColorError)
		return
	}

	oldAreaUUID := ctx.Character.Room().ParentArea.ID()
	ctx.Character.Move(
		newRoom,
		TextStyle(fmt.Sprintf("You walk %s.", misc.MoveToStringFromDir("to the", normDir)), WithUserColor(ctx.Character, ColorMovement)),
		TextStyle(fmt.Sprintf("%s walks %s.", ctx.Character.FormattedName(), misc.MoveToStringFromDir("to the", normDir)), WithUserColor(ctx.Character, ColorMovement)),
		TextStyle(fmt.Sprintf("%s walked in from %s.", ctx.Character.FormattedName(), misc.MoveFromStringFromDir("the", misc.OppositeDirection(normDir))), WithUserColor(ctx.Character, ColorMovement)),
		"",
	)

	if newRoom.ParentArea.ID() != oldAreaUUID {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("You've just entered %s.", TextStyle(newRoom.ParentArea.Name(), WithBold())),
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
	attr := AttributeCasing(ctx.Args["property"])
	if !misc.Contains(AttributeList(ObjectTypeRoom), attr) {
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
		fmt.Sprintf("You modified the %s property of the room (%s).", TextStyle(attr, WithBold()), ta),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(tr.EditorData())
	}
}

func handleRoomMoveCommand(ctx *CommandContext) {
	dir := strings.ToLower(ctx.Args["direction"])

	if dir == "up" || dir == "down" {
		ctx.Player.client.ShowColorizedText("Rooms cannot be moved up or down.", ColorError)
		return
	}

	oppositeDir := misc.OppositeDirection(dir)
	if len(oppositeDir) == 0 {
		ctx.Player.client.ShowColorizedText("That's not a valid direction to move a room to.", ColorError)
		return
	}

	rm := ctx.Character.Room()
	offsets := misc.DirectionOffsets(dir)
	newCoords := &Coords{
		UnsafeX: rm.Coords.X() + offsets["x"],
		UnsafeY: rm.Coords.Y() + offsets["y"],
		UnsafeZ: rm.Coords.Z(),
		UnsafeI: rm.Coords.I(),
	}

	// Check if there is already a room at the new coords.
	if rm.ParentArea.RoomAt(newCoords) != nil {
		ctx.Player.client.ShowColorizedText("There's already a room at the intended new location.", ColorError)
		return
	}

	// Move the room.
	rm.Coords.SetFrom(newCoords)

	// Link the rooms (if applicable).
	oppositeRm := rm.ConnectedRoom(oppositeDir)
	if oppositeRm != nil {
		oppositeRm.SetAttribute(dir, rm.Coords.String())
		rm.SetAttribute(oppositeDir, oppositeRm.Coords.String())
	}

	// Sync the minimap for anyone in the area.
	for _, char := range rm.ParentArea.Characters() {
		char.Player().client.SyncMap()
		if char.Room() == rm {
			char.Player().client.SyncMapLocation()
		}
	}

	ctx.Player.client.ShowColorizedText("The room has been moved.", ColorSuccess)
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

	rm := Armeria.worldManager.CreateRoom(ctx.Character.Room().ParentArea, c)

	// Match room colors.
	rm.SetAttribute(AttributeColor, ctx.Character.Room().Attribute(AttributeColor))

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
	ctx.Player.client.SyncMoney()
	ctx.Player.client.SyncSettings()
	ctx.Player.client.SyncCommands()
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

	rows := []string{TableRow(
		TableCell{content: "Character", header: true},
	)}

	for _, c := range Armeria.characterManager.Characters() {
		if len(f) == 0 || strings.Contains(strings.ToLower(c.Name()), strings.ToLower(f)) {
			rows = append(rows, TableRow(
				TableCell{content: c.Name()},
			))
		}
	}

	if len(f) > 0 && len(rows) == 1 {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("There are no characters matching \"%s\".", f),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleCharacterCreateCommand(ctx *CommandContext) {
	charName := ctx.Args["character"]
	charPass := ctx.Args["password"]

	if char := Armeria.characterManager.CharacterByName(charName); char != nil {
		ctx.Player.client.ShowColorizedText("A character with that name already exists.", ColorError)
		return
	}

	Armeria.characterManager.CreateCharacter(charName, charPass)

	ctx.Player.client.ShowColorizedText("The character has been created!", ColorSuccess)
}

func handleCharacterSetCommand(ctx *CommandContext) {
	char := ctx.Args["character"]
	attr := ctx.Args["property"]
	val := ctx.Args["value"]

	c := Armeria.characterManager.CharacterByName(char)
	if c == nil {
		ctx.Player.client.ShowColorizedText("That character doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(AttributeList(ObjectTypeCharacter), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid character attribute.", ColorError)
		return
	}

	if len(val) > 0 {
		valid := AttributeValidate(ObjectTypeCharacter, attr, val)
		if !valid.Result {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", valid), ColorError)
			return
		}
	}

	_ = c.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the character %s.", TextStyle(attr, WithBold()), c.FormattedName()),
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

	rows := []string{TableRow(
		TableCell{content: "Mob", header: true},
		TableCell{content: "Instances", header: true},
	)}

	for _, m := range Armeria.mobManager.Mobs() {
		if len(f) == 0 || strings.Contains(strings.ToLower(m.Name()), strings.ToLower(f)) {
			rows = append(rows, TableRow(
				TableCell{
					content: TextStyle(m.Name(), WithLinkCmd("/mob edit "+m.Name())),
				},
				TableCell{
					content: TextStyle(
						fmt.Sprintf("%d instances", len(m.Instances())),
						WithLinkCmd(fmt.Sprintf("/mob instances %s", m.Name())),
					),
				},
			))
		}
	}

	if len(f) > 0 && len(rows) == 1 {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("There are no mobs matching \"%s\".", f),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(TextTable(rows...))
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
		fmt.Sprintf("A mob named %s has been created.", TextStyle(n, WithBold())),
		ColorSuccess,
	)
}

func handleMobDeleteCommand(ctx *CommandContext) {
	n := ctx.Args["name"]

	mob := Armeria.mobManager.MobByName(n)
	if mob == nil {
		ctx.Player.client.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	if len(mob.Instances()) > 0 {
		ctx.Player.client.ShowColorizedText("That mob has one or more instances in the game world and cannot be deleted.", ColorError)
		return
	}

	Armeria.mobManager.RemoveMob(mob)

	ctx.Player.client.ShowColorizedText("The mob has been removed from the game.", ColorSuccess)
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
	mob := ctx.Args["mob"]
	attr := AttributeCasing(ctx.Args["property"])
	val := ctx.Args["value"]

	m := Armeria.mobManager.MobByName(mob)
	if m == nil {
		ctx.Player.client.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(AttributeList(ObjectTypeMob), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid mob attribute.", ColorError)
		return
	}

	if len(val) > 0 {
		valid := AttributeValidate(ObjectTypeMob, attr, val)
		if !valid.Result {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", valid), ColorError)
			return
		}
	}

	m.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the mob %s.",
			TextStyle(attr, WithBold()),
			TextStyle(m.UnsafeName, WithBold()),
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
	attr := AttributeCasing(ctx.Args["property"])
	val := ctx.Args["value"]

	if !misc.Contains(AttributeList(ObjectTypeMob), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid mob attribute.", ColorError)
		return
	} else if !misc.Contains(AttributeList(ObjectTypeMobInstance), attr) {
		ctx.Player.client.ShowColorizedText("That's not an attribute you can set on a mob instance.", ColorError)
		return
	}

	if len(val) > 0 {
		valid := AttributeValidate(ObjectTypeMob, attr, val)
		if !valid.Result {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", valid), ColorError)
			return
		}
	}

	_ = mi.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the mob instace %s (%s).",
			TextStyle(attr, WithBold()),
			TextStyle(mi.Name(), WithBold()),
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
			TableCell{
				content: TextStyle(
					fmt.Sprintf("%s (%s)", mi.Room().LocationString(), mi.Room().Attribute("title")),
					WithLinkCmd(fmt.Sprintf("/tp %s", mi.Room().LocationString())),
				),
			},
		))
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleWipeCommand(ctx *CommandContext) {
	filter := ctx.Args["filter"]
	matches := 0

	for _, o := range ctx.Character.Room().Here().All() {
		obj := o.(ContainerObject)

		// using a filter?
		if len(filter) > 0 && !strings.Contains(strings.ToLower(obj.Name()), strings.ToLower(filter)) {
			continue
		}

		switch obj.Type() {
		case ContainerObjectTypeMob:
			m := Armeria.mobManager.MobByName(obj.Name())
			ctx.Character.Room().Here().Remove(obj.ID())
			if m != nil {
				m.DeleteInstance(obj.(*MobInstance))
				matches = matches + 1
			}
		case ContainerObjectTypeItem:
			i := Armeria.itemManager.ItemByName(obj.Name())
			ctx.Character.Room().Here().Remove(obj.ID())
			if i != nil {
				i.DeleteInstance(obj.(*ItemInstance))
				matches = matches + 1
			}
		}
	}

	if len(filter) > 0 && matches == 0 {
		ctx.Player.client.ShowColorizedText("The filter did not match anything in the room.", ColorError)
		return
	}

	for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
		c.Player().client.ShowText(
			fmt.Sprintf("%s wiped one or more things from the room.", ctx.Character.FormattedName()),
		)
		c.Player().client.SyncRoomObjects()
	}

	ctx.Player.client.ShowColorizedText(fmt.Sprintf("You wiped %d things from the room.", matches), ColorSuccess)
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
		fmt.Sprintf("An item named %s has been created.", TextStyle(n, WithBold())),
		ColorSuccess,
	)
}

func handleItemDeleteCommand(ctx *CommandContext) {
	n := ctx.Args["name"]

	item := Armeria.itemManager.ItemByName(n)
	if item == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist.", ColorError)
		return
	}

	if len(item.Instances()) > 0 {
		ctx.Player.client.ShowColorizedText("That item has one or more instances in the game world and cannot be deleted.", ColorError)
		return
	}

	Armeria.itemManager.RemoveItem(item)

	ctx.Player.client.ShowColorizedText("The item has been removed from the game.", ColorSuccess)
}

func handleItemListCommand(ctx *CommandContext) {
	f := ctx.Args["filter"]

	rows := []string{TableRow(
		TableCell{content: "Item", header: true},
		TableCell{content: "Instances", header: true},
		TableCell{content: "Type", header: true},
	)}

	for _, i := range Armeria.itemManager.Items() {
		if len(f) == 0 || strings.Contains(strings.ToLower(i.Name()), strings.ToLower(f)) {
			rows = append(rows, TableRow(
				TableCell{
					content: TextStyle(
						i.Name(),
						WithLinkCmd(fmt.Sprintf("/item edit %s", i.Name())),
					),
				},
				TableCell{
					content: TextStyle(
						fmt.Sprintf("x%d", len(i.Instances())),
						WithLinkCmd(fmt.Sprintf("/item instances %s", i.Name())),
					),
				},
				TableCell{
					content: i.Attribute(AttributeType),
				},
			))
		}
	}

	if len(f) > 0 && len(rows) == 1 {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("There are no items matching \"%s\".", f),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(TextTable(rows...))
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
	item := ctx.Args["item"]
	attr := AttributeCasing(ctx.Args["property"])
	val := ctx.Args["value"]

	i := Armeria.itemManager.ItemByName(item)
	if i == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist.", ColorError)
		return
	}

	if !misc.Contains(AttributeList(ObjectTypeItem), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid item attribute.", ColorError)
		return
	}

	if len(val) > 0 {
		valid := AttributeValidate(ObjectTypeItem, attr, val)
		if !valid.Result {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", valid), ColorError)
			return
		}
	}

	i.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the item %s.",
			TextStyle(attr, WithBold()),
			TextStyle(i.Name(), WithBold()),
		),
		ColorSuccess,
	)

	// Sync possible locations of the item being edited.
	ctx.Player.client.SyncInventory()
	ctx.Player.client.SyncRoomObjects()

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
	attr := AttributeCasing(ctx.Args["property"])
	val := ctx.Args["value"]

	if !misc.Contains(AttributeList(ObjectTypeItem), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid item attribute.", ColorError)
		return
	} else if !misc.Contains(AttributeList(ObjectTypeItemInstance), attr) {
		ctx.Player.client.ShowColorizedText("That's not an attribute you can set on an item instance.", ColorError)
		return
	}

	if len(val) > 0 {
		valid := AttributeValidate(ObjectTypeItem, attr, val)
		if !valid.Result {
			ctx.Player.client.ShowColorizedText(fmt.Sprintf("The attribute value could not be validated: %s.", valid), ColorError)
			return
		}
	}

	_ = ii.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the %s property of the item instace %s (%s).",
			TextStyle(attr, WithBold()),
			TextStyle(ii.Name(), WithBold()),
			ii.ID(),
		),
		ColorSuccess,
	)

	// Sync possible locations of the item being edited.
	ctx.Player.client.SyncInventory()
	ctx.Player.client.SyncRoomObjects()

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
				TableCell{
					content: TextStyle(
						fmt.Sprintf("%s (%s)", ii.Room().LocationString(), ii.Room().Attribute("title")),
						WithLinkCmd(fmt.Sprintf("/tp %s", ii.Room().LocationString())),
					),
				},
			))
		} else if ctr.ParentType() == ContainerParentTypeCharacter {
			rows = append(rows, TableRow(
				TableCell{content: ii.FormattedName()},
				TableCell{content: ii.ID()},
				TableCell{
					content: fmt.Sprintf(
						"Character: %s (slot %d)",
						ii.Character().FormattedName(),
						ii.Character().Inventory().Slot(ii.ID()),
					),
					styling: "",
				},
			))
		} else if ctr.ParentType() == ContainerParentTypeMobInstance {
			rows = append(rows, TableRow(
				TableCell{content: ii.FormattedName()},
				TableCell{content: ii.ID()},
				TableCell{
					content: fmt.Sprintf("Mob: %s (%s)", ii.MobInstance().FormattedName(), ii.MobInstance().ID()),
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
		fmt.Sprintf("An area named %s has been created!", TextStyle(a.Name(), WithBold())),
		ColorSuccess,
	)
}

func handleAreaListCommand(ctx *CommandContext) {
	f := ctx.Args["filter"]

	rows := []string{TableRow(
		TableCell{content: "Area", header: true},
		TableCell{content: "Rooms", header: true},
	)}

	for _, a := range Armeria.worldManager.Areas() {
		if len(f) == 0 || strings.Contains(strings.ToLower(a.Name()), strings.ToLower(f)) {
			rows = append(rows, TableRow(
				TableCell{content: a.Name()},
				TableCell{content: fmt.Sprintf("%d rooms", len(a.Rooms()))},
			))
		}
	}

	if len(f) > 0 && len(rows) == 1 {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("There are no areas matching \"%s\".", f),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(TextTable(rows...))
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

	charToMove := ctx.Character
	var destination *Room
	var moveMsg string
	if t[0:2] == "@@" {
		cn := t[2:]
		c := Armeria.characterManager.CharacterByName(cn)
		if c == nil {
			ctx.Player.client.ShowColorizedText("There is no character by that name.", ColorError)
			return
		}
		charToMove = c
		destination = ctx.Character.Room()
		moveMsg = fmt.Sprintf("You were teleported far away by %s!", ctx.Character.FormattedName())
		ctx.Player.client.ShowColorizedText(fmt.Sprintf("You summoned %s here.", c.FormattedName()), ColorMovement)
	} else if t[0:1] == "@" {
		cn := t[1:]
		c := Armeria.characterManager.CharacterByName(cn)
		if c == nil {
			ctx.Player.client.ShowColorizedText("There is no character by that name.", ColorError)
			return
		} else if !c.Online() {
			ctx.Player.client.ShowColorizedText("That character is not online.", ColorError)
			return
		}

		destination = c.Room()
		moveMsg = fmt.Sprintf("You teleported to %s.", c.FormattedName())
	} else {
		loc := strings.Split(t, ",")
		if len(loc) == 1 {
			loc = append(loc, "0", "0", "0")
		} else if len(loc) != 4 {
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
		moveMsg = fmt.Sprintf("You teleported to %s at %d, %d, %d.", TextStyle(a.Name(), WithBold()), x, y, z)
	}

	if destination == nil {
		ctx.Player.client.ShowColorizedText("You cannot teleport there!", ColorError)
		return
	}

	charToMove.Move(
		destination,
		TextStyle(moveMsg, WithUserColor(charToMove, ColorMovement)),
		TextStyle(fmt.Sprintf("%s teleported away!", charToMove.FormattedName()), WithUserColor(charToMove, ColorMovement)),
		TextStyle(fmt.Sprintf("%s teleported here!", charToMove.FormattedName()), WithUserColor(charToMove, ColorMovement)),
		sfx.Teleport,
	)

	if charToMove.Online() {
		Armeria.commandManager.ProcessCommand(charToMove.Player(), "look", false)
	}
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
			TableCell{content: TextStyle("/"+cmd.Name, WithBold()), styling: "padding:0px 2px"},
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
			attrs = AttributeList(ObjectTypeRoom)
		}
		for _, attr := range attrs {
			if !misc.Contains(AttributeList(ObjectTypeRoom), attr) {
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
		ctx.Player.client.SyncMap()
		ctx.Player.client.SyncRoomTitle()
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
	searchString := ctx.Args["item"]

	roomObjects := ctx.Character.Room().Here()
	o, _, rt := roomObjects.GetByAny(searchString)
	if o == nil {
		ctx.Player.client.ShowColorizedText("There is nothing here by that name.", ColorError)
		return
	} else if rt != RegistryTypeItemInstance {
		ctx.Player.client.ShowColorizedText("You cannot pick that up.", ColorError)
		return
	}

	item := o.(*ItemInstance)

	if item.Attribute(AttributeHoldable) == "false" {
		ctx.Player.client.ShowColorizedText("You are not able to pick that up.", ColorError)
		return
	}

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
	ctx.Player.client.PlaySFX(sfx.PickupItem)
	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You picked up a %s.", item.FormattedName()),
		ColorSuccess,
	)

	for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
		c.Player().client.SyncRoomObjects()
		c.Player().client.ShowText(
			fmt.Sprintf("%s picked up a %s.", ctx.Character.FormattedName(), item.FormattedName()),
		)
	}
}

func handleDropCommand(ctx *CommandContext) {
	searchString := ctx.Args["item"]

	i, _, rt := ctx.Character.Inventory().GetByAny(searchString)
	if rt == RegistryTypeUnknown {
		ctx.Player.client.ShowColorizedText("You don't have that item in your inventory.", ColorError)
		return
	}

	item := i.(*ItemInstance)

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
		c.Player().client.ShowText(
			fmt.Sprintf("%s dropped a %s.", ctx.Character.FormattedName(), item.FormattedName()),
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
					TableCell{content: TextStyle(c.Name, WithBold())},
					TableCell{content: c.Description},
					TableCell{content: TextStyle("Yes", WithUserColor(ctx.Character, ColorSuccess))},
				))
			} else {
				rows = append(rows, TableRow(
					TableCell{content: TextStyle(c.Name, WithBold())},
					TableCell{content: c.Description},
					TableCell{content: TextStyle("No", WithUserColor(ctx.Character, ColorError))},
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
			TextStyle(ch.Name, WithBold()),
			TextStyle(ch.SlashCommand, WithBold()),
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
	value := ctx.Args["value"]

	if len(setting) == 0 {
		rows := []string{TableRow(
			TableCell{content: "Name", header: true},
			TableCell{content: "Description", header: true},
			TableCell{content: "Current", header: true},
			TableCell{content: "Default", header: true},
		)}

		valid := ValidSettings()

		for _, s := range valid {
			reqPerm := SettingPermission(s)
			if len(reqPerm) > 0 && !ctx.Character.HasPermission(reqPerm) {
				continue
			}

			rows = append(rows, TableRow(
				TableCell{content: s},
				TableCell{content: SettingDesc(s), styling: "padding:0px 2px"},
				TableCell{content: ctx.Character.Setting(s), styling: "padding:0px 2px"},
				TableCell{content: SettingDefault(s), styling: "padding:0px 2px;color:#666"},
			))
		}
		ctx.Player.client.ShowText(TextTable(rows...))
		return
	} else if !misc.Contains(ValidSettings(), setting) {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf(
				"%s is not a valid setting name.",
				TextStyle(setting, WithBold()),
			),
			ColorError,
		)
		return
	}

	// Check if the character has permission to modify the setting.
	reqPerm := SettingPermission(setting)
	if len(reqPerm) > 0 && !ctx.Character.HasPermission(reqPerm) {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf(
				"%s is not a valid setting name.",
				TextStyle(setting, WithBold()),
			),
			ColorError,
		)
		return
	}

	// If no value is specified, either swap the bools or set the default value.
	if len(value) == 0 {
		if misc.IsStringBool(ctx.Character.Setting(setting)) {
			value = misc.ToggleStringBool(ctx.Character.Setting(setting))
		} else {
			value = SettingDefault(setting)
		}
	} else {
		// Validate the setting value.
		valid := validate.Check(value, SettingValidationString(setting))
		if !valid.Result {
			ctx.Player.client.ShowColorizedText(
				fmt.Sprintf("You cannot use that value due to: %s", valid),
				ColorError,
			)
			return
		}
	}

	if err := ctx.Character.SetSetting(setting, value); err != nil {
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf(
				"That setting cannot be changed: %s",
				err,
			),
			ColorError,
		)
	}

	ctx.Player.client.SyncSettings()
	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("Setting %s has been set to '%s'.", TextStyle(setting, WithBold()), value),
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
					"submit a bug report %s.",
				TextStyle("here", WithLink("https://github.com/heyitsmdr/armeria/issues/new")),
			),
			ColorError,
		)
		return
	}

	ctx.Player.client.ShowText(
		fmt.Sprintf("Thank you for your submission! You can view/track it %s.",
			TextStyle("here", WithLink(issue.GetHTMLURL())),
		),
	)
}

func handleGiveCommand(ctx *CommandContext) {
	target := ctx.Args["target"]
	item := ctx.Args["item"]

	ctr := ctx.Character.Room().Here()
	tobj, _, trt := ctr.GetByAny(target)
	if trt == RegistryTypeUnknown {
		ctx.Player.client.ShowColorizedText(CommonTargetNotFoundHere, ColorError)
		return
	}

	iobj, _, irt := ctx.Character.Inventory().GetByAny(item)
	if irt != RegistryTypeItemInstance {
		ctx.Player.client.ShowColorizedText(CommonItemNotFoundOnCharacter, ColorError)
		return
	}

	if trt == RegistryTypeItemInstance && tobj.(*ItemInstance).Attribute(AttributeType) == ItemTypeTrashCan {
		// Destroy the item.
		ii := iobj.(*ItemInstance)
		ctx.Character.Inventory().Remove(ii.ID())
		ii.Delete()
		ctx.Player.client.ShowColorizedText(
			fmt.Sprintf("You put a %s into the %s. Goodbye!",
				ii.FormattedName(),
				tobj.(*ItemInstance).FormattedName(),
			),
			ColorSuccess,
		)
		ctx.Player.client.SyncInventory()
		for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
			c.Player().client.ShowText(
				fmt.Sprintf(
					"%s put an item into the %s.",
					ctx.Character.FormattedName(),
					tobj.(*ItemInstance).FormattedName(),
				),
			)
		}
		return
	} else if trt != RegistryTypeCharacter && trt != RegistryTypeMobInstance {
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

	// check if the target object container can hold it
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

	// check if the mob can handle it?
	if trt == RegistryTypeMobInstance {
		if !misc.Contains(toc.ParentMobInstance().Parent.ScriptFuncs(), "received_item") {
			ctx.Player.client.ShowColorizedText(
				fmt.Sprintf(
					"%s does not want that.",
					toc.ParentMobInstance().FormattedName(),
				),
				ColorError,
			)
			return
		}
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

func handleEmoteCommand(ctx *CommandContext) {
	emotion := ctx.Args["emote"]

	if emotion[len(emotion)-1:] == "." {
		emotion = emotion[:len(emotion)-1]
	}

	for _, c := range ctx.Character.Room().Here().Characters(true) {
		c.Player().client.ShowText(
			fmt.Sprintf("%s %s.", ctx.Character.FormattedName(), emotion),
		)
	}
}

func handleLedgerListCommand(ctx *CommandContext) {
	rows := []string{TableRow(
		TableCell{content: "Ledger", header: true},
		TableCell{content: "Items", header: true},
	)}

	for _, l := range Armeria.ledgerManager.Ledgers() {
		rows = append(rows, TableRow(
			TableCell{content: fmt.Sprintf("[cmd=/ledger show %[1]s]%[1]s[/cmd]", l.Name())},
			TableCell{content: fmt.Sprintf("%d items", len(l.Entries()))},
		))
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleLedgerCreateCommand(ctx *CommandContext) {
	name := ctx.Args["name"]

	if strings.Contains(name, " ") {
		ctx.Player.client.ShowColorizedText("The ledger name cannot contain a space.", ColorError)
		return
	}

	exists := Armeria.ledgerManager.LedgerByName(name)
	if exists != nil {
		ctx.Player.client.ShowColorizedText("A ledger already exists with that name.", ColorError)
		return
	}

	l := Armeria.ledgerManager.CreateLedger(name)
	Armeria.ledgerManager.AddLedger(l)

	ctx.Player.client.ShowColorizedText("The ledger has been created.", ColorSuccess)
}

func handleLedgerRenameCommand(ctx *CommandContext) {
	ledgerName := ctx.Args["ledger_name"]
	newName := ctx.Args["new_name"]

	ledger := Armeria.ledgerManager.LedgerByName(ledgerName)
	if ledger == nil {
		ctx.Player.client.ShowColorizedText("A ledger by that name doesn't exist.", ColorError)
		return
	}

	if strings.Contains(newName, " ") {
		ctx.Player.client.ShowColorizedText("The new ledger name cannot contain a space.", ColorError)
		return
	}

	ledger.SetName(newName)

	ctx.Player.client.ShowColorizedText("The ledger has been renamed.", ColorSuccess)
}

func handleLedgerAddCommand(ctx *CommandContext) {
	ledgerName := ctx.Args["ledger_name"]
	itemName := ctx.Args["item_name"]

	ledger := Armeria.ledgerManager.LedgerByName(ledgerName)
	if ledger == nil {
		ctx.Player.client.ShowColorizedText("A ledger by that name doesn't exist.", ColorError)
		return
	}

	item := Armeria.itemManager.ItemByName(itemName)
	if item == nil {
		ctx.Player.client.ShowColorizedText("An item by that name doesn't exist.", ColorError)
		return
	}

	if ledger.Contains(item.Name()) != nil {
		ctx.Player.client.ShowColorizedText("That item is already on that ledger.", ColorError)
		return
	}

	ledger.AddEntry(&LedgerEntry{
		ItemName:  item.Name(),
		BuyPrice:  0.00,
		SellPrice: 0.00,
	})

	ctx.Player.client.ShowColorizedText("Entry has been added to the ledger.", ColorSuccess)
}

func handleLedgerRemoveCommand(ctx *CommandContext) {
	ledgerName := ctx.Args["ledger_name"]
	itemName := ctx.Args["item_name"]

	ledger := Armeria.ledgerManager.LedgerByName(ledgerName)
	if ledger == nil {
		ctx.Player.client.ShowColorizedText("A ledger by that name doesn't exist.", ColorError)
		return
	}

	entry := ledger.Contains(itemName)
	if entry == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist on that ledger.", ColorError)
		return
	}

	ledger.RemoveEntry(entry)

	ctx.Player.client.ShowColorizedText("Entry has been removed from the ledger.", ColorSuccess)
}

func handleLedgerShowCommand(ctx *CommandContext) {
	ledgerName := ctx.Args["ledger_name"]

	ledger := Armeria.ledgerManager.LedgerByName(ledgerName)
	if ledger == nil {
		ctx.Player.client.ShowColorizedText("A ledger by that name doesn't exist.", ColorError)
		return
	}

	rows := []string{TableRow(
		TableCell{content: "Item", header: true},
		TableCell{content: "Buy", header: true},
		TableCell{content: "Sell", header: true},
	)}

	for _, entry := range ledger.Entries() {
		rows = append(rows, TableRow(
			TableCell{content: entry.ItemName},
			TableCell{content: misc.Money.FormatMoney(entry.BuyPrice)},
			TableCell{content: misc.Money.FormatMoney(entry.SellPrice)},
		))
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleLedgerSearchCommand(ctx *CommandContext) {
	itemName := ctx.Args["item_name"]

	matches := make(map[string][]string)

	for _, ledger := range Armeria.ledgerManager.Ledgers() {
		for _, entry := range ledger.Entries() {
			if strings.Contains(strings.ToLower(entry.ItemName), strings.ToLower(itemName)) {
				matches[ledger.Name()] = append(matches[ledger.Name()], entry.ItemName)
			}
		}
	}

	if len(matches) == 0 {
		ctx.Player.client.ShowColorizedText("No matches found across all ledgers.", ColorError)
		return
	}

	rows := []string{TableRow(
		TableCell{content: "Ledger", header: true},
		TableCell{content: "Matched", header: true},
	)}

	for ledger, matches := range matches {
		for _, item := range matches {
			rows = append(rows, TableRow(
				TableCell{content: fmt.Sprintf("[cmd=/ledger show %[1]s]%[1]s[/cmd]", ledger)},
				TableCell{content: item},
			))
		}
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleLedgerSetCommand(ctx *CommandContext) {
	buyOrSell := strings.ToLower(ctx.Args["buy_or_sell"])
	ledgerName := ctx.Args["ledger_name"]
	itemName := ctx.Args["item_name"]
	price := ctx.Args["price"]

	if buyOrSell != "buy" && buyOrSell != "sell" {
		ctx.Player.client.ShowColorizedText("You must set either a BUY or SELL price.", ColorError)
		return
	}

	ledger := Armeria.ledgerManager.LedgerByName(ledgerName)
	if ledger == nil {
		ctx.Player.client.ShowColorizedText("A ledger by that name doesn't exist.", ColorError)
		return
	}

	entry := ledger.Contains(itemName)
	if entry == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist on that ledger.", ColorError)
		return
	}

	amount, err := strconv.ParseFloat(price, 64)
	if err != nil {
		ctx.Player.client.ShowColorizedText("You must set a numerical price.", ColorError)
		return
	}

	if buyOrSell == "buy" {
		entry.BuyPrice = amount
	} else {
		entry.SellPrice = amount
	}

	ctx.Player.client.ShowColorizedText("The price has been set on the ledger.", ColorSuccess)
}

func handleBuyCommand(ctx *CommandContext) {
	mobName := ctx.Args["npc"]
	itemName := ctx.Args["item"]

	// Ensure mob is present in the room
	m, _, rt := ctx.Character.Room().Here().GetByName(mobName)
	if rt != RegistryTypeMobInstance {
		ctx.Player.client.ShowColorizedText(CommonTargetNotFoundHere, ColorError)
		return
	}
	mobInstance := m.(*MobInstance)

	// Ensure mob is aware of a ledger that contains the item
	var item *ItemInstance
	var itemLedger *LedgerEntry
	for _, ledger := range mobInstance.ItemLedgers() {
		ledgerEntry := ledger.Contains(itemName)
		if ledgerEntry != nil {
			itemLedger = ledgerEntry
			mobInstance.Inventory().PopulateFromLedger(ledger)
			if i, _, rt := mobInstance.Inventory().GetByName(ledgerEntry.ItemName); rt == RegistryTypeItemInstance {
				item = i.(*ItemInstance)
				break
			}
		}
	}
	if item == nil || itemLedger == nil || itemLedger.BuyPrice == 0 {
		ctx.Player.client.ShowColorizedText(fmt.Sprintf("%s does not have that to sell.", mobInstance.Name()), ColorError)
		return
	}

	// Ensure character has room in their inventory
	if ctx.Character.Inventory().Count() >= ctx.Character.Inventory().MaxSize() {
		ctx.Player.client.ShowColorizedText(CommonInventoryFilled, ColorError)
		return
	}

	// Remove money from character
	if !ctx.Character.RemoveMoney(itemLedger.BuyPrice) {
		ctx.Player.client.ShowColorizedText("You can't afford that.", ColorError)
		return
	}

	// Transfer the item
	mobInstance.Inventory().Remove(item.ID())
	if err := ctx.Character.Inventory().Add(item.ID()); err != nil {
		// Something went wrong, let's destroy the item instance and return the money
		item.Parent.DeleteInstance(item)
		ctx.Character.AddMoney(itemLedger.BuyPrice)
		ctx.Player.client.ShowColorizedText("Something went wrong with the transaction.", ColorError)
		return
	}

	ctx.Player.client.SyncMoney()
	ctx.Player.client.SyncInventory()
	ctx.Player.client.PlaySFX(sfx.SellBuyItem)
	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf(
			"You bought a %s from %s for %s.",
			item.FormattedName(),
			mobInstance.FormattedName(),
			ctx.Character.Colorize(misc.Money.FormatMoney(itemLedger.BuyPrice), ColorMoney),
		),
		ColorSuccess,
	)

	for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
		c.Player().client.ShowText(
			fmt.Sprintf(
				"%s bought something from %s.",
				ctx.Character.FormattedName(),
				mobInstance.FormattedName(),
			),
		)
	}
}

func handleSellCommand(ctx *CommandContext) {
	mobName := ctx.Args["npc"]
	itemName := ctx.Args["item"]

	// Ensure mob is present in the room
	m, _, rt := ctx.Character.Room().Here().GetByName(mobName)
	if rt != RegistryTypeMobInstance {
		ctx.Player.client.ShowColorizedText(CommonTargetNotFoundHere, ColorError)
		return
	}
	mobInstance := m.(*MobInstance)

	// Ensure item exists in the character's inventory
	var item *ItemInstance
	if i, _, rt := ctx.Character.Inventory().GetByAny(itemName); rt == RegistryTypeItemInstance {
		item = i.(*ItemInstance)
	}
	if item == nil {
		ctx.Player.client.ShowColorizedText(CommonItemNotFoundOnCharacter, ColorError)
		return
	}

	// Ensure mob is aware of a ledger that contains the item
	var itemLedger *LedgerEntry
	for _, ledger := range mobInstance.ItemLedgers() {
		ledgerEntry := ledger.Contains(item.Name())
		if ledgerEntry != nil {
			itemLedger = ledgerEntry
			break
		}
	}
	if itemLedger == nil || itemLedger.SellPrice == 0 {
		ctx.Player.client.ShowColorizedText(fmt.Sprintf("%s does not want that item.", mobInstance.Name()), ColorError)
		return
	}

	// Add money to the character
	ctx.Character.AddMoney(itemLedger.SellPrice)

	// Destroy the item
	ctx.Character.Inventory().Remove(item.ID())
	item.Parent.DeleteInstance(item)

	ctx.Player.client.SyncMoney()
	ctx.Player.client.SyncInventory()
	ctx.Player.client.PlaySFX(sfx.SellBuyItem)
	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf(
			"You sold a %s to %s for %s.",
			item.FormattedName(),
			mobInstance.FormattedName(),
			ctx.Character.Colorize(misc.Money.FormatMoney(itemLedger.SellPrice), ColorMoney),
		),
		ColorSuccess,
	)

	for _, c := range ctx.Character.Room().Here().Characters(true, ctx.Character) {
		c.Player().client.ShowText(
			fmt.Sprintf(
				"%s sold something to %s.",
				ctx.Character.FormattedName(),
				mobInstance.FormattedName(),
			),
		)
	}
}

func handleDestroyCommand(ctx *CommandContext) {
	searchString := ctx.Args["object"]

	if i, _, rt := ctx.Character.Inventory().GetByAny(searchString); rt == RegistryTypeItemInstance {
		item := i.(*ItemInstance)
		ctx.Character.Inventory().Remove(item.ID())
		item.Delete()
		ctx.Player.client.ShowColorizedText("The item has been destroyed!", ColorSuccess)
		ctx.Player.client.SyncInventory()
		return
	} else if i, _, rt := ctx.Character.Room().Here().GetByAny(searchString); rt == RegistryTypeItemInstance {
		item := i.(*ItemInstance)
		ctx.Character.Room().Here().Remove(item.ID())
		item.Delete()
	} else if m, _, rt := ctx.Character.Room().Here().GetByAny(searchString); rt == RegistryTypeMobInstance {
		mob := m.(*MobInstance)
		ctx.Character.Room().Here().Remove(mob.ID())
		mob.Delete()
	} else {
		ctx.Player.client.ShowColorizedText("There were no matches in the room or your inventory.", ColorError)
		return
	}

	for _, c := range ctx.Character.Room().Here().Characters(true) {
		c.Player().client.ShowText("The atmosphere around the room feels different.")
		c.Player().client.SyncRoomObjects()
	}
}

func handleTickersCommand(ctx *CommandContext) {
	rows := []string{TableRow(
		TableCell{content: "Ticker", header: true},
		TableCell{content: "Interval", header: true},
		TableCell{content: "Last Ran", header: true},
		TableCell{content: "Execution Time", header: true},
		TableCell{content: "Iterations", header: true},
	)}

	for _, t := range Armeria.tickManager.Tickers {
		rows = append(rows, TableRow(
			TableCell{content: t.Name},
			TableCell{content: t.Interval.String()},
			TableCell{content: t.LastRanString()},
			TableCell{content: t.LastDurationString()},
			TableCell{content: t.IterationsString()},
		))
	}

	ctx.Player.client.ShowText(TextTable(rows...))
}

func handleSelectCommand(ctx *CommandContext) {
	mob := ctx.Args["mob"]
	optionId := ctx.Args["option_id"]

	mobGeneric, _, rt := ctx.Character.Room().Here().GetByAny(mob)
	if rt != RegistryTypeMobInstance {
		ctx.Player.client.ShowColorizedText("There are no mobs here with that name.", ColorError)
		return
	}
	mobInst := mobGeneric.(*MobInstance)

	if len(mobInst.ConvoText(optionId)) == 0 {
		ctx.Player.client.ShowColorizedText("That is not a valid selection.", ColorError)
		return
	}

	Armeria.commandManager.ProcessCommand(ctx.Player, fmt.Sprintf("say %s", mobInst.ConvoText(optionId)), false)

	go CallMobFunc(
		ctx.Character,
		mobInst,
		"conversation_select",
		lua.LString(optionId),
	)
}

func handleInteractCommand(ctx *CommandContext) {
	mob := ctx.Args["mob"]

	mobGeneric, _, rt := ctx.Character.Room().Here().GetByAny(mob)
	if rt != RegistryTypeMobInstance {
		ctx.Player.client.ShowColorizedText("There are no mobs here with that name.", ColorError)
		return
	}
	mobInst := mobGeneric.(*MobInstance)

	go CallMobFunc(
		ctx.Character,
		mobInst,
		"interact",
	)
}
