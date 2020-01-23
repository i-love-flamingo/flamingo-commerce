package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Success state
	Success struct{}
)

var _ process.State = Success{}

// Name get state name
func (s Success) Name() string {
	return "Success"
}

// Run the state operations
func (s Success) Run(context.Context, *process.Process) process.RunResult {
	return process.RunResult{}
}

// Rollback the state operations
func (s Success) Rollback(process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (s Success) IsFinal() bool {
	return true
}
