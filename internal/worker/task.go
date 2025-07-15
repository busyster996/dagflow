package worker

import (
	"context"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/internal/worker/common"
	"github.com/busyster996/dagflow/internal/worker/event"
	"github.com/busyster996/dagflow/pkg/dagcuter"
	"github.com/busyster996/dagflow/pkg/logx"
)

type sTask struct {
	// 生命周期控制（强杀）
	lcCtx    context.Context
	lcCancel context.CancelFunc

	// 控制上下文, 控制挂起或解卦
	ctrlCtx    context.Context
	ctrlCancel context.CancelFunc

	stg       storage.ITask
	kind      string
	taskName  string
	workspace string
	scriptDir string
	dagTasks  map[string]dagcuter.Task
	state     int32 // 0: 正常, 1: 挂起
}

func newTask(taskName string, workerSpace, scriptDir string) (*sTask, error) {
	t := &sTask{
		stg:       storage.Task(taskName),
		taskName:  taskName,
		dagTasks:  make(map[string]dagcuter.Task),
		workspace: filepath.Join(workerSpace, taskName),
		scriptDir: filepath.Join(scriptDir, taskName),
	}
	t.lcCtx, t.lcCancel = context.WithCancel(context.WithValue(context.Background(), "ctx", "task"))
	var err error
	defer func() {
		if err != nil {
			state := models.StateFailed
			if t.stg.IsDisable() {
				state = models.StateStopped
			}
			_ = t.stg.Update(&models.STaskUpdate{
				Message:  err.Error(),
				State:    models.Pointer(state),
				OldState: models.Pointer(models.StatePending),
				STime:    models.Pointer(time.Now()),
				ETime:    models.Pointer(time.Now()),
			})
			// 清理资源
			t.clearDir()
		}
	}()

	// 如果任务被禁用，则更新所有步骤状态并返回错误
	if t.stg.IsDisable() {
		logx.Infoln("the task is disabled, no execution required", taskName)
		for _, sName := range t.stg.StepNameList("") {
			_ = t.stg.Step(sName).Update(&models.SStepUpdate{
				Message:  "the task is disabled, no execution required",
				State:    models.Pointer(models.StateStopped),
				OldState: models.Pointer(models.StatePending),
				STime:    models.Pointer(time.Now()),
				ETime:    models.Pointer(time.Now()),
			})
		}
		err = errors.New("the task is disabled, no execution required")
		return nil, err
	}

	// 获取任务类型
	t.kind, err = t.stg.Kind()
	if err != nil {
		logx.Errorln(t.taskName, err)
		return nil, err
	}

	for _, s := range t.stg.StepList("") {
		if t.stg.Step(s.Name).IsDisable() {
			logx.Infoln("the step is disabled, no execution required", s.Name)
			_ = t.stg.Step(s.Name).Update(&models.SStepUpdate{
				Message:  "the step is disabled, no execution required",
				State:    models.Pointer(models.StateStopped),
				OldState: models.Pointer(models.StatePending),
				STime:    models.Pointer(time.Now()),
				ETime:    models.Pointer(time.Now()),
			})
			continue
		}
		t.dagTasks[s.Name] = t.newStep(s.Name)
	}
	if dagcuter.HasCycle(t.dagTasks) {
		err = errors.New("the task has a cycle")
		return nil, err
	}
	// 加入管理
	taskManager.Store(t.taskName, t)
	return t, nil
}

func (t *sTask) updateNode(node string) error {
	return t.stg.UpdateNode(node)
}

