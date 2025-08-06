package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"github.com/spf13/viper"
	"go.uber.org/multierr"

	"github.com/busyster996/dagflow/internal/pubsub"
	"github.com/busyster996/dagflow/internal/server/types"
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/internal/utils"
	"github.com/busyster996/dagflow/internal/worker/common"
	"github.com/busyster996/dagflow/pkg/logx"
)

// 只允许中文,英文(含大小写),0-9,-_.~字符
var reg = regexp.MustCompile("[^a-zA-Z\\p{Han}0-9\\-_.~]")

type STaskService struct {
	name string
}

func Task(name string) *STaskService {
	return &STaskService{
		name: name,
	}
}

func TaskList(req *types.SPageReq) *types.STaskListDetailRes {
	tasks, total := storage.TaskList(req.Page, req.Size, req.Prefix)
	if tasks == nil {
		return nil
	}
	pageTotal := total / req.Size
	if total%req.Size != 0 {
		pageTotal += 1
	}
	var list = &types.STaskListDetailRes{
		Page: &types.SPageRes{
			Current: req.Page,
			Size:    req.Size,
			Total:   pageTotal,
		},
	}
	for _, task := range tasks {
		res := &types.STaskRes{
			Kind:    task.Kind,
			Name:    task.Name,
			State:   models.StateMap[*task.State],
			Message: task.Message,
			Time: &types.STimeRes{
				Start: task.STimeStr(),
				End:   task.ETimeStr(),
			},
		}
		st := storage.Task(task.Name)

		// 获取当前进行到那些步骤
		steps := st.StepStateList(storage.All)
		res.Count = int64(len(steps))
		var groups = make(map[models.State][]string)
		for name, state := range steps {
			groups[state] = append(groups[state], name)
		}
		res.Message = GenerateStateMessage(res.Message, groups)
		list.Tasks = append(list.Tasks, res)
	}
	return list
}

func (ts *STaskService) Create(task *types.STaskReq) (err error) {
	task.Kind = strings.ToLower(task.Kind)
	// 检查请求内容
	if err = ts.review(task); err != nil {
		logx.Errorln("task review", ts.name, err)
		return err
	}

	var db = storage.Task(task.Name)
	// 检查全局
	state, err := db.State()
	if err != nil {
		logx.Errorln("task state", ts.name, err)
		return err
	}
	if state != models.StateStopped && state != models.StateSkipped && state != models.StateUnknown && state != models.StateFailed {
		return errors.New("task is running")
	}

	// 清理旧数据
	_ = db.ClearAll()

	defer func() {
		if err != nil {
			// rollback
			_ = db.ClearAll()
		}
	}()

	if err = ts.saveTask(task); err != nil {
		logx.Errorln("task save", ts.name, err)
		return err
	}

	err = ts.reviewStep(task.Kind, task.Step)
	if err != nil {
		logx.Errorln("task review step", ts.name, err)
		return err
	}

	for k, step := range task.Step {
		var seqNo = int64(k + 1)
		// save step
		stepSvc := Step(task.Name, step.Name)
		if _err := stepSvc.Create(seqNo, task.Timeout, step); _err != nil {
			logx.Errorln("task create step", task.Name, step.Name, _err)
			err = multierr.Append(err, fmt.Errorf("save step error: %s", _err))
		}
	}
	if err != nil {
		logx.Errorln("task create", ts.name, err)
		return err
	}
	node := task.Node
	if node == "" {
		node = "random"
	}
	// 提交任务
	delayedSub := task.Delayed.Sub(time.Now().UTC())
	if !task.Delayed.IsZero() && delayedSub > 0 {
		return pubsub.PublishTaskDelayed(node, ts.name, delayedSub)
	}
	return pubsub.PublishTask(node, ts.name)
}

func (ts *STaskService) review(task *types.STaskReq) error {
	if task.Step == nil || len(task.Step) == 0 {
		return errors.New("steps can not be empty")
	}

	// 校验env是否重复
	var envKeys []string
	for _, v := range task.Env {
		envKeys = append(envKeys, v.Name)
	}
	dup := utils.CheckDuplicate(envKeys)
	if dup != nil {
		return fmt.Errorf("duplicate keys %v", dup)
	}

	task.Name = reg.ReplaceAllString(task.Name, "")
	if task.Name == "" {
		task.Name = ksuid.New().String()
	}
	ts.name = task.Name
	// 确保超时时间在合理范围内
	if task.Timeout <= 0 || task.Timeout >= viper.GetDuration("exec_timeout") {
		task.Timeout = viper.GetDuration("exec_timeout")
	}
	return nil
}

