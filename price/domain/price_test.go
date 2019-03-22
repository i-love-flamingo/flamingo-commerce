package domain_test

import (
	"bytes"
	"encoding/gob"
	"math/big"
	"testing"

	"flamingo.me/flamingo-commerce/v3/price/domain"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, domain.NewFromFloat(12.35, "EUR").GetPayable().Amount(), payable.Amount())

}

func TestNewFromInt(t *testing.T) {
	price1 := domain.NewFromInt(1245, 100, "EUR")
	price2 := domain.NewFromFloat(12.45, "EUR")
	assert.Equal(t, price2.GetPayable().Amount(), price1.GetPayable().Amount())
	pricePayable := price1.GetPayable()
	assert.True(t, price2.GetPayable().Equal(pricePayable))
}

func TestPrice_SplitInPayables(t *testing.T) {
	originalPrice := domain.NewFromFloat(12.456, "EUR")
	payableSplitPrices, _ := originalPrice.SplitInPayables(6)

	sumPrice := domain.NewZero("EUR")
	for _, price := range payableSplitPrices {
		sumPrice, _ = sumPrice.Add(price)
	}
	//sum of the splitted payable need to match original price payable
	assert.Equal(t, originalPrice.GetPayable().Amount(), sumPrice.GetPayable().Amount())
}

func TestPrice_Discounted(t *testing.T) {
	originalPrice := domain.NewFromFloat(12.45, "EUR")
	discountedPrice := originalPrice.Discounted(10).GetPayable()
	//10% of - expected rounded value of 11.21
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
	//119 €
	price := domain.NewFromInt(119, 1, "EUR")
	tax := price.TaxFromGross(*new(big.Float).SetInt64(19))
	assert.Equal(t, tax, domain.NewFromInt(19, 1, "EUR"))
}

func TestPrice_TaxFromNet(t *testing.T) {
	//100 €
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
