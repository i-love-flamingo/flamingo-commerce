package placeorder

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"net/url"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (

	// TryLocker port for a locking implementation
	TryLocker interface {
		// TryLock tries to get the lock for the provided key, if lock is already taken or couldn't be acquired function
		// returns an error. If the lock could be acquired a unlock function is returned which should be called to release the lock.
		// The provided duration is used in case that the node which required the lock dies so that the lock can released anyways.
		// If the node stays alive the lock time is not restricted in any way.
		TryLock(ctx context.Context, key string, maxLockDuration time.Duration) (Unlock, error)
	}

	// Unlock function to release the previously acquired lock, should be called within defer
	Unlock func() error

	// Coordinator ensures that certain parts of the place order process are only done once at a time
	Coordinator struct {
		locker         TryLocker
		logger         flamingo.Logger
		cartService    *application.CartService
		processFactory *process.Factory
		contextStore   process.ContextStore
		sessionStore   *web.SessionStore
		sessionName    string
		area           string
	}
)

// maxRunCount specifies the limit how often the coordinator should try to proceed in the state machine for a single call to Run / RunBlocking
const maxRunCount = 100

// waitForLockThrottle specifies the time to wait between attempts to get the lock for all blocking operations (cancel / runBlocking)
const waitForLockThrottle = 50 * time.Millisecond

var (
	// ErrLockTaken to indicate the lock is taken (by another running process)
	ErrLockTaken = errors.New("lock already taken")
	// ErrNoPlaceOrderProcess if a requested process not running
	ErrNoPlaceOrderProcess = errors.New("ErrNoPlaceOrderProcess")
	// ErrAnotherPlaceOrderProcessRunning if a process runs
	ErrAnotherPlaceOrderProcessRunning = errors.New("ErrAnotherPlaceOrderProcessRunning")

	maxLockDuration = 2 * time.Minute

	// startCount counts starts of new place order processes
	startCount = stats.Int64("flamingo-commerce/checkout/placeorder/starts", "Counts how often a new place order process was started", stats.UnitDimensionless)
)

func init() {
	gob.Register(process.Context{})
	err := opencensus.View("flamingo-commerce/checkout/placeorder/starts", startCount, view.Sum())
	if err != nil {
		panic(err)
	}

	stats.Record(context.Background(), startCount.M(0))
}

// Inject dependencies
func (c *Coordinator) Inject(
	locker TryLocker,
	logger flamingo.Logger,
	processFactory *process.Factory,
	contextStore process.ContextStore,
	sessionStore *web.SessionStore,
	cartService *application.CartService,
	cfg *struct {
		SessionName string `inject:"config:flamingo.session.name,optional"`
		Area        string `inject:"config:area"`
	},
) {
	c.locker = locker
	c.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "placeorder")
	c.processFactory = processFactory
	c.contextStore = contextStore
	c.sessionStore = sessionStore
	c.cartService = cartService

	if cfg != nil {
		c.area = cfg.Area
		c.sessionName = cfg.SessionName
	}
}

// New acquires lock if possible and creates new process with first run call blocking
// returns error if already locked or error during run
func (c *Coordinator) New(ctx context.Context, cart cartDomain.Cart, returnURL *url.URL) (*process.Context, error) {
	ctx, span := trace.StartSpan(ctx, "placeorder/coordinator/New")
	defer span.End()

	unlock, err := c.locker.TryLock(ctx, determineLockKeyForCart(cart), maxLockDuration)
	if err != nil {
		if err == ErrLockTaken {
			return nil, ErrAnotherPlaceOrderProcessRunning
		}
		return nil, err
	}
	defer func() {
		_ = unlock()
	}()

	var runErr error
	var runPCtx *process.Context
	web.RunWithDetachedContext(ctx, func(ctx context.Context) {
		has, err := c.HasUnfinishedProcess(ctx)
		if err != nil {
			runErr = err
			c.logger.Error(err)
			return
		}
		if has {
			runErr = ErrAnotherPlaceOrderProcessRunning
			c.logger.Info(runErr)
			return
		}

		censusCtx, _ := tag.New(ctx, tag.Upsert(opencensus.KeyArea, c.area))
		stats.Record(censusCtx, startCount.M(1))
		newProcess, err := c.processFactory.New(returnURL, cart)
		if err != nil {
			runErr = err
			c.logger.Error(err)
			return
		}
		pctx := newProcess.Context()
		runPCtx = &pctx
		err = c.storeProcessContext(ctx, pctx)
		if err != nil {
			runErr = err
			c.logger.Error(err)
			return
		}

		c.Run(ctx)
	})

	return runPCtx, runErr
}

