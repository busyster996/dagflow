package worker

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/busyster996/dagflow/internal/common"
	"github.com/busyster996/dagflow/internal/pubsub"
	"github.com/busyster996/dagflow/internal/runner"
	"github.com/busyster996/dagflow/internal/storage"
	"github.com/busyster996/dagflow/internal/storage/models"
	"github.com/busyster996/dagflow/internal/utility"
	"github.com/busyster996/dagflow/internal/worker/event"
	"github.com/busyster996/dagflow/pkg/logx"
	"github.com/busyster996/dagflow/pkg/tunny"
)

var (
	pool        = tunny.NewCallback(1)
	taskManager = new(sync.Map)
	stepManager = new(sync.Map)
)

func Start(ctx context.Context) error {
	logx.Infoln("number of workers", GetSize())
	if err := storage.FixDatabase(viper.GetString("node_name")); err != nil {
		return err
	}

	// 清理当前节点的残留文件
	for _, t := range storage.NodeTasks(viper.GetString("node_name")) {
		// clear old script
		utility.ClearDir(filepath.Join(viper.GetString("script_dir"), t.Name()))

		// clear old workspace
		utility.ClearDir(filepath.Join(viper.GetString("workspace_dir"), t.Name()))
	}

	// 打印当前支持的runner
	logx.Infoln("runner", runner.ListAvailable())
	if err := pubsub.SubscribeTask(ctx, viper.GetString("node_name"), func(data string) {
		if data == "" {
			return
		}
		t, err := newTask(data)
		if err != nil {
			logx.Errorln(err)
			return
		}
		err = pool.Submit(t.Execute)
		if err != nil {
			logx.Errorln(err)
		}
	}); err != nil {
		return err
	}
	if err := pubsub.SubscribeTask(ctx, "random", func(data string) {
		if data == "" {
			return
		}
		t, err := newTask(data)
		if err != nil {
			logx.Errorln(err)
			return
		}
		if err = t.updateNode(viper.GetString("node_name")); err != nil {
			logx.Errorln(err)
			return
		}
		err = pool.Submit(t.Execute)
		if err != nil {
			logx.Errorln(err)
		}
	}); err != nil {
		return err
	}

	if err := pubsub.SubscribeManager(ctx, viper.GetString("node_name"), func(data string) {
		if !utility.ContainsInvisibleChar(data) {
			return
		}
		slice := utility.SplitByInvisibleChar(data)
		switch len(slice) {
		case 3:
			taskName := slice[0]
			action := slice[1]
			duration := slice[2]
			if err := managerTask(taskName, action, duration); err != nil {
				logx.Errorln(err)
			}
		case 4:
			taskName := slice[0]
			stepName := slice[1]
			action := slice[2]
			duration := slice[3]
			if err := managerStep(taskName, stepName, action, duration); err != nil {
				logx.Errorln(err)
			}
		}
	}); err != nil {
		return err
	}
	_event, id, err := event.Subscribe()
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				event.Unsubscribe(id)
				return
			case e := <-_event:
				_ = pubsub.PublishEvent(e)
			}
		}
	}()
	return nil
}

func managerTask(taskName, action, duration string) error {
	t, err := storage.Task(taskName).Get()
	if err != nil {
		return err
	}

	value, ok := taskManager.Load(taskName)
	if !ok {
		return errors.New("task not found")
	}
	task, ok := value.(*sTask)
	switch action {
	case "kill":
		task.Stop()
		return storage.Task(taskName).Update(&models.STaskUpdate{
			State:    models.Pointer(models.StateFailed),
			OldState: t.State,
			Message:  "has been killed",
		})
	case "pause":
		if *t.State == models.StateRunning {
			return errors.New("step is running")
		}
		if atomic.CompareAndSwapInt32(&task.state, 0, 1) {
			var d time.Duration
			d, err = time.ParseDuration(duration)
			if err == nil && d > 0 {
				task.ctrlCtx, task.ctrlCancel = context.WithTimeout(context.Background(), d)
			} else {
				task.ctrlCtx, task.ctrlCancel = context.WithCancel(context.Background())
			}
			return storage.Task(taskName).Update(&models.STaskUpdate{
				State:    models.Pointer(models.StatePaused),
				OldState: t.State,
				Message:  "has been paused",
			})
		}
	case "resume":
		if atomic.CompareAndSwapInt32(&task.state, 1, 0) {
			if task.ctrlCancel != nil {
				task.ctrlCancel()
			}
			return storage.Task(taskName).Update(&models.STaskUpdate{
				State:    t.OldState,
				OldState: t.State,
				Message:  "has been resumed",
			})
		}
	}
	return nil
}

func managerStep(taskName, stepName, action, duration string) error {
	value, ok := stepManager.Load(fmt.Sprintf("%s/%s", taskName, stepName))
	if !ok {
		return errors.New("step not found")
	}
	step, ok := value.(*sStep)
	if !ok {
		return errors.New("step not found")
	}
	s, err := storage.Task(taskName).Step(stepName).Get()
	if err != nil {
		return err
	}
	switch action {
	case "kill":
		step.Stop()
		return storage.Task(taskName).Step(stepName).Update(&models.SStepUpdate{
			Code:     models.Pointer(common.ExecCodeKilled),
			State:    models.Pointer(models.StateFailed),
			OldState: s.State,
			Message:  "has been killed",
		})
	case "pause":
		if *s.State == models.StateRunning {
			return errors.New("step is running")
		}
		if atomic.CompareAndSwapInt32(&step.state, 0, 1) {
			var d time.Duration
			d, err = time.ParseDuration(duration)
			if err == nil && d > 0 {
				step.ctrlCtx, step.ctrlCancel = context.WithTimeout(context.Background(), d)
			} else {
				step.ctrlCtx, step.ctrlCancel = context.WithCancel(context.Background())
			}
			return storage.Task(taskName).Step(stepName).Update(&models.SStepUpdate{
				State:    models.Pointer(models.StatePaused),
				OldState: s.State,
				Message:  "has been paused",
			})
		}
	case "resume":
		if atomic.CompareAndSwapInt32(&step.state, 1, 0) {
			if step.ctrlCancel != nil {
				step.ctrlCancel()
			}
			return storage.Task(taskName).Step(stepName).Update(&models.SStepUpdate{
				State:    s.OldState,
				OldState: s.State,
				Message:  "has been resumed",
			})
		}
	}
	return nil
}

func SetSize(n int) {
	pool.SetSize(n)
}

func GetSize() int {
	return pool.GetSize()
}

func Shutdown() {
	pool.Close()
}
