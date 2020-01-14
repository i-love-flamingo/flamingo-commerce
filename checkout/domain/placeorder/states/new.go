package states

import (
	"context"
	"encoding/gob"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// New state
	New struct {
	}

	Test struct {
		A string
	}
)

var _ process.State = New{}

func init() {
	gob.Register(New{})
	gob.Register(Test{})
}

// Name get state name
func (New) Name() string {
	return "New"
}

// Run the state operations
func (n New) Run(context.Context, *process.Process) process.RunResult {
	return process.RunResult{
		RollbackData: &Test{
			A: "A",
		},
		Failed: process.ErrorOccurredReason{Error: "not implemented"},
	}
}

// Rollback the state operations
func (n New) Rollback(data process.RollbackData) error {
	fmt.Println(data.(*Test))

	return nil
}

// IsFinal if state is a final state
func (n New) IsFinal() bool {
	return false
}
