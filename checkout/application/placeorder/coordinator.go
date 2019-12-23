package placeorder

import (
	"context"
	"time"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"
)

type (
	// todo use/compare w/ domain interface
	locker interface {
		Lock(string, time.Duration) (func(), error)
	}

	Coordinator struct {
		Locker locker
	}
)

// New acquires lock if possible and creates new process with first run call blocking
// returns error if already locked or error during run
func (c *Coordinator) New(ctx context.Context, cart cartDomain.Cart) (placeorder.Context, error) {

	return placeorder.Context{}, nil
}

// Current State of the process if it exists
func (c *Coordinator) Current(ctx context.Context, cart cartDomain.Cart) (placeorder.Context, error) {
	return placeorder.Context{}, nil
}

// Cancel the process if it exists (blocking)
// be aware that all rollback callbacks are executed
func (c *Coordinator) Cancel(ctx context.Context, cart cartDomain.Cart) (placeorder.Context, error) {
	return placeorder.Context{}, nil
}

// Run starts the next processing if not already running
// Run is NOP if the process is locked
// Run returns immediately
func (c *Coordinator) Run(ctx context.Context, cart cartDomain.Cart) {
	return
}

// RunBlocking waits for the lock and starts the next processing
// RunBlocking waits until the process is finished and returns its result
func (c *Coordinator) RunBlocking(ctx context.Context, cart cartDomain.Cart) (placeorder.Context, error) {
	return placeorder.Context{}, nil
}

func determineLockKey(ctx context.Context, cart cartDomain.Cart) string {
	return "myKey"
}
