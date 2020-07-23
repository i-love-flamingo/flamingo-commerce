package states_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
)

func TestValidatePaymentSelection_IsFinal(t *testing.T) {
	assert.False(t, states.ValidatePaymentSelection{}.IsFinal())
}

func TestValidatePaymentSelection_Name(t *testing.T) {
	assert.Equal(t, "ValidatePaymentSelection", states.ValidatePaymentSelection{}.Name())
}

func TestValidatePaymentSelection_Rollback(t *testing.T) {
	s := states.ValidatePaymentSelection{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestValidatePaymentSelection_Run(t *testing.T) {
	isCalled := false
	validatorFactory := func(expectedPaymentSelection cart.PaymentSelection, err error) states.PaymentSelectionValidator {
		return func(paymentSelection cart.PaymentSelection) error {
			isCalled = true
			assert.Equal(t, expectedPaymentSelection, paymentSelection)

			return err
		}
	}

	cartWithPaymentSelection := provideCartWithPaymentSelection(t)

	tests := []struct {
		name                    string
		cart                    cart.Cart
		validator               states.PaymentSelectionValidator
		expectedResult          process.RunResult
		expectedValidatorCalled bool
		expectedState           string
	}{
		{
			name:      "no payment selection",
			cart:      cart.Cart{},
			validator: nil,
			expectedResult: process.RunResult{
				Failed: process.ErrorOccurredReason{Error: "no payment selection on cart"},
			},
			expectedValidatorCalled: false,
			expectedState:           states.New{}.Name(),
		},
		{
			name:                    "no validator",
			cart:                    cartWithPaymentSelection,
			validator:               nil,
			expectedResult:          process.RunResult{},
			expectedValidatorCalled: false,
			expectedState:           states.CreatePayment{}.Name(),
		},
		{
			name:                    "call validator",
			cart:                    cartWithPaymentSelection,
			validator:               validatorFactory(cartWithPaymentSelection.PaymentSelection, nil),
			expectedResult:          process.RunResult{},
			expectedValidatorCalled: true,
			expectedState:           states.CreatePayment{}.Name(),
		},
		{
			name:      "call validator with error",
			cart:      cartWithPaymentSelection,
			validator: validatorFactory(cartWithPaymentSelection.PaymentSelection, errors.New("validator error")),
			expectedResult: process.RunResult{
				Failed: process.ErrorOccurredReason{Error: "validator error"},
			},
			expectedValidatorCalled: true,
			expectedState:           states.New{}.Name(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCalled = false

			factory := provideProcessFactory(t)
			p, _ := factory.New(&url.URL{}, tt.cart)

			s := states.ValidatePaymentSelection{}
			s.Inject(&struct {
				Validator states.PaymentSelectionValidator `inject:",optional"`
			}{Validator: tt.validator})

			result := s.Run(context.Background(), p)
			assert.Equal(t, result, tt.expectedResult)
			assert.Equal(t, p.Context().CurrentStateName, tt.expectedState)
			assert.Equal(t, tt.expectedValidatorCalled, isCalled)
		})
	}
}
