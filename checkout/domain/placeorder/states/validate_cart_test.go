package states_test

import (
	"context"
	"testing"

	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
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
		name           string
		isValid        bool
		expectedState  string
		expectedResult process.RunResult
	}{
		{
			name:           "Valid",
			isValid:        true,
			expectedState:  states.CreatePayment{}.Name(),
			expectedResult: process.RunResult{},
		},
		{
			name:          "Invalid",
			isValid:       false,
			expectedState: states.ValidateCart{}.Name(),
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
	guestCartService.On("GetNewCart", mock.Anything).Return(&cartDomain.Cart{ID: "mock_guest_cart"}, nil)
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
		new(authApplication.AuthManager),
		func() authApplication.UserServiceInterface {
			us := new(authApplication.UserService)
			us.Inject(new(authApplication.AuthManager), nil)
			return us
		}(),
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
			p.UpdateState(state.Name(), nil)
			ctx := web.ContextWithSession(context.Background(), web.EmptySession())
			result := state.Run(ctx, p, nil)

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
