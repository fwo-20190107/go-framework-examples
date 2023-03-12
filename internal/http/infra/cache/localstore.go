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

type storeValue struct {
	ch    chan struct{}
	value any
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
	if sv.ch != nil {
		sv.ch <- struct{}{}
	} else {
		s.delete(key)
	}
}

func (s *localStore) delete(key any) {
	sv, ok := s.getsv(key)
	if !ok {
		return
	}
	if sv.ch != nil {
		close(sv.ch)
	}
	s.store.Delete(key)
}

var _ repository.LocalStore = (*localStore)(nil)
