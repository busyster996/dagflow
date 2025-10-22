package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/pubsub"
	"github.com/busyster996/dagflow/internal/server/router/base"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/pkg/logx"
)

type SStepService struct {
	taskName string
	stepName string
}

func Step(taskName string, stepName string) *SStepService {
	return &SStepService{
		taskName: taskName,
		stepName: stepName,
	}
}

func (ss *SStepService) Create(seqNo int64, globalTimeout time.Duration, step *types.SStepReq) error {
	if globalTimeout <= 0 {
		return errors.New("global timeout must be greater than 0")
	}
	if err := ss.review(step); err != nil {
		logx.Errorln("step review", ss.taskName, ss.stepName, err)
		return err
	}
	// 确保超时时间在合理范围内
	if step.Timeout <= 0 || step.Timeout > globalTimeout {
		step.Timeout = globalTimeout
	}
	if err := ss.saveStep(seqNo, step); err != nil {
		logx.Errorln("step save", ss.taskName, ss.stepName, err)
		return err
	}
	return nil
}

func (ss *SStepService) review(step *types.SStepReq) error {
	step.Name = reg.ReplaceAllString(step.Name, "")
	if step.Name == "" {
		step.Name = ksuid.New().String()
	}
	ss.stepName = step.Name

	if step.Type == "" {
		return errors.New("step type is empty")
	}

	if step.Content == "" {
		return errors.New("step content is empty")
	}

	// 校验env是否重复
	var envKeys []string
	for _, v := range step.Env {
		envKeys = append(envKeys, v.Name)
	}
	dup := utility.CheckDuplicate(envKeys)
	if dup != nil {
		return fmt.Errorf("duplicate key %v", dup)
	}

	step.Depends = utility.RemoveDuplicate(step.Depends)
	return nil
}

