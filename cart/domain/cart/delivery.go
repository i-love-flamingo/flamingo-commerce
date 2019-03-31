package cart

import (
	"encoding/json"
	"errors"
	"time"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

type (

	// Delivery - represents the DeliveryInfo and the assigned Items
	Delivery struct {
		//DeliveryInfo - The details for this delivery - normaly completed during checkout
		DeliveryInfo DeliveryInfo
		//Cartitems - list of cartitems
		Cartitems []Item
		//ShippingItem	- The Shipping Costs that may be involved in this delivery
		ShippingItem ShippingItem
	}

	// DeliveryInfo - represents the Delivery
	DeliveryInfo struct {
		// Code - is a project specific idendifier for the Delivery - you need it for the AddToCart Request for example
		// The code can follow the convention in the Readme: Type_Method_LocationType_LocationCode
		Code string
		//Type - The Type of the Delivery - e.g. delivery or pickup - this might trigger different workflows
		Workflow string
		//Method - The shippingmethod something that is project specific and that can mean different delivery qualities with different deliverycosts
		Method string
		//Carrier - Optional the name of the Carrier that should be responsible for executing the delivery
		Carrier string
		//DeliveryLocation The target Location for the delivery
		DeliveryLocation DeliveryLocation
		//DesiredTime - Optional - the desired time of the delivery
		DesiredTime time.Time
		//AdditionalData  - Possibility for key value based information on the delivery - can be used flexible by each project
		AdditionalData map[string]string
		//AdditionalDeliveryInfos - similar to AdditionalData this can be used to store "any" other object on a delivery encoded as json.RawMessage
		AdditionalDeliveryInfos map[string]json.RawMessage
	}

	// ShippingItem value object
	ShippingItem struct {
		Title          string
		Price          priceDomain.Price
		TaxAmount      priceDomain.Price
		DiscountAmount priceDomain.Price
	}

	//AdditionalDeliverInfo is an interface that allows to store "any" additional objects on the cart
	// see DeliveryInfoUpdateCommand
	AdditionalDeliverInfo interface {
		Marshal() (json.RawMessage, error)
		Unmarshal(json.RawMessage) error
	}

	// DeliveryLocation value object
	DeliveryLocation struct {
		//Type - the type of the delivery - use some of the constant defined in the package like DeliverylocationTypeAddress
		Type string
		//Address -  (only relevant for type adress)
		Address *Address
		//UseBillingAddress - the adress should be taken from billing (only relevant for type adress)
		UseBillingAddress bool
		//Code - optional idendifier of this location/destination - is used in special destination Types
		Code string
	}

	//DeliveryBuilder - the Builder (factory) to build new deliveries by making sure the invariants are ok
	DeliveryBuilder struct {
		deliveryInBuilding *Delivery
	}

	// DeliveryBuilderProvider should be used to create a Delivery
	DeliveryBuilderProvider func() *DeliveryBuilder
)

const (
	//DeliveryWorkflowPickup - constant for common delivery workflows
	DeliveryWorkflowPickup = "pickup"
	//DeliveryWorkflowDelivery - workflow constant
	DeliveryWorkflowDelivery = "delivery"
	//DeliveryWorkflowUnspecified - workflow constant
	DeliveryWorkflowUnspecified = "unspecified"

	//DeliverylocationTypeUnspecified - constant
	DeliverylocationTypeUnspecified = "unspecified"
	//DeliverylocationTypeCollectionpoint - constant
	DeliverylocationTypeCollectionpoint = "collection-point"
	//DeliverylocationTypeStore - constant
	DeliverylocationTypeStore = "store"
	//DeliverylocationTypeAddress - constant
	DeliverylocationTypeAddress = "address"
	//DeliverylocationTypeFreightstation - constant
	DeliverylocationTypeFreightstation = "freight-station"
)

//LoadAdditionalInfo - returns the additional Data
func (di *DeliveryInfo) LoadAdditionalInfo(key string, info AdditionalDeliverInfo) error {

	if di.AdditionalDeliveryInfos == nil {
		return ErrAdditionalInfosNotFound
	}
	if val, ok := di.AdditionalDeliveryInfos[key]; ok {
		return info.Unmarshal(val)
	}
	return ErrAdditionalInfosNotFound
}

//SubTotalGross - returns SubTotalGross
func (d Delivery) SubTotalGross() priceDomain.Price {
	var prices []priceDomain.Price
	for _, item := range d.Cartitems {
		prices = append(prices, item.RowPriceGross)
	}
	result, _ := priceDomain.SumAll(prices...)
	return result
}

//SumRowTaxes - returns SumRowTaxes
func (d Delivery) SumRowTaxes() Taxes {
	var taxes Taxes
	for _, item := range d.Cartitems {
		for _, tax := range item.RowTaxes {
			taxes = taxes.AddTaxWithMerge(tax)
		}
	}
	return taxes
}

//SumTotalTaxAmount - returns SumTotalTaxAmount
func (d Delivery) SumTotalTaxAmount() priceDomain.Price {
	var prices []priceDomain.Price
	for _, item := range d.Cartitems {
		prices = append(prices, item.TotalTaxAmount())
	}
	result, _ := priceDomain.SumAll(prices...)
	return result
}

//SubTotalNet - returns SubTotalNet
func (d Delivery) SubTotalNet() priceDomain.Price {
	var prices []priceDomain.Price
	for _, item := range d.Cartitems {
		prices = append(prices, item.RowPriceNet)
	}
	result, _ := priceDomain.SumAll(prices...)
	return result
}

//SumTotalDiscountAmount - returns SumTotalDiscountAmount
func (d Delivery) SumTotalDiscountAmount() priceDomain.Price {
	var prices []priceDomain.Price
	for _, item := range d.Cartitems {
		prices = append(prices, item.TotalDiscountAmount())
	}
	result, _ := priceDomain.SumAll(prices...)
	return result
}

//SumNonItemRelatedDiscountAmount returns SumNonItemRelatedDiscountAmount
func (d Delivery) SumNonItemRelatedDiscountAmount() priceDomain.Price {
	var prices []priceDomain.Price
	for _, item := range d.Cartitems {
		prices = append(prices, item.NonItemRelatedDiscountAmount())
	}
	result, _ := priceDomain.SumAll(prices...)
	return result
}

//SumItemRelatedDiscountAmount - returns SumItemRelatedDiscountAmount
func (d Delivery) SumItemRelatedDiscountAmount() priceDomain.Price {
	var prices []priceDomain.Price
	for _, item := range d.Cartitems {
		prices = append(prices, item.ItemRelatedDiscountAmount())
	}
	result, _ := priceDomain.SumAll(prices...)
	return result
}

//SubTotalGrossWithDiscounts returns SubTotalGrossWithDiscounts
func (d Delivery) SubTotalGrossWithDiscounts() priceDomain.Price {
	price, _ := d.SubTotalGross().Add(d.SumTotalDiscountAmount())
	return price
}

//SubTotalNetWithDiscounts - returns SubTotalNet With Discounts
func (d Delivery) SubTotalNetWithDiscounts() priceDomain.Price {
	price, _ := d.SubTotalNet().Add(d.SumTotalDiscountAmount())
	return price
}

//HasItems - returns true if there are items under the delivery
func (d Delivery) HasItems() bool {
	return len(d.Cartitems) > 0
}

//Copy - use to set the values for the new delivery from an existing delivery by copying it
func (f *DeliveryBuilder) Copy(d *Delivery) *DeliveryBuilder {
	f.init()
	f.deliveryInBuilding.Cartitems = d.Cartitems
	f.deliveryInBuilding.ShippingItem = d.ShippingItem
	f.deliveryInBuilding.DeliveryInfo = d.DeliveryInfo

	return f
}

//AddItem adds an item to the delivery
func (f *DeliveryBuilder) AddItem(i Item) *DeliveryBuilder {
	f.init()
	f.deliveryInBuilding.Cartitems = append(f.deliveryInBuilding.Cartitems, i)
	return f
}

//SetShippingItem - sets the delivery ShippingItem
func (f *DeliveryBuilder) SetShippingItem(i ShippingItem) *DeliveryBuilder {
	f.init()
	f.deliveryInBuilding.ShippingItem = i
	return f
}

//SetDeliveryInfo - sets DeliveryInfo
func (f *DeliveryBuilder) SetDeliveryInfo(i DeliveryInfo) *DeliveryBuilder {
	f.init()
	f.deliveryInBuilding.DeliveryInfo = i
	return f
}

//SetDeliveryCode - sets the deliverycode (dont need to be called if SetDeliveryInfo has a code set already)
func (f *DeliveryBuilder) SetDeliveryCode(code string) *DeliveryBuilder {
	f.init()
	f.deliveryInBuilding.DeliveryInfo.Code = code
	return f
}

//Build - main Factory method
func (f *DeliveryBuilder) Build() (*Delivery, error) {
	if f.deliveryInBuilding == nil {
		return nil, errors.New("Nothing in building")
	}
	if f.deliveryInBuilding.DeliveryInfo.Code == "" {
		return nil, errors.New("DeliveryInfo.Code is not allowed empty")
	}

	return f.deliveryInBuilding, nil
}

func (f *DeliveryBuilder) init() {
	if f.deliveryInBuilding == nil {
		f.deliveryInBuilding = &Delivery{}
	}
}

func (f *DeliveryBuilder) reset() {
	f.deliveryInBuilding = nil
}
