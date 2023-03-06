package cache

import (
	"context"
	"examples/internal/http/interface/repository"
	"sync"
	"time"
)

type localStore struct {
	store sync.Map
}

type timer struct {
	t   *time.Timer
	key any
}

func NewLocalStore() *localStore {
	return &localStore{}
}

func (s *localStore) Set(ctx context.Context, key, value any, ttl time.Duration) {
	if _, ok := s.Get(ctx, key); ok {
		return
	}

	if ttl > 0 {
		timer := &timer{
			t:   time.NewTimer(ttl),
			key: key,
		}
		go func() {
			select {
			case <-ctx.Done():
			case <-timer.t.C:
				// nothing todo
			}

			s.Drop(ctx, timer.key)
			if !timer.t.Stop() {
				<-timer.t.C
			}
		}()
	}
	s.store.Store(key, value)
}

func (s *localStore) Get(ctx context.Context, key any) (any, bool) {
	return s.store.Load(key)
}

func (s *localStore) Drop(ctx context.Context, key any) {
	s.store.Delete(key)
}

var _ repository.LocalStore = (*localStore)(nil)
