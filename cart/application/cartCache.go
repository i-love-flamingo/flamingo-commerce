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
	CartCache interface {
		GetCart(context.Context, *sessions.Session, CartCacheIdentifier) (*cart.Cart, error)
		CacheCart(context.Context, *sessions.Session, CartCacheIdentifier, *cart.Cart) error
		Invalidate(context.Context, *sessions.Session, CartCacheIdentifier) error
		Delete(context.Context, *sessions.Session, CartCacheIdentifier) error
		DeleteAll(context.Context, *sessions.Session) error
	}

	CartCacheIdentifier struct {
		GuestCartId    string
		IsCustomerCart bool
		CustomerId     string
	}

	CartSessionCache struct {
		Logger flamingo.Logger `inject:""`
	}

	CachedCartEntry struct {
		IsInvalid bool
		Entry     cart.Cart
	}
)

const (
	CartSessionCache_CacheKeyPrefix = "cart.sessioncache."
)

var (
	CacheIsInvalidError error     = errors.New("Cache is invalid")
	_                   CartCache = new(CartSessionCache)
)

func init() {
	gob.Register(CachedCartEntry{})
}

func (ci *CartCacheIdentifier) CacheKey() string {
	return fmt.Sprintf(
		"cart_%v_%v",
		ci.CustomerId,
		ci.GuestCartId,
	)
}

func BuildIdentifierFromCart(cart *cart.Cart) (*CartCacheIdentifier, error) {
	if cart == nil {
		return nil, errors.New("no cart")
	}
	if cart.BelongsToAuthenticatedUser {
		return &CartCacheIdentifier{
			CustomerId:     cart.AuthenticatedUserId,
			IsCustomerCart: true,
		}, nil
	}

	return &CartCacheIdentifier{
		GuestCartId:    cart.ID,
		CustomerId:     cart.AuthenticatedUserId,
		IsCustomerCart: false,
	}, nil
}

func (c *CartSessionCache) GetCart(ctx context.Context, session *sessions.Session, id CartCacheIdentifier) (*cart.Cart, error) {
	if cache, ok := session.Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()]; ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			c.Logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").Debug("Found cached cart %v", id.CacheKey())
			if cachedCartsEntry.IsInvalid {
				return &cachedCartsEntry.Entry, CacheIsInvalidError
			}
			return &cachedCartsEntry.Entry, nil
		} else {
			c.Logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").Error("Cannot Cast Cache Entry %v", id.CacheKey())
		}
	}
	c.Logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").Debug("Did not Found cached cart %v", id.CacheKey())

	return nil, errors.New("no cart in cache")
}

func (c *CartSessionCache) CacheCart(ctx context.Context, session *sessions.Session, id CartCacheIdentifier, cartForCache *cart.Cart) error {
	if cartForCache == nil {
		return errors.New("No cart given to cache")
	}
	entry := CachedCartEntry{
		Entry: *cartForCache,
	}
	c.Logger.WithField(flamingo.LogKeyCategory, "CartSessionCache").Debug("Caching cart %v", id.CacheKey())
	session.Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()] = entry
	return nil
}

func (c *CartSessionCache) Invalidate(ctx context.Context, session *sessions.Session, id CartCacheIdentifier) error {
	if cache, ok := session.Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()]; ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			cachedCartsEntry.IsInvalid = false
			session.Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()] = cachedCartsEntry
			return nil
		}
	}

	return errors.New("not found for invalidate")
}

func (c *CartSessionCache) Delete(ctx context.Context, session *sessions.Session, id CartCacheIdentifier) error {
	if _, ok := session.Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()]; ok {
		delete(session.Values, CartSessionCache_CacheKeyPrefix+id.CacheKey())
		// ok deleted something
		return nil
	}

	return errors.New("not found for delete")
}

func (c *CartSessionCache) DeleteAll(ctx context.Context, session *sessions.Session) error {
	deleted := false
	for k, _ := range session.Values {
		if stringKey, ok := k.(string); ok {
			if strings.Contains(stringKey, CartSessionCache_CacheKeyPrefix) {
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
