package domain_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"testing"

	"flamingo.me/flamingo-commerce/v3/price/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrice_IsLessThen(t *testing.T) {
	type fields struct {
		Amount   float64
		Currency string
	}
	type args struct {
		amount big.Float
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "simple is less",
			fields: fields{
				Amount: 11.0,
			},
			args: args{
				amount: *big.NewFloat(12.2),
			},
			want: true,
		},
		{
			name: "simple is not less",
			fields: fields{
				Amount: 13.0,
			},
			args: args{
				amount: *big.NewFloat(12.2),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := domain.NewFromFloat(tt.fields.Amount, tt.fields.Currency)

			if got := p.IsLessThenValue(tt.args.amount); got != tt.want {
				t.Errorf("Amount.IsLessThen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrice_Multiply(t *testing.T) {
	p := domain.NewFromFloat(2.5, "EUR")
	resultPrice := p.Multiply(3)
	assert.Equal(t, domain.NewFromFloat(7.5, "EUR").GetPayable().Amount(), resultPrice.GetPayable().Amount())
}

func TestPrice_GetPayable(t *testing.T) {
	price := domain.NewFromFloat(12.34567, "EUR")

	payable := price.GetPayable()
	assert.Equal(t, float64(12.35), payable.FloatAmount())

	price = domain.NewFromFloat(math.MaxInt64, "EUR").GetPayable()
	assert.Equal(t, float64(math.MaxInt64), price.FloatAmount())
}

func TestNewFromInt(t *testing.T) {
	price1 := domain.NewFromInt(1245, 100, "EUR")
	price2 := domain.NewFromFloat(12.45, "EUR")
	assert.Equal(t, price2.GetPayable().Amount(), price1.GetPayable().Amount())
	pricePayable := price1.GetPayable()
	assert.True(t, price2.GetPayable().Equal(pricePayable))
}

func TestPrice_SplitInPayables(t *testing.T) {
	originalPrice := domain.NewFromFloat(32.1, "EUR") // float edge case
	payableSplitPrices, _ := originalPrice.SplitInPayables(1)

	sumPrice := domain.NewZero("EUR")
	for _, price := range payableSplitPrices {
		sumPrice, _ = sumPrice.Add(price)
	}
	// sum of the splitted payable need to match original price payable
	assert.Equal(t, originalPrice.GetPayable().Amount(), sumPrice.GetPayable().Amount())

	originalPrice = domain.NewFromFloat(12.456, "EUR")
	payableSplitPrices, _ = originalPrice.SplitInPayables(6)

	sumPrice = domain.NewZero("EUR")
	for _, price := range payableSplitPrices {
		sumPrice, _ = sumPrice.Add(price)
	}
	// sum of the splitted payable need to match original price payable
	assert.Equal(t, originalPrice.GetPayable().Amount(), sumPrice.GetPayable().Amount())

	// edge case for negative input (happens when discounts are split)
	originalPrice = domain.NewFromFloat(-152.99, "EUR")
	payableSplitPrices, _ = originalPrice.SplitInPayables(3)

	sumPrice = domain.NewZero("EUR")
	for _, price := range payableSplitPrices {
		sumPrice, _ = sumPrice.Add(price)
	}
	assert.Equal(t, originalPrice.GetPayable().Amount(), sumPrice.GetPayable().Amount())
}

func TestPrice_Discounted(t *testing.T) {
	originalPrice := domain.NewFromFloat(12.45, "EUR")
	discountedPrice := originalPrice.Discounted(10).GetPayable()
	// 10% of - expected rounded value of 11.21
	assert.Equal(t, domain.NewFromInt(1121, 100, "").Amount(), discountedPrice.Amount())
}

func TestPrice_IsZero(t *testing.T) {
	var price domain.Price
	assert.Equal(t, domain.NewZero("").Amount(), price.GetPayable().Amount())
}

func TestSumAll(t *testing.T) {
	price1 := domain.NewFromInt(1200, 100, "EUR")
	price2 := domain.NewFromInt(1200, 100, "EUR")
	price3 := domain.NewFromInt(1200, 100, "EUR")

	result, err := domain.SumAll(price1, price2, price3)
	assert.NoError(t, err)
	assert.Equal(t, result, domain.NewFromInt(3600, 100, "EUR"))

}

func TestPrice_TaxFromGross(t *testing.T) {
	// 119 €
	price := domain.NewFromInt(119, 1, "EUR")
	tax := price.TaxFromGross(*new(big.Float).SetInt64(19))
	assert.Equal(t, tax, domain.NewFromInt(19, 1, "EUR"))
}

func TestPrice_TaxFromNet(t *testing.T) {
	// 100 €
	price := domain.NewFromInt(100, 1, "EUR")
	tax := price.TaxFromNet(*new(big.Float).SetInt64(19))
	assert.Equal(t, tax, domain.NewFromInt(19, 1, "EUR"), "expect 19 € tax fromm 100€")

	taxedPrice := price.Taxed(*new(big.Float).SetInt64(19))
	assert.Equal(t, taxedPrice, domain.NewFromInt(119, 1, "EUR"))
}

func TestPrice_LikelyEqual(t *testing.T) {
	price1 := domain.NewFromFloat(100, "EUR")
	price2 := domain.NewFromFloat(100.000000000000001, "EUR")
	price3 := domain.NewFromFloat(100.1, "EUR")
	assert.True(t, price1.LikelyEqual(price2))
	assert.False(t, price1.LikelyEqual(price3))
}

func TestPrice_MarshalBinaryForGob(t *testing.T) {
	type (
		SomeTypeWithPrice struct {
			Price domain.Price
		}
	)
	gob.Register(SomeTypeWithPrice{})
	var network bytes.Buffer
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.

	err := enc.Encode(&SomeTypeWithPrice{Price: domain.NewFromInt(1111, 100, "EUR")})
	if err != nil {
		t.Fatal("encode error:", err)
	}
	var receivedPrice SomeTypeWithPrice
	err = dec.Decode(&receivedPrice)
	if err != nil {
		t.Fatal("decode error 1:", err)
	}
	float, _ := receivedPrice.Price.Amount().Float64()
	assert.Equal(t, 11.11, float)
}

func TestPrice_GetPayableByRoundingMode(t *testing.T) {
	price := domain.NewFromFloat(12.34567, "EUR")

	payable := price.GetPayableByRoundingMode(domain.RoundingModeCeil, 100)
	assert.Equal(t, domain.NewFromInt(1235, 100, "EUR").Amount(), payable.Amount())

	payable = price.GetPayableByRoundingMode(domain.RoundingModeFloor, 100)
	assert.Equal(t, domain.NewFromInt(1234, 100, "EUR").Amount(), payable.Amount())

	payable = price.GetPayableByRoundingMode(domain.RoundingModeFloor, 1)
	assert.Equal(t, domain.NewFromInt(12, 1, "EUR").Amount(), payable.Amount())

	price = domain.NewFromFloat(-0.119, "EUR")
	payable = price.GetPayableByRoundingMode(domain.RoundingModeFloor, 100)
	assert.Equal(t, domain.NewFromInt(-12, 100, "EUR").Amount(), payable.Amount())
}

func TestPrice_GetPayableByRoundingMode_RoundingModeCeil(t *testing.T) {
	tests := []struct {
		price     float64
		precision int
		expected  int64
	}{
		{price: 12.34567, precision: 100, expected: 1235},
		{price: -12.34567, precision: 100, expected: -1234},
		{price: 5.5, precision: 1, expected: 6},
		{price: 2.5, precision: 1, expected: 3},
		{price: 1.6, precision: 1, expected: 2},
		{price: 1.1, precision: 1, expected: 2},
		{price: 1.0, precision: 1, expected: 1},
		{price: -1.0, precision: 1, expected: -1},
		{price: -1.1, precision: 1, expected: -1},
		{price: -1.6, precision: 1, expected: -1},
		{price: -2.5, precision: 1, expected: -2},
		{price: -5.5, precision: 1, expected: -5},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("rounding %f", tt.price), func(t *testing.T) {
			price := domain.NewFromFloat(tt.price, "EUR")

			payable := price.GetPayableByRoundingMode(domain.RoundingModeCeil, tt.precision)
			assert.Equal(t, domain.NewFromInt(tt.expected, tt.precision, "EUR").Amount(), payable.Amount())
		})
	}
}

func TestPrice_GetPayableByRoundingMode_RoundingModeFloor(t *testing.T) {
	tests := []struct {
		price     float64
		precision int
		expected  int64
	}{
		{price: 12.34567, precision: 100, expected: 1234},
		{price: -12.34567, precision: 100, expected: -1235},
		{price: 5.5, precision: 1, expected: 5},
		{price: 2.5, precision: 1, expected: 2},
		{price: 1.6, precision: 1, expected: 1},
		{price: 1.1, precision: 1, expected: 1},
		{price: 1.0, precision: 1, expected: 1},
		{price: -1.0, precision: 1, expected: -1},
		{price: -1.1, precision: 1, expected: -2},
		{price: -1.6, precision: 1, expected: -2},
		{price: -2.5, precision: 1, expected: -3},
		{price: -5.5, precision: 1, expected: -6},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("rounding %f", tt.price), func(t *testing.T) {
			price := domain.NewFromFloat(tt.price, "EUR")

			payable := price.GetPayableByRoundingMode(domain.RoundingModeFloor, tt.precision)
			assert.Equal(t, domain.NewFromInt(tt.expected, tt.precision, "EUR").Amount(), payable.Amount())
		})
	}
}

func TestPrice_GetPayableByRoundingMode_RoundingModeHalfUp(t *testing.T) {
	tests := []struct {
		price     float64
		precision int
		expected  int64
		msg       string
	}{
		{price: 7.6, precision: 1, expected: 8, msg: "7.6 should be rounded to 8"},
		{price: 7.5, precision: 1, expected: 8, msg: "7.5 should be rounded to 8"},
		{price: 7.4, precision: 1, expected: 7, msg: "7.4 should be rounded to 7"},
		{price: 12.34567, precision: 100, expected: 1235, msg: "12.34567 should be rounded to 12.35"},
		{price: -7.4, precision: 1, expected: -7, msg: "-7.4 should be rounded to -7"},
		{price: -7.45, precision: 1, expected: -7, msg: "-7.45 should be rounded to -7"},
		{price: -7.5, precision: 1, expected: -8, msg: "-7.5 should be rounded to -8"},
		{price: -7.55, precision: 1, expected: -8, msg: "-7.55 should be rounded to -8"},
		{price: -7.6, precision: 1, expected: -8, msg: "-7.6 should be rounded to -8"},

		{price: 5.5, precision: 1, expected: 6, msg: "5.5 should be rounded to 6"},
		{price: 2.5, precision: 1, expected: 3, msg: "2.5 should be rounded to 3"},
		{price: 1.6, precision: 1, expected: 2, msg: "1.6 should be rounded to 2"},
		{price: 1.1, precision: 1, expected: 1, msg: "1.1 should be rounded to 1"},
		{price: 1.0, precision: 1, expected: 1, msg: "1.0 should be rounded to 1"},
		{price: -1.0, precision: 1, expected: -1, msg: "-1.0 should be rounded to -1"},
		{price: -1.1, precision: 1, expected: -1, msg: "-1.1 should be rounded to -1"},
		{price: -1.6, precision: 1, expected: -2, msg: "-1.6 should be rounded to -2"},
		{price: -2.5, precision: 1, expected: -3, msg: "-2.5 should be rounded to -3"},
		{price: -5.5, precision: 1, expected: -6, msg: "-5.5 should be rounded to -6"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("rounding %f", tt.price), func(t *testing.T) {
			price := domain.NewFromFloat(tt.price, "EUR")

			payable := price.GetPayableByRoundingMode(domain.RoundingModeHalfUp, tt.precision)
			assert.Equal(t, domain.NewFromInt(tt.expected, tt.precision, "EUR").Amount(), payable.Amount(), tt.msg)
		})
	}
}

func TestPrice_GetPayableByRoundingMode_RoundingModeHalfDown(t *testing.T) {
	tests := []struct {
		price     float64
		precision int
		expected  int64
		msg       string
	}{
		{price: 7.6, precision: 1, expected: 8, msg: "7.6 should be rounded to 8"},
		{price: 7.5, precision: 1, expected: 7, msg: "7.5 should be rounded to 7"},
		{price: 7.4, precision: 1, expected: 7, msg: "7.4 should be rounded to 7"},
		{price: 12.34567, precision: 100, expected: 1235, msg: "12.34567 should be rounded to 12.35"},

		{price: -7.4, precision: 1, expected: -7, msg: "-7.4 should be rounded to -7"},
		{price: -7.5, precision: 1, expected: -7, msg: "-7.5 should be rounded to -7"},
		{price: -7.6, precision: 1, expected: -8, msg: "-7.6 should be rounded to -8"},

		{price: 5.5, precision: 1, expected: 5, msg: "5.5 should be rounded to 5"},
		{price: 2.5, precision: 1, expected: 2, msg: "2.5 should be rounded to 2"},
		{price: 1.6, precision: 1, expected: 2, msg: "1.6 should be rounded to 2"},
		{price: 1.1, precision: 1, expected: 1, msg: "1.1 should be rounded to 1"},
		{price: 1.0, precision: 1, expected: 1, msg: "1.0 should be rounded to 1"},
		{price: -1.0, precision: 1, expected: -1, msg: "-1.0 should be rounded to -1"},
		{price: -1.1, precision: 1, expected: -1, msg: "-1.1 should be rounded to -1"},
		{price: -1.6, precision: 1, expected: -2, msg: "-1.6 should be rounded to -2"},
		{price: -2.5, precision: 1, expected: -2, msg: "-2.5 should be rounded to -2"},
		{price: -5.5, precision: 1, expected: -5, msg: "-5.5 should be rounded to -5"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("rounding %f", tt.price), func(t *testing.T) {
			price := domain.NewFromFloat(tt.price, "EUR")

			payable := price.GetPayableByRoundingMode(domain.RoundingModeHalfDown, tt.precision)
			assert.Equal(t, domain.NewFromInt(tt.expected, tt.precision, "EUR").Amount(), payable.Amount(), tt.msg)
		})
	}

}

func TestCharges_Add(t *testing.T) {
	c1 := domain.Charges{}

	byType := make(map[string]domain.Charge)
	byType["main"] = domain.Charge{
		Type:  "main",
		Price: domain.NewFromInt(100, 1, "EUR"),
		Value: domain.NewFromInt(50, 1, "EUR"),
	}
	c2 := domain.NewCharges(byType)

	byType = make(map[string]domain.Charge)
	byType["main"] = domain.Charge{
		Type:  "main",
		Price: domain.NewFromInt(100, 1, "EUR"),
		Value: domain.NewFromInt(100, 1, "EUR"),
	}
	c3 := domain.NewCharges(byType)

	c1and2 := c1.Add(*c2)
	charge, found := c1and2.GetByType("main")
	assert.True(t, found)
	assert.Equal(t, domain.Charge{
		Price: domain.NewFromInt(100, 1, "EUR"),
		Value: domain.NewFromInt(50, 1, "EUR"),
		Type:  "main",
	}, charge)

	c2and3 := c2.Add(*c3)
	charge, found = c2and3.GetByType("main")
	assert.True(t, found)
	assert.Equal(t, domain.Charge{
		Price: domain.NewFromInt(200, 1, "EUR"),
		Value: domain.NewFromInt(150, 1, "EUR"),
		Type:  "main",
	}, charge)
}

func TestCharges_GetAllByType(t *testing.T) {
	charges := domain.Charges{}
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: "type-a", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: "type-x", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: "type-c", Reference: "HUHUWHHUHX", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: "type-a", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Reference: "ABC123", Price: domain.NewFromInt(200, 1, "€")})

	assert.Len(t, charges.GetAllByType(domain.ChargeTypeMain), 3)
	assert.Len(t, charges.GetAllByType("type-a"), 1)
	assert.Len(t, charges.GetAllByType("type-c"), 1)
	assert.Len(t, charges.GetAllByType("type-x"), 1)
}

