package domain

import (
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo-commerce/v3/price/domain"
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
	t.Run("[]interface{} raw value", func(t *testing.T) {
		a := Attribute{RawValue: "some string"}
		assert.False(t, a.HasMultipleValues())

		var rawValue []interface{}
		for _, val := range []string{"some", "string"} {
			rawValue = append(rawValue, val)
		}
		a.RawValue = rawValue

		assert.True(t, a.HasMultipleValues())
	})

	t.Run("Translated multi value", func(t *testing.T) {
		a := Attribute{RawValue: []Attribute{{Code: "foo"}}}
		assert.True(t, a.HasMultipleValues())
	})
}

func TestAttributeValues(t *testing.T) {
	t.Run("string values", func(t *testing.T) {
		a := Attribute{RawValue: "some string"}
		result := a.Values()
		assert.IsType(t, []string{}, result)
		assert.Len(t, result, 0)

		var rawValue []interface{}
		for _, val := range []string{"some", "  string    "} {
			rawValue = append(rawValue, val)
		}
		a.RawValue = rawValue
		result = a.Values()
		assert.IsType(t, []string{}, result)
		assert.Len(t, result, 2)
		assert.Equal(t, "some", result[0])
		assert.Equal(t, "string", result[1])
	})

	t.Run("translated values", func(t *testing.T) {
		a := Attribute{RawValue: []Attribute{{Label: "translation-A", RawValue: "raw-A"}, {Label: "translation-B", RawValue: "raw-B"}}}
		values := a.Values()
		assert.Len(t, values, 2)
		assert.Equal(t, "raw-A", values[0])
		assert.Equal(t, "raw-B", values[1])
	})
}

func TestAttributeLabels(t *testing.T) {
	t.Run("translated values", func(t *testing.T) {
		a := Attribute{RawValue: []Attribute{{Label: "translation-A", RawValue: "raw-A"}, {Label: "translation-B", RawValue: "raw-B"}}}
		labels := a.Labels()
		assert.Len(t, labels, 2)
		assert.Equal(t, "translation-A", labels[0])
		assert.Equal(t, "translation-B", labels[1])
	})

	t.Run("no translated values will fallback to raw values", func(t *testing.T) {
		a := Attribute{}
		var rawValue []interface{}
		for _, val := range []string{"raw-1", "raw-2"} {
			rawValue = append(rawValue, val)
		}
		a.RawValue = rawValue

		labels := a.Labels()
		assert.Len(t, labels, 2)
		assert.Equal(t, "raw-1", labels[0])
		assert.Equal(t, "raw-2", labels[1])
	})
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

	b.Attributes = map[string]Attribute{"code": {Code: "code"}}
	assert.True(t, b.HasAttribute("code"))
	assert.False(t, b.HasAttribute("Code"))
}

func TestBasicProductHasAllAttributes(t *testing.T) {
	b := BasicProductData{}
	assert.False(t, b.HasAllAttributes([]string{"code", "color"}))

	b.Attributes = map[string]Attribute{"code": {Code: "code"}, "color": {Code: "color"}}
	assert.True(t, b.HasAllAttributes([]string{"code", "color"}))
	assert.False(t, b.HasAllAttributes([]string{"Code", "Color"}))
}

func TestBasicProductHasGetAttributesByKey(t *testing.T) {
	b := BasicProductData{}
	b.Attributes = map[string]Attribute{
		"foo": {Code: "foo"},
		"bar": {Code: "bar"},
	}
	assert.Equal(t, []Attribute{
		{Code: "foo"},
		{Code: "bar"},
	}, b.Attributes.AttributesByKey([]string{"foo", "bar"}))

	assert.Equal(t, []Attribute{
		{Code: "bar"},
		{Code: "foo"},
	}, b.Attributes.AttributesByKey([]string{"bar", "foo"}))

	assert.Equal(t, []Attribute{
		{Code: "foo"},
		{Code: "bar"},
	}, b.Attributes.AttributesByKey([]string{"foo", "baz", "bar"}))
}

