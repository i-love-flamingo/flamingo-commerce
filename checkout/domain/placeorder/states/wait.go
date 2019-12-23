package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"
)

type (
	Wait struct {
		ctx *placeorder.Context
	}
)

var _ State = Wait{}

func (n Wait) SetContext(ctx *placeorder.Context) {
	n.ctx = ctx
}

func (n Wait) Run() (Rollback, error) {
	// n.ctx.UpdateState()
	return nil, nil
}

func (n Wait) IsFinal() bool {
	return false
}
