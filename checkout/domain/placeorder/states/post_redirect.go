package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// PostRedirect state
	PostRedirect struct {
	}
)

var _ process.State = PostRedirect{}

// Name get state name
func (PostRedirect) Name() string {
	return "PostRedirect"
}

// Run the state operations
func (pr PostRedirect) Run(_ context.Context, p *process.Process) process.RunResult {
	p.UpdateState(ValidatePayment{}.Name())
	return process.RunResult{}
}

// Rollback the state operations
func (pr PostRedirect) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (pr PostRedirect) IsFinal() bool {
	return false
}
