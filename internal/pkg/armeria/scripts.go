package armeria

import (
	"armeria/internal/pkg/misc"
	"fmt"
	"io/ioutil"
	"time"

	"go.uber.org/zap"

	lua "github.com/yuin/gopher-lua"
)

// ReadMobScript returns the script contents for a mob from disk.
func ReadMobScript(m *Mob) string {
	return m.Script()
}

// WriteMobScript writes a mob script to disk.
func WriteMobScript(m *Mob, script string) {
	_ = ioutil.WriteFile(m.ScriptFile(), []byte(script), 0644)
	m.CacheScript()
}

// LuaInvoker returns the Character that invoked the lua function.
func LuaInvoker(L *lua.LState) *Character {
	cuuid := lua.LVAsString(L.GetGlobal("invoker_uuid"))
	if o, rt := Armeria.registry.Get(cuuid); rt == RegistryTypeCharacter {
		return o.(*Character)
	}

	return nil
}

// LuaMobInstance returns the MobInstance that the script belongs to.
func LuaMobInstance(L *lua.LState) *MobInstance {
	muuid := lua.LVAsString(L.GetGlobal("mob_uuid"))
	if o, rt := Armeria.registry.Get(muuid); rt == RegistryTypeMobInstance {
		return o.(*MobInstance)
	}

	return nil
}

// LuaMobSay (mob_say) causes the mob to say something to the room.
func LuaMobSay(L *lua.LState) int {
	text := L.ToString(1)
	mname := lua.LVAsString(L.GetGlobal("mob_name"))
	mid := lua.LVAsString(L.GetGlobal("mob_uuid"))

	m := Armeria.mobManager.MobByName(mname)
	mi := m.InstanceByUUID(mid)

	normalizedText, textType := TextPunctuation(text)

	var verb string
	switch textType {
	case TextQuestion:
		verb = "asks"
	case TextExclaim:
		verb = "exclaims"
	default:
		verb = "says"
	}

	for _, c := range mi.Room().Here().Characters(true) {
		c.Player().client.ShowColorizedText(
			fmt.Sprintf("%s %s, \"%s\"", mi.FormattedName(), verb, normalizedText),
			ColorSay,
		)
	}

	return 0
}

// LuaSleep (sleep) causes the script to sleep for a particular time.Duration (ie: 30s).
func LuaSleep(L *lua.LState) int {
	duration := L.ToString(1)

	if d, err := time.ParseDuration(duration); err == nil {
		time.Sleep(d)
	}

	return 0
}

// LuaCharacterAttribute (c_attr) retrieves a permanent or temporary unsafeCharacter attribute.
func LuaCharacterAttribute(L *lua.LState) int {
	character := L.ToString(1)
	attr := L.ToString(2)
	tmp := L.ToBool(3)

	c := Armeria.characterManager.CharacterByName(character)
	if c == nil {
		L.Push(lua.LNumber(-1))
		return 1
	}

	var attrValue string
	if !tmp {
		attrValue = c.Attribute(attr)
	} else {
		attrValue = c.TempAttribute(attr)
	}

	L.Push(lua.LString(attrValue))
	return 1
}

// LuaSetCharacterAttribute (c_set_attr) sets a permanent or temporary unsafeCharacter attribute.
func LuaSetCharacterAttribute(L *lua.LState) int {
	character := L.ToString(1)
	attr := L.ToString(2)
	val := L.ToString(3)
	tmp := L.ToBool(4)

	c := Armeria.characterManager.CharacterByName(character)
	if c == nil {
		L.Push(lua.LNumber(-1))
		return 1
	}

	if !tmp {
		err := c.SetAttribute(attr, val)
		if err != nil {
			L.Push(lua.LNumber(-2))
			return 1
		}
	} else {
		c.SetTempAttribute(attr, val)
	}

	L.Push(lua.LNumber(0))
	return 1
}

// LuaItemName (i_name) returns the formatted item name from an item uuid.
func LuaItemName(L *lua.LState) int {
	uuid := L.ToString(1)

	o, rt := Armeria.registry.Get(uuid)
	if rt != RegistryTypeItemInstance {
		L.Push(lua.LNumber(-1))
		return 1
	}

	L.Push(lua.LString(o.(*ItemInstance).FormattedName()))
	return 1
}

