package infrastructure_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os/exec"
	"testing"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	"github.com/go-redis/redis/v8"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stvp/tempredis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	existingKey  = "test"
	wrongDataKey = "wrong data"
)

var (
	testCart = &cartDomain.Cart{
		ID:       "test",
		EntityID: "1",
	}
)

func getRedisStorage(network, address string) *infrastructure.RedisStorage {
	return new(infrastructure.RedisStorage).Inject(
		&infrastructure.GobSerializer{},
		&struct {
			RedisKeyPrefix       string  `inject:"config:commerce.cart.redis.keyPrefix"`
			RedisTTL             int     `inject:"config:commerce.cart.redis.ttl"`
			RedisNetwork         string  `inject:"config:commerce.cart.redis.network"`
			RedisAddress         string  `inject:"config:commerce.cart.redis.address"`
			RedisPassword        string  `inject:"config:commerce.cart.redis.password"`
			RedisIdleConnections float64 `inject:"config:commerce.cart.redis.idle.connections"`
			RedisDatabase        int     `inject:"config:commerce.cart.redis.database,optional"`
			RedisTLS             bool    `inject:"config:commerce.cart.redis.tls,optional"`
			RedisClusterMode     bool    `inject:"config:commerce.cart.redis.clusterMode,optional"`
		}{RedisIdleConnections: 3, RedisNetwork: network, RedisAddress: address, RedisDatabase: 0, RedisTTL: 60})
}

func prepareData(t *testing.T, ctx context.Context, client redis.UniversalClient) {
	t.Helper()

	buffer := new(bytes.Buffer)
	require.NoError(t, gob.NewEncoder(buffer).Encode(&testCart))
	err := client.Set(ctx, existingKey, buffer.Bytes(), 0).Err()
	require.NoError(t, err)
	err = client.Set(ctx, wrongDataKey, "wrong data", 0).Err()
	require.NoError(t, err)
}

func startUpLocalRedis(t *testing.T) (*tempredis.Server, redis.UniversalClient) {
	t.Helper()

	server, err := tempredis.Start(tempredis.Config{})
	if err != nil {
		t.Fatal(err)
	}

	client := redis.NewClient(&redis.Options{Network: "unix", Addr: server.Socket()})

	prepareData(t, context.Background(), client)

	return server, client
}

func startUpDockerRedis(t *testing.T) (func(), string, redis.UniversalClient) {
	t.Helper()

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	port, err := redisC.MappedPort(ctx, "6379")
	require.NoError(t, err)

	host, err := redisC.Host(ctx)
	require.NoError(t, err)

	address := fmt.Sprintf("%s:%s", host, port.Port())
	client := redis.NewClient(&redis.Options{Network: "tcp", Addr: address})

	prepareData(t, ctx, client)

	return func() { _ = redisC.Terminate(ctx) }, address, client
}

func TestRedisStorage_GetCart(t *testing.T) {
	t.Parallel()

	runTestCases := func(t *testing.T, storage *infrastructure.RedisStorage) {
		t.Helper()

		tests := []struct {
			name        string
			key         string
			expected    *cartDomain.Cart
			expectedErr bool
		}{
			{
				name:        "load existing",
				key:         existingKey,
				expected:    testCart,
				expectedErr: false,
			},
			{
				name:        "load existing with wrong data",
				key:         wrongDataKey,
				expected:    nil,
				expectedErr: true,
			},
			{
				name:        "load non existing",
				key:         "non",
				expected:    nil,
				expectedErr: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := storage.GetCart(context.Background(), tt.key)
				assert.Equal(t, tt.expectedErr, err != nil)
				if diff := deep.Equal(got, tt.expected); diff != nil {
					t.Error("expected response is wrong: ", diff)
				}
			})
		}
	}

	t.Run("local-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("redis-server"); err != nil {
			t.Skip("redis-server not installed")
		}

		server, _ := startUpLocalRedis(t)
		store := getRedisStorage("unix", server.Socket())
		runTestCases(t, store)
	})

	t.Run("docker-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("docker"); err != nil {
			t.Skip("docker not installed")
		}

		shutdown, address, _ := startUpDockerRedis(t)
		defer shutdown()
		store := getRedisStorage("tcp", address)
		runTestCases(t, store)
	})
}

func TestRedisStorage_HasCart(t *testing.T) {
	t.Parallel()

	runTestCases := func(t *testing.T, storage *infrastructure.RedisStorage) {
		t.Helper()

		tests := []struct {
			name     string
			key      string
			expected bool
		}{
			{
				name:     "has existing",
				key:      existingKey,
				expected: true,
			},
			{
				name:     "has not non existing",
				key:      "non",
				expected: false,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := storage.HasCart(context.Background(), tt.key)
				assert.Equal(t, tt.expected, got)
			})
		}
	}

	t.Run("local-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("redis-server"); err != nil {
			t.Skip("redis-server not installed")
		}

		server, _ := startUpLocalRedis(t)
		store := getRedisStorage("unix", server.Socket())
		runTestCases(t, store)
	})

	t.Run("docker-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("docker"); err != nil {
			t.Skip("docker not installed")
		}

		shutdown, address, _ := startUpDockerRedis(t)
		defer shutdown()
		store := getRedisStorage("tcp", address)
		runTestCases(t, store)
	})
}

