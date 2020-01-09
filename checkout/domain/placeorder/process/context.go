package process

import "flamingo.me/flamingo-commerce/v3/cart/domain/cart"

type (
	//Context contains information (state etc) about a placeorder process
	Context struct {
		state              State
		cart               cart.Cart
		rollbackReferences []RollbackReference
	}
)

//CurrentState returns current state
func (c *Context) CurrentState() State {
	return c.state
}
