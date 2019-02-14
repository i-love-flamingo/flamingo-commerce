package application

import (
	"context"
	"flamingo.me/flamingo/v3/framework/web"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/order/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// OrderService provides methods to place/fetch orders
	OrderService struct {
		logger               flamingo.Logger
		cartService          *cartApplication.CartService
		orderReceiverService *OrderReceiverService
		eventPublisher       EventPublisher
	}
)

// Inject the order service dependencies
func (os *OrderService) Inject(
	Logger flamingo.Logger,
	CartService *cartApplication.CartService,
	OrderReceiverService *OrderReceiverService,
	EventPublisher EventPublisher,
) {
	os.logger = Logger
	os.cartService = CartService
	os.orderReceiverService = OrderReceiverService
	os.eventPublisher = EventPublisher
}

// PlaceOrder submits an order
func (os *OrderService) PlaceOrder(ctx context.Context, session *web.Session, payment *cart.CartPayment) (domain.PlacedOrderInfos, error) {
	cart, _, err := os.cartService.GetCartReceiverService().GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	behaviour, err := os.orderReceiverService.GetBehaviour(ctx, session)
	if err != nil {
		return nil, err
	}

	orderNumbers, err := behaviour.PlaceOrder(ctx, cart, payment)
	if err != nil {
		os.handleCartNotFound(session, err)
		os.logger.WithField("category", "orderService").WithField("subCategory", "PlaceOrder").Error(err)

		return nil, err
	}

	os.eventPublisher.PublishOrderPlacedEvent(ctx, cart, orderNumbers)
	os.cartService.DeleteSavedSessionGuestCartID(session)
	os.cartService.DeleteCartInCache(ctx, session, cart)

	return orderNumbers, err
}

func (os *OrderService) handleCartNotFound(session *web.Session, err error) {
	if err == cart.CartNotFoundError {
		os.cartService.DeleteSavedSessionGuestCartID(session)
	}
}
