package application_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo/framework/flamingo"
	"github.com/gorilla/sessions"
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

func TestBuildIdentifierFromCart(t *testing.T) {
	type args struct {
		cart *cart.Cart
	}
	tests := []struct {
		name           string
		args           args
		want           *application.CartCacheIdentifier
		wantErr        bool
		wantMessageErr string
	}{
		{
			name: "error for no cart",
			args: args{
				cart: nil,
			},
			want:           nil,
			wantErr:        true,
			wantMessageErr: "no cart",
		},
		{
			name: "authenticated user returned",
			args: args{
				cart: &cart.Cart{
					BelongsToAuthenticatedUser: true,
					AuthenticatedUserId:        "test_user_id",
				},
			},
			want: &application.CartCacheIdentifier{
				GuestCartID:    "",
				IsCustomerCart: true,
				CustomerID:     "test_user_id",
			},
			wantErr:        false,
			wantMessageErr: "",
		},
		{
			name: "guest user returned",
			args: args{
				cart: &cart.Cart{
					ID:                  "test_cart_id",
					AuthenticatedUserId: "test_user_id",
				},
			},
			want: &application.CartCacheIdentifier{
				GuestCartID:    "test_cart_id",
				IsCustomerCart: false,
				CustomerID:     "test_user_id",
			},
			wantErr:        false,
			wantMessageErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := application.BuildIdentifierFromCart(tt.args.cart)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildIdentifierFromCart() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildIdentifierFromCart() = %v, wantType0 %v", got, tt.want)
			}
		})
	}
}

