package domain

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"fmt"
	"math"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	// SourcingService helps in retrieving item sources
	SourcingService interface {
		//GetSourceID  returns one source location code where the product should be sourced
		//@todo will be Depricated in future in favor of SourcingServiceDetail interface
		GetSourceID(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, item *decorator.DecoratedCartItem) (string, error)
	}
	// SourcingServiceDetail additional interface to return
	// @todo - the methods in the interface will replace the methods in interface above (SourcingServiceDetail will be deleted then)
	SourcingServiceDetail interface {
		//GetSourcesForItem returns Sources for the given item in the cart
		GetSourcesForItem(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, item *decorator.DecoratedCartItem) (Sources, error)
		//GetAvailableSources returns Sources for the product - containing the maximum possible qty per source
		GetAvailableSources(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart, deliveryCode string, product domain.BasicProduct) (Sources, error)
	}

	//Sources - result value object containing all sources (for the request item or product)
	Sources []Source

	// Source - represents the Sourceing info
	Source struct {
		//LocationCode - idendifies the warehouse or stocklocation
		LocationCode string
		//Qty - for the sources items
		Qty int
	}

	// SourcingEngine computes item sources
	SourcingEngine struct {
		SourcingService SourcingService          `inject:",optional"`
		Logger          flamingo.Logger          `inject:""`
		Cartservice     *application.CartService `inject:""`
	}
)

var (
	//ErrInsufficientSourceQty - use to indicate that the requested qty exceeds the available qty
	ErrInsufficientSourceQty = errors.New("Available Source Qty insufficient")
	//ErrNoSourceAvailable - use to indicate that no source for item is available at all
	ErrNoSourceAvailable = errors.New("No Available Source Qty")
)

// MainLocation returns first sourced location (or empty string)
func (s Sources) MainLocation() string {
	if len(s) < 1 {
		return ""
	}
	return s[0].LocationCode
}

// QtySum returns the sum of all sourced items
func (s Sources) QtySum() int {
	qty := int(0)
	for _, source := range s {
		if source.Qty == math.MaxInt64 {
			return math.MaxInt64
		}
		qty = qty + source.Qty
	}
	return qty
}

// Reduce returns new Source
func (s Sources) Reduce(reduceby Sources) Sources {
	for k, source := range s {
		for _, reducebySource := range reduceby {
			if source.LocationCode == reducebySource.LocationCode {
				s[k].Qty = s[k].Qty - reducebySource.Qty
			}
		}
	}
	return s
}

// SetSourcesForCartItems gets Sources and modifies the Cart Items
// todo move to application layer ?
func (se *SourcingEngine) SetSourcesForCartItems(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart) error {
	if se.SourcingService == nil {
		return nil
	}
	for _, decoratedDelivery := range decoratedCart.DecoratedDeliveries {
		for _, decoratedCartItem := range decoratedDelivery.DecoratedItems {
			sourceID, err := se.SourcingService.GetSourceID(ctx, session, decoratedCart, decoratedDelivery.Delivery.DeliveryInfo.Code, &decoratedCartItem)
			if err != nil {
				se.Logger.WithContext(ctx).WithField("subcategory", "SourcingEngine").Error(err)
				return fmt.Errorf("checkout.application.sourcingengine error: %v", err)
			}
			se.Logger.WithContext(ctx).WithField("category", "checkout").WithField("subcategory", "SourcingEngine").Debug("SourcingEngine detected source %v for item %v", sourceID, decoratedCartItem.Item.ID)
			err = se.Cartservice.UpdateItemSourceID(ctx, session, decoratedCartItem.Item.ID, decoratedDelivery.Delivery.DeliveryInfo.Code, sourceID)
			if err != nil {
				return errors.Wrap(err, "Could not update cart item")
			}
		}
	}
	return nil
}
