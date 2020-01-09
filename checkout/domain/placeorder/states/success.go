package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//Success state
	Success struct {
	}
)

var _ process.State = Success{}

//SetProcess set process reference
func (s Success) SetProcess(ctx *process.Process) {

}

//Run run state
func (s Success) Run() (*process.RollbackReference, error) {
	// todo
	return nil, nil
}

//IsFinal if state is a final state
func (s Success) IsFinal() bool {
	return true
}
