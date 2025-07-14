package wasm

import (
	"context"

	"github.com/busyster996/dagflow/internal/storage"
)

type sWasm struct {
	storage   storage.IStep
	workspace string
}

func (w *sWasm) Run(ctx context.Context) (exit int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (w *sWasm) Clear() error {
	//TODO implement me
	panic("implement me")
}
