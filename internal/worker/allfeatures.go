//go:build dagflow.all_runner_features

package worker

import (
	_ "github.com/busyster996/dagflow/internal/worker/runner/docker"
	_ "github.com/busyster996/dagflow/internal/worker/runner/git"
	_ "github.com/busyster996/dagflow/internal/worker/runner/kubectl"
	_ "github.com/busyster996/dagflow/internal/worker/runner/lua"
	_ "github.com/busyster996/dagflow/internal/worker/runner/scp"
	_ "github.com/busyster996/dagflow/internal/worker/runner/sftp"
	_ "github.com/busyster996/dagflow/internal/worker/runner/ssh"
	_ "github.com/busyster996/dagflow/internal/worker/runner/svn"
	_ "github.com/busyster996/dagflow/internal/worker/runner/wasm"
	_ "github.com/busyster996/dagflow/internal/worker/runner/yaegi"
)
