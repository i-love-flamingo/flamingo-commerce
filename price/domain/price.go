package domain

import (
	"encoding"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type (
	// Price is a Type that represents a Amount - it is immutable
	// DevHint: We use Amount and Charge as Value - so we do not pass pointers. (According to Go Wiki's code review comments page suggests passing by value when structs are small and likely to stay that way)
	Price struct {
		amount   big.Float `swaggertype:"string"`
		currency string
	}

	// Charge is a Amount of a certain Type. Charge is used as value object
	Charge struct {
		// Price that is paid, can be in a certain currency
		Price Price
		// Value of the "Price" in another (base) currency
		Value Price
		// Type of the charge - can be ChargeTypeMain or something else. Used to differentiate between different charges of a single thing
		Type string
		// Reference contains further information to distinguish charges of the same type
		Reference string
	}

	// Charges - Represents the Charges the product need to be paid with
	Charges struct {
		chargesByQualifier map[ChargeQualifier]Charge
	}

	// ChargeQualifier distinguishes charges by type and reference
	ChargeQualifier struct {
		// Type represents charge type
		Type string
		// Reference contains further information to distinguish charges of the same type
		Reference string
	}

	// priceEncodeAble is a type that we need to allow marshalling the price values. The type itself is unexported
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
	// ChargeTypeGiftCard  used as a charge type for gift cards
	ChargeTypeGiftCard = "giftcard"
	// ChargeTypeMain used as default for a Charge
	ChargeTypeMain = "main"

	// RoundingModeFloor use if you want to cut (round down)
	RoundingModeFloor = "floor"
	// RoundingModeCeil use if you want to round up always
	RoundingModeCeil = "ceil"
	// RoundingModeHalfUp round up if the discarded fraction is ≥ 0.5, otherwise round down. Default for GetPayable()
	RoundingModeHalfUp = "halfup"
	// RoundingModeHalfDown round up if the discarded fraction is > 0.5, otherwise round down.
	RoundingModeHalfDown = "halfdown"
)

// NewFromFloat - factory method
func NewFromFloat(amount float64, currency string) Price {
	return Price{
		amount:   *big.NewFloat(amount),
		currency: currency,
	}
}

// NewFromBigFloat - factory method
func NewFromBigFloat(amount big.Float, currency string) Price {
	return Price{
		amount:   amount,
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

// NewFromInt use to set money by smallest payable unit - e.g. to set 2.45 EUR you should use NewFromInt(245, 100, "EUR")
func NewFromInt(amount int64, precision int, currency string) Price {
	amountF := new(big.Float).SetInt64(amount)
	if precision == 0 {
		return Price{
			amount:   *new(big.Float).SetInt64(0),
			currency: currency,
		}
	}
	precicionF := new(big.Float).SetInt64(int64(precision))
	return Price{
		amount:   *new(big.Float).Quo(amountF, precicionF),
		currency: currency,
	}
}

// Add the given price to the current price and returns a new price
func (p Price) Add(add Price) (Price, error) {
	newPrice, err := p.currencyGuard(add)
	if err != nil {
		return newPrice, err
	}
	newPrice.amount.Add(&p.amount, &add.amount)
	return newPrice, nil
}

// ForceAdd tries to add the given price to the current price - will not return errors
func (p Price) ForceAdd(add Price) Price {
	newPrice, err := p.currencyGuard(add)
	if err != nil {
		return p
	}
	newPrice.amount.Add(&p.amount, &add.amount)
	return newPrice
}

// currencyGuard is a common Guard that protects price calculations of prices with different currency.
// Robust: if original is Zero and the currencies are different we take the given currency
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
	return NewZero(p.currency), errors.New("cannot calculate prices in different currencies")
}

// Discounted returns new price reduced by given percent
func (p Price) Discounted(percent float64) Price {
	newPrice := Price{
		currency: p.currency,
		amount:   *new(big.Float).Mul(&p.amount, big.NewFloat((100-percent)/100)),
	}
	return newPrice
}

// Taxed returns new price added with Tax (assuming current price is net)
func (p Price) Taxed(percent big.Float) Price {
	newPrice := Price{
		currency: p.currency,
		amount:   *new(big.Float).Add(&p.amount, p.TaxFromNet(percent).Amount()),
	}
	return newPrice
}

// TaxFromNet returns new price representing the tax amount (assuming the current price is net 100%)
func (p Price) TaxFromNet(percent big.Float) Price {
	quo := new(big.Float).Mul(&percent, &p.amount)
	newPrice := Price{
		currency: p.currency,
		amount:   *new(big.Float).Quo(quo, new(big.Float).SetInt64(100)),
	}
	return newPrice
}

// TaxFromGross returns new price representing the tax amount (assuming the current price is gross 100+percent)
func (p Price) TaxFromGross(percent big.Float) Price {
	quo := new(big.Float).Mul(&percent, &p.amount)
	percent100 := new(big.Float).Add(&percent, new(big.Float).SetInt64(100))
	newPrice := Price{
		currency: p.currency,
		amount:   *new(big.Float).Quo(quo, percent100),
	}
	return newPrice
}

// Sub the given price from the current price and returns a new price
func (p Price) Sub(sub Price) (Price, error) {
	newPrice, err := p.currencyGuard(sub)
	if err != nil {
		return newPrice, err
	}
	newPrice.amount.Sub(&p.amount, &sub.amount)
	return newPrice, nil
}

// Inverse returns the price multiplied with -1
func (p Price) Inverse() Price {
	p.amount = *new(big.Float).Mul(&p.amount, big.NewFloat(-1))
	return p
}

// Multiply returns a new price with the amount Multiply
func (p Price) Multiply(qty int) Price {
	newPrice := Price{
		currency: p.currency,
	}
	newPrice.amount.Mul(&p.amount, new(big.Float).SetInt64(int64(qty)))
	return newPrice
}

// Divided returns a new price with the amount Divided
func (p Price) Divided(qty int) Price {
	newPrice := Price{
		currency: p.currency,
	}
	if qty == 0 {
		return NewZero(p.currency)
	}
	newPrice.amount.Quo(&p.amount, new(big.Float).SetInt64(int64(qty)))
	return newPrice
}

// Equal compares the prices exact
func (p Price) Equal(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	return p.amount.Cmp(&cmp.amount) == 0
}

// LikelyEqual compares the prices with some tolerance
func (p Price) LikelyEqual(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	diff := new(big.Float).Sub(&p.amount, &cmp.amount)
	absDiff := new(big.Float).Abs(diff)
	return absDiff.Cmp(big.NewFloat(0.000000001)) == -1
}

// IsLessThen compares the current price with a given one
func (p Price) IsLessThen(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	return p.amount.Cmp(&cmp.amount) == -1
}

// IsGreaterThen compares the current price with a given one
func (p Price) IsGreaterThen(cmp Price) bool {
	if p.currency != cmp.currency {
		return false
	}
	return p.amount.Cmp(&cmp.amount) == 1
}

// IsLessThenValue compares the price with a given amount value (assuming same currency)
func (p Price) IsLessThenValue(amount big.Float) bool {
	return p.amount.Cmp(&amount) == -1
}

// IsGreaterThenValue compares the price with a given amount value (assuming same currency)
func (p Price) IsGreaterThenValue(amount big.Float) bool {
	return p.amount.Cmp(&amount) == 1
}

// IsNegative returns true if the price represents a negative value
func (p Price) IsNegative() bool {
	return p.IsLessThenValue(*big.NewFloat(0.0))
}

// IsPositive returns true if the price represents a positive value
func (p Price) IsPositive() bool {
	return p.IsGreaterThenValue(*big.NewFloat(0.0))
}

// IsPayable returns true if the price represents a payable (rounded) value
func (p Price) IsPayable() bool {
	return p.GetPayable().Equal(p)
}

// IsZero returns true if the price represents zero value
func (p Price) IsZero() bool {
	return p.LikelyEqual(NewZero(p.Currency())) || p.LikelyEqual(NewFromFloat(0, p.Currency()))
}

// FloatAmount gets the current amount as float
func (p Price) FloatAmount() float64 {
	a, _ := p.amount.Float64()
	return a
}

// GetPayable rounds the price with the precision required by the currency in a price that can actually be paid
// e.g. an internal amount of 1,23344 will get rounded to 1,23
func (p Price) GetPayable() Price {
	mode, precision := p.payableRoundingPrecision()
	return p.GetPayableByRoundingMode(mode, precision)
}

// GetPayableByRoundingMode returns the price rounded you can pass the used rounding mode and precision
// Example for precision 100:
//
//	1.115 >  1.12 (RoundingModeHalfUp)  / 1.11 (RoundingModeFloor)
//	-1.115 > -1.11 (RoundingModeHalfUp) / -1.12 (RoundingModeFloor)
func (p Price) GetPayableByRoundingMode(mode string, precision int) Price {
	newPrice := Price{
		currency: p.currency,
	}

	amountForRound := new(big.Float).Copy(&p.amount)
	negative := int64(1)
	if p.IsNegative() {
		negative = -1
	}

	amountTruncatedFloat, _ := new(big.Float).Mul(amountForRound, p.precisionF(precision)).Float64()
	integerPart, fractionalPart := math.Modf(amountTruncatedFloat)
	amountTruncatedInt := int64(integerPart)
	valueAfterPrecision := (math.Round(fractionalPart*1000) / 100) * float64(negative)
	if amountTruncatedFloat >= float64(math.MaxInt64) {
		// will not work if we are already above MaxInt - so we return unrounded price:
		newPrice.amount = p.amount
		return newPrice
	}

	switch mode {
	case RoundingModeCeil:
		if negative == 1 && valueAfterPrecision > 0 {
			amountTruncatedInt = amountTruncatedInt + negative
		}
	case RoundingModeHalfUp:
		if valueAfterPrecision >= 5 {
			amountTruncatedInt = amountTruncatedInt + negative
		}
	case RoundingModeHalfDown:
		if valueAfterPrecision > 5 {
			amountTruncatedInt = amountTruncatedInt + negative
		}
	case RoundingModeFloor:
		if negative == -1 && valueAfterPrecision > 0 {
			amountTruncatedInt = amountTruncatedInt + negative
		}
	default:
		// nothing to round
	}

	amountRounded := new(big.Float).Quo(new(big.Float).SetInt64(amountTruncatedInt), p.precisionF(precision))
	newPrice.amount = *amountRounded
	return newPrice
}

// precisionF returns big.Float from int
func (p Price) precisionF(precision int) *big.Float {
	return new(big.Float).SetInt64(int64(precision))
}

// precisionF - 10 * n - n is the amount of decimal numbers after comma
// - can be currency specific (for now defaults to 2)
// - TODO - use currency configuration or registry
func (p Price) payableRoundingPrecision() (string, int) {
	if strings.ToLower(p.currency) == "miles" || strings.ToLower(p.currency) == "points" {
		return RoundingModeFloor, int(1)
	}
	return RoundingModeHalfUp, int(100)
}

// SplitInPayables returns "count" payable prices (each rounded) that in sum matches the given price
//   - Given a price of 12.456 (Payable 12,46)  - Splitted in 6 will mean: 6 * 2.076
//   - but having them payable requires rounding them each (e.g. 2.07) which would mean we have 0.03 difference (=12,45-6*2.07)
//   - so that the sum is as close as possible to the original value   in this case the correct return will be:
//   - 2.07 + 2.07+2.08 +2.08 +2.08 +2.08
func (p Price) SplitInPayables(count int) ([]Price, error) {
	if count <= 0 {
		return nil, errors.New("split must be higher than zero")
	}
	// guard clause invert negative values
	_, precision := p.payableRoundingPrecision()
	amount := p.GetPayable().Amount()
	// we have to invert negative numbers, otherwise split is not correct
	if p.IsNegative() {
		amount = p.GetPayable().Inverse().Amount()
	}
	amountToMatchFloat, _ := new(big.Float).Mul(amount, p.precisionF(precision)).Float64()
	amountToMatchInt := int64(amountToMatchFloat)

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
		_, precision := p.payableRoundingPrecision()
		splittedAmount := splittedAmounts[i]
		// invert prices again to keep negative values
		if p.IsNegative() {
			splittedAmount *= -1
		}
		prices[i] = NewFromInt(splittedAmount, precision, p.Currency())
	}

	return prices, nil
}

// Clone returns a copy of the price - the amount gets Excat acc
func (p Price) Clone() Price {
	return Price{
		amount:   *new(big.Float).Set(&p.amount),
		currency: p.currency,
	}
}

// Currency returns currency
func (p Price) Currency() string {
	return p.currency
}

// Amount returns exact amount as bigFloat
func (p Price) Amount() *big.Float {
	return &p.amount
}

// SumAll returns new price with sum of all given prices
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

// MarshalJSON implements interface required by json marshal
func (p Price) MarshalJSON() (data []byte, err error) {
	type priceJSON struct {
		Amount   string
		Currency string
	}

	pn := priceJSON{
		Amount:   strconv.FormatFloat(p.GetPayable().FloatAmount(), 'f', 2, 64),
		Currency: p.currency,
	}

	r, e := json.Marshal(&pn)
	return r, e
}

// MarshalBinary implements interface required by gob
func (p Price) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}

