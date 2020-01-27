package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Wait state
	Wait struct {
	}
)

var _ process.State = Wait{}

// TODO: remove state since its not an internal state..

// Name get state name
func (w Wait) Name() string {
	return "Wait"
}

// Run the state operations
func (w Wait) Run(context.Context, *process.Process) process.RunResult {
	panic("implement me")
}

// Rollback the state operations
func (w Wait) Rollback(context.Context, process.RollbackData) error {
	panic("implement me")
}

// IsFinal if state is a final state
func (w Wait) IsFinal() bool {
	return false
}
