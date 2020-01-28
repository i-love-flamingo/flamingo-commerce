package states

import (
	"context"
	"encoding/gob"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Redirect state
	Redirect struct {
	}
)

var _ process.State = Redirect{}

func init() {
	gob.Register(url.URL{})
}

//NewRedirectStateData - creates data required for this state
func NewRedirectStateData(url url.URL) process.StateData {
	return process.StateData(url)
}

// Name get state name
func (Redirect) Name() string {
	return "Redirect"
}

// Run the state operations
func (r Redirect) Run(_ context.Context, p *process.Process, data process.StateData) process.RunResult {
	p.UpdateState(ValidatePayment{}.Name(), nil)
	return process.RunResult{}
}

// Rollback the state operations
func (r Redirect) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (r Redirect) IsFinal() bool {
	return false
}
