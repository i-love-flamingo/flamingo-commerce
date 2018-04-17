package templatefunctions

import (
	"go.aoe.com/flamingo/core/product/application"
	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	// GetFormattedAttributeValue is exported as a template function
	GetFormattedAttributeValue struct {
		FormatService *application.AttributeValueFormatService `inject:""`
		Logger        flamingo.Logger                          `inject:""`
	}
)

// Name alias for use in tempalte
func (tf GetFormattedAttributeValue) Name() string {
	return "getFormattedAttributeValue"
}

// Func returns the json object
func (tf GetFormattedAttributeValue) Func() interface{} {
	return func(a domain.Attribute) string {
		value, err := tf.FormatService.FormatValue(&a)
		if err != nil {
			tf.Logger.Infof("Unable to format attribute value: %v", err.Error())
			return ""
		}

		return value
	}
}
