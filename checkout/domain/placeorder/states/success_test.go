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

type (
	eventRouter struct {
		validator func(flamingo.Event)
	}
)

func (e *eventRouter) Dispatch(_ context.Context, event flamingo.Event) {
	e.validator(event)
}

func TestSuccess_IsFinal(t *testing.T) {
	s := states.Success{}
	assert.True(t, s.IsFinal())
}

func TestSuccess_Name(t *testing.T) {
	s := states.Success{}
	assert.Equal(t, "Success", s.Name())
}

func TestSuccess_Rollback(t *testing.T) {
	s := states.Success{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestSuccess_Run(t *testing.T) {
	p := &process.Process{}
	s := new(states.Success).Inject(&eventRouter{
		validator: func(event flamingo.Event) {
			require.IsType(t, &states.SuccessEvent{}, event)
			assert.Equal(t, event.(*states.SuccessEvent).ProcessContext, p.Context())
		},
	})
	assert.Equal(t, s.Run(context.Background(), p), process.RunResult{})
}
