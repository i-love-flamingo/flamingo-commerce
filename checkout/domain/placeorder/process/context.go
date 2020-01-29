package process

import (
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
)

type (
	// Context contains information (state etc) about a place order process
	Context struct {
		UUID               string
		CurrentStateName   string
		CurrentStateData   StateData
		PlaceOrderInfo     *application.PlaceOrderInfo
		Cart               cart.Cart
		ReturnURL          *url.URL
		RollbackReferences []RollbackReference
		FailedReason       FailedReason
	}
	// StateData holding state relevant data
	StateData interface{}

	// ContextStore can persist process Context instances
	ContextStore interface {
		Store(key string, value Context) error
		Get(key string) (Context, bool)
		Delete(key string) error
	}
)