func (ts *STaskService) reviewStep(kind string, steps types.SStepsReq) error {
	// 检查步骤名称是否重复
	if err := ts.uniqStepsName(steps); err != nil {
		logx.Errorln("task review step", ts.name, err)
		return err
	}
	if kind != common.KindDag {
		// 非编排模式,按顺序执行
		for k := range steps {
			if k == 0 {
				steps[k].Depends = nil
				continue
			}
			steps[k].Depends = []string{steps[k-1].Name}
		}
	}
	return nil
}

func (ts *STaskService) uniqStepsName(steps types.SStepsReq) error {
	counts := make(map[string]int)
	for _, v := range steps {
		if v.Name == "" {
			continue
		}
		counts[v.Name]++
	}
	var errs []error
	for name, count := range counts {
		if count > 1 {
			errs = append(errs, fmt.Errorf("%s repeat count %d", name, count))
		}
	}
	if errs == nil {
		return nil
	}
	return fmt.Errorf("%v", errs)
}
func (ts *STaskService) saveTask(task *types.STaskReq) error {
	// save task
	err := storage.TaskCreate(&models.STask{
		Kind:    task.Kind,
		Name:    task.Name,
		Desc:    task.Desc,
		Node:    task.Node,
		Timeout: task.Timeout,
		Disable: models.Pointer(task.Disable),
		STaskUpdate: models.STaskUpdate{
			Message:  "the task is waiting to be scheduled for execution",
			State:    models.Pointer(models.StatePending),
			OldState: models.Pointer(models.StatePending),
		},
	})
	if err != nil {
		logx.Errorln("task save", ts.name, err)
		return fmt.Errorf("save task error: %s", err)
	}

	// save task env
	var envs models.SEnvs
	for _, env := range task.Env {
		envs = append(envs, &models.SEnv{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	if err = storage.Task(task.Name).Env().Insert(envs...); err != nil {
		logx.Errorln("task save envs", ts.name, err)
		return fmt.Errorf("save task env error: %s", err)
	}
	return nil
}

func (ts *STaskService) Delete() error {
	task, err := storage.Task(ts.name).Get()
	if err != nil {
		logx.Errorln("task delete", ts.name, err)
		return errors.New("task not found")
	}
	// 尝试强杀任务
	err = pubsub.PublishManager(task.Node, utils.JoinWithInvisibleChar(ts.name, "kill", "0"))
	if err != nil {
		logx.Errorln("task delete", ts.name, "kill error", err)
		return err
	}
	return storage.Task(ts.name).ClearAll()
}

func (ts *STaskService) Exist() bool {
	_, err := storage.Task(ts.name).Get()
	return err == nil
}

func (ts *STaskService) Detail() (types.Code, *types.STaskRes, error) {
	db := storage.Task(ts.name)
	var task *models.STask
	var err error
	task, err = db.Get()
	if err != nil {
		logx.Errorln("task detail", ts.name, err)
		return types.CodeFailed, nil, errors.New("task not found")
	}

	data := &types.STaskRes{
		Kind:    task.Kind,
		Name:    task.Name,
		Desc:    task.Desc,
		Node:    task.Node,
		State:   models.StateMap[*task.State],
		Message: task.Message,
		Timeout: task.Timeout,
		Disable: *task.Disable,
		Time: &types.STimeRes{
			Start: task.STimeStr(),
			End:   task.ETimeStr(),
		},
	}
	for _, env := range db.Env().List() {
		data.Env = append(data.Env, &types.SEnv{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	// 获取当前进行到那些步骤
	steps := db.StepStateList(storage.All)
	data.Count = int64(len(steps))
	var groups = make(map[models.State][]string)
	for name, state := range steps {
		groups[state] = append(groups[state], name)
	}
	data.Message = GenerateStateMessage(data.Message, groups)
	return ConvertState(*task.State), data, nil
}

func (ts *STaskService) StepCount() (res int64) {
	return storage.Task(ts.name).StepCount()
}

func (ts *STaskService) Manager(action string, duration string) error {
	task, err := storage.Task(ts.name).Get()
	if err != nil {
		logx.Errorln("task manager", ts.name, err)
		return errors.New("task not found")
	}
	if *task.State != models.StateRunning && *task.State != models.StatePending && *task.State != models.StatePaused {
		return errors.New("task is no running")
	}
	return pubsub.PublishManager(task.Node, utils.JoinWithInvisibleChar(ts.name, action, duration))
}

func (ts *STaskService) Dump() (*types.STaskReq, error) {
	task, err := storage.Task(ts.name).Get()
	if err != nil {
		logx.Errorln("task dump", ts.name, err)
		return nil, errors.New("task not found")
	}
	res := &types.STaskReq{
		Kind:    task.Kind,
		Name:    task.Name,
		Desc:    task.Desc,
		Node:    task.Node,
		Timeout: task.Timeout,
		Disable: *task.Disable,
	}
	for _, env := range storage.Task(ts.name).Env().List() {
		res.Env = append(res.Env, &types.SEnv{
			Name:  env.Name,
			Value: env.Value,
		})
	}
	steps := storage.Task(ts.name).StepList(storage.All)
	for _, step := range steps {
		stepRes := &types.SStepReq{
			Name:    step.Name,
			Desc:    step.Desc,
			Type:    step.Type,
			Content: step.Content,
			Timeout: step.Timeout,
			Disable: *step.Disable,
			Action:  step.Action,
			Rule:    step.Rule,
			RetryPolicy: &types.SRetryPolicy{
				Interval:    step.RetryPolicy.Data().Interval,
				MaxAttempts: step.RetryPolicy.Data().MaxAttempts,
				MaxInterval: step.RetryPolicy.Data().MaxInterval,
				Multiplier:  step.RetryPolicy.Data().Multiplier,
			},
		}
		envs := storage.Task(ts.name).Step(step.Name).Env().List()
		for _, env := range envs {
			stepRes.Env = append(stepRes.Env, &types.SEnv{
				Name:  env.Name,
				Value: env.Value,
			})
		}
		stepRes.Depends = storage.Task(ts.name).Step(step.Name).Depend().List()
		res.Step = append(res.Step, stepRes)
	}
	return res, nil
}

func (ts *STaskService) Steps() (code types.Code, data types.SStepsRes, err error) {
	db := storage.Task(ts.name)
	task, err := db.Get()
	if err != nil {
		logx.Errorln("task get steps", ts.name, err)
		return types.CodeNoData, nil, err
	}

	steps := db.StepList(storage.All)
	if steps == nil {
		return types.CodeNoData, nil, errors.New("steps not found")
	}

	// 用于分组和构建任务数据
	var groups = make(map[models.State][]string)
	taskMap := make(map[string]*types.SStepRes, len(steps))
	for _, step := range steps {
		groups[*step.State] = append(groups[*step.State], step.Name)
		taskMap[step.Name] = &types.SStepRes{
			Name:    step.Name,
			State:   models.StateMap[*step.State],
			Code:    *step.Code,
			Message: step.Message,
			Time: &types.STimeRes{
				Start: step.STimeStr(),
				End:   step.ETimeStr(),
			},
			Depends: db.Step(step.Name).Depend().List(),
		}
	}

	// 按深度排序
	data = ts.sortTasksByDepth(taskMap)

	// 生成任务的状态消息
	task.Message = GenerateStateMessage(task.Message, groups)

	return ConvertState(*task.State), data, errors.New(task.Message)
}

// 按深度排序
func (ts *STaskService) sortTasksByDepth(taskMap map[string]*types.SStepRes) types.SStepsRes {
	visited := make(map[string]bool)
	sorted := make([]*types.SStepRes, 0, len(taskMap))

	var visit func(name string)
	visit = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true

		if task, exists := taskMap[name]; exists {
			for _, dep := range task.Depends {
				visit(dep) // 递归访问依赖
			}
			sorted = append(sorted, task)
		}
	}

	for name := range taskMap {
		visit(name) // 遍历所有任务
	}

	return sorted
}
