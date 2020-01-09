package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//Wait state
	Wait struct {
	}
)

var _ process.State = Wait{}

//SetProcess set process reference
func (n Wait) SetProcess(ctx *process.Process) {

}

//Run run state
func (n Wait) Run() (*process.RollbackReference, error) {
	/* Todo: maybe split in multiple states
	1. Reserve Order id
	2. Start Payment
	3. Reserve Order if EarlyPlace
	4. Get Payment Status
	// n.ctx.UpdateState(WaitingForPaymentInformation)
	*/

	return nil, nil
}

//IsFinal if state is a final state
func (n Wait) IsFinal() bool {
	return false
}
