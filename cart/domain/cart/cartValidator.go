package cart

import (
	"context"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// CartValidationResult represents the validation outcome
	CartValidationResult struct {
		HasCommonError        bool
		CommonErrorMessageKey string
		ItemResults           []ItemValidationError
	}

	// ItemValidationError represents a single item's error
	ItemValidationError struct {
		ItemId          string
		UniqueItemID    string
		ErrorMessageKey string
	}

	// CartValidator provides a validation of all items
	CartValidator interface {
		Validate(ctx context.Context, session *web.Session, cart *DecoratedCart) CartValidationResult
	}
)

// IsValid checks if there is any error in the result at all
func (c CartValidationResult) IsValid() bool {
	if c.HasCommonError {
		return false
	}
	if len(c.ItemResults) > 0 {
		return false
	}
	return true
}

// HasErrorForItem checks if there is an error for the given item
func (c CartValidationResult) HasErrorForItem(id string) bool {
	for _, itemMessage := range c.ItemResults {
		if itemMessage.UniqueItemID == id {
			return true
		}
	}
	return false
}

// GetErrorMessageKeyForItem returns the error message key for the given item
func (c CartValidationResult) GetErrorMessageKeyForItem(id string) string {
	for _, itemMessage := range c.ItemResults {
		if itemMessage.UniqueItemID == id {
			return itemMessage.ErrorMessageKey
		}
	}
	return ""
}