func TestBasicProductGetFinalPrice(t *testing.T) {
	p := PriceInfo{
		IsDiscounted: false,
		Discounted:   domain.NewFromFloat(0.99, "EUR"),
		Default:      domain.NewFromFloat(1.99, "EUR"),
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

func TestBasicProductBadgesGetFirst(t *testing.T) {
	t.Parallel()
	var badges Badges

	assert.Nil(t, badges.First(), "get nil if badges are nil - don't fail")

	badges = Badges{}
	assert.Nil(t, badges.First(), "get nil if badges are empty")

	badges = Badges{
		{
			Code:  "first",
			Label: "First",
		},
		{
			Code:  "second",
			Label: "Second",
		},
	}
	assert.Equal(t, &Badge{Code: "first", Label: "First"}, badges.First(), "get the first badge")
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

func TestSaleable_GetLoyaltyChargeSplit(t *testing.T) {
	t.Parallel()

	t.Run("get correct charges", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				// 100EUR value
				Default: domain.NewFromInt(100, 1, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: new(big.Float).SetInt64(50),
				// 10 is the minimum to pay in miles (=20EUR value)
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				// 50 miles == 100EUR meaning 1Mile = 2EUR
				Default: domain.NewFromInt(50, 1, "Miles"),
			},
		}

		// Test default charges (the min price in points should be evaluated)
		charges := p.GetLoyaltyChargeSplit(nil, nil, 1)

		chargeLoyaltyMiles, found := charges.GetByType("loyalty.miles")
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(10, 1, "Miles"), chargeLoyaltyMiles.Price, "only minimum points expected")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(80, 1, "EUR"), chargeMain.Price)

		// Test when we pass 15 miles as wish
		wished := NewWishedToPay().Add("loyalty.miles", domain.NewFromInt(15, 1, "Miles"))
		charges = p.GetLoyaltyChargeSplit(nil, &wished, 1)
		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)

		assert.Equal(t, domain.NewFromInt(70, 1, "EUR"), chargeMain.Price)

		chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(15, 1, "Miles"), chargeLoyaltyMiles.Price, "the whished 15 points expected")

		// Test when we pass 100 miles as wish
		wished = NewWishedToPay().Add("loyalty.miles", domain.NewFromInt(100, 1, "Miles"))
		charges = p.GetLoyaltyChargeSplit(nil, &wished, 1)
		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)

		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "Main charge should be 0")

		chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(50, 1, "Miles"), chargeLoyaltyMiles.Price, "50 points expected as max")

		// Test when we pass 30 miles as desired payment and wish for qty 2
		wished = NewWishedToPay().Add("loyalty.miles", domain.NewFromInt(30, 1, "Miles"))
		doublePrice := p.ActivePrice.GetFinalPrice().Multiply(2)
		charges = p.GetLoyaltyChargeSplit(&doublePrice, &wished, 1)
		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)

		assert.Equal(t, domain.NewFromInt(140, 1, "EUR"), chargeMain.Price)

		chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(30, 1, "Miles"), chargeLoyaltyMiles.Price, "the whished 30 points expected")
	})

	t.Run("empty charges, when active loyalty is nil", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				// 100EUR value
				Default: domain.NewFromInt(100, 1, "EUR"),
			},
			ActiveLoyaltyPrice: nil,
		}

		// Test default charges (the min price in points should be evaluated)
		charges := p.GetLoyaltyChargeSplit(nil, nil, 1)

		chargeLoyaltyMiles, found := charges.GetByType("loyalty.miles")
		assert.False(t, found)
		assert.Equal(t, domain.Price{}, chargeLoyaltyMiles.Price)

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(100, 1, "EUR"), chargeMain.Price)

		// Test when we pass 15 miles as wish
		wished := NewWishedToPay().Add("loyalty.miles", domain.NewFromInt(15, 1, "Miles"))
		charges = p.GetLoyaltyChargeSplit(nil, &wished, 1)
		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)

		assert.Equal(t, domain.NewFromInt(100, 1, "EUR"), chargeMain.Price)

		chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
		assert.False(t, found)
		assert.Equal(t, domain.Price{}, chargeLoyaltyMiles.Price)

		// Test when we pass 100 miles as wish
		wished = NewWishedToPay().Add("loyalty.miles", domain.NewFromInt(100, 1, "Miles"))
		charges = p.GetLoyaltyChargeSplit(nil, &wished, 1)
		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)

		assert.Equal(t, domain.NewFromInt(100, 1, "EUR"), chargeMain.Price)

		chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
		assert.False(t, found)
		assert.Equal(t, domain.Price{}, chargeLoyaltyMiles.Price)

		// Test when we pass 30 miles as desired payment and wish for qty 2
		wished = NewWishedToPay().Add("loyalty.miles", domain.NewFromInt(30, 1, "Miles"))
		doublePrice := p.ActivePrice.GetFinalPrice().Multiply(2)
		charges = p.GetLoyaltyChargeSplit(&doublePrice, &wished, 1)
		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)

		assert.Equal(t, domain.NewFromInt(200, 1, "EUR"), chargeMain.Price)

		chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
		assert.False(t, found)
		assert.Equal(t, domain.Price{}, chargeLoyaltyMiles.Price)
	})
}

