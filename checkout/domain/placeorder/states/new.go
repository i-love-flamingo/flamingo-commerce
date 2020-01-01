package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/context"
)

type (
	New struct {
		ctx          *context.Context
		runFunctions []func(context.Rollback, error) //
	}
)

var _ context.State = New{}

func (n New) SetContext(ctx *context.Context) {
	n.ctx = ctx
}

func (n New) Run() (context.Rollback, error) {
	n.ctx.UpdateState(Wait{
		ctx: n.ctx,
	})

	return nil, nil
}

func (n New) IsFinal() bool {
	return false
}
