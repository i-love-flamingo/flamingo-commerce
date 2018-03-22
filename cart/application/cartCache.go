package application

import (
	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/web"
)

type (
	CartCache interface {
		GetCart(web.Context, CartCacheIdendifier) (*cart.Cart, error)
		CacheCart(web.Context, CartCacheIdendifier, *cart.Cart) error
		Invalidate(web.Context, CartCacheIdendifier) error
	}

	CartCacheIdendifier struct {
		CartId string
	}

	CartSessionCache struct {
	}

	CachedCartEntry struct {
		IsInvalid bool
		Entry     cart.Cart
	}
	CachedCartsMap map[string]CachedCartEntry
)

var (
	CacheIsInvalidError error = errors.New("Cache is invalid")
)

func (ci *CartCacheIdendifier) CacheKey() string {
	return ci.CartId
}

func (c *CartSessionCache) GetCart(ctx web.Context, id CartCacheIdendifier) (*cart.Cart, error) {
	if cache, ok := ctx.Session().Values["cart.cache"]; ok {
		if cachedCartsMap, ok := cache.(CachedCartsMap); ok {
			if cart, ok := cachedCartsMap[id.CacheKey()]; ok {
				if cart.IsInvalid {
					return &cart.Entry, CacheIsInvalidError
				}
				return &cart.Entry, nil
			}
		}
	}
	return nil, errors.New("no cart in cache")
}

func (c *CartSessionCache) CacheCart(ctx web.Context, id CartCacheIdendifier, cartForCache *cart.Cart) error {
	entry := CachedCartEntry{
		Entry: *cartForCache,
	}

	cachedCartsMap := make(CachedCartsMap)

	//if there is a map in the session already save it there:
	if cache, ok := ctx.Session().Values["cart.cache"]; ok {
		if cachedCartsMap, ok := cache.(CachedCartsMap); ok {
			cachedCartsMap[id.CacheKey()] = entry
		}
	}

	cachedCartsMap[id.CacheKey()] = entry
	ctx.Session().Values["cart.cache"] = cachedCartsMap
	return nil
}

func (c *CartSessionCache) Invalidate(ctx web.Context, id CartCacheIdendifier) error {

	//if there is a map in the session already save it there:
	if cache, ok := ctx.Session().Values["cart.cache"]; ok {
		if cachedCartsMap, ok := cache.(CachedCartsMap); ok {
			if entry, ok := cachedCartsMap[id.CacheKey()]; ok {
				entry.IsInvalid = false
				cachedCartsMap[id.CacheKey()] = entry
				ctx.Session().Values["cart.cache"] = cachedCartsMap
				return nil
			}

		}
	}
	return errors.New("not found")
}
