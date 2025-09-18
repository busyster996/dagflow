package cmd

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds cmd to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local cmd = require("cmd")
func Preload(L *lua.LState) {
	L.PreloadModule("cmd", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	// Register the encodings offered by base64 go module.
	t := L.NewTable()
	L.SetFuncs(t, map[string]lua.LGFunction{
		"exec":            Exec,
		"exec_by_timeout": ExecByTimeout,
	})
	L.Push(t)
	return 1
}
