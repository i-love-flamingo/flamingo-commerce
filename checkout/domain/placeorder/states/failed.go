package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//Failed state
	Failed struct {
	}
)

var _ process.State = Failed{}

//SetProcess set process reference
func (f Failed) SetProcess(ctx *process.Process) {
}

//Run run state
func (f Failed) Run() (*process.RollbackReference, error) {
	// todo
	return nil, nil
}

//IsFinal if state is a final state
func (f Failed) IsFinal() bool {
	return true
}