func TestCharges_GetByType(t *testing.T) {
	charges := domain.Charges{}
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: "type-a", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: "type-x", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: "type-c", Reference: "HUHUWHHUHX", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: "type-a", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Reference: "ABC123", Price: domain.NewFromInt(200, 1, "€")})

	charge, found := charges.GetByType(domain.ChargeTypeMain)
	assert.True(t, found)
	assert.Equal(t, charge, domain.Charge{Type: domain.ChargeTypeMain, Price: domain.NewFromInt(600, 1, "€")})

	charge, found = charges.GetByType("type-a")
	assert.True(t, found)
	want := domain.Charge{Type: "type-a", Price: domain.NewFromInt(400, 1, "€")}
	assert.Equal(t, charge.Price, want.Price)
}

func TestCharges_GetByTypeForced(t *testing.T) {
	charges := domain.Charges{}

	charge := charges.GetByTypeForced(domain.ChargeTypeMain)
	assert.Equal(t, charge, domain.Charge{})

	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82", Price: domain.NewFromInt(200, 1, "€")})
	charge = charges.GetByTypeForced(domain.ChargeTypeMain)

	assert.Equal(t, charge, domain.Charge{Type: domain.ChargeTypeMain, Price: domain.NewFromInt(200, 1, "€")})
}

