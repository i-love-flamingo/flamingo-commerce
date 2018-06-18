package cart

import "flamingo.me/flamingo/framework/web"

type (
	CartValidationResult struct {
		HasCommonError        bool
		CommonErrorMessageKey string
		ItemResults           []ItemValidationError
	}

	ItemValidationError struct {
		ItemId          string
		ErrorMessageKey string
	}

	CartValidator interface {
		Validate(ctx web.Context, cart *DecoratedCart) CartValidationResult
	}
)

func (c CartValidationResult) IsValid() bool {
	if c.HasCommonError {
		return false
	}
	if len(c.ItemResults) > 0 {
		return false
	}
	return true
}

func (c CartValidationResult) HasErrorForItem(id string) bool {
	for _, itemMessage := range c.ItemResults {
		if itemMessage.ItemId == id {
			return true
		}
	}
	return false
}

func (c CartValidationResult) GetErrorMessageKeyForItem(id string) string {
	for _, itemMessage := range c.ItemResults {
		if itemMessage.ItemId == id {
			return itemMessage.ErrorMessageKey
		}
	}
	return ""
}
