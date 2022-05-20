package application_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"flamingo.me/flamingo/v3/core/auth"
	authMock "flamingo.me/flamingo/v3/core/auth/mock"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	cartInfrastructure "flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// MockGuestCartServiceAdapter
	MockGuestCartServiceAdapter struct {
		Behaviour cartDomain.ModifyBehaviour
	}
)

var (
	// test interface implementation
	_ cartDomain.GuestCartService = (*MockGuestCartServiceAdapter)(nil)
)

func (m *MockGuestCartServiceAdapter) GetCart(ctx context.Context, cartID string) (*cartDomain.Cart, error) {
	return &cartDomain.Cart{
		ID: "mock_guest_cart",
	}, nil
}

func (m *MockGuestCartServiceAdapter) GetNewCart(ctx context.Context) (*cartDomain.Cart, error) {
	return &cartDomain.Cart{
		ID: "mock_guest_cart",
	}, nil
}

func (m *MockGuestCartServiceAdapter) GetModifyBehaviour(context.Context) (cartDomain.ModifyBehaviour, error) {
	return m.Behaviour, nil
}

func (m *MockGuestCartServiceAdapter) RestoreCart(ctx context.Context, cart cartDomain.Cart) (*cartDomain.Cart, error) {
	restoredCart := cart
	restoredCart.ID = "1111"
	return &restoredCart, nil
}

var (
	// test interface implementation
	_ cartDomain.GuestCartService = (*MockGuestCartServiceAdapterError)(nil)
)

type (
	// MockGuestCartServiceAdapter with error on GetCart
	MockGuestCartServiceAdapterError struct{}
)

func (m *MockGuestCartServiceAdapterError) GetCart(ctx context.Context, cartID string) (*cartDomain.Cart, error) {
	return nil, errors.New("defective")
}

func (m *MockGuestCartServiceAdapterError) GetNewCart(ctx context.Context) (*cartDomain.Cart, error) {
	return nil, errors.New("defective")
}

func (m *MockGuestCartServiceAdapterError) GetModifyBehaviour(context.Context) (cartDomain.ModifyBehaviour, error) {
	// get guest behaviour
	return new(cartInfrastructure.DefaultCartBehaviour), nil
}

func (m *MockGuestCartServiceAdapterError) RestoreCart(ctx context.Context, cart cartDomain.Cart) (*cartDomain.Cart, error) {
	return nil, errors.New("defective")
}

// MockCustomerCartService

type (
	MockCustomerCartService struct {
		Behaviour cartDomain.ModifyBehaviour
	}
)

var (
	// test interface implementation
	_ cartDomain.CustomerCartService = (*MockCustomerCartService)(nil)
)

func (m *MockCustomerCartService) GetModifyBehaviour(context.Context, auth.Identity) (cartDomain.ModifyBehaviour, error) {
	// customer behaviour
	return m.Behaviour, nil
}

func (m *MockCustomerCartService) GetCart(ctx context.Context, identity auth.Identity, cartID string) (*cartDomain.Cart, error) {
	return &cartDomain.Cart{
		ID: "mock_customer_cart",
	}, nil
}

func (m *MockCustomerCartService) RestoreCart(ctx context.Context, identity auth.Identity, cart cartDomain.Cart) (*cartDomain.Cart, error) {
	panic("implement me")
}

// MockProductService

type (
	MockProductService struct{}
)

func (m *MockProductService) Get(ctx context.Context, marketplaceCode string) (productDomain.BasicProduct, error) {
	mockProduct := new(productDomain.SimpleProduct)

	mockProduct.Identifier = "mock_product"

	return mockProduct, nil
}

// MockCartCache

type (
	MockCartCache struct {
		CachedCart *cartDomain.Cart
	}
)

func (m *MockCartCache) GetCart(context.Context, *web.Session, cartApplication.CartCacheIdentifier) (*cartDomain.Cart, error) {
	return m.CachedCart, nil
}

func (m *MockCartCache) CacheCart(ctx context.Context, s *web.Session, cci cartApplication.CartCacheIdentifier, cart *cartDomain.Cart) error {
	m.CachedCart = cart
	return nil
}

func (m *MockCartCache) Invalidate(context.Context, *web.Session, cartApplication.CartCacheIdentifier) error {
	return nil
}

func (m *MockCartCache) Delete(context.Context, *web.Session, cartApplication.CartCacheIdentifier) error {
	return nil
}

func (m *MockCartCache) DeleteAll(context.Context, *web.Session) error {
	return nil
}

