package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// New state
	New struct {
	}
)

var _ process.State = New{}

// Name get state name
func (New) Name() string {
	return "New"
}

// Run the state operations
func (n New) Run(_ context.Context, p *process.Process) process.RunResult {
	p.UpdateState(CreatePayment{}.Name())

	return process.RunResult{}
}

// Rollback the state operations
func (n New) Rollback(process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (n New) IsFinal() bool {
	return false
}
