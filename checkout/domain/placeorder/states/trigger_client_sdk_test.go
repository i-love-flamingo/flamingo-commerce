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

func TestTriggerClientSDK_IsFinal(t *testing.T) {
	s := states.TriggerClientSDK{}
	assert.False(t, s.IsFinal())
}

func TestTriggerClientSDK_Name(t *testing.T) {
	s := states.TriggerClientSDK{}
	assert.Equal(t, "TriggerClientSDK", s.Name())
}

func TestTriggerClientSDK_Rollback(t *testing.T) {
	s := states.TriggerClientSDK{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestTriggerClientSDK_Run(t *testing.T) {
	s := states.TriggerClientSDK{}
	isCalled := false
	s.Inject(nil, func(_ context.Context, _ *process.Process, _ *application.PaymentService) process.RunResult {
		isCalled = true
		return process.RunResult{}
	})

	s.Run(context.Background(), nil)

	assert.True(t, isCalled)
}

func TestNewTriggerClientSDKStateData(t *testing.T) {
	assert.Equal(t,
		process.StateData(states.TriggerClientSDKData{URL: &url.URL{Host: "test.com"}, Data: "data"}),
		states.NewTriggerClientSDKStateData(&url.URL{Host: "test.com"}, "data"))
}
