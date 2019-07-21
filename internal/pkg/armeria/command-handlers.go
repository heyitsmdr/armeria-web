package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func handleLoginCommand(ctx *CommandContext) {
	if len(ctx.Args) != 2 {
		return
	}

	c := Armeria.characterManager.CharacterByName(ctx.Args["character"])
	if c == nil {
		ctx.Player.client.ShowText("Character not found.")
		return
	}

	if !c.CheckPassword(ctx.Args["password"]) {
		ctx.Player.client.ShowColorizedText("Password incorrect for that character.", ColorError)
		return
	}

	if c.Player() != nil {
		ctx.Player.client.ShowColorizedText("This character is already logged in.", ColorError)
		return
	}

	if c.Location().Room() == nil {
		ctx.Player.client.ShowColorizedText("This character logged out of a room which no longer exists.", ColorError)
		return
	}

	ctx.Player.AttachCharacter(c)
	c.SetPlayer(ctx.Player)

	ctx.Player.client.ShowColorizedText(fmt.Sprintf("You've entered Armeria as %s!", c.FormattedName()), ColorSuccess)

	c.LoggedIn()
}

func handleLookCommand(ctx *CommandContext) {
	rm := ctx.Character.Location().Room()

	var objNames []string
	for _, o := range rm.Objects() {
		if o.Type() != ObjectTypeCharacter || o.Name() != ctx.Character.Name() {
			objNames = append(objNames, o.FormattedName())
		}
	}

	var withYou string
	if len(objNames) > 0 {
		withYou = fmt.Sprintf("\nHere with you: %s.", strings.Join(objNames, ", "))
	}

	ar := ctx.Character.Location().Area().AdjacentRooms(rm)
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

	ctx.Player.client.ShowText(
		ctx.Character.Colorize(rm.Attribute("title"), ColorRoomTitle) + "\n" +
			rm.Attribute("description") +
			ctx.Character.Colorize(validDirString, ColorRoomDirs) +
			withYou,
	)

	if ctx.PlayerInitiated {
		for _, c := range ctx.Character.Location().Room().Characters(ctx.Character) {
			c.Player().client.ShowText(
				fmt.Sprintf("%s takes a look around.", ctx.Character.FormattedName()),
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

	verbs := []string{"say", "says"}
	lastChar := ctx.Args["text"][len(ctx.Args["text"])-1:]
	if lastChar == "?" {
		verbs = []string{"ask", "asks"}
	} else if lastChar == "!" {
		verbs = []string{"exclaim", "exclaims"}
	}

	ctx.Player.client.ShowText(
		ctx.Player.Character().Colorize(fmt.Sprintf("You %s, \"%s\".", verbs[0], ctx.Args["text"]), ColorSay),
	)

	room := ctx.Character.Location().Room()
	for _, c := range room.Characters(ctx.Character) {
		c.Player().client.ShowText(
			c.Player().Character().Colorize(
				fmt.Sprintf("%s %s, \"%s\".", ctx.Character.FormattedName(), verbs[1], ctx.Args["text"]),
				ColorSay,
			),
		)
	}

	for _, o := range room.Objects() {
		if o.Type() == ObjectTypeMob {
			go CallMobFunc(
				ctx.Character,
				o.(*MobInstance),
				"character_said",
				lua.LString(ctx.Character.Name()),
				lua.LString(ctx.Args["text"]),
			)
		}
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
		ctx.Player.client.ShowText("That's not a valid direction to move in.")
		return
	}

	loc := ctx.Character.Location()
	o := misc.DirectionOffsets(d)
	x := loc.Coords.X() + o["x"]
	y := loc.Coords.Y() + o["y"]
	z := loc.Coords.Z() + o["z"]

	newLocation := NewLocation(loc.AreaUUID(), x, y, z)
	moveAllowed, moveError := ctx.Character.MoveAllowed(newLocation)
	if !moveAllowed {
		ctx.Player.client.ShowColorizedText(moveError, ColorError)
		return
	}

	ctx.Character.Move(
		newLocation,
		ctx.Character.Colorize(fmt.Sprintf("You walk to %s.", walkDir), ColorMovement),
		ctx.Character.Colorize(fmt.Sprintf("%s walks to %s.", ctx.Character.FormattedName(), walkDir), ColorMovement),
		ctx.Character.Colorize(fmt.Sprintf("%s walked in from %s.", ctx.Character.FormattedName(), arriveDir), ColorMovement),
	)

	Armeria.commandManager.ProcessCommand(ctx.Player, "look", false)
}

func handleRoomEditCommand(stx *CommandContext) {
	stx.Player.client.ShowObjectEditor(stx.Character.Location().Room().EditorData())
}

func handleRoomSetCommand(ctx *CommandContext) {
	attr := strings.ToLower(ctx.Args["property"])

	if !misc.Contains(ValidRoomAttributes(), attr) {
		ctx.Player.client.ShowColorizedText("That's not a valid room attribute.", ColorError)
		return
	}

	ctx.Character.Location().Room().SetAttribute(attr, ctx.Args["value"])

	for _, c := range ctx.Character.Location().Room().Characters(ctx.Character) {
		c.Player().client.ShowText(
			fmt.Sprintf("%s modified the room.", ctx.Character.FormattedName()),
		)
	}

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the [b]%s[/b] property of the room.", attr),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(ctx.Character.Location().Room().EditorData())
	}
}

func handleRoomCreateCommand(ctx *CommandContext) {
	d := ctx.Args["direction"]

	o := misc.DirectionOffsets(d)
	if o == nil {
		ctx.Player.client.ShowColorizedText("That's not a valid direction to create a room in.", ColorError)
		return
	}

	loc := ctx.Character.Location()
	x := loc.Coords.X() + o["x"]
	y := loc.Coords.Y() + o["y"]
	z := loc.Coords.Z() + o["z"]

	newLoc := NewLocation(loc.AreaUUID(), x, y, z)
	if newLoc.Room() != nil {
		ctx.Player.client.ShowColorizedText("There's already a room in that direction.", ColorError)
		return
	}

	room := Armeria.worldManager.CreateRoom(NewCoords(x, y, z, loc.Coords.I()))
	ctx.Character.Location().Area().AddRoom(room)

	for _, c := range ctx.Character.Location().Area().Characters(nil) {
		c.Player().client.RenderMap()
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

	loc := ctx.Character.Location()
	x := loc.Coords.X() + o["x"]
	y := loc.Coords.Y() + o["y"]
	z := loc.Coords.Z() + o["z"]

	l := NewLocation(loc.AreaUUID(), x, y, z)

	r := l.Room()
	if r == nil {
		ctx.Player.client.ShowColorizedText("There's no room in that direction.", ColorError)
		return
	}

	if len(r.Objects()) > 0 {
		ctx.Player.client.ShowColorizedText("There is something in the room you're attempting to destroy.", ColorError)
		return
	}

	loc.Area().RemoveRoom(r)

	for _, c := range ctx.Character.Location().Area().Characters(nil) {
		c.Player().client.RenderMap()
	}

	ctx.Player.client.ShowColorizedText("The room has been destroyed.", ColorSuccess)
}

func handleSaveCommand(ctx *CommandContext) {
	Armeria.Save()
	ctx.Player.client.ShowText("The game data has been saved to disk.")
}

func handleRefreshCommand(ctx *CommandContext) {
	ctx.Player.client.RenderMap()
	ctx.Player.client.SyncRoomObjects()
	ctx.Player.client.SyncRoomTitle()
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

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You whisper to %s, \"%s\".", c.FormattedNameWithTitle(), m),
		ColorWhisper,
	)

	c.Player().client.ShowColorizedText(
		fmt.Sprintf("%s whispers to you from %s, \"%s\".",
			ctx.Character.FormattedNameWithTitle(),
			c.Location().Area().Name(),
			m,
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
		fmt.Sprintf("There are [b]%d[/b] characters%s: %s.", len(chars), matchingText, strings.Join(chars, ", ")),
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

	c.SetAttribute(attr, val)

	ctx.Player.client.ShowColorizedText(
		fmt.Sprintf("You modified the [b]%s[/b] property of the character %s.", attr, c.FormattedName()),
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
		fmt.Sprintf("There are [b]%d[/b] mobs%s: %s.", len(mobs), matchingText, strings.Join(mobs, ", ")),
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
		fmt.Sprintf("A mob named [b]%s[/b] has been created.", n),
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

func handleMobSetCommand(ctx *CommandContext) {
	mob := strings.ToLower(ctx.Args["mob"])
	attr := strings.ToLower(ctx.Args["property"])
	val := strings.ToLower(ctx.Args["value"])

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
		fmt.Sprintf("You modified the [b]%s[/b] property of the mob [b]%s[/b].", attr, m.UnsafeName),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(m.EditorData())
	}
}

func handleMobSpawnCommand(ctx *CommandContext) {
	m := Armeria.mobManager.MobByName(ctx.Args["mob"])
	if m == nil {
		ctx.Player.client.ShowColorizedText("That mob doesn't exist.", ColorError)
		return
	}

	l := ctx.Character.Location()

	mi := m.CreateInstance(l)
	ctx.Character.Location().Room().AddObjectToRoom(mi)

	for _, c := range ctx.Character.Location().Room().Characters(nil) {
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

	var mobLocations []string
	for i, mi := range m.Instances() {
		a := mi.Location.Area()
		mobLocations = append(
			mobLocations,
			fmt.Sprintf(
				"  %d) %s (%s) is currently at %s,%d,%d,%d (%s).",
				i+1,
				mi.FormattedName(),
				mi.Id(),
				a.Name(),
				mi.Location.Coords.X(),
				mi.Location.Coords.Y(),
				mi.Location.Coords.Z(),
				mi.Location.Room().Attribute("title"),
			),
		)
	}

	ctx.Player.client.ShowText(
		fmt.Sprintf(
			"Instances of %s:\n%s",
			m.Name(),
			strings.Join(mobLocations, "\n"),
		),
	)
}

func handleWipeCommand(ctx *CommandContext) {
	for _, o := range ctx.Character.Location().Room().Objects() {
		switch o.Type() {
		case ObjectTypeMob:
			m := Armeria.mobManager.MobByName(o.Name())
			s := ctx.Character.Location().Room().RemoveObjectFromRoom(o)
			if m != nil && s {
				m.DeleteInstance(o.(*MobInstance))
			}
		case ObjectTypeItem:
			i := Armeria.itemManager.ItemByName(o.Name())
			s := ctx.Character.Location().Room().RemoveObjectFromRoom(o)
			if i != nil && s {
				i.DeleteInstance(o.(*ItemInstance))
			}
		}
	}

	for _, c := range ctx.Character.Location().Room().Characters(ctx.Character) {
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
		fmt.Sprintf("An item named [b]%s[/b] has been created.", n),
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
		fmt.Sprintf("There are [b]%d[/b] items%s: %s.", len(items), matchingText, strings.Join(items, ", ")),
	)
}

func handleItemSpawnCommand(ctx *CommandContext) {
	i := Armeria.itemManager.ItemByName(ctx.Args["item"])
	if i == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist.", ColorError)
		return
	}

	l := ctx.Character.Location()

	ii := i.CreateInstance()
	ii.SetLocationType(ItemLocationRoom)
	ii.Location.SetAreaUUID(l.AreaUUID())
	ii.Location.Coords.Set(l.Coords.X(), l.Coords.Y(), l.Coords.Z(), l.Coords.I())
	ctx.Character.Location().Room().AddObjectToRoom(ii)

	for _, c := range ctx.Character.Location().Room().Characters(nil) {
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

func handleItemSetCommand(ctx *CommandContext) {
	item := strings.ToLower(ctx.Args["item"])
	attr := strings.ToLower(ctx.Args["property"])
	val := strings.ToLower(ctx.Args["value"])

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
		fmt.Sprintf("You modified the [b]%s[/b] property of the item [b]%s[/b].", attr, i.Name()),
		ColorSuccess,
	)

	editorOpen := ctx.Character.TempAttribute(TempAttributeEditorOpen)
	if editorOpen == "true" {
		ctx.Player.client.ShowObjectEditor(i.EditorData())
	}
}

func handleItemInstancesCommand(ctx *CommandContext) {
	i := Armeria.itemManager.ItemByName(ctx.Args["item"])
	if i == nil {
		ctx.Player.client.ShowColorizedText("That item doesn't exist.", ColorError)
		return
	}

	var itemLocations []string
	for idx, ii := range i.Instances() {
		if ii.LocationType() == ItemLocationRoom {
			a := ii.Location.Area()
			itemLocations = append(
				itemLocations,
				fmt.Sprintf(
					"  %d) %s (%s) is currently at %s,%d,%d,%d (%s).",
					idx+1,
					ii.FormattedName(),
					ii.Id(),
					a.Name(),
					ii.Location.Coords.X(),
					ii.Location.Coords.Y(),
					ii.Location.Coords.Z(),
					ii.Location.Room().Attribute("title"),
				),
			)
		} else if ii.LocationType() == ItemLocationCharacter {
			itemLocations = append(
				itemLocations,
				fmt.Sprintf(
					"  %d) %s (%s) is currently on the character %s.",
					idx+1,
					ii.FormattedName(),
					ii.Id(),
					ii.Character().FormattedName(),
				),
			)
		}
	}

	ctx.Player.client.ShowText(
		fmt.Sprintf(
			"Instances of %s:\n%s",
			i.Name(),
			strings.Join(itemLocations, "\n"),
		),
	)
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
		fmt.Sprintf("An area named [b]%s[/b] has been created!", a.Name()),
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
		fmt.Sprintf("There are [b]%d[/b] areas%s: %s.", len(areas), matchingText, strings.Join(areas, ", ")),
	)
}

func handleAreaEditCommand(ctx *CommandContext) {
	area := ctx.Args["area"]
	var a *Area
	if len(area) == 0 {
		a = ctx.Character.Location().Area()
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

	var l *Location
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

		cl := c.Location().Coords
		l = NewLocation(c.Location().Area().Id(), cl.X(), cl.Y(), cl.Z())
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

		l = NewLocation(a.Id(), x, y, z)
		moveMsg = fmt.Sprintf("You teleported to [b]%s[/b] at %d, %d, %d.", a.Name(), x, y, z)
	}

	if l.Room() == nil {
		ctx.Player.client.ShowColorizedText("You cannot teleport there!", ColorError)
		return
	}

	ctx.Character.Move(
		l,
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
		if cmd.CheckPermissions(ctx.Player) && len(cmd.Alias) == 0 {
			valid = append(valid, cmd)

			if len(cmd.Name) > largest {
				largest = len(cmd.Name)
			}
		}
	}

	var list []string
	for _, cmd := range valid {
		list = append(list,
			fmt.Sprintf("  /%-"+strconv.Itoa(largest)+"v %s", cmd.Name, cmd.Help),
		)
	}

	ctx.Player.client.ShowColorizedText(
		TextStyle(
			fmt.Sprintf("[b]Commands you can use:[/b]\n%s", strings.Join(list, "\n")),
			TextStyleMonospace,
		),
		ColorCmdHelp,
	)
}
