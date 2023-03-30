package cache

import (
	"context"
	"examples/pkg/adapter/infra"
	"sync"
	"time"
)

type localStore struct {
	store sync.Map
}

type storeValue struct {
	ch    chan struct{}
	value any
	once  sync.Once
}

type timer struct {
	t   *time.Timer
	key any
}

func NewLocalStore() *localStore {
	return &localStore{}
}

func (s *localStore) Set(ctx context.Context, key, value any, ttl time.Duration) {
	if _, ok := s.getsv(key); ok {
		return
	}

	sv := storeValue{
		value: value,
	}
	if ttl > 0 {
		sv.ch = make(chan struct{})
		go func(ch chan struct{}) {
			timer := &timer{
				t:   time.NewTimer(ttl),
				key: key,
			}

			select {
			case <-timer.t.C:
			case <-ch:
				// nothing todo
			}

			s.delete(timer.key)
			if !timer.t.Stop() {
				<-timer.t.C
			}
		}(sv.ch)
	}
	s.store.Store(key, sv)
}

func (s *localStore) Get(ctx context.Context, key any) (any, bool) {
	sv, ok := s.getsv(key)
	if !ok {
		return nil, ok
	}
	return sv.value, ok
}

func (s *localStore) getsv(key any) (*storeValue, bool) {
	v, ok := s.store.Load(key)
	if !ok {
		return nil, ok
	}
	sv, ok := v.(storeValue)
	if !ok {
		return nil, ok
	}
	return &sv, ok
}

func (s *localStore) Drop(ctx context.Context, key any) {
	sv, ok := s.getsv(key)
	if !ok {
		return
	}
	sv.close()
}

func (s *localStore) delete(key any) {
	sv, ok := s.getsv(key)
	if !ok {
		return
	}
	sv.close()
	s.store.Delete(key)
}

func (v *storeValue) close() {
	v.once.Do(func() { close(v.ch) })
}

var _ infra.LocalStore = (*localStore)(nil)
