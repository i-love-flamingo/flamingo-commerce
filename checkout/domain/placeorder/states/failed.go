package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"
)

type (
	Failed struct {
		ctx          *placeorder.Context
		runFunctions []func(Rollback, error) //
	}
)

var _ State = Failed{}

func (f Failed) SetContext(ctx *placeorder.Context) {
	f.ctx = ctx
}

func (f Failed) Run() (Rollback, error) {
	// todo
	return nil, nil
}

func (f Failed) IsFinal() bool {
	return true
}
