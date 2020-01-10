package states

import (
	"context"
	"encoding/gob"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Failed state
	Failed struct {
		Reason process.FailedReason
	}
)

var _ process.FailedState = Failed{}

func init() {
	gob.Register(Failed{})
}

// Name get state name
func (f Failed) Name() string {
	return "Failed"
}

// Run the state operations
func (f Failed) Run(context.Context, *process.Process) process.RunResult {
	return process.RunResult{}
}

// Rollback the state operations
func (f Failed) Rollback(process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (f Failed) IsFinal() bool {
	return true
}

// SetFailedReason for the state
func (f Failed) SetFailedReason(reason process.FailedReason) process.FailedState {
	f.Reason = reason

	return f
}
