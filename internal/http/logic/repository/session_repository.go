package repository

import (
	"context"
	"time"
)

type SessionRepository interface {
	Get(ctx context.Context, token string) (int, bool)
	Set(ctx context.Context, token string, userID int, ttl time.Duration)
	Drop(ctx context.Context, token string)
}
