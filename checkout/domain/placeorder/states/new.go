package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//New state
	New struct {
	}
)

var _ process.State = New{}

//SetProcess set process reference
func (n New) SetProcess(ctx *process.Process) {

}

//Run run state
func (n New) Run() (*process.RollbackReference, error) {
	return nil, nil
}

//IsFinal if state is a final state
func (n New) IsFinal() bool {
	return false
}