func TestRedisStorage_StoreCart(t *testing.T) {
	t.Parallel()

	runTestCases := func(t *testing.T, storage *infrastructure.RedisStorage, client redis.UniversalClient) {
		t.Helper()

		tests := []struct {
			name  string
			key   string
			value *cartDomain.Cart
		}{
			{
				name:  "store new value",
				key:   "another-test",
				value: &cartDomain.Cart{ID: "another-test", EntityID: "1"},
			},
			{
				name:  "overwrite existing",
				key:   "test",
				value: &cartDomain.Cart{ID: "test", EntityID: "2"},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				require.NoError(t, storage.StoreCart(context.Background(), tt.value))

				cmd := client.Get(context.Background(), tt.key)
				require.NoError(t, cmd.Err())

				result, err := cmd.Bytes()
				require.NoError(t, err)

				buffer := new(bytes.Buffer)
				require.NoError(t, gob.NewEncoder(buffer).Encode(tt.value))

				assert.Equal(t, buffer.Bytes(), result)
			})
		}
	}

	t.Run("local-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("redis-server"); err != nil {
			t.Skip("redis-server not installed")
		}

		server, client := startUpLocalRedis(t)
		store := getRedisStorage("unix", server.Socket())
		runTestCases(t, store, client)
	})

	t.Run("docker-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("docker"); err != nil {
			t.Skip("docker not installed")
		}

		shutdown, address, client := startUpDockerRedis(t)
		defer shutdown()
		store := getRedisStorage("tcp", address)
		runTestCases(t, store, client)
	})
}

func TestRedisStorage_RemoveCart(t *testing.T) {
	t.Parallel()

	runTestCases := func(t *testing.T, storage *infrastructure.RedisStorage, client redis.UniversalClient) {
		t.Helper()

		tests := []struct {
			name        string
			key         string
			value       *cartDomain.Cart
			expectedErr bool
		}{
			{
				name:        "delete existing",
				key:         existingKey,
				value:       testCart,
				expectedErr: false,
			},
			{
				name:        "delete non existing",
				key:         "non",
				value:       nil,
				expectedErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := storage.RemoveCart(context.Background(), tt.value)
				assert.Equal(t, tt.expectedErr, err != nil)
				assert.False(t, client.Exists(context.Background(), tt.key).Val() > 0)
			})
		}
	}

	t.Run("local-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("redis-server"); err != nil {
			t.Skip("redis-server not installed")
		}

		server, client := startUpLocalRedis(t)
		store := getRedisStorage("unix", server.Socket())
		runTestCases(t, store, client)
	})

	t.Run("docker-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("docker"); err != nil {
			t.Skip("docker not installed")
		}

		shutdown, address, client := startUpDockerRedis(t)
		defer shutdown()
		store := getRedisStorage("tcp", address)
		runTestCases(t, store, client)
	})
}

func TestRedisStorage_MultipleStoragesSingleRedis(t *testing.T) {
	t.Parallel()

	runTestCase := func(t *testing.T, storageA, storageB *infrastructure.RedisStorage) {
		t.Helper()

		assert.True(t, storageA.HasCart(context.Background(), existingKey))
		assert.True(t, storageB.HasCart(context.Background(), existingKey))

		cartA, err := storageA.GetCart(context.Background(), existingKey)
		assert.NoError(t, err)

		cartA.EntityID = "storageA"

		err = storageA.StoreCart(context.Background(), cartA)
		assert.NoError(t, err)

		cartB, err := storageB.GetCart(context.Background(), existingKey)
		assert.NoError(t, err)
		assert.Equal(t, "storageA", cartB.EntityID)

		err = storageB.RemoveCart(context.Background(), cartB)
		assert.NoError(t, err)

		assert.False(t, storageA.HasCart(context.Background(), existingKey))
		assert.False(t, storageB.HasCart(context.Background(), existingKey))
	}

	t.Run("local-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("redis-server"); err != nil {
			t.Skip("redis-server not installed")
		}

		server, _ := startUpLocalRedis(t)
		storeA := getRedisStorage("unix", server.Socket())
		storeB := getRedisStorage("unix", server.Socket())
		runTestCase(t, storeA, storeB)
	})

	t.Run("docker-redis", func(t *testing.T) {
		t.Parallel()

		if _, err := exec.LookPath("docker"); err != nil {
			t.Skip("docker not installed")
		}

		shutdown, address, _ := startUpDockerRedis(t)
		defer shutdown()
		storeA := getRedisStorage("tcp", address)
		storeB := getRedisStorage("tcp", address)
		runTestCase(t, storeA, storeB)
	})
}
