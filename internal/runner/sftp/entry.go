package sftp

import (
	"github.com/busyster996/dagflow/internal/runner"
	"github.com/busyster996/dagflow/internal/storage"
)

func init() {
	runner.Register("sftp", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sSftp{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