func TestCharges_GetByChargeQualifier(t *testing.T) {
	charges := domain.Charges{}
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82", Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Price: domain.NewFromInt(200, 1, "€")})
	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Reference: "ABC123", Price: domain.NewFromInt(200, 1, "€")})

	charge, found := charges.GetByChargeQualifier(domain.ChargeQualifier{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82"})
	assert.True(t, found)
	assert.Equal(t, charge, domain.Charge{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82", Price: domain.NewFromInt(200, 1, "€")})

}

func TestCharges_GetByChargeQualifierForced(t *testing.T) {
	charges := domain.Charges{}
	charge := charges.GetByChargeQualifierForced(domain.ChargeQualifier{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82"})
	assert.Equal(t, charge, domain.Charge{})

	charges = charges.AddCharge(domain.Charge{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82", Price: domain.NewFromInt(200, 1, "€")})
	charge = charges.GetByChargeQualifierForced(domain.ChargeQualifier{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82"})
	assert.Equal(t, charge, domain.Charge{Type: domain.ChargeTypeMain, Reference: "SJHHQWAXX6HJSDZ82", Price: domain.NewFromInt(200, 1, "€")})

}

func TestPrice_MarshalJSON(t *testing.T) {
	t.Run("JSON marshalling works", func(t *testing.T) {
		price := domain.NewFromFloat(55.111111, "USD")

		priceJSON, err := price.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"Amount":"55.11","Currency":"USD"}`, string(priceJSON))
	})

	t.Run("JSON price is rounded", func(t *testing.T) {
		price := domain.NewFromFloat(55.119, "USD")

		priceJSON, err := price.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"Amount":"55.12","Currency":"USD"}`, string(priceJSON))
	})
}

func TestPrice_UnmarshalJSON(t *testing.T) {
	var p domain.Price

	err := json.Unmarshal([]byte(`{"Amount":"55.12","Currency":"USD"}`), &p)
	require.NoError(t, err)
	assert.Equal(t, domain.NewFromFloat(55.12, "USD").GetPayable(), p.GetPayable())
}
