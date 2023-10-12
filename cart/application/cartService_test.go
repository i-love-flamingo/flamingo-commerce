package application_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"flamingo.me/flamingo/v3/core/auth"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocks2 "flamingo.me/flamingo-commerce/v3/cart/domain/cart/mocks"
	mocks3 "flamingo.me/flamingo-commerce/v3/cart/domain/validation/mocks"
	"flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	"flamingo.me/flamingo-commerce/v3/product/domain/mocks"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
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
						nil,
						flamingo.NullLogger{},
						new(MockEventRouter),
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
				nil,
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

func TestCartService_AdjustItemsToRestrictedQty(t *testing.T) {
	type fields struct {
		CartReceiverService *cartApplication.CartReceiverService
		ProductService      productDomain.ProductService
		Logger              flamingo.Logger
		EventPublisher      *MockEventPublisher
		EventRouter         *MockEventRouter
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
		ctx     context.Context
		session *web.Session
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   cartApplication.QtyAdjustmentResults
	}{
		{
			name: "restrictors higher than qty dont reduce qty",
			fields: fields{
				CartReceiverService: func() *cartApplication.CartReceiverService {
					result := &cartApplication.CartReceiverService{}
					result.Inject(
						new(MockGuestCartServiceWithModifyBehaviour),
						new(MockCustomerCartService),
						func() *decorator.DecoratedCartFactory {
							result := &decorator.DecoratedCartFactory{}
							result.Inject(
								&MockProductService{},
								flamingo.NullLogger{},
							)

							return result
						}(),
						nil,
						flamingo.NullLogger{},
						new(MockEventRouter),
						&struct {
							CartCache cartApplication.CartCache `inject:",optional"`
						}{
							CartCache: new(MockCartWithItemCache),
						},
					)
					return result
				}(),
				ProductService: &MockProductService{},
				Logger:         flamingo.NullLogger{},
				EventRouter:    new(MockEventRouter),
				EventPublisher: new(MockEventPublisher),
				config: &struct {
					DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
					DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
				}{
					DefaultDeliveryCode: "default_delivery_code",
					DeleteEmptyDelivery: false,
				},
				DeliveryInfoBuilder: new(MockDeliveryInfoBuilder),
				CartCache:           new(MockCartWithItemCache),
				RestrictionService: func() *validation.RestrictionService {
					rs := &validation.RestrictionService{}
					rs.Inject(
						[]validation.MaxQuantityRestrictor{
							&MockRestrictor{IsRestricted: true, MaxQty: 10, DifferenceQty: 0},
						},
					)
					return rs
				}(),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want: cartApplication.QtyAdjustmentResults{},
		},
		{
			name: "restrictors lower than qty reduce qty",
			fields: fields{
				CartReceiverService: func() *cartApplication.CartReceiverService {
					result := &cartApplication.CartReceiverService{}
					result.Inject(
						new(MockGuestCartServiceWithModifyBehaviour),
						new(MockCustomerCartService),
						func() *decorator.DecoratedCartFactory {
							result := &decorator.DecoratedCartFactory{}
							result.Inject(
								&MockProductService{},
								flamingo.NullLogger{},
							)

							return result
						}(),
						nil,
						flamingo.NullLogger{},
						new(MockEventRouter),
						&struct {
							CartCache cartApplication.CartCache `inject:",optional"`
						}{
							CartCache: new(MockCartWithItemCache),
						},
					)
					return result
				}(),
				ProductService: &MockProductService{},
				Logger:         flamingo.NullLogger{},
				EventRouter:    new(MockEventRouter),
				EventPublisher: new(MockEventPublisher),
				config: &struct {
					DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
					DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
				}{
					DefaultDeliveryCode: "default_delivery_code",
					DeleteEmptyDelivery: false,
				},
				DeliveryInfoBuilder: new(MockDeliveryInfoBuilder),
				CartCache:           new(MockCartWithItemCache),
				RestrictionService: func() *validation.RestrictionService {
					rs := &validation.RestrictionService{}
					rs.Inject(
						[]validation.MaxQuantityRestrictor{
							&MockRestrictor{IsRestricted: true, MaxQty: 5, DifferenceQty: -2},
						},
					)
					return rs
				}(),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want: cartApplication.QtyAdjustmentResults{
				cartApplication.QtyAdjustmentResult{
					OriginalItem: cartDomain.Item{
						ID:  "mock_item",
						Qty: 7,
					},
					DeliveryCode: "default_delivery_code",
					WasDeleted:   false,
					RestrictionResult: &validation.RestrictionResult{
						IsRestricted:        true,
						MaxAllowed:          5,
						RemainingDifference: -2,
						RestrictorName:      "",
					},
					NewQty: 5,
				},
			},
		},
		{
			name: "maxAllowed of 0 deletes item",
			fields: fields{
				CartReceiverService: func() *cartApplication.CartReceiverService {
					result := &cartApplication.CartReceiverService{}
					result.Inject(
						new(MockGuestCartServiceWithModifyBehaviour),
						new(MockCustomerCartService),
						func() *decorator.DecoratedCartFactory {
							result := &decorator.DecoratedCartFactory{}
							result.Inject(
								&MockProductService{},
								flamingo.NullLogger{},
							)

							return result
						}(),
						nil,
						flamingo.NullLogger{},
						new(MockEventRouter),
						&struct {
							CartCache cartApplication.CartCache `inject:",optional"`
						}{
							CartCache: new(MockCartWithItemCache),
						},
					)
					return result
				}(),
				ProductService: &MockProductService{},
				Logger:         flamingo.NullLogger{},
				EventRouter:    new(MockEventRouter),
				EventPublisher: new(MockEventPublisher),
				config: &struct {
					DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
					DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
				}{
					DefaultDeliveryCode: "default_delivery_code",
					DeleteEmptyDelivery: false,
				},
				DeliveryInfoBuilder: new(MockDeliveryInfoBuilder),
				CartCache:           new(MockCartWithItemCache),
				RestrictionService: func() *validation.RestrictionService {
					rs := &validation.RestrictionService{}
					rs.Inject(
						[]validation.MaxQuantityRestrictor{
							&MockRestrictor{IsRestricted: true, MaxQty: 0, DifferenceQty: -7},
						},
					)
					return rs
				}(),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want: cartApplication.QtyAdjustmentResults{
				cartApplication.QtyAdjustmentResult{
					OriginalItem: cartDomain.Item{
						ID:  "mock_item",
						Qty: 7,
					},
					DeliveryCode: "default_delivery_code",
					WasDeleted:   true,
					RestrictionResult: &validation.RestrictionResult{
						IsRestricted:        true,
						MaxAllowed:          0,
						RemainingDifference: -7,
						RestrictorName:      "",
					},
					NewQty: 0,
				},
			},
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
				nil,
				tt.fields.Logger,
				tt.fields.config,
				nil,
			)
			got, _ := cs.AdjustItemsToRestrictedQty(tt.args.ctx, tt.args.session)

			if diff := deep.Equal(got, tt.want); diff != nil {
				t.Errorf("CartService.AdjustItemsToRestrictedQty() got!=want, diff: %#v", diff)
			}

		})
	}
}

func TestCartService_ReserveOrderIDAndSave(t *testing.T) {
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
		ctx     context.Context
		session *web.Session
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "reserve order id, reserved before",
			fields: fields{
				CartReceiverService: func() *cartApplication.CartReceiverService {
					result := &cartApplication.CartReceiverService{}
					result.Inject(
						new(MockGuestCartServiceWithModifyBehaviour),
						new(MockCustomerCartService),
						func() *decorator.DecoratedCartFactory {
							result := &decorator.DecoratedCartFactory{}
							result.Inject(
								&MockProductService{},
								flamingo.NullLogger{},
							)

							return result
						}(),
						nil,
						flamingo.NullLogger{},
						new(MockEventRouter),
						&struct {
							CartCache cartApplication.CartCache `inject:",optional"`
						}{
							CartCache: new(MockCartWithItemCacheWithAdditionalData),
						},
					)
					return result
				}(),
				ProductService: &MockProductService{},
				Logger:         flamingo.NullLogger{},
				EventRouter:    new(MockEventRouter),
				EventPublisher: new(MockEventPublisher),
				config: &struct {
					DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
					DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
				}{
					DefaultDeliveryCode: "default_delivery_code",
					DeleteEmptyDelivery: false,
				},
				DeliveryInfoBuilder: new(MockDeliveryInfoBuilder),
				CartCache:           new(MockCartWithItemCacheWithAdditionalData),
				RestrictionService:  nil,
				PlaceOrderService:   &MockPlaceOrderService{},
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want: "201910251128792ZM",
		},
		{
			name: "reserved order id, not reserved before",
			fields: fields{
				CartReceiverService: func() *cartApplication.CartReceiverService {
					result := &cartApplication.CartReceiverService{}
					result.Inject(
						new(MockGuestCartServiceWithModifyBehaviour),
						new(MockCustomerCartService),
						func() *decorator.DecoratedCartFactory {
							result := &decorator.DecoratedCartFactory{}
							result.Inject(
								&MockProductService{},
								flamingo.NullLogger{},
							)

							return result
						}(),
						nil,
						flamingo.NullLogger{},
						new(MockEventRouter),
						&struct {
							CartCache cartApplication.CartCache `inject:",optional"`
						}{
							CartCache: new(MockCartWithItemCache),
						},
					)
					return result
				}(),
				ProductService: &MockProductService{},
				Logger:         flamingo.NullLogger{},
				EventPublisher: new(MockEventPublisher),
				EventRouter:    new(MockEventRouter),
				config: &struct {
					DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
					DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
				}{
					DefaultDeliveryCode: "default_delivery_code",
					DeleteEmptyDelivery: false,
				},
				DeliveryInfoBuilder: new(MockDeliveryInfoBuilder),
				CartCache:           new(MockCartWithItemCache),
				RestrictionService:  nil,
				PlaceOrderService:   &MockPlaceOrderService{},
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want: "201910251128788TD",
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
				nil,
				tt.fields.Logger,
				tt.fields.config,
				&struct {
					CartValidator     validation.Validator      `inject:",optional"`
					ItemValidator     validation.ItemValidator  `inject:",optional"`
					CartCache         cartApplication.CartCache `inject:",optional"`
					PlaceOrderService placeorder.Service        `inject:",optional"`
				}{
					PlaceOrderService: tt.fields.PlaceOrderService,
				},
			)

			got, _ := cs.ReserveOrderIDAndSave(tt.args.ctx, tt.args.session)

			reservedOrderIDFromGot := got.AdditionalData.ReservedOrderID
			if reservedOrderIDFromGot != tt.want {
				t.Errorf("CartService.ReserveOrderIDAndSave() got!=want, got: %s, want: %s", reservedOrderIDFromGot, tt.want)
			}
		})
	}
}

