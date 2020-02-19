package locker

import (
	"context"
	"sync"
	"time"

	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
)

type (
	// Memory TryLocker for non clustered applications
	Memory struct {
		mainLocker sync.Locker
		locks      map[string]struct{}
	}
)

var _ placeorder.TryLocker = &Memory{}

// NewMemory creates a new memory based lock
func NewMemory() *Memory {
	return &Memory{mainLocker: &sync.Mutex{}, locks: make(map[string]struct{})}
}

func (m *Memory) locked(id string) (ok bool) { _, ok = m.locks[id]; return }

// TryLock unblocking implementation see https://github.com/LK4D4/trylock/blob/master/trylock.go
func (m *Memory) TryLock(_ context.Context, key string, _ time.Duration) (placeorder.Unlock, error) {
	m.mainLocker.Lock()
	defer m.mainLocker.Unlock()
	if m.locked(key) {
		return nil, placeorder.ErrLockTaken
	}
	m.locks[key] = struct{}{}
	return func() error {
		m.mainLocker.Lock()
		defer m.mainLocker.Unlock()
		delete(m.locks, key)
		return nil
	}, nil
}
