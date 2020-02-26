package locker

import (
	"context"
	"errors"
	"time"

	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"go.opencensus.io/trace"
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
		healthcheck func() error
	}
)

var _ placeorder.TryLocker = &Redis{}
var _ healthcheck.Status = &Redis{}

// NewRedis creates a new distributed mutex using multiple Redis connection pools.
func NewRedis(
	cfg *struct {
		MaxIdle                 int    `inject:"config:commerce.checkout.placeorder.lock.redis.maxIdle"`
		IdleTimeoutMilliseconds int    `inject:"config:commerce.checkout.placeorder.lock.redis.idleTimeoutMilliseconds"`
		Network                 string `inject:"config:commerce.checkout.placeorder.lock.redis.network"`
		Address                 string `inject:"config:commerce.checkout.placeorder.lock.redis.address"`
		Database                int    `inject:"config:commerce.checkout.placeorder.lock.redis.database"`
	},
) *Redis {
	r := new(Redis)

	if cfg != nil {
		r.maxIdle = cfg.MaxIdle
		r.idleTimeout = time.Duration(cfg.IdleTimeoutMilliseconds) * time.Millisecond
		r.network = cfg.Network
		r.address = cfg.Address
		r.database = cfg.Database
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

	r.healthcheck = func() error {
		_, err := pools[0].Get().Do("PING")
		return err
	}

	r.redsync = redsync.New(pools)

	return r
}

// TryLock ties once to acquire a lock and returns the unlock func if successful
func (r *Redis) TryLock(ctx context.Context, key string, maxlockduration time.Duration) (placeorder.Unlock, error) {
	_, span := trace.StartSpan(ctx, "placeorder/lock/TryLock")
	defer span.End()
	mutex := r.redsync.NewMutex(
		key,
		redsync.SetExpiry(maxlockduration),
		redsync.SetTries(1),
		redsync.SetRetryDelayFunc(func(int) time.Duration { return 50 * time.Millisecond }),
	)
	err := mutex.Lock()
	if err != nil {
		alive, _ := r.Status()
		if !alive {
			return nil, errors.New("redis not reachable, see health-check")
		}
		return nil, placeorder.ErrLockTaken
	}
	ticker := time.NewTicker(maxlockduration / 3)
	go func() {
		for range ticker.C {
			mutex.Extend()
		}
	}()

	return func() error {
		ticker.Stop()
		ok := mutex.Unlock()
		if !ok {
			return errors.New("unlock unsuccessful")
		}
		return nil
	}, nil
}

// Status is the health check
func (r *Redis) Status() (alive bool, details string) {
	err := r.healthcheck()

	if err == nil {
		return true, "redis for place order lock replies to PING"
	}

	return false, err.Error()
}