func TestCartService_ForceReserveOrderIDAndSave(t *testing.T) {
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
		ctx     context.Context
		session *web.Session
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "force reserved order id, reserved before",
			fields: fields{
				CartReceiverService: func() *cartApplication.CartReceiverService {
					result := &cartApplication.CartReceiverService{}
					result.Inject(
						new(MockGuestCartServiceWithModifyBehaviour),
						new(MockCustomerCartService),
						func() *decorator.DecoratedCartFactory {
							result := &decorator.DecoratedCartFactory{}
							result.Inject(
								&MockProductService{},
								flamingo.NullLogger{},
							)

							return result
						}(),
						nil,
						flamingo.NullLogger{},
						new(MockEventRouter),
						&struct {
							CartCache cartApplication.CartCache `inject:",optional"`
						}{
							CartCache: new(MockCartWithItemCacheWithAdditionalData),
						},
					)
					return result
				}(),
				ProductService: &MockProductService{},
				Logger:         flamingo.NullLogger{},
				EventPublisher: new(MockEventPublisher),
				EventRouter:    new(MockEventRouter),
				config: &struct {
					DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
					DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
				}{
					DefaultDeliveryCode: "default_delivery_code",
					DeleteEmptyDelivery: false,
				},
				DeliveryInfoBuilder: new(MockDeliveryInfoBuilder),
				CartCache:           new(MockCartWithItemCacheWithAdditionalData),
				RestrictionService:  nil,
				PlaceOrderService:   &MockPlaceOrderService{},
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want: "201910251128788TD",
		},
		{
			name: "force reserved order id, not reserved before",
			fields: fields{
				CartReceiverService: func() *cartApplication.CartReceiverService {
					result := &cartApplication.CartReceiverService{}
					result.Inject(
						new(MockGuestCartServiceWithModifyBehaviour),
						new(MockCustomerCartService),
						func() *decorator.DecoratedCartFactory {
							result := &decorator.DecoratedCartFactory{}
							result.Inject(
								&MockProductService{},
								flamingo.NullLogger{},
							)

							return result
						}(),
						nil,
						flamingo.NullLogger{},
						new(MockEventRouter),
						&struct {
							CartCache cartApplication.CartCache `inject:",optional"`
						}{
							CartCache: new(MockCartWithItemCache),
						},
					)
					return result
				}(),
				ProductService: &MockProductService{},
				Logger:         flamingo.NullLogger{},
				EventPublisher: new(MockEventPublisher),
				EventRouter:    new(MockEventRouter),
				config: &struct {
					DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
					DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
				}{
					DefaultDeliveryCode: "default_delivery_code",
					DeleteEmptyDelivery: false,
				},
				DeliveryInfoBuilder: new(MockDeliveryInfoBuilder),
				CartCache:           new(MockCartWithItemCache),
				RestrictionService:  nil,
				PlaceOrderService:   &MockPlaceOrderService{},
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want: "201910251128788TD",
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
				nil,
				tt.fields.Logger,
				tt.fields.config,
				&struct {
					CartValidator     validation.Validator      `inject:",optional"`
					ItemValidator     validation.ItemValidator  `inject:",optional"`
					CartCache         cartApplication.CartCache `inject:",optional"`
					PlaceOrderService placeorder.Service        `inject:",optional"`
				}{
					PlaceOrderService: tt.fields.PlaceOrderService,
				},
			)

			got, _ := cs.ForceReserveOrderIDAndSave(tt.args.ctx, tt.args.session)

			reservedOrderIDFromGot := got.AdditionalData.ReservedOrderID
			if reservedOrderIDFromGot != tt.want {
				t.Errorf("CartService.ReserveOrderIDAndSave() got!=want, got: %s, want: %s", reservedOrderIDFromGot, tt.want)
			}
		})
	}
}

