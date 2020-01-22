package process

import (
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// Context contains information (state etc) about a placeorder process
	Context struct {
		UUID               string
		State              string
		Cart               cart.Cart
		ReturnURL          *url.URL
		RollbackReferences []RollbackReference
		FailedReason       FailedReason
		Data               interface{}
	}
)
