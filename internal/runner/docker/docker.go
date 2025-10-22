package docker

import (
	"context"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/storage"
)

type sDocker struct {
	storage    storage.IStep
	subCommand string
	workspace  string
}

func (s *sDocker) Run(ctx context.Context) (exit common.ExecCode, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *sDocker) Clear() error {
	//TODO implement me
	panic("implement me")
}
