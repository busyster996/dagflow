package pty

import (
	"context"
)

type Terminal interface {
	Read([]byte) (int, error)
	Write(p []byte) (int, error)
	Resize(width, height int16) error
	Wait(ctx context.Context) (uint32, error)
}
