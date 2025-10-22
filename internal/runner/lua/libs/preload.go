package libs

import (
	lua "github.com/yuin/gopher-lua"

	"github.com/busyster996/dagflow/internal/runner/lua/libs/base64"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/bit"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/cmd"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/crypto"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/filepath"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/hex"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/humanize"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/inspect"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/ioutil"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/json"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/regexp"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/runtime"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/shellescape"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/ssh"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/strings"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/tac"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/time"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/xmlpath"
	"github.com/busyster996/dagflow/internal/runner/lua/libs/yaml"
)

// Preload preload all gopher lua packages
func Preload(L *lua.LState) {
	L.SetGlobal("I64", L.NewFunction(I64))

	base64.Preload(L)
	bit.Preload(L)
	cmd.Preload(L)
	crypto.Preload(L)
	filepath.Preload(L)
	hex.Preload(L)
	humanize.Preload(L)
	inspect.Preload(L)
	ioutil.Preload(L)
	json.Preload(L)
	regexp.Preload(L)
	runtime.Preload(L)
	shellescape.Preload(L)
	ssh.Preload(L)
	strings.Preload(L)
	tac.Preload(L)
	time.Preload(L)
	xmlpath.Preload(L)
	yaml.Preload(L)
}
