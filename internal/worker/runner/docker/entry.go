package docker

import (
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/worker/runner"
)

func init() {
	runner.Register("docker", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sDocker{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
