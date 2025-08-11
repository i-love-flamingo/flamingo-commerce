package locker_test

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stvp/tempredis"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/goleak"

	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure/locker"
)

func startUp(t *testing.T) *tempredis.Server {
	t.Helper()
	server, err := tempredis.Start(tempredis.Config{})
	if err != nil {
		t.Fatal(err)
	}

	return server
}

func getRedisLocker(network, address, username, password string) *locker.Redis {
	redis := locker.NewRedis(&struct {
		MaxIdle                 int    `inject:"config:commerce.checkout.placeorder.lock.redis.maxIdle"`
		IdleTimeoutMilliseconds int    `inject:"config:commerce.checkout.placeorder.lock.redis.idleTimeoutMilliseconds"`
		Network                 string `inject:"config:commerce.checkout.placeorder.lock.redis.network"`
		Address                 string `inject:"config:commerce.checkout.placeorder.lock.redis.address"`
		Database                int    `inject:"config:commerce.checkout.placeorder.lock.redis.database"`
		Username                string `inject:"config:commerce.checkout.placeorder.lock.redis.username,optional"`
		Password                string `inject:"config:commerce.checkout.placeorder.lock.redis.password,optional"`
		UseTLS                  bool   `inject:"config:commerce.checkout.placeorder.lock.redis.useTLS,optional"`
	}{MaxIdle: 3, IdleTimeoutMilliseconds: 240000, Network: network, Address: address, Database: 0, Username: username, Password: password})

	return redis
}

func TestRedis_TryLockDocker(t *testing.T) {
	ctx := context.Background()

	username := "myuser"
	password := "MySecurePassword"

	req := getContainerRequest(username, password)
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer func() { _ = redisC.Terminate(ctx) }()
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	port, err := redisC.MappedPort(ctx, "6379")
	require.NoError(t, err)

	host, err := redisC.Host(ctx)
	require.NoError(t, err)
	address := fmt.Sprintf("%s:%s", host, port.Port())

	redisLocker := getRedisLocker("tcp", address, username, password)
	runTestCases(t, redisLocker)
}

func TestRedis_TryLock(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	username := "myuser"
	password := "MySecurePassword"

	if _, err := exec.LookPath("redis-server"); err != nil {
		t.Skip("redis-server not installed")
	}

	server := startUp(t)
	defer func() { _ = server.Term() }()

	redisLocker := getRedisLocker("unix", server.Socket(), username, password)

	runTestCases(t, redisLocker)

}

func runTestCases(t *testing.T, redisLocker *locker.Redis) {
	t.Run("really locked", func(t *testing.T) {
		key := "test"
		start := time.Now()
		// get a long lock
		unlock, err := redisLocker.TryLock(context.Background(), key, time.Minute)
		require.NoError(t, err)

		// try to get same lock
		_, err = redisLocker.TryLock(context.Background(), key, time.Second)
		assert.Error(t, err)
		// assert if we were really in the lock period
		assert.True(t, time.Now().Before(start.Add(time.Minute)))

		// unlock
		require.NoError(t, unlock())

		// get the lock successfully again after unlock
		unlock, err = redisLocker.TryLock(context.Background(), key, time.Minute)
		require.NoError(t, err)
		require.NoError(t, unlock())
	})

	t.Run("lock should be expanded", func(t *testing.T) {
		key := "test_expanded"

		unlock, err := redisLocker.TryLock(context.Background(), key, 100*time.Millisecond)
		require.NoError(t, err)
		defer func() { _ = unlock() }()

		time.Sleep(200 * time.Millisecond)
		// try to get same lock
		_, err = redisLocker.TryLock(context.Background(), key, time.Second)
		assert.Error(t, err)
	})
}

func TestRedis_StatusDocker(t *testing.T) {
	t.Parallel()

	t.Run("redis not ready", func(t *testing.T) {
		t.Parallel()

		username := "myuser"
		password := "MySecurePassword"

		redisLocker := getRedisLocker("tcp", "127.0.0.1:80", username, password)
		alive, _ := redisLocker.Status()

		assert.False(t, alive)
	})

	t.Run("redis is there", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		username := "myuser"
		password := "MySecurePassword"

		req := getContainerRequest(username, password)

		redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		require.NoError(t, err)
		defer func() { _ = redisC.Terminate(ctx) }()
		defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

		port, err := redisC.MappedPort(ctx, "6379")
		require.NoError(t, err)

		host, err := redisC.Host(ctx)
		require.NoError(t, err)
		address := fmt.Sprintf("%s:%s", host, port.Port())

		redisLocker := getRedisLocker("tcp", address, username, password)
		alive, _ := redisLocker.Status()

		assert.True(t, alive)
	})

	t.Run("incorrect password", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		username := "myuser"
		password := "MySecurePassword"

		req := getContainerRequest(username, password)

		redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		require.NoError(t, err)

		defer func() { _ = redisC.Terminate(ctx) }()
		defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

		port, err := redisC.MappedPort(ctx, "6379")
		require.NoError(t, err)

		host, err := redisC.Host(ctx)
		require.NoError(t, err)

		address := fmt.Sprintf("%s:%s", host, port.Port())

		redisLocker := getRedisLocker("tcp", address, "username", "password")
		alive, _ := redisLocker.Status()

		assert.False(t, alive)
	})
}

func getContainerRequest(username, password string) testcontainers.ContainerRequest {
	return testcontainers.ContainerRequest{
		Image:        "valkey/valkey:7",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready to accept connections"),
			wait.ForListeningPort("6379/tcp")),
		Cmd: []string{
			"valkey-server",
			"--user default off",
			fmt.Sprintf("--user %s on >%s allcommands allkeys", username, password),
		},
	}
}
