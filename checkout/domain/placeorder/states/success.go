package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"
)

type (
	Success struct {
		ctx          *placeorder.Context
		runFunctions []func(Rollback, error) //
	}
)

var _ State = Success{}

func (s Success) SetContext(ctx *placeorder.Context) {
	s.ctx = ctx
}

func (s Success) Run() (Rollback, error) {
	// todo
	return nil, nil
}

func (s Success) IsFinal() bool {
	return true
}
