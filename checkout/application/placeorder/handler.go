package placeorder

import (
	"errors"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"

	placeorderContext "flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/context"
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

func (h *Handler) Handle(command interface{}) (placeorderContext.Context, error) {
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
		return placeorderContext.Context{}, errors.New("invalid command")
	}
}
