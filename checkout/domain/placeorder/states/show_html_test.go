package states_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"github.com/stretchr/testify/assert"
)

func TestShowHTML_IsFinal(t *testing.T) {
	s := states.ShowHTML{}
	assert.False(t, s.IsFinal())
}

func TestShowHTML_Name(t *testing.T) {
	s := states.ShowHTML{}
	assert.Equal(t, "ShowHTML", s.Name())
}

func TestShowHTML_Rollback(t *testing.T) {
	s := states.ShowHTML{}
	assert.Nil(t, s.Rollback(nil, nil))
}

func TestShowHTML_Run(t *testing.T) {
	s := states.ShowHTML{}
	p := &process.Process{}

	s.Run(context.Background(), p)

	assert.Equal(t, states.ValidatePayment{}.Name(), p.Context().State, "Next state should be ValidatePayment.")

	assert.Equal(t, s.Run(context.Background(), &process.Process{}), process.RunResult{})
}
