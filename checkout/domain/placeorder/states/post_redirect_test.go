package states_test

import (
	"context"
	"net/url"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/payment/application"

	"github.com/stretchr/testify/assert"
)

func TestPostRedirect_IsFinal(t *testing.T) {
	s := states.PostRedirect{}
	assert.False(t, s.IsFinal())
}

func TestPostRedirect_Name(t *testing.T) {
	s := states.PostRedirect{}
	assert.Equal(t, "PostRedirect", s.Name())
}

func TestPostRedirect_Rollback(t *testing.T) {
	s := states.PostRedirect{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestPostRedirect_Run(t *testing.T) {
	s := states.PostRedirect{}
	isCalled := false
	s.Inject(nil, func(_ context.Context, _ *process.Process, _ *application.PaymentService) process.RunResult {
		isCalled = true
		return process.RunResult{}
	})

	s.Run(context.Background(), nil)

	assert.True(t, isCalled)
}

func TestNewPostRedirectStateData(t *testing.T) {
	redirectURL := &url.URL{Host: "test.com"}
	formParameter := map[string]states.FormField{
		"test": {Value: []string{"abc"}},
	}

	assert.Equal(t,
		process.StateData(states.PostRedirectData{
			FormFields: formParameter,
			URL:        redirectURL,
		}),
		states.NewPostRedirectStateData(redirectURL, formParameter),
	)
}
