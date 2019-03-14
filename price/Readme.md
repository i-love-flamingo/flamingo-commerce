# Price Module

## Price Type:
The price module offers a Price Type with useful methods.

Price calculation is not a trivial topic and multiple solutions exist. 
The implementation details of the price object is:

* internal we use big.Float to hold the amount, this is to be able to calculate excat
* however a float like representation of an amount cannot be payed, that is why the price has a method "GetPayablePrice" that returns a Price that can be payed in the given currency, using correct roundings and amount representation


### Example:

```

//get 2.45 EUR Price
price := NewFromInt(245,100,"EUR")

rowPrice := price.Multiply(10)

//10% discount
discountedRowPrice := rowPrice.Discounted(10.0)

//What need to be payed:
priceToPay := discountedRowPrice.GetPayable()

//what need to be payed by item()
singleItemsPrices := discountedRowPrice.SplitInPayables(10)


//you can also set price from float:
price2 := NewFromFloat(2.45,"EUR")

//Be aware that

price.Equals(price2) may be false but due to float arithmentics
but
price.GetPayable().Equals(price2.GetPayable()) will be true
```

## Charge Type:
Represents a price together with a type.
Can be used in places where you need to give the price value a certain extra semantic information.