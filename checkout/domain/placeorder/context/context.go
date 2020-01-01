package context

type (
	Context struct {
		// todo more fields, see graphql result
		state State
	}
)

func (c *Context) CurrentState() State {
	return c.state
}

func (c *Context) UpdateState(s State) {
	c.state = s
}
