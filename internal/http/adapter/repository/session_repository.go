package repository

import (
	"context"
	"examples/internal/http/interface/infra"
	"time"
)

type sessionRepository struct {
	store infra.LocalStore
}

func NewSessionRepository(localStore infra.LocalStore) *sessionRepository {
	return &sessionRepository{store: localStore}
}

func (m *sessionRepository) Get(ctx context.Context, token string) (int, bool) {
	v, ok := m.store.Get(ctx, token)
	if !ok {
		return 0, false
	}
	userID, ok := v.(int)
	if !ok {
		return 0, false
	}
	return userID, true
}

func (m *sessionRepository) Set(ctx context.Context, token string, userID int, ttl time.Duration) {
	m.store.Set(ctx, token, userID, ttl)
}

func (m *sessionRepository) Drop(ctx context.Context, token string) {
	m.store.Drop(ctx, token)
}
