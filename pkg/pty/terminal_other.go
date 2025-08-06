//go:build !windows

package pty

import (
	"context"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
)

type terminal struct {
	*os.File
	cmd *exec.Cmd
}

func (t *terminal) Resize(width, height int16) error {
	window := struct {
		row uint16
		col uint16
		x   uint16
		y   uint16
	}{
		uint16(height),
		uint16(width),
		0,
		0,
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		t.Fd(),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&window)),
	)
	if errno != 0 {
		return errno
	} else {
		return nil
	}
}

func (t *terminal) Wait(ctx context.Context) (uint32, error) {
	defer t.Close()
	if err := t.cmd.Wait(); err != nil {
		code := t.cmd.ProcessState.ExitCode()
		return uint32(code), err
	}
	code := t.cmd.ProcessState.ExitCode()
	return uint32(code), nil
}

func New(cmd string) (Terminal, error) {
	if cmd == "" {
		cmd = "bash"
	}
	c := exec.Command(cmd)
	f, err := pty.Start(c)
	if err != nil {
		return nil, err
	}
	return &terminal{
		File: f,
		cmd:  c,
	}, err
}
