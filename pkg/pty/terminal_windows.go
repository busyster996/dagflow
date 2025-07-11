//go:build windows

package pty

import (
	"github.com/busyster996/dagflow/pkg/pty/conpty"
)

type terminal struct {
	*conpty.ConPty
}

func New(cmd string) (Terminal, error) {
	if cmd == "" {
		cmd = "cmd.exe"
	}
	pty, err := conpty.Start(cmd)
	if err != nil {
		return nil, err
	}
	return &terminal{
		ConPty: pty,
	}, nil
}
