package svn

import (
	"context"

	"github.com/busyster996/dagflow/internal/storage"
)

type sSvn struct {
	storage    storage.IStep
	subCommand string
	workspace  string
}

func (s *sSvn) Run(ctx context.Context) (exit int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *sSvn) Clear() error {
	//TODO implement me
	panic("implement me")
}
