package application_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	cartInfrastructure "flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	"flamingo.me/flamingo-commerce/v3/order/application"
	"flamingo.me/flamingo-commerce/v3/order/domain"
	domainMocks "flamingo.me/flamingo-commerce/v3/order/domain/mocks"
	authApplication "flamingo.me/flamingo/v3/core/auth/application"
	coreApplicationMocks "flamingo.me/flamingo/v3/core/auth/application/mocks"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/mock"
)

type (
	EventPublisherMock struct {
		mock.Mock
	}

	GuestCartServiceMock struct {
		mock.Mock
	}
)

func (epm *EventPublisherMock) PublishOrderPlacedEvent(ctx context.Context, c *cart.Cart, poi domain.PlacedOrderInfos) {
	epm.Called(ctx, c, poi)
}

func (gcsm *GuestCartServiceMock) GetBehaviour(ctx context.Context) (cart.Behaviour, error) {
	args := gcsm.Called(ctx)

	return args.Get(0).(cart.Behaviour), args.Error(1)
}

func (gcsm *GuestCartServiceMock) GetCart(ctx context.Context, cartID string) (*cart.Cart, error) {
	args := gcsm.Called(ctx, cartID)

	return args.Get(0).(*cart.Cart), args.Error(1)
}

func (gcsm *GuestCartServiceMock) GetNewCart(ctx context.Context) (*cart.Cart, error) {
	args := gcsm.Called(ctx)

	return args.Get(0).(*cart.Cart), args.Error(1)
}

