package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/robertkrimen/otto"
	"github.com/tyler-sommer/squircy2/squircy/event"
	"github.com/veonik/go-lisp/lisp"
)

func newJavascriptVm(m *ScriptManager) *otto.Otto {
	jsVm := otto.New()
	jsVm.Set("Http", &m.httpHelper)
	jsVm.Set("Data", &m.dataHelper)
	jsVm.Set("Irc", &m.ircHelper)
	jsVm.Set("bind", func(call otto.FunctionCall) otto.Value {
		eventType := call.Argument(0).String()
		fnName := call.Argument(1).String()
		m.scriptHelper.Bind(Javascript, event.EventType(eventType), fnName)
		return otto.UndefinedValue()
	})
	jsVm.Set("unbind", func(call otto.FunctionCall) otto.Value {
		eventType := call.Argument(0).String()
		fnName := call.Argument(1).String()
		m.scriptHelper.Unbind(Javascript, event.EventType(eventType), fnName)
		return otto.UndefinedValue()
	})

	return jsVm
}

func newLuaVm(m *ScriptManager) *lua.State {
	luaVm := lua.NewState()
	luaVm.OpenLibs()
	luaVm.Register("typename", func(vm *lua.State) int {
		o := vm.Typename(int(vm.Type(1)))
		vm.PushString(o)
		return 1
	})
	luaVm.Register("setex", func(vm *lua.State) int {
		key := vm.ToString(1)
		value := vm.ToString(2)
		m.dataHelper.Set(key, value)
		return 0
	})
	luaVm.Register("getex", func(vm *lua.State) int {
		key := vm.ToString(1)
		value := m.dataHelper.Get(key)
		if value != nil {
			vm.PushString(value.(string))
			return 1
		}
		vm.PushNil()
		return 1
	})
	luaVm.Register("joinchan", func(vm *lua.State) int {
		channel := vm.ToString(1)
		m.ircHelper.Join(channel)
		return 0
	})
	luaVm.Register("partchan", func(vm *lua.State) int {
		channel := vm.ToString(1)
		m.ircHelper.Part(channel)
		return 0
	})
	luaVm.Register("privmsg", func(vm *lua.State) int {
		target := vm.ToString(1)
		message := vm.ToString(2)
		m.ircHelper.Privmsg(target, message)
		return 0
	})
	luaVm.Register("httpget", func(vm *lua.State) int {
		url := vm.ToString(1)
		res := m.httpHelper.Get(url)
		vm.PushString(res)
		return 1
	})
	luaVm.Register("bind", func(vm *lua.State) int {
		eventType := vm.ToString(1)
		fnName := vm.ToString(2)
		m.scriptHelper.Bind(Lua, event.EventType(eventType), fnName)
		return 0
	})
	luaVm.Register("unbind", func(vm *lua.State) int {
		eventType := vm.ToString(1)
		fnName := vm.ToString(2)
		m.scriptHelper.Unbind(Lua, event.EventType(eventType), fnName)
		return 0
	})

	return luaVm
}

func newLispVm(m *ScriptManager) {
	lisp.SetHandler("setex", func(vars ...lisp.Value) (lisp.Value, error) {
		if len(vars) != 2 {
			return lisp.Nil, nil
		}
		key := vars[0].String()
		value := vars[1].String()
		m.dataHelper.Set(key, value)
		return lisp.Nil, nil
	})
	lisp.SetHandler("getex", func(vars ...lisp.Value) (lisp.Value, error) {
		if len(vars) != 1 {
			return lisp.Nil, nil
		}
		key := vars[0].String()
		if val := m.dataHelper.Get(key); val != nil {
			return lisp.StringValue(val.(string)), nil
		}
		return lisp.Nil, nil
	})
	lisp.SetHandler("joinchan", func(vars ...lisp.Value) (lisp.Value, error) {
		if len(vars) != 1 {
			return lisp.Nil, nil
		}
		channel := vars[0].String()
		m.ircHelper.Join(channel)
		return lisp.Nil, nil
	})
	lisp.SetHandler("partchan", func(vars ...lisp.Value) (lisp.Value, error) {
		if len(vars) != 1 {
			return lisp.Nil, nil
		}
		channel := vars[0].String()
		m.ircHelper.Part(channel)
		return lisp.Nil, nil
	})
	lisp.SetHandler("privmsg", func(vars ...lisp.Value) (lisp.Value, error) {
		if len(vars) != 2 {
			return lisp.Nil, nil
		}
		target := vars[0].String()
		message := vars[1].String()
		m.ircHelper.Privmsg(target, message)
		return lisp.Nil, nil
	})
	lisp.SetHandler("httpget", func(vars ...lisp.Value) (lisp.Value, error) {
		if len(vars) != 1 {
			return lisp.Nil, nil
		}
		url := vars[0].String()
		res := m.httpHelper.Get(url)
		return lisp.StringValue(res), nil
	})
	lisp.SetHandler("bind", func(vars ...lisp.Value) (lisp.Value, error) {
		if len(vars) != 2 {
			return lisp.Nil, nil
		}
		eventType := vars[0].String()
		fnName := vars[1].String()
		m.scriptHelper.Bind(Lisp, event.EventType(eventType), fnName)
		return lisp.Nil, nil
	})
	lisp.SetHandler("unbind", func(vars ...lisp.Value) (lisp.Value, error) {
		if len(vars) != 2 {
			return lisp.Nil, nil
		}
		eventType := vars[0].String()
		fnName := vars[1].String()
		m.scriptHelper.Unbind(Lisp, event.EventType(eventType), fnName)
		return lisp.Nil, nil
	})
}