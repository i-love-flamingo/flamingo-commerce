package contextstore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"runtime"
	"time"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/gomodule/redigo/redis"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Redis saves all contexts in a simple map
	Redis struct {
		pool   *redis.Pool
		logger flamingo.Logger
	}
)

var _ process.ContextStore = new(Redis)
var _ healthcheck.Status = &Redis{}

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
	}) *Redis {
	r.logger = logger
	if cfg != nil {
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
func (r *Redis) Store(key string, value process.Context) error {
	conn := r.pool.Get()
	defer conn.Close()

	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(value)
	if err != nil {
		return err
	}
	_, err = conn.Do(
		"SET",
		key,
		buffer,
	)

	return err
}

// Get a stored context
func (r *Redis) Get(key string) (process.Context, bool) {
	conn := r.pool.Get()
	defer conn.Close()

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
func (r *Redis) Delete(key string) error {
	conn := r.pool.Get()
	defer conn.Close()

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