func TestOrderService_PlaceOrder(t *testing.T) {
	type fields struct {
		logger flamingo.Logger
	}
	type args struct {
		ctx                context.Context
		session            *sessions.Session
		payment            *cart.CartPayment
		cart               *cart.Cart
		cartError          error
		cartBehaviour      cart.Behaviour
		cartBehaviourErr   error
		orderBehaviourErr  error
		loggedIn           bool
		placedOrderNumbers domain.PlacedOrderInfos
		placeOrderError    error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.PlacedOrderInfos
		wantErr bool
	}{
		{
			name: "not logged in, no cart",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				session:   &sessions.Session{},
				cartError: errors.New("no cart"),
			},
			wantErr: true,
		},
		{
			name: "not logged in, cart, no cart behaviour",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				session: &sessions.Session{
					Values: map[interface{}]interface{}{},
				},
				cart:             &cart.Cart{},
				cartBehaviour:    new(cartInfrastructure.InMemoryBehaviour),
				cartBehaviourErr: errors.New("no behaviour"),
			},
			wantErr: true,
		},
		{
			name: "not logged in, cart, cart behaviour, no order behaviour",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				session: &sessions.Session{
					Values: map[interface{}]interface{}{},
				},
				cart:              &cart.Cart{},
				cartBehaviour:     new(cartInfrastructure.InMemoryBehaviour),
				orderBehaviourErr: errors.New("no order behaviour"),
			},
			wantErr: true,
		},
		{
			name: "not logged in, cart, cart behaviour, order behaviour, any place order error",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				session: &sessions.Session{
					Values: map[interface{}]interface{}{},
				},
				cart:               &cart.Cart{},
				cartBehaviour:      new(cartInfrastructure.InMemoryBehaviour),
				placedOrderNumbers: domain.PlacedOrderInfos{},
				placeOrderError:    errors.New("some order place error"),
			},
			wantErr: true,
		},
		{
			name: "not logged in, cart, cart behaviour, order behaviour, cart not found error",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				session: &sessions.Session{
					Values: map[interface{}]interface{}{},
				},
				cart:               &cart.Cart{},
				cartBehaviour:      new(cartInfrastructure.InMemoryBehaviour),
				placedOrderNumbers: domain.PlacedOrderInfos{},
				placeOrderError:    cart.CartNotFoundError,
			},
			wantErr: true,
		},
		{
			name: "not logged in, cart, cart behaviour, order behaviour, order placed",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				session: &sessions.Session{
					Values: map[interface{}]interface{}{},
				},
				cart:          &cart.Cart{},
				cartBehaviour: new(cartInfrastructure.InMemoryBehaviour),
				placedOrderNumbers: domain.PlacedOrderInfos{
					domain.PlacedOrderInfo{
						OrderNumber:  "1",
						DeliveryCode: "delivery_1",
					},
					domain.PlacedOrderInfo{
						OrderNumber:  "2",
						DeliveryCode: "delivery_2",
					},
				},
			},
			wantErr: false,
			want: domain.PlacedOrderInfos{
				domain.PlacedOrderInfo{
					OrderNumber:  "1",
					DeliveryCode: "delivery_1",
				},
				domain.PlacedOrderInfo{
					OrderNumber:  "2",
					DeliveryCode: "delivery_2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := new(coreApplicationMocks.UserServiceInterface)
			userServiceMock.On("IsLoggedIn", mock.Anything, mock.Anything).Return(tt.args.loggedIn)

			guestCartServiceMock := new(GuestCartServiceMock)
			guestCartServiceMock.On("GetNewCart", mock.Anything).Return(tt.args.cart, tt.args.cartError)
			guestCartServiceMock.On("GetBehaviour", mock.Anything).Return(tt.args.cartBehaviour, tt.args.cartBehaviourErr)

			cartReceiverService := &cartApplication.CartReceiverService{}
			cartReceiverService.Inject(
				guestCartServiceMock,
				nil,
				userServiceMock,
				flamingo.NullLogger{},
				nil,
			)

			cartService := &cartApplication.CartService{}
			cartService.Inject(cartReceiverService, nil, nil, nil, nil, nil, nil, nil, nil)

			orderBehaviourMock := new(domainMocks.Behaviour)
			orderBehaviourMock.On("PlaceOrder", mock.Anything, tt.args.cart, tt.args.payment).Return(tt.args.placedOrderNumbers, tt.args.placeOrderError)

			guestOrderServiceMock := new(domainMocks.GuestOrderService)
			guestOrderServiceMock.On("GetBehaviour", mock.Anything).Return(orderBehaviourMock, tt.args.orderBehaviourErr)

			customerOrderServiceMock := new(domainMocks.CustomerOrderService)

			authManager := &authApplication.AuthManager{}

			orderReceiverService := &application.OrderReceiverService{}
			orderReceiverService.Inject(
				guestOrderServiceMock,
				customerOrderServiceMock,
				userServiceMock,
				authManager,
				flamingo.NullLogger{},
			)

			eventPublisherMock := new(EventPublisherMock)
			eventPublisherMock.On("PublishOrderPlacedEvent", tt.args.ctx, tt.args.cart, tt.args.placedOrderNumbers)

			os := application.OrderService{}
			os.Inject(
				tt.fields.logger,
				cartService,
				orderReceiverService,
				eventPublisherMock,
			)
			got, err := os.PlaceOrder(tt.args.ctx, tt.args.session, tt.args.payment)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.PlaceOrder() = %v, want %v", got, tt.want)
			}

			userServiceMock.AssertExpectations(t)

			if tt.args.cartError == nil {
				guestCartServiceMock.AssertExpectations(t)
			}

			if tt.args.cartError == nil && tt.args.cartBehaviourErr == nil {
				guestOrderServiceMock.AssertExpectations(t)
			}

			if tt.args.cartError == nil && tt.args.cartBehaviourErr == nil && tt.args.orderBehaviourErr == nil {
				orderBehaviourMock.AssertExpectations(t)
			}

			if tt.args.cartError == nil && tt.args.cartBehaviourErr == nil && tt.args.orderBehaviourErr == nil && tt.args.placeOrderError == nil {
				eventPublisherMock.AssertExpectations(t)
			}
		})
	}
}
