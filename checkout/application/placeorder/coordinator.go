package placeorder

import (
	"context"
	"encoding/gob"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/gorilla/sessions"
	"go.opencensus.io/trace"

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
		contextStore   process.ContextStore
		sessionStore   sessions.Store
		sessionName    string
	}

	//Unlock func
	Unlock func() error

	emptyResponseWriter struct{}
)

var (
	//ErrLockTaken to indicate the lock is taken (by another running process)
	ErrLockTaken = errors.New("lock already taken")
	//ErrNoPlaceOrderProcess if a requested process not running
	ErrNoPlaceOrderProcess = errors.New("ErrNoPlaceOrderProcess")
	//ErrAnotherPlaceOrderProcessRunning if a process runs
	ErrAnotherPlaceOrderProcessRunning = errors.New("ErrAnotherPlaceOrderProcessRunning")

	maxLockDuration = 2 * time.Minute
)

// emptyResponseWriter to be able to properly persist sessions
func (emptyResponseWriter) Header() http.Header       { return http.Header{} }
func (emptyResponseWriter) Write([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }
func (emptyResponseWriter) WriteHeader(int)           {}

func init() {
	gob.Register(process.Context{})
}

//Inject dependencies
func (c *Coordinator) Inject(locker TryLock, logger flamingo.Logger, processFactory *process.Factory,
	contextStore process.ContextStore, sessionStore sessions.Store,
	cfg *struct {
		SessionName string `inject:"config:flamingo.session.name,optional"`
	},
) {
	c.locker = locker
	c.logger = logger.WithField(flamingo.LogKeyModule, "checkout").WithField(flamingo.LogKeyCategory, "placeorder")
	c.processFactory = processFactory
	c.contextStore = contextStore
	c.sessionStore = sessionStore
	c.sessionName = cfg.SessionName
}

// New acquires lock if possible and creates new process with first run call blocking
// returns error if already locked or error during run
func (c *Coordinator) New(ctx context.Context, cart cartDomain.Cart, returnURL *url.URL) (*process.Context, error) {
	ctx, span := trace.StartSpan(ctx, "placeorder/coordinator/New")
	defer span.End()

	unlock, err := c.locker.TryLock(determineLockKeyForCart(cart), maxLockDuration)
	if err != nil {
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

	return c.contextStore.Store(session.ID(), pctx)
}

// LastProcess current place order process
func (c *Coordinator) LastProcess(ctx context.Context) (*process.Process, error) {
	session := web.SessionFromContext(ctx)
	if session == nil {
		return nil, errors.New("session not available to check for last place order context")
	}
	poContext, found := c.contextStore.Get(session.ID())
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
				unlock, err = c.locker.TryLock(determineLockKeyForProcess(p), maxLockDuration)
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
			_ = c.storeProcessContext(ctx, p.Context())
			c.updateGuestCartIDInSessionIfNeeded(ctx)
		})
	}(ctx)
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
				unlock, err = c.locker.TryLock(determineLockKeyForProcess(p), maxLockDuration)
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

		p.Run(ctx)
		err = c.storeProcessContext(ctx, p.Context())
		if err != nil {
			returnErr = err
			return
		}
		runPctx := p.Context()
		pctx = &runPctx
	})

	return pctx, returnErr
}

func (c *Coordinator) updateGuestCartIDInSessionIfNeeded(ctx context.Context) {
	rollbackSession := web.SessionFromContext(ctx)
	rollbackGuestCartID := rollbackSession.Try(application.GuestCartSessionKey)

	mostCurrentSession, err := c.sessionStore.Get(web.RequestFromContext(ctx).Request(), c.sessionName)
	if err != nil {
		c.logger.Error("couldn't receive current session from session store:", err)
		return
	}

	// id of rollback session and latest session are equal, nothing to do
	if rollbackGuestCartID == mostCurrentSession.Values[application.GuestCartSessionKey] {
		return
	}

	mostCurrentSession.Values[application.GuestCartSessionKey] = rollbackGuestCartID

	if rollbackGuestCartID == nil {
		delete(mostCurrentSession.Values, application.GuestCartSessionKey)
	}

	err = c.sessionStore.Save(web.RequestFromContext(ctx).Request(), new(emptyResponseWriter), mostCurrentSession)
	if err != nil {
		c.logger.Error("couldn't save session in the session store:", err)
		return
	}
}

func determineLockKeyForCart(cart cartDomain.Cart) string {
	return "checkout_placeorder_lock_" + cart.ID
}

func determineLockKeyForProcess(p *process.Process) string {
	return "checkout_placeorder_lock_" + p.Context().UUID
}
