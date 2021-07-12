package contextstore

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"runtime"
	"time"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/gomodule/redigo/redis"
	"go.opencensus.io/trace"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Redis saves all contexts in a simple map
	Redis struct {
		pool   *redis.Pool
		logger flamingo.Logger
		ttl    time.Duration
	}
)

var (
	_ process.ContextStore = new(Redis)
	_ healthcheck.Status   = &Redis{}
	// ErrNoRedisConnection is returned if the underlying connection is erroneous
	ErrNoRedisConnection = errors.New("no redis connection, see healthcheck")
)

func init() {
	gob.Register(process.Context{})
}

// Inject dependencies
func (r *Redis) Inject(
	logger flamingo.Logger,
	cfg *struct {
		MaxIdle                 int    `inject:"config:commerce.checkout.placeorder.contextstore.redis.maxIdle"`
		IdleTimeoutMilliseconds int    `inject:"config:commerce.checkout.placeorder.contextstore.redis.idleTimeoutMilliseconds"`
		Network                 string `inject:"config:commerce.checkout.placeorder.contextstore.redis.network"`
		Address                 string `inject:"config:commerce.checkout.placeorder.contextstore.redis.address"`
		Database                int    `inject:"config:commerce.checkout.placeorder.contextstore.redis.database"`
		TTL                     string `inject:"config:commerce.checkout.placeorder.contextstore.redis.ttl"`
	}) *Redis {
	r.logger = logger
	if cfg != nil {
		var err error
		r.ttl, err = time.ParseDuration(cfg.TTL)
		if err != nil {
			panic("can't parse commerce.checkout.placeorder.contextstore.redis.ttl")
		}

		r.pool = &redis.Pool{
			MaxIdle:     cfg.MaxIdle,
			IdleTimeout: time.Duration(cfg.IdleTimeoutMilliseconds) * time.Millisecond,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Dial: func() (redis.Conn, error) {
				return redis.Dial(cfg.Network, cfg.Address, redis.DialDatabase(cfg.Database))
			},
		}
		runtime.SetFinalizer(r, func(r *Redis) { r.pool.Close() }) // close all connections on destruction
	}

	return r
}

// Store a given context
func (r *Redis) Store(ctx context.Context, key string, placeOrderContext process.Context) error {
	_, span := trace.StartSpan(ctx, "placeorder/contextstore/Store")
	defer span.End()
	conn := r.pool.Get()
	defer conn.Close()
	if conn.Err() != nil {
		r.logger.Error("placeorder/contextstore/Store:", conn.Err())
		return ErrNoRedisConnection
	}

	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(placeOrderContext)
	if err != nil {
		return err
	}
	_, err = conn.Do(
		"SETEX",
		key,
		int(r.ttl.Round(time.Second).Seconds()),
		buffer,
	)

	return err
}

// Get a stored context
func (r *Redis) Get(ctx context.Context, key string) (process.Context, bool) {
	_, span := trace.StartSpan(ctx, "placeorder/contextstore/Get")
	defer span.End()
	conn := r.pool.Get()
	defer conn.Close()
	if conn.Err() != nil {
		r.logger.Error("placeorder/contextstore/Get:", conn.Err())
	}

	content, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return process.Context{}, false
	}

	buffer := bytes.NewBuffer(content)
	decoder := gob.NewDecoder(buffer)
	pctx := new(process.Context)
	err = decoder.Decode(pctx)
	if err != nil {
		r.logger.Error(fmt.Sprintf("context in key %q is not decodable: %s", key, err))
	}

	return *pctx, err == nil
}

// Delete a stored context, nop if it doesn't exist
func (r *Redis) Delete(ctx context.Context, key string) error {
	_, span := trace.StartSpan(ctx, "placeorder/contextstore/Delete")
	defer span.End()
	conn := r.pool.Get()
	defer conn.Close()
	if conn.Err() != nil {
		r.logger.Error("placeorder/contextstore/Delete:", conn.Err())
		return ErrNoRedisConnection
	}

	_, err := conn.Do("DEL", key)

	return err
}

// Status handles the health check of redis
func (r *Redis) Status() (alive bool, details string) {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	if err == nil {
		return true, "redis for place order context store replies to PING"
	}

	return false, err.Error()
}