type (
	MockGuestCartServiceWithModifyBehaviour struct {
		cart *cartDomain.Cart
	}
)

func (m *MockGuestCartServiceWithModifyBehaviour) GetCart(ctx context.Context, cartID string) (*cartDomain.Cart, error) {
	if m.cart != nil {
		return m.cart, nil
	}

	return &cartDomain.Cart{
		ID: "mock_guest_cart",
	}, nil
}

func (m *MockGuestCartServiceWithModifyBehaviour) GetNewCart(ctx context.Context) (*cartDomain.Cart, error) {
	return &cartDomain.Cart{
		ID: "mock_guest_cart",
	}, nil
}

func (m *MockGuestCartServiceWithModifyBehaviour) GetModifyBehaviour(context.Context) (cartDomain.ModifyBehaviour, error) {
	cob := &infrastructure.DefaultCartBehaviour{}

	storage := &infrastructure.InMemoryCartStorage{}
	storage.Inject()

	cart := cartDomain.Cart{
		ID: "mock_guest_cart",
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{Code: "default_delivery_code"},
				Cartitems: []cartDomain.Item{
					{
						ID:  "mock_item",
						Qty: 7,
					},
				},
			},
		},
	}

	_ = storage.StoreCart(context.Background(), &cart)

	cob.Inject(
		storage,
		&MockProductService{},
		flamingo.NullLogger{},
		nil,
		nil,
		nil,
	)

	return cob, nil
}