// UnmarshalBinary implements interface required by gob.
// Modifies the receiver so it must take a pointer receiver!
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

// UnmarshalJSON implements encode Unmarshaler
func (p *Price) UnmarshalJSON(data []byte) error {
	return p.UnmarshalBinary(data)
}

// Add the given Charge to the current Charge and returns a new Charge
func (p Charge) Add(add Charge) (Charge, error) {
	if p.Type != add.Type {
		return Charge{}, errors.New("charge type mismatch")
	}
	newPrice, err := p.Price.Add(add.Price)
	if err != nil {
		return Charge{}, err
	}
	p.Price = newPrice

	newPrice, err = p.Value.Add(add.Value)
	if err != nil {
		return Charge{}, err
	}
	p.Value = newPrice
	return p, nil
}

// GetPayable rounds the charge
func (p Charge) GetPayable() Charge {
	p.Value = p.Value.GetPayable()
	p.Price = p.Price.GetPayable()
	return p
}

// Mul the given Charge and returns a new Charge
func (p Charge) Mul(qty int) Charge {
	p.Price = p.Price.Multiply(qty)
	p.Value = p.Value.Multiply(qty)
	return p
}

// NewCharges creates a new Charges object
func NewCharges(chargesByType map[string]Charge) *Charges {
	charges := addChargeQualifier(chargesByType)
	return &charges
}

