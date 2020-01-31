package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"go.opencensus.io/trace"
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
func (s Success) Run(ctx context.Context, _ *process.Process) process.RunResult {
	_, span := trace.StartSpan(ctx, "placeorder/state/Success/Run")
	defer span.End()

	return process.RunResult{}
}

// Rollback the state operations
func (s Success) Rollback(ctx context.Context, _ process.RollbackData) error {
	_, span := trace.StartSpan(ctx, "placeorder/state/Success/Rollback")
	defer span.End()

	return nil
}

// IsFinal if state is a final state
func (s Success) IsFinal() bool {
	return true
}
