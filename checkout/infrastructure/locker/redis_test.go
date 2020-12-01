package locker_test

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stvp/tempredis"
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

func getRedisLocker(network, address string) *locker.Redis {
	redis := locker.NewRedis(&struct {
		MaxIdle                 int    `inject:"config:commerce.checkout.placeorder.lock.redis.maxIdle"`
		IdleTimeoutMilliseconds int    `inject:"config:commerce.checkout.placeorder.lock.redis.idleTimeoutMilliseconds"`
		Network                 string `inject:"config:commerce.checkout.placeorder.lock.redis.network"`
		Address                 string `inject:"config:commerce.checkout.placeorder.lock.redis.address"`
		Database                int    `inject:"config:commerce.checkout.placeorder.lock.redis.database"`
	}{MaxIdle: 3, IdleTimeoutMilliseconds: 240000, Network: network, Address: address, Database: 0})
	return redis
}

func TestRedis_TryLockDocker(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Skip("docker not installed")
	}

	res, err := pool.Run("redis", "5.0", nil)
	require.NoError(t, err, "failed to run redis docker image")
	defer func() { _ = pool.Purge(res) }()
	address := fmt.Sprintf("%s:%s", "127.0.0.1", res.GetPort("6379/tcp"))

	require.NoError(t, pool.Retry(func() error {
		conn, err := redis.Dial("tcp", address)
		if err != nil {
			return err
		}
		_, err = conn.Do("PING")
		return err
	}))

	redisLocker := getRedisLocker("tcp", address)
	runTestCases(t, redisLocker)
}

func TestRedis_TryLock(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	if _, err := exec.LookPath("redis-server"); err != nil {
		t.Skip("redis-server not installed")
	}

	server := startUp(t)
	defer func() { _ = server.Term() }()
	redisLocker := getRedisLocker("unix", server.Socket())

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
