package runner

import (
	"context"

	"github.com/busyster996/dagflow/internal/common"
)

type IRunner interface {
	Run(ctx context.Context) (exit common.ExecCode, err error)
	Clear() error
}
