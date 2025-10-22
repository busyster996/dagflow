package cmd

import (
	"context"
	"time"

	"github.com/go-cmd/cmd"
	lua "github.com/yuin/gopher-lua"
)

const (
	//Timeout default execution timeout in seconds
	Timeout = 3 * time.Minute
)

// ExecByTimeout lua cmd.exec_by_timeout(timeout, command, args...) return ({status=0, stdout="", stderr=""}, err)
func ExecByTimeout(L *lua.LState) int {
	timeout := time.Duration(L.CheckInt64(1)) * time.Second
	command := L.CheckString(2)
	var args []string
	for i := 3; i <= L.GetTop(); i++ {
		args = append(args, L.CheckAny(i).String())
	}
	return _exec(L, timeout, command, args...)
}

// Exec lua cmd.exec(command, args...) return ({status=0, stdout="", stderr=""}, err)
func Exec(L *lua.LState) int {
	var command = L.CheckString(1)
	var args []string
	for i := 2; i <= L.GetTop(); i++ {
		args = append(args, L.CheckAny(i).String())
	}
	return _exec(L, Timeout, command, args...)
}

func _exec(L *lua.LState, timeout time.Duration, command string, args ...string) int {
	execCmd := cmd.NewCmdOptions(cmd.Options{
		Buffered:   true,
		Streaming:  false,
		BeforeExec: BeforeExec,
	}, command, args...)
	ctx, cancel := context.WithTimeout(L.Context(), timeout)
	defer cancel()

	select {
	case status := <-execCmd.Start():
		if status.Error != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(status.Error.Error()))
			return 2
		}

		stdout := L.NewTable()
		for i, line := range status.Stdout {
			L.SetTable(stdout, lua.LNumber(i+1), lua.LString(line))
		}

		stderr := L.NewTable()
		for i, line := range status.Stderr {
			L.SetTable(stderr, lua.LNumber(i+1), lua.LString(line))
		}

		result := L.NewTable()
		L.SetField(result, "status", lua.LNumber(status.Exit))
		if stdout.Len() > 0 {
			L.SetField(result, "stdout", stdout)
		} else {
			L.SetField(result, "stdout", lua.LNil)
		}
		if stderr.Len() > 0 {
			L.SetField(result, "stderr", stderr)
		} else {
			L.SetField(result, "stderr", lua.LNil)
		}

		L.Push(result)
		L.Push(lua.LNil)
		return 2
	case <-ctx.Done():
		_ = execCmd.Stop()
		L.Push(lua.LNil)
		L.Push(lua.LString(`execute timeout`))
		return 2
	}
}