func (m *MockCartCache) BuildIdentifier(context.Context, *web.Session) (cartApplication.CartCacheIdentifier, error) {
	return cartApplication.CartCacheIdentifier{}, nil
}

// MockEventPublisher

type (
	MockEventPublisher struct {
		mock.Mock
	}
)

var (
	_ events.EventPublisher = (*MockEventPublisher)(nil)
)

func (m *MockEventPublisher) PublishAddToCartEvent(ctx context.Context, cart *cartDomain.Cart, marketPlaceCode string, variantMarketPlaceCode string, qty int) {
	m.Called()
}

func (m *MockEventPublisher) PublishChangedQtyInCartEvent(ctx context.Context, cart *cartDomain.Cart, item *cartDomain.Item, qtyBefore int, qtyAfter int) {
	m.Called()
}

func (m *MockEventPublisher) PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, placedOrderInfos placeorder.PlacedOrderInfos) {
	m.Called()
}

type (
	MockDeliveryInfoBuilder struct{}
)

func (m *MockDeliveryInfoBuilder) BuildByDeliveryCode(deliveryCode string) (*cartDomain.DeliveryInfo, error) {
	return &cartDomain.DeliveryInfo{}, nil
}

type (
	MockEventRouter struct {
		mock.Mock
	}
)

var _ flamingo.EventRouter = new(MockEventRouter)

func (m *MockEventRouter) Dispatch(ctx context.Context, event flamingo.Event) {
	// we just write the event type and the marketplace code to the mock, so we don't have to compare
	// the complete cart
	switch eventType := event.(type) {
	case *events.AddToCartEvent:
		m.Called(ctx, fmt.Sprintf("%T", event), eventType.MarketplaceCode)
	}
}

func TestCartReceiverService_ShouldHaveGuestCart(t *testing.T) {
	type fields struct {
		GuestCartService     cartDomain.GuestCartService
		CustomerCartService  cartDomain.CustomerCartService
		CartDecoratorFactory *decorator.DecoratedCartFactory
		Logger               flamingo.Logger
		CartCache            cartApplication.CartCache
	}
	type args struct {
		session *web.Session
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "has session key",
			fields: fields{
				GuestCartService:    new(MockGuestCartServiceAdapter),
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, struct{}{}),
			},
			want: true,
		}, {
			name: "doesn't have session key",
			fields: fields{
				GuestCartService:    new(MockGuestCartServiceAdapter),
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				session: web.EmptySession().Store("arbitrary_and_wrong_key", struct{}{}),
			},

			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cartApplication.CartReceiverService{}
			cs.Inject(
				tt.fields.GuestCartService,
				tt.fields.CustomerCartService,
				tt.fields.CartDecoratorFactory,
				nil,
				tt.fields.Logger,
				nil,
				&struct {
					CartCache cartApplication.CartCache `inject:",optional"`
				}{
					CartCache: tt.fields.CartCache,
				},
			)

			got := cs.ShouldHaveGuestCart(tt.args.session)

			if got != tt.want {
				t.Errorf("CartReceiverService.ShouldHaveGuestCart() = %v, wantType0 %v", got, tt.want)
			}
		})
	}
}

func TestCartReceiverService_ViewGuestCart(t *testing.T) {
	type fields struct {
		GuestCartService     cartDomain.GuestCartService
		CustomerCartService  cartDomain.CustomerCartService
		CartDecoratorFactory *decorator.DecoratedCartFactory
		Logger               flamingo.Logger
		CartCache            cartApplication.CartCache
	}
	type args struct {
		ctx     context.Context
		session *web.Session
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *cartDomain.Cart
		wantErr        bool
		wantMessageErr string
	}{
		{
			name: "has no guest cart",
			fields: fields{
				GuestCartService:    new(MockGuestCartServiceAdapter),
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store("stuff", "some_malformed_id"),
			},
			want:           &cartDomain.Cart{},
			wantErr:        false,
			wantMessageErr: "",
		}, {
			name: "failed guest cart get",
			fields: fields{
				GuestCartService:    new(MockGuestCartServiceAdapterError),
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want:           nil,
			wantErr:        true,
			wantMessageErr: cartApplication.ErrTemporaryCartService.Error(),
		}, {
			name: "guest cart get without error",
			fields: fields{
				GuestCartService:    new(MockGuestCartServiceAdapter),
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store(cartApplication.GuestCartSessionKey, "some_guest_id"),
			},
			want: &cartDomain.Cart{
				ID: "mock_guest_cart",
			},
			wantErr:        false,
			wantMessageErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cartApplication.CartReceiverService{}
			cs.Inject(
				tt.fields.GuestCartService,
				tt.fields.CustomerCartService,
				tt.fields.CartDecoratorFactory,
				nil,
				tt.fields.Logger,
				nil,
				&struct {
					CartCache cartApplication.CartCache `inject:",optional"`
				}{
					CartCache: tt.fields.CartCache,
				},
			)

			got, err := cs.ViewGuestCart(tt.args.ctx, tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartReceiverService.ViewGuestCart() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CartReceiverService.ViewGuestCart() = %v, wantType0 %v", got, tt.want)
			}
		})
	}
}

