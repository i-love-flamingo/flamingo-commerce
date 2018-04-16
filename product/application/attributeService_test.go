package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/flamingo"
)

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

func getAttributeValueFormatService() *AttributeValueFormatService {
	return &AttributeValueFormatService{
		Precision:          2,
		Thousand:           ",",
		Decimal:            ".",
		TranslationService: &application.TranslationService{Logger: flamingo.NullLogger{}},
	}
}
