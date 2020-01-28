package states

import (
	"context"
	"encoding/gob"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// PostRedirect state
	PostRedirect struct {
	}

	// PostRedirectData holds details regarding the redirect
	PostRedirectData struct {
		FormFields map[string]FormField
		URL        url.URL
	}

	// FormField represents a form field to be displayed to the user
	FormField struct {
		Value []string
	}
)

func init() {
	gob.Register(PostRedirectData{})
}

var _ process.State = PostRedirect{}

//NewPostRedirectStateData creates new StateData with (persisted) Data required for this state
func NewPostRedirectStateData(url url.URL, formParameter map[string]FormField) process.StateData {
	return process.StateData(PostRedirectData{
		FormFields: formParameter,
		URL:        url,
	})
}

// Name get state name
func (PostRedirect) Name() string {
	return "PostRedirect"
}

// Run the state operations
func (pr PostRedirect) Run(_ context.Context, p *process.Process, stateData process.StateData) process.RunResult {
	p.UpdateState(ValidatePayment{}.Name(), nil)
	return process.RunResult{}
}

// Rollback the state operations
func (pr PostRedirect) Rollback(context.Context, process.RollbackData) error {
	return nil
}

// IsFinal if state is a final state
func (pr PostRedirect) IsFinal() bool {
	return false
}