func TestCartReceiverService_DecorateCart(t *testing.T) {
	type fields struct {
		GuestCartService     cartDomain.GuestCartService
		CustomerCartService  cartDomain.CustomerCartService
		CartDecoratorFactory *decorator.DecoratedCartFactory
		Logger               flamingo.Logger
		CartCache            cartApplication.CartCache
	}
	type args struct {
		ctx  context.Context
		cart *cartDomain.Cart
	}
	tests := []struct {
		name                     string
		fields                   fields
		args                     args
		want                     *decorator.DecoratedCart
		wantErr                  bool
		wantMessageErr           string
		wantDecoratedItemsLength int
	}{
		{
			name: "error b/c no cart supplied",
			fields: fields{
				GuestCartService:    new(MockGuestCartServiceAdapter),
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				ctx:  context.Background(),
				cart: nil,
			},
			wantErr:                  true,
			wantMessageErr:           "no cart given",
			wantDecoratedItemsLength: 0,
		}, {
			name: "basic decoration of cart",
			fields: fields{
				GuestCartService:    new(MockGuestCartServiceAdapter),
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				ctx: context.Background(),
				cart: &cartDomain.Cart{
					ID: "some_test_cart",
					Deliveries: []cartDomain.Delivery{
						{
							Cartitems: []cartDomain.Item{
								{
									ID: "test_id",
								},
							},
						},
					},
				},
			},
			wantErr:                  false,
			wantMessageErr:           "",
			wantDecoratedItemsLength: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cartApplication.CartReceiverService{}
			cs.Inject(
				tt.fields.GuestCartService,
				tt.fields.CustomerCartService,
				tt.fields.CartDecoratorFactory,
				nil,
				tt.fields.Logger,
				nil,
				&struct {
					CartCache cartApplication.CartCache `inject:",optional"`
				}{
					CartCache: tt.fields.CartCache,
				},
			)

			got, err := cs.DecorateCart(tt.args.ctx, tt.args.cart)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartReceiverService.DecorateCart() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)

				return
			}

			if tt.wantDecoratedItemsLength > 0 {
				for _, decoratedDeliveryItem := range got.DecoratedDeliveries {
					if len(decoratedDeliveryItem.DecoratedItems) != tt.wantDecoratedItemsLength {
						t.Errorf("Mismatch of expected Decorated Items, got %d, expected %d", len(decoratedDeliveryItem.DecoratedItems), tt.wantDecoratedItemsLength)
					}
				}
			}
		})
	}
}

