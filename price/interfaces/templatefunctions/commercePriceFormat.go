package templatefunctions

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/price/application"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

// CommercePriceFormatFunc for formatting prices
type CommercePriceFormatFunc struct {
	priceService *application.Service
}

// Inject dependencies
func (pff *CommercePriceFormatFunc) Inject(priceService *application.Service) {
	pff.priceService = priceService
}

// Func as implementation of debug method
// todo fix
func (pff *CommercePriceFormatFunc) Func(context.Context) interface{} {
	return func(price domain.Price) string {
		return pff.priceService.FormatPrice(price)
	}
}
