package placeorder

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

// Handler for handling PlaceOrder related commands
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

// StartPlaceOrder handles start place order command
func (h *Handler) StartPlaceOrder(ctx context.Context, command StartPlaceOrderCommand) (*process.Context, error) {
	_ = h.coordinator.ClearLastProcess(ctx)
	return h.coordinator.New(ctx, command.Cart, command.ReturnURL)
}

// CurrentContext returns the last saved state
func (h *Handler) CurrentContext(ctx context.Context) (*process.Context, error) {
	p, err := h.coordinator.LastProcess(ctx)
	if err != nil {
		return nil, err
	}
	currentContext := p.Context()

	return &currentContext, nil
}

// ClearPlaceOrder clears the last placed order from the context store, only possible if order in final state
func (h *Handler) ClearPlaceOrder(ctx context.Context) error {
	return h.coordinator.ClearLastProcess(ctx)
}

// RefreshPlaceOrder handles RefreshPlaceOrder command
func (h *Handler) RefreshPlaceOrder(ctx context.Context, _ RefreshPlaceOrderCommand) (*process.Context, error) {
	p, err := h.coordinator.LastProcess(ctx)
	if err != nil {
		return nil, err
	}
	lastPlaceOrderCtx := p.Context()

	// proceed in state
	h.coordinator.Run(ctx)

	return &lastPlaceOrderCtx, nil
}

// RefreshPlaceOrderBlocking handles RefreshPlaceOrder blocking
func (h *Handler) RefreshPlaceOrderBlocking(ctx context.Context, _ RefreshPlaceOrderCommand) (*process.Context, error) {
	return h.coordinator.RunBlocking(ctx)
}

// HasUnfinishedProcess checks for processes not in final state
func (h *Handler) HasUnfinishedProcess(ctx context.Context) (bool, error) {
	return h.coordinator.HasUnfinishedProcess(ctx)
}

// CancelPlaceOrder handles order cancellation, is blocking
func (h *Handler) CancelPlaceOrder(ctx context.Context, _ CancelPlaceOrderCommand) error {
	return h.coordinator.Cancel(ctx)
}
