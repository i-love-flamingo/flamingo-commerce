package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/flamingo"
)

type TranslationServiceMock struct {
	mock.Mock
}

func (t *TranslationServiceMock) Translate(key string, defaultLabel string, localeCode string, count int, translationArguments map[string]interface{}) string {
	args := t.Called()

	return args.String(0)
}

func TestFormatNilAttribute(t *testing.T) {
	as := getAttributeValueFormatService()

	v, e := as.FormatValue(nil)
	assert.Empty(t, v, "formatted value of nil attribute is empty")
	assert.Equal(t, e.Error(), "no attribute provided")
}

func TestFormatSingleStringValueAttribute(t *testing.T) {
	a := &domain.Attribute{
		RawValue: "some string",
	}

	as := getAttributeValueFormatService()
	v, e := as.FormatValue(a)
	assert.Nil(t, e, "no format error")
	assert.Equal(t, v, "Some String")
}

func TestFormatMultipleValueAttribute(t *testing.T) {
	var values []interface{}
	values = append(values, "some")
	values = append(values, "string")
	values = append(values, 1)
	values = append(values, 10.990000)
	values = append(values, "10.990000")
	values = append(values, "1000000.990000")
	values = append(values, false)
	a := &domain.Attribute{
		RawValue: values,
	}

	as := getAttributeValueFormatService()
	v, e := as.FormatValue(a)
	assert.Nil(t, e, "no format error")
	assert.Equal(t, "Some, String, 1, 10.99, 10.99, 1,000,000.99, False", v)
}

func TestFormatWithUnit(t *testing.T) {
	a := &domain.Attribute{
		RawValue: "10.99",
		UnitCode: domain.MILLILITER,
	}

	as := getAttributeValueFormatService()
	v, e := as.FormatValue(a)
	assert.Nil(t, e, "no format error")
	assert.Equal(t, "10.99 ml", v)
}

func TestRendererNochange(t *testing.T) {
	a := &domain.Attribute{
		Code:     "test",
		RawValue: "10.99000",
	}

	as := getAttributeValueFormatService()
	as.RendererConfig = config.Map{
		"test": config.Map{
			"renderer": []string{"nochange"},
		},
	}

	v, e := as.FormatValue(a)
	assert.Nil(t, e, "rendered without error")
	assert.Equal(t, "10.99000", v)
}

func TestRendererTitle(t *testing.T) {
	a := &domain.Attribute{
		Code:     "test",
		RawValue: "UnFORTunate TiTLE!",
	}

	as := getAttributeValueFormatService()
	as.RendererConfig = config.Map{
		"test": config.Map{
			"renderer": []string{"title"},
		},
	}

	v, e := as.FormatValue(a)
	assert.Nil(t, e, "rendered without error")
	assert.Equal(t, "Unfortunate Title!", v)
}

func TestRendererTranslate(t *testing.T) {
	a := &domain.Attribute{
		Code:     "test",
		RawValue: "untranslated",
	}

	as := getAttributeValueFormatService()
	as.RendererConfig = config.Map{
		"test": config.Map{
			"renderer": []string{"translate"},
		},
	}
	translationServiceMock := new(TranslationServiceMock)
	translationServiceMock.On("Translate", mock.Anything).Return("translated")

	as.TranslationService = translationServiceMock

	v, e := as.FormatValue(a)
	assert.Nil(t, e, "rendered without error")
	assert.Equal(t, "translated", v)
}

func getAttributeValueFormatService() *AttributeValueFormatService {
	return &AttributeValueFormatService{
		Precision:          2,
		Thousand:           ",",
		Decimal:            ".",
		TranslationService: &application.TranslationService{Logger: flamingo.NullLogger{}},
	}
}
