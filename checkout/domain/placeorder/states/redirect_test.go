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

func TestRedirect_IsFinal(t *testing.T) {
	s := states.Redirect{}
	assert.False(t, s.IsFinal())
}

func TestRedirect_Name(t *testing.T) {
	s := states.Redirect{}
	assert.Equal(t, "Redirect", s.Name())
}

func TestRedirect_Rollback(t *testing.T) {
	s := states.Redirect{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestRedirect_Run(t *testing.T) {
	s := states.Redirect{}
	isCalled := false
	s.Inject(nil, func(_ context.Context, _ *process.Process, _ *application.PaymentService) process.RunResult {
		isCalled = true
		return process.RunResult{}
	})

	s.Run(context.Background(), nil)

	assert.True(t, isCalled)
}

func TestNewRedirectStateData(t *testing.T) {
	assert.Equal(t, process.StateData(&url.URL{Host: "test.com"}), states.NewRedirectStateData(&url.URL{Host: "test.com"}))
}
