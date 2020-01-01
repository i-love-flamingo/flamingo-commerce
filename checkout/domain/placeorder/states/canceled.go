package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/context"
)

type (
	Canceled struct {
		ctx          *context.Context
		runFunctions []func(context.Rollback, error) //
	}
)

var _ context.State = Canceled{}

func (c Canceled) SetContext(ctx *context.Context) {
	c.ctx = ctx
}

func (c Canceled) Run() (context.Rollback, error) {
	// todo
	return nil, nil
}

func (c Canceled) IsFinal() bool {
	return true
}
