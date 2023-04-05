package states_test

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart/mocks"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
)

func TestCompleteCart_IsFinal(t *testing.T) {
	assert.False(t, states.CompleteCart{}.IsFinal())
}

func TestCompleteCart_Name(t *testing.T) {
	assert.Equal(t, "CompleteCart", states.CompleteCart{}.Name())
}

func TestCompleteCart_Run(t *testing.T) {
	cart := cartDomain.Cart{ID: "mock_guest_cart"}

	tests := []struct {
		name           string
		behaviour      cartDomain.ModifyBehaviour
		behaviourError error
		validator      func(*testing.T, interface{})
		expectedState  string
		expectedResult process.RunResult
	}{
		{
			name: "successful completion",
			behaviour: func() cartDomain.ModifyBehaviour {
				behaviour := new(mocks.AllBehaviour)
				behaviour.CompleteBehaviour.EXPECT().Complete(mock.Anything, &cart).Return(&cart, nil, nil)
				return behaviour
			}(),
			behaviourError: nil,
			validator: func(t *testing.T, behaviour interface{}) {
				t.Helper()
				behaviour.(*mocks.AllBehaviour).CompleteBehaviour.AssertNumberOfCalls(t, "Complete", 1)
			},
			expectedState: states.PlaceOrder{}.Name(),
			expectedResult: process.RunResult{
				RollbackData: states.CompleteCartRollbackData{
					CompletedCart: &cart,
				},
			},
		},
		{
			name: "error on completion",
			behaviour: func() cartDomain.ModifyBehaviour {
				behaviour := new(mocks.AllBehaviour)
				behaviour.CompleteBehaviour.EXPECT().Complete(mock.Anything, &cart).Return(nil, nil, errors.New("test error"))
				return behaviour
			}(),
			behaviourError: nil,
			validator: func(t *testing.T, behaviour interface{}) {
				t.Helper()
				behaviour.(*mocks.AllBehaviour).CompleteBehaviour.AssertNumberOfCalls(t, "Complete", 1)
			},
			expectedState: states.New{}.Name(),
			expectedResult: process.RunResult{
				Failed: process.ErrorOccurredReason{Error: "test error"},
			},
		},
		{
			name: "no complete behaviour",
			behaviour: func() cartDomain.ModifyBehaviour {
				behaviour := new(mocks.ModifyBehaviour)
				return behaviour
			}(),
			behaviourError: nil,
			validator:      nil,
			expectedState:  states.PlaceOrder{}.Name(),
			expectedResult: process.RunResult{},
		},
		{
			name: "no behaviour",
			behaviour: func() cartDomain.ModifyBehaviour {
				return nil
			}(),
			behaviourError: errors.New("no behaviour"),
			validator:      nil,
			expectedState:  states.New{}.Name(),
			expectedResult: process.RunResult{
				Failed: process.ErrorOccurredReason{Error: "no behaviour"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := provideProcessFactory(t)

			p, _ := factory.New(&url.URL{}, cart)

			cartReceiverService := &cartApplication.CartReceiverService{}
			guestCartService := new(mocks.GuestCartService)
			guestCartService.EXPECT().GetModifyBehaviour(mock.Anything).Return(tt.behaviour, tt.behaviourError)
			guestCartService.EXPECT().GetNewCart(mock.Anything).Return(&cart, nil)
			cartReceiverService.Inject(
				guestCartService,
				nil,
				nil,
				nil,
				new(flamingo.NullLogger),
				nil,
				nil,
			)
			cartService := &cartApplication.CartService{}
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
				nil,
			)

			state := states.CompleteCart{}
			state.Inject(cartService, cartReceiverService)
			ctx := web.ContextWithSession(context.Background(), web.EmptySession())

			result := state.Run(ctx, p)
			assert.Equal(t, tt.expectedState, p.Context().CurrentStateName)

			if diff := deep.Equal(result, tt.expectedResult); diff != nil {
				t.Error("expected result is wrong: ", diff)
			}

			if tt.validator != nil {
				tt.validator(t, tt.behaviour)
			}
		})
	}
}

func TestCompleteCart_Rollback(t *testing.T) {
	cart := cartDomain.Cart{ID: "mock_guest_cart"}

	tests := []struct {
		name          string
		behaviour     cartDomain.ModifyBehaviour
		rollbackData  process.RollbackData
		validator     func(*testing.T, interface{})
		expectedError error
	}{
		{
			name: "successful restore",
			behaviour: func() cartDomain.ModifyBehaviour {
				behaviour := new(mocks.AllBehaviour)
				behaviour.CompleteBehaviour.EXPECT().Restore(mock.Anything, &cart).Return(&cart, nil, nil)
				return behaviour
			}(),
			rollbackData: states.CompleteCartRollbackData{
				CompletedCart: &cart,
			},
			validator: func(t *testing.T, behaviour interface{}) {
				t.Helper()
				behaviour.(*mocks.AllBehaviour).CompleteBehaviour.AssertNumberOfCalls(t, "Restore", 1)
			},
			expectedError: nil,
		},
		{
			name: "error on restore",
			behaviour: func() cartDomain.ModifyBehaviour {
				behaviour := new(mocks.AllBehaviour)
				behaviour.CompleteBehaviour.EXPECT().Restore(mock.Anything, &cart).Return(nil, nil, errors.New("test error"))
				return behaviour
			}(),
			rollbackData: states.CompleteCartRollbackData{
				CompletedCart: &cart,
			},
			validator: func(t *testing.T, behaviour interface{}) {
				t.Helper()
				behaviour.(*mocks.AllBehaviour).CompleteBehaviour.AssertNumberOfCalls(t, "Restore", 1)
			},
			expectedError: errors.New("test error"),
		},
		{
			name: "wrong rollback data",
			behaviour: func() cartDomain.ModifyBehaviour {
				behaviour := new(mocks.AllBehaviour)
				return behaviour
			}(),
			rollbackData:  nil,
			validator:     nil,
			expectedError: errors.New("rollback data not of expected type 'CompleteCartRollbackData', but <nil>"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cartReceiverService := &cartApplication.CartReceiverService{}
			guestCartService := new(mocks.GuestCartService)
			guestCartService.EXPECT().GetModifyBehaviour(mock.Anything).Return(tt.behaviour, nil)
			guestCartService.EXPECT().GetNewCart(mock.Anything).Return(&cart, nil)
			cartReceiverService.Inject(
				guestCartService,
				nil,
				nil,
				nil,
				new(flamingo.NullLogger),
				nil,
				nil,
			)
			cartService := &cartApplication.CartService{}
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
				nil,
			)

			state := states.CompleteCart{}
			state.Inject(cartService, cartReceiverService)
			ctx := web.ContextWithSession(context.Background(), web.EmptySession())
			err := state.Rollback(ctx, tt.rollbackData)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
