package k8s

import (
	"github.com/busyster996/dagflow/storage"
	"github.com/busyster996/dagflow/worker/runner"
)

func init() {
	runner.Register("kubectl", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sKubectl{
			storage:    storage,
			subCommand: subCmd,
			workspace:  workspace,
		}, nil
	})
}