// HasType returns a true if any charges include a charge with given type
func (c Charges) HasType(ctype string) bool {
	for qualifier := range c.chargesByQualifier {
		if qualifier.Type == ctype {
			return true
		}
	}
	return false
}

// GetByType returns a charge of given type. If it was not found a Zero amount
// is returned and the second return value is false
// sums up charges by a certain type if there are multiple
func (c Charges) GetByType(ctype string) (Charge, bool) {
	// guard in case type is not available
	if !c.HasType(ctype) {
		return Charge{}, false
	}
	result := Charge{
		Type: ctype,
	}
	// sum up all charges with certain type to one charge
	for qualifier, charge := range c.chargesByQualifier {
		if qualifier.Type == ctype {
			result, _ = result.Add(charge)
		}
	}
	return result, true
}

// HasChargeQualifier returns a true if any charges include a charge with given type
// and concrete key values provided by additional
func (c Charges) HasChargeQualifier(qualifier ChargeQualifier) bool {
	if _, ok := c.chargesByQualifier[qualifier]; ok {
		return true
	}
	return false
}

// GetByChargeQualifier returns a charge of given qualifier.
// If it was not found a Zero amount is returned and the second return value is false
func (c Charges) GetByChargeQualifier(qualifier ChargeQualifier) (Charge, bool) {
	// guard in case type is not available
	if !c.HasChargeQualifier(qualifier) {
		return Charge{}, false
	}

	if charge, ok := c.chargesByQualifier[qualifier]; ok {
		return charge, true
	}
	return Charge{}, false
}

