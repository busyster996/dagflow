package lualibs

import (
	lua "github.com/yuin/gopher-lua"

	"github.com/busyster996/dagflow/pkg/lualibs/base64"
	"github.com/busyster996/dagflow/pkg/lualibs/bit"
	"github.com/busyster996/dagflow/pkg/lualibs/cmd"
	"github.com/busyster996/dagflow/pkg/lualibs/crypto"
	"github.com/busyster996/dagflow/pkg/lualibs/filepath"
	"github.com/busyster996/dagflow/pkg/lualibs/hex"
	"github.com/busyster996/dagflow/pkg/lualibs/humanize"
	"github.com/busyster996/dagflow/pkg/lualibs/inspect"
	"github.com/busyster996/dagflow/pkg/lualibs/ioutil"
	"github.com/busyster996/dagflow/pkg/lualibs/json"
	"github.com/busyster996/dagflow/pkg/lualibs/regexp"
	"github.com/busyster996/dagflow/pkg/lualibs/runtime"
	"github.com/busyster996/dagflow/pkg/lualibs/shellescape"
	"github.com/busyster996/dagflow/pkg/lualibs/ssh"
	"github.com/busyster996/dagflow/pkg/lualibs/strings"
	"github.com/busyster996/dagflow/pkg/lualibs/tac"
	"github.com/busyster996/dagflow/pkg/lualibs/time"
	"github.com/busyster996/dagflow/pkg/lualibs/xmlpath"
	"github.com/busyster996/dagflow/pkg/lualibs/yaml"
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
