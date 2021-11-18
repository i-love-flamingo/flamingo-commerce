package application_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"flamingo.me/flamingo-commerce/v3/price/domain"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

var (
	_ application.CartCache = new(application.CartSessionCache)
)

func TestCartCacheIdentifier_CacheKey(t *testing.T) {
	type fields struct {
		GuestCartID    string
		IsCustomerCart bool
		CustomerID     string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic test",
			fields: fields{
				GuestCartID:    "guest_cart_id",
				IsCustomerCart: false,
				CustomerID:     "customer_id",
			},
			want: "cart_customer_id_guest_cart_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := &application.CartCacheIdentifier{
				GuestCartID:    tt.fields.GuestCartID,
				IsCustomerCart: tt.fields.IsCustomerCart,
				CustomerID:     tt.fields.CustomerID,
			}

			if got := ci.CacheKey(); got != tt.want {
				t.Errorf("CartCacheIdentifier.CacheKey() = %v, wantType0 %v", got, tt.want)
			}
		})
	}
}

func TestCartSessionCache_GetCart(t *testing.T) {
	type args struct {
		ctx     context.Context
		session *web.Session
		id      application.CartCacheIdentifier
	}
	tests := []struct {
		name           string
		args           args
		want           *cart.Cart
		wantErr        bool
		wantMessageErr string
	}{
		{
			name: "error for no cart in cache",
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession(),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           nil,
			wantErr:        true,
			wantMessageErr: "cache entry not found",
		},
		{
			name: "cached cart found/invalid cache entry",
			args: args{
				ctx: context.Background(),
				session: web.EmptySession().Store(
					application.CartSessionCacheCacheKeyPrefix+"cart_customer_id_guest_cart_id",
					application.CachedCartEntry{IsInvalid: true},
				),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           new(cart.Cart),
			wantErr:        true,
			wantMessageErr: "cache is invalid",
		},
		{
			name: "cached cart found/valid cache entry",
			args: args{
				ctx: context.Background(),
				session: web.EmptySession().Store(
					application.CartSessionCacheCacheKeyPrefix+"cart_customer_id_guest_cart_id",
					application.CachedCartEntry{
						IsInvalid: false,
						ExpiresOn: time.Now().Add(5 * time.Minute),
						Entry:     *getFixtureCartForTest(t),
					},
				),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           getFixtureCartForTest(t),
			wantErr:        false,
			wantMessageErr: "",
		},
		{
			name: "session contains invalid data at cache key",
			args: args{
				ctx: context.Background(),
				session: web.EmptySession().Store(
					application.CartSessionCacheCacheKeyPrefix+"cart_customer_id_guest_cart_id",
					struct {
						invalidProperty bool
					}{
						invalidProperty: true,
					},
				),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           nil,
			wantErr:        true,
			wantMessageErr: "cart cache contains invalid data at cache key",
		},
		{
			name: "session contains expired cart cache",
			args: args{
				ctx: context.Background(),
				session: web.EmptySession().Store(application.CartSessionCacheCacheKeyPrefix+"cart_customer_id_guest_cart_id",
					application.CachedCartEntry{
						IsInvalid: false,
						ExpiresOn: time.Now().Add(-1 * time.Second),
					},
				),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           nil,
			wantErr:        true,
			wantMessageErr: "cache is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &application.CartSessionCache{}
			c.Inject(flamingo.NullLogger{}, nil, nil)

			got, err := c.GetCart(tt.args.ctx, tt.args.session, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartSessionCache.GetCart() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CartSessionCache.GetCart() = %v, wantType0 %v", got, tt.want)
			}
		})
	}
}

func TestCartSessionCache_CacheCart(t *testing.T) {
	type fields struct {
		config *struct {
			LifetimeSeconds float64 `inject:"config:commerce.cart.cacheLifetime"` // in seconds
		}
	}
	type args struct {
		ctx          context.Context
		session      *web.Session
		id           application.CartCacheIdentifier
		cartForCache *cart.Cart
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantMessageErr string
	}{
		{
			name: "no cart given",
			fields: fields{
				config: &struct {
					LifetimeSeconds float64 `inject:"config:commerce.cart.cacheLifetime"` // in seconds
				}{
					LifetimeSeconds: 300,
				},
			},
			args: args{
				ctx:     context.Background(),
				session: nil,
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "",
				},
				cartForCache: nil,
			},
			wantErr:        true,
			wantMessageErr: "no cart given to cache",
		},
		{
			name:   "cart is cached",
			fields: fields{},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession(),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "",
				},
				cartForCache: new(cart.Cart),
			},
			wantErr:        false,
			wantMessageErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &application.CartSessionCache{}
			c.Inject(flamingo.NullLogger{}, nil, nil)

			err := c.CacheCart(tt.args.ctx, tt.args.session, tt.args.id, tt.args.cartForCache)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartSessionCache.CacheCart() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)
			}
		})
	}
}

