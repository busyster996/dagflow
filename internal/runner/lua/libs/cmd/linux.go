//go:build linux

package cmd

import (
	"os"
	"os/exec"
	"syscall"
)

var BeforeExec = []func(cmd *exec.Cmd){
	func(cmd *exec.Cmd) {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Pgid: os.Getpid(),
			//Setpgid:   true,
			Setsid:    true,
			Pdeathsig: syscall.SIGKILL,
		}
	},
}