func (t *sTask) Execute() (err error) {
	defer func() {
		// 清理资源
		t.Stop()
	}()
	// 更新任务状态为运行中
	if err = t.stg.Update(&models.STaskUpdate{
		State:    models.Pointer(models.StateRunning),
		OldState: models.Pointer(models.StatePending),
		STime:    models.Pointer(time.Now()),
		Message:  "task is running",
	}); err != nil {
		logx.Errorln(t.taskName, err)
		return
	}

	if err = t.checkCtx(); err != nil {
		logx.Errorln(t.taskName, err)
		return
	}

	if err = t.initDir(); err != nil {
		logx.Errorln(t.taskName, err)
		return
	}
	res := new(models.STaskUpdate)
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			logx.Errorln(r, string(stack))
			err = errors.Errorf("task panic: %v", r)
		}
		res.State = models.Pointer(models.StateStopped)
		res.Message = "task has stopped"
		res.ETime = models.Pointer(time.Now())
		res.OldState = models.Pointer(models.StateRunning)
		if err != nil {
			res.State = models.Pointer(models.StateFailed)
			res.Message = err.Error()
		}
		if updErr := t.stg.Update(res); updErr != nil {
			logx.Warnln(t.taskName, updErr)
		}
	}()

	timeout, err := t.stg.Timeout()
	if err != nil {
		logx.Errorln(t.taskName, err)
		return
	}
	var ctx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeoutCause(t.lcCtx, timeout+time.Minute, common.ErrTimeOut)
	} else {
		ctx, cancel = context.WithCancel(t.lcCtx)
	}
	defer cancel()

	_dag, err := dagcuter.New(t.dagTasks)
	if err != nil {
		logx.Errorln(t.taskName, err)
		return
	}
	_, err = _dag.Execute(ctx)
	if err != nil {
		logx.Errorln(t.taskName, err)
		return
	}
	// 策略模式下，获取最后一个非待执行状态的步骤状态
	if t.kind == common.KindStrategy {
		steps := t.stg.StepList("")
		for i := len(steps) - 1; i >= 0; i-- {
			switch *steps[i].State {
			case models.StateFailed:
				err = errors.New(steps[i].Message)
				return
			case models.StateStopped:
				return
			default:
				continue
			}
		}
		logx.Warnln(t.taskName, "no steps executed????")
	}
	return
}

func (t *sTask) checkCtx() error {
	// 挂起, 则等待解挂
	if t.ctrlCtx != nil {
		// 等待控制信号
		select {
		case <-t.ctrlCtx.Done():
		case <-t.lcCtx.Done():
			return t.lcCtx.Err()
		}
	}

	if t.lcCtx.Err() != nil {
		return t.lcCtx.Err()
	}
	return nil
}

func (t *sTask) Stop() {
	defer func() {
		if _err := recover(); _err != nil {
			logx.Errorln(_err)
		}
	}()
	if t.lcCancel != nil {
		logx.Infoln(t.taskName, "Stop")
		event.SendEventf("%s Stop", t.taskName)
		t.lcCancel()
	}
	// 删除manager
	taskManager.Delete(t.taskName)
	for _, step := range t.dagTasks {
		stepManager.Delete(step.Name())
	}
	t.clearDir()
}

func (t *sTask) initDir() error {
	if err := utils.EnsureDirExist(t.workspace); err != nil {
		logx.Errorln(t.taskName, t.workspace, t.scriptDir, err)
		return err
	}
	if err := utils.EnsureDirExist(t.scriptDir); err != nil {
		logx.Errorln(t.taskName, err)
		return err
	}
	return nil
}

func (t *sTask) clearDir() {
	if err := os.RemoveAll(t.scriptDir); err != nil {
		logx.Errorln(t.taskName, err)
	}
	if err := os.RemoveAll(t.workspace); err != nil {
		logx.Errorln(t.taskName, err)
	}
}

func (t *sTask) newStep(stepName string) *sStep {
	s := &sStep{
		kind:      t.kind,
		taskName:  t.taskName,
		stepName:  stepName,
		stg:       t.stg.Step(stepName),
		workspace: t.workspace,
		scriptDir: t.scriptDir,
	}
	s.lcCtx, s.lcCancel = context.WithCancel(context.WithValue(context.Background(), "ctx", "step"))
	stepManager.Store(s.Name(), s)
	return s
}
