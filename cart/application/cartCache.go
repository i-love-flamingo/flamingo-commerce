package application

import (
	"context"
	"encoding/gob"
	"fmt"
	"strings"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo/framework/flamingo"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

type (
	// CartCache describes a cart caches methods
	CartCache interface {
		GetCart(context.Context, *sessions.Session, CartCacheIdentifier) (*cart.Cart, error)
		CacheCart(context.Context, *sessions.Session, CartCacheIdentifier, *cart.Cart) error
		Invalidate(context.Context, *sessions.Session, CartCacheIdentifier) error
		Delete(context.Context, *sessions.Session, CartCacheIdentifier) error
		DeleteAll(context.Context, *sessions.Session) error
	}

	// CartCacheIdentifier identifies Cart Caches
	CartCacheIdentifier struct {
		GuestCartID    string
		IsCustomerCart bool
		CustomerID     string
	}

	// CartSessionCache defines a Cart Cache
	CartSessionCache struct {
		Logger flamingo.Logger `inject:""`
	}

	// CachedCartEntry defines a single Cart Cache Entry
	CachedCartEntry struct {
		IsInvalid bool
		Entry     cart.Cart
	}
)

const (
	// CartSessionCacheCacheKeyPrefix is a string prefix for Cart Cache Keys
	CartSessionCacheCacheKeyPrefix = "cart.sessioncache."
)

var (
	// ErrCacheIsInvalid sets generalized invalid Cache Error
	ErrCacheIsInvalid = errors.New("cache is invalid")
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
func BuildIdentifierFromCart(cart *cart.Cart) (*CartCacheIdentifier, error) {
	if cart == nil {
		return nil, errors.New("no cart")
	}

	if cart.BelongsToAuthenticatedUser {
		return &CartCacheIdentifier{
			CustomerID:     cart.AuthenticatedUserId,
			IsCustomerCart: true,
		}, nil
	}

	return &CartCacheIdentifier{
		GuestCartID:    cart.ID,
		CustomerID:     cart.AuthenticatedUserId,
		IsCustomerCart: false,
	}, nil
}

// GetCart fetches a Cart from the Cache
func (c *CartSessionCache) GetCart(ctx context.Context, session *sessions.Session, id CartCacheIdentifier) (*cart.Cart, error) {
	if cache, ok := session.Values[CartSessionCacheCacheKeyPrefix+id.CacheKey()]; ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			c.Logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").Debug("Found cached cart %v", id.CacheKey())

			if cachedCartsEntry.IsInvalid {
				return &cachedCartsEntry.Entry, ErrCacheIsInvalid
			}

			return &cachedCartsEntry.Entry, nil
		}
		c.Logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").Error("Cannot Cast Cache Entry %v", id.CacheKey())

		return nil, errors.New("cart cache contains invalid data at cache key")
	}
	c.Logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").Debug("Did not Found cached cart %v", id.CacheKey())

	return nil, errors.New("no cart in cache")
}

// CacheCart adds a Cart to the Cache
func (c *CartSessionCache) CacheCart(ctx context.Context, session *sessions.Session, id CartCacheIdentifier, cartForCache *cart.Cart) error {
	if cartForCache == nil {
		return errors.New("no cart given to cache")
	}

	entry := CachedCartEntry{
		Entry: *cartForCache,
	}

	c.Logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").Debug("Caching cart %v", id.CacheKey())
	session.Values[CartSessionCacheCacheKeyPrefix+id.CacheKey()] = entry

	return nil
}

// Invalidate a Cache Entry
func (c *CartSessionCache) Invalidate(ctx context.Context, session *sessions.Session, id CartCacheIdentifier) error {
	if cache, ok := session.Values[CartSessionCacheCacheKeyPrefix+id.CacheKey()]; ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			cachedCartsEntry.IsInvalid = true
			session.Values[CartSessionCacheCacheKeyPrefix+id.CacheKey()] = cachedCartsEntry

			return nil
		}
	}

	return errors.New("not found for invalidate")
}

// Delete a Cache entry
func (c *CartSessionCache) Delete(ctx context.Context, session *sessions.Session, id CartCacheIdentifier) error {
	if _, ok := session.Values[CartSessionCacheCacheKeyPrefix+id.CacheKey()]; ok {
		delete(session.Values, CartSessionCacheCacheKeyPrefix+id.CacheKey())

		// ok deleted something
		return nil
	}

	return errors.New("not found for delete")
}

// DeleteAll empties the Cache
func (c *CartSessionCache) DeleteAll(ctx context.Context, session *sessions.Session) error {
	deleted := false
	for k := range session.Values {
		if stringKey, ok := k.(string); ok {
			if strings.Contains(stringKey, CartSessionCacheCacheKeyPrefix) {
				delete(session.Values, k)
				deleted = true
			}
		}
	}

	if deleted {
		// successfully deleted something
		return nil
	}

	return errors.New("not found for delete")
}
