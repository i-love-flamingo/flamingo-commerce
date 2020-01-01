package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/context"
)

type (
	Success struct {
		ctx          *context.Context
		runFunctions []func(context.Rollback, error) //
	}
)

var _ context.State = Success{}

func (s Success) SetContext(ctx *context.Context) {
	s.ctx = ctx
}

func (s Success) Run() (context.Rollback, error) {
	// todo
	return nil, nil
}

func (s Success) IsFinal() bool {
	return true
}
