package dagcuter

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/busyster996/dagflow/pkg/tunny"
)

type Dagcuter struct {
	Tasks          map[string]Task
	results        *sync.Map
	inDegrees      map[string]int
	dependents     map[string][]string
	executionOrder []string
	worker         *tunny.Pool
	mu             *sync.Mutex
	wg             *sync.WaitGroup
}

func New(tasks map[string]Task) (*Dagcuter, error) {
	return NewWithWorkers(tasks, 150) // 默认150个工作者
}

func NewWithWorkers(tasks map[string]Task, maxWorkers int) (*Dagcuter, error) {
	if HasCycle(tasks) {
		return nil, fmt.Errorf("circular dependency detected")
	}
	dag := &Dagcuter{
		mu:         new(sync.Mutex),
		wg:         new(sync.WaitGroup),
		results:    new(sync.Map),
		inDegrees:  make(map[string]int),
		dependents: make(map[string][]string),
		Tasks:      tasks,
	}
	if maxWorkers <= 0 {
		// 默认150个工作者
		maxWorkers = 150
	} else if maxWorkers > 1000 {
		// 最大不超过1000个工作者
		maxWorkers = 1000
	}
	dag.worker = tunny.NewCallback(maxWorkers)

	for name, task := range dag.Tasks {
		dag.inDegrees[name] = len(task.Dependencies())
		for _, dep := range task.Dependencies() {
			dag.dependents[dep] = append(dag.dependents[dep], name)
		}
	}

	return dag, nil
}

func (d *Dagcuter) WorkerStatus() (queueLength int64) {
	if d.worker == nil {
		return 0
	}
	return d.worker.QueueLength()
}

func (d *Dagcuter) ResizeWorker(newSize int) {
	// 如果新容量小于等于当前容量且最大不能超过1000，则不需要调整
	if d.worker != nil && d.worker.GetSize() < newSize && newSize <= 1000 {
		d.worker.SetSize(newSize)
	}
}

func (d *Dagcuter) Execute(ctx context.Context) (map[string]map[string]any, error) {
	defer d.results.Clear()
	defer d.worker.Close()
	errCh := make(chan error, 1)

	for name, deg := range d.inDegrees {
		if deg == 0 {
			d.wg.Add(1)
			go d.runTask(ctx, name, errCh)
		}
	}

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		results := make(map[string]map[string]any)
		d.results.Range(func(key, value any) bool {
			results[key.(string)] = value.(map[string]any)
			return true
		})
		return results, nil
	case err := <-errCh:
		return nil, err
	}
}

func (d *Dagcuter) runTask(ctx context.Context, name string, errCh chan error) {
	defer d.wg.Done()
	task := d.Tasks[name]

	d.mu.Lock()
	inputs := d.prepareInputs(task)
	d.mu.Unlock()

	output, err := d.executeTask(ctx, name, task, inputs)
	if err != nil {
		select {
		case errCh <- err:
		default:
		}
		return
	}

	d.mu.Lock()
	d.executionOrder = append(d.executionOrder, name)
	d.mu.Unlock()

	d.mu.Lock()
	d.results.Store(name, output)
	for _, child := range d.dependents[name] {
		d.inDegrees[child]--
		if d.inDegrees[child] == 0 {
			d.wg.Add(1)
			//go d.runTask(ctx, child, errCh)

			// 使用 tunny 池来执行任务
			_ = d.worker.Submit(func() error {
				d.runTask(ctx, child, errCh)
				return nil
			})
		}
	}
	d.mu.Unlock()
}

func (d *Dagcuter) executeTask(ctx context.Context, name string, task Task, inputs map[string]any) (map[string]any, error) {
	// 获取任务的重试策略
	retryExecutor := d.newRetryExecutor(task.RetryPolicy())

	var result map[string]any

	// 使用重试机制执行任务
	err := retryExecutor.ExecuteWithRetry(ctx, name, func(n int) error {
		inputs["attempt"] = n // 将当前尝试次数传递给任务
		// PreExecution
		if err := task.PreExecution(ctx, inputs); err != nil {
			return fmt.Errorf("pre execution failed: %w", err)
		}

		// Execute
		output, err := task.Execute(ctx, inputs)
		if err != nil {
			return fmt.Errorf("execution failed: %w", err)
		}

		// PostExecution
		if err = task.PostExecution(ctx, output); err != nil {
			return fmt.Errorf("post execution failed: %w", err)
		}

		result = output
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("task %s failed: %w", name, err)
	}

	return result, nil
}

func (d *Dagcuter) prepareInputs(task Task) map[string]any {
	inputs := make(map[string]any)
	for _, dep := range task.Dependencies() {
		if value, ok := d.results.Load(dep); ok {
			inputs[dep] = value
		}
	}
	return inputs
}

func (d *Dagcuter) ExecutionOrder() string {
	var sb = strings.Builder{}
	sb.WriteString("\n")
	for i, step := range d.executionOrder {
		_, _ = fmt.Fprintf(&sb, "%d. %s\n", i+1, step)
	}
	return sb.String()
}

// PrintGraph 输出链式依赖
func (d *Dagcuter) PrintGraph() {
	// 1. 找到所有根节点（入度为 0）
	var roots []string
	for name, deg := range d.inDegrees {
		if deg == 0 {
			roots = append(roots, name)
		}
	}
	// 2. 分别从每个根节点开始打印
	for _, root := range roots {
		fmt.Println(root)        // 先打印根
		d.printChain(root, "  ") // 从根的下一层开始缩进两格
		fmt.Println()            // 不同根之间空行
	}
}

// printChain 递归打印子依赖，
// name: 当前节点；
// prefix: 当前缩进前缀（已经包含了箭头前需要的空格）
func (d *Dagcuter) printChain(name, prefix string) {
	children := d.dependents[name]
	for _, child := range children {
		// 打印箭头和子节点
		fmt.Printf("%s└─> %s\n", prefix, child)
		// 递归打印子节点的子依赖，缩进再多四格
		d.printChain(child, prefix+"    ")
	}
}