func (m *MockGuestCartServiceWithModifyBehaviour) RestoreCart(ctx context.Context, cart cartDomain.Cart) (*cartDomain.Cart, error) {
	return &cart, nil
}

type (
	MockCartWithItemCache struct {
		CachedCart cartDomain.Cart
	}
)

func (m *MockCartWithItemCache) GetCart(context.Context, *web.Session, cartApplication.CartCacheIdentifier) (*cartDomain.Cart, error) {
	m.CachedCart = cartDomain.Cart{
		ID: "mock_guest_cart",
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{Code: "default_delivery_code"},
				Cartitems: []cartDomain.Item{
					{
						ID:  "mock_item",
						Qty: 7,
					},
				},
			},
		},
	}

	return &m.CachedCart, nil
}

func (m *MockCartWithItemCache) CacheCart(ctx context.Context, s *web.Session, cci cartApplication.CartCacheIdentifier, cart *cartDomain.Cart) error {
	m.CachedCart = *cart
	return nil
}

func (m *MockCartWithItemCache) Invalidate(context.Context, *web.Session, cartApplication.CartCacheIdentifier) error {
	return nil
}

func (m *MockCartWithItemCache) Delete(context.Context, *web.Session, cartApplication.CartCacheIdentifier) error {
	return nil
}

func (m *MockCartWithItemCache) DeleteAll(context.Context, *web.Session) error {
	return nil
}

func (m *MockCartWithItemCache) BuildIdentifier(context.Context, *web.Session) (cartApplication.CartCacheIdentifier, error) {
	return cartApplication.CartCacheIdentifier{}, nil
}

type (
	MockCartWithItemCacheWithAdditionalData struct {
		CachedCart cartDomain.Cart
	}
)

func (m *MockCartWithItemCacheWithAdditionalData) GetCart(context.Context, *web.Session, cartApplication.CartCacheIdentifier) (*cartDomain.Cart, error) {
	m.CachedCart = cartDomain.Cart{
		ID: "mock_guest_cart",
		Deliveries: []cartDomain.Delivery{
			{
				DeliveryInfo: cartDomain.DeliveryInfo{Code: "default_delivery_code"},
				Cartitems: []cartDomain.Item{
					{
						ID:  "mock_item",
						Qty: 7,
					},
				},
			},
		},
		AdditionalData: struct {
			CustomAttributes map[string]string
			ReservedOrderID  string
		}{CustomAttributes: nil, ReservedOrderID: "201910251128792ZM"},
	}

	return &m.CachedCart, nil
}

func (m *MockCartWithItemCacheWithAdditionalData) CacheCart(ctx context.Context, s *web.Session, cci cartApplication.CartCacheIdentifier, cart *cartDomain.Cart) error {
	m.CachedCart = *cart
	return nil
}

func (m *MockCartWithItemCacheWithAdditionalData) Invalidate(context.Context, *web.Session, cartApplication.CartCacheIdentifier) error {
	return nil
}

func (m *MockCartWithItemCacheWithAdditionalData) Delete(context.Context, *web.Session, cartApplication.CartCacheIdentifier) error {
	return nil
}

func (m *MockCartWithItemCacheWithAdditionalData) DeleteAll(context.Context, *web.Session) error {
	return nil
}

