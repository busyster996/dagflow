package ssh

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds ssh to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local ssh = require("ssh")
func Preload(L *lua.LState) {
	L.PreloadModule(`ssh`, Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	ssh := L.NewTypeMetatable(`ssh`)
	L.SetGlobal(`ssh`, ssh)
	L.SetField(ssh, `__index`, L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		`execute`: Execute,
		`copy`:    Copy,
	}))
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1

}

var api = map[string]lua.LGFunction{
	`client`: Client,
}
