package process

import (
	"net/url"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
)

type (
	// Context contains information (state etc) about a place order process
	Context struct {
		UUID  string
		State string
		// todo: maybe add SateData so that the current state can expose data
		Cart               cart.Cart
		ReturnURL          *url.URL
		RollbackReferences []RollbackReference
		FailedReason       FailedReason
		// URL is used to pass URL data to the user if the current state needs some
		URL *url.URL
		// DisplayData holds data, normally HTML to be displayed to the user
		DisplayData   string
		FormParameter map[string]FormField
	}

	// FormField represents a form field to be displayed to the user
	FormField struct {
		Value []string
	}
)
