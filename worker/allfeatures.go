//go:build dagflow.allfeatures

package worker

import (
	_ "github.com/busyster996/dagflow/worker/runner/docker"
	_ "github.com/busyster996/dagflow/worker/runner/kubectl"
	_ "github.com/busyster996/dagflow/worker/runner/lua"
	_ "github.com/busyster996/dagflow/worker/runner/scp"
	_ "github.com/busyster996/dagflow/worker/runner/sftp"
	_ "github.com/busyster996/dagflow/worker/runner/ssh"
	_ "github.com/busyster996/dagflow/worker/runner/wasm"
	_ "github.com/busyster996/dagflow/worker/runner/yaegi"
)
