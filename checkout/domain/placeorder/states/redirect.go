package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Redirect state
	Redirect struct {
	}
)

var _ process.State = Redirect{}

// Name get state name
func (Redirect) Name() string {
	return "Redirect"
}

// Run the state operations
func (r Redirect) Run(_ context.Context, p *process.Process) process.RunResult {
	p.UpdateState(ValidatePayment{}.Name())
	return process.RunResult{}
}

// Rollback the state operations
func (r Redirect) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (r Redirect) IsFinal() bool {
	return false
}
