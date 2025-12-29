package runner

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/runner/xexec"
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/pkg/logx"
)

var (
	digitRe = regexp.MustCompile(`\d+`)
)

const (
	// 最大缓冲区大小 100MB
	maxCapacity = 100 * 1024 * 1024
)

type sCmd struct {
	storage    storage.IStep
	ctx        context.Context
	cancel     context.CancelFunc
	envPath    string
	shell      string
	workspace  string
	scriptPath string
	timeout    time.Duration
}

func (c *sCmd) scriptSuffix() string {
	switch c.shell {
	case "python", "python2", "python3", "py", "py2", "py3":
		return ".py"
	}
	return c.selfScriptSuffix()
}

func (c *sCmd) Clear() error {
	c.cancel()
	_ = os.Remove(c.scriptPath)
	return nil
}

func (c *sCmd) Run(ctx context.Context) (exit common.ExecCode, err error) {
	defer func() {
		c.cancel()
		if _r := recover(); _r != nil {
			err = fmt.Errorf("panic during execution %v", _r)
			exit = common.ExecCodeSystemErr
			stack := debug.Stack()
			c.storage.Log().Write(err.Error(), string(stack))
		}
	}()
	timeout, err := c.storage.Timeout()
	if err != nil {
		return common.ExecCodeFailed, err
	}
	c.ctx, c.cancel = context.WithCancel(ctx)
	if timeout > 0 {
		c.ctx, c.cancel = context.WithTimeoutCause(ctx, timeout, common.ExecErrTimeOut)
	}
	code, err := xexec.ExecScript(
		c.ctx, c.scriptPath,
		xexec.WithScriptEnv(c.envs()),
		xexec.WithScriptWorkdir(c.workspace),
		xexec.WithOutput(c.storage.Log().Write),
	)
	exit = common.ExecCode(code)
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

func (c *sCmd) envs() []string {
	var envs []string
	taskEnv := c.storage.GlobalEnv().List()
	for _, env := range taskEnv {
		envs = append(envs, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}
	stepEnv := c.storage.Env().List()
	for _, env := range stepEnv {
		envs = append(envs, fmt.Sprintf("%s=%s", env.Name, env.Value))
	}

	return append(envs,
		fmt.Sprintf("TASK_NAME=%s", c.storage.TaskName()),
		fmt.Sprintf("TASK_STEP_NAME=%s", c.storage.Name()),
		fmt.Sprintf("TASK_WORKSPACE=%s", c.workspace),
	)
}

func (c *sCmd) parseEnvFileFromFile() {
	// 打开源文件
	file, err := os.Open(c.envPath)
	if err != nil {
		logx.Warnln(err)
		return
	}
	defer func() {
		_ = file.Close()
	}()
	s := bufio.NewScanner(file)
	firstLine := true
	for s.Scan() {
		line := s.Text()

		if firstLine {
			firstLine = false
			// skip utf8 bom, powershell 5 legacy uses it for utf8
			if len(line) >= 3 && line[0] == 239 && line[1] == 187 && line[2] == 191 {
				line = line[3:]
			}
		}

		// 处理单行和多行环境变量
		singleLineEnv := strings.Index(line, "=")
		multiLineEnv := strings.Index(line, "<<")
		if singleLineEnv != -1 && (multiLineEnv == -1 || singleLineEnv < multiLineEnv) {
			// TODO: write to storage
			logx.Debugf("parsed env: %v=%v", line[:singleLineEnv], line[singleLineEnv+1:])
		} else if multiLineEnv != -1 {
			multiLineEnvContent := ""
			multiLineEnvDelimiter := line[multiLineEnv+2:]
			delimiterFound := false
			for s.Scan() {
				content := s.Text()
				if content == multiLineEnvDelimiter {
					delimiterFound = true
					break
				}
				if multiLineEnvContent != "" {
					multiLineEnvContent += "\n"
				}
				multiLineEnvContent += content
			}
			if !delimiterFound {
				logx.Errorf("invalid format delimiter '%v' not found before end of file", multiLineEnvDelimiter)
				return
			}
			// TODO: write to storage
			logx.Debugf("parsed env: %v=%v", line[:multiLineEnv], multiLineEnvContent)
		} else {
			logx.Errorf("invalid format '%v', expected a line with '=' or '<<'", line)
			return
		}
	}

	if err = s.Err(); err != nil {
		logx.Errorf("error reading file: %v", err)
		return
	}
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