// LuaInventoryGive (inv_give) gives an item to a unsafeCharacter from the mob's inventory.
func LuaInventoryGive(L *lua.LState) int {
	cuuid := L.ToString(1)
	iuuid := L.ToString(2)

	muuid := lua.LVAsString(L.GetGlobal("mob_uuid"))

	var mi *MobInstance
	if o, rt := Armeria.registry.Get(muuid); rt == RegistryTypeMobInstance {
		mi = o.(*MobInstance)
	} else {
		return 0
	}

	var c *Character
	if o, rt := Armeria.registry.Get(cuuid); rt == RegistryTypeCharacter {
		c = o.(*Character)
	} else {
		return 0
	}

	var ii *ItemInstance
	if o, rt := Armeria.registry.Get(iuuid); rt == RegistryTypeItemInstance {
		ii = o.(*ItemInstance)
	} else {
		return 0
	}

	if !mi.Inventory().Contains(iuuid) {
		return 0
	}

	if c.Inventory().Count() >= c.Inventory().MaxSize() {
		return 0
	}

	mi.Inventory().Remove(iuuid)

	_ = c.Inventory().Add(iuuid)

	if c.Online() {
		c.Player().client.SyncInventory()
		c.Player().client.ShowText(
			fmt.Sprintf(
				"%s gave you a %s.",
				mi.FormattedName(),
				ii.FormattedName(),
			),
		)

		for _, char := range c.Room().Here().Characters(true, c) {
			char.Player().client.ShowText(
				fmt.Sprintf(
					"%s gave something to %s.",
					mi.FormattedName(),
					c.FormattedName(),
				),
			)
		}
	}

	return 0
}

// LuaStartConvo (start_convo) starts a new conversation and begins conversation ticks.
func LuaStartConvo(L *lua.LState) int {
	cuuid := lua.LVAsString(L.GetGlobal("invoker_uuid"))
	muuid := lua.LVAsString(L.GetGlobal("mob_uuid"))

	var c *Character
	var mi *MobInstance

	if o, rt := Armeria.registry.Get(cuuid); rt == RegistryTypeCharacter {
		c = o.(*Character)
	} else {
		return 0
	}

	if o, rt := Armeria.registry.Get(muuid); rt == RegistryTypeMobInstance {
		mi = o.(*MobInstance)
	} else {
		return 0
	}

	if !c.Online() {
		return 0
	}

	// check if the unsafeCharacter is already in a conversation; cancel it if so
	if c.MobConvo() != nil {
		c.MobConvo().Cancel()
	}

	convo := Armeria.convoManager.NewConversation()
	convo.SetCharacter(c)
	convo.SetMobInstance(mi)
	c.SetMobConvo(convo)
	convo.Start()

	return 0
}

// LuaEndConvo (end_convo) ends a conversation and stops conversation ticks.
func LuaEndConvo(L *lua.LState) int {
	cuuid := lua.LVAsString(L.GetGlobal("invoker_uuid"))

	var c *Character
	if o, rt := Armeria.registry.Get(cuuid); rt == RegistryTypeCharacter {
		c = o.(*Character)
	} else {
		return 0
	}

	if c.MobConvo() != nil {
		c.MobConvo().Cancel()
	}

	return 0
}

// LuaShop (shop) displays shop items to the character that invoked the command.
func LuaShop(L *lua.LState) int {
	ledgerName := L.ToString(1)

	ledger := Armeria.ledgerManager.LedgerByName(ledgerName)
	if ledger == nil {
		panic(fmt.Sprintf("ledger not found: %s", ledgerName))
	}

	c := LuaInvoker(L)
	if c == nil {
		return 0
	}

	mi := LuaMobInstance(L)
	if mi == nil {
		return 0
	}

	// Ensure mob's inventory has at least one item instance from everything on the ledger.
	mi.Inventory().PopulateFromLedger(ledger)

	// Ensure mob knows about this ledger
	mi.AddItemLedger(ledger)

	// Generate buy table
	buyTable := []string{TableRow(
		TableCell{content: "Item", header: true},
		TableCell{content: "Description", header: true},
		TableCell{content: "Buy", header: true},
	)}
	for _, ii := range mi.Inventory().Items() {
		ledgerEntry := ledger.Contains(ii.Name())
		if ledgerEntry != nil && ledgerEntry.BuyPrice > 0 {
			buyTable = append(buyTable, TableRow(
				TableCell{content: ii.FormattedName()},
				TableCell{content: ii.Attribute(AttributeDescription)},
				TableCell{content: TextStyle(
					fmt.Sprintf("Buy %s <%s>", ii.Name(), misc.Money.FormatMoney(ledgerEntry.BuyPrice)),
					WithLinkCmd(fmt.Sprintf("/buy \"%s\" \"%s\"", mi.Name(), ii.Name())),
				)},
			))
		}
	}
	c.Player().client.ShowText(TextTable(buyTable...))

	// Generate sell table
	sellTable := []string{TableRow(
		TableCell{content: "Inventory", header: true},
		TableCell{content: "Sell", header: true},
	)}
	for _, ii := range c.Inventory().Items() {
		ledgerEntry := ledger.Contains(ii.Name())
		if ledgerEntry != nil && ledgerEntry.SellPrice > 0 {
			sellTable = append(sellTable, TableRow(
				TableCell{content: ii.FormattedName()},
				TableCell{content: TextStyle(
					fmt.Sprintf("Sell %s <%s>", ii.Name(), misc.Money.FormatMoney(ledgerEntry.SellPrice)),
					WithLinkCmd(fmt.Sprintf("/sell \"%s\" \"%s\"", mi.Name(), ii.ID())),
				)},
			))
		}
	}
	if len(sellTable) > 1 {
		c.Player().client.ShowText(TextTable(sellTable...))
	}

	return 0
}

