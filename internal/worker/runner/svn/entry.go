package svn

import (
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/worker/runner"
)

func init() {
	runner.Register("svn", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sSvn{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
