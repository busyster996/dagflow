package lua

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/runner/lua/libs"
	"github.com/busyster996/dagflow/internal/storage"
)

type sLua struct {
	storage   storage.IStep
	workspace string
}

func (l *sLua) Run(ctx context.Context) (exit common.ExecCode, err error) {
	defer func() {
		if _r := recover(); _r != nil {
			err = fmt.Errorf("panic during execution %v", _r)
			exit = common.ExecCodeSystemErr
			stack := debug.Stack()
			if _err, ok := _r.(error); ok && strings.Contains(_err.Error(), context.Canceled.Error()) {
				exit = common.ExecCodeKilled
				err = common.ExecErrKilled
			}
			l.storage.Log().Write(err.Error(), string(stack))
		}
	}()
	params, err := l.getParams()
	if err != nil {
		return common.ExecCodeFailed, err
	}
	content, err := l.storage.Content()
	if err != nil {
		return common.ExecCodeFailed, err
	}
	vm := lua.NewState(lua.Options{
		IncludeGoStackTrace: true,
	})
	defer vm.Close()
	l.loadLibs(vm)
	vm.SetGlobal("params", luar.New(vm, params))
	vm.SetGlobal("storage", luar.New(vm, l.storage))
	if err = vm.DoString(content); err != nil {
		return common.ExecCodeFailed, err
	}
	return common.ExecCodeSuccess, nil
}

func (l *sLua) loadLibs(vm *lua.LState) {
	_printf := func(format string, args ...any) {
		l.storage.Log().Write(fmt.Sprintf(format, args...))
	}
	_print := func(args ...any) {
		l.storage.Log().Write(fmt.Sprint(args...))
	}
	_println := func(args ...any) {
		l.storage.Log().Write(fmt.Sprintln(args...))
	}
	vm.SetGlobal("printf", luar.New(vm, _printf))
	vm.SetGlobal("print", luar.New(vm, _print))
	vm.SetGlobal("println", luar.New(vm, _println))

	vm.PreloadModule("log", l.log)
	libs.Preload(vm)
}

func (l *sLua) getParams() (gjson.Result, error) {
	var rawJSON string
	var err error
	taskEnv := l.storage.GlobalEnv().List()
	for _, v := range taskEnv {
		rawJSON, err = sjson.Set(rawJSON, v.Name, []byte(v.Value))
		if err != nil {
			return gjson.Result{}, err
		}
	}
	stepEnv := l.storage.Env().List()
	for _, v := range stepEnv {
		rawJSON, err = sjson.Set(rawJSON, v.Name, []byte(v.Value))
		if err != nil {
			return gjson.Result{}, err
		}
	}
	return gjson.Parse(rawJSON), nil
}

func (l *sLua) Clear() error {
	return nil
}

func (l *sLua) log(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, map[string]lua.LGFunction{
		"debug":  l.logFunc("[debug]"),
		"info":   l.logFunc("[info]"),
		"warn":   l.logFunc("[warn]"),
		"error":  l.logFunc("[error]"),
		"debugf": l.logFFunc("[debug]"),
		"infof":  l.logFFunc("[info]"),
		"warnf":  l.logFFunc("[warn]"),
		"errorf": l.logFFunc("[error]"),
	})
	L.Push(t)
	return 1
}

func (l *sLua) logFunc(level string) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		var v = []string{
			level,
			L.Where(1),
		}
		for i := 1; i <= L.GetTop(); i++ {
			v = append(v, L.Get(i).String())
		}
		l.storage.Log().Write(v...)
		return 0
	}
}

func (l *sLua) logFFunc(level string) func(L *lua.LState) int {
	return func(L *lua.LState) int {
		format := "%s %s" + L.CheckString(1)
		var v = []interface{}{
			level,
			L.Where(1),
		}
		for i := 2; i <= L.GetTop(); i++ {
			v = append(v, L.Get(i))
		}
		l.storage.Log().Writef(format, v...)
		return 0
	}
}
