package application_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"

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
	MockGuestCartServiceAdapter struct{}
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
	return new(cartInfrastructure.InMemoryBehaviour), nil
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
	return new(cartInfrastructure.InMemoryBehaviour), nil
}

func (m *MockGuestCartServiceAdapterError) RestoreCart(ctx context.Context, cart cartDomain.Cart) (*cartDomain.Cart, error) {
	return nil, errors.New("defective")
}

// MockCustomerCartService

type (
	MockCustomerCartService struct{}
)

var (
	// test interface implementation
	_ cartDomain.CustomerCartService = (*MockCustomerCartService)(nil)
)

func (m *MockCustomerCartService) GetModifyBehaviour(context.Context, domain.Auth) (cartDomain.ModifyBehaviour, error) {
	return new(cartInfrastructure.InMemoryBehaviour), nil
}

func (m *MockCustomerCartService) GetCart(ctx context.Context, auth domain.Auth, cartID string) (*cartDomain.Cart, error) {
	return &cartDomain.Cart{
		ID: "mock_customer_cart",
	}, nil
}

func (m *MockCustomerCartService) RestoreCart(ctx context.Context, auth domain.Auth, cart cartDomain.Cart) (*cartDomain.Cart, error) {
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
	MockUserService struct {
		LoggedIn bool
	}
)

var _ authApplication.UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) GetUser(ctx context.Context, session *web.Session) *domain.User {
	return &domain.User{
		Name: "Test",
	}
}

func (m *MockUserService) IsLoggedIn(ctx context.Context, session *web.Session) bool {
	return m.LoggedIn
}

type (
	MockAuthManager struct {
		ShouldReturnError bool
	}
	MockTokenSource struct {
	}
)

var _ cartApplication.AuthManagerInterface = &MockAuthManager{}
var _ oauth2.TokenSource = &MockTokenSource{}

func (m MockTokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

func (m MockAuthManager) Auth(_ context.Context, _ *web.Session) (domain.Auth, error) {
	if m.ShouldReturnError {
		return domain.Auth{}, errors.New("generic auth error")
	}

	return domain.Auth{
		TokenSource: &MockTokenSource{},
		IDToken:     &oidc.IDToken{},
	}, nil
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
		AuthManager          *authApplication.AuthManager
		UserService          *authApplication.UserService
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
				AuthManager: &authApplication.AuthManager{},
				UserService: &authApplication.UserService{},
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
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
				AuthManager: &authApplication.AuthManager{},
				UserService: &authApplication.UserService{},
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
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
				tt.fields.AuthManager,
				tt.fields.UserService,
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
		AuthManager          *authApplication.AuthManager
		UserService          *authApplication.UserService
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
				AuthManager: &authApplication.AuthManager{},
				UserService: &authApplication.UserService{},
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
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
				AuthManager: &authApplication.AuthManager{},
				UserService: &authApplication.UserService{},
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
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
				AuthManager: &authApplication.AuthManager{},
				UserService: &authApplication.UserService{},
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
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
				tt.fields.AuthManager,
				tt.fields.UserService,
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
		AuthManager          *authApplication.AuthManager
		UserService          *authApplication.UserService
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
				AuthManager: &authApplication.AuthManager{},
				UserService: &authApplication.UserService{},
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
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
				AuthManager: &authApplication.AuthManager{},
				UserService: &authApplication.UserService{},
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
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
				tt.fields.AuthManager,
				tt.fields.UserService,
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
	authmanager := &authApplication.AuthManager{}
	authmanager.Inject(flamingo.NullLogger{}, nil, nil)
	userservice := &authApplication.UserService{}
	userservice.Inject(authmanager, nil)
	type fields struct {
		GuestCartService     cartDomain.GuestCartService
		CustomerCartService  cartDomain.CustomerCartService
		CartDecoratorFactory *decorator.DecoratedCartFactory
		AuthManager          *authApplication.AuthManager
		UserService          *authApplication.UserService
		Logger               flamingo.Logger
		CartCache            cartApplication.CartCache
	}
	type args struct {
		ctx     context.Context
		session *web.Session
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantType0 *decorator.DecoratedCart
		wantType1 cartDomain.ModifyBehaviour
		wantErr   bool
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
				AuthManager: authmanager,
				UserService: userservice,
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store("some_nonvalid_key", "some_guest_id"),
			},
			wantType0: nil,
			wantType1: nil,
			wantErr:   true,
		}, {
			name: "decorated cart found",
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
				AuthManager: authmanager,
				UserService: userservice,
				Logger:      flamingo.NullLogger{},
				CartCache:   new(MockCartCache),
			},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession().Store("some_nonvalid_key", "some_guest_id"),
			},
			wantType0: &decorator.DecoratedCart{},
			wantType1: &cartInfrastructure.InMemoryBehaviour{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &cartApplication.CartReceiverService{}
			cs.Inject(
				tt.fields.GuestCartService,
				tt.fields.CustomerCartService,
				tt.fields.CartDecoratorFactory,
				tt.fields.AuthManager,
				tt.fields.UserService,
				tt.fields.Logger,
				nil,
				&struct {
					CartCache cartApplication.CartCache `inject:",optional"`
				}{
					CartCache: tt.fields.CartCache,
				},
			)

			got, got1, err := cs.GetDecoratedCart(tt.args.ctx, tt.args.session)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartReceiverService.GetDecoratedCart() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantType0 == nil {
				if !reflect.DeepEqual(got, tt.wantType0) {
					t.Errorf("CartReceiverService.GetDecoratedCart() got = %v, wantType0 %v", got, tt.wantType0)

					return
				}
			} else {
				gotType := reflect.TypeOf(got).Elem()
				wantType := reflect.TypeOf(tt.wantType0).Elem()
				if wantType != gotType {
					t.Error("Return Type for wantType0 doesn't match")

					return
				}
			}

			if tt.wantType1 == nil {
				if !reflect.DeepEqual(got1, tt.wantType1) {
					t.Errorf("CartReceiverService.GetDecoratedCart() got = %v, wantType0 %v", got1, tt.wantType1)

					return
				}
			} else {
				gotType1 := reflect.TypeOf(got1).Elem()
				wantType1 := reflect.TypeOf(tt.wantType1).Elem()
				if wantType1 != gotType1 {
					t.Errorf("Return Type for wantType0 doesn't match, got = %v, want = %v", gotType1, wantType1)

					return
				}
			}

			if !reflect.DeepEqual(got1, tt.wantType1) {
				t.Errorf("CartReceiverService.GetDecoratedCart() got1 = %v, wantType0 %v", got1, tt.wantType1)

				return
			}
		})
	}
}