func (ss *SStepService) saveStep(seqNo int64, step *types.SStepReq) (err error) {
	stepStorage := storage.Task(ss.taskName).Step(step.Name)
	defer func() {
		if err != nil {
			_ = stepStorage.ClearAll()
		}
	}()
	data := &models.SStep{
		TaskName: ss.taskName,
		Name:     step.Name,
		Desc:     step.Desc,
		Type:     step.Type,
		Content:  step.Content,
		Action:   step.Action,
		Rule:     step.Rule,
		SeqNo:    seqNo,
		Timeout:  step.Timeout,
		Disable:  models.Pointer(step.Disable),
		SStepUpdate: models.SStepUpdate{
			Message:  "the step is waiting to be scheduled for execution",
			Code:     models.Pointer(common.ExecCode(0)),
			State:    models.Pointer(models.StatePending),
			OldState: models.Pointer(models.StatePending),
		},
	}
	if step.RetryPolicy != nil {
		data.RetryPolicy = datatypes.NewJSONType(models.SRetryPolicy{
			Interval:    step.RetryPolicy.Interval,
			MaxAttempts: step.RetryPolicy.MaxAttempts,
			MaxInterval: step.RetryPolicy.MaxInterval,
			Multiplier:  step.RetryPolicy.Multiplier,
		})
	}
	err = storage.Task(ss.taskName).StepCreate(data)
	if err != nil {
		logx.Errorln("step save", ss.taskName, ss.stepName, err)
		return fmt.Errorf("save step error: %s", err)
	}
	// save step env
	var envs models.SEnvs
	for _, env := range step.Env {
		envs = append(envs, &models.SEnv{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	if err = stepStorage.Env().Insert(envs...); err != nil {
		logx.Errorln("step save envs", ss.taskName, ss.stepName, err)
		return fmt.Errorf("save step env error: %s", err)
	}
	// save a step depend
	err = stepStorage.Depend().Insert(step.Depends...)
	if err != nil {
		logx.Errorln("step save depends", ss.taskName, ss.stepName, err)
		return fmt.Errorf("save step depend error: %s", err)
	}
	return
}

func (ss *SStepService) Detail() (base.Code, *types.SStepRes, error) {
	stepStorage := storage.Task(ss.taskName).Step(ss.stepName)
	step, err := stepStorage.Get()
	if err != nil {
		logx.Errorln("step detail", ss.taskName, ss.stepName, err)
		return base.CodeFailed, nil, errors.New("step not found")
	}
	data := &types.SStepRes{
		Name:    step.Name,
		Desc:    step.Desc,
		State:   models.StateMap[*step.State],
		Code:    step.Code.Int64(),
		Message: step.Message,
		Timeout: step.Timeout,
		Disable: *step.Disable,
		Type:    step.Type,
		Content: step.Content,
		Action:  step.Action,
		Rule:    step.Rule,
		RetryPolicy: &types.SRetryPolicy{
			Interval:    step.RetryPolicy.Data().Interval,
			MaxAttempts: step.RetryPolicy.Data().MaxAttempts,
			MaxInterval: step.RetryPolicy.Data().MaxInterval,
			Multiplier:  step.RetryPolicy.Data().Multiplier,
		},
		Time: &types.STimeRes{
			Start: step.STimeStr(),
			End:   step.ETimeStr(),
		},
	}
	data.Depends = storage.Task(ss.taskName).Step(step.Name).Depend().List()
	envs := stepStorage.Env().List()
	for _, env := range envs {
		data.Env = append(data.Env, &types.SEnv{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	return base.Code(data.Code), data, nil
}

func (ss *SStepService) Manager(action string, duration string) error {
	task, err := storage.Task(ss.taskName).Get()
	if err != nil {
		logx.Errorln("step manager", ss.taskName, ss.stepName, err)
		return errors.New("task not found")
	}
	if *task.State != models.StateRunning && *task.State != models.StatePending && *task.State != models.StatePaused {
		return errors.New("task is no running")
	}
	step, err := storage.Task(ss.taskName).Step(ss.stepName).Get()
	if err != nil {
		logx.Errorln("step manager", ss.taskName, ss.stepName, err)
		return errors.New("step not found")
	}
	if *step.State != models.StateRunning && *step.State != models.StatePending && *step.State != models.StatePaused {
		return errors.New("step is no running")
	}
	return pubsub.PublishManager(task.Node, utility.JoinWithInvisibleChar(ss.taskName, ss.stepName, action, duration))
}

func (ss *SStepService) Delete() error {
	return storage.Task(ss.taskName).Step(ss.stepName).ClearAll()
}

func (ss *SStepService) Log() (base.Code, types.SStepLogsRes, error) {
	step, err := storage.Task(ss.taskName).Step(ss.stepName).Get()
	if err != nil {
		logx.Errorln("step log", ss.taskName, ss.stepName, err)
		return base.CodeFailed, nil, errors.New("step not found")
	}
	switch *step.State {
	case models.StatePending:
		return base.CodePending, types.SStepLogsRes{
			{
				Timestamp: time.Now().UnixNano(),
				Line:      1,
				Content:   step.Message,
			},
		}, errors.New(step.Message)
	case models.StatePaused:
		return base.CodePaused, types.SStepLogsRes{
			{
				Timestamp: time.Now().UnixNano(),
				Line:      1,
				Content:   "step is paused",
			},
		}, errors.New(step.Message)
	default:
		res, _ := ss.log(nil)
		return ConvertState(*step.State), res, errors.New(step.Message)
	}
}

func (ss *SStepService) log(latestLine *int64) (res types.SStepLogsRes, done bool) {
	logs := storage.Task(ss.taskName).Step(ss.stepName).Log().List(latestLine)
	for _, v := range logs {
		if v.Content == common.ExecConsoleStart {
			continue
		}
		if v.Content == common.ExecConsoleDone {
			done = true
			continue
		}
		res = append(res, &types.SStepLogRes{
			Timestamp: v.Timestamp,
			Line:      *v.Line,
			Content:   v.Content,
		})
	}
	// 如果查询到有新日志，更新 latestLine 为最后一条日志的行号
	if len(logs) > 0 && latestLine != nil {
		*latestLine = *logs[len(logs)-1].Line // 更新 latestLine
	}
	return
}

type outputFn func(code base.Code, data types.SStepLogsRes, err error) error
type stateHandlerFn func(output outputFn, latest *int64) (bool, error)

func (ss *SStepService) LogStream(ctx context.Context, output outputFn) error {
	db := storage.Task(ss.taskName).Step(ss.stepName)
	step, err := db.Get()
	if err != nil {
		logx.Errorln("step logstream", ss.taskName, ss.stepName, err)
		return errors.New("step not found")
	}

	var latestLine int64
	// 用于防止某些状态下的重复推送
	var onceMap = map[models.State]*sync.Once{
		models.StatePending: new(sync.Once),
		models.StatePaused:  new(sync.Once),
		models.StateUnknown: new(sync.Once),
	}
	// 状态处理函数映射
	handlers := map[models.State]stateHandlerFn{
		models.StatePending: ss.createOnceHandler(onceMap[models.StatePending], base.CodePending, "step is pending"),
		models.StatePaused:  ss.createOnceHandler(onceMap[models.StatePaused], base.CodePaused, "step is paused"),
		models.StateUnknown: ss.createOnceHandler(onceMap[models.StateUnknown], base.CodeNoData, "step status unknown"),
		models.StateRunning: ss.handleRunningState,
		models.StateStopped: ss.handleFinalState(base.CodeSuccess),
		models.StateFailed:  ss.handleFinalState(base.CodeFailed),
		models.StateSkipped: ss.handleFinalState(base.CodeSkipped),
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		// 避免多次推送
		// Paused, Pending,Unknown 状态只发送一次, 然后继续等到状态变化再继续推送
		// Running 状态会一直推送, 直到状态推送完成.
		// Stop, Failed 推送后结束.

		if handler, exists := handlers[*step.State]; exists {
			var shouldContinue bool
			shouldContinue, err = handler(output, &latestLine)
			if err != nil {
				logx.Errorln("step logstream", ss.taskName, ss.stepName, err)
				return err
			}
			if !shouldContinue {
				return nil
			}
		} else {
			return errors.New("unhandled step state")
		}

		step, err = db.Get()
		if err != nil {
			logx.Errorln("step logstream", ss.taskName, ss.stepName, err)
			return errors.New("step not found")
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func (ss *SStepService) createOnceHandler(once *sync.Once, code base.Code, message string) stateHandlerFn {
	return func(output outputFn, latest *int64) (bool, error) {
		once.Do(func() {
			_ = output(code, types.SStepLogsRes{
				{
					Timestamp: time.Now().UnixNano(),
					Line:      1,
					Content:   message,
				},
			}, nil)
		})
		return true, nil
	}
}

func (ss *SStepService) handleRunningState(output outputFn, latestLine *int64) (bool, error) {
	res, done := ss.log(latestLine)
	err := output(base.CodeRunning, res, errors.New("in progress"))
	if err != nil {
		logx.Errorln("step logstream", ss.taskName, ss.stepName, err)
		return false, err
	}
	if done {
		return false, nil
	}
	return true, nil
}

func (ss *SStepService) handleFinalState(code base.Code) stateHandlerFn {
	return func(output outputFn, latestLine *int64) (bool, error) {
		db := storage.Task(ss.taskName).Step(ss.stepName)
		step, err := db.Get()
		if err != nil {
			logx.Errorln("step logstream", ss.taskName, ss.stepName, err)
			return false, errors.New("step not found")
		}
		res, _ := ss.log(latestLine)
		var errMsg error
		if code == base.CodeFailed {
			errMsg = fmt.Errorf("exit code: %d", step.Code)
			if step.Message != "" {
				errMsg = errors.New(step.Message)
			}
		}
		err = output(code, res, errMsg)
		if err != nil {
			logx.Errorln("step logstream", ss.taskName, ss.stepName, err)
			return false, err
		}
		return false, nil
	}
}
