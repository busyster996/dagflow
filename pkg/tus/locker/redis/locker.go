package redis

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"

	"github.com/busyster996/dagflow/pkg/tus/locker"
)

var (
	LockExchangeChannel = "tusd_lock_release_request_"
	LockReleaseChannel  = "tusd_lock_released_"
	LockExpiry          = 8 * time.Second
)

type ILockExchange interface {
	Listen(ctx context.Context, id string)
	Request(ctx context.Context, id string) error
}

type IBidirectionalLockExchange interface {
	ILockExchange
	Release(ctx context.Context, id string) error
}

type IMutexLock interface {
	TryLockContext(context.Context) error
	ExtendContext(context.Context) (bool, error)
	UnlockContext(context.Context) (bool, error)
	Until() time.Time
}

type Locker struct {
	LockExchangeChannel string
	LockReleaseChannel  string
	CreateMutex         func(id string) IMutexLock
	Exchange            IBidirectionalLockExchange
}

func (locker *Locker) NewLock(id string) (locker.ILock, error) {
	mutex := locker.CreateMutex(id)
	return &redisLock{
		id:       id,
		mutex:    mutex,
		exchange: locker.Exchange,
	}, nil
}

type redisLock struct {
	id       string
	mutex    IMutexLock
	ctx      context.Context
	cancel   func()
	exchange IBidirectionalLockExchange
}

func (l *redisLock) Lock(ctx context.Context) error {
	if err := l.requestLock(ctx); err != nil {
		return err
	}
	go l.exchange.Listen(l.ctx, l.id)
	go func() {
		if err := l.keepAlive(l.ctx); err != nil {
			l.cancel()
		}
	}()
	return nil
}

func (l *redisLock) aquireLock(ctx context.Context) error {
	if err := l.mutex.TryLockContext(ctx); err != nil {
		// Currently there aren't any errors
		// defined by redsync we don't want to retry.
		// If there are any return just that error without
		// handler.ErrFileLocked to show it's non-recoverable.
		return err
	}

	l.ctx, l.cancel = context.WithCancel(context.Background())

	return nil
}

func (l *redisLock) requestLock(ctx context.Context) error {
	err := l.aquireLock(ctx)
	if err == nil {
		return nil
	}
	if err = l.exchange.Request(ctx, l.id); err != nil {
		return err
	}
	return l.aquireLock(ctx)
}

func (l *redisLock) keepAlive(ctx context.Context) error {
	//insures that an extend will be canceled if it's unlocked in the middle of an attempt
	for {
		select {
		case <-time.After(time.Until(l.mutex.Until()) / 2):
			_, err := l.mutex.ExtendContext(ctx)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (l *redisLock) Unlock() {
	if l.cancel != nil {
		defer l.cancel()
	}
	_, _ = l.mutex.UnlockContext(l.ctx)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_ = l.exchange.Release(ctx, l.id)
	return
}

func NewFromClient(client redis.UniversalClient) (*Locker, error) {
	rs := redsync.New(goredis.NewPool(client))

	_locker := &Locker{
		CreateMutex: func(id string) IMutexLock {
			return rs.NewMutex(id, redsync.WithExpiry(LockExpiry))
		},
		Exchange: &LockExchange{
			client: client,
		},
	}

	return _locker, nil
}

func New(uri string) (*Locker, error) {
	connection, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(connection)
	if res := client.Ping(context.Background()); res.Err() != nil {
		return nil, res.Err()
	}
	return NewFromClient(client)
}
