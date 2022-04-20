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
	Config config.Map `inject:"config:locale.accounting"`
}) {
	s.labelService = labelService
	s.config = config.Config
}

// GetConfigForCurrency get configuration for currency
func (s *Service) getConfigForCurrency(currency string) config.Map {
	if configForCurrency, ok := s.config[currency]; ok {
		return configForCurrency.(config.Map)
	}

	if defaultConfig, ok := s.config["default"].(config.Map); ok {
		return defaultConfig
	}
	return nil
}

// FormatPrice by price
func (s *Service) FormatPrice(price domain.Price) string {
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

	return ac.FormatMoney(price.GetPayable().FloatAmount())
}
