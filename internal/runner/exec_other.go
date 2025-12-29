//go:build !windows

package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"syscall"

	"github.com/busyster996/dagflow/pkg/logx"
)

func (c *sCmd) selfScriptSuffix() string {
	switch c.shell {
	case "ash":
		return ".ash"
	case "bash":
		return ".bash"
	case "csh":
		return ".csh"
	case "dash":
		return ".dash"
	case "ksh":
		return ".ksh"
	case "shell", "sh":
		return ".sh"
	case "tcsh":
		return ".tcsh"
	case "zsh":
		return ".zsh"
	default:
		return ".bash"
	}
}

func (c *sCmd) execSysScript() (code int, err error) {
	_ = os.Chmod(c.scriptPath, os.ModePerm)
	wrapperContent := fmt.Sprintf(`#!/bin/sh
_save_exit_code() {
	local code=$?
	[ ! -f "${XEXEC_CODE_FILE_PATH}" ] && echo $code > "${XEXEC_CODE_FILE_PATH}"
	exit $code
}
trap _save_exit_code EXIT INT TERM
"%s"
`, c.scriptPath)
	wrapperPath := filepath.Join(os.TempDir(), c.randomFilename("wrapper", ".sh"))
	defer func() {
		_ = os.Remove(wrapperPath)
	}()
	if err = os.WriteFile(wrapperPath, []byte(wrapperContent), os.ModePerm); err != nil {
		return 255, err
	}
	return c.exec("sh", wrapperPath)
}

func (c *sCmd) beforeExec() {
	c.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL,
	}
	c.cmd.Cancel = func() error {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorf("cmd cancel panic: %v\n%s", r, string(debug.Stack()))
			}
			_ = os.WriteFile(c.codeFilePath, []byte("137"), os.ModePerm)
		}()
		if c.cmd.Process == nil {
			return fmt.Errorf("process not started")
		}
		var errs error
		if err := c.cmd.Process.Kill(); err != nil {
			errs = errors.Join(errs, err)
		}
		if err := syscall.Kill(-c.cmd.Process.Pid, syscall.SIGKILL); err != nil {
			errs = errors.Join(errs, err)
		}
		if errs != nil {
			return fmt.Errorf("error cancelling pid=%d, errors=%v", c.cmd.Process.Pid, errs)
		}
		return nil
	}
}

func (c *sCmd) utf8ToGb2312(line string) string {
	return line
}

func (c *sCmd) transform(line string) string {
	return line
}
