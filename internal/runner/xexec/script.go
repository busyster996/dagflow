package xexec

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log/slog"
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
)

var (
	digitRe = regexp.MustCompile(`\d+`)
)

const (
	// 最大缓冲区大小 100MB
	maxCapacity = 100 * 1024 * 1024
)

// Logger 日志接口，允许自定义日志实现
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// defaultLogger 默认日志实现，使用 xlog
type defaultLogger struct{}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	slog.Info(fmt.Sprintf(format, args...))
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	slog.Error(fmt.Sprintf(format, args...))
}

type script struct {
	cmd          *exec.Cmd
	ctx          context.Context
	cancel       context.CancelFunc
	env          []string
	args         []string
	dir          string
	path         string
	stdin        io.Reader
	codeFilePath string
	logger       Logger
}

type ScriptOption func(*script)

func WithScriptWorkdir(path string) ScriptOption {
	return func(s *script) {
		s.dir = path
	}
}

func WithScriptEnv(env ...string) ScriptOption {
	return func(s *script) {
		s.env = env
	}
}

func WithScriptArgs(args ...string) ScriptOption {
	return func(s *script) {
		s.args = args
	}
}

func WithScriptStdin(stdin io.Reader) ScriptOption {
	return func(s *script) {
		s.stdin = stdin
	}
}

func WithScriptLogger(logger Logger) ScriptOption {
	return func(s *script) {
		s.logger = logger
	}
}

func ExecScriptContent(ctx context.Context, kind, content string, options ...ScriptOption) (code int, err error) {
	scriptPath := filepath.Join(os.TempDir(), fmt.Sprintf(".script_%d.%s", time.Now().UnixNano(), kind))
	defer func() {
		_ = os.Remove(scriptPath)
	}()
	err = os.WriteFile(scriptPath, []byte(content), os.ModePerm)
	if err != nil {
		return 255, err
	}
	return ExecScript(ctx, scriptPath, options...)
}

func ExecScript(ctx context.Context, path string, options ...ScriptOption) (code int, err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	s := &script{
		path:   path,
		logger: &defaultLogger{}, // 使用默认日志实现
	}
	s.codeFilePath = filepath.Join(os.TempDir(), s.randomFilename("exitcode", ".txt"))
	s.abs()
	_ = os.Chmod(path, os.ModePerm)
	_, err = os.Stat(path)
	if err != nil {
		return 255, fmt.Errorf("stat script: %w", err)
	}
	s.ctx, s.cancel = context.WithCancel(ctx)
	defer func() {
		s.cancel()
		_ = os.Remove(s.codeFilePath)
	}()
	for _, option := range options {
		option(s)
	}
	defer func() {
		if s.ctx.Err() != nil {
			err = s.ctx.Err()
		}
	}()
	code, err = s.execScript()
	// 从文件读取的退出码更准确（因为包装脚本会捕获真实退出码）
	code1, err1 := s.parseExitCodeFromFile()
	if err1 == nil {
		code = code1
	}
	return code, err
}

func (s *script) abs() {
	if filepath.IsAbs(s.path) {
		s.path = filepath.Clean(s.path)
	}
	s.path = filepath.Join(s.dir, s.path)
}

func (s *script) execScript() (code int, err error) {
	ext := filepath.Ext(s.path)
	switch ext {
	case ".py", ".py3":
		return s.execPythonScript()
	default:
		return s.execSysScript()
	}
}

func (s *script) randomFilename(prefix, ext string) string {
	h := sha256.Sum256([]byte(s.path))
	return fmt.Sprintf("%s-%x-%d%s", prefix, h[:8], time.Now().UnixNano(), ext)
}

func (s *script) execPythonScript() (code int, err error) {
	// 使用 Python 包装脚本以正确透传标准输入和参数
	wrapperContent := fmt.Sprintf(`# -*- coding: utf-8 -*-
import sys
import os
import runpy

exit_code = 0
script_path = r'%s'

# 设置 sys.argv，第一个参数是脚本路径，后续是传入的参数
sys.argv = [script_path] + sys.argv[1:]

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
`, s.path)
	wrapperPath := filepath.Join(os.TempDir(), s.randomFilename("wrapper", ".py"))
	defer func() {
		_ = os.Remove(wrapperPath)
	}()
	if err = os.WriteFile(wrapperPath, []byte(wrapperContent), os.ModePerm); err != nil {
		return 255, err
	}
	// 将包装脚本路径和用户参数合并传递给 python
	args := append([]string{"-u", wrapperPath}, s.args...)
	return s.exec("python", args...)
}

func (s *script) exec(command string, args ...string) (code int, err error) {
	s.cmd = exec.CommandContext(s.ctx, command, args...)
	s.cmd.Env = s.mergeEnv()
	if s.dir != "" {
		s.cmd.Dir = s.dir
	}
	s.beforeExec()
	if s.stdin != nil {
		s.cmd.Stdin = s.stdin
	}
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return 255, err
	}
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return 255, err
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.consoleOutput("STDOUT", stdout)
	}()
	go func() {
		defer wg.Done()
		s.consoleOutput("STDERR", stderr)
	}()
	err = s.cmd.Run()
	wg.Wait()

	// 僵尸进程收割会触发no child processes
	if err != nil && strings.Contains(err.Error(), "waitid: no child processes") {
		err = nil
	}

	if err != nil {
		code = 255
	}
	if s.cmd.ProcessState != nil {
		code = s.cmd.ProcessState.ExitCode()
	}
	return code, err
}

func (s *script) consoleOutput(title string, reader io.ReadCloser) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("[SYSTEM] consoleOutput panic:%v\n%s", r, string(debug.Stack()))
		}
		s.logger.Infof("[SYSTEM] stop %s console print", title)
	}()

	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		line = s.transform(line)
		s.logger.Infof("[%s] %s", title, line)
	}
}

// parseExitCodeFromFile 按行读取文件，获取到退出码就结束
func (s *script) parseExitCodeFromFile() (int, error) {
	data, err := os.ReadFile(s.codeFilePath)
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
func (s *script) mergeEnv() []string {
	// 创建 map 来存储环境变量，键是变量名，值是完整的 "KEY=VALUE" 字符串
	envMap := make(map[string]string)

	// 先添加系统环境变量
	for _, env := range syscall.Environ() {
		key := s.getEnvKey(env)
		if key != "" {
			envMap[key] = env
		}
	}

	// 用自定义环境变量覆盖（优先级更高）
	for _, env := range s.env {
		key := s.getEnvKey(env)
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
	result = append(result, fmt.Sprintf("XEXEC_CODE_FILE_PATH=%s", s.codeFilePath))

	return result
}

// getEnvKey 从 "KEY=VALUE" 格式的字符串中提取键名
func (s *script) getEnvKey(env string) string {
	idx := strings.IndexByte(env, '=')
	if idx == -1 {
		return ""
	}
	return env[:idx]
}
