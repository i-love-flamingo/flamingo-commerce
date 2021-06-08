package validation

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Result groups the validation result
	Result struct {
		HasCommonError        bool
		CommonErrorMessageKey string
		ItemResults           []ItemValidationError
	}

	// ItemValidationError applies for a single item
	ItemValidationError struct {
		ItemID          string
		ErrorMessageKey string
	}

	// Validator checks a complete decorated cart
	Validator interface {
		Validate(ctx context.Context, session *web.Session, cart *decorator.DecoratedCart) Result
	}
)

// IsValid is valid is true if no errors occurred
func (c Result) IsValid() bool {
	if c.HasCommonError {
		return false
	}
	if len(c.ItemResults) > 0 {
		return false
	}
	return true
}

// HasErrorForItem checks if a specified item has an error
func (c Result) HasErrorForItem(id string) bool {
	for _, itemMessage := range c.ItemResults {
		if itemMessage.ItemID == id {
			return true
		}
	}
	return false
}

// GetErrorMessageKeyForItem returns the specific error message for that item
func (c Result) GetErrorMessageKeyForItem(id string) string {
	for _, itemMessage := range c.ItemResults {
		if itemMessage.ItemID == id {
			return itemMessage.ErrorMessageKey
		}
	}
	return ""
}
