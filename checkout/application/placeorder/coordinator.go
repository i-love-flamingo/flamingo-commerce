package placeorder

import (
	"context"
	"encoding/gob"
	"errors"
	"net/url"
	"time"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (

	//TryLock interface
	TryLock interface {
		TryLock(string, time.Duration) (Unlock, error)
	}

	// Coordinator ensures that certain parts of the place order process are only done once at a time
	Coordinator struct {
		locker         TryLock
		logger         flamingo.Logger
		processFactory *process.Factory
	}

	//Unlock func
	Unlock func() error
)

var (
	//ErrLockTaken to indicate the lock is taken (by another running process)
	ErrLockTaken = errors.New("Lock already taken")
	//ErrNoPlaceOrderProcess if a requested process not running
	ErrNoPlaceOrderProcess = errors.New("ErrNoPlaceOrderProcess")
	//ErrAnotherPlaceOrderProcessRunning if a process runs
	ErrAnotherPlaceOrderProcessRunning = errors.New("ErrAnotherPlaceOrderProcessRunning")

	maxLockDuration = 2 * time.Minute
)

const (
	contextSessionStorageKey = "checkout_placeorder_context"
)

func init() {
	gob.Register(process.Context{})
}

//Inject dependencies
func (c *Coordinator) Inject(locker TryLock, logger flamingo.Logger, processFactory *process.Factory) {
	c.locker = locker
	c.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "placeorder")
	c.processFactory = processFactory
}

// New acquires lock if possible and creates new process with first run call blocking
// returns error if already locked or error during run
func (c *Coordinator) New(ctx context.Context, cart cartDomain.Cart, returnURL *url.URL) (*process.Context, error) {
	unlock, err := c.locker.TryLock(determineLockKeyForCart(cart), maxLockDuration)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = unlock()
	}()

	var runerr error
	var runpctx *process.Context
	web.RunWithDetachedContext(ctx, func(ctx context.Context) {
		has, err := c.HasUnfinishedProcess(ctx)
		if err != nil {
			runerr = err
			c.logger.Error(err)
			return
		}
		if has {
			runerr = ErrAnotherPlaceOrderProcessRunning
			c.logger.Info(runerr)
			return
		}

		newProcess, err := c.processFactory.New(returnURL, cart)
		if err != nil {
			runerr = err
			c.logger.Error(err)
			return
		}
		newProcess.Run(ctx)
		pctx := newProcess.Context()
		runpctx = &pctx
		err = c.storeProcessContext(ctx, pctx)
		if err != nil {
			runerr = err
			c.logger.Error(err)
			return
		}
	})

	return runpctx, runerr

	/* Todo:
	1. determine lock key based on session/cart
	2. check if place order process already in running if
	3. try to acquire lock
	4. create new place order context with state: new (Default start state should be configurable for project stuff)
	5. run first state transition by State->Run()
	6. store rollback callback
	*/
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
		return errors.New("Session not available to check for last placeorder context")
	}
	session.Store(contextSessionStorageKey, pctx)
	return nil
}

// LastProcess current place order process
func (c *Coordinator) LastProcess(ctx context.Context) (*process.Process, error) {
	session := web.SessionFromContext(ctx)
	if session == nil {
		return nil, errors.New("session not available to check for last placeorder context")
	}
	data, found := session.Load(contextSessionStorageKey)
	if !found {
		return nil, ErrNoPlaceOrderProcess
	}
	poContext, ok := data.(process.Context)
	if !ok {
		err := errors.New("context could not be read from session")
		c.logger.Error(err)
		return nil, err
	}

	p, err := c.processFactory.NewFromProcessContext(poContext)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Cancel the process if it exists (blocking)
// be aware that all rollback callbacks are executed
func (c *Coordinator) Cancel(ctx context.Context, cart cartDomain.Cart) error {
	/* Todo:
	1. Check if there is a place order process running
	2. Check if place order is in final state
	3. Wait for lock acquiring
	4. Run Rollbacks
	5. Set state to canceled
	6. Return result
	*/
	return nil
}

// Run starts the next processing if not already running
// Run is NOP if the process is locked
// Run returns immediately
func (c *Coordinator) Run(ctx context.Context) {
	// todo move to go routine
	web.RunWithDetachedContext(ctx, func(ctx context.Context) {
		has, err := c.HasUnfinishedProcess(ctx)
		if err != nil || has == false {
			return
		}

		p, err := c.LastProcess(ctx)
		if err != nil {
			return
		}

		unlock, err := c.locker.TryLock(determineLockKeyForProcess(p), maxLockDuration)
		if err != nil {
			return
		}
		defer func() {
			_ = unlock()
		}()

		p.Run(ctx)
		c.storeProcessContext(ctx, p.Context())
	})

	/* Todo: Do stuff in a go routine to be non blocking
	1. check if process is there
	2. check if lock can acquired
	3. hit state->run()
	*/
	return
}

// RunBlocking waits for the lock and starts the next processing
// RunBlocking waits until the process is finished and returns its result
func (c *Coordinator) RunBlocking(ctx context.Context) (*process.Context, error) {
	/* Todo:
	1. check if process is there
	2. get lock
	3. State->run()
	4. Return result
	*/
	return &process.Context{}, nil
}

func determineLockKeyForCart(cart cartDomain.Cart) string {
	return "checkout_placeorder_lock_" + cart.ID
}

func determineLockKeyForProcess(p *process.Process) string {
	return "checkout_placeorder_lock_" + p.Context().UUID
}
