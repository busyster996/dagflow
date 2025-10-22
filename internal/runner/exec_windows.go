//go:build windows

package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/busyster996/dagflow/internal/common"
)

var acp = windows.GetACP()

func (c *sCmd) selfScriptSuffix() string {
	switch c.shell {
	case "cmd", "bat":
		return ".bat"
	case "powershell", "ps", "ps1":
		return ".ps1"
	default:
		return ".bat"
	}
}

func (c *sCmd) selfCmd() *exec.Cmd {
	var cmd *exec.Cmd
	switch c.shell {
	case "cmd", "bat":
		cmd = exec.CommandContext(c.ctx, "cmd", "/D", "/E:ON", "/V:OFF", "/Q", "/S", "/C", c.scriptName)
	case "powershell", "ps", "ps1":
		// 解决用户不写exit时, powershell进程外获取不到退出码
		command := fmt.Sprintf("$ErrorActionPreference='Continue';%s;exit $LASTEXITCODE", c.scriptName)
		// 激进方式, 强制用户脚本没问题
		// command := fmt.Sprintf("$ErrorActionPreference='Stop';%s;exit $LASTEXITCODE", c.absFilePath)
		cmd = exec.CommandContext(c.ctx, "powershell", "-NoLogo", "-NonInteractive", "-Command", command)
	default:
		cmd = exec.CommandContext(c.ctx, "cmd", "/D", "/E:ON", "/V:OFF", "/Q", "/S", "/C", c.scriptName)
	}
	return cmd
}

func (c *sCmd) Run(ctx context.Context) (exit common.ExecCode, err error) {
	defer func() {
		c.cancel()
		if _r := recover(); _r != nil {
			err = fmt.Errorf("panic during execution %v", _r)
			exit = common.ExecCodeSystemErr
			stack := debug.Stack()
			if _err, ok := _r.(error); ok && strings.Contains(_err.Error(), context.Canceled.Error()) {
				exit = common.ExecCodeKilled
				err = common.ExecErrKilled
			}
			c.storage.Log().Write(err.Error(), string(stack))
		}
	}()

	cmd, err := c.newCmd(ctx)
	if err != nil {
		c.storage.Log().Write(err.Error())
		return common.ExecCodeSystemErr, err
	}
	// 设置输出

	reader, err := cmd.StdoutPipe()
	if err != nil {
		c.storage.Log().Write(err.Error())
		return common.ExecCodeSystemErr, err
	}
	cmd.Stderr = cmd.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		HideWindow:    true,
	}
	go c.copyOutput(reader)
	err = cmd.Run()
	if cmd.ProcessState != nil {
		exit = common.ExecCode(cmd.ProcessState.ExitCode())
		if cmd.ProcessState.Pid() != 0 {
			_ = c.kill(cmd.ProcessState.Pid())
		}
	}
	if err != nil && exit == 0 {
		exit = common.ExecCodeFailed
	}
	if c.ctx.Err() != nil {
		switch {
		case errors.Is(context.Cause(c.ctx), common.ExecErrTimeOut):
			err = common.ExecErrTimeOut
			exit = common.ExecCodeTimeout
		default:
			err = common.ExecErrKilled
			exit = common.ExecCodeKilled
		}
	}
	return
}

func (c *sCmd) kill(pid int) error {
	if pid == 0 {
		return nil
	}
	kill := exec.Command("TASKKILL.exe", "/T", "/F", "/PID", strconv.Itoa(pid))
	return kill.Run()
}

func (c *sCmd) transform(line string) string {
	if c.isGBK(line) || acp == 936 {
		line = string(c.gbkToUtf8([]byte(line)))
	}
	return line
}

func (c *sCmd) gbkToUtf8(s []byte) []byte {
	defer func() {
		recover()
	}()
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	b, err := io.ReadAll(reader)
	if err != nil {
		return s
	}
	return b
}

func (c *sCmd) utf8ToGb2312(s string) string {
	defer func() {
		recover()
	}()
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		return s
	}

	return string(d)
}

func (c *sCmd) isGBK(data string) bool {
	defer func() {
		recover()
	}()
	length := len(data)
	var i = 0
	for i < length {
		if data[i] <= 0x7f {
			// 编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			// 大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}
