package process

import (
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// Context contains information (state etc) about a placeorder process
	Context struct {
		UUID               string
		State              State
		Cart               cart.Cart
		ReturnURL          *url.URL
		RollbackReferences []RollbackReference
	}
)

// CurrentState returns current state
func (c *Context) CurrentState() State {
	return c.State
}