func TestSaleable_GetLoyaltyChargeSplitWithAdjustedValue(t *testing.T) {

	p := Saleable{
		ActivePrice: PriceInfo{
			// 100EUR value
			Default: domain.NewFromInt(100, 1, "EUR"),
		},
		ActiveLoyaltyPrice: &LoyaltyPriceInfo{
			Type:             "loyalty.miles",
			MaxPointsToSpent: nil,
			// 10 is the minimum to pay in miles (=20EUR value)
			MinPointsToSpent: *new(big.Float).SetInt64(10),
			// 50 miles == 100EUR meaning 1Mile = 2EUR
			Default: domain.NewFromInt(50, 1, "Miles"),
		},
	}

	// we need to pay 150EUR (e,g. because some tax are added)
	newValue := domain.NewFromInt(150, 1, "EUR")
	charges := p.GetLoyaltyChargeSplit(&newValue, nil, 1)

	chargeLoyaltyMiles, found := charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(10, 1, "Miles"), chargeLoyaltyMiles.Price, "only minimum points expected")

	chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(130, 1, "EUR"), chargeMain.Price)

	// pay 150 - and want to spend less then min
	newValue = domain.NewFromInt(150, 1, "EUR")
	wished := NewWishedToPay()
	wished.Add("loyalty.miles", domain.NewFromInt(8, 1, "Miles"))
	charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 1)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(10, 1, "Miles"), chargeLoyaltyMiles.Price, "only minimum points expected")

	chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(130, 1, "EUR"), chargeMain.Price)

	// we need to pay 50EUR (e,g. because some discounts are applied)
	newValue = domain.NewFromInt(50, 1, "EUR")
	charges = p.GetLoyaltyChargeSplit(&newValue, nil, 1)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(10, 1, "Miles"), chargeLoyaltyMiles.Price, "only minimum points expected")

	chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(30, 1, "EUR"), chargeMain.Price)

	// we need to pay 150EUR and wish to pay everything with miles
	newValue = domain.NewFromInt(150, 1, "EUR")

	wished = NewWishedToPay()
	wished.Add("loyalty.miles", domain.NewFromInt(200000, 1, "Miles"))
	charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 1)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(75, 1, "Miles"), chargeLoyaltyMiles.Price, "adjusted  points expected as charge (not more than total value)")

	chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price)

}

