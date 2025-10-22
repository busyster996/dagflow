package yaegi

import (
	"github.com/busyster996/dagflow/internal/runner"
	"github.com/busyster996/dagflow/internal/storage"
)

func init() {
	runner.Register("yaegi", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sYaegi{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
