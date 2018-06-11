package application

import (
	"encoding/gob"
	"fmt"

	"strings"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	CartCache interface {
		GetCart(web.Context, CartCacheIdentifier) (*cart.Cart, error)
		CacheCart(web.Context, CartCacheIdentifier, *cart.Cart) error
		Invalidate(web.Context, CartCacheIdentifier) error
		Delete(web.Context, CartCacheIdentifier) error
		DeleteAll(web.Context) error
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
	CacheIsInvalidError error = errors.New("Cache is invalid")
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

func (c *CartSessionCache) GetCart(ctx web.Context, id CartCacheIdentifier) (*cart.Cart, error) {
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

func (c *CartSessionCache) CacheCart(ctx web.Context, id CartCacheIdentifier, cartForCache *cart.Cart) error {
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

func (c *CartSessionCache) Invalidate(ctx web.Context, id CartCacheIdentifier) error {
	if cache, ok := ctx.Session().Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()]; ok {
		if cachedCartsEntry, ok := cache.(CachedCartEntry); ok {
			cachedCartsEntry.IsInvalid = false
			ctx.Session().Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()] = cachedCartsEntry
			return nil
		}
	}

	return errors.New("not found for invalidate")
}

func (c *CartSessionCache) Delete(ctx web.Context, id CartCacheIdentifier) error {
	if _, ok := ctx.Session().Values[CartSessionCache_CacheKeyPrefix+id.CacheKey()]; ok {
		delete(ctx.Session().Values, CartSessionCache_CacheKeyPrefix+id.CacheKey())
		// ok deleted something
		return nil
	}

	return errors.New("not found for delete")
}

func (c *CartSessionCache) DeleteAll(ctx web.Context) error {
	deleted := false
	for k, _ := range ctx.Session().Values {
		if stringKey, ok := k.(string); ok {
			if strings.Contains(stringKey, CartSessionCache_CacheKeyPrefix) {
				delete(ctx.Session().Values, k)
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
