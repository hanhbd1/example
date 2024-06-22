package lock

import (
	"context"
	"errors"
	"time"
)

var (
	ErrLockObtained = errors.New("already obtained")
	ErrLockNotHeld  = errors.New("lock not held")
)

type Lock interface {
	Release(ctx context.Context) error
}

type Locker interface {
	Obtain(ctx context.Context, key string, value string, obtainTime ...time.Duration) (Lock, error)
}
