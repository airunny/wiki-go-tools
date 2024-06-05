package locker

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAlreadyLocked = errors.New("already locked")
	ErrLockTimeout   = errors.New("lock timeout")
)

type Release func() error

type Locker interface {
	Lock(ctx context.Context, key string, expires, timeout time.Duration) (Release, error)
	TryLock(ctx context.Context, key string, expires time.Duration) (Release, error)
}
