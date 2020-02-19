package states_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
)

func TestNew_Name(t *testing.T) {
	s := states.New{}
	assert.Equal(t, "New", s.Name())
}

func TestNew_Rollback(t *testing.T) {
	s := states.New{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestNew_Run(t *testing.T) {
	p := &process.Process{}
	state := states.New{}

	state.Run(context.Background(), p)

	assert.Equal(t, states.PrepareCart{}.Name(), p.Context().CurrentStateName, "Next state after New should be PrepareCart.")
}

func TestNew_IsFinal(t *testing.T) {
	state := states.New{}
	assert.False(t, state.IsFinal())
}