// LuaRoomText (room_text) sends arbitrary text to the room.
func LuaRoomText(L *lua.LState) int {
	text := L.ToString(1)
	muuid := lua.LVAsString(L.GetGlobal("mob_uuid"))

	var mi *MobInstance
	if o, rt := Armeria.registry.Get(muuid); rt == RegistryTypeMobInstance {
		mi = o.(*MobInstance)
	} else {
		return 0
	}

	for _, c := range mi.Room().Here().Characters(true) {
		c.Player().client.ShowText(text)
	}

	return 0
}

// CallMobFunc handles executing mob scripts within the Lua environment.
func CallMobFunc(invoker *Character, mi *MobInstance, funcName string, args ...lua.LValue) {
	L := lua.NewState()
	defer L.Close()

	// global variables
	L.SetGlobal("invoker_uuid", lua.LString(invoker.ID()))
	L.SetGlobal("invoker_name", lua.LString(invoker.Name()))
	L.SetGlobal("mob_uuid", lua.LString(mi.UUID))
	L.SetGlobal("mob_name", lua.LString(mi.Name()))
	// global functions
	L.SetGlobal("say", L.NewFunction(LuaMobSay))
	L.SetGlobal("sleep", L.NewFunction(LuaSleep))
	L.SetGlobal("start_convo", L.NewFunction(LuaStartConvo))
	L.SetGlobal("end_convo", L.NewFunction(LuaEndConvo))
	L.SetGlobal("c_attr", L.NewFunction(LuaCharacterAttribute))
	L.SetGlobal("c_set_attr", L.NewFunction(LuaSetCharacterAttribute))
	L.SetGlobal("i_name", L.NewFunction(LuaItemName))
	L.SetGlobal("inv_give", L.NewFunction(LuaInventoryGive))
	L.SetGlobal("room_text", L.NewFunction(LuaRoomText))
	L.SetGlobal("shop", L.NewFunction(LuaShop))

	err := L.DoString(mi.Parent.Script())
	if err != nil {
		Armeria.log.Error("error compiling lua script",
			zap.String("script", mi.Parent.ScriptFile()),
			zap.Error(err),
		)
		if invoker.HasPermission("CAN_BUILD") {
			invoker.Player().client.ShowColorizedText(
				fmt.Sprintf(
					"There was an error compiling %s() on mob %s:\n%s",
					funcName,
					mi.Name(),
					err.Error(),
				),
				ColorError,
			)
		}
		return
	}

	lv := L.GetGlobal(funcName)
	if lv.Type() == lua.LTNil {
		return
	}

	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal(funcName),
		NRet:    0,
		Protect: true,
	}, args...)
	if err != nil {
		Armeria.log.Error("error executing function in lua script",
			zap.String("script", mi.Parent.ScriptFile()),
			zap.String("function", funcName),
			zap.Error(err),
		)
		if invoker.HasPermission("CAN_BUILD") {
			invoker.Player().client.ShowColorizedText(
				fmt.Sprintf(
					"There was an error running %s on mob %s.\n\n%s",
					TextStyle(funcName+"()", WithBold()),
					TextStyle(mi.Name(), WithBold()),
					err.Error(),
				),
				ColorError,
			)
		}
	}
}
