package templatefunctions

import (
	"go.aoe.com/flamingo/core/product/application"
	"go.aoe.com/flamingo/core/product/domain"
)

type (
	// GetFormattedAttributeValue is exported as a template function
	GetFormattedAttributeValue struct {
		FormatService *application.AttributeValueFormatService `inject:""`
	}
)

// Name alias for use in tempalte
func (tf GetFormattedAttributeValue) Name() string {
	return "getFormattedValue"
}

// Func returns the json object
func (tf GetFormattedAttributeValue) Func() interface{} {
	return func(a *domain.Attribute) string {
		if a == nil {
			return ""
		}

		value, err := tf.FormatService.FormatValue(a)
		if err != nil {
			return ""
		}

		return value
	}
}
