package ssh

import (
	"context"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/storage"
)

type sSsh struct {
	storage    storage.IStep
	subCommand string
	workspace  string
}

func (s *sSsh) Run(ctx context.Context) (exit common.ExecCode, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *sSsh) Clear() error {
	//TODO implement me
	panic("implement me")
}
