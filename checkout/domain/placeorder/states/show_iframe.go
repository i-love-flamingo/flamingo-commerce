package states

import (
	"context"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// ShowIframe state
	ShowIframe struct {
	}
)

var _ process.State = ShowIframe{}

//NewShowIframeStateData creates new state data for this state
func NewShowIframeStateData(url url.URL) process.StateData {
	return process.StateData(url)
}

// Name get state name
func (ShowIframe) Name() string {
	return "ShowIframe"
}

// Run the state operations
func (si ShowIframe) Run(_ context.Context, p *process.Process, stateData process.StateData) process.RunResult {
	// todo: add extracted validate switch here and only update state if something happened
	p.UpdateState(ValidatePayment{}.Name(), nil)
	return process.RunResult{}
}

// Rollback the state operations
func (si ShowIframe) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (si ShowIframe) IsFinal() bool {
	return false
}
