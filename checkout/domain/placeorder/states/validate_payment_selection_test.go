package states_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
)

type paymentSelectionValidator struct {
	t                        *testing.T
	isCalled                 bool
	expectedPaymentSelection cart.PaymentSelection
	returnedError            error
}

func (p *paymentSelectionValidator) Validate(_ context.Context, _ *decorator.DecoratedCart, selection cart.PaymentSelection) error {
	p.isCalled = true
	assert.Equal(p.t, p.expectedPaymentSelection, selection)

	return p.returnedError
}

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
	cartWithPaymentSelection := provideCartWithPaymentSelection(t)
	tests := []struct {
		name                    string
		cart                    cart.Cart
		validator               validation.PaymentSelectionValidator
		expectedResult          process.RunResult
		expectedValidatorCalled bool
		expectedState           string
	}{
		{
			name:      "no payment selection",
			cart:      cart.Cart{},
			validator: nil,
			expectedResult: process.RunResult{
				Failed: process.PaymentErrorOccurredReason{Error: cart.ErrPaymentSelectionNotSet.Error()},
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
			name: "call validator",
			cart: cartWithPaymentSelection,
			validator: &paymentSelectionValidator{
				t:                        t,
				expectedPaymentSelection: cartWithPaymentSelection.PaymentSelection,
				returnedError:            nil,
			},
			expectedResult:          process.RunResult{},
			expectedValidatorCalled: true,
			expectedState:           states.CreatePayment{}.Name(),
		},
		{
			name: "call validator with error",
			cart: cartWithPaymentSelection,
			validator: &paymentSelectionValidator{
				t:                        t,
				expectedPaymentSelection: cartWithPaymentSelection.PaymentSelection,
				returnedError:            errors.New("validator error"),
			},
			expectedResult: process.RunResult{
				Failed: process.PaymentErrorOccurredReason{Error: "validator error"},
			},
			expectedValidatorCalled: true,
			expectedState:           states.New{}.Name(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := provideProcessFactory(t)
			p, _ := factory.New(&url.URL{}, tt.cart)

			s := states.ValidatePaymentSelection{}
			s.Inject(
				func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						nil,
						flamingo.NullLogger{},
					)

					return result
				}(),
				&struct {
					Validator validation.PaymentSelectionValidator `inject:",optional"`
				}{Validator: tt.validator})

			result := s.Run(context.Background(), p)
			assert.Equal(t, result, tt.expectedResult)
			assert.Equal(t, p.Context().CurrentStateName, tt.expectedState)
			if tt.validator != nil {
				assert.Equal(t, tt.expectedValidatorCalled, tt.validator.(*paymentSelectionValidator).isCalled)
			}
		})
	}
}
