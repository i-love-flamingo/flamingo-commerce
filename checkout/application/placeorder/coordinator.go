package placeorder

import (
	"context"
	"time"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	placeorderContext "flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/context"
)

type (
	// todo use/compare w/ domain interface
	locker interface {
		Lock(string, time.Duration) (func(), error)
	}

	// Coordinator ensures that certain parts of the place order process are only done once at a time
	Coordinator struct {
		Locker locker
	}
)

// New acquires lock if possible and creates new process with first run call blocking
// returns error if already locked or error during run
func (c *Coordinator) New(ctx context.Context, cart cartDomain.Cart) (placeorderContext.Context, error) {
	/* Todo:
	1. determine lock key based on session/cart
	2. check if place order process already in running if
	3. try to acquire lock
	4. create new place order context with state: new (Default start state should be configurable for project stuff)
	5. run first state transition by State->Run()
	6. store rollback callback
	*/

	return placeorderContext.Context{}, nil
}

// Current State of the process if it exists
func (c *Coordinator) Current(ctx context.Context, cart cartDomain.Cart) (placeorderContext.Context, error) {
	/* Todo:
	1. Check if there is a previous place order context/result
	2. Return if available, not available -> error
	*/
	return placeorderContext.Context{}, nil
}

// Cancel the process if it exists (blocking)
// be aware that all rollback callbacks are executed
func (c *Coordinator) Cancel(ctx context.Context, cart cartDomain.Cart) (placeorderContext.Context, error) {
	/* Todo:
	1. Check if there is a place order process running
	2. Check if place order is in final state
	3. Wait for lock acquiring
	4. Run Rollbacks
	5. Set state to canceled
	6. Return result
	*/
	return placeorderContext.Context{}, nil
}

// Run starts the next processing if not already running
// Run is NOP if the process is locked
// Run returns immediately
func (c *Coordinator) Run(ctx context.Context, cart cartDomain.Cart) {
	/* Todo: Do stuff in a go routine to be non blocking
	1. check if process is there
	2. check if lock can acquired
	3. hit state->run()
	*/
	return
}

// RunBlocking waits for the lock and starts the next processing
// RunBlocking waits until the process is finished and returns its result
func (c *Coordinator) RunBlocking(ctx context.Context, cart cartDomain.Cart) (placeorderContext.Context, error) {
	/* Todo:
	1. check if process is there
	2. get lock
	3. State->run()
	4. Return result
	*/
	return placeorderContext.Context{}, nil
}

func determineLockKey(ctx context.Context, cart cartDomain.Cart) string {
	return "myKey"
}
