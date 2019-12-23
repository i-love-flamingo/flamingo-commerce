package states

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder"
)

type (
	Rollback func() error

	State interface {
		SetContext(ctx *placeorder.Context)
		Run() (Rollback, error)
		IsFinal() bool
	}
)
