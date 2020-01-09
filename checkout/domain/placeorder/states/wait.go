package states

import (
	"encoding/gob"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//Wait state
	Wait struct {
	}
)

var _ process.State = Wait{}

func init() {
	gob.Register(Wait{})
}

//Name get state name
func (s Wait) Name() string {
	return "Wait"
}

//Run run state
func (s Wait) Run(process *process.Process) *process.RollbackReference {

	/* Todo: maybe split in multiple states
	1. Reserve Order id
	2. Start Payment
	3. Reserve Order if EarlyPlace
	4. Get Payment Status
	// n.ctx.UpdateState(WaitingForPaymentInformation)
	*/

	return nil
}

//IsFinal if state is a final state
func (s Wait) IsFinal() bool {
	return false
}
