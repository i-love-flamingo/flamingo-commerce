package states

import (
	"encoding/gob"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	//Canceled state
	Canceled struct {
	}
)

var _ process.State = Canceled{}

func init() {
	gob.Register(Canceled{})
}

//Name get state name
func (c Canceled) Name() string {
	return "Canceled"
}

//Run run state
func (c Canceled) Run(process *process.Process) *process.RollbackReference {
	// todo
	return nil
}

//IsFinal if state is a final state
func (c Canceled) IsFinal() bool {
	return true
}
