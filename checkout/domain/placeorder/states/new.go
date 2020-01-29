package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"go.opencensus.io/trace"
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
func (n New) Run(ctx context.Context, p *process.Process, _ process.StateData) process.RunResult {
	_, span := trace.StartSpan(ctx, "placeorder/state/New/Run")
	defer span.End()

	p.UpdateState(CreatePayment{}.Name(), nil)

	return process.RunResult{}
}

// Rollback the state operations
func (n New) Rollback(ctx context.Context, _ process.RollbackData) error {
	_, span := trace.StartSpan(ctx, "placeorder/state/New/Rollback")
	defer span.End()

	return nil
}

// IsFinal if state is a final state
func (n New) IsFinal() bool {
	return false
}
