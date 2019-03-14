package domain

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAttributeValue(t *testing.T) {
	a := Attribute{RawValue: "testValue"}

	assert.Equal(t, a.Value(), "testValue")
}

func TestAttributeIsEnabledValue(t *testing.T) {
	a := Attribute{RawValue: "Yes"}
	assert.True(t, a.IsEnabledValue())

	a.RawValue = "yes"
	assert.True(t, a.IsEnabledValue())

	a.RawValue = "true"
	assert.True(t, a.IsEnabledValue())

	a.RawValue = true
	assert.True(t, a.IsEnabledValue())

	a.RawValue = "1"
	assert.True(t, a.IsEnabledValue())

	a.RawValue = 1
	assert.True(t, a.IsEnabledValue())

	a.RawValue = "anything"
	assert.False(t, a.IsEnabledValue())
}

func TestAttributeIsDisabledValue(t *testing.T) {
	a := Attribute{RawValue: "No"}
	assert.True(t, a.IsDisabledValue())

	a.RawValue = "no"
	assert.True(t, a.IsDisabledValue())

	a.RawValue = "false"
	assert.True(t, a.IsDisabledValue())

	a.RawValue = false
	assert.True(t, a.IsDisabledValue())

	a.RawValue = "0"
	assert.True(t, a.IsDisabledValue())

	a.RawValue = 0
	assert.True(t, a.IsDisabledValue())

	a.RawValue = "anything"
	assert.False(t, a.IsDisabledValue())
}

func TestAttributeHasMultipleValues(t *testing.T) {
	a := Attribute{RawValue: "some string"}
	assert.False(t, a.HasMultipleValues())

	var rawValue []interface{}
	for _, val := range []string{"some", "string"} {
		rawValue = append(rawValue, val)
	}
	a.RawValue = rawValue

	assert.True(t, a.HasMultipleValues())
}

func TestAttributeValues(t *testing.T) {
	a := Attribute{RawValue: "some string"}
	result := a.Values()
	assert.IsType(t, []string{}, result)
	assert.Empty(t, result)

	var rawValue []interface{}
	for _, val := range []string{"some", "string"} {
		rawValue = append(rawValue, val)
	}
	a.RawValue = rawValue
	result = a.Values()
	assert.IsType(t, []string{}, result)
	assert.Len(t, result, 2)
}

func TestAttributeHasUnit(t *testing.T) {
	a := Attribute{}
	assert.False(t, a.HasUnitCode())

	a.UnitCode = "UNIT"
	assert.True(t, a.HasUnitCode())
}

func TestAttributeGetUnit(t *testing.T) {
	a := Attribute{}
	assert.IsType(t, Unit{}, a.GetUnit(), "Get an empty unit if not available")

	a.UnitCode = "PCS"
	assert.IsType(t, Unit{}, a.GetUnit(), "Get a unit if available")
	assert.Equal(t, "PCS", a.GetUnit().Code, "Fetched unit contains the correct content")
}

func TestBasicProductHasAttribute(t *testing.T) {
	b := BasicProductData{}
	assert.False(t, b.HasAttribute("code"))

	b.Attributes = map[string]Attribute{"code": Attribute{Code: "code"}}
	assert.True(t, b.HasAttribute("code"))
	assert.False(t, b.HasAttribute("Code"))
}

func TestBasicProductGetFinalPrice(t *testing.T) {
	p := PriceInfo{
		IsDiscounted: false,
		Discounted:   domain.NewFromFloat(0.99,"€"),
		Default:      domain.NewFromFloat(1.99, "€"),
	}
	assert.Equal(t, 1.99, p.GetFinalPrice().FloatAmount())

	p.IsDiscounted = true
	assert.Equal(t, 0.99, p.GetFinalPrice().FloatAmount())
}

func TestBasicProductGetMedia(t *testing.T) {
	var m []Media
	p := BasicProductData{Media: m}

	result := p.GetMedia("something")
	assert.IsType(t, Media{}, result, "Media returned on unknown type")
	assert.Empty(t, result.Usage, "empty media returned on unknown type")

	result = p.GetListMedia()
	assert.Empty(t, result.Usage, "empty list media returned if it does not exist")

	p.Media = append(m, Media{Usage: MediaUsageList})
	result = p.GetMedia(MediaUsageList)
	assert.Equal(t, "list", result.Usage, "media with correct usage returned if exists")

	result = p.GetListMedia()
	assert.Equal(t, "list", result.Usage, "list media returned if exists")
}

func TestIsSaleableNow(t *testing.T) {
	s := Saleable{}
	assert.False(t, s.IsSaleableNow(), "not saleable now if nothing is set")

	s.IsSaleable = true
	assert.True(t, s.IsSaleableNow(), "saleable now if just saleable")

	future := time.Now().Add(time.Hour)
	past := time.Now().Add(-time.Hour)
	s.SaleableFrom = future
	assert.False(t, s.IsSaleableNow(), "not saleable now if saleable from is in future")
	s.SaleableFrom = past
	assert.True(t, s.IsSaleableNow(), "saleable now if saleable from is in the past")

	s.SaleableTo = past
	assert.False(t, s.IsSaleableNow(), "not saleable now if saleable to is in the past")
	s.SaleableTo = future
	assert.True(t, s.IsSaleableNow(), "saleable now if saleable to is in the future")
}
