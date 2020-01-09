package placeorder

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

//Handler for handling PlaceOrder related commands
type Handler struct {
	coordinator *Coordinator
}

// Inject dependencies
func (h *Handler) Inject(
	c *Coordinator,
) *Handler {
	h.coordinator = c

	return h
}

//StartPlaceOrder handles start place order command
func (h *Handler) StartPlaceOrder(ctx context.Context, command StartPlaceOrderCommand) (*process.Context, error) {
	return h.coordinator.New(ctx, command.Cart)
}

//RefreshPlaceOrder handles start RefreshPlaceOrder command
func (h *Handler) RefreshPlaceOrder(ctx context.Context, command RefreshPlaceOrderCommand) (*process.Context, error) {
	return h.coordinator.Last(ctx)
}