func TestCartReceiverService_GetDecoratedCart(t *testing.T) {
	type fields struct {
		GuestCartService     cartDomain.GuestCartService
		CustomerCartService  cartDomain.CustomerCartService
		CartDecoratorFactory *decorator.DecoratedCartFactory
		Logger               flamingo.Logger
		CartCache            cartApplication.CartCache
	}
	type args struct {
		ctx     context.Context
		session *web.Session
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          *decorator.DecoratedCart
		wantBehaviour cartDomain.ModifyBehaviour
		wantErr       bool
	}{
		{
			name: "decorated cart not found",
			fields: fields{
				GuestCartService:    new(MockGuestCartServiceAdapterError),
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store("some_nonvalid_key", "some_guest_id"),
			},
			want:          nil,
			wantBehaviour: nil,
			wantErr:       true,
		}, {
			name: "decorated cart found",
			fields: fields{
				GuestCartService:    &MockGuestCartServiceAdapter{Behaviour: &cartInfrastructure.DefaultCartBehaviour{}},
				CustomerCartService: new(MockCustomerCartService),
				CartDecoratorFactory: func() *decorator.DecoratedCartFactory {
					result := &decorator.DecoratedCartFactory{}
					result.Inject(
						&MockProductService{},
						flamingo.NullLogger{},
					)

					return result
				}(),
				Logger:    flamingo.NullLogger{},
				CartCache: new(MockCartCache),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store("some_nonvalid_key", "some_guest_id"),
			},
			want:          &decorator.DecoratedCart{},
			wantBehaviour: &cartInfrastructure.DefaultCartBehaviour{},
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cartApplication.CartReceiverService{}
			cs.Inject(
				tt.fields.GuestCartService,
				tt.fields.CustomerCartService,
				tt.fields.CartDecoratorFactory,
				nil,
				tt.fields.Logger,
				nil,
				&struct {
					CartCache cartApplication.CartCache `inject:",optional"`
				}{
					CartCache: tt.fields.CartCache,
				},
			)

			got, gotBehaviour, err := cs.GetDecoratedCart(tt.args.ctx, tt.args.session)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartReceiverService.GetDecoratedCart() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want == nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("CartReceiverService.GetDecoratedCart() got = %v, want %v", got, tt.want)

					return
				}
			} else {
				gotType := reflect.TypeOf(got).Elem()
				wantType := reflect.TypeOf(tt.want).Elem()
				if wantType != gotType {
					t.Error("Return Type for want doesn't match")

					return
				}
			}

			if tt.wantBehaviour == nil {
				if !reflect.DeepEqual(gotBehaviour, tt.wantBehaviour) {
					t.Errorf("CartReceiverService.GetDecoratedCart() gotBehaviour = %v, wantBehaviour %v", gotBehaviour, tt.wantBehaviour)

					return
				}
			} else {
				gotType1 := reflect.TypeOf(gotBehaviour).Elem()
				wantType1 := reflect.TypeOf(tt.wantBehaviour).Elem()
				if wantType1 != gotType1 {
					t.Errorf("Return Type for want doesn't match, got = %v, want = %v", gotType1, wantType1)

					return
				}
			}

			if !reflect.DeepEqual(gotBehaviour, tt.wantBehaviour) {
				t.Errorf("CartReceiverService.GetDecoratedCart() gotBehaviour = %v, wantBehaviour %v", gotBehaviour, tt.wantBehaviour)

				return
			}
		})
	}
}

func TestCartReceiverService_RestoreCart(t *testing.T) {
	type fields struct {
		guestCartService     cartDomain.GuestCartService
		customerCartService  cartDomain.CustomerCartService
		cartDecoratorFactory *decorator.DecoratedCartFactory
		logger               flamingo.Logger
		cartCache            cartApplication.CartCache
	}
	type args struct {
		ctx           context.Context
		session       *web.Session
		cartToRestore cartDomain.Cart
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		want                   *cartDomain.Cart
		wantErr                bool
		wantGuestCartSessionID bool
		wantCartStoredInCache  bool
	}{
		{
			name: "restore guest cart without error",
			fields: fields{
				guestCartService: &MockGuestCartServiceAdapter{},
				logger:           &flamingo.NullLogger{},
				cartCache:        &MockCartCache{},
			},
			args: args{
				ctx:     web.ContextWithSession(context.Background(), web.EmptySession()),
				session: web.EmptySession(),
				cartToRestore: cartDomain.Cart{
					ID: "1234",
					BillingAddress: &cartDomain.Address{
						Firstname: "Test",
						Lastname:  "Test",
						Email:     "test@test.xy",
					},
					Deliveries: []cartDomain.Delivery{
						{
							DeliveryInfo: cartDomain.DeliveryInfo{
								Code:     "pickup",
								Workflow: "pickup",
							},
							Cartitems: []cartDomain.Item{
								{
									ID:                     "1",
									ExternalReference:      "sku-1",
									MarketplaceCode:        "sku-1",
									VariantMarketPlaceCode: "",
									ProductName:            "Product #1",
									SourceID:               "",
									Qty:                    2,
								},
							},
							ShippingItem: cartDomain.ShippingItem{},
						},
					},
				},
			},
			want: &cartDomain.Cart{
				ID: "1111",
				BillingAddress: &cartDomain.Address{
					Firstname: "Test",
					Lastname:  "Test",
					Email:     "test@test.xy",
				},
				Deliveries: []cartDomain.Delivery{
					{
						DeliveryInfo: cartDomain.DeliveryInfo{
							Code:     "pickup",
							Workflow: "pickup",
						},
						Cartitems: []cartDomain.Item{
							{
								ID:                     "1",
								ExternalReference:      "sku-1",
								MarketplaceCode:        "sku-1",
								VariantMarketPlaceCode: "",
								ProductName:            "Product #1",
								SourceID:               "",
								Qty:                    2,
							},
						},
						ShippingItem: cartDomain.ShippingItem{},
					},
				},
			},
			wantErr:                false,
			wantGuestCartSessionID: true,
			wantCartStoredInCache:  true,
		},
		{
			name: "restore guest cart with error",
			fields: fields{
				guestCartService: &MockGuestCartServiceAdapterError{},

				logger:    &flamingo.NullLogger{},
				cartCache: &MockCartCache{},
			},
			args: args{
				ctx:     web.ContextWithSession(context.Background(), web.EmptySession()),
				session: web.EmptySession(),
				cartToRestore: cartDomain.Cart{
					ID: "1234",
					BillingAddress: &cartDomain.Address{
						Firstname: "Test",
						Lastname:  "Test",
						Email:     "test@test.xy",
					},
					Deliveries: []cartDomain.Delivery{
						{
							DeliveryInfo: cartDomain.DeliveryInfo{
								Code:     "pickup",
								Workflow: "pickup",
							},
							Cartitems: []cartDomain.Item{
								{
									ID:                     "1",
									ExternalReference:      "sku-1",
									MarketplaceCode:        "sku-1",
									VariantMarketPlaceCode: "",
									ProductName:            "Product #1",
									SourceID:               "",
									Qty:                    2,
								},
							},
							ShippingItem: cartDomain.ShippingItem{},
						},
					},
				},
			},
			want:                   nil,
			wantErr:                true,
			wantGuestCartSessionID: false,
			wantCartStoredInCache:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cartApplication.CartReceiverService{}
			cs.Inject(
				tt.fields.guestCartService,
				tt.fields.customerCartService,
				tt.fields.cartDecoratorFactory,
				nil,
				tt.fields.logger,
				nil,
				&struct {
					CartCache cartApplication.CartCache `inject:",optional"`
				}{
					CartCache: tt.fields.cartCache,
				},
			)
			got, err := cs.RestoreCart(tt.args.ctx, tt.args.session, tt.args.cartToRestore)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestoreCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RestoreCart() got = %v, want %v", got, tt.want)
			}

			sessionGot, found := tt.args.session.Load(cartApplication.GuestCartSessionKey)
			if found != tt.wantGuestCartSessionID {
				t.Error("GuestCartID not found in session")
			}

			if found == true && tt.want != nil {
				if !reflect.DeepEqual(tt.want.ID, sessionGot) {
					t.Errorf("GuestCartID in session does not match restored cart got = %v, want %v", got, tt.wantGuestCartSessionID)
				}
			}

			if tt.wantCartStoredInCache && tt.fields.cartCache != nil {
				cart, _ := tt.fields.cartCache.GetCart(context.Background(), nil, cartApplication.CartCacheIdentifier{})
				if cart == nil {
					t.Error("Cart not found in cart cache")
				}
			}
		})
	}
}

