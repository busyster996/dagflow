package docker

import (
	"github.com/busyster996/dagflow/internal/runner"
	"github.com/busyster996/dagflow/internal/storage"
)

func init() {
	runner.Register("docker", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sDocker{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
