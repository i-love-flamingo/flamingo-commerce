package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"go.opencensus.io/trace"
)

type (
	// Success state
	Success struct {
		eventRouter flamingo.EventRouter
	}

	// SuccessEvent is dispatched when the final success state runs
	SuccessEvent struct {
		ProcessContext process.Context
	}
)

var _ process.State = Success{}

// Inject dependencies
func (s *Success) Inject(
	eventRouter flamingo.EventRouter,
) *Success {
	s.eventRouter = eventRouter

	return s
}

// Name get state name
func (s Success) Name() string {
	return "Success"
}

// Run the state operations
func (s Success) Run(ctx context.Context, p *process.Process) process.RunResult {
	_, span := trace.StartSpan(ctx, "placeorder/state/Success/Run")
	defer span.End()

	s.eventRouter.Dispatch(ctx, &SuccessEvent{
		ProcessContext: p.Context(),
	})

	return process.RunResult{}
}

// Rollback the state operations
func (s Success) Rollback(_ context.Context, _ process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (s Success) IsFinal() bool {
	return true
}