func TestSaleable_GetLoyaltyChargeSplitCentRoundingCheck(t *testing.T) {
	t.Parallel()

	t.Run("9.99 EUR and 53 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				// 100EUR value
				Default: domain.NewFromFloat(9.99, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				// 10 is the minimum to pay in miles (=20EUR value)
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				// 50 miles == 100EUR meaning 1Mile = 2EUR
				Default: domain.NewFromInt(53, 1, "Miles"), // one mile = 5.305305305305305 EUR
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		// 107.06 would be 567.98 miles - so  we pay 568 miles we expect to pay everything in miles.
		expectedMilesMax := int64(568)
		newValue := domain.NewFromFloat(107.06, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "107.06 expected to be 568 miles")

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewFromInt(expectedMilesMax, 1, "Miles"))

		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 1)

		chargeLoyaltyMiles, found := charges.GetByType("loyalty.miles")
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(568, 1, "Miles"), chargeLoyaltyMiles.Price, "adjusted  points expected as charge (not more than total value)")
		assert.Equal(t, 107.06, chargeLoyaltyMiles.Value.FloatAmount(), "adjusted  points expected as charge (not more than total value)")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price)

		// 106.89 would be 567.084084084084084 miles - so  we also pay 567 miles (rounded) we expect to pay everything in miles.
		newValue = domain.NewFromFloat(106.89, "EUR")

		wished = NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewFromInt(567, 1, "Miles"))
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 1)

		chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(567, 1, "Miles"), chargeLoyaltyMiles.Price, "adjusted  points expected as charge (not more than total value)")
		assert.Equal(t, 106.89, chargeLoyaltyMiles.Value.FloatAmount(), "adjusted  points expected as charge (not more than total value)")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price)
	})

	t.Run("10 items by 9.50 EUR and 53 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(9.50, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				Default:          domain.NewFromInt(53, 1, "Miles"), // one mile = 5.57894737 EUR
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		// 107.06 would be 567.98 miles - so  we pay 568 miles we expect to pay everything in miles.
		expectedMilesMax := int64(530)
		newValue := domain.NewFromFloat(95, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 10)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "95 expected to be 530 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price)
	})

	t.Run("10 items by 9.49 EUR and 53 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(9.49, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				Default:          domain.NewFromInt(53, 1, "Miles"), // one mile = 5.58482613 EUR
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		expectedMilesMax := int64(530)
		newValue := domain.NewFromFloat(94.9, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 10)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "94.9 expected to be 530 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price)
	})

	t.Run("item for 1.50 EUR or 450 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.50, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(30),
				Default:          domain.NewFromInt(450, 1, "Miles"),
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		expectedMilesMax := int64(450)
		newValue := domain.NewFromFloat(1.5, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1.5 expected to be 450 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(300)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1 expected to be 300 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(900)
		newValue = domain.NewFromFloat(3.00, "EUR") // reduced price still calculated perfectly
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "3 expected to be 900 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMiles := int64(90)

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewZero("Miles"))

		newValue = domain.NewFromFloat(4.50, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMiles, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "you cannot pay less then 90 Miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(4.200000123, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should de left to pay 4.20, but got %f", chargeMain.Price.FloatAmount()))
	})

	t.Run("item for 1.49 EUR or 447 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.49, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(30),
				Default:          domain.NewFromInt(447, 1, "Miles"),
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		// 1.49 would be 447 miles - so  we pay 447 miles we expect to pay everything in miles.
		expectedMilesMax := int64(447)
		newValue := domain.NewFromFloat(1.49, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1.49 expected to be 447 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(300)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1 expected to be 300 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(1341)
		newValue = domain.NewFromFloat(4.47, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "4.47 expected to be 1341 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMiles := int64(90)

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewZero("Miles"))

		newValue = domain.NewFromFloat(4.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMiles, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "you cannot pay less then 90 Miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(3.70000123, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should de left to pay 3.7, but got %f", chargeMain.Price.FloatAmount()))
	})

	t.Run("item for 1.01 EUR or 303 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.01, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(30),
				Default:          domain.NewFromInt(303, 1, "Miles"),
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		// 1.01 would be 303 miles - so  we pay 303 miles we expect to pay everything in miles.
		expectedMilesMax := int64(303)
		newValue := domain.NewFromFloat(1.01, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1.01 expected to be 303 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		// 1.0 would be 300 miles - so  we pay 300 miles we expect to pay everything in miles.
		expectedMilesMax = int64(300)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1 expected to be 300 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(909)
		newValue = domain.NewFromFloat(3.03, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "3.03 expected to be 909 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMiles := int64(90)

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewZero("Miles"))

		newValue = domain.NewFromFloat(3.03, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMiles, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "you cannot pay less then 90 Miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(2.729999999, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should de left to pay 2.73, but got %f", chargeMain.Price.FloatAmount()))
	})

	t.Run("item for 1.99 EUR or 597 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.99, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(30),
				Default:          domain.NewFromInt(597, 1, "Miles"), // 300 miles = 1 EUR
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		// 1.99 would be 597 miles - so  we pay 597 miles we expect to pay everything in miles.
		expectedMilesMax := int64(597)
		newValue := domain.NewFromFloat(1.99, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1.99 expected to be 597 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(300)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1 expected to be 300 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(1791)
		newValue = domain.NewFromFloat(5.97, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "5.97 expected to be 1791 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMiles := int64(90)

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewZero("Miles"))

		newValue = domain.NewFromFloat(5.97, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMiles, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "you cannot pay less then 90 Miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(5.6666666, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should de left to pay 5.67, but got %f", chargeMain.Price.FloatAmount()))
	})

	t.Run("item for 1.50 EUR or 450 Miles with max 90", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.50, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: new(big.Float).SetInt64(90),
				MinPointsToSpent: *new(big.Float).SetInt64(30),
				Default:          domain.NewFromInt(450, 1, "Miles"),
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles")) // we wish to pay everything in miles

		expectedMilesMax := int64(90)
		newValue := domain.NewFromFloat(1.5, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount())

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromFloat(1.20, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should be left to pay 1.20 but got %f", chargeMain.Price.FloatAmount()))

		expectedMilesMax = int64(90)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount())

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(0.70, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should be left to pay 1.20 but got %f", chargeMain.Price.FloatAmount()))

		expectedMilesMax = int64(270)
		newValue = domain.NewFromFloat(3.00, "EUR") // reduced price still calculated perfectly
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount())

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(2.10, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should be left to pay 2.10 but got %f", chargeMain.Price.FloatAmount()))
	})

	t.Run("item for 1.50 EUR or 2 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.50, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(1),
				//Default:          domain.NewFromInt(2, 1, "Miles"), // Should come unrounded!!!
				Default: domain.NewFromFloat(1.5, "Miles"),
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		expectedMilesMax := int64(2)
		newValue := domain.NewFromFloat(1.5, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1.5 expected to be 2 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(1)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1 expected to be 1 mile")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(5)
		newValue = domain.NewFromFloat(4.50, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "4.50 expected to be 5 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMiles := int64(3)

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewZero("Miles"))

		newValue = domain.NewFromFloat(4.10, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMiles, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "you cannot pay less then 3 Miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(1.10, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should de left to pay 1.10, but got %f", chargeMain.Price.FloatAmount()))
	})

	t.Run("item for 1.49 EUR or 1 Mile", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.49, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(1),
				//Default:          domain.NewFromInt(447, 1, "Miles"), // Should come unrounded!!!
				Default: domain.NewFromFloat(1.49, "Miles"),
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		expectedMilesMax := int64(1)
		newValue := domain.NewFromFloat(1.49, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1.49 expected to be 1 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(1)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1 expected to be 300 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(4)
		newValue = domain.NewFromFloat(4.47, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "4.47 expected to be 4 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMiles := int64(3)

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewZero("Miles"))

		newValue = domain.NewFromFloat(3.47, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMiles, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "you cannot pay less then 3 Miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewZero("EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("we should gift 0.47 cents in this case, but got %f to pay", chargeMain.Price.FloatAmount()))
	})

	t.Run("item for 1.01 EUR or 1 Mile", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.01, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(1),
				//Default:          domain.NewFromInt(303, 1, "Miles"), // Should come unrounded!!!
				Default: domain.NewFromFloat(1.01, "Miles"),
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		expectedMilesMax := int64(1)
		newValue := domain.NewFromFloat(1.01, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1.01 expected to be 1 mile")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(1)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1 expected to be 1 mile")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(3)
		newValue = domain.NewFromFloat(3.03, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "3.03 expected to be 3 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMiles := int64(3)

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewZero("Miles"))

		newValue = domain.NewFromFloat(3.50, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMiles, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "you cannot pay less then 3 Miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(0.50, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should de left to pay 0.50, but got %f", chargeMain.Price.FloatAmount()))
	})

	t.Run("item for 1.99 EUR or 2 Miles", func(t *testing.T) {
		t.Parallel()

		p := Saleable{
			ActivePrice: PriceInfo{
				Default: domain.NewFromFloat(1.99, "EUR"),
			},
			ActiveLoyaltyPrice: &LoyaltyPriceInfo{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				MinPointsToSpent: *new(big.Float).SetInt64(1),
				Default:          domain.NewFromFloat(1.99, "Miles"),
			},
		}

		wishedMax := NewWishedToPay()
		wishedMax.Add("loyalty.miles", domain.NewFromInt(math.MaxInt64, 1, "Miles"))

		expectedMilesMax := int64(2)
		newValue := domain.NewFromFloat(1.99, "EUR")
		charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1.01 expected to be 2 miles")

		chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
		require.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(1)
		newValue = domain.NewFromFloat(1.00, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "1 expected to be 1 mile")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMilesMax = int64(6)
		newValue = domain.NewFromFloat(5.97, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "5.97 expected to be 6 miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromInt(0, 1, "EUR"), chargeMain.Price, "should be nothing left to pay")

		expectedMiles := int64(3)

		wished := NewWishedToPay()
		wished.Add("loyalty.miles", domain.NewZero("Miles"))

		newValue = domain.NewFromFloat(5.97, "EUR")
		charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 3)
		chargeLoyaltyMiles, _ = charges.GetByType("loyalty.miles")
		assert.Equal(t, domain.NewFromInt(expectedMiles, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "you cannot pay less then 3 Miles")

		chargeMain, found = charges.GetByType(domain.ChargeTypeMain)
		assert.True(t, found)
		assert.Equal(t, domain.NewFromFloat(2.97, "EUR").GetPayable(), chargeMain.Price.GetPayable(), fmt.Sprintf("should de left to pay 2.97, but got %f", chargeMain.Price.FloatAmount()))
	})
}

