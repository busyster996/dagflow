package yaegi

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/traefik/yaegi/interp"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/runner/yaegi/libs"
	"github.com/busyster996/dagflow/internal/storage"
)

type sYaegi struct {
	interp    *interp.Interpreter
	storage   storage.IStep
	workspace string
}

func (y *sYaegi) Run(ctx context.Context) (exit common.ExecCode, err error) {
	defer func() {
		if _r := recover(); _r != nil {
			err = fmt.Errorf("panic during execution %v", _r)
			exit = common.ExecCodeSystemErr
			stack := debug.Stack()
			if _err, ok := _r.(error); ok && strings.Contains(_err.Error(), context.Canceled.Error()) {
				exit = common.ExecCodeKilled
				err = common.ExecErrKilled
			}
			y.storage.Log().Write(err.Error(), string(stack))
		}
	}()

	params, err := y.getParams()
	if err != nil {
		return common.ExecCodeFailed, err
	}

	if err = y.createVM(ctx); err != nil {
		return common.ExecCodeSystemErr, err
	}

	evalFnval, err := y.interp.EvalWithContext(ctx, "EvalCall")
	if err != nil {
		return common.ExecCodeFailed, err
	}
	evalFn, ok := evalFnval.Interface().(func(ctx context.Context, params gjson.Result))
	if !ok {
		return common.ExecCodeFailed, errors.New("not found EvalCall")
	}
	evalFn(ctx, params)
	return common.ExecCodeSuccess, nil
}

func (y *sYaegi) getParams() (gjson.Result, error) {
	var rawJSON string
	var err error
	taskEnv := y.storage.GlobalEnv().List()
	for _, v := range taskEnv {
		rawJSON, err = sjson.Set(rawJSON, v.Name, []byte(v.Value))
		if err != nil {
			return gjson.Result{}, err
		}
	}
	stepEnv := y.storage.Env().List()
	for _, v := range stepEnv {
		rawJSON, err = sjson.Set(rawJSON, v.Name, []byte(v.Value))
		if err != nil {
			return gjson.Result{}, err
		}
	}
	return gjson.Parse(rawJSON), nil
}

func (y *sYaegi) createVM(ctx context.Context) (err error) {
	y.interp = interp.New(interp.Options{
		Env: []string{
			fmt.Sprintf("WORKSPACE=%s", y.workspace),
		},
		Stdout: y.output(),
		Stderr: y.output(),
	})
	if err = y.interp.Use(libs.Symbols); err != nil {
		return err
	}
	y.interp.ImportUsed()

	content, err := y.storage.Content()
	if err != nil {
		return err
	}

	_, err = y.interp.EvalWithContext(ctx, content)
	if err != nil {
		return err
	}

	return
}

func (y *sYaegi) Clear() error {
	return nil
}

type sYeagiOutput struct {
	storage storage.IStep
}

func (s *sYeagiOutput) Write(p []byte) (n int, err error) {
	n = len(p)
	s.storage.Log().Write(string(p))
	return
}

func (y *sYaegi) output() io.Writer {
	return &sYeagiOutput{
		storage: y.storage,
	}
}
