package locker

import "context"

type ILocker interface {
	// NewLock creates a new unlocked lock object for the given upload ID.
	NewLock(id string) (ILock, error)
}

type ILock interface {
	Lock(ctx context.Context) error
	Unlock()
}
