package process

import (
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// Context contains information (state etc) about a place order process
	Context struct {
		UUID               string
		CurrrentStateName  string
		CurrrentStateData  StateData
		Cart               cart.Cart
		ReturnURL          *url.URL
		RollbackReferences []RollbackReference
		FailedReason       FailedReason
	}

	StateData interface{}

	// ContextStore can persist process Context instances
	ContextStore interface {
		Store(key string, value Context) error
		Get(key string) (Context, bool)
		Delete(key string) error
	}
)
