//go:build dagflow.all_runner_features

package worker

import (
	_ "github.com/busyster996/dagflow/internal/runner/docker"
	_ "github.com/busyster996/dagflow/internal/runner/git"
	_ "github.com/busyster996/dagflow/internal/runner/kubectl"
	_ "github.com/busyster996/dagflow/internal/runner/lua"
	_ "github.com/busyster996/dagflow/internal/runner/scp"
	_ "github.com/busyster996/dagflow/internal/runner/sftp"
	_ "github.com/busyster996/dagflow/internal/runner/ssh"
	_ "github.com/busyster996/dagflow/internal/runner/svn"
	_ "github.com/busyster996/dagflow/internal/runner/wasm"
	_ "github.com/busyster996/dagflow/internal/runner/yaegi"
)