func TestSaleable_GetLoyaltyChargeSplitIgnoreMin(t *testing.T) {

	p := Saleable{
		ActivePrice: PriceInfo{
			// 100EUR value
			Default: domain.NewFromInt(100, 1, "EUR"),
		},
		LoyaltyPrices: []LoyaltyPriceInfo{
			{
				Type:             "loyalty.miles",
				MaxPointsToSpent: new(big.Float).SetInt64(50),
				// 10 is the minimum to pay in miles (=20EUR value)
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				// 50 miles == 100EUR meaning 1Mile = 2EUR
				Default: domain.NewFromInt(50, 1, "Miles"),
			},
		},
	}

	// Test default charges (the min price in points should be evaluated)
	charges := p.GetLoyaltyChargeSplitIgnoreMin(nil, nil, 1)

	_, found := charges.GetByType("loyalty.miles")
	assert.False(t, found)

	chargeMain, found := charges.GetByType(domain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, domain.NewFromInt(100, 1, "EUR"), chargeMain.Price)
}

func TestSaleable_GetLoyaltyEarningByType(t *testing.T) {

	tests := []struct {
		name            string
		loyaltyEarnings []LoyaltyEarningInfo
		leType          string
		wantBool        bool
		wantEarning     *LoyaltyEarningInfo
	}{
		{
			name:            "empty loyalty info",
			loyaltyEarnings: nil,
			leType:          "dontCare",
			wantBool:        false,
			wantEarning:     nil,
		},
		{
			name: "matching loyalty info",
			loyaltyEarnings: []LoyaltyEarningInfo{
				{
					Type:    "MilesAndMore",
					Default: domain.NewFromFloat(23.23, "NZD"),
				},
				{
					Type:    "TheOtherThing",
					Default: domain.NewFromFloat(24.24, "NZD"),
				},
			},
			leType:   "MilesAndMore",
			wantBool: true,
			wantEarning: &LoyaltyEarningInfo{
				Type:    "MilesAndMore",
				Default: domain.NewFromFloat(23.23, "NZD"),
			},
		},
		{
			name: "no matching loyalty info",
			loyaltyEarnings: []LoyaltyEarningInfo{
				{
					Type:    "MilesAndMoreX",
					Default: domain.NewFromFloat(23.23, "NZD"),
				},
				{
					Type:    "TheOtherThing",
					Default: domain.NewFromFloat(24.24, "NZD"),
				},
			},
			leType:      "MilesAndMore",
			wantBool:    false,
			wantEarning: nil,
		},
	}

	for _, tt := range tests {
		saleable := new(Saleable)
		saleable.LoyaltyEarnings = tt.loyaltyEarnings

		resultEarning, resultBool := saleable.GetLoyaltyEarningByType(tt.leType)
		assert.Equal(t, tt.wantBool, resultBool, tt.name)
		assert.Equal(t, tt.wantEarning, resultEarning, tt.name)
	}

}

