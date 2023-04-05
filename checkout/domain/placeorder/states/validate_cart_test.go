package states_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart/mocks"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	validator struct {
		Valid bool
	}
)

func (v *validator) Validate(_ context.Context, _ *web.Session, _ *decorator.DecoratedCart) validation.Result {
	return validation.Result{HasCommonError: !v.Valid}
}

func TestValidateCart_Name(t *testing.T) {
	s := states.ValidateCart{}
	assert.Equal(t, "ValidateCart", s.Name())
}

func TestValidateCart_Rollback(t *testing.T) {
	s := states.ValidateCart{}
	assert.Nil(t, s.Rollback(context.Background(), nil))
}

func TestValidateCart_Run(t *testing.T) {
	tests := []struct {
		name             string
		isValid          bool
		isGrandTotalZero bool
		expectedState    string
		expectedResult   process.RunResult
	}{
		{
			name:             "Valid cart that requires payment",
			isValid:          true,
			isGrandTotalZero: false,
			expectedState:    states.ValidatePaymentSelection{}.Name(),
			expectedResult:   process.RunResult{},
		},
		{
			name:             "Valid cart that is fully discounted, no payment needed",
			isValid:          true,
			isGrandTotalZero: true,
			expectedState:    states.CompleteCart{}.Name(),
			expectedResult:   process.RunResult{},
		},
		{
			name:             "Invalid",
			isValid:          false,
			isGrandTotalZero: false,
			expectedState:    states.ValidateCart{}.Name(),
			expectedResult: process.RunResult{
				RollbackData: nil,
				Failed: process.CartValidationErrorReason{
					ValidationResult: validation.Result{
						HasCommonError: true,
					},
				},
			},
		},
	}

	// global service setup
	cartReceiverService := &application.CartReceiverService{}
	guestCartService := new(mocks.GuestCartService)
	guestCartService.EXPECT().GetNewCart(mock.Anything).Return(&cartDomain.Cart{ID: "mock_guest_cart"}, nil)
	cartReceiverService.Inject(
		guestCartService,
		new(mocks.CustomerCartService),
		func() *decorator.DecoratedCartFactory {
			result := &decorator.DecoratedCartFactory{}
			result.Inject(
				nil,
				flamingo.NullLogger{},
			)

			return result
		}(),
		nil,
		new(flamingo.NullLogger),
		nil,
		nil,
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cartService := application.CartService{}
			cartService.Inject(
				cartReceiverService,
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				new(flamingo.NullLogger),
				nil,
				&struct {
					CartValidator     validation.Validator     `inject:",optional"`
					ItemValidator     validation.ItemValidator `inject:",optional"`
					CartCache         application.CartCache    `inject:",optional"`
					PlaceOrderService placeorder.Service       `inject:",optional"`
				}{CartValidator: &validator{Valid: tt.isValid}, ItemValidator: nil, CartCache: nil, PlaceOrderService: nil},
			)
			state := new(states.ValidateCart).Inject(&cartService)
			p := &process.Process{}
			cart := cartDomain.Cart{
				ID:         "cart-id",
				EntityID:   "entity-id",
				GrandTotal: domain.NewFromInt(1, 1, "EUR"),
			}

			if tt.isGrandTotalZero {
				cart.GrandTotal = domain.NewFromInt(0, 1, "EUR")
			}

			p.UpdateCart(cart)
			p.UpdateState(state.Name(), nil)
			ctx := web.ContextWithSession(context.Background(), web.EmptySession())
			result := state.Run(ctx, p)

			assert.Equal(t, tt.expectedState, p.Context().CurrentStateName, "Next state after ValidateCart should be CreatePayment.")
			if diff := deep.Equal(result, tt.expectedResult); diff != nil {
				t.Error("expected re is wrong: ", diff)
			}
		})
	}

}

func TestValidateCart_IsFinal(t *testing.T) {
	state := states.ValidateCart{}
	assert.False(t, state.IsFinal())
}
