package placeorder_test

import (
	"context"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure/contextstore"
	"flamingo.me/flamingo-commerce/v3/checkout/infrastructure/locker"
)

// slowUnlockLocker delays each unlock, so a run that starts before the lock is released finds it taken
type slowUnlockLocker struct {
	placeorder.TryLocker
}

func (l slowUnlockLocker) TryLock(ctx context.Context, key string, maxLockDuration time.Duration) (placeorder.Unlock, error) {
	unlock, err := l.TryLocker.TryLock(ctx, key, maxLockDuration)
	if err != nil {
		return nil, err
	}

	return func() error {
		time.Sleep(100 * time.Millisecond)

		return unlock()
	}, nil
}

// recordingState is a non-final state that records that it ran
type recordingState struct {
	ran atomic.Bool
}

func (s *recordingState) Run(context.Context, *process.Process) process.RunResult {
	s.ran.Store(true)

	return process.RunResult{}
}

func (s *recordingState) Rollback(context.Context, process.RollbackData) error { return nil }

func (s *recordingState) IsFinal() bool { return false }

func (s *recordingState) Name() string { return "start" }

func TestCoordinator_New(t *testing.T) {
	t.Parallel()

	t.Run("first run starts after the lock is released", func(t *testing.T) {
		t.Parallel()

		start := &recordingState{}

		factory := &process.Factory{}
		factory.Inject(
			func() *process.Process {
				return new(process.Process).Inject(map[string]process.State{start.Name(): start}, flamingo.NullLogger{}, nil)
			},
			&struct {
				StartState  process.State `inject:"startState"`
				FailedState process.State `inject:"failedState"`
			}{
				StartState:  start,
				FailedState: &states.Failed{},
			},
		)

		coordinator := &placeorder.Coordinator{}
		coordinator.Inject(slowUnlockLocker{locker.NewMemory()}, flamingo.NullLogger{}, factory, new(contextstore.Memory).Inject(), nil, nil, nil)

		ctx := web.ContextWithSession(context.Background(), web.EmptySession())

		_, err := coordinator.New(ctx, cartDomain.Cart{ID: "cart"}, &url.URL{})
		require.NoError(t, err)

		assert.Eventually(t, start.ran.Load, 5*time.Second, 10*time.Millisecond, "first state did not run")
	})
}
