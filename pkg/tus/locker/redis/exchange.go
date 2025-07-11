package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type LockExchange struct {
	client redis.UniversalClient
}

func (e *LockExchange) Listen(ctx context.Context, id string) {
	psub := e.client.PSubscribe(ctx, LockExchangeChannel+id)
	defer psub.Close()
	c := psub.Channel()
	select {
	case <-c:
		return
	case <-ctx.Done():
		return
	}
}

func (e *LockExchange) Request(ctx context.Context, id string) error {
	psub := e.client.PSubscribe(ctx, LockReleaseChannel+id)
	defer psub.Close()
	res := e.client.Publish(ctx, LockExchangeChannel+id, id)
	if res.Err() != nil {
		return res.Err()
	}
	select {
	case <-psub.Channel():
		return nil
	case <-ctx.Done():
		return fmt.Errorf("timeout")
	}
}

func (e *LockExchange) Release(ctx context.Context, id string) error {
	res := e.client.Publish(ctx, LockReleaseChannel+id, id)
	return res.Err()
}
