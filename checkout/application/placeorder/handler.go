package placeorder

import (
	"errors"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"
)

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

func (h *Handler) Handle(command interface{}) (placeorder.Context, error) {
	switch c := command.(type) {
	case placeorder.StartPlaceOrder:
		return h.coordinator.New(c.Ctx, c.Cart)
	case placeorder.RefreshPlaceOrder:
		h.coordinator.Run(c.Ctx, c.Cart)
		return h.coordinator.Current(c.Ctx, c.Cart)
	case placeorder.RefreshBlockingPlaceOrder:
		return h.coordinator.RunBlocking(c.Ctx, c.Cart)
	case placeorder.CancelPlaceOrder:
		return h.coordinator.Cancel(c.Ctx, c.Cart)
	default:
		return placeorder.Context{}, errors.New("invalid command")
	}
}
