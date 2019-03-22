package domain

import (
	"encoding"
	"encoding/json"
	"errors"
	"math/big"
)

type (
	//Price is a Type that represents a Amount - it is immutable
	// DevHint: We use Amount and Charge as Value - so we do not pass pointers. (According to Go Wiki's code review comments page suggests passing by value when structs are small and likely to stay that way)
	Price struct {
		amount   big.Float
		currency string
	}

	//Charge is a Amount of a certain Type. Charge is used to indicate that this need to be payed
	Charge struct {
		Price Price
		Type  string
	}

	//priceEncodeAble is a type that we need to allow marshalling the price values. The type itself is unexported
	priceEncodeAble struct {
		Amount   big.Float
		Currency string
	}
)

var (
	_ encoding.BinaryMarshaler   = Price{}
	_ encoding.BinaryUnmarshaler = &Price{}
)

const (
	//ChargeTypeMain used as default for a Charge
	ChargeTypeMain = "main"
)

//NewFromFloat - factory method
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
		amount:   *new(big.Float).Quo(amountF, precicionF),
		currency: currency,
	}
}

//Add - Adds the given price to the current price and returns a new price
func (p Price) Add(add Price) (Price, error) {
	newPrice, err := p.currencyGuard(add)
	if err != nil {
		return newPrice, err
	}
	newPrice.amount.Add(&p.amount, &add.amount)
	return newPrice, nil
}

//ForceAdd - tries to Adds the given price to the current price - will not return errors
func (p Price) ForceAdd(add Price) Price {
	newPrice, err := p.currencyGuard(add)
	if err != nil {
		return p
	}
	newPrice.amount.Add(&p.amount, &add.amount)
	return newPrice
}

//currencyGuard - common Guard that protects price calculations of prices with different currency.
// 	Robust: if original is Zero and the currencies are different we take the given currency
func (p Price) currencyGuard(check Price) (Price, error) {
	if p.currency == check.currency {
		return Price{
			currency: check.currency,
		}, nil
	}
	if p.IsZero() {
		return Price{
			currency: check.currency,
		}, nil
	}

	if check.IsZero() {
		return Price{
			currency: p.currency,
		}, nil
	}
	return NewZero(p.currency), errors.New("Cannot calculate prices in different currencies")
}

//Discounted - returns new price reduced by given percent
func (p Price) Discounted(percent float64) Price {
	newPrice := Price{
		currency: p.currency,
		amount:   *new(big.Float).Mul(&p.amount, big.NewFloat((100-percent)/100)),
	}
	return newPrice
}

//Taxed - returns new price added with Tax (assuming current price is net)
func (p Price) Taxed(percent big.Float) Price {
	newPrice := Price{
		currency: p.currency,
		amount:   *new(big.Float).Add(&p.amount, p.TaxFromNet(percent).Amount()),
	}
	return newPrice
}

//TaxFromNet - returns new price representing the taxamount (assuming the current price is net 100%)
func (p Price) TaxFromNet(percent big.Float) Price {
	quo := new(big.Float).Mul(&percent, &p.amount)
	newPrice := Price{
		currency: p.currency,
		amount:   *new(big.Float).Quo(quo, new(big.Float).SetInt64(100)),
	}
	return newPrice
}

//TaxFromGross - returns new price representing the taxamount (assuming the current price is gross 100+percent)
func (p Price) TaxFromGross(percent big.Float) Price {
	quo := new(big.Float).Mul(&percent, &p.amount)
	percent100 := new(big.Float).Add(&percent, new(big.Float).SetInt64(100))
	newPrice := Price{
		currency: p.currency,
		amount:   *new(big.Float).Quo(quo, percent100),
	}
	return newPrice
}

//Sub - Subtract the given price from the current price and returns a new price
func (p Price) Sub(sub Price) (Price, error) {
	newPrice, err := p.currencyGuard(sub)
	if err != nil {
		return newPrice, err
	}
	newPrice.amount.Sub(&p.amount, &sub.amount)
	return newPrice, nil
}

