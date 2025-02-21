package persistence

import (
	"context"
	"time"
)

type Store interface {
	Incr(ctx context.Context, key string, expiration time.Duration) (int, error)
}
