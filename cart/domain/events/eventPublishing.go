package events

import (
	"context"

	"flamingo.me/flamingo/v3/framework/flamingo"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

type (

	// EventPublisher technology free interface to publish events that might be interesting for outside (Publish)
	EventPublisher interface {
		PublishAddToCartEvent(ctx context.Context, cart *cartDomain.Cart, marketPlaceCode string, variantMarketPlaceCode string, qty int)
		PublishChangedQtyInCartEvent(ctx context.Context, cart *cartDomain.Cart, item *cartDomain.Item, qtyBefore int, qtyAfter int)
		PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, placedOrderInfos placeorder.PlacedOrderInfos)
	}

	// DefaultEventPublisher implements the event publisher of the domain and uses the framework event router
	DefaultEventPublisher struct {
		logger         flamingo.Logger
		productService productDomain.ProductService
		eventRouter    flamingo.EventRouter
	}
)

var (
	_ EventPublisher = (*DefaultEventPublisher)(nil)
	_ flamingo.Event = (*OrderPlacedEvent)(nil)
	_ flamingo.Event = (*AddToCartEvent)(nil)
	_ flamingo.Event = (*PaymentSelectionHasBeenResetEvent)(nil)
	_ flamingo.Event = (*ChangedQtyInCartEvent)(nil)
)

// Inject dependencies
func (d *DefaultEventPublisher) Inject(
	logger flamingo.Logger,
	productService productDomain.ProductService,
	eventRouter flamingo.EventRouter,
) {
	d.logger = logger
	d.productService = productService
	d.eventRouter = eventRouter
}

// PublishAddToCartEvent publishes an event for add to cart actions
func (d *DefaultEventPublisher) PublishAddToCartEvent(ctx context.Context, cart *cartDomain.Cart, marketPlaceCode string, variantMarketPlaceCode string, qty int) {
	product, err := d.productService.Get(ctx, marketPlaceCode)
	if err != nil {
		return
	}

	eventObject := AddToCartEvent{
		MarketplaceCode:        marketPlaceCode,
		VariantMarketplaceCode: variantMarketPlaceCode,
		ProductName:            product.TeaserData().ShortTitle,
		Qty:                    qty,
		Cart:                   cart,
	}

	d.logger.WithContext(ctx).Info("Publish Event PublishAddToCartEvent: ", eventObject)
	d.eventRouter.Dispatch(ctx, &eventObject)
}

// PublishChangedQtyInCartEvent publishes an event for cart item quantity change actions
func (d *DefaultEventPublisher) PublishChangedQtyInCartEvent(ctx context.Context, cart *cartDomain.Cart, item *cartDomain.Item, qtyBefore int, qtyAfter int) {
	eventObject := ChangedQtyInCartEvent{
		CartID:                 cart.ID,
		MarketplaceCode:        item.MarketplaceCode,
		VariantMarketplaceCode: item.VariantMarketPlaceCode,
		ProductName:            item.ProductName,
		QtyBefore:              qtyBefore,
		QtyAfter:               qtyAfter,
		Cart:                   cart,
	}

	d.logger.WithContext(ctx).Info("Publish Event PublishCartChangedQtyEvent: ", eventObject)
	d.eventRouter.Dispatch(ctx, &eventObject)
}

// PublishOrderPlacedEvent publishes an event for placed orders
func (d *DefaultEventPublisher) PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, placedOrderInfos placeorder.PlacedOrderInfos) {
	eventObject := OrderPlacedEvent{
		Cart:             cart,
		PlacedOrderInfos: placedOrderInfos,
	}

	d.logger.WithContext(ctx).Info("Publish Event OrderPlacedEvent for Order: ", placedOrderInfos)

	// For now we publish only to Flamingo default Event Router
	d.eventRouter.Dispatch(ctx, &eventObject)
}
