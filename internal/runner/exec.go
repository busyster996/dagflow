package runner

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/busyster996/dagflow/internal/common"
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
	storage      storage.IStep
	cmd          *exec.Cmd
	ctx          context.Context
	cancel       context.CancelFunc
	envPath      string
	shell        string
	workspace    string
	scriptPath   string
	codeFilePath string
	timeout      time.Duration
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
	_ = os.Remove(c.codeFilePath)
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
	code, err := c.execScript()
	// 从文件读取的退出码更准确（因为包装脚本会捕获真实退出码）
	code1, err1 := c.parseExitCodeFromFile()
	if err1 == nil {
		code = code1
	}
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

func (c *sCmd) execScript() (code int, err error) {
	ext := filepath.Ext(c.scriptPath)
	switch ext {
	case ".py", ".py3":
		return c.execPythonScript()
	default:
		return c.execSysScript()
	}
}

func (c *sCmd) randomFilename(prefix, ext string) string {
	h := sha256.Sum256([]byte(c.scriptPath))
	return fmt.Sprintf("%s-%x-%d%s", prefix, h[:8], time.Now().UnixNano(), ext)
}

func (c *sCmd) execPythonScript() (code int, err error) {
	// 使用 Python 包装脚本以正确透传标准输入
	wrapperContent := fmt.Sprintf(`# -*- coding: utf-8 -*-
import sys
import os
import runpy

exit_code = 0
script_path = r'%s'

try:
    runpy.run_path(script_path, run_name='__main__')
except SystemExit as e:
    exit_code = e.code if e.code is not None else 0
except Exception:
    import traceback
    traceback.print_exc()
    exit_code = 255
finally:
    code_file = os.environ.get('XEXEC_CODE_FILE_PATH')
    if code_file and not os.path.exists(code_file):
        try:
            with open(code_file, 'w') as f:
                f.write(str(exit_code if isinstance(exit_code, int) else 1))
        except:
            pass
    sys.exit(exit_code if isinstance(exit_code, int) else 1)
`, c.scriptPath)
	wrapperPath := filepath.Join(os.TempDir(), c.randomFilename("wrapper", ".py"))
	defer func() {
		_ = os.Remove(wrapperPath)
	}()
	if err = os.WriteFile(wrapperPath, []byte(wrapperContent), os.ModePerm); err != nil {
		return 255, err
	}
	return c.exec("python", "-u", wrapperPath)
}

func (c *sCmd) exec(command string, args ...string) (code int, err error) {
	c.cmd = exec.CommandContext(c.ctx, command, args...)
	c.cmd.Env = c.mergeEnv()
	if c.workspace != "" {
		c.cmd.Dir = c.workspace
	}
	c.beforeExec()
	stdout, err := c.cmd.StdoutPipe()
	if err != nil {
		return 255, err
	}
	stderr, err := c.cmd.StderrPipe()
	if err != nil {
		return 255, err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		c.consoleOutput("STDOUT", stdout)
	}()
	go func() {
		defer wg.Done()
		c.consoleOutput("STDERR", stderr)
	}()
	err = c.cmd.Run()
	wg.Wait()

	// 僵尸进程收割会触发no child processes
	if err != nil && strings.Contains(err.Error(), "waitid: no child processes") {
		err = nil
	}

	if err != nil {
		code = 255
	}
	if c.cmd.ProcessState != nil {
		code = c.cmd.ProcessState.ExitCode()
	}
	return code, err
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

func (c *sCmd) consoleOutput(title string, reader io.ReadCloser) {
	defer func() {
		_ = reader.Close()
	}()
	// 按行读取输出写入到日志中
	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		line = c.transform(line)
		c.storage.Log().Write(title, line)
	}
}

// parseExitCodeFromFile 按行读取文件，获取到退出码就结束
func (c *sCmd) parseExitCodeFromFile() (int, error) {
	data, err := os.ReadFile(c.codeFilePath)
	if err != nil {
		return 255, fmt.Errorf("read code file error: %w", err)
	}

	matches := digitRe.FindString(string(bytes.TrimSpace(data)))
	if matches == "" {
		return 255, fmt.Errorf("no valid exit code found")
	}

	return strconv.Atoi(matches)
}

// mergeEnv 合并系统环境变量和自定义环境变量，去除重复
// 自定义环境变量 s.env 优先级更高，会覆盖系统环境变量中的同名键
func (c *sCmd) mergeEnv() []string {
	// 创建 map 来存储环境变量，键是变量名，值是完整的 "KEY=VALUE" 字符串
	envMap := make(map[string]string)

	// 先添加系统环境变量
	for _, env := range syscall.Environ() {
		key := c.getEnvKey(env)
		if key != "" {
			envMap[key] = env
		}
	}

	// 用自定义环境变量覆盖（优先级更高）
	for _, env := range c.envs() {
		key := c.getEnvKey(env)
		if key != "" {
			envMap[key] = env
		}
	}

	// 转换回 []string
	result := make([]string, 0, len(envMap))
	for _, env := range envMap {
		result = append(result, env)
	}

	// 增加自定义退出码环境变量
	result = append(result, fmt.Sprintf("XEXEC_CODE_FILE_PATH=%s", c.codeFilePath))

	return result
}

// getEnvKey 从 "KEY=VALUE" 格式的字符串中提取键名
func (c *sCmd) getEnvKey(env string) string {
	idx := strings.IndexByte(env, '=')
	if idx == -1 {
		return ""
	}
	return env[:idx]
}