func TestBasicProductData_IsInStockForDeliveryCode(t *testing.T) {
	t.Parallel()

	t.Run("when in stock for delivery code then return true", func(t *testing.T) {
		t.Parallel()

		product := BasicProductData{
			Stock: []Stock{
				{
					InStock:      true,
					DeliveryCode: "test",
				},
				{
					InStock:      false,
					DeliveryCode: "not this one",
				},
			},
		}

		result := product.IsInStockForDeliveryCode("test")

		assert.True(t, result)
	})

	t.Run("when not in stock for delivery code then return false", func(t *testing.T) {
		t.Parallel()

		product := BasicProductData{
			Stock: []Stock{
				{
					InStock:      true,
					DeliveryCode: "not this one",
				},
				{
					InStock:      false,
					DeliveryCode: "test",
				},
			},
		}

		result := product.IsInStockForDeliveryCode("test")

		assert.False(t, result)
	})

	t.Run("when delivery code not found return false", func(t *testing.T) {
		t.Parallel()

		product := BasicProductData{
			Stock: []Stock{
				{
					InStock:      true,
					DeliveryCode: "not this one",
				},
				{
					InStock:      false,
					DeliveryCode: "nope",
				},
			},
		}

		result := product.IsInStockForDeliveryCode("test")

		assert.False(t, result)
	})
}

