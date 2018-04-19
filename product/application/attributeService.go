package application

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/leekchan/accounting"
	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/config"
)

// Renderer codes
const (
	RendererDefault   = "default"
	RendererNochange  = "nochange"
	RendererTranslate = "translate"
	RendererTitle     = "title"
)

type (
	// AttributeValueFormatService provides format service on product attribute value
	AttributeValueFormatService struct {
		Precision          float64                                 `inject:"config:locale.numbers.precision"`
		Decimal            string                                  `inject:"config:locale.numbers.decimal"`
		Thousand           string                                  `inject:"config:locale.numbers.thousand"`
		RendererConfig     config.Map                              `inject:"config.templating.product.attributeRenderer"`
		TranslationService application.TranslationServiceInterface `inject:""`
	}

	rendererConfig struct {
		Renderer []string
	}
)

// FormatValue formats the attribute value
func (as *AttributeValueFormatService) FormatValue(a *domain.Attribute) (string, error) {
	if a == nil {
		return "", errors.New("no attribute provided")
	}

	formattedValue := as.formatValue(a)

	if a.HasUnitCode() {
		return fmt.Sprintf(
			"%v %v",
			formattedValue,
			as.TranslationService.Translate(a.GetUnit().Symbol, a.GetUnit().Symbol, "", 1, nil),
		), nil
	}

	return formattedValue, nil
}

func (as *AttributeValueFormatService) formatValue(a *domain.Attribute) string {
	renderer := as.getAttributeRenderer(a)

	switch {
	case a.HasMultipleValues():
		result := make([]string, len(a.Values()))
		for i, v := range a.Values() {
			result[i] = as.render(renderer, a, v)
		}
		return strings.Join(result, ", ")
	default:
		return as.render(renderer, a, a.RawValue)
	}
}

func (as *AttributeValueFormatService) render(config rendererConfig, a *domain.Attribute, i interface{}) string {
	var result = fmt.Sprintf("%v", i)
	for _, code := range config.Renderer {
		switch code {
		case RendererNochange:
			// just do nothing
			break

		case RendererTranslate:
			result = as.translate(a.Code, result)

		case RendererTitle:
			result = as.title(result)

		default:
			result = as.format(a.Code, result)
		}
	}

	return result
}

// format returns the default formatted value with
// - float format per config
// - titled and translated value
func (as *AttributeValueFormatService) format(code string, v string) string {
	// check if value ist an int
	_, e := strconv.ParseInt(v, 10, 64)
	if e == nil {
		return v
	}

	// check if value is a float
	f, e := strconv.ParseFloat(v, 64)
	if e == nil {
		return accounting.FormatNumber(f, int(as.Precision), as.Thousand, as.Decimal)
	}

	return as.translate(code, strings.Title(v))
}

func (as *AttributeValueFormatService) translate(code string, v string) string {
	// return the nicified value - also try to translate it
	return as.TranslationService.Translate(
		"product.attribute."+code+".value."+v,
		v,
		"",
		0,
		nil,
	)
}

// title is the super forced title (Camel Case) of v
func (as *AttributeValueFormatService) title(v string) string {
	return strings.Title(strings.ToLower(v))
}

// formatLabel returns the (translated) value

func (as *AttributeValueFormatService) getAttributeRenderer(a *domain.Attribute) rendererConfig {
	c, ok := as.RendererConfig[a.Code]
	if !ok {
		return rendererConfig{
			Renderer: []string{RendererDefault},
		}
	}

	result := rendererConfig{}
	c.(config.Map).MapInto(&result)

	return result
}
