package clustersync

import (
	"errors"
	"time"
)

type (
	TryLock interface {
		//TryLock
		//  - requests a (distributed) lock - idendified by "key"
		//  - if ythe lock could be get it returns true
		//  - Unlocking:
		//  	- The Lock is automatically released if the process is killed
		//		- Or Unlock is called
		//	- if somethings fails during Lock or Unlock the lock will be available again after brokenLockTimeout
		TryLock(key string, brokenLockTimeout time.Duration) (Unlock, error)
	}

	Unlock func() error
)

var (
	LockTaken = errors.New("Lock is already taken")
)
