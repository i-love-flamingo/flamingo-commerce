package application

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/leekchan/accounting"
	"go.aoe.com/flamingo/core/locale/application"
	"go.aoe.com/flamingo/core/product/domain"
)

type (
	// AttributeValueFormatService provides format service on product attribute value
	AttributeValueFormatService struct {
		Precision          float64                         `inject:"config:locale.numbers.precision"`
		Decimal            string                          `inject:"config:locale.numbers.decimal"`
		Thousand           string                          `inject:"config:locale.numbers.thousand"`
		TranslationService *application.TranslationService `inject:""`
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
	switch {
	case a.HasMultipleValues():
		result := make([]string, len(a.Values()))
		for i, v := range a.Values() {
			result[i] = as.format(a.Code, v)
		}
		return strings.Join(result, ", ")
	default:
		return as.format(a.Code, a.RawValue)
	}
}

func (as *AttributeValueFormatService) format(code string, i interface{}) string {
	stringValue := fmt.Sprintf("%v", i)

	// check if value ist an int
	_, e := strconv.ParseInt(stringValue, 10, 64)
	if e == nil {
		return stringValue
	}

	// check if value is a float
	f, e := strconv.ParseFloat(stringValue, 64)
	if e == nil {
		return accounting.FormatNumber(f, int(as.Precision), as.Thousand, as.Decimal)
	}

	// return the nicified value - also try to translate it
	return as.TranslationService.Translate(
		"product.attribute."+code+".value."+stringValue,
		strings.Title(stringValue),
		"",
		0,
		nil,
	)
}
