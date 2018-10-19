package application

import (
	"context"

	cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo/framework/event"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	// OrderPlacedEvent defines event properties
	OrderPlacedEvent struct {
		Cart             *cartDomain.Cart
		PlacedOrderInfos cartDomain.PlacedOrderInfos
	}

	// AddToCartEvent defines event properties
	AddToCartEvent struct {
		MarketplaceCode        string
		VariantMarketplaceCode string
		ProductName            string
		Qty                    int
	}

	// ChangedQtyInCartEvent defines event properties
	ChangedQtyInCartEvent struct {
		CartID                 string
		MarketplaceCode        string
		VariantMarketplaceCode string
		ProductName            string
		QtyBefore              int
		QtyAfter               int
	}

	//EventPublisher - technology free interface to publish events that might be interesting for outside (Publish)
	EventPublisher interface {
		PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, placedOrderInfos cartDomain.PlacedOrderInfos)
		PublishAddToCartEvent(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode string, qty int)
		PublishChangedQtyInCartEvent(ctx context.Context, item *cartDomain.Item, qtyBefore int, qtyAfter int, cartID string)
	}

	//DefaultEventPublisher implements the event publisher of the domain and uses the framework event router
	DefaultEventPublisher struct {
		Logger         flamingo.Logger              `inject:""`
		ProductService productDomain.ProductService `inject:""`
		EventRouter    event.Router                 `inject:""`
	}
)

// PublishOrderPlacedEvent publishes an event for placed orders
func (d *DefaultEventPublisher) PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, placedOrderInfos cartDomain.PlacedOrderInfos) {
	eventObject := OrderPlacedEvent{
		Cart:             cart,
		PlacedOrderInfos: placedOrderInfos,
	}

	d.Logger.Info("Publish Event OrderPlacedEvent for Order: %#v", placedOrderInfos)

	//For now we publish only to Flamingo default Event Router
	d.EventRouter.Dispatch(ctx, &eventObject)
}

// PublishAddToCartEvent publishes an event for add to cart actions
func (d *DefaultEventPublisher) PublishAddToCartEvent(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode string, qty int) {
	product, err := d.ProductService.Get(ctx, marketPlaceCode)
	if err != nil {
		return
	}

	eventObject := AddToCartEvent{
		MarketplaceCode:        marketPlaceCode,
		VariantMarketplaceCode: variantMarketPlaceCode,
		ProductName:            product.TeaserData().ShortTitle,
		Qty:                    qty,
	}

	d.Logger.Info("Publish Event PublishAddToCartEvent: %v", eventObject)
	d.EventRouter.Dispatch(ctx, &eventObject)
}

// PublishChangedQtyInCartEvent publishes an event for cart item quantity change actions
func (d *DefaultEventPublisher) PublishChangedQtyInCartEvent(ctx context.Context, item *cartDomain.Item, qtyBefore int, qtyAfter int, cartID string) {
	eventObject := ChangedQtyInCartEvent{
		CartID:                 cartID,
		MarketplaceCode:        item.MarketplaceCode,
		VariantMarketplaceCode: item.VariantMarketPlaceCode,
		ProductName:            item.ProductName,
		QtyBefore:              qtyBefore,
		QtyAfter:               qtyAfter,
	}

	d.Logger.Info("Publish Event PublishCartChangedQtyEvent: %v", eventObject)
	d.EventRouter.Dispatch(ctx, &eventObject)
}
