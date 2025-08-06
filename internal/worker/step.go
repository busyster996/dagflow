package worker

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/expr-lang/expr"
	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/internal/worker/common"
	"github.com/busyster996/dagflow/internal/worker/event"
	"github.com/busyster996/dagflow/internal/worker/runner"
	"github.com/busyster996/dagflow/pkg/dagcuter"
	"github.com/busyster996/dagflow/pkg/logx"
)

type sStep struct {
	// 生命周期控制（强杀）
	lcCtx    context.Context
	lcCancel context.CancelFunc

	// 控制上下文, 控制挂起或解卦
	ctrlCtx    context.Context
	ctrlCancel context.CancelFunc

	stg       storage.IStep
	kind      string
	taskName  string
	stepName  string
	workspace string
	scriptDir string
	state     int32 // 0: 正常, 1: 挂起
}

func (s *sStep) Name() string {
	return fmt.Sprintf("%s/%s", s.taskName, s.stepName)
}

func (s *sStep) Dependencies() []string {
	return s.stg.Depend().List()
}

func (s *sStep) RetryPolicy() *dagcuter.RetryPolicy {
	retryPolicy, err := s.stg.RetryPolicy()
	if err != nil {
		logx.Errorln(s.taskName, s.stepName, "RetryPolicy error:", err)
		return nil
	}
	return &dagcuter.RetryPolicy{
		Interval:    retryPolicy.Interval,
		MaxAttempts: retryPolicy.MaxAttempts,
		MaxInterval: retryPolicy.MaxInterval,
		Multiplier:  retryPolicy.Multiplier,
	}
}

func (s *sStep) PreExecution(ctx context.Context, input map[string]any) error {
	logx.Infoln(s.taskName, s.stepName, s.workspace, "PreExecution")
	event.SendEventf("%s %s PreExecution", s.taskName, s.stepName)
	return nil
}

func (s *sStep) Execute(ctx context.Context, input map[string]any) (map[string]any, error) {
	if err := s.checkCtx(ctx); err != nil {
		return nil, err
	}
	logx.Infoln(s.taskName, s.stepName, s.workspace, "Execute")
	event.SendEventf("%s %s Execute", s.taskName, s.stepName)
	var err error
	if err = s.stg.Update(&models.SStepUpdate{
		State:    models.Pointer(models.StateRunning),
		OldState: models.Pointer(models.StatePending),
		Message:  "step is running",
		STime:    models.Pointer(time.Now()),
	}); err != nil {
		logx.Errorln(s.taskName, s.stepName, err)
		event.SendEventf("%s %s error %v", s.taskName, s.stepName, err)
		return nil, err
	}

	// proc step
	var res = new(models.SStepUpdate)
	defer func() {
		if _err := recover(); _err != nil {
			stack := string(debug.Stack())
			logx.Errorln(stack, _err)
			res.Code = models.Pointer(common.CodeSystemErr)
			res.Message = fmt.Sprint(_err)
		}
		res.ETime = models.Pointer(time.Now())
		res.OldState = models.Pointer(models.StateRunning)
		if _err := s.stg.Update(res); _err != nil {
			logx.Errorln(_err)
		}
		event.SendEventf("%s %s %v", s.taskName, s.stepName, res.Message)
	}()

	// 日志写入
	s.stg.Log().Write(common.ConsoleStart)
	defer s.stg.Log().Write(common.ConsoleDone)

	if s.kind == common.KindStrategy {
		// 评估规则, 使用expr
		var action common.Action
		action, err = s.evaluateExprRule()
		if err != nil {
			logx.Errorln(s.taskName, s.stepName, err)
			// 写入步骤日志
			s.stg.Log().Write(err.Error())
			res.State = models.Pointer(models.StateFailed)
			res.Code = models.Pointer(common.CodeSystemErr)
			res.Message = err.Error()
			return nil, err
		}
		switch action {
		case common.ActionSkip:
			res.State = models.Pointer(models.StateSkipped)
			res.Code = models.Pointer(common.CodeSkipped)
			res.Message = "skipped due to rule"
			return nil, nil
		default:
			// 继续执行
		}
		defer func() {
			// 策略模式下需要当前步骤成功才会触发
			err = nil
		}()
	}
	var _runner runner.IRunner
	_runner, err = s.newExecutorRunner()
	if err != nil {
		logx.Errorln(s.taskName, s.stepName, err)
		res.State = models.Pointer(models.StateFailed)
		res.Message = err.Error()
		res.Code = models.Pointer(common.CodeSystemErr)
		return nil, err
	}

	defer func() {
		if cErr := _runner.Clear(); cErr != nil {
			logx.Warnln(cErr)
		}
	}()

	res.Message = "execution succeed"
	var code int64
	_ctx, cancel := utils.MergerContext(ctx, s.lcCtx)
	defer cancel()
	code, err = _runner.Run(_ctx)
	res.Code = models.Pointer(code)
	res.State = models.Pointer(models.StateStopped)
	if err != nil {
		res.State = models.Pointer(models.StateFailed)
		res.Message = err.Error()
		return nil, err
	}
	if code != 0 {
		res.State = models.Pointer(models.StateFailed)
		res.Message = fmt.Sprintf("execution failed with code: %d", code)
		return nil, errors.New(res.Message)
	}
	return nil, nil
}

func (s *sStep) PostExecution(ctx context.Context, output map[string]any) error {
	logx.Infoln(s.taskName, s.stepName, s.workspace, "PostExecution")
	event.SendEventf("%s %s PostExecution", s.taskName, s.stepName)
	stepManager.Delete(s.Name())
	return nil
}

func (s *sStep) Stop() {
	defer func() {
		if _err := recover(); _err != nil {
			logx.Errorln(_err)
		}
	}()
	if s.lcCancel != nil {
		logx.Infoln(s.taskName, s.stepName, s.workspace, "Stop")
		event.SendEventf("%s %s Stop", s.taskName, s.stepName)
		s.lcCancel()
	}
	stepManager.Delete(s.Name())
}

func (s *sStep) checkCtx(ctx context.Context) error {
	// 挂起, 则等待解挂
	if s.ctrlCtx != nil {
		// 等待控制信号
		select {
		case <-s.ctrlCtx.Done():
		case <-s.lcCtx.Done():
			return s.lcCtx.Err()
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}
	if s.lcCtx.Err() != nil {
		return s.lcCtx.Err()
	}
	return nil
}

func (s *sStep) evaluateExprRule() (common.Action, error) {
	// 查询规则
	rule, err := s.stg.Rule()
	if err != nil {
		logx.Errorln(s.taskName, s.stepName, err)
		return common.ActionUnknown, err
	}
	action, err := s.stg.Action()
	if err != nil {
		logx.Errorln(s.taskName, s.stepName, err)
		return common.ActionUnknown, err
	}
	if rule == "" || action == "" {
		logx.Infoln(s.taskName, s.stepName, "no rule or no action")
		return common.ActionAllow, nil
	}
	program, err := expr.Compile(rule, s.exprBuiltins()...)
	if err != nil {
		logx.Errorln(s.taskName, s.stepName, err)
		return common.ActionUnknown, err
	}
	result, err := expr.Run(program, nil)
	if err != nil {
		logx.Errorln(s.taskName, s.stepName, err)
		return common.ActionUnknown, err
	}
	matched, ok := result.(bool)
	if !ok {
		return common.ActionUnknown, fmt.Errorf("rule result is not a boolean")
	}
	// 如果不匹配, 继续执行
	if !matched {
		return common.ActionAllow, nil
	}
	return common.ActionConvert(action), nil
}
