package states_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"github.com/stretchr/testify/assert"
)

func TestValidatePayment_IsFinal(t *testing.T) {
	assert.False(t, states.ValidatePayment{}.IsFinal())
}

func TestValidatePayment_Name(t *testing.T) {
	assert.Equal(t, "ValidatePayment", states.ValidatePayment{}.Name())
}

func TestValidatePayment_Rollback(t *testing.T) {

}

func TestValidatePayment_Run(t *testing.T) {

}
