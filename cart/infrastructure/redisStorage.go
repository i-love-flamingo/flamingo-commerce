package infrastructure

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/gob"
	"errors"
	"fmt"
	"runtime"
	"time"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"github.com/go-redis/redis/v9"
)

type (
	// RedisStorage stores carts in redis
	RedisStorage struct {
		// client to connect to redis
		client redis.UniversalClient
		// cart serializer
		serializer CartSerializer
		// key prefix with which the cart will be stored
		keyPrefix string
		// time to live
		ttlGuest    time.Duration
		ttlCustomer time.Duration
	}

	// CartSerializer serializes carts in order to store them in redis
	CartSerializer interface {
		Serialize(c *cartDomain.Cart) ([]byte, error)
		Deserialize(b []byte) (*cartDomain.Cart, error)
	}

	// GobSerializer serializes carts using gob
	GobSerializer struct{}
)

var (
	_ CartStorage        = &RedisStorage{}
	_ healthcheck.Status = &RedisStorage{}
	_ CartSerializer     = &GobSerializer{}

	ErrCartIsNil = errors.New("cart is nil")
)

// Inject dependencies and build redis client
func (r *RedisStorage) Inject(
	serializer CartSerializer,
	config *struct {
		RedisKeyPrefix       string  `inject:"config:commerce.cart.redis.keyPrefix"`
		RedisTTLGuest        string  `inject:"config:commerce.cart.redis.ttl.guest"`
		RedisTTLCustomer     string  `inject:"config:commerce.cart.redis.ttl.customer"`
		RedisNetwork         string  `inject:"config:commerce.cart.redis.network"`
		RedisAddress         string  `inject:"config:commerce.cart.redis.address"`
		RedisPassword        string  `inject:"config:commerce.cart.redis.password"`
		RedisIdleConnections float64 `inject:"config:commerce.cart.redis.idleConnections"`
		RedisDatabase        int     `inject:"config:commerce.cart.redis.database,optional"`
		RedisTLS             bool    `inject:"config:commerce.cart.redis.tls,optional"`
	},
) *RedisStorage {
	r.serializer = serializer
	if config != nil {
		r.keyPrefix = config.RedisKeyPrefix

		var err error
		r.ttlGuest, err = time.ParseDuration(config.RedisTTLGuest)
		if err != nil {
			panic("can't parse commerce.cart.redis.ttl.guest")
		}

		r.ttlCustomer, err = time.ParseDuration(config.RedisTTLCustomer)
		if err != nil {
			panic("can't parse commerce.cart.redis.ttl.customer")
		}

		var tlsConfig *tls.Config
		if config.RedisTLS {
			tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}

		r.client = redis.NewClient(&redis.Options{
			Network:   config.RedisNetwork,
			Addr:      config.RedisAddress,
			Password:  config.RedisPassword,
			DB:        config.RedisDatabase,
			PoolSize:  int(config.RedisIdleConnections),
			TLSConfig: tlsConfig,
		})

		// close redis client
		runtime.SetFinalizer(r, func(r *RedisStorage) { _ = r.client.Close() })
	}

	return r
}

// GetCart fetches a cart from redis and deserializes it
func (r *RedisStorage) GetCart(ctx context.Context, id string) (*cartDomain.Cart, error) {
	cmd := r.client.Get(ctx, r.keyPrefix+id)
	if cmd.Err() != nil {
		return nil, fmt.Errorf("could not get cart: %w", cmd.Err())
	}

	b, err := cmd.Bytes()
	if err != nil {
		return nil, fmt.Errorf("could not get cart: %w", err)
	}

	cart, err := r.serializer.Deserialize(b)
	if err != nil {
		return nil, fmt.Errorf("could not get cart: %w", err)
	}

	return cart, nil
}

// HasCart checks if the cart id exists as a key in redis
func (r *RedisStorage) HasCart(ctx context.Context, id string) bool {
	return r.client.Exists(ctx, r.keyPrefix+id).Val() > 0
}

// StoreCart serializes a cart and stores it in redis
func (r *RedisStorage) StoreCart(ctx context.Context, cart *cartDomain.Cart) error {
	if cart == nil {
		return ErrCartIsNil
	}

	b, err := r.serializer.Serialize(cart)
	if err != nil {
		return fmt.Errorf("could not store cart: %w", err)
	}

	err = r.client.Set(ctx, r.keyPrefix+cart.ID, b, r.ttl(cart)).Err()
	if err != nil {
		return fmt.Errorf("could not store cart: %w", err)
	}

	return nil
}

// ttl may differ for guest and customer carts
func (r *RedisStorage) ttl(cart *cartDomain.Cart) time.Duration {
	if cart.BelongsToAuthenticatedUser {
		return r.ttlCustomer
	}

	return r.ttlGuest
}

// RemoveCart deletes a cart from redis
func (r *RedisStorage) RemoveCart(ctx context.Context, cart *cartDomain.Cart) error {
	if cart == nil {
		return ErrCartIsNil
	}

	err := r.client.Del(ctx, r.keyPrefix+cart.ID).Err()
	if err != nil {
		return fmt.Errorf("could not remove cart: %w", err)
	}

	return nil
}

// Status healthcheck via ping
func (r *RedisStorage) Status() (alive bool, details string) {
	err := r.client.Ping(context.Background()).Err()
	if err != nil {
		return false, err.Error()
	}

	return true, "redis for cart storage replies to PING"
}

// Serialize a cart using gob
func (gs GobSerializer) Serialize(c *cartDomain.Cart) ([]byte, error) {
	buf := new(bytes.Buffer)

	err := gob.NewEncoder(buf).Encode(&c)
	if err != nil {
		return nil, fmt.Errorf("could not serialize cart: %w", err)
	}

	return buf.Bytes(), nil
}

// Deserialize a cart using gob
func (gs GobSerializer) Deserialize(d []byte) (*cartDomain.Cart, error) {
	var cart cartDomain.Cart

	err := gob.NewDecoder(bytes.NewBuffer(d)).Decode(&cart)
	if err != nil {
		return nil, fmt.Errorf("could not deserialize cart: %w", err)
	}

	return &cart, nil
}
