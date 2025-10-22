package libs

import (
	"math"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

type luai64 uintptr

func (id luai64) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func getI64Arg(L *lua.LState, idx int) luai64 {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(luai64); ok {
		return v
	}
	L.ArgError(idx, "i64 expected")
	return 0
}

func registerI64Type(L *lua.LState) {
	mt := L.NewTypeMetatable("i64")

	// tostring
	L.SetField(mt, "__tostring", L.NewFunction(func(L *lua.LState) int {
		v := getI64Arg(L, 1)
		L.Push(lua.LString(strconv.FormatInt(int64(v), 10)))
		return 1
	}))

	// eq
	L.SetField(mt, "__eq", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		L.Push(lua.LBool(v1 == v2))
		return 1
	}))

	// lt
	L.SetField(mt, "__lt", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		L.Push(lua.LBool(v1 < v2))
		return 1
	}))

	// le
	L.SetField(mt, "__le", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		L.Push(lua.LBool(v1 <= v2))
		return 1
	}))

	// add
	L.SetField(mt, "__add", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		L.Push(newI64(L, int64(v1+v2)))
		return 1
	}))

	// sub
	L.SetField(mt, "__sub", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		L.Push(newI64(L, int64(v1-v2)))
		return 1
	}))

	// mul
	L.SetField(mt, "__mul", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		L.Push(newI64(L, int64(v1*v2)))
		return 1
	}))

	// div
	L.SetField(mt, "__div", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		if v2 == 0 {
			L.ArgError(2, "division by zero")
			return 0
		}
		L.Push(newI64(L, int64(v1/v2)))
		return 1
	}))

	// mod
	L.SetField(mt, "__mod", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		L.Push(newI64(L, int64(v1%v2)))
		return 1
	}))

	// unm (-a)
	L.SetField(mt, "__unm", L.NewFunction(func(L *lua.LState) int {
		v := getI64Arg(L, 1)
		L.Push(newI64(L, int64(-v)))
		return 1
	}))

	// pow
	L.SetField(mt, "__pow", L.NewFunction(func(L *lua.LState) int {
		v1 := getI64Arg(L, 1)
		v2 := getI64Arg(L, 2)
		L.Push(newI64(L, int64(math.Pow(float64(v1), float64(v2)))))
		return 1
	}))
}

func newI64(L *lua.LState, v int64) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = luai64(v)
	L.SetMetatable(ud, L.GetTypeMetatable("i64"))
	return ud
}

// I64 Lua: I64("123456")
func I64(L *lua.LState) int {
	registerI64Type(L)
	s := L.CheckString(1)
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		L.ArgError(1, "invalid int64 string")
		return 0
	}
	L.Push(newI64(L, v))
	return 1
}
