package queue

import (
	"context"
	"time"

	"github.com/busyster996/dagflow/internal/pubsub/common"
	"github.com/busyster996/dagflow/internal/utils"
)

type IBroker interface {
	PublishEvent(data string) error
	PublishTask(node string, data string) error
	PublishTaskDelayed(node string, data string, delay time.Duration) error
	PublishManager(node string, data string) error

	SubscribeEvent(ctx context.Context, handler common.HandleFn) error
	SubscribeTask(ctx context.Context, node string, handler common.HandleFn) error
	SubscribeManager(ctx context.Context, node string, handler common.HandleFn) error

	Shutdown(ctx context.Context)
}

func TaskRoutingKey() string {
	return utils.ServiceName + ".task"
}

func EventRoutingKey() string {
	return utils.ServiceName + ".event"
}

func ManagerRoutingKey() string {
	return utils.ServiceName + ".manager"
}