func TestCartReceiverService_RestoreCart(t *testing.T) {
	authmanager := &authApplication.AuthManager{}
	authmanager.Inject(flamingo.NullLogger{}, nil, nil)
	userservice := &authApplication.UserService{}
	userservice.Inject(authmanager, nil)
	type fields struct {
		guestCartService     cartDomain.GuestCartService
		customerCartService  cartDomain.CustomerCartService
		cartDecoratorFactory *decorator.DecoratedCartFactory
		authManager          *authApplication.AuthManager
		userService          authApplication.UserServiceInterface
		eventRouter          flamingo.EventRouter
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
				userService:      userservice,
				logger:           &flamingo.NullLogger{},
				cartCache:        &MockCartCache{},
				authManager:      authmanager,
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
				userService:      userservice,
				logger:           &flamingo.NullLogger{},
				cartCache:        &MockCartCache{},
				authManager:      authmanager,
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
				tt.fields.authManager,
				tt.fields.userService,
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
				cart, _ := tt.fields.cartCache.GetCart(nil, nil, cartApplication.CartCacheIdentifier{})
				if cart == nil {
					t.Error("Cart not found in cart cache")
				}
			}
		})
	}
}

func TestCartReceiverService_ModifyBehaviour(t *testing.T) {
	t.Run("get guest behaviour", func(t *testing.T) {
		cs := &cartApplication.CartReceiverService{}
		cs.Inject(
			&MockGuestCartServiceAdapter{},
			&MockCustomerCartService{},
			&decorator.DecoratedCartFactory{},
			&MockAuthManager{},
			&MockUserService{LoggedIn: false},
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
		assert.IsType(t, behaviour, &cartInfrastructure.InMemoryBehaviour{})
	})

	t.Run("get customer behaviour", func(t *testing.T) {
		cs := &cartApplication.CartReceiverService{}
		cs.Inject(
			&MockGuestCartServiceAdapter{},
			&MockCustomerCartService{},
			&decorator.DecoratedCartFactory{},
			&MockAuthManager{},
			&MockUserService{LoggedIn: true},
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
		assert.IsType(t, behaviour, &cartInfrastructure.InMemoryBehaviour{})
	})

	t.Run("error during customer auth should lead to error", func(t *testing.T) {
		cs := &cartApplication.CartReceiverService{}
		cs.Inject(
			&MockGuestCartServiceAdapter{},
			&MockCustomerCartService{},
			&decorator.DecoratedCartFactory{},
			&MockAuthManager{ShouldReturnError: true},
			&MockUserService{LoggedIn: true},
			flamingo.NullLogger{},
			nil,
			&struct {
				CartCache cartApplication.CartCache `inject:",optional"`
			}{
				CartCache: &MockCartCache{},
			},
		)

		behaviour, err := cs.ModifyBehaviour(context.Background())

		assert.Error(t, err)
		assert.Nil(t, behaviour)
	})
}
