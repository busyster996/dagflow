//go:build !windows

package xexec

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"syscall"
)

func (s *script) execSysScript() (code int, err error) {
	_ = os.Chmod(s.path, os.ModePerm)
	wrapperContent := fmt.Sprintf(`#!/bin/sh
_save_exit_code() {
	local code=$?
	[ ! -f "${XEXEC_CODE_FILE_PATH}" ] && echo $code > "${XEXEC_CODE_FILE_PATH}"
	exit $code
}
trap _save_exit_code EXIT INT TERM
"%s" "$@"
`, s.path)
	wrapperPath := filepath.Join(os.TempDir(), s.randomFilename("wrapper", ".sh"))
	defer func() {
		_ = os.Remove(wrapperPath)
	}()
	if err = os.WriteFile(wrapperPath, []byte(wrapperContent), os.ModePerm); err != nil {
		return 255, err
	}
	// 将包装脚本路径和用户参数合并传递给 sh
	args := append([]string{wrapperPath}, s.args...)
	return s.exec("sh", args...)
}

func (s *script) beforeExec() {
	s.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
	}
	s.cmd.Cancel = func() error {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Errorf("[SYSTEM] cmd cancel panic: %v\n%s", r, string(debug.Stack()))
			}
			_ = os.WriteFile(s.codeFilePath, []byte("137"), os.ModePerm)
		}()
		if s.cmd.Process == nil {
			return fmt.Errorf("process not started")
		}
		var errs error
		if err := s.cmd.Process.Kill(); err != nil {
			errs = errors.Join(errs, err)
		}
		if err := syscall.Kill(-s.cmd.Process.Pid, syscall.SIGKILL); err != nil {
			errs = errors.Join(errs, err)
		}
		if errs != nil {
			return fmt.Errorf("error cancelling pid=%d, errors=%v", s.cmd.Process.Pid, errs)
		}
		return nil
	}
}

func (s *script) utf8ToGb2312(line string) string {
	return line
}

func (s *script) transform(line string) string {
	return line
}
