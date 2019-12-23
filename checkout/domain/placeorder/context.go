package placeorder

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
)

type (
	Context struct {
		// todo more fields, see graphql result
		state states.State
	}
)

func (c *Context) CurrentState() states.State {
	return r.state
}

func (c *Context) UpdateState(s states.State) {
	c.state = s
}
