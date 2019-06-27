package cart

import (
	"encoding/json"
	"errors"
	"strings"

	price "flamingo.me/flamingo-commerce/v3/price/domain"
)

const (
	splitQualifierSeparator = "-"
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
		ChargeType string
		Method     string
	}

	//PaymentSplit - the Charges qualified by Type and PaymentMethod
	PaymentSplit map[SplitQualifier]price.Charge

	//PaymentSplitByItem - simelar to value object that contains items of the different possible types, that have a price
	PaymentSplitByItem struct {
		CartItems     map[string]PaymentSplit
		ShippingItems map[string]PaymentSplit
		TotalItems    map[string]PaymentSplit
	}

	//PaymentSplitByItemBuilder - Builder to get valid PaymentSplitByItem instances
	PaymentSplitByItemBuilder struct {
		inBuilding *PaymentSplitByItem
	}

	// DefaultPaymentSelection value object - that implements the PaymentSelection interface
	DefaultPaymentSelection struct {
		//GatewayProp - the selected Gateway
		GatewayProp      string
		ChargedItemsProp PaymentSplitByItem
	}
)

//NewSimplePaymentSelection - returns a PaymentSelection that can be used to update the cart.
// 	multiple charges by item are not used here: The complete grandtotal is selected to be payed in one charge with the given paymentgateway and paymentmethod
func NewSimplePaymentSelection(gateway string, method string, pricedItems PricedItems) PaymentSelection {
	selection := DefaultPaymentSelection{
		GatewayProp: gateway,
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
	selection.ChargedItemsProp = builder.Build()
	return selection
}

// NewPaymentSelection - with the passed PaymentSplitByItem
func NewPaymentSelection(gateway string, chargedItems PaymentSplitByItem) PaymentSelection {
	selection := DefaultPaymentSelection{
		GatewayProp:      gateway,
		ChargedItemsProp: chargedItems,
	}
	return selection
}

//Gateway - returns the selected Gateway code
func (d DefaultPaymentSelection) Gateway() string {
	return d.GatewayProp
}

//CartSplit - the selected split per ChargeType and PaymentMethod
func (d DefaultPaymentSelection) CartSplit() PaymentSplit {
	return d.ChargedItemsProp.Sum()
}

//ItemSplit - the selected split per ChargeType and PaymentMethod
func (d DefaultPaymentSelection) ItemSplit() PaymentSplitByItem {
	return d.ChargedItemsProp
}

//TotalValue - returns Valued price sum
func (d DefaultPaymentSelection) TotalValue() price.Price {
	return d.ChargedItemsProp.Sum().TotalValue()
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
	for _, itemSplit := range c.CartItems {
		addToSum(itemSplit)
	}
	for _, itemSplit := range c.ShippingItems {
		addToSum(itemSplit)
	}
	for _, itemSplit := range c.TotalItems {
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

// MarshalJSON serialize to json
func (s PaymentSplit) MarshalJSON() ([]byte, error) {
	result := make(map[string]price.Charge)
	for qualifier, charge := range s {
		// explicit method and chargetype is necessary, otherwise keys could be overwritten
		if qualifier.Method == "" || qualifier.ChargeType == "" {
			return nil, errors.New("Method or ChargeType is empty")
		}
		// SplitQualifier is parsed to a string method___chargetype
		result[qualifier.Method+splitQualifierSeparator+qualifier.ChargeType] = charge
	}
	return json.Marshal(result)
}

// UnmarshalJSON deserialize from json
func (s *PaymentSplit) UnmarshalJSON(data []byte) error {
	var input map[string]price.Charge
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	result := PaymentSplit{}
	// parse string method___chargetype back to splitqualifier
	for key, charge := range input {
		splitted := strings.Split(key, splitQualifierSeparator)
		// guard in case cannot be splitted
		if len(splitted) < 2 {
			return errors.New("SplitQualifier cannot be parsed for paymentsplit")
		}
		qualifier := SplitQualifier{
			Method:     splitted[0],
			ChargeType: splitted[1],
		}
		result[qualifier] = charge
	}
	*s = result
	return nil
}

//AddCartItem - adds a cartitems charge to the PaymentSplitByItem
func (pb *PaymentSplitByItemBuilder) AddCartItem(id string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	if pb.inBuilding.CartItems[id] == nil {
		pb.inBuilding.CartItems[id] = make(PaymentSplit)
	}
	pb.inBuilding.CartItems[id][SplitQualifier{
		Method:     method,
		ChargeType: charge.Type,
	}] = charge
	return pb
}

//AddShippingItem - adds shipping charge
func (pb *PaymentSplitByItemBuilder) AddShippingItem(deliveryCode string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	if pb.inBuilding.ShippingItems[deliveryCode] == nil {
		pb.inBuilding.ShippingItems[deliveryCode] = make(PaymentSplit)
	}
	pb.inBuilding.ShippingItems[deliveryCode][SplitQualifier{
		Method:     method,
		ChargeType: charge.Type,
	}] = charge
	return pb
}

//AddTotalItem - adds totalitem charge
func (pb *PaymentSplitByItemBuilder) AddTotalItem(totalType string, method string, charge price.Charge) *PaymentSplitByItemBuilder {
	pb.init()
	if pb.inBuilding.TotalItems[totalType] == nil {
		pb.inBuilding.TotalItems[totalType] = make(PaymentSplit)
	}
	pb.inBuilding.TotalItems[totalType][SplitQualifier{
		Method:     method,
		ChargeType: charge.Type,
	}] = charge
	return pb
}

//Build - returns the instance of PaymentSplitByItem
func (pb *PaymentSplitByItemBuilder) Build() PaymentSplitByItem {
	pb.init()
	return *pb.inBuilding
}

func (pb *PaymentSplitByItemBuilder) init() {
	if pb.inBuilding != nil {
		return
	}
	pb.inBuilding = &PaymentSplitByItem{
		CartItems:     make(map[string]PaymentSplit),
		ShippingItems: make(map[string]PaymentSplit),
		TotalItems:    make(map[string]PaymentSplit),
	}
}
