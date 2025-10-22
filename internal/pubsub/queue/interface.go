package queue

import (
	"context"
	"time"

	"github.com/busyster996/dagflow/internal/utility"
)

type IBroker interface {
	PublishEvent(data string) error
	PublishTask(node string, data string) error
	PublishTaskDelayed(node string, data string, delay time.Duration) error
	PublishManager(node string, data string) error

	SubscribeEvent(ctx context.Context, handler utility.QueueHandleFn) error
	SubscribeTask(ctx context.Context, node string, handler utility.QueueHandleFn) error
	SubscribeManager(ctx context.Context, node string, handler utility.QueueHandleFn) error

	Shutdown(ctx context.Context)
}

func TaskRoutingKey() string {
	return utility.ServiceName + ".task"
}

func EventRoutingKey() string {
	return utility.ServiceName + ".event"
}

func ManagerRoutingKey() string {
	return utility.ServiceName + ".manager"
}
