package locker

import (
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type (
	//Simple TryLock for nin clustered applications
	Simple struct {
		m sync.Mutex
	}
)

var _ placeorder.TryLock = &Simple{}

const mutexLocked = 1 << iota

//TryLock unblocking implementation - see https://github.com/LK4D4/trylock/blob/master/trylock.go
func (s *Simple) TryLock(key string, maxlockduration time.Duration) (placeorder.Unlock, error) {
	haveLock := atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&s.m)), 0, mutexLocked)
	if !haveLock {
		return nil, placeorder.ErrLockTaken
	}
	return func() error {
		s.m.Unlock()
		return nil
	}, nil
}
