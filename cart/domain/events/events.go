package events

import (
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
)

type (
	// OrderPlacedEvent defines event properties
	OrderPlacedEvent struct {
		Cart             *cartDomain.Cart
		PlacedOrderInfos placeorder.PlacedOrderInfos
	}

	// AddToCartEvent defines event properties
	AddToCartEvent struct {
		Cart                   *cartDomain.Cart
		MarketplaceCode        string
		VariantMarketplaceCode string
		ProductName            string
		Qty                    int
	}

	// ChangedQtyInCartEvent defines event properties
	ChangedQtyInCartEvent struct {
		Cart *cartDomain.Cart
		// Deprecated: CartID exists for compatibility, use Cart instead
		CartID                 string
		MarketplaceCode        string
		VariantMarketplaceCode string
		ProductName            string
		QtyBefore              int
		QtyAfter               int
	}

	// PaymentSelectionHasBeenResetEvent defines event properties
	PaymentSelectionHasBeenResetEvent struct {
		Cart                     *cartDomain.Cart
		ResettedPaymentSelection *cartDomain.PaymentSelection
	}
)
