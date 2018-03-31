package application

import (
	"encoding/gob"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	CartCache interface {
		GetCart(web.Context, CartCacheIdendifier) (*cart.Cart, error)
		CacheCart(web.Context, CartCacheIdendifier, *cart.Cart) error
		Invalidate(web.Context, CartCacheIdendifier) error
		Delete(web.Context, CartCacheIdendifier) error
	}

	CartCacheIdendifier struct {
		GuestCartId    string
		IsCustomerCart bool
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
	CacheIsInvalidError error = errors.New("Cache is invalid")
)

func init() {
	gob.Register(CachedCartEntry{})
}

func (ci *CartCacheIdendifier) CacheKey() string {
	if ci.IsCustomerCart {
		return "customer_" + ci.GuestCartId
	}
	return ci.GuestCartId
}

func BuildIdendifierFromCart(cart *cart.Cart) (*CartCacheIdendifier, error) {
	if cart == nil {
		return nil, errors.New("no cart")
	}
	if cart.IsCustomerCart {
		return &CartCacheIdendifier{
			IsCustomerCart: true,
		}, nil
	}
	return &CartCacheIdendifier{
		GuestCartId: cart.ID,
	}, nil
}

func (c *CartSessionCache) GetCart(ctx web.Context, id CartCacheIdendifier) (*cart.Cart, error) {
	if cache, ok := ctx.Session().Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()]; ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			c.Logger.WithField("category", "CartSessionCache").Debugf("Found cached cart %v", id.CacheKey())
			if cachedCartsEntry.IsInvalid {
				return &cachedCartsEntry.Entry, CacheIsInvalidError
			}
			return &cachedCartsEntry.Entry, nil
		} else {
			c.Logger.WithField("category", "CartSessionCache").Errorf("Cannot Cast Cache Entry %v", id.CacheKey())
		}
	}
	c.Logger.WithField("category", "CartSessionCache").Debugf("Did not Found cached cart %v", id.CacheKey())

	return nil, errors.New("no cart in cache")
}

func (c *CartSessionCache) CacheCart(ctx web.Context, id CartCacheIdendifier, cartForCache *cart.Cart) error {
	if cartForCache == nil {
		return errors.New("No cart given to cache")
	}
	entry := CachedCartEntry{
		Entry: *cartForCache,
	}
	c.Logger.WithField("category", "CartSessionCache").Debugf("Caching cart %v", id.CacheKey())
	ctx.Session().Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()] = entry
	return nil
}

func (c *CartSessionCache) Invalidate(ctx web.Context, id CartCacheIdendifier) error {
	if cache, ok := ctx.Session().Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()]; ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			cachedCartsEntry.IsInvalid = false
			ctx.Session().Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()] = cachedCartsEntry
			return nil
		}
	}
	return errors.New("not found")
}

func (c *CartSessionCache) Delete(ctx web.Context, id CartCacheIdendifier) error {
	if _, ok := ctx.Session().Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()]; ok {
		delete(ctx.Session().Values, CartSessionCache_CacheKeyPrefix+id.CacheKey())
	}
	return errors.New("not found")
}