func TestCartSessionCache_CartExpiry(t *testing.T) {
	c := &application.CartSessionCache{}
	c.Inject(
		flamingo.NullLogger{},
		nil,
		&struct {
			LifetimeSeconds float64 `inject:"config:commerce.cart.cacheLifetime"` // in seconds
		}{
			LifetimeSeconds: 1,
		},
	)

	ctx := context.Background()
	session := web.EmptySession()

	id := application.CartCacheIdentifier{
		GuestCartID:    "guest_cart_id",
		IsCustomerCart: false,
		CustomerID:     "",
	}

	cartForCache := new(cart.Cart)

	err := c.CacheCart(ctx, session, id, cartForCache)
	if err != nil {
		t.Errorf("Failed to cache the cart. Error: %v", err)
	}

	time.Sleep(1 * time.Second)

	_, err = c.GetCart(context.Background(), session, id)
	if err != application.ErrCacheIsInvalid {
		t.Error("Expected cache invalid error. Received nil")
		return
	}
}

func TestCartSessionCache_Invalidate(t *testing.T) {
	type fields struct {
	}
	type args struct {
		ctx     context.Context
		session *web.Session
		id      application.CartCacheIdentifier
	}
	tests := []struct {
		name                  string
		fields                fields
		args                  args
		wantErr               bool
		wantMessageErr        string
		wantCacheEntryInvalid bool
	}{
		{
			name:   "no cache entry",
			fields: fields{},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession(),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "",
				},
			},
			wantErr:               true,
			wantMessageErr:        "cache entry not found",
			wantCacheEntryInvalid: false,
		},
		{
			name:   "invalidate cache entry",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				session: web.EmptySession().Store(
					application.CartSessionCacheCacheKeyPrefix+"cart__guest_cart_id",
					application.CachedCartEntry{IsInvalid: false},
				),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "",
				},
			},
			wantErr:               false,
			wantMessageErr:        "",
			wantCacheEntryInvalid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &application.CartSessionCache{}
			c.Inject(flamingo.NullLogger{}, nil, nil)

			err := c.Invalidate(tt.args.ctx, tt.args.session, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartSessionCache.Invalidate() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)
			}

			if tt.wantCacheEntryInvalid == true {
				for key := range tt.args.session.Keys() {
					if cacheEntry, ok := tt.args.session.Load(key); ok {
						if entry, ok := cacheEntry.(application.CachedCartEntry); ok {
							if entry.IsInvalid != tt.wantCacheEntryInvalid {
								t.Errorf("Cache validity doesnt match - got %v, wantCacheEntryInvalid %v", entry.IsInvalid, tt.wantCacheEntryInvalid)
							}
						}
					}
				}
			}
		})
	}
}

func TestCartSessionCache_Delete(t *testing.T) {
	type fields struct {
	}
	type args struct {
		ctx     context.Context
		session *web.Session
		id      application.CartCacheIdentifier
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        bool
		wantMessageErr string
	}{
		{
			name:   "cache entry not found for delete",
			fields: fields{},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession(),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "",
				},
			},
			wantErr:        true,
			wantMessageErr: "cache entry not found",
		},
		{
			name:   "deleted correctly",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				session: web.EmptySession().Store(
					application.CartSessionCacheCacheKeyPrefix+"cart__guest_cart_id",
					application.CachedCartEntry{IsInvalid: false},
				),
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "",
				},
			},
			wantErr:        false,
			wantMessageErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &application.CartSessionCache{}
			c.Inject(flamingo.NullLogger{}, nil, nil)

			err := c.Delete(tt.args.ctx, tt.args.session, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartSessionCache.Delete() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)
			}
		})
	}
}

func TestCartSessionCache_DeleteAll(t *testing.T) {
	type fields struct {
	}
	type args struct {
		ctx     context.Context
		session *web.Session
	}
	tests := []struct {
		name                   string
		fields                 fields
		args                   args
		wantErr                bool
		wantMessageErr         string
		wantSessionValuesEmpty bool
	}{
		{
			name:   "no cachekey found/nothing deleted",
			fields: fields{},
			args: args{
				ctx:     context.Background(),
				session: web.EmptySession(),
			},
			wantErr:                true,
			wantMessageErr:         "cache entry not found",
			wantSessionValuesEmpty: false,
		},
		{
			name:   "deleted an entry",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				session: web.EmptySession().Store(
					application.CartSessionCacheCacheKeyPrefix+"cart__guest_cart_id",
					application.CachedCartEntry{IsInvalid: false},
				),
			},
			wantErr:                false,
			wantMessageErr:         "",
			wantSessionValuesEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &application.CartSessionCache{}
			c.Inject(flamingo.NullLogger{}, nil, nil)

			err := c.DeleteAll(tt.args.ctx, tt.args.session)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartSessionCache.DeleteAll() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)
			}

			if tt.wantSessionValuesEmpty == true {
				if len(tt.args.session.Keys()) > 0 {
					t.Error("Session Values should have been emptied, but aren't")
				}
			}
		})
	}
}

func getFixtureCartForTest(t *testing.T) *cart.Cart {
	t.Helper()
	item := cart.Item{
		ID:                "1",
		ExternalReference: "1ext",
		Qty:               2,
		SinglePriceGross:  domain.NewFromInt(2050, 100, "EUR"),
		SinglePriceNet:    domain.NewFromInt(2050, 100, "EUR"),
	}

	delivery := cart.Delivery{DeliveryInfo: cart.DeliveryInfo{Code: "code"}, Cartitems: []cart.Item{item}}
	cartResult := &cart.Cart{ID: "1", EntityID: "1", DefaultCurrency: "EUR", Deliveries: []cart.Delivery{delivery}}

	return cartResult
}
