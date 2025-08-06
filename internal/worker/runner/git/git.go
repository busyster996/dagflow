package git

import (
	"context"

	"github.com/busyster996/dagflow/internal/storage"
)

type sGit struct {
	storage    storage.IStep
	subCommand string
	workspace  string
}

func (s *sGit) Run(ctx context.Context) (exit int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *sGit) Clear() error {
	//TODO implement me
	panic("implement me")
}
