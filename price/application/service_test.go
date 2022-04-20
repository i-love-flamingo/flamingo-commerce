package application_test

import (
	"testing"

	localeApplication "flamingo.me/flamingo/v3/core/locale/application"
	"flamingo.me/flamingo/v3/core/locale/domain"
	"flamingo.me/flamingo/v3/core/locale/infrastructure/fake"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/price/application"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

func TestService_FormatPrice(t *testing.T) {
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
			Config config.Map `inject:"config:locale.accounting"`
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

	price := priceDomain.NewFromFloat(-161.92, "USD")
	formatted := service.FormatPrice(price)
	assert.Equal(t, "-USD161.92", formatted)

	price = priceDomain.NewFromFloat(-161.92, "JPY")
	formatted = service.FormatPrice(price)
	assert.Equal(t, "-JPY162", formatted)
}
