package cart

import (
	"encoding/json"
	price "flamingo.me/flamingo-commerce/v3/price/domain"
	"log"
)

type (
	// PaymentSelection value object - that represents the payment selection on the cart
	PaymentSelection interface {
		Gateway() string
		//ChargeSplits - the selected split per ChargeType and PaymentMethod
		CartSplit() PaymentSplit
		//ChargeSplits - the selected split per ChargeType and PaymentMethod
		ItemSplit() PaymentSplitByItem
		TotalValue() price.Price
	}

	//SplitQualifier qualifies by Type and PaymentMethod
	SplitQualifier struct {
		chargeType string
		method     string
	}

	//PaymentSplit - the Charges qualified by Type and PaymentMethod
	PaymentSplit map[SplitQualifier]price.Charge

	//PaymentSplitByItem - simelar to value object that contains items of the different possible types, that have a price
	PaymentSplitByItem struct {
		cartItems     map[string]PaymentSplit
		shippingItems map[string]PaymentSplit
		totalItems    map[string]PaymentSplit
	}

	PaymentSplitByItemBuilder struct {
		inBuilding *PaymentSplitByItem
	}

	// DefaultPaymentSelection value object - that implements the PaymentSelection interface
	DefaultPaymentSelection struct {
		//Gateway - the selected Gateway
		gateway      string
		chargedItems PaymentSplitByItem
	}

	//CartChargeAssignment.GetForCartItem(itemId) map[string]Charge
	//CartChargeAssignment.GetForShippingItem(itemId) map[string]Charge
	//CartChargeAssignment.GetForTotalItem(itemId) map[string]Charge
	//CartChargeAssignment.GroupedSum() map[string]Charge

)

//NewSimplePaymentSelection - returns a PaymentSelection that can be used to update the cart.
// 	multiple charges by item are not used here: The complete grandtotal is selected to be payed in one charge with the given paymentgateway and paymentmethod
func NewSimplePaymentSelection(gateway string, method string, pricedItems PricedItems) PaymentSelection {
	selection := DefaultPaymentSelection{
		gateway: gateway,
	}
	builder := PaymentSplitByItemBuilder{}

	for k, itemPrice := range pricedItems.CartItems() {
		builder.AddCartItem(k, method, price.Charge{
			Price: itemPrice,
			Value: itemPrice,
			Type:  price.ChargeTypeMain,
		})

	}
	for k, itemPrice := range pricedItems.ShippingItems() {
		builder.AddShippingItem(k, method, price.Charge{
			Price: itemPrice,
			Value: itemPrice,
			Type:  price.ChargeTypeMain,
		})

	}
	for k, itemPrice := range pricedItems.TotalItems() {
		builder.AddTotalItem(k, method, price.Charge{
			Price: itemPrice,
			Value: itemPrice,
			Type:  price.ChargeTypeMain,
		})
	}
	selection.chargedItems = builder.Build()
	return selection
}

// NewPaymentSelection - with the passed PaymentSplitByItem
func NewPaymentSelection(gateway string, chargedItems PaymentSplitByItem) PaymentSelection {
	selection := DefaultPaymentSelection{
		gateway:      gateway,
		chargedItems: chargedItems,
	}
	return selection
}

//Gateway - returns the selected Gateway code
func (d DefaultPaymentSelection) Gateway() string {
	return d.gateway
}

//CartSplit - the selected split per ChargeType and PaymentMethod
func (d DefaultPaymentSelection) CartSplit() PaymentSplit {
	return d.chargedItems.Sum()
}

//ItemSplit - the selected split per ChargeType and PaymentMethod
func (d DefaultPaymentSelection) ItemSplit() PaymentSplitByItem {
	return d.chargedItems
}

func (d DefaultPaymentSelection) TotalValue() price.Price {
	return d.chargedItems.Sum().TotalValue()
}

//Sum - the resulting Split after sum all the included item split
func (c PaymentSplitByItem) Sum() PaymentSplit {
	sum := make(PaymentSplit)
	addToSum := func(splits PaymentSplit) {
		for qualifier, charge := range splits {
			_, ok := sum[qualifier]
			if ok {
				sum[qualifier], _ = sum[qualifier].Add(charge)
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
	for _, v := range s {
		prices = append(prices, v.Value)
	}
	sum, _ := price.SumAll(prices...)
	return sum
}

//ChargesByType returns Charges (a list of Charges summed by Type)
func (s PaymentSplit) ChargesByType() price.Charges {
	charges := price.Charges{}
	for _, charge := range s {
		charges = charges.AddCharge(charge)
	}
	return charges
}

//ChargeType - returns the ChargeType of the Qualifier
func (s SplitQualifier) ChargeType() string {
	return s.chargeType
}

//Method - return Method
func (s SplitQualifier) Method() string {
	return s.method
}

func (p PaymentSplitByItem) ShippingItems() map[string]PaymentSplit {
	return p.shippingItems
}

func (p PaymentSplitByItem) TotalItems() map[string]PaymentSplit {
	return p.totalItems
}

func (p PaymentSplitByItem) CartItems() map[string]PaymentSplit {
	return p.cartItems
}

func (pb *PaymentSplitByItemBuilder) AddCartItem(id string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	log.Printf("%#v",pb.inBuilding)
	if pb.inBuilding.cartItems[id] == nil {
		pb.inBuilding.cartItems[id] = make(PaymentSplit)
	}
	pb.inBuilding.cartItems[id][SplitQualifier{
		method:     method,
		chargeType: charge.Type,
	}] = charge
	return pb
}

func (pb *PaymentSplitByItemBuilder) AddShippingItem(deliveryCode string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	if pb.inBuilding.shippingItems[deliveryCode] == nil {
		pb.inBuilding.shippingItems[deliveryCode] = make(PaymentSplit)
	}
	pb.inBuilding.shippingItems[deliveryCode][SplitQualifier{
		method:     method,
		chargeType: charge.Type,
	}] = charge
	return pb
}

func (pb *PaymentSplitByItemBuilder) AddTotalItem(totalType string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	if pb.inBuilding.totalItems[totalType] == nil {
		pb.inBuilding.totalItems[totalType] = make(PaymentSplit)
	}
	pb.inBuilding.totalItems[totalType][SplitQualifier{
		method:     method,
		chargeType: charge.Type,
	}] = charge
	return pb
}

func (pb *PaymentSplitByItemBuilder) Build() PaymentSplitByItem {
	pb.init()
	return *pb.inBuilding
}

func (pb *PaymentSplitByItemBuilder) init() {
	if pb.inBuilding != nil {
		return
	}
	pb.inBuilding = &PaymentSplitByItem{
		cartItems:     make(map[string]PaymentSplit),
		shippingItems: make(map[string]PaymentSplit),
		totalItems:    make(map[string]PaymentSplit),
	}
}



//MarshalBinary - implements interface required by gob
func (d DefaultPaymentSelection) MarshalBinary() (data []byte, err error) {
	return json.Marshal(d)
}

//UnmarshalBinary - implements interace required by gob.
//UnmarshalBinary - modifies the receiver so it must take a pointer receiver!
func (d *DefaultPaymentSelection) UnmarshalBinary(data []byte) error {
	type encodeAbleDefaultPaymentSelection DefaultPaymentSelection
	var encodeAble encodeAbleDefaultPaymentSelection
	err := json.Unmarshal(data, &encodeAble)
	if err != nil {
		return err
	}
	newSelection := DefaultPaymentSelection(encodeAble)
	d = &newSelection
	return nil
}