package locker

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
)

type (
	// Memory TryLocker for non clustered applications
	Memory struct {
		m sync.Mutex
	}
)

var _ placeorder.TryLocker = &Memory{}

const mutexLocked = 1 << iota

// TryLock unblocking implementation see https://github.com/LK4D4/trylock/blob/master/trylock.go
func (s *Memory) TryLock(key string, _ time.Duration) (placeorder.Unlock, error) {
	// Todo: support multiple locks based on lock key
	haveLock := atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&s.m)), 0, mutexLocked)
	if !haveLock {
		return nil, placeorder.ErrLockTaken
	}
	return func() error {
		s.m.Unlock()
		return nil
	}, nil
}
