package application

import (
	"context"

	"go.aoe.com/flamingo/core/auth/domain"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/event"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	//DomainEventPublisher implements the event publisher of the domain and uses the framework event router
	DomainEventPublisher struct {
		Logger         flamingo.Logger              `inject:""`
		ProductService productDomain.ProductService `inject:""`
	}

	//EventReceiver
	EventReceiver struct {
		Logger      flamingo.Logger `inject:""`
		CartService *CartService    `inject:""`
	}
)

func (d *DomainEventPublisher) PublishOrderPlacedEvent(ctx context.Context, cart *cartDomain.Cart, orderId string) {
	eventObject := cartDomain.OrderPlacedEvent{
		Cart:    cart,
		OrderId: orderId,
	}
	if webContext, ok := ctx.(web.Context); ok {
		d.Logger.Infof("Publish Event OrderPlacedEvent for Order: %v", orderId)
		//For now we publish only to Flamingo default Event Router
		webContext.EventRouter().Dispatch(ctx, &eventObject)
	}
}

func (d *DomainEventPublisher) PublishAddToCartEvent(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode string, qty int) {

	product, err := d.ProductService.Get(ctx, marketPlaceCode)
	if err != nil {
		return
	}
	eventObject := cartDomain.AddToCartEvent{
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

func (d *DomainEventPublisher) PublishChangedQtyInCartEvent(ctx context.Context, item *cartDomain.Item, qtyBefore int, qtyAfter int, cartId string) {

	eventObject := cartDomain.ChangedQtyInCartEvent{
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

//Notify should get called by flamingo Eventlogic
func (e *EventReceiver) Notify(event event.Event) {
	switch eventType := event.(type) {
	//Handle LoginEvent and Merge Cart
	case *domain.LoginEvent:
		if eventType == nil {
			return
		}
		if eventType.Context == nil {
			return
		}
		if !e.CartService.HasSessionAGuestCart(eventType.Context) {
			return
		}
		guestCart, err := e.CartService.GetSessionGuestCart(eventType.Context)
		if err != nil {
			e.Logger.WithField("category", "cart").Errorf("LoginEvent - Guestcart cannot be received %v", err)
			return
		}
		if !e.CartService.UserService.IsLoggedIn(eventType.Context) {
			e.Logger.WithField("category", "cart").Error("Received LoginEvent but user is not logged in!!!")
			return
		}

		for _, item := range guestCart.Cartitems {
			e.Logger.WithField("category", "cart").Debugf("Merging item from guest to user cart %v", item)

			e.CartService.AddProduct(eventType.Context, cartDomain.AddRequest{MarketplaceCode: item.MarketplaceCode, Qty: item.Qty, VariantMarketplaceCode: item.VariantMarketPlaceCode})
		}
		e.CartService.DeleteSessionGuestCart(eventType.Context)

	}
}