//Multiply  returns a new price with the amount Multiply
func (p Price) Multiply(qty int) Price {
	newPrice := Price{
		currency: p.currency,
	}
	newPrice.amount.Mul(&p.amount, new(big.Float).SetInt64(int64(qty)))
	return newPrice
}

//Divided  returns a new price with the amount Divided
func (p Price) Divided(qty int) Price {
	newPrice := Price{
		currency: p.currency,
	}
	if qty == 0 {
		//TODO log
		return NewZero(p.currency)
	}
	newPrice.amount.Quo(&p.amount, new(big.Float).SetInt64(int64(qty)))
	return newPrice
}

//Equal - compares the prices exact
func (p Price) Equal(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	return p.amount.Cmp(&cmp.amount) == 0
}

//LikelyEqual - compares the prices with some tolerance
func (p Price) LikelyEqual(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	diff := new(big.Float).Sub(&p.amount, &cmp.amount)
	absDiff := new(big.Float).Abs(diff)
	return absDiff.Cmp(big.NewFloat(0.000000001)) == -1
}

//IsLessThen - compares the prices
func (p Price) IsLessThen(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	return p.amount.Cmp(&cmp.amount) == -1
}

//IsGreaterThen - compares the prices
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

//IsPositive - returns true if the price represents a positive value
func (p Price) IsPositive() bool {
	return p.IsGreaterThenValue(*big.NewFloat(0.0))
}

//IsPayable - returns true if the price represents a payable (rounded) value
func (p Price) IsPayable() bool {
	return p.GetPayable().Equal(p)
}

//IsZero - returns true if the price represents zero value
func (p Price) IsZero() bool {
	return p.Equal(NewZero(p.Currency())) || p.Equal(NewFromFloat(0, p.Currency()))
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

	offsetToCheckRounding := new(big.Float).Mul(p.payableRoundingPrecisionF(), new(big.Float).SetInt64(10))

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

	splittedAmounts := make([]int64, count)
	for i := 0; i < count; i++ {
		splittedAmounts[i] = splittedAmount
	}

	for i := 0; i < int(splittedAmountModulo); i++ {
		splittedAmounts[i] = splittedAmounts[i] + 1
	}

	prices := make([]Price, count)
	for i := 0; i < count; i++ {
		prices[i] = NewFromInt(splittedAmounts[i], p.payableRoundingPrecision(), p.Currency())
	}

	return prices, nil
}

//Clone returns a copy of the price - the amount gets Excat acc
func (p Price) Clone() Price {
	return Price{
		amount:   *new(big.Float).Set(&p.amount),
		currency: p.currency,
	}
}

//Currency returns currency
func (p Price) Currency() string {
	return p.currency
}

//Amount - returns exact amount as bigFloat
func (p Price) Amount() *big.Float {
	return &p.amount
}

//SumAll - retruns new price with sum of all given prices
func SumAll(prices ...Price) (Price, error) {
	if len(prices) == 0 {
		return NewZero(""), errors.New("no price given")
	}
	result := prices[0].Clone()
	var err error
	for _, price := range prices[1:] {
		result, err = result.Add(price)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

//MarshalJSON - implements interace required by json marshal
func (p Price) MarshalJSON() (data []byte, err error) {
	pn := priceEncodeAble{
		Amount:   p.amount,
		Currency: p.currency,
	}
	r, e := json.Marshal(&pn)
	return r, e
}

//MarshalBinary - implements interface required by gob
func (p Price) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}

//UnmarshalBinary - implements interace required by gob.
//UnmarshalBinary - modifies the receiver so it must take a pointer receiver!
func (p *Price) UnmarshalBinary(data []byte) error {
	var pe priceEncodeAble
	err := json.Unmarshal(data, &pe)
	if err != nil {
		return err
	}
	p.amount = pe.Amount
	p.currency = pe.Currency
	return nil
}
