package states_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
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
	p := &process.Process{}

	s.Run(context.Background(), p)

	assert.Equal(t, states.ValidatePayment{}.Name(), p.Context().State, "Next state should be ValidatePayment.")

	assert.Equal(t, s.Run(context.Background(), &process.Process{}), process.RunResult{})
}
