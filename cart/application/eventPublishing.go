package application

import (
	"context"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
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
		PublishAddToCartEvent(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode string, qty int)
		PublishChangedQtyInCartEvent(ctx context.Context, item *cartDomain.Item, qtyBefore int, qtyAfter int, cartID string)
		PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, placedOrderInfos cartDomain.PlacedOrderInfos)
	}

	//DefaultEventPublisher implements the event publisher of the domain and uses the framework event router
	DefaultEventPublisher struct {
		logger         flamingo.Logger
		productService productDomain.ProductService
		eventRouter    flamingo.EventRouter
	}
)

var (
	_ EventPublisher = (*DefaultEventPublisher)(nil)
	_ flamingo.Event = (*OrderPlacedEvent)(nil)
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
func (d *DefaultEventPublisher) PublishAddToCartEvent(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode string, qty int) {
	product, err := d.productService.Get(ctx, marketPlaceCode)
	if err != nil {
		return
	}

	eventObject := AddToCartEvent{
		MarketplaceCode:        marketPlaceCode,
		VariantMarketplaceCode: variantMarketPlaceCode,
		ProductName:            product.TeaserData().ShortTitle,
		Qty:                    qty,
	}

	d.logger.Info("Publish Event PublishAddToCartEvent: %v", eventObject)
	d.eventRouter.Dispatch(ctx, &eventObject)
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

	d.logger.Info("Publish Event PublishCartChangedQtyEvent: %v", eventObject)
	d.eventRouter.Dispatch(ctx, &eventObject)
}

// PublishOrderPlacedEvent publishes an event for placed orders
func (d *DefaultEventPublisher) PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, placedOrderInfos cartDomain.PlacedOrderInfos) {
	eventObject := OrderPlacedEvent{
		Cart:             cart,
		PlacedOrderInfos: placedOrderInfos,
	}

	d.logger.Info("Publish Event OrderPlacedEvent for Order: %#v", placedOrderInfos)

	//For now we publish only to Flamingo default Event Router
	d.eventRouter.Dispatch(ctx, &eventObject)
}