func (m *MockCartWithItemCacheWithAdditionalData) BuildIdentifier(context.Context, *web.Session) (cartApplication.CartCacheIdentifier, error) {
	return cartApplication.CartCacheIdentifier{}, nil
}

type MockRestrictor struct {
	IsRestricted  bool
	MaxQty        int
	DifferenceQty int
}

func (r *MockRestrictor) Name() string {
	return "MockRestrictor"
}

func (r *MockRestrictor) Restrict(ctx context.Context, session *web.Session, product productDomain.BasicProduct, currentCart *cartDomain.Cart, deliveryCode string) *validation.RestrictionResult {
	return &validation.RestrictionResult{
		IsRestricted:        r.IsRestricted,
		MaxAllowed:          r.MaxQty,
		RemainingDifference: r.DifferenceQty,
	}
}

type (
	MockPlaceOrderService struct{}
)

func (mpos *MockPlaceOrderService) PlaceGuestCart(ctx context.Context, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	return nil, nil
}

func (mpos *MockPlaceOrderService) PlaceCustomerCart(ctx context.Context, identity auth.Identity, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	return nil, nil
}

func (mpos *MockPlaceOrderService) ReserveOrderID(ctx context.Context, cart *cartDomain.Cart) (string, error) {
	return "201910251128788TD", nil
}

func (mpos *MockPlaceOrderService) CancelGuestOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos) error {
	return nil
}

func (mpos *MockPlaceOrderService) CancelCustomerOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos, identity auth.Identity) error {
	return nil
}

func TestCartService_CartInEvent(t *testing.T) {
	// bootstrap cart service
	cartReceiverService := func() *cartApplication.CartReceiverService {
		result := &cartApplication.CartReceiverService{}
		result.Inject(
			new(MockGuestCartServiceWithModifyBehaviour),
			new(MockCustomerCartService),
			func() *decorator.DecoratedCartFactory {
				result := &decorator.DecoratedCartFactory{}
				result.Inject(
					&MockProductService{},
					flamingo.NullLogger{},
				)

				return result
			}(),
			nil,
			flamingo.NullLogger{},
			new(MockEventRouter),
			&struct {
				CartCache cartApplication.CartCache `inject:",optional"`
			}{
				CartCache: new(MockCartWithItemCacheWithAdditionalData),
			},
		)
		return result
	}()
	productService := &MockProductService{}
	logger := flamingo.NullLogger{}
	eventPublisher := new(MockEventPublisher)
	eventRouter := new(MockEventRouter)
	eventRouter.On("Dispatch", mock.Anything, mock.Anything, mock.Anything).Return()
	config := &struct {
		DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
		DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
	}{
		DefaultDeliveryCode: "default_delivery_code",
		DeleteEmptyDelivery: false,
	}
	deliveryInfoBuilder := new(MockDeliveryInfoBuilder)
	restrictionService := func() *validation.RestrictionService {
		rs := &validation.RestrictionService{}
		rs.Inject(
			[]validation.MaxQuantityRestrictor{
				&MockRestrictor{},
			},
		)
		return rs
	}()

	// init cart service with dependencies
	cartService := cartApplication.CartService{}
	cartService.Inject(
		cartReceiverService,
		productService,
		eventPublisher,
		eventRouter,
		deliveryInfoBuilder,
		restrictionService,
		nil,
		logger,
		config,
		nil,
	)

	// add product to cart, we expect event to be thrown with updated cart
	addRequest := cartDomain.AddRequest{
		MarketplaceCode: "code-1",
		Qty:             1,
	}
	ctx := context.Background()
	session := web.EmptySession()
	_, err := cartService.AddProduct(ctx, session, "default_delivery_code", addRequest)
	assert.Nil(t, err)
	// white box tests that event router has been called as expected (once)
	eventRouter.AssertNumberOfCalls(t, "Dispatch", 1)
	// white box test that ensures router has been called with expected parameter (add to cart event)
	// with the expected marketplace code of the item
	eventRouter.AssertCalled(t, "Dispatch", ctx, fmt.Sprintf("%T", new(events.AddToCartEvent)), addRequest.MarketplaceCode)
}

