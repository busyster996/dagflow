//go:build windows

package cmd

import (
	"os/exec"
	"syscall"
)

var BeforeExec = []func(cmd *exec.Cmd){
	func(cmd *exec.Cmd) {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
			HideWindow:    true,
		}
	},
}
