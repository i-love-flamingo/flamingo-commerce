package states

import (
	"encoding/gob"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//New state
	New struct {
	}
)

var _ process.State = New{}

func init() {
	gob.Register(New{})
}

//Name get state name
func (c New) Name() string {
	return "New"
}

//Run run state
func (n New) Run(process *process.Process) *process.RollbackReference {
	// todo
	return nil
}

//IsFinal if state is a final state
func (n New) IsFinal() bool {
	return false
}
