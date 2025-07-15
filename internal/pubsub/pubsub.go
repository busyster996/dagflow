package pubsub

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/busyster996/dagflow/internal/pubsub/common"
	"github.com/busyster996/dagflow/internal/pubsub/queue"
	"github.com/busyster996/dagflow/pkg/logx"
)

var (
	broker queue.IBroker
)

func New(rawURL string) error {
	// 打印当前支持的队列
	logx.Infoln("queue", queue.ListAvailable())
	scheme, _, found := strings.Cut(rawURL, "://")
	if !found {
		return fmt.Errorf("invalid message queue url")
	}

	factory, err := queue.Get(scheme)
	if err != nil {
		return err
	}

	broker, err = factory(rawURL)
	if err != nil {
		return errors.Wrap(err, "failed to create broker")
	}
	return nil
}

func PublishTask(node string, data string) error {
	return broker.PublishTask(node, data)
}

func PublishTaskDelayed(node string, data string, delay time.Duration) error {
	return broker.PublishTaskDelayed(node, data, delay)
}

func SubscribeTask(ctx context.Context, node string, handler common.HandleFn) error {
	return broker.SubscribeTask(ctx, node, handler)
}

func PublishEvent(data string) error {
	return broker.PublishEvent(data)
}
func SubscribeEvent(ctx context.Context, handler common.HandleFn) error {
	return broker.SubscribeEvent(ctx, handler)
}
func PublishManager(node string, data string) error {
	return broker.PublishManager(node, data)
}
func SubscribeManager(ctx context.Context, node string, handler common.HandleFn) error {
	return broker.SubscribeManager(ctx, node, handler)
}

func Shutdown(ctx context.Context) {
	broker.Shutdown(ctx)
}