// GetByChargeQualifierForced returns a charge of given qualifier.
// If it was not found a Zero amount is returned. This method might be useful to call in View (template) directly.
func (c Charges) GetByChargeQualifierForced(qualifier ChargeQualifier) Charge {
	result, ok := c.GetByChargeQualifier(qualifier)
	if !ok {
		return Charge{}
	}
	return result
}

// GetByTypeForced returns a charge of given type. If it was not found a Zero amount is returned.
// This method might be useful to call in View (template) directly where you need one return value
// sums up charges by a certain type if there are multiple
func (c Charges) GetByTypeForced(ctype string) Charge {
	result, ok := c.GetByType(ctype)
	if !ok {
		return Charge{}
	}
	return result
}

// GetAllCharges returns all charges
func (c Charges) GetAllCharges() map[ChargeQualifier]Charge {
	return c.chargesByQualifier
}

// GetAllByType returns all charges of type
func (c Charges) GetAllByType(ctype string) map[ChargeQualifier]Charge {
	chargesByType := make(map[ChargeQualifier]Charge)

	for qualifier, charge := range c.chargesByQualifier {
		if qualifier.Type == ctype {
			chargesByType[ChargeQualifier{
				qualifier.Type,
				qualifier.Reference,
			}] = charge
		}
	}

	return chargesByType
}

