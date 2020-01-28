package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Failed state
	Failed struct {
		Reason process.FailedReason
	}
)

var _ process.State = Failed{}

// Name get state name
func (f Failed) Name() string {
	return "Failed"
}

// Run the state operations
func (f Failed) Run(context.Context, *process.Process, process.StateData) process.RunResult {
	return process.RunResult{}
}

// Rollback the state operations
func (f Failed) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (f Failed) IsFinal() bool {
	return true
}
