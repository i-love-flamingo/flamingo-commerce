package application

import (
	"context"

	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	OrderPlacedEvent struct {
		Cart    *cartDomain.Cart
		OrderId string
	}

	AddToCartEvent struct {
		MarketplaceCode        string
		VariantMarketplaceCode string
		ProductName            string
		Qty                    int
	}

	ChangedQtyInCartEvent struct {
		CartId                 string
		MarketplaceCode        string
		VariantMarketplaceCode string
		ProductName            string
		QtyBefore              int
		QtyAfter               int
	}

	//EventPublisher - technology free interface  to publish events that might be interesting for outside (Publish)
	EventPublisher interface {
		PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, orderId string)
		PublishAddToCartEvent(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode string, qty int)
		PublishChangedQtyInCartEvent(ctx context.Context, item *cartDomain.Item, qtyBefore int, qtyAfter int, cartId string)
	}

	//DefaultEventPublisher implements the event publisher of the domain and uses the framework event router
	DefaultEventPublisher struct {
		Logger         flamingo.Logger              `inject:""`
		ProductService productDomain.ProductService `inject:""`
	}
)

func (d *DefaultEventPublisher) PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, orderId string) {
	eventObject := OrderPlacedEvent{
		Cart:    cart,
		OrderId: orderId,
	}
	if webContext, ok := ctx.(web.Context); ok {
		d.Logger.Infof("Publish Event OrderPlacedEvent for Order: %v", orderId)
		//For now we publish only to Flamingo default Event Router
		webContext.EventRouter().Dispatch(ctx, &eventObject)
	}
}

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
	if webContext, ok := ctx.(web.Context); ok {
		d.Logger.Infof("Publish Event PublishAddToCartEvent: %v", eventObject)
		webContext.EventRouter().Dispatch(ctx, &eventObject)
	}
}

func (d *DefaultEventPublisher) PublishChangedQtyInCartEvent(ctx context.Context, item *cartDomain.Item, qtyBefore int, qtyAfter int, cartId string) {

	eventObject := ChangedQtyInCartEvent{
		CartId:                 cartId,
		MarketplaceCode:        item.MarketplaceCode,
		VariantMarketplaceCode: item.VariantMarketPlaceCode,
		ProductName:            item.ProductName,
		QtyBefore:              qtyBefore,
		QtyAfter:               qtyAfter,
	}
	if webContext, ok := ctx.(web.Context); ok {
		d.Logger.Infof("Publish Event PublishCartChangedQtyEvent: %v", eventObject)
		webContext.EventRouter().Dispatch(ctx, &eventObject)
	}
}