func TestCartReceiverService_ModifyBehaviour(t *testing.T) {
	t.Run("get guest behaviour", func(t *testing.T) {
		mockBehaviour := &cartInfrastructure.DefaultCartBehaviour{}
		cs := &cartApplication.CartReceiverService{}
		cs.Inject(
			&MockGuestCartServiceAdapter{Behaviour: mockBehaviour},
			&MockCustomerCartService{},
			&decorator.DecoratedCartFactory{},
			&auth.WebIdentityService{},
			flamingo.NullLogger{},
			nil,
			&struct {
				CartCache cartApplication.CartCache `inject:",optional"`
			}{
				CartCache: &MockCartCache{},
			},
		)

		behaviour, err := cs.ModifyBehaviour(context.Background())

		assert.NoError(t, err)
		require.NotNil(t, behaviour)
		assert.Same(t, behaviour, mockBehaviour)
	})

	t.Run("get customer behaviour", func(t *testing.T) {
		mockIdentifier := new(authMock.Identifier).SetIdentifyMethod(
			func(identifier *authMock.Identifier, ctx context.Context, request *web.Request) (auth.Identity, error) {
				return &authMock.Identity{Sub: "foo"}, nil
			},
		)
		cs := &cartApplication.CartReceiverService{}
		mockBehaviour := &cartInfrastructure.DefaultCartBehaviour{}
		cs.Inject(
			&MockGuestCartServiceAdapter{},
			&MockCustomerCartService{Behaviour: mockBehaviour},
			&decorator.DecoratedCartFactory{},
			new(auth.WebIdentityService).Inject([]auth.RequestIdentifier{mockIdentifier}, nil, nil, nil),
			flamingo.NullLogger{},
			nil,
			&struct {
				CartCache cartApplication.CartCache `inject:",optional"`
			}{
				CartCache: &MockCartCache{},
			},
		)

		behaviour, err := cs.ModifyBehaviour(context.Background())

		assert.NoError(t, err)
		require.NotNil(t, behaviour)
		assert.Same(t, behaviour, mockBehaviour)
	})
}
