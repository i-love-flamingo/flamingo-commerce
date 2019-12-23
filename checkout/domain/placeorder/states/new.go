package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"
)

type (
	New struct {
		ctx          *placeorder.Context
		runFunctions []func(Rollback, error) //
	}
)

var _ State = New{}

func (n New) SetContext(ctx *placeorder.Context) {
	n.ctx = ctx
}

func (n New) Run() (Rollback, error) {
	n.ctx.UpdateState(Wait{
		ctx: n.ctx,
	})

	return nil, nil
}

func (n New) IsFinal() bool {
	return false
}
