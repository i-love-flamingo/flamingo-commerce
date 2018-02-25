package cart

import (
	"context"
)

type (
	OrderPlacedEvent struct {
		Cart    *Cart
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

	//EventPublisher - technology free interface which is used in the Domain Layer to publish events that might be interesting for outside (Publish)
	EventPublisher interface {
		PublishOrderPlacedEvent(ctx context.Context, cart *Cart, orderId string)
		PublishAddToCartEvent(ctx context.Context, marketPlaceCode string, variantMarketPlaceCode string, qty int)
		PublishChangedQtyInCartEvent(ctx context.Context, item *Item, qtyBefore int, qtyAfter int, cartId string)
	}
)
