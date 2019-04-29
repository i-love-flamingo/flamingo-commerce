package cart

import (
	price "flamingo.me/flamingo-commerce/v3/price/domain"
)

type (

	// PaymentSelection value object - that represents the payment selection on the cart
	PaymentSelection interface {
		Gateway() string
		//ChargeSplits - the selected split per ChargeType and PaymentMethod
		CartSplit() PaymentSplit
		//ChargeSplits - the selected split per ChargeType and PaymentMethod
		ItemSplit() ChargedItems
		TotalValue() price.Price
	}

	//SplitQualifier qualifies by Type and PaymentMethod
	SplitQualifier struct {
		chargeType string
		method     string
	}

	//PaymentSplit - the Charges qualified by Type and PaymentMethod
	PaymentSplit map[SplitQualifier]price.Charge

	//ChargedItems - simelar to value object that contains items of the different possible types, that have a price
	ChargedItems struct {
		cartItems map[string]PaymentSplit
		shippingItems map[string]PaymentSplit
		totalItems map[string]PaymentSplit
	}

	// defaultSelection value object - that implements the PaymentSelection interface
	defaultSelection struct {
		//Gateway - the selected Gateway
		gateway string
		chargedItems ChargedItems
	}




	//CartChargeAssignment.GetForCartItem(itemId) map[string]Charge
	//CartChargeAssignment.GetForShippingItem(itemId) map[string]Charge
	//CartChargeAssignment.GetForTotalItem(itemId) map[string]Charge
	//CartChargeAssignment.GroupedSum() map[string]Charge


)

//NewSimplePaymentSelection - returns a PaymentSelection that can be used to update the cart.
// 	multiple charges to pay the cart are not used here: The complete grandtotal is selected to be payed in one charge with the given paymentgateway and paymentmethod
func NewSimplePaymentSelection(gateway string, method string, pricedItems PricedItems) PaymentSelection {
	selection := defaultSelection{
		gateway: gateway,
		chargedItems:ChargedItems{
			cartItems: make(map[string]PaymentSplit),
			shippingItems: make(map[string]PaymentSplit),
			totalItems: make(map[string]PaymentSplit),
		},
	}


	//addPrice - adds the price as Main Charge to the given Split
	addPrice := func(items PaymentSplit,itemprice price.Price) {
		items[SplitQualifier{
			chargeType: price.ChargeTypeMain,
			method:method,
		}] = price.Charge{
			Price: itemprice,
			Value:itemprice,
			Type:price.ChargeTypeMain,
		}
	}

	for k, itemPrice := range pricedItems.CartItems() {
		addPrice(selection.chargedItems.cartItems[k],itemPrice)
	}
	for k, itemPrice := range pricedItems.ShippingItems() {
		addPrice(selection.chargedItems.shippingItems[k],itemPrice)

	}
	for k, itemPrice := range pricedItems.TotalItems() {
		addPrice(selection.chargedItems.totalItems[k],itemPrice)
	}
	return selection
}

//Gateway - returns the selected Gateway code
func (d defaultSelection) Gateway() string {
	return d.gateway
}
//CartSplit - the selected split per ChargeType and PaymentMethod
func (d defaultSelection)  CartSplit() PaymentSplit {
	return d.chargedItems.Sum()
}
//ItemSplit - the selected split per ChargeType and PaymentMethod
func (d defaultSelection)  ItemSplit() ChargedItems{
	return d.chargedItems
}

func (d defaultSelection) TotalValue() price.Price{
	return d.chargedItems.Sum().TotalValue()
}

//Sum - the resulting Split after sum all the included item split
func (c ChargedItems) Sum() PaymentSplit {
	sum := make(PaymentSplit)
	addToSum := func(splits PaymentSplit) {
		for qualifier, charge := range splits {
			_, ok := sum[qualifier]
			if ok {
				sum[qualifier],_ = sum[qualifier].Add(charge)
			} else {
				sum[qualifier] = charge
			}
		}
	}
	for _, itemSplit := range c.cartItems {
		addToSum(itemSplit)
	}
	for _, itemSplit := range c.shippingItems {
		addToSum(itemSplit)
	}
	for _, itemSplit := range c.totalItems {
		addToSum(itemSplit)
	}
	return sum
}

//TotalValue returns the sum of the valued Price in the included Charges in this Split
func (s PaymentSplit) TotalValue() price.Price {
	var prices []price.Price
	for _,v := range s {
		prices = append(prices,v.Value)
	}
	sum, _ :=  price.SumAll(prices...)
	return sum
}