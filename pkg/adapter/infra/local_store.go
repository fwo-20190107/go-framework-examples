package infra

import (
	"context"
	"time"
)

type LocalStore interface {
	Get(ctx context.Context, key any) (value any, ok bool)
	Set(ctx context.Context, key, value any, ttl time.Duration)
	Drop(ctx context.Context, key any)
}