// HasUnfinishedProcess checks for processes not in final state
func (c *Coordinator) HasUnfinishedProcess(ctx context.Context) (bool, error) {
	last, err := c.LastProcess(ctx)
	if err == ErrNoPlaceOrderProcess {
		return false, nil
	}
	if err != nil {
		return true, err
	}

	currentState, err := last.CurrentState()
	if err != nil {
		return true, err
	}

	return !currentState.IsFinal(), nil
}

func (c *Coordinator) storeProcessContext(ctx context.Context, pctx process.Context) error {
	session := web.SessionFromContext(ctx)
	if session == nil {
		return errors.New("session not available to check for last place order context")
	}

	return c.contextStore.Store(ctx, session.ID(), pctx)
}

func (c *Coordinator) clearProcessContext(ctx context.Context) error {
	session := web.SessionFromContext(ctx)
	if session == nil {
		return errors.New("session not available to check for last place order context")
	}

	return c.contextStore.Delete(ctx, session.ID())
}

// LastProcess current place order process
func (c *Coordinator) LastProcess(ctx context.Context) (*process.Process, error) {
	session := web.SessionFromContext(ctx)
	if session == nil {
		return nil, errors.New("session not available to check for last place order context")
	}
	poContext, found := c.contextStore.Get(ctx, session.ID())
	if !found {
		return nil, ErrNoPlaceOrderProcess
	}

	p, err := c.processFactory.NewFromProcessContext(poContext)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Cancel the process if it exists (blocking)
// be aware that all rollback callbacks are executed
func (c *Coordinator) Cancel(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "placeorder/coordinator/Cancel")
	defer span.End()

	var returnErr error
	web.RunWithDetachedContext(ctx, func(ctx context.Context) {
		{
			// scope things here to avoid using old process later
			p, err := c.LastProcess(ctx)
			if err != nil {
				returnErr = err
				return
			}
			var unlock Unlock
			err = ErrLockTaken
			for err == ErrLockTaken {
				unlock, err = c.locker.TryLock(ctx, determineLockKeyForProcess(p), maxLockDuration)
				// todo: add proper throttling

				time.Sleep(waitForLockThrottle)
			}
			if err != nil {
				returnErr = err
				return
			}
			defer func() {
				_ = unlock()
			}()
		}

		// lock acquired get fresh process state
		p, err := c.LastProcess(ctx)
		if err != nil {
			returnErr = err
			return
		}

		currentState, err := p.CurrentState()
		if err != nil {
			returnErr = err
			return
		}

		if currentState.IsFinal() {
			err = errors.New("process already in final state, cancel not possible")
			returnErr = err
			return
		}

		p.Failed(ctx, process.CanceledByCustomerReason{})
		err = c.storeProcessContext(ctx, p.Context())
		if err != nil {
			returnErr = err
		}
	})
	return returnErr
}

// ClearLastProcess removes last stored process
func (c *Coordinator) ClearLastProcess(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "placeorder/coordinator/Clear")
	defer span.End()

	var returnErr error
	web.RunWithDetachedContext(ctx, func(ctx context.Context) {
		err := c.clearProcessContext(ctx)
		if err != nil {
			returnErr = err
		}
	})
	return returnErr
}

