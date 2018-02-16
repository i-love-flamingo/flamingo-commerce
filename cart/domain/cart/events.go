package cart

import (
	"context"

	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	OrderPlacedEvent struct {
		Cart           *Cart
		OrderId        string
		CurrentContext web.Context
	}

	AddToCartEvent struct {
		Product         productDomain.BasicProduct
		Qty             int
	}

	ChangedQtyInCartEvent struct {
		CartId          string
		Product         productDomain.BasicProduct
		QtyBefore       int
		QtyAfter        int
	}

	//EventPublisher - technology free interface which is used in the Domain Layer to publish events that might be interesting for outside (Publish)
	EventPublisher interface {
		PublishOrderPlacedEvent(ctx context.Context, cart *Cart, orderId string)
		PublishAddToCartEvent(ctx context.Context, product productDomain.BasicProduct, qty int)
		PublishChangedQtyInCartEvent(ctx context.Context, item *Item, qtyBefore int, qtyAfter int, cartId string)
	}
)
