package domain

import (
	"math"
	"math/big"
	"testing"
	"time"

	"flamingo.me/flamingo-commerce/v3/price/domain"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
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
	assert.Len(t, result, 1)
	assert.Equal(t, "some string", result[0])

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
		Discounted:   domain.NewFromFloat(0.99, "€"),
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

func TestSaleable_GetLoyaltyChargeSplit(t *testing.T) {

	p := Saleable{
		ActivePrice: PriceInfo{
			//100€ value
			Default: priceDomain.NewFromInt(100, 1, "€"),
		},
		LoyaltyPrices: []LoyaltyPriceInfo{
			{
				Type:             "loyalty.miles",
				MaxPointsToSpent: new(big.Float).SetInt64(50),
				//10 is the minimum to pay in miles (=20€ value)
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				//50 miles == 100€ meaning 1Mile = 2€
				Default: priceDomain.NewFromInt(50, 1, "Miles"),
			},
		},
	}

	//Test default charges (the min price in points should be evaluated)
	charges := p.GetLoyaltyChargeSplit(nil, nil, 1)

	chargeLoyaltyMiles, found := charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(10, 1, "Miles"), chargeLoyaltyMiles.Price, "only minimum points expected")

	chargeMain, found := charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(80, 1, "€"), chargeMain.Price)

	//Test when we pass 15 miles as wish
	wished := NewWishedToPay().Add("loyalty.miles", priceDomain.NewFromInt(15, 1, "Miles"))
	charges = p.GetLoyaltyChargeSplit(nil, &wished, 1)
	chargeMain, found = charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)

	assert.Equal(t, priceDomain.NewFromInt(70, 1, "€"), chargeMain.Price)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(15, 1, "Miles"), chargeLoyaltyMiles.Price, "the whished 15 points expected")

	//Test when we pass 100 miles as wish
	wished = NewWishedToPay().Add("loyalty.miles", priceDomain.NewFromInt(100, 1, "Miles"))
	charges = p.GetLoyaltyChargeSplit(nil, &wished, 1)
	chargeMain, found = charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)

	assert.Equal(t, priceDomain.NewFromInt(0, 1, "€"), chargeMain.Price, "Main charge should be 0")

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(50, 1, "Miles"), chargeLoyaltyMiles.Price, "50 points expected as max")

	//Test when we pass 30 miles as desired payment and wish for qty 2
	wished = NewWishedToPay().Add("loyalty.miles", priceDomain.NewFromInt(30, 1, "Miles"))
	doublePrice := p.ActivePrice.GetFinalPrice().Multiply(2)
	charges = p.GetLoyaltyChargeSplit(&doublePrice, &wished, 1)
	chargeMain, found = charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)

	assert.Equal(t, priceDomain.NewFromInt(140, 1, "€"), chargeMain.Price)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(30, 1, "Miles"), chargeLoyaltyMiles.Price, "the whished 30 points expected")

}

func TestSaleable_GetLoyaltyChargeSplitWithAdjustedValue(t *testing.T) {

	p := Saleable{
		ActivePrice: PriceInfo{
			//100€ value
			Default: priceDomain.NewFromInt(100, 1, "€"),
		},
		LoyaltyPrices: []LoyaltyPriceInfo{
			{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				//10 is the minimum to pay in miles (=20€ value)
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				//50 miles == 100€ meaning 1Mile = 2€
				Default: priceDomain.NewFromInt(50, 1, "Miles"),
			},
		},
	}

	//we need to pay 150€ (e,g. because some tax are added)
	newValue := priceDomain.NewFromInt(150, 1, "€")
	charges := p.GetLoyaltyChargeSplit(&newValue, nil, 1)

	chargeLoyaltyMiles, found := charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(10, 1, "Miles"), chargeLoyaltyMiles.Price, "only minimum points expected")

	chargeMain, found := charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(130, 1, "€"), chargeMain.Price)

	//pay 150 - and want to spend less then min
	newValue = priceDomain.NewFromInt(150, 1, "€")
	wished := NewWishedToPay()
	wished.Add("loyalty.miles", priceDomain.NewFromInt(8, 1, "Miles"))
	charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 1)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(10, 1, "Miles"), chargeLoyaltyMiles.Price, "only minimum points expected")

	chargeMain, found = charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(130, 1, "€"), chargeMain.Price)

	//we need to pay 50€ (e,g. because some discounts are applied)
	newValue = priceDomain.NewFromInt(50, 1, "€")
	charges = p.GetLoyaltyChargeSplit(&newValue, nil, 1)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(10, 1, "Miles"), chargeLoyaltyMiles.Price, "only minimum points expected")

	chargeMain, found = charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(30, 1, "€"), chargeMain.Price)

	//we need to pay 150€ and wish to pay everything with miles
	newValue = priceDomain.NewFromInt(150, 1, "€")

	wished = NewWishedToPay()
	wished.Add("loyalty.miles", priceDomain.NewFromInt(200000, 1, "Miles"))
	charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 1)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(75, 1, "Miles"), chargeLoyaltyMiles.Price, "adjusted  points expected as charge (not more than total value)")

	chargeMain, found = charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(0, 1, "€"), chargeMain.Price)

}

