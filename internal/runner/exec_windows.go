//go:build windows

package runner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/busyster996/dagflow/pkg/logx"
)

var (
	acp = windows.GetACP()
)
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


func (c *sCmd) execSysScript() (code int, err error) {
	ext := filepath.Ext(c.scriptPath)
	var (
		command        string
		args           []string
		wrapperPath    string
		wrapperContent string
	)
	switch ext {
	case ".ps1", ".ps":
		command = "powershell.exe"
		args = []string{"-NoLogo", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File"}
		wrapperPath = filepath.Join(os.TempDir(), c.randomFilename("wrapper", ".ps1"))
		wrapperContent = fmt.Sprintf(`$ErrorActionPreference='Continue'
try {
    & "%s"
    if ($LASTEXITCODE -eq $null) {
        $exitCode = 0
    } else {
        $exitCode = $LASTEXITCODE
    }
} catch {
    $exitCode = 255
    Write-Error "[$($_.Exception.GetType().FullName)] $($_.Exception.Message)"
} finally {
    if (-not (Test-Path "$env:XEXEC_CODE_FILE_PATH")) {
        $exitCode | Out-File -FilePath "$env:XEXEC_CODE_FILE_PATH" -Encoding ASCII
    }
    exit $exitCode
}`, c.scriptPath)
	default:
		command = "cmd.exe"
		args = []string{"/C"}
		wrapperPath = filepath.Join(os.TempDir(), c.randomFilename("wrapper", ".bat"))
		wrapperContent = fmt.Sprintf(`@echo off
call "%s"
set exitcode=%%ERRORLEVEL%%
if not exist "%%XEXEC_CODE_FILE_PATH%%" (
    echo %%exitcode%% > "%%XEXEC_CODE_FILE_PATH%%"
)
exit /b %%exitcode%%`, c.scriptPath)
	}
	defer func() {
		_ = os.Remove(wrapperPath)
	}()
	if err = os.WriteFile(wrapperPath, []byte(wrapperContent), os.ModePerm); err != nil {
		return 255, err
	}
	args = append(args, wrapperPath)
	return c.exec(command, args...)
}

func (c *sCmd) beforeExec() {
	c.cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		HideWindow:    true,
	}
	c.cmd.Cancel = func() error {
		defer func() {
			if r := recover(); r != nil {
				logx.Errorln("cmd cancel panic: %v\n%s", r, string(debug.Stack()))
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
		kill := exec.Command("TASKKILL.exe", "/T", "/F", "/PID", strconv.Itoa(c.cmd.Process.Pid))
		if err := kill.Run(); err != nil {
			errs = errors.Join(errs, err)
		}
		if errs != nil {
			return fmt.Errorf("error cancelling pid=%d, errors=%v", c.cmd.Process.Pid, errs)
		}
		return nil
	}
}

func (c *sCmd) transform(line string) string {
	if c.isGBK(line) || acp == 936 {
		line = c.gbkToUtf8(line)
	}
	return line
}

func (c *sCmd) gbkToUtf8(line string) string {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("gbkToUtf8 panic:%v\n%s", r, string(debug.Stack()))
		}
	}()
	reader := transform.NewReader(strings.NewReader(line), simplifiedchinese.GBK.NewDecoder())
	b, err := io.ReadAll(reader)
	if err != nil {
		logx.Errorln(err)
		return line
	}
	return string(b)
}

func (c *sCmd) isGBK(data string) bool {
	defer func() {
		if r := recover(); r != nil {
			logx.Errorf("isGBK panic:%v\n%s", r, string(debug.Stack()))
		}
	}()

	length := len(data)
	for i := 0; i < length; {
		if data[i] <= 0x7f {
			i++
			continue
		}

		// 边界检查
		if i+1 >= length {
			return false
		}

		// GBK 双字节检测
		if data[i] >= 0x81 && data[i] <= 0xfe &&
			data[i+1] >= 0x40 && data[i+1] <= 0xfe &&
			data[i+1] != 0xf7 {
			i += 2
			continue
		}
		return false
	}
	return true
}
