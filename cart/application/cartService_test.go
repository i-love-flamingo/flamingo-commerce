package application_test

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"testing"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

func TestCartService_DeleteSavedSessionGuestCartID(t *testing.T) {
	type fields struct {
		CartReceiverService *cartApplication.CartReceiverService
		ProductService      productDomain.ProductService
		Logger              flamingo.Logger
		EventPublisher      events.EventPublisher
		EventRouter         flamingo.EventRouter
		RestrictionService  *validation.RestrictionService
		config              *struct {
			DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
			DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
		}
		DeliveryInfoBuilder cartDomain.DeliveryInfoBuilder
		CartCache           cartApplication.CartCache
		PlaceOrderService   placeorder.Service
	}
	type args struct {
		session *web.Session
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantErr       bool
		valuesCleared bool
	}{
		{
			name: "basic clearing of guest cart session value",
			fields: fields{
				CartReceiverService: func() *cartApplication.CartReceiverService {
					result := &cartApplication.CartReceiverService{}
					result.Inject(
						new(MockGuestCartServiceAdapter),
						new(MockCustomerCartService),
						func() *decorator.DecoratedCartFactory {
							result := &decorator.DecoratedCartFactory{}
							result.Inject(
								&MockProductService{},
								flamingo.NullLogger{},
							)

							return result
						}(),
						&authApplication.AuthManager{},
						&authApplication.UserService{},
						flamingo.NullLogger{},
						nil,
						&struct {
							CartCache cartApplication.CartCache `inject:",optional"`
						}{
							CartCache: new(MockCartCache),
						},
					)
					return result
				}(),
				ProductService: &MockProductService{},
				Logger:         flamingo.NullLogger{},
				EventPublisher: new(MockEventPublisher),
				config: &struct {
					DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
					DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
				}{
					DefaultDeliveryCode: "default_delivery_code",
					DeleteEmptyDelivery: false,
				},
				DeliveryInfoBuilder: new(MockDeliveryInfoBuilder),
				CartCache:           new(MockCartCache),
			},
			args: args{
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			wantErr:       false,
			valuesCleared: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cartApplication.CartService{}
			cs.Inject(
				tt.fields.CartReceiverService,
				tt.fields.ProductService,
				tt.fields.EventPublisher,
				tt.fields.EventRouter,
				tt.fields.DeliveryInfoBuilder,
				tt.fields.RestrictionService,
				tt.fields.Logger,
				tt.fields.config,
				nil,
			)

			err := cs.DeleteSavedSessionGuestCartID(tt.args.session)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.DeleteSavedSessionGuestCartID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.valuesCleared == true {
				if len(tt.args.session.Keys()) > 0 {
					t.Error("Session Values should be empty, but aren't")
				}
			}
		})
	}
}
