package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/context"
)

type (
	Failed struct {
		ctx          *context.Context
		runFunctions []func(context.Rollback, error) //
	}
)

var _ context.State = Failed{}

func (f Failed) SetContext(ctx *context.Context) {
	f.ctx = ctx
}

func (f Failed) Run() (context.Rollback, error) {
	// todo
	return nil, nil
}

func (f Failed) IsFinal() bool {
	return true
}
