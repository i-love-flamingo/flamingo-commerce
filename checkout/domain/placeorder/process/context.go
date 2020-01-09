package process

import "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

type (
	//Context contains information (state etc) about a placeorder process
	Context struct {
		State              State
		Cart               cart.Cart
		RollbackReferences []RollbackReference
	}
)

//CurrentState returns current state
func (c *Context) CurrentState() State {
	return c.State
}
