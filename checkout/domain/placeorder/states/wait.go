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
	/* Todo: maybe split in multiple states
	1. Reserve Order id
	2. Start Payment
	3. Reserve Order if EarlyPlace
	4. Get Payment Status
	// n.ctx.UpdateState(WaitingForPaymentInformation)
	*/

	return nil, nil
}

func (n Wait) IsFinal() bool {
	return false
}