func createCartServiceWithDependencies() *cartApplication.CartService {
	eventRouter := new(MockEventRouter)
	eventRouter.On("Dispatch", mock.Anything, mock.Anything, mock.Anything).Return()
	cartCache := new(MockCartCache)

	crs := &cartApplication.CartReceiverService{}
	crs.Inject(new(MockGuestCartServiceWithModifyBehaviour),
		nil,
		func() *decorator.DecoratedCartFactory {
			result := &decorator.DecoratedCartFactory{}
			result.Inject(
				&MockProductService{},
				flamingo.NullLogger{},
			)

			return result
		}(),
		nil,
		flamingo.NullLogger{},
		eventRouter,
		&struct {
			CartCache cartApplication.CartCache `inject:",optional"`
		}{
			CartCache: cartCache,
		})

	cs := &cartApplication.CartService{}
	cs.Inject(
		crs,
		&MockProductService{},
		new(MockEventPublisher),
		eventRouter,
		new(MockDeliveryInfoBuilder),
		func() *validation.RestrictionService {
			rs := &validation.RestrictionService{}
			rs.Inject(
				[]validation.MaxQuantityRestrictor{
					&MockRestrictor{},
				},
			)
			return rs
		}(),
		nil,
		flamingo.NullLogger{},
		&struct {
			DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
			DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
		}{
			DefaultDeliveryCode: "default_delivery_code",
			DeleteEmptyDelivery: false,
		},
		&struct {
			CartValidator     validation.Validator      `inject:",optional"`
			ItemValidator     validation.ItemValidator  `inject:",optional"`
			CartCache         cartApplication.CartCache `inject:",optional"`
			PlaceOrderService placeorder.Service        `inject:",optional"`
		}{
			CartCache: cartCache,
		},
	)
	return cs
}

func TestCartService_SetAdditionalData(t *testing.T) {
	cs := createCartServiceWithDependencies()

	cart, err := cs.UpdateAdditionalData(context.Background(), web.EmptySession(), map[string]string{"test": "data", "foo": "bar"})
	assert.NoError(t, err)
	assert.Equal(t, "data", cart.AdditionalData.CustomAttributes["test"])
	assert.Equal(t, "bar", cart.AdditionalData.CustomAttributes["foo"])
}

func TestCartService_SetAdditionalDataForDelivery(t *testing.T) {
	cs := createCartServiceWithDependencies()

	ctx := context.Background()
	session := web.EmptySession()
	addRequest := cartDomain.AddRequest{
		MarketplaceCode: "code-1",
		Qty:             1,
	}
	_, err := cs.AddProduct(ctx, session, "default_delivery_code", addRequest)
	assert.Nil(t, err)
	cart, err := cs.UpdateDeliveryAdditionalData(ctx, session, "default_delivery_code", map[string]string{"test": "data", "foo": "bar"})
	assert.NoError(t, err)
	var deliveryInfo *cartDomain.DeliveryInfo
	for _, delivery := range cart.Deliveries {
		if delivery.DeliveryInfo.Code == "default_delivery_code" {
			deliveryInfo = &delivery.DeliveryInfo
		}
	}

	require.NotNil(t, deliveryInfo)
	assert.Equal(t, "data", deliveryInfo.AdditionalData["test"])
	assert.Equal(t, "bar", deliveryInfo.AdditionalData["foo"])
}