// Add returns new Charges with the given added
func (c Charges) Add(toadd Charges) Charges {
	if c.chargesByQualifier == nil {
		c.chargesByQualifier = make(map[ChargeQualifier]Charge)
	}
	for addk, addCharge := range toadd.chargesByQualifier {
		if existingCharge, ok := c.chargesByQualifier[addk]; ok {
			chargeSum, _ := existingCharge.Add(addCharge)
			c.chargesByQualifier[addk] = chargeSum.GetPayable()
		} else {
			c.chargesByQualifier[addk] = addCharge
		}
	}
	return c
}

// AddCharge returns new Charges with the given Charge added
func (c Charges) AddCharge(toadd Charge) Charges {
	if c.chargesByQualifier == nil {
		c.chargesByQualifier = make(map[ChargeQualifier]Charge)
	}
	qualifier := ChargeQualifier{
		Type:      toadd.Type,
		Reference: toadd.Reference,
	}
	if existingCharge, ok := c.chargesByQualifier[qualifier]; ok {
		chargeSum, _ := existingCharge.Add(toadd)
		c.chargesByQualifier[qualifier] = chargeSum.GetPayable()
	} else {
		c.chargesByQualifier[qualifier] = toadd
	}

	return c
}

// Mul returns new Charges with the given multiplied
func (c Charges) Mul(qty int) Charges {
	if c.chargesByQualifier == nil {
		return c
	}
	for t, charge := range c.chargesByQualifier {
		c.chargesByQualifier[t] = charge.Mul(qty)
	}
	return c
}

// Items returns all charges items
func (c Charges) Items() []Charge {
	var charges []Charge

	for _, charge := range c.chargesByQualifier {
		charges = append(charges, charge)
	}

	return charges
}

// addChargeQualifier parse string keys to charge qualifier for backwards compatibility
func addChargeQualifier(chargesByType map[string]Charge) Charges {
	withQualifier := make(map[ChargeQualifier]Charge)
	for chargeType, charge := range chargesByType {
		qualifier := ChargeQualifier{
			Type:      chargeType,
			Reference: charge.Reference,
		}
		withQualifier[qualifier] = charge
	}
	return Charges{chargesByQualifier: withQualifier}
}
