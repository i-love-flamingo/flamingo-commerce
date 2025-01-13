package application

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"

	"flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/leekchan/accounting"
)

// Service for formatting prices
type Service struct {
	config       config.Map
	labelService *application.LabelService
}

// Inject dependencies
func (s *Service) Inject(labelService *application.LabelService, config *struct {
	Config config.Map `inject:"config:core.locale.accounting"`
}) {
	s.labelService = labelService
	s.config = config.Config
}

// GetConfigForCurrency get configuration for currency
func (s *Service) getConfigForCurrency(currency string) config.Map {
	if configForCurrency, ok := s.config[currency]; ok {
		return configForCurrency.(config.Map)
	}

	if info, ok := accounting.LocaleInfo[currency]; ok {
		format := "%s %v"
		if !info.Pre {
			format = "%v %s"
		}

		return config.Map{
			"decimal":   info.DecSep,
			"thousand":  info.ThouSep,
			"precision": float64(info.FractionLength),
			"format":    format,
		}
	}

	if defaultConfig, ok := s.config["default"].(config.Map); ok {
		return defaultConfig
	}
	return nil
}

// WithComSymbol tries to get the commercial symbol from LocaleInfo and overrides the currency code if found
func WithComSymbol(ac *accounting.Accounting) {
	if info, ok := accounting.LocaleInfo[ac.Symbol]; ok {
		ac.Symbol = info.ComSymbol
	}
}

// FormatPrice by price
func (s *Service) FormatPrice(price domain.Price, options ...func(*accounting.Accounting)) string {
	currency := s.labelService.NewLabel(price.Currency()).String()

	configForCurrency := s.getConfigForCurrency(price.Currency())

	ac := accounting.Accounting{
		Symbol:    currency,
		Precision: 2,
	}

	precision, ok := configForCurrency["precision"].(float64)
	if ok {
		ac.Precision = int(precision)
	}

	decimal, ok := configForCurrency["decimal"].(string)
	if ok {
		ac.Decimal = decimal
	}

	thousand, ok := configForCurrency["thousand"].(string)
	if ok {
		ac.Thousand = thousand
	}

	formatZero, ok := configForCurrency["formatZero"].(string)
	if ok {
		ac.FormatZero = formatZero
	}

	format, ok := configForCurrency["format"].(string)
	if ok {
		ac.Format = format
	}

	for _, option := range options {
		option(&ac)
	}

	return ac.FormatMoney(price.GetPayable().FloatAmount())
}