func TestSaleable_GetLoyaltyPriceByType(t *testing.T) {
	t.Parallel()

	t.Run("returns loyalty from active loyalty price", func(t *testing.T) {
		t.Parallel()

		activeLoyaltyPrice := LoyaltyPriceInfo{
			Type:    "valid_type",
			Default: domain.NewFromFloat(10.00, "LOYALTY"),
		}

		availablePrice1 := LoyaltyPriceInfo{
			Type:    "valid_type",
			Default: domain.NewFromFloat(5.00, "LOYALTY"),
		}

		availablePrice2 := LoyaltyPriceInfo{
			Type:    "invalid_type",
			Default: domain.NewFromFloat(100.00, "LOYALTY"),
		}

		saleable := Saleable{
			ActiveLoyaltyPrice: &activeLoyaltyPrice,
			LoyaltyPrices: []LoyaltyPriceInfo{
				availablePrice1,
				availablePrice2,
			},
		}

		resultPrice, found := saleable.GetLoyaltyPriceByType("valid_type")
		require.True(t, found)
		assert.Equal(t, activeLoyaltyPrice, *resultPrice)
	})

	t.Run("returns loyalty from available loyalty prices", func(t *testing.T) {
		t.Parallel()

		activeLoyaltyPrice := LoyaltyPriceInfo{
			Type:    "invalid_type",
			Default: domain.NewFromFloat(10.00, "LOYALTY"),
		}

		availablePrice1 := LoyaltyPriceInfo{
			Type:    "valid_type",
			Default: domain.NewFromFloat(5.00, "LOYALTY"),
		}

		availablePrice2 := LoyaltyPriceInfo{
			Type:    "invalid_type",
			Default: domain.NewFromFloat(100.00, "LOYALTY"),
		}

		saleable := Saleable{
			ActiveLoyaltyPrice: &activeLoyaltyPrice,
			LoyaltyPrices: []LoyaltyPriceInfo{
				availablePrice1,
				availablePrice2,
			},
		}

		resultPrice, found := saleable.GetLoyaltyPriceByType("valid_type")
		require.True(t, found)
		assert.Equal(t, availablePrice1, *resultPrice)
	})

	t.Run("returns loyalty from available loyalty prices when active is nil", func(t *testing.T) {
		t.Parallel()

		availablePrice1 := LoyaltyPriceInfo{
			Type:    "valid_type",
			Default: domain.NewFromFloat(5.00, "LOYALTY"),
		}

		availablePrice2 := LoyaltyPriceInfo{
			Type:    "invalid_type",
			Default: domain.NewFromFloat(100.00, "LOYALTY"),
		}

		saleable := Saleable{
			ActiveLoyaltyPrice: nil,
			LoyaltyPrices: []LoyaltyPriceInfo{
				availablePrice1,
				availablePrice2,
			},
		}

		resultPrice, found := saleable.GetLoyaltyPriceByType("valid_type")
		require.True(t, found)
		assert.Equal(t, availablePrice1, *resultPrice)
	})

	t.Run("returns loyalty from available loyalty prices", func(t *testing.T) {
		t.Parallel()

		activeLoyaltyPrice := LoyaltyPriceInfo{
			Type:    "invalid_type",
			Default: domain.NewFromFloat(10.00, "LOYALTY"),
		}

		availablePrice1 := LoyaltyPriceInfo{
			Type:    "invalid_type",
			Default: domain.NewFromFloat(5.00, "LOYALTY"),
		}

		availablePrice2 := LoyaltyPriceInfo{
			Type:    "invalid_type",
			Default: domain.NewFromFloat(100.00, "LOYALTY"),
		}

		saleable := Saleable{
			ActiveLoyaltyPrice: &activeLoyaltyPrice,
			LoyaltyPrices: []LoyaltyPriceInfo{
				availablePrice1,
				availablePrice2,
			},
		}

		resultPrice, found := saleable.GetLoyaltyPriceByType("valid_type")
		require.False(t, found)
		assert.Nil(t, resultPrice)
	})
}
