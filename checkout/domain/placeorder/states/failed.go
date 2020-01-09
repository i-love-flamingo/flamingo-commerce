package states

import (
	"encoding/gob"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//Failed state
	Failed struct {
	}
)

var _ process.State = Failed{}

func init() {
	gob.Register(Failed{})
}

//Name get state name
func (f Failed) Name() string {
	return "Failed"
}

//Run run state
func (f Failed) Run(process *process.Process) *process.RollbackReference {
	// todo
	return nil
}

//IsFinal if state is a final state
func (f Failed) IsFinal() bool {
	return true
}
