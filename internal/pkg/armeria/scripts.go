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

func LuaMobSay(L *lua.LState) int {
	text := L.ToString(1)
	mname := lua.LVAsString(L.GetGlobal("mob_name"))
	mid := lua.LVAsString(L.GetGlobal("mob_uuid"))

	m := Armeria.mobManager.MobByName(mname)
	mi := m.InstanceByUUID(mid)

	for _, c := range mi.Location.Room().Characters(nil) {
		c.Player().clientActions.ShowColorizedText(
			fmt.Sprintf("%s says, \"%s\".", mi.FormattedName(), text),
			ColorSay,
		)
	}

	return 0
}

func CallMobFunc(invoker *Character, mi *MobInstance, funcName string, args ...lua.LValue) {
	L := lua.NewState()
	defer L.Close()

	// global variables
	L.SetGlobal("invoker_name", lua.LString(invoker.Name()))
	L.SetGlobal("mob_uuid", lua.LString(mi.UUID))
	L.SetGlobal("mob_name", lua.LString(mi.Name()))
	// global functions
	L.SetGlobal("say", L.NewFunction(LuaMobSay))

	err := L.DoFile(mi.Parent().ScriptFile())
	if err != nil {
		Armeria.log.Error("error compiling lua script",
			zap.String("script", mi.Parent().ScriptFile()),
			zap.Error(err),
		)
		if invoker.HasPermission("CAN_BUILD") {
			invoker.Player().clientActions.ShowColorizedText(
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
			zap.String("script", mi.Parent().ScriptFile()),
			zap.String("function", funcName),
			zap.Error(err),
		)
		if invoker.HasPermission("CAN_BUILD") {
			invoker.Player().clientActions.ShowColorizedText(
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