func TestCartService_UpdateItemBundleConfig(t *testing.T) {
	t.Parallel()

	eventRouter := new(MockEventRouter)
	eventRouter.On("Dispatch", mock.Anything, mock.Anything, mock.Anything).Return()

	t.Run("error when bundle configuration not provided", func(t *testing.T) {
		t.Parallel()

		behaviour := &mocks2.ModifyBehaviour{}
		behaviour.EXPECT().UpdateItem(mock.Anything, mock.Anything, mock.Anything)

		cart := &cartDomain.Cart{
			Deliveries: []cartDomain.Delivery{
				{
					Cartitems: []cartDomain.Item{
						{
							ID:              "fakeID",
							MarketplaceCode: "fake",
						},
					},
				},
			},
		}

		guestCartService := &mocks2.GuestCartService{}
		guestCartService.EXPECT().GetModifyBehaviour(mock.Anything).Return(behaviour, nil)

		cartReceiverService, _ := getCartReceiverServiceForBundleUpdateTest(cart, guestCartService)

		cartService := &cartApplication.CartService{}
		cartService.Inject(
			cartReceiverService,
			&mocks.ProductService{},
			new(MockEventPublisher),
			eventRouter,
			new(MockDeliveryInfoBuilder),
			nil,
			nil,
			flamingo.NullLogger{},
			&struct {
				DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
				DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
			}{},
			&struct {
				CartValidator     validation.Validator      `inject:",optional"`
				ItemValidator     validation.ItemValidator  `inject:",optional"`
				CartCache         cartApplication.CartCache `inject:",optional"`
				PlaceOrderService placeorder.Service        `inject:",optional"`
			}{},
		)

		session := web.EmptySession()
		session.Store(cartApplication.GuestCartSessionKey, "fakeCartSession")

		updateCommand := cartDomain.ItemUpdateCommand{ItemID: "fakeID"}

		err := cartService.UpdateItemBundleConfig(context.Background(), session, updateCommand)
		assert.ErrorIs(t, err, cartApplication.ErrBundleConfigNotProvided)
	})

	t.Run("error when trying to update item that is not a bundle", func(t *testing.T) {
		t.Parallel()

		behaviour := &mocks2.ModifyBehaviour{}
		behaviour.EXPECT().UpdateItem(mock.Anything, mock.Anything, mock.Anything)

		cart := &cartDomain.Cart{
			Deliveries: []cartDomain.Delivery{
				{
					Cartitems: []cartDomain.Item{
						{
							ID:              "fakeID",
							MarketplaceCode: "fake",
						},
					},
				},
			},
		}

		guestCartService := &mocks2.GuestCartService{}
		guestCartService.EXPECT().GetModifyBehaviour(mock.Anything).Return(behaviour, nil)

		cartReceiverService, _ := getCartReceiverServiceForBundleUpdateTest(cart, guestCartService)

		product := productDomain.SimpleProduct{}

		productService := &mocks.ProductService{}
		productService.EXPECT().Get(mock.Anything, mock.Anything).Return(product, nil)

		cartService := &cartApplication.CartService{}
		cartService.Inject(
			cartReceiverService,
			productService,
			new(MockEventPublisher),
			eventRouter,
			new(MockDeliveryInfoBuilder),
			nil,
			nil,
			flamingo.NullLogger{},
			&struct {
				DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
				DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
			}{},
			&struct {
				CartValidator     validation.Validator      `inject:",optional"`
				ItemValidator     validation.ItemValidator  `inject:",optional"`
				CartCache         cartApplication.CartCache `inject:",optional"`
				PlaceOrderService placeorder.Service        `inject:",optional"`
			}{},
		)

		session := web.EmptySession()
		session.Store(cartApplication.GuestCartSessionKey, "fakeCartSession")

		updateCommand := cartDomain.ItemUpdateCommand{
			ItemID: "fakeID",
			BundleConfiguration: productDomain.BundleConfiguration{
				"choice1": {
					MarketplaceCode: "test",
					Qty:             1,
				},
			},
		}

		err := cartService.UpdateItemBundleConfig(context.Background(), session, updateCommand)
		assert.ErrorIs(t, err, cartApplication.ErrProductNotTypeBundle)
	})

	t.Run("if item validator is not provided, call update item on cart behaviour", func(t *testing.T) {
		t.Parallel()

		cart := &cartDomain.Cart{
			Deliveries: []cartDomain.Delivery{
				{
					Cartitems: []cartDomain.Item{
						{
							ID:              "fakeID",
							MarketplaceCode: "fake",
						},
					},
				},
			},
		}

		behaviour := &mocks2.ModifyBehaviour{}
		behaviour.EXPECT().UpdateItem(mock.Anything, mock.Anything, mock.Anything).Return(cart, nil, nil)

		guestCartService := &mocks2.GuestCartService{}
		guestCartService.EXPECT().GetModifyBehaviour(mock.Anything).Return(behaviour, nil)

		cartReceiverService, _ := getCartReceiverServiceForBundleUpdateTest(cart, guestCartService)

		product := productDomain.BundleProduct{}

		productService := &mocks.ProductService{}
		productService.EXPECT().Get(mock.Anything, mock.Anything).Return(product, nil)

		cartService := &cartApplication.CartService{}
		cartService.Inject(
			cartReceiverService,
			productService,
			new(MockEventPublisher),
			eventRouter,
			new(MockDeliveryInfoBuilder),
			nil,
			nil,
			flamingo.NullLogger{},
			&struct {
				DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
				DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
			}{},
			&struct {
				CartValidator     validation.Validator      `inject:",optional"`
				ItemValidator     validation.ItemValidator  `inject:",optional"`
				CartCache         cartApplication.CartCache `inject:",optional"`
				PlaceOrderService placeorder.Service        `inject:",optional"`
			}{},
		)

		session := web.EmptySession()
		session.Store(cartApplication.GuestCartSessionKey, "fakeCartSession")

		updateCommand := cartDomain.ItemUpdateCommand{
			ItemID: "fakeID",
			BundleConfiguration: productDomain.BundleConfiguration{
				"choice1": {
					MarketplaceCode: "test",
					Qty:             1,
				},
			},
		}

		err := cartService.UpdateItemBundleConfig(context.Background(), session, updateCommand)
		assert.NoError(t, err)

		behaviour.MethodCalled("UpdateItem")
	})

	t.Run("if item validator does not complain, call update item on cart behaviour", func(t *testing.T) {
		t.Parallel()

		cart := &cartDomain.Cart{
			Deliveries: []cartDomain.Delivery{
				{
					Cartitems: []cartDomain.Item{
						{
							ID:              "fakeID",
							MarketplaceCode: "fake",
						},
					},
				},
			},
		}

		behaviour := &mocks2.ModifyBehaviour{}
		behaviour.EXPECT().UpdateItem(mock.Anything, mock.Anything, mock.Anything).Return(cart, nil, nil)

		guestCartService := &mocks2.GuestCartService{}
		guestCartService.EXPECT().GetModifyBehaviour(mock.Anything).Return(behaviour, nil)

		cartReceiverService, _ := getCartReceiverServiceForBundleUpdateTest(cart, guestCartService)

		product := productDomain.BundleProduct{}

		productService := &mocks.ProductService{}
		productService.EXPECT().Get(mock.Anything, mock.Anything).Return(product, nil)

		itemValidator := &mocks3.ItemValidator{}
		itemValidator.EXPECT().Validate(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil)

		cartService := &cartApplication.CartService{}
		cartService.Inject(
			cartReceiverService,
			productService,
			new(MockEventPublisher),
			eventRouter,
			new(MockDeliveryInfoBuilder),
			nil,
			nil,
			flamingo.NullLogger{},
			&struct {
				DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
				DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
			}{},
			&struct {
				CartValidator     validation.Validator      `inject:",optional"`
				ItemValidator     validation.ItemValidator  `inject:",optional"`
				CartCache         cartApplication.CartCache `inject:",optional"`
				PlaceOrderService placeorder.Service        `inject:",optional"`
			}{
				ItemValidator: itemValidator,
			},
		)

		session := web.EmptySession()
		session.Store(cartApplication.GuestCartSessionKey, "fakeCartSession")

		updateCommand := cartDomain.ItemUpdateCommand{
			ItemID: "fakeID",
			BundleConfiguration: productDomain.BundleConfiguration{
				"choice1": {
					MarketplaceCode: "test",
					Qty:             1,
				},
			},
		}

		err := cartService.UpdateItemBundleConfig(context.Background(), session, updateCommand)
		assert.NoError(t, err)

		behaviour.MethodCalled("UpdateItem")
	})

	t.Run("if item validator fails, return an error and don't call update item on cart behaviour", func(t *testing.T) {
		t.Parallel()

		validationError := errors.New("some error")

		cart := &cartDomain.Cart{
			Deliveries: []cartDomain.Delivery{
				{
					Cartitems: []cartDomain.Item{
						{
							ID:              "fakeID",
							MarketplaceCode: "fake",
						},
					},
				},
			},
		}

		behaviour := &mocks2.ModifyBehaviour{}
		behaviour.EXPECT().UpdateItem(mock.Anything, mock.Anything, mock.Anything).Return(cart, nil, nil)

		guestCartService := &mocks2.GuestCartService{}
		guestCartService.EXPECT().GetModifyBehaviour(mock.Anything).Return(behaviour, nil)

		cartReceiverService, _ := getCartReceiverServiceForBundleUpdateTest(cart, guestCartService)

		product := productDomain.BundleProduct{}

		productService := &mocks.ProductService{}
		productService.EXPECT().Get(mock.Anything, mock.Anything).Return(product, nil)

		itemValidator := &mocks3.ItemValidator{}
		itemValidator.EXPECT().Validate(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(validationError)

		cartService := &cartApplication.CartService{}
		cartService.Inject(
			cartReceiverService,
			productService,
			new(MockEventPublisher),
			eventRouter,
			new(MockDeliveryInfoBuilder),
			nil,
			nil,
			flamingo.NullLogger{},
			&struct {
				DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
				DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
			}{},
			&struct {
				CartValidator     validation.Validator      `inject:",optional"`
				ItemValidator     validation.ItemValidator  `inject:",optional"`
				CartCache         cartApplication.CartCache `inject:",optional"`
				PlaceOrderService placeorder.Service        `inject:",optional"`
			}{
				ItemValidator: itemValidator,
			},
		)

		session := web.EmptySession()
		session.Store(cartApplication.GuestCartSessionKey, "fakeCartSession")

		updateCommand := cartDomain.ItemUpdateCommand{
			ItemID: "fakeID",
			BundleConfiguration: productDomain.BundleConfiguration{
				"choice1": {
					MarketplaceCode: "test",
					Qty:             1,
				},
			},
		}

		err := cartService.UpdateItemBundleConfig(context.Background(), session, updateCommand)
		assert.ErrorIs(t, err, validationError)
	})
}

func getCartReceiverServiceForBundleUpdateTest(cart *cartDomain.Cart, guestCartService cartDomain.GuestCartService) (*cartApplication.CartReceiverService, *MockCartCache) {
	cartCache := new(MockCartCache)
	cartCache.CachedCart = cart

	return func() (*cartApplication.CartReceiverService, *MockCartCache) {
		result := &cartApplication.CartReceiverService{}
		result.Inject(
			guestCartService,
			nil,
			func() *decorator.DecoratedCartFactory {
				result := &decorator.DecoratedCartFactory{}
				result.Inject(
					&MockProductService{},
					flamingo.NullLogger{},
				)

				return result
			}(),
			nil,
			flamingo.NullLogger{},
			new(MockEventRouter),
			&struct {
				CartCache cartApplication.CartCache `inject:",optional"`
			}{
				CartCache: cartCache,
			},
		)
		return result, cartCache
	}()
}
