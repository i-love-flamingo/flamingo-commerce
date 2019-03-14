package domain

import (
	"errors"
	"math/big"
)

type (
	//Price is a Type that represents a Price - it is immutable
	// DevHint: We use Price and Charge as Value - so we do not pass pointers. (According to Go Wiki's code review comments page suggests passing by value when structs are small and likely to stay that way)
	Price struct {
		amount   big.Float
		currency string
	}

	//Charge is a Price of a certain Type. Charge is used to indicate that this need to be payed
	Charge struct {
		Price Price
		Type  string
	}
)

const (
	ChargeTypeMain = "main"
)

func NewFromFloat(amount float64, currency string) Price {
	return Price{
		amount:   *big.NewFloat(amount),
		currency: currency,
	}
}


// NewZero Zero price
func NewZero(currency string) Price {
	return Price{
		amount:   *new(big.Float).SetInt64(0),
		currency: currency,
	}
}

// NewFromInt use to set money by smallest payable unit - e.g. to set 2.45 EUR you should use NewFromInt(245,100)
func NewFromInt(amount int64, precicion int, currency string) Price {
	amountF := new(big.Float).SetInt64(amount)
	if precicion == 0 {
		return Price{
			amount:   *new(big.Float).SetInt64(0),
			currency: currency,
		}
	}
	precicionF := new(big.Float).SetInt64(int64(precicion))
	return Price{
		amount:   *new(big.Float).Quo(amountF,precicionF),
		currency: currency,
	}
}

//Add - Adds the given price to the current price and returns a new price
func (p Price) Add(add Price) (Price, error) {
	if p.currency != add.currency {
		return p, errors.New("Cannot add prices in different currencies")
	}
	newPrice := Price{
		currency: p.currency,
	}
	newPrice.amount.Add(&p.amount, &add.amount)
	return newPrice, nil
}


//Discounted - returns new price reduced by given percent
func (p Price) Discounted(percent float64) (Price) {
	newPrice := Price{
		currency: p.currency,
		amount: *new(big.Float).Mul(&p.amount,big.NewFloat((100-percent)/100)),
	}
	return newPrice
}


//Taxed - returns new price added with Tax
func (p Price) Taxed(percent float64) (Price) {
	newPrice := Price{
		currency: p.currency,
		amount: *new(big.Float).Mul(&p.amount,big.NewFloat((100+percent)/100)),
	}
	return newPrice
}

//Sub - Subtract the given price from the current price and returns a new price
func (p Price) Sub(sub Price) (Price, error) {
	if p.currency != sub.currency {
		return p, errors.New("Cannot add prices in different currencies")
	}
	newPrice := Price{
		currency: p.currency,
	}
	newPrice.amount.Sub(&p.amount, &sub.amount)
	return newPrice, nil
}

//Multiply  returns a new price with the amount Multiply
func (p Price) Multiply(qty int) (Price) {
	newPrice := Price{
		currency: p.currency,
	}
	newPrice.amount.Mul(&p.amount, new(big.Float).SetInt64(int64(qty)))
	return newPrice
}

//Sub - Subtract the given price from the current price and returns a new price
func (p Price) Equals(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	return p.amount.Cmp(&cmp.amount) == 0
}

func (p Price) IsLessThen(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	return p.amount.Cmp(&cmp.amount) == -1
}

func (p Price) IsGreaterThen(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	return p.amount.Cmp(&cmp.amount) == 1
}

//IsLessThenValue compares the price with a given amount value (assuming same currency)
func (p Price) IsLessThenValue(amount big.Float) bool {
	if p.amount.Cmp(&amount) == -1 {
		return true
	}
	return false
}

//IsGreaterThenValue compares the price with a given amount value (assuming same currency)
func (p Price) IsGreaterThenValue(amount big.Float) bool {
	if p.amount.Cmp(&amount) == 1 {
		return true
	}
	return false
}

//IsNegative - returns true if the price represents a negative value
func (p Price) IsNegative() bool {
	return p.IsLessThenValue(*big.NewFloat(0.0))
}

//IsNegative - returns true if the price represents a negative value
func (p Price) IsPositive() bool {
	return p.IsGreaterThenValue(*big.NewFloat(0.0))
}

//FloatAmount gets the current amount as float
func (p Price) FloatAmount() float64 {
	a, _ := p.amount.Float64()
	return a
}

// GetPayable - rounds the price with the precision required by the currency in a price that can actually be payed
// e.g. an internal amount of 1,23344 will get rounded to 1,23
func (p Price) GetPayable() Price {
	newPrice := Price{
		currency: p.currency,
	}

	amountForRound := new(big.Float).Copy(&p.amount)

	offsetToCheckRounding := new(big.Float).Mul(p.payableRoundingPrecisionF(),new(big.Float).SetInt64(10))

	amountTruncatedInt, _ := new(big.Float).Mul(amountForRound, p.payableRoundingPrecisionF()).Int64()
	amountRoundingCheckInt, _ := new(big.Float).Mul(amountForRound, offsetToCheckRounding).Int64()
	if (amountRoundingCheckInt - (amountTruncatedInt * 10)) >= 5 {
		amountTruncatedInt = amountTruncatedInt + 1
	}

	amountRounded := new(big.Float).Quo(new(big.Float).SetInt64(amountTruncatedInt), p.payableRoundingPrecisionF())
	newPrice.amount = *amountRounded
	return newPrice
}

//payableRoundingPrecisionF - 10 * n - n is the amount of decimal numbers after comma
// - can be currency specific (for now defaults to 2)
func (p Price) payableRoundingPrecisionF() *big.Float {
	return new(big.Float).SetInt64(int64(p.payableRoundingPrecision()))
}

//payableRoundingPrecisionF - 10 * n - n is the amount of decimal numbers after comma
// - can be currency specific (for now defaults to 2)
func (p Price) payableRoundingPrecision() int {
	return int(100)
}

// SplitInPayables - returns "count" payable prices (each rounded) that in sum matches the given price
//  - Given a price of 12.456 (Payable 12,46)  - Splitted in 6 will mean: 6 * 2.076
//  - but having them payable requires rounding them each (e.g. 2.07) which would mean we have 0.03 difference (=12,45-6*2.07)
//  - so that the sum is as close as possible to the original value   in this case the correct return will be:
//  - 	 2.07 + 2.07+2.08 +2.08 +2.08 +2.08
func (p Price) SplitInPayables(count int) ([]Price, error) {
	if count <= 0 {
		return nil, errors.New("Split must be higher than zero")
	}

	amountToMatchInt, _ := new(big.Float).Mul(p.GetPayable().Amount(), p.payableRoundingPrecisionF()).Int64()

	splittedAmountModulo := amountToMatchInt % int64(count)
	splittedAmount := amountToMatchInt / int64(count)

	splittedAmounts := make([]int64,count)
	for i:=0; i < count; i++ {
		splittedAmounts[i] = splittedAmount
	}

	for i:=0; i < int(splittedAmountModulo); i++ {
		splittedAmounts[i] = splittedAmounts[i] +1
	}

	prices := make([]Price,count)
	for i:=0; i < count; i++ {
		prices[i] = NewFromInt(splittedAmounts[i],p.payableRoundingPrecision(),p.Currency())
	}


	return prices, nil
}


func (p Price) Clone() Price {
	return Price{
		amount:p.amount,
		currency:p.currency,
	}
}

func (p Price) Currency() string {
	return p.currency
}

func (p Price) Amount() *big.Float {
	return &p.amount
}
