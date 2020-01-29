package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"go.opencensus.io/trace"
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
func (f Failed) Run(ctx context.Context, _ *process.Process, _ process.StateData) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/Failed/Run")
	defer span.End()

	return process.RunResult{}
}

// Rollback the state operations
func (f Failed) Rollback(ctx context.Context, _ process.RollbackData) error {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/Failed/Rollback")
	defer span.End()

	return nil
}

// IsFinal if state is a final state
func (f Failed) IsFinal() bool {
	return true
}
