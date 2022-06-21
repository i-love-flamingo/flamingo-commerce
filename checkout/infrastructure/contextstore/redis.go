package contextstore

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"runtime"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/go-redis/redis/v8"
	"go.opencensus.io/trace"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Redis saves all contexts in a simple map
	Redis struct {
		client *redis.Client
		logger flamingo.Logger
		ttl    time.Duration
	}
)

var (
	_ process.ContextStore = new(Redis)
	_ healthcheck.Status   = &Redis{}
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

		r.client = redis.NewClient(&redis.Options{
			PoolSize:    cfg.MaxIdle,
			IdleTimeout: time.Duration(cfg.IdleTimeoutMilliseconds) * time.Millisecond,
			DB:          cfg.Database,
			Addr:        cfg.Address,
			Network:     cfg.Network,
		})

		runtime.SetFinalizer(r, func(r *Redis) { _ = r.client.Close() }) // close all connections on destruction
	}

	return r
}

// Store a given context
func (r *Redis) Store(ctx context.Context, key string, placeOrderContext process.Context) error {
	_, span := trace.StartSpan(ctx, "placeorder/contextstore/Store")
	defer span.End()

	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(placeOrderContext)
	if err != nil {
		return err
	}

	return r.client.SetEX(context.Background(), key, buffer.String(), r.ttl).Err()
}

// Get a stored context
func (r *Redis) Get(ctx context.Context, key string) (process.Context, bool) {
	_, span := trace.StartSpan(ctx, "placeorder/contextstore/Get")
	defer span.End()

	content, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return process.Context{}, false
	}

	decoder := gob.NewDecoder(strings.NewReader(content))
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

	return r.client.Del(context.Background(), key).Err()
}

// Status handles the health check of redis
func (r *Redis) Status() (alive bool, details string) {
	err := r.client.Ping(context.Background()).Err()
	if err == nil {
		return true, "redis for place order context store replies to PING"
	}

	return false, err.Error()
}
