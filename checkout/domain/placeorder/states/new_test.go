package states_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"github.com/stretchr/testify/assert"
)

func TestNew_Run(t *testing.T) {
	p := &process.Process{}
	state := states.New{}

	state.Run(context.Background(), p, nil)

	assert.Equal(t, states.CreatePayment{}.Name(), p.Context().CurrrentStateName, "Next state after New should be CreatePayment.")
}

func TestNew_IsFinal(t *testing.T) {
	state := states.New{}
	assert.False(t, state.IsFinal())
}
