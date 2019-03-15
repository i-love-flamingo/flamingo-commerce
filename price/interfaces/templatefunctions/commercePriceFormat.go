package templatefunctions

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/price/domain"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/leekchan/accounting"
)

// CommercePriceFormatFunc for formatting prices
type CommercePriceFormatFunc struct {
	config       config.Map
	labelService *application.LabelService
}

// Inject dependencies
func (pff *CommercePriceFormatFunc) Inject(labelService *application.LabelService, config *struct {
	Config config.Map `inject:"config:locale.accounting"`
}) {
	pff.labelService = labelService
	pff.config = config.Config
}

// Func as implementation of debug method
// todo fix
func (pff *CommercePriceFormatFunc) Func(context.Context) interface{} {
	return func(price domain.Price) string {
		currency := pff.labelService.NewLabel(price.Currency()).String()
		ac := accounting.Accounting{
			Symbol:    currency,
			Precision: 2,
		}
		decimal, ok := pff.config["decimal"].(string)
		if ok {
			ac.Decimal = decimal
		}
		thousand, ok := pff.config["thousand"].(string)
		if ok {
			ac.Thousand = thousand
		}
		formatZero, ok := pff.config["formatZero"].(string)
		if ok {
			ac.FormatZero = formatZero
		}
		format, ok := pff.config["format"].(string)
		if ok {
			ac.Format = format
		}

		return ac.FormatMoney(price.GetPayable().FloatAmount())
	}
}
