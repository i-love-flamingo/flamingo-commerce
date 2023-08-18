package application

import (
	"context"
	"encoding/gob"
	"fmt"
	"strings"
	"time"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name CartCache --case snake

type (
	// CartCache describes a cart caches methods
	CartCache interface {
		GetCart(context.Context, *web.Session, CartCacheIdentifier) (*cart.Cart, error)
		CacheCart(context.Context, *web.Session, CartCacheIdentifier, *cart.Cart) error
		Invalidate(context.Context, *web.Session, CartCacheIdentifier) error
		Delete(context.Context, *web.Session, CartCacheIdentifier) error
		DeleteAll(context.Context, *web.Session) error
		BuildIdentifier(context.Context, *web.Session) (CartCacheIdentifier, error)
	}

	// CartCacheIdentifier identifies Cart Caches
	CartCacheIdentifier struct {
		GuestCartID    string
		IsCustomerCart bool
		CustomerID     string
	}

	// CartSessionCache defines a Cart Cache
	CartSessionCache struct {
		logger             flamingo.Logger
		webIdentityService *auth.WebIdentityService
		lifetimeSeconds    float64
	}

	// CachedCartEntry defines a single Cart Cache Entry
	CachedCartEntry struct {
		IsInvalid bool
		Entry     cart.Cart
		ExpiresOn time.Time
	}
)

const (
	// CartSessionCacheCacheKeyPrefix is a string prefix for Cart Cache Keys
	CartSessionCacheCacheKeyPrefix = "cart.sessioncache."
)

var (
	_ CartCache = (*CartSessionCache)(nil)
	// ErrCacheIsInvalid sets generalized invalid Cache Error
	ErrCacheIsInvalid = errors.New("cache is invalid")
	// ErrNoCacheEntry - used if cache is not found
	ErrNoCacheEntry = errors.New("cache entry not found")
)

func init() {
	gob.Register(CachedCartEntry{})
}

// CacheKey creates a Cache Key Identifier string
func (ci *CartCacheIdentifier) CacheKey() string {
	return fmt.Sprintf(
		"cart_%v_%v",
		ci.CustomerID,
		ci.GuestCartID,
	)
}

// BuildIdentifierFromCart creates a Cache Identifier from Cart Data
// Deprecated: use BuildIdentifier function of concrete implementation
func BuildIdentifierFromCart(cart *cart.Cart) (*CartCacheIdentifier, error) {
	if cart == nil {
		return nil, errors.New("no cart")
	}

	if cart.BelongsToAuthenticatedUser {
		return &CartCacheIdentifier{
			CustomerID:     cart.AuthenticatedUserID,
			IsCustomerCart: true,
		}, nil
	}

	return &CartCacheIdentifier{
		GuestCartID:    cart.ID,
		CustomerID:     cart.AuthenticatedUserID,
		IsCustomerCart: false,
	}, nil
}

// Inject the dependencies
func (cs *CartSessionCache) Inject(
	logger flamingo.Logger,
	webIdentityService *auth.WebIdentityService,
	config *struct {
		LifetimeSeconds float64 `inject:"config:commerce.cart.cacheLifetime"` // in seconds
	},
) {
	cs.webIdentityService = webIdentityService
	cs.logger = logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").WithField(flamingo.LogKeyModule, "cart")

	if config != nil {
		cs.lifetimeSeconds = config.LifetimeSeconds
	}
}

// BuildIdentifier creates a CartCacheIdentifier based on the login state
func (cs *CartSessionCache) BuildIdentifier(ctx context.Context, session *web.Session) (CartCacheIdentifier, error) {
	identity := cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx))
	if identity != nil {
		return CartCacheIdentifier{
			CustomerID:     identity.Subject(),
			IsCustomerCart: true,
		}, nil
	}

	guestCartID, ok := session.Load(GuestCartSessionKey)
	if !ok {
		return CartCacheIdentifier{}, errors.New("Fatal - ShouldHaveGuestCart returned true but got no GuestCartSessionKey?")
	}

	guestCartIDString, ok := guestCartID.(string)
	if !ok {
		return CartCacheIdentifier{}, errors.New("Fatal - ShouldHaveGuestCart returned true but got no GuestCartSessionKey string")
	}

	return CartCacheIdentifier{
		GuestCartID: guestCartIDString,
	}, nil
}

// GetCart fetches a Cart from the Cache
func (cs *CartSessionCache) GetCart(ctx context.Context, session *web.Session, id CartCacheIdentifier) (*cart.Cart, error) {
	if cache, ok := session.Load(CartSessionCacheCacheKeyPrefix + id.CacheKey()); ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			cs.logger.WithContext(ctx).Debugf("Found cached cart: %v  InValid: %v", id.CacheKey(), cachedCartsEntry.IsInvalid)
			if cachedCartsEntry.IsInvalid {
				return &cachedCartsEntry.Entry, ErrCacheIsInvalid
			}

			if time.Now().After(cachedCartsEntry.ExpiresOn) {
				err := cs.Invalidate(ctx, session, id)
				if err != nil {
					return nil, err
				}

				return nil, ErrCacheIsInvalid
			}

			return &cachedCartsEntry.Entry, nil
		}
		cs.logger.WithContext(ctx).Error("Cannot Cast Cache Entry %v", id.CacheKey())

		return nil, errors.New("cart cache contains invalid data at cache key")
	}
	cs.logger.WithContext(ctx).Debug("Did not Found cached cart %v", id.CacheKey())

	return nil, ErrNoCacheEntry
}

// CacheCart adds a Cart to the Cache
func (cs *CartSessionCache) CacheCart(ctx context.Context, session *web.Session, id CartCacheIdentifier, cartForCache *cart.Cart) error {
	if cartForCache == nil {
		return errors.New("no cart given to cache")
	}
	entry := CachedCartEntry{
		Entry:     *cartForCache,
		ExpiresOn: time.Now().Add(time.Duration(cs.lifetimeSeconds * float64(time.Second))),
	}

	cs.logger.WithContext(ctx).Debug("Caching cart %v", id.CacheKey())
	session.Store(CartSessionCacheCacheKeyPrefix+id.CacheKey(), entry)
	return nil
}

// Invalidate a Cache Entry
func (cs *CartSessionCache) Invalidate(ctx context.Context, session *web.Session, id CartCacheIdentifier) error {
	if cache, ok := session.Load(CartSessionCacheCacheKeyPrefix + id.CacheKey()); ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			cachedCartsEntry.IsInvalid = true
			session.Store(CartSessionCacheCacheKeyPrefix+id.CacheKey(), cachedCartsEntry)

			return nil
		}
	}

	return ErrNoCacheEntry
}

// Delete a Cache entry
func (cs *CartSessionCache) Delete(ctx context.Context, session *web.Session, id CartCacheIdentifier) error {
	if _, ok := session.Load(CartSessionCacheCacheKeyPrefix + id.CacheKey()); ok {
		session.Delete(CartSessionCacheCacheKeyPrefix + id.CacheKey())

		// ok deleted something
		return nil
	}

	return ErrNoCacheEntry
}

// DeleteAll empties the Cache
func (cs *CartSessionCache) DeleteAll(ctx context.Context, session *web.Session) error {
	deleted := false
	for _, k := range session.Keys() {
		if stringKey, ok := k.(string); ok {
			if strings.Contains(stringKey, CartSessionCacheCacheKeyPrefix) {
				session.Delete(k)
				deleted = true
			}
		}
	}

	if deleted {
		// successfully deleted something
		return nil
	}

	return ErrNoCacheEntry
}
