package entv1

import (
	"context"
	"runtime/debug"

	"github.com/busyster996/dagflow/pkg/logx"
)

//go:generate go run -mod=mod entgo.io/ent/cmd/ent@latest generate --feature sql/lock,sql/modifier,sql/execquery,sql/upsert --template ./setstruct.tmpl ./schema

func SafeOnlyX[T any](ctx context.Context, fn func(ctx context.Context) T) T {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			logx.Errorln(r, string(stack))
		}
	}()
	return fn(ctx)
}

func Pointer[T any](v T) *T {
	return &v
}

func UnPointer[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
