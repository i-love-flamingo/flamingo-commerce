package states_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
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
	assert.Nil(t, s.Rollback(nil))
}

func TestPostRedirect_Run(t *testing.T) {
	s := states.PostRedirect{}
	p := &process.Process{}

	s.Run(context.Background(), p)

	assert.Equal(t, states.ValidatePayment{}.Name(), p.Context().State, "Next state should be ValidatePayment.")

	assert.Equal(t, s.Run(context.Background(), &process.Process{}), process.RunResult{})
}
