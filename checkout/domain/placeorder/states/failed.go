package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"go.opencensus.io/trace"
)

type (
	// Failed state
	Failed struct {
		eventRouter flamingo.EventRouter
	}

	// FailedEvent is dispatched when the final failed state runs
	FailedEvent struct {
		ProcessContext process.Context
	}
)

var _ process.State = Failed{}

// Inject dependencies
func (f *Failed) Inject(
	eventRouter flamingo.EventRouter,
) *Failed {
	f.eventRouter = eventRouter

	return f
}

// Name get state name
func (f Failed) Name() string {
	return "Failed"
}

// Run the state operations
func (f Failed) Run(ctx context.Context, p *process.Process) process.RunResult {
	_, span := trace.StartSpan(ctx, "placeorder/state/Failed/Run")
	defer span.End()

	f.eventRouter.Dispatch(ctx, &FailedEvent{
		ProcessContext: p.Context(),
	})

	return process.RunResult{}
}

// Rollback the state operations
func (f Failed) Rollback(_ context.Context, _ process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (f Failed) IsFinal() bool {
	return true
}
