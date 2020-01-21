package states

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Wait state
	Wait struct {
	}
)

var _ process.State = Wait{}

// Name get state name
func (w Wait) Name() string {
	return "Wait"
}

// Run the state operations
func (w Wait) Run(context.Context, *process.Process) process.RunResult {

	/* Todo: maybe split in multiple states
	1. Reserve Order id
	2. Start Payment
	3. Reserve Order if EarlyPlace
	4. Get Payment Status
	// n.ctx.UpdateState(WaitingForPaymentInformation)
	*/

	panic("implement me")
}

// Rollback the state operations
func (w Wait) Rollback(process.RollbackData) error {
	panic("implement me")
}

// IsFinal if state is a final state
func (w Wait) IsFinal() bool {
	return false
}
