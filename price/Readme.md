# Price Module

The price module defines value objects to deal with prices and charges.
The following types exist:
* Price: Represents simply a money value. See below for more details. For example: 5€
* Charge: Represents a Price together with its Type. This is used to describe that something need to be paid in a certain currency and represents a certain value. For example it might be you need (or want) to pay the value of 5€ in 500 Loyaltypoints. This can be described in the Charge Type.
* Charges: Is a list of the type Charge. For example to indicate that you are paying to value of 5€ with two Charges: 2€ and 300 Loyaltypoints

## Price Type:
The price module offers a Price Type with useful methods.

Price calculation is not a trivial topic and multiple solutions exist. 
The implementation details of the price object is:

* internally we use big.Float to hold the amount, this is to be able to calculate exactly
* however, a float like representation of an amount cannot be paid, that is why the price has a method "GetPayablePrice" that returns a Price that can be paid in the given currency, using correct rounding and amount representation


### Example:

```go
// get 2.45 EUR Price
price := NewFromInt(245,100,"EUR")

rowPrice := price.Multiply(10)

// 10% discount
discountedRowPrice := rowPrice.Discounted(10.0)

// What needs to be paid:
priceToPay := discountedRowPrice.GetPayable()

// what needs to be paid by item()
singleItemsPrices := discountedRowPrice.SplitInPayables(10)

// you can also set price from float:
price2 := NewFromFloat(2.45,"EUR")
```

Be aware that `price.Equals(price2)` may be false but due to float arithmetic but
`price.GetPayable().Equals(price2.GetPayable())` will be true

#### Rounding Modes

We support a variety of rounding modes. These can be used with the `GetPayableByRoundingMode` function by specifying a mode and precision. The following modes are available:

* `RoundingModeFloor`
  
  Rounding mode to round towards negative infinity.
  _Note that this rounding mode never increases the calculated value._

  | `price` | `price.GetPayableByRoundingMode(RoundingModeFloor, 1)` |
  |---------|--------------------------------------------------------|
  |     5.5 |                                                      5 |
  |     2.5 |                                                      2 |
  |     1.6 |                                                      1 |
  |     1.1 |                                                      1 |
  |     1.0 |                                                      1 |
  |    -1.0 |                                                     -1 |
  |    -1.1 |                                                     -2 |
  |    -1.6 |                                                     -2 |
  |    -2.5 |                                                     -3 |
  |    -5.5 |                                                     -6 |
* `RoundingModeCeil`
  
  Rounding mode to round towards positive infinity.
  _Note that this rounding mode never decreases the calculated value._

  | `price` | `price.GetPayableByRoundingMode(RoundingModeCeil, 1)` |
  |---------|-------------------------------------------------------|
  |     5.5 |                                                     6 |
  |     2.5 |                                                     3 |
  |     1.6 |                                                     2 |
  |     1.1 |                                                     2 |
  |     1.0 |                                                     1 |
  |    -1.0 |                                                    -1 |
  |    -1.1 |                                                    -1 |
  |    -1.6 |                                                    -1 |
  |    -2.5 |                                                    -2 |
  |    -5.5 |                                                    -5 |
* `RoundingModeHalfUp` default mode for `GetPayable()`
  
  Rounding mode to round towards "nearest neighbor" unless both neighbors are equidistant, in which case round up.
  _Note that this is the rounding mode commonly taught at school._

  | `price` | `price.GetPayableByRoundingMode(RoundingModeHalfUp, 1)` |
  |---------|---------------------------------------------------------|
  |     5.5 |                                                       6 |
  |     2.5 |                                                       3 |
  |     1.6 |                                                       2 |
  |     1.1 |                                                       1 |
  |     1.0 |                                                       1 |
  |    -1.0 |                                                      -1 |
  |    -1.1 |                                                      -1 |
  |    -1.6 |                                                      -2 |
  |    -2.5 |                                                      -3 |
  |    -5.5 |                                                      -6 |
  
* `RoundingModeHalfDown`
  
  Rounding mode to round towards "nearest neighbor" unless both neighbors are equidistant, in which case round down.

  | `price` | `price.GetPayableByRoundingMode(RoundingModeHalfDown, 1)` |
  |---------|-----------------------------------------------------------|
  |     5.5 |                                                         5 |
  |     2.5 |                                                         2 |
  |     1.6 |                                                         2 |
  |     1.1 |                                                         1 |
  |     1.0 |                                                         1 |
  |    -1.0 |                                                        -1 |
  |    -1.1 |                                                        -1 |
  |    -1.6 |                                                        -2 |
  |    -2.5 |                                                        -2 |
  |    -5.5 |                                                        -5 |


## Charge:
Represents a price together with a type. A charge has a values price (normally in default currency) and a the price that is paid that might be in a different currency.
Can be used in places where you need to give the price value a certain extra semantic information or to represent something that need to be paid (charged).

## Template Func - Formatting a Price Object

Just use the template function commercePriceFormat like this: `commercePriceFormat(priceObject)` 
The template functions used the configurations of the Flamingo "locale" package. For more details on the configuration options please read there.
