package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//Canceled state
	Canceled struct {
	}
)

var _ process.State = Canceled{}

//SetProcess set process reference
func (c Canceled) SetProcess(ctx *process.Process) {

}

//Run run state
func (c Canceled) Run() (*process.RollbackReference, error) {
	// todo
	return nil, nil
}

//IsFinal if state is a final state
func (c Canceled) IsFinal() bool {
	return true
}