func TestCartSessionCache_GetCart(t *testing.T) {
	type args struct {
		ctx     context.Context
		session *sessions.Session
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
				ctx: context.Background(),
				session: &sessions.Session{
					ID: "test_session",
				},
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           nil,
			wantErr:        true,
			wantMessageErr: "no cart in cache",
		}, {
			name: "cached cart found/invalid cache entry",
			args: args{
				ctx: context.Background(),
				session: &sessions.Session{
					ID: "test_session",
					Values: map[interface{}]interface{}{
						application.CartSessionCacheCacheKeyPrefix + "cart_customer_id_guest_cart_id": application.CachedCartEntry{
							IsInvalid: true,
						},
					},
				},
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           new(cart.Cart),
			wantErr:        true,
			wantMessageErr: "cache is invalid",
		}, {
			name: "cached cart found/valid cache entry",
			args: args{
				ctx: context.Background(),
				session: &sessions.Session{
					ID: "test_session",
					Values: map[interface{}]interface{}{
						application.CartSessionCacheCacheKeyPrefix + "cart_customer_id_guest_cart_id": application.CachedCartEntry{
							IsInvalid: false,
							ExpiresOn: time.Now().Add(5 * time.Minute),
						},
					},
				},
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           new(cart.Cart),
			wantErr:        false,
			wantMessageErr: "",
		}, {
			name: "session contains invalid data at cache key",
			args: args{
				ctx: context.Background(),
				session: &sessions.Session{
					ID: "test_session",
					Values: map[interface{}]interface{}{
						application.CartSessionCacheCacheKeyPrefix + "cart_customer_id_guest_cart_id": struct {
							invalidProperty bool
						}{
							invalidProperty: true,
						},
					},
				},
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "customer_id",
				},
			},
			want:           nil,
			wantErr:        true,
			wantMessageErr: "cart cache contains invalid data at cache key",
		}, {
			name: "session contains expired cart cache",
			args: args{
				ctx: context.Background(),
				session: &sessions.Session{
					ID: "test_session",
					Values: map[interface{}]interface{}{
						application.CartSessionCacheCacheKeyPrefix + "cart_customer_id_guest_cart_id": application.CachedCartEntry{
							IsInvalid: false,
							ExpiresOn: time.Now().Add(-1 * time.Second),
						},
					},
				},
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
			c.Inject(nil, nil, flamingo.NullLogger{}, nil)

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
			LifetimeSeconds float64 `inject:"config:cart.cacheLifetime"` // in seconds
		}
	}
	type args struct {
		ctx          context.Context
		session      *sessions.Session
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
					LifetimeSeconds float64 `inject:"config:cart.cacheLifetime"` // in seconds
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
		}, {
			name:   "cart is cached",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				session: &sessions.Session{
					ID:     "test_session",
					Values: map[interface{}]interface{}{},
				},
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
			c.Inject(nil, nil, flamingo.NullLogger{}, nil)

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
	ctx := context.Background()
	c := &application.CartSessionCache{}
	c.Inject(
		nil,
		nil,
		flamingo.NullLogger{},
		&struct {
			LifetimeSeconds float64 `inject:"config:cart.cacheLifetime"` // in seconds
		}{
			LifetimeSeconds: 1,
		},
	)

	session := &sessions.Session{
		ID:     "test_session",
		Values: map[interface{}]interface{}{},
	}

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
		session *sessions.Session
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
				ctx: context.Background(),
				session: &sessions.Session{
					ID:     "test_session",
					Values: map[interface{}]interface{}{},
				},
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "",
				},
			},
			wantErr:               true,
			wantMessageErr:        "not found for invalidate",
			wantCacheEntryInvalid: false,
		}, {
			name:   "invalidate cache entry",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				session: &sessions.Session{
					ID: "test_session",
					Values: map[interface{}]interface{}{
						application.CartSessionCacheCacheKeyPrefix + "cart__guest_cart_id": application.CachedCartEntry{
							IsInvalid: false,
						},
					},
				},
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
			c.Inject(nil, nil, flamingo.NullLogger{}, nil)

			err := c.Invalidate(tt.args.ctx, tt.args.session, tt.args.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartSessionCache.Invalidate() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)
			}

			if tt.wantCacheEntryInvalid == true {
				for _, sessionValueValue := range tt.args.session.Values {
					if cacheEntry, ok := sessionValueValue.(application.CachedCartEntry); ok {
						if cacheEntry.IsInvalid != tt.wantCacheEntryInvalid {
							t.Errorf("Cache validity doesnt match - got %v, wantCacheEntryInvalid %v", cacheEntry.IsInvalid, tt.wantCacheEntryInvalid)
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
		session *sessions.Session
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
				ctx: context.Background(),
				session: &sessions.Session{
					ID:     "test_session",
					Values: map[interface{}]interface{}{},
				},
				id: application.CartCacheIdentifier{
					GuestCartID:    "guest_cart_id",
					IsCustomerCart: false,
					CustomerID:     "",
				},
			},
			wantErr:        true,
			wantMessageErr: "not found for delete",
		}, {
			name:   "deleted correclty",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				session: &sessions.Session{
					ID: "test_session",
					Values: map[interface{}]interface{}{
						application.CartSessionCacheCacheKeyPrefix + "cart__guest_cart_id": application.CachedCartEntry{
							IsInvalid: false,
						},
					},
				},
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
			c.Inject(nil, nil, flamingo.NullLogger{}, nil)

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
		session *sessions.Session
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
				ctx: context.Background(),
				session: &sessions.Session{
					ID:     "test_session",
					Values: map[interface{}]interface{}{},
				},
			},
			wantErr:                true,
			wantMessageErr:         "not found for delete",
			wantSessionValuesEmpty: false,
		},
		{
			name:   "deleted an entry",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				session: &sessions.Session{
					ID: "test_session",
					Values: map[interface{}]interface{}{
						application.CartSessionCacheCacheKeyPrefix + "cart__guest_cart_id": application.CachedCartEntry{
							IsInvalid: false,
						},
					},
				},
			},
			wantErr:                false,
			wantMessageErr:         "",
			wantSessionValuesEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &application.CartSessionCache{}
			c.Inject(nil, nil, flamingo.NullLogger{}, nil)

			err := c.DeleteAll(tt.args.ctx, tt.args.session)

			if (err != nil) != tt.wantErr {
				t.Errorf("CartSessionCache.DeleteAll() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil && err.Error() != tt.wantMessageErr {
				t.Errorf("Error doesn't match - error = %v, wantMessageErr %v", err, tt.wantMessageErr)
			}

			if tt.wantSessionValuesEmpty == true {
				if len(tt.args.session.Values) > 0 {
					t.Error("Session Values should have been emptied, but aren't")
				}
			}
		})
	}
}
