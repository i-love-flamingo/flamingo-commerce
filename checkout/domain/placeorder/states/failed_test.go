package states_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFailed_IsFinal(t *testing.T) {
	s := states.Failed{}
	assert.True(t, s.IsFinal())
}

func TestFailed_Name(t *testing.T) {
	s := states.Failed{}
	assert.Equal(t, "Failed", s.Name())
}

func TestFailed_Rollback(t *testing.T) {
	s := states.Failed{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestFailed_Run(t *testing.T) {
	p := &process.Process{}
	s := new(states.Failed).Inject(&eventRouter{
		validator: func(event flamingo.Event) {
			require.IsType(t, &states.FailedEvent{}, event)
			assert.Equal(t, event.(*states.FailedEvent).ProcessContext, p.Context())
		},
	})
	assert.Equal(t, s.Run(context.Background(), &process.Process{}), process.RunResult{})
}
