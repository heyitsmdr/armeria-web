package armeria

import (
	"fmt"
	"io/ioutil"

	"go.uber.org/zap"

	lua "github.com/yuin/gopher-lua"
)

// ReadMobScript returns the script contents for a mob from disk.
func ReadMobScript(m *Mob) string {
	b, err := ioutil.ReadFile(m.ScriptFile())
	if err != nil {
		return ""
	}
	return string(b)
}

// WriteMobScript writes a mob script to disk.
func WriteMobScript(m *Mob, script string) {
	_ = ioutil.WriteFile(m.ScriptFile(), []byte(script), 0644)
}

// LuaMobSay is the handler for lua function: mob_say.
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

// LuaCharacterAttribute is the handler for lua function: c_attr.
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

// LuaSetCharacterAttribute is the handler for lua function: c_set_attr.
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

// LuaItemName is the handler for lua function: i_name.
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

// LuaInventoryGive is the handler for lua function: inv_give.
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

// CallMobFunc handles executing mob scripts within the Lua environment.
func CallMobFunc(invoker *Character, mi *MobInstance, funcName string, args ...lua.LValue) {
	L := lua.NewState()
	defer L.Close()

	// global variables
	L.SetGlobal("invoker_name", lua.LString(invoker.Name()))
	L.SetGlobal("mob_uuid", lua.LString(mi.UUID))
	L.SetGlobal("mob_name", lua.LString(mi.Name()))
	// global functions
	L.SetGlobal("say", L.NewFunction(LuaMobSay))
	L.SetGlobal("c_attr", L.NewFunction(LuaCharacterAttribute))
	L.SetGlobal("c_set_attr", L.NewFunction(LuaSetCharacterAttribute))
	L.SetGlobal("i_name", L.NewFunction(LuaItemName))
	L.SetGlobal("inv_give", L.NewFunction(LuaInventoryGive))

	err := L.DoFile(mi.Parent.ScriptFile())
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
					"There was an error running %s() on mob %s:\n%s",
					funcName,
					mi.Name(),
					err.Error(),
				),
				ColorError,
			)
		}
	}
}
