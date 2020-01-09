package states

import (
	"encoding/gob"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//Success state
	Success struct {
	}
)

var _ process.State = Success{}

func init() {
	gob.Register(Success{})
}

//Name get state name
func (c Success) Name() string {
	return "Success"
}

//Run run state
func (s Success) Run(process *process.Process) *process.RollbackReference {
	// todo
	return nil
}

//IsFinal if state is a final state
func (s Success) IsFinal() bool {
	return true
}
