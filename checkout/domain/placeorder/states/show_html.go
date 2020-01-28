package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// ShowHTML state
	ShowHTML struct {
	}
)

var _ process.State = ShowHTML{}

//NewShowHTMLStateData creates new StateData required for this ShowHTML state
func NewShowHTMLStateData(html string) process.StateData {
	return process.StateData(html)
}

// Name get state name
func (ShowHTML) Name() string {
	return "ShowHTML"
}

// Run the state operations
func (sh ShowHTML) Run(_ context.Context, p *process.Process, stateData process.StateData) process.RunResult {
	p.UpdateState(ValidatePayment{}.Name(), nil)
	return process.RunResult{}
}

// Rollback the state operations
func (sh ShowHTML) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (sh ShowHTML) IsFinal() bool {
	return false
}
