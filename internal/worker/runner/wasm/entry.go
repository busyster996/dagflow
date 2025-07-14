package wasm

import (
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/worker/runner"
)

func init() {
	runner.Register("wasm", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sWasm{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
