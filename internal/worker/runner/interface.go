package runner

import (
	"context"
)

type IRunner interface {
	Run(ctx context.Context) (exit int64, err error)
	Clear() error
}