// Run starts the next processing if not already running
// Run is NOP if the process is locked
// Run returns immediately
func (c *Coordinator) Run(ctx context.Context) {
	go func(ctx context.Context) {
		ctx, span := trace.StartSpan(ctx, "placeorder/coordinator/Run")
		defer span.End()

		web.RunWithDetachedContext(ctx, func(ctx context.Context) {
			has, err := c.HasUnfinishedProcess(ctx)
			if err != nil || !has {
				return
			}

			p, err := c.LastProcess(ctx)
			if err != nil {
				c.logger.Error("no last process on run: ", err)
				return
			}

			unlock, err := c.locker.TryLock(ctx, determineLockKeyForProcess(p), maxLockDuration)
			if err != nil {
				return
			}
			defer func() {
				_ = unlock()
			}()

			p, err = c.LastProcess(ctx)
			if err != nil {
				c.logger.Error("no last process on run: ", err)
				return
			}

			err = c.proceedInStateMachineUntilNoStateChange(ctx, p)
			if err != nil {
				c.logger.Error("proceeding in state machine failed: ", err)
				return
			}
		})
	}(ctx)
}

func (c *Coordinator) proceedInStateMachineUntilNoStateChange(ctx context.Context, p *process.Process) error {
	stateBeforeRun := p.Context().CurrentStateName
	for i := 0; i < maxRunCount; i++ {

		p.Run(ctx)
		err := c.storeProcessContext(ctx, p.Context())
		if err != nil {
			return err
		}
		c.forceSessionUpdate(ctx)
		stateAfterRun := p.Context().CurrentStateName
		if stateBeforeRun == stateAfterRun {
			return nil
		}
		stateBeforeRun = stateAfterRun
	}

	p.Failed(ctx, process.ErrorOccurredReason{
		Error: fmt.Sprintf("max run count %d of state machine reached", maxRunCount),
	})
	return nil
}

// RunBlocking waits for the lock and starts the next processing
// RunBlocking waits until the process is finished and returns its result
func (c *Coordinator) RunBlocking(ctx context.Context) (*process.Context, error) {
	ctx, span := trace.StartSpan(ctx, "placeorder/coordinator/RunBlocking")
	defer span.End()

	var pctx *process.Context
	var returnErr error
	web.RunWithDetachedContext(ctx, func(ctx context.Context) {
		{
			// scope things here to avoid continuing with an old process state
			p, err := c.LastProcess(ctx)
			if err != nil {
				returnErr = err
				return
			}

			var unlock Unlock
			err = ErrLockTaken
			for err == ErrLockTaken {
				unlock, err = c.locker.TryLock(ctx, determineLockKeyForProcess(p), maxLockDuration)
				// todo: add proper throttling
				time.Sleep(waitForLockThrottle)
			}
			if err != nil {
				returnErr = err
				return
			}

			defer func() {
				_ = unlock()
			}()
		}

		// lock acquired fetch everything new
		has, err := c.HasUnfinishedProcess(ctx)
		if err != nil {
			returnErr = err
			return
		}

		p, err := c.LastProcess(ctx)
		if err != nil {
			returnErr = err
			return
		}

		if !has {
			lastPctx := p.Context()
			pctx = &lastPctx
			return
		}

		// Load the most recent session, as we could have waited quite a while for the TryLock.
		session, err := c.sessionStore.LoadByID(ctx, web.SessionFromContext(ctx).ID())
		if err != nil {
			returnErr = err
			return
		}

		ctx = web.ContextWithSession(ctx, session)

		err = c.proceedInStateMachineUntilNoStateChange(ctx, p)
		if err != nil {
			returnErr = err
			return
		}
		runPctx := p.Context()
		pctx = &runPctx
	})

	return pctx, returnErr
}

func (c *Coordinator) forceSessionUpdate(ctx context.Context) {
	session := web.SessionFromContext(ctx)
	_, err := c.sessionStore.Save(ctx, session)
	if err != nil {
		c.logger.Error(err)
	}
}

func determineLockKeyForCart(cart cartDomain.Cart) string {
	return "checkout_placeorder_lock_" + cart.ID
}

func determineLockKeyForProcess(p *process.Process) string {
	return "checkout_placeorder_lock_" + p.Context().Cart.ID
}
