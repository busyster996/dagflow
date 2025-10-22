package wasm

import (
	"github.com/busyster996/dagflow/internal/runner"
	"github.com/busyster996/dagflow/internal/storage"
)

func init() {
	runner.Register("wasm", func(storage storage.IStep, subCmd, workspace, scriptDir string) (runner.IRunner, error) {
		return &sWasm{
			storage:   storage,
			workspace: workspace,
		}, nil
	})
}