func TestSaleable_GetLoyaltyChargeSplitCentRoundingCheck(t *testing.T) {

	p := Saleable{
		ActivePrice: PriceInfo{
			//100€ value
			Default: priceDomain.NewFromFloat(9.99, "€"),
		},
		LoyaltyPrices: []LoyaltyPriceInfo{
			{
				Type:             "loyalty.miles",
				MaxPointsToSpent: nil,
				//10 is the minimum to pay in miles (=20€ value)
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				//50 miles == 100€ meaning 1Mile = 2€
				Default: domain.NewFromInt(53, 1, "Miles"), // one mile = 5.305305305305305 €
			},
		},
	}

	wishedMax := NewWishedToPay()
	wishedMax.Add("loyalty.miles", priceDomain.NewFromInt(math.MaxInt64, 1, "Miles"))

	//107.06 would be 567.98 miles - so  we pay 567 miles (rounded floor always) we expect to pay everything in miles.
	expectedMilesMax := int64(567)
	newValue := priceDomain.NewFromFloat(107.06, "€")
	charges := p.GetLoyaltyChargeSplit(&newValue, &wishedMax, 1)
	chargeLoyaltyMiles, _ := charges.GetByType("loyalty.miles")
	assert.Equal(t, priceDomain.NewFromInt(expectedMilesMax, 1, "Miles").FloatAmount(), chargeLoyaltyMiles.Price.FloatAmount(), "107.06 expected to be 567 miles")

	wished := NewWishedToPay()
	wished.Add("loyalty.miles", priceDomain.NewFromInt(expectedMilesMax, 1, "Miles"))

	charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 1)

	chargeLoyaltyMiles, found := charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(567, 1, "Miles"), chargeLoyaltyMiles.Price, "adjusted  points expected as charge (not more than total value)")
	assert.Equal(t, 107.06, chargeLoyaltyMiles.Value.FloatAmount(), "adjusted  points expected as charge (not more than total value)")

	chargeMain, found := charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(0, 1, "€"), chargeMain.Price)

	//106.89 would be 567.084084084084084 miles - so  we also pay 567 miles (rounded) we expect to pay everything in miles.
	newValue = priceDomain.NewFromFloat(106.89, "€")

	wished = NewWishedToPay()
	wished.Add("loyalty.miles", priceDomain.NewFromInt(567, 1, "Miles"))
	charges = p.GetLoyaltyChargeSplit(&newValue, &wished, 1)

	chargeLoyaltyMiles, found = charges.GetByType("loyalty.miles")
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(567, 1, "Miles"), chargeLoyaltyMiles.Price, "adjusted  points expected as charge (not more than total value)")
	assert.Equal(t, 106.89, chargeLoyaltyMiles.Value.FloatAmount(), "adjusted  points expected as charge (not more than total value)")

	chargeMain, found = charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(0, 1, "€"), chargeMain.Price)

}

func TestSaleable_GetLoyaltyChargeSplitIgnoreMin(t *testing.T) {

	p := Saleable{
		ActivePrice: PriceInfo{
			//100€ value
			Default: priceDomain.NewFromInt(100, 1, "€"),
		},
		LoyaltyPrices: []LoyaltyPriceInfo{
			{
				Type:             "loyalty.miles",
				MaxPointsToSpent: new(big.Float).SetInt64(50),
				//10 is the minimum to pay in miles (=20€ value)
				MinPointsToSpent: *new(big.Float).SetInt64(10),
				//50 miles == 100€ meaning 1Mile = 2€
				Default: priceDomain.NewFromInt(50, 1, "Miles"),
			},
		},
	}

	//Test default charges (the min price in points should be evaluated)
	charges := p.GetLoyaltyChargeSplitIgnoreMin(nil, nil, 1)

	_, found := charges.GetByType("loyalty.miles")
	assert.False(t, found)

	chargeMain, found := charges.GetByType(priceDomain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, priceDomain.NewFromInt(100, 1, "€"), chargeMain.Price)
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
					Default: priceDomain.NewFromFloat(23.23, "NZD"),
				},
				{
					Type:    "TheOtherThing",
					Default: priceDomain.NewFromFloat(24.24, "NZD"),
				},
			},
			leType:   "MilesAndMore",
			wantBool: true,
			wantEarning: &LoyaltyEarningInfo{
				Type:    "MilesAndMore",
				Default: priceDomain.NewFromFloat(23.23, "NZD"),
			},
		},
		{
			name: "no matching loyalty info",
			loyaltyEarnings: []LoyaltyEarningInfo{
				{
					Type:    "MilesAndMoreX",
					Default: priceDomain.NewFromFloat(23.23, "NZD"),
				},
				{
					Type:    "TheOtherThing",
					Default: priceDomain.NewFromFloat(24.24, "NZD"),
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
