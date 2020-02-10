package application_test

import (
	"context"
	"fmt"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"flamingo.me/flamingo-commerce/v3/cart/infrastructure"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"

	"testing"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
	oauth "flamingo.me/flamingo/v3/core/oauth/domain"
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
			authmanager := &authApplication.AuthManager{}
			authmanager.Inject(
				flamingo.NullLogger{},
				nil, nil,
			)
			cs.Inject(
				tt.fields.CartReceiverService,
				tt.fields.ProductService,
				tt.fields.EventPublisher,
				tt.fields.EventRouter,
				tt.fields.DeliveryInfoBuilder,
				tt.fields.RestrictionService,
				authmanager,
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
	authmanager := &authApplication.AuthManager{}
	authmanager.Inject(&flamingo.NullLogger{}, nil, nil)
	userservice := &authApplication.UserService{}
	userservice.Inject(authmanager, nil)
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
						authmanager,
						userservice,
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
						authmanager,
						userservice,
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
						authmanager,
						userservice,
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
				authmanager,
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
	authmanager := &authApplication.AuthManager{}
	authmanager.Inject(&flamingo.NullLogger{}, nil, nil)
	userservice := &authApplication.UserService{}
	userservice.Inject(authmanager, nil)
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
						authmanager,
						userservice,
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
						authmanager,
						userservice,
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
				authmanager,
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
	authmanager := &authApplication.AuthManager{}
	authmanager.Inject(&flamingo.NullLogger{}, nil, nil)
	userservice := &authApplication.UserService{}
	userservice.Inject(authmanager, nil)
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
						authmanager,
						userservice,
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
						authmanager,
						userservice,
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
				authmanager,
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
	MockGuestCartServiceWithModifyBehaviour struct{}
)

func (m *MockGuestCartServiceWithModifyBehaviour) GetCart(ctx context.Context, cartID string) (*cartDomain.Cart, error) {
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
	cob := &infrastructure.InMemoryBehaviour{}

	storage := &infrastructure.InMemoryCartStorage{}
	cart := cartDomain.Cart{
		ID: "mock_guest_cart",
		Deliveries: []cartDomain.Delivery{
			cartDomain.Delivery{
				DeliveryInfo: cartDomain.DeliveryInfo{Code: "default_delivery_code"},
				Cartitems: []cartDomain.Item{
					cartDomain.Item{
						ID:  "mock_item",
						Qty: 7,
					},
				},
			},
		},
	}

	_ = storage.StoreCart(&cart)

	cob.Inject(
		storage,
		&MockProductService{},
		flamingo.NullLogger{},
		func() *cartDomain.ItemBuilder {
			return &cartDomain.ItemBuilder{}
		},
		func() *cartDomain.DeliveryBuilder {
			return &cartDomain.DeliveryBuilder{}
		},
		func() *cartDomain.Builder {
			return &cartDomain.Builder{}
		},
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
			cartDomain.Delivery{
				DeliveryInfo: cartDomain.DeliveryInfo{Code: "default_delivery_code"},
				Cartitems: []cartDomain.Item{
					cartDomain.Item{
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
			cartDomain.Delivery{
				DeliveryInfo: cartDomain.DeliveryInfo{Code: "default_delivery_code"},
				Cartitems: []cartDomain.Item{
					cartDomain.Item{
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
	return fmt.Sprintf("MockRestrictor")
}

func (r *MockRestrictor) Restrict(ctx context.Context, product productDomain.BasicProduct, currentCart *cartDomain.Cart, deliveryCode string) *validation.RestrictionResult {
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

func (mpos *MockPlaceOrderService) PlaceCustomerCart(ctx context.Context, auth oauth.Auth, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	return nil, nil
}

func (mpos *MockPlaceOrderService) ReserveOrderID(ctx context.Context, cart *cartDomain.Cart) (string, error) {
	return "201910251128788TD", nil
}

func (mpos *MockPlaceOrderService) CancelGuestOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos) error {
	return nil
}

func (mpos *MockPlaceOrderService) CancelCustomerOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos, auth oauth.Auth) error {
	return nil
}

func TestCartService_CartInEvent(t *testing.T) {
	// bootstrap cart service
	authmanager := &authApplication.AuthManager{}
	authmanager.Inject(&flamingo.NullLogger{}, nil, nil)
	userservice := &authApplication.UserService{}
	userservice.Inject(authmanager, nil)
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
			authmanager,
			userservice,
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
		authmanager,
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
