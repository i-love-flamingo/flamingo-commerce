package application

import (
	"fmt"

	"github.com/pkg/errors"
	"go.aoe.com/flamingo/core/cart/application"
	"go.aoe.com/flamingo/core/cart/domain/cart"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	SourcingService interface {
		GetSourceId(ctx web.Context, decoratedCart *cart.DecoratedCart, item *cart.DecoratedCartItem) (string, error)
	}

	SourcingEngine struct {
		SourcingService SourcingService          `inject:",optional"`
		Logger          flamingo.Logger          `inject:""`
		Cartservice     *application.CartService `inject:""`
	}
)

// SetSourcesForCartItems gets Sources and modifies the Cart Items
func (se *SourcingEngine) SetSourcesForCartItems(ctx web.Context, decoratedCart *cart.DecoratedCart, behaviour cart.CartBehaviour) error {
	if se.SourcingService == nil {
		return nil
	}
	for _, decoratedCartItem := range decoratedCart.DecoratedItems {
		sourceId, err := se.SourcingService.GetSourceId(ctx, decoratedCart, &decoratedCartItem)
		if err != nil {
			return fmt.Errorf("checkout.application.sourcingengine error: %v", err)
		}
		cartItem := decoratedCartItem.Item
		cartItem.SourceId = sourceId
		err = behaviour.UpdateItem(ctx, &decoratedCart.Cart, cartItem.ID, cartItem)
		if err != nil {
			return errors.Wrap(err, "Could not update cart item")
		}
	}
	return nil
}
