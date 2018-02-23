package cart

import (
	"context"
	"encoding/gob"
)

type (
	OrderPlacedEvent struct {
		Cart           *Cart
		OrderId        string
	}

	AddToCartEvent struct {
		ProductIdentifier string
		Qty               int
	}

	ChangedQtyInCartEvent struct {
		CartId            string
		ProductIdentifier string
		QtyBefore         int
		QtyAfter          int
	}

	//EventPublisher - technology free interface which is used in the Domain Layer to publish events that might be interesting for outside (Publish)
	EventPublisher interface {
		PublishOrderPlacedEvent(ctx context.Context, cart *Cart, orderId string)
		PublishAddToCartEvent(ctx context.Context, marketPlaceCode string, qty int)
		PublishChangedQtyInCartEvent(ctx context.Context, item *Item, qtyBefore int, qtyAfter int, cartId string)
	}
)

func init() {
	gob.Register(AddToCartEvent{})
	gob.Register(ChangedQtyInCartEvent{})
}
