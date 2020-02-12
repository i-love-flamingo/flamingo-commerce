package locker

import (
	"errors"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"

	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
)

type (
	// Redis TryLocker for clustered applications
	Redis struct {
		redsync     *redsync.Redsync
		network     string
		address     string
		maxIdle     int
		idleTimeout time.Duration
		database    int
	}
)

var _ placeorder.TryLocker = &Redis{}

// NewRedis creates a new distributed mutex using multiple Redis connection pools.
func NewRedis(
	cfg *struct {
		MaxIdle                 float64 `inject:"config:commerce.checkout.placeorder.lock.redis.maxIdle"`
		IdleTimeoutMilliseconds float64 `inject:"config:commerce.checkout.placeorder.lock.redis.idleTimeoutMilliseconds"`
		Network                 string  `inject:"config:commerce.checkout.placeorder.lock.redis.network"`
		Address                 string  `inject:"config:commerce.checkout.placeorder.lock.redis.address"`
		Database                float64 `inject:"config:commerce.checkout.placeorder.lock.redis.database"`
	},
) *Redis {
	r := new(Redis)

	if cfg != nil {
		r.maxIdle = int(cfg.MaxIdle)
		r.idleTimeout = time.Duration(cfg.IdleTimeoutMilliseconds) * time.Millisecond
		r.network = cfg.Network
		r.address = cfg.Address
		r.database = int(cfg.Database)
	}

	pools := []redsync.Pool{&redis.Pool{
		MaxIdle:     r.maxIdle,
		IdleTimeout: r.idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(r.network, r.address, redis.DialDatabase(r.database))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}}

	r.redsync = redsync.New(pools)

	return r
}

// TryLock ties once to acquire a lock and returns the unlock func if successful
func (r *Redis) TryLock(key string, maxlockduration time.Duration) (placeorder.Unlock, error) {
	mutex := r.redsync.NewMutex(
		key,
		redsync.SetExpiry(maxlockduration),
		redsync.SetTries(1),
		redsync.SetRetryDelayFunc(func(int) time.Duration { return 50 * time.Millisecond }),
	)
	err := mutex.Lock()
	if err != nil {
		return nil, err
	}
	ticker := time.NewTicker(maxlockduration / 3)
	go func() {
		for {
			<-ticker.C
			mutex.Extend()
		}
	}()

	return func() error {
		ok := mutex.Unlock()
		if !ok {
			return errors.New("unlock unsuccessful")
		}
		ticker.Stop()
		return nil
	}, nil
}
