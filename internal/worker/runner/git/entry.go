package git

import (
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/worker/runner"
)

func init() {
	runner.Register("git", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sGit{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
