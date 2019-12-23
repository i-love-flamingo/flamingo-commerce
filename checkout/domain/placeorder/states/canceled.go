package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"
)

type (
	Canceled struct {
		ctx          *placeorder.Context
		runFunctions []func(Rollback, error) //
	}
)

var _ State = Canceled{}

func (c Canceled) SetContext(ctx *placeorder.Context) {
	c.ctx = ctx
}

func (c Canceled) Run() (Rollback, error) {
	// todo
	return nil, nil
}

func (c Canceled) IsFinal() bool {
	return true
}
