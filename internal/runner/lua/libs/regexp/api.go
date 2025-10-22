package regexp

import (
	"regexp"

	"github.com/dlclark/regexp2"
	lua "github.com/yuin/gopher-lua"
)

type luaRegexp struct {
	*regexp2.Regexp
}

func checkRegexp(L *lua.LState, n int) *luaRegexp {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaRegexp); ok {
		return v
	}
	L.ArgError(n, "regexp_ud expected")
	return nil
}

// Compile regexp.compile(string) returns (regexp_ud, error)
func Compile(L *lua.LState) int {
	expr := L.CheckString(1)
	reg, err := regexp2.Compile(expr, regexp2.RE2)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = &luaRegexp{Regexp: reg}
	L.SetMetatable(ud, L.GetTypeMetatable(`regexp_ud`))
	L.Push(ud)
	return 1
}

// Match regexp_ud:match(string) returns bool
func Match(L *lua.LState) int {
	reg := checkRegexp(L, 1)
	str := L.CheckString(2)
	ok, _ := reg.MatchString(str)
	L.Push(lua.LBool(ok))
	return 1
}

// SimpleMatch regexp.match(regular expression string, string) returns (bool, error)
func SimpleMatch(L *lua.LState) int {
	expr := L.CheckString(1)
	str := L.CheckString(2)
	reg, err := regexp.Compile(expr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LBool(reg.MatchString(str)))
	return 1
}

// FindAllStringSubmatch regexp_ud:find_all_string_submatch(string) returns table of table of strings
func FindAllStringSubmatch(L *lua.LState) int {
	reg := checkRegexp(L, 1)
	str := L.CheckString(2)
	result := L.NewTable()
	m, _ := reg.FindStringMatch(str)
	for m != nil {
		row := L.NewTable()
		for _, v := range m.Groups() {
			row.Append(lua.LString(v.String()))
		}
		result.Append(row)
		m, _ = reg.FindNextMatch(m)
	}
	L.Push(result)
	return 1
}

// SimpleFindAllStringSubmatch regexp.find_all_string_submatch(regular expression string, string) returns (table of table strings, error)
func SimpleFindAllStringSubmatch(L *lua.LState) int {
	expr := L.CheckString(1)
	str := L.CheckString(2)
	reg, err := regexp.Compile(expr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	for _, t := range reg.FindAllStringSubmatch(str, -1) {
		row := L.NewTable()
		for _, v := range t {
			row.Append(lua.LString(v))
		}
		result.Append(row)
	}
	L.Push(result)
	return 1
}
