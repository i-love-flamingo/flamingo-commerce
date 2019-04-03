package cart

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// ValidationResult groups the validation result
	ValidationResult struct {
		HasCommonError        bool
		CommonErrorMessageKey string
		ItemResults           []ItemValidationError
	}

	// ItemValidationError applies for a single item
	ItemValidationError struct {
		ItemID          string
		UniqueItemID    string
		ErrorMessageKey string
	}

	// Validator checks a complete decorated cart
	Validator interface {
		Validate(ctx context.Context, session *web.Session, cart *DecoratedCart) ValidationResult
	}
)

// IsValid is valid is true if no errors occurred
func (c ValidationResult) IsValid() bool {
	if c.HasCommonError {
		return false
	}
	if len(c.ItemResults) > 0 {
		return false
	}
	return true
}

// HasErrorForItem checks if a specified item has an error
func (c ValidationResult) HasErrorForItem(id string) bool {
	for _, itemMessage := range c.ItemResults {
		if itemMessage.UniqueItemID == id {
			return true
		}
	}
	return false
}

// GetErrorMessageKeyForItem returns the specific error message for that item
func (c ValidationResult) GetErrorMessageKeyForItem(id string) string {
	for _, itemMessage := range c.ItemResults {
		if itemMessage.UniqueItemID == id {
			return itemMessage.ErrorMessageKey
		}
	}
	return ""
}
