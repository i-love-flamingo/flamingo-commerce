package application_test

import (
	"testing"

	localeApplication "flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure/fake"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/leekchan/accounting"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/price/application"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestService_FormatPrice(t *testing.T) {
	t.Parallel()

	translationService := &fake.TranslationService{}

	labelService := &localeApplication.LabelService{}
	labelService.Inject(func() *domain.Label {
		label := &domain.Label{}
		label.Inject(translationService)
		return label
	}, translationService, nil, &struct {
		DefaultLocaleCode string       `inject:"config:core.locale.locale"`
		FallbackLocalCode config.Slice `inject:"config:core.locale.fallbackLocales,optional"`
	}{
		DefaultLocaleCode: "en",
		FallbackLocalCode: config.Slice{"de"},
	})

	service := application.Service{}
	service.Inject(
		labelService,
		&struct {
			Config config.Map `inject:"config:core.locale.accounting"`
		}{
			Config: config.Map{
				"JPY": config.Map{
					"precision":  0.0,
					"decimal":    "",
					"thousand":   ",",
					"formatLong": "%s %s",
					"formatZero": "%s%v",
					"format":     "%s%v",
				},
				"default": config.Map{
					"precision":  2.0,
					"decimal":    ".",
					"thousand":   ",",
					"formatLong": "%s %s",
					"formatZero": "%s%v",
					"format":     "%s%v",
				},
			},
		},
	)

	t.Run("standard USD format from LocaleInfo", func(t *testing.T) {
		t.Parallel()

		price := priceDomain.NewFromFloat(-161.92, "USD")
		formatted := service.FormatPrice(price)
		assert.Equal(t, "-USD 161.92", formatted)
	})

	t.Run("override USD format from LocaleInfo using options", func(t *testing.T) {
		t.Parallel()

		price := priceDomain.NewFromFloat(-1161.92, "USD")
		formatted := service.FormatPrice(price,
			application.WithComSymbol,
			func(ac *accounting.Accounting) {
				ac.Thousand = "#"
				ac.Decimal = "«"
				ac.Precision = 1
			})
		assert.Equal(t, "-$ 1#161«9", formatted)
	})

	t.Run("unknown currency fall back", func(t *testing.T) {
		t.Parallel()

		// unknown currency will fall back to locale default settings as well as to fraction length 0 / precision 1
		price := priceDomain.NewFromFloat(-161.92, "DEFAULT")
		formatted := service.FormatPrice(price)
		assert.Equal(t, "-DEFAULT162.00", formatted)
	})

	t.Run("use core.locale.accounting config", func(t *testing.T) {
		t.Parallel()

		price := priceDomain.NewFromFloat(-161.92, "JPY")
		formatted := service.FormatPrice(price)
		assert.Equal(t, "-JPY162", formatted)
	})

}
