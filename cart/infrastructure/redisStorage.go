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
	"github.com/go-redis/redis/v8"
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
		ttl int
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
	_ CartStorage    = &RedisStorage{}
	_ CartSerializer = &GobSerializer{}

	errCartIsNil = errors.New("cart is nil")
)

// Inject dependencies and build redis client
func (r *RedisStorage) Inject(
	serializer CartSerializer,
	config *struct {
		RedisKeyPrefix       string  `inject:"config:commerce.cart.redis.keyPrefix"`
		RedisTTL             int     `inject:"config:commerce.cart.redis.ttl"`
		RedisNetwork         string  `inject:"config:commerce.cart.redis.network"`
		RedisAddress         string  `inject:"config:commerce.cart.redis.address"`
		RedisPassword        string  `inject:"config:commerce.cart.redis.password"`
		RedisIdleConnections float64 `inject:"config:commerce.cart.redis.idle.connections"`
		RedisDatabase        int     `inject:"config:commerce.cart.redis.database,optional"`
		RedisTLS             bool    `inject:"config:commerce.cart.redis.tls,optional"`
		RedisClusterMode     bool    `inject:"config:commerce.cart.redis.clusterMode,optional"`
	},
) *RedisStorage {
	r.serializer = serializer
	if config != nil {
		r.keyPrefix = config.RedisKeyPrefix
		r.ttl = config.RedisTTL

		var tlsConfig *tls.Config
		if config.RedisTLS {
			tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}

		if config.RedisClusterMode {
			r.client = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:     []string{config.RedisAddress},
				Password:  config.RedisPassword,
				PoolSize:  int(config.RedisIdleConnections),
				TLSConfig: tlsConfig,
			})
		} else {
			r.client = redis.NewClient(&redis.Options{
				Network:   config.RedisNetwork,
				Addr:      config.RedisAddress,
				Password:  config.RedisPassword,
				DB:        config.RedisDatabase,
				PoolSize:  int(config.RedisIdleConnections),
				TLSConfig: tlsConfig,
			})
		}

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
		return errCartIsNil
	}

	b, err := r.serializer.Serialize(cart)
	if err != nil {
		return fmt.Errorf("could not store cart: %w", err)
	}

	err = r.client.Set(ctx, r.keyPrefix+cart.ID, b, time.Duration(r.ttl)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("could not store cart: %w", err)
	}

	return nil
}

// RemoveCart deletes a cart from redis
func (r *RedisStorage) RemoveCart(ctx context.Context, cart *cartDomain.Cart) error {
	if cart == nil {
		return errCartIsNil
	}

	err := r.client.Del(ctx, r.keyPrefix+cart.ID).Err()
	if err != nil {
		return fmt.Errorf("could not remove cart: %w", err)
	}

	return nil
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
