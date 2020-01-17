package states

import (
	"context"
	"encoding/gob"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// CreatePayment state
	CreatePayment struct {
	}
)

var _ process.State = CreatePayment{}

func init() {
	gob.Register(CreatePayment{})
}

// Name get state name
func (CreatePayment) Name() string {
	return "CreatePayment"
}

// Run the state operations
func (c CreatePayment) Run(context.Context, *process.Process) process.RunResult {
	return process.RunResult{
		Failed: process.ErrorOccurredReason{Error: "not implemented"},
	}
}

// Rollback the state operations
func (c CreatePayment) Rollback(data process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (c CreatePayment) IsFinal() bool {
	return false
}
