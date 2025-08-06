package storage

import (
	"time"

	"github.com/busyster996/dagflow/internal/storage/models"
)

const All = ""

type IStorage interface {
	Name() (name string)

	// FixDatabase fix database
	FixDatabase(node string) (err error)

	// NodeTasks 节点任务接口
	NodeTasks(node string) (tasks []ITask)
	// Task 任务接口
	Task(name string) (task ITask)
	// TaskCreate 创建任务
	TaskCreate(task *models.STask) (err error)
	// TaskCount 指定状态任务总数, -1为所有
	TaskCount(state models.State) (res int64)
	// TaskList 获取任务,支持分页, 模糊匹配
	TaskList(page, pageSize int64, str string) (res models.STasks, total int64)

	// Pipeline 流水线接口
	Pipeline(name string) (pipeline IPipeline)
	// PipelineCreate 创建流水线
	PipelineCreate(pipeline *models.SPipeline) (err error)
	// PipelineList 获取流水线,支持分页, 模糊匹配
	PipelineList(page, pageSize int64, str string) (res models.SPipelines, total int64)
}

type IBase interface {
	// Name 名称
	Name() (name string)
	// ClearAll 清理
	ClearAll() (err error)
	// Remove 删除
	Remove() (err error)
}

type ITask interface {
	IBase

	// Kind 获取类型
	Kind() (res string, err error)
	// IsDisable 是否禁用
	IsDisable() (disable bool)
	// State 获取状态
	State() (state models.State, err error)
	// Env 环境变量接口
	Env() (env IEnv)

	// Timeout 超时时间
	Timeout() (res time.Duration, err error)
	// Get 根据名称获取指定任务
	Get() (res *models.STask, err error)
	// Update 更新
	Update(value *models.STaskUpdate) (err error)
	// UpdateNode 更新节点
	UpdateNode(node string) error

	// Step 步骤接口
	Step(name string) IStep
	// StepCreate 创建步骤
	StepCreate(step *models.SStep) (err error)
	// StepCount 任务总数
	StepCount() (res int64)
	// StepNameList 所有步骤名称
	StepNameList(str string) (res []string)
	// StepStateList 获取任务下所有步骤状态
	StepStateList(str string) (res map[string]models.State)
	// StepList 获取任务下所有步骤
	StepList(str string) (res models.SSteps)
}

type IStep interface {
	IBase

	// IsDisable 是否禁用
	IsDisable() (disable bool)
	// State 获取状态
	State() (state models.State, err error)
	// Env 环境变量接口
	Env() (env IEnv)

	// TaskName 任务名称
	TaskName() (taskName string)
	// Timeout 超时时间
	Timeout() (res time.Duration, err error)
	// Type 类型
	Type() (res string, err error)
	// Content 内容
	Content() (res string, err error)
	// Action 动作
	Action() (res string, err error)
	// Rule 规则
	Rule() (res string, err error)
	// RetryPolicy 重试策略
	RetryPolicy() (res models.SRetryPolicy, err error)
	// Get 根据名称获取指定步骤
	Get() (res *models.SStep, err error)
	// Update 更新
	Update(value *models.SStepUpdate) (err error)
	// GlobalEnv 全局环境变量接口
	GlobalEnv() (env IEnv)
	// Depend 依赖接口
	Depend() (depend IDepend)
	// Log 日志接口
	Log() (log ILog)
}

type ILog interface {
	// List 获取指定任务指定步骤所有日志, 增量查询
	List(latestLine *int64) (res models.SStepLogs)
	// Insert 插入
	Insert(log *models.SStepLog) (err error)
	Write(contents ...string)
	Writef(format string, args ...interface{})
	RemoveAll() (err error)
}

type IEnv interface {
	List() (res models.SEnvs)
	Insert(env ...*models.SEnv) (err error)
	Update(env *models.SEnv) (err error)
	Get(name string) (string, error)
	Remove(name string) (err error)
	RemoveAll() (err error)
}

type IDepend interface {
	List() (res []string)
	Insert(depends ...string) (err error)
	RemoveAll() (err error)
}

type IPipeline interface {
	IBase

	// Build 执行相关
	Build() (build IPipelineBuild)
	// Task 任务接口
	Task(name string) (task ITask)

	// Update 更新
	Update(value *models.SPipelineUpdate) (err error)
	// Get 获取
	Get() (res *models.SPipeline, err error)
	// IsDisable 是否禁用
	IsDisable() (disable bool)
	// TplType 类型
	TplType() (res string, err error)
	// Content 内容
	Content() (res string, err error)
}

type IPipelineBuild interface {
	// Get 根据名称获取指定构建
	Get(name string) (res *models.SPipelineBuildRes, err error)
	// Insert 插入
	Insert(build *models.SPipelineBuild) (err error)
	// List 获取所有
	List(page, size int64) (res models.SPipelineBuilds, total int64)
	// Remove 删除
	Remove(name string) (err error)
	// ClearAll 清理
	ClearAll() (err error)
}
