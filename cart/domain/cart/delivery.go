package cart

import (
	"encoding/json"
	"time"

	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// Delivery - represents the DeliveryInfo and the assigned Items
	Delivery struct {
		// DeliveryInfo contains details for this delivery e.g. how and where the delivery should be delivered to
		DeliveryInfo DeliveryInfo
		// Cartitems is the list of items belonging to this delivery
		Cartitems []Item
		// ShippingItem	represent the shipping cost that may be involved in this delivery
		ShippingItem ShippingItem

		// SubTotalGross contains the sum of row gross prices, without shipping/discounts
		SubTotalGross priceDomain.Price
		// SubTotalNet contains the sum of row net prices, without shipping/discounts
		SubTotalNet priceDomain.Price
		// TotalDiscountAmount contains the sum of all discounts (incl. shipping)
		TotalDiscountAmount priceDomain.Price
		// TotalDiscountAmount contains the sum of all discounts (excl. shipping)
		SubTotalDiscountAmount priceDomain.Price
		// NonItemRelatedDiscountAmount contains the sum of discounts that are not related to the item, e.g. a general promo
		NonItemRelatedDiscountAmount priceDomain.Price
		// ItemRelatedDiscountAmount contains the sum of discounts that are related to the item, e.g. promo due to product attribute
		ItemRelatedDiscountAmount priceDomain.Price
		// SubTotalGrossWithDiscounts contains the sum of row gross prices reduced by the applied discounts
		SubTotalGrossWithDiscounts priceDomain.Price
		// SubTotalNetWithDiscounts contains the sum of row net prices reduced by the net value of the applied discounts
		SubTotalNetWithDiscounts priceDomain.Price
		// GrandTotal contains the final price to pay
		GrandTotal priceDomain.Price
	}

	// DeliveryInfo - represents the Delivery
	DeliveryInfo struct {
		// Code is a project specific identifier for the Delivery - you need it for the AddToCart Request for example
		// the code can follow the convention in the Readme: Type_Method_LocationType_LocationCode
		Code string
		// Workflow of the Delivery e.g. delivery or pickup, see DeliveryWorkflowPickup, DeliveryWorkflowDelivery or DeliveryWorkflowUnspecified
		Workflow string
		// Method is the shipping method something that is project specific and that can mean different delivery qualities with different delivery costs
		Method string
		// Carrier optional name of the Carrier that should be responsible for executing the delivery
		Carrier string
		// DeliveryLocation is the target location for the delivery
		DeliveryLocation DeliveryLocation
		// DesiredTime is an optional desired time for the delivery
		DesiredTime time.Time
		// AdditionalData can be used to store project specific information on the delivery
		AdditionalData map[string]string
		// AdditionalDeliveryInfos is similar to AdditionalData but can be used to store "any" other object on a delivery encoded as json.RawMessage
		AdditionalDeliveryInfos map[string]json.RawMessage `swaggerignore:"true"`
	}

	// ShippingItem represents shipping costs that need to be paid by the customer
	ShippingItem struct {
		Title                   string
		PriceNet                priceDomain.Price
		PriceNetWithDiscounts   priceDomain.Price
		PriceGross              priceDomain.Price
		PriceGrossWithDiscounts priceDomain.Price
		TaxAmount               priceDomain.Price
		AppliedDiscounts        AppliedDiscounts
	}

	// AdditionalDeliverInfo is an interface that allows to store "any" additional objects on the cart
	// see DeliveryInfoUpdateCommand
	AdditionalDeliverInfo interface {
		Marshal() (json.RawMessage, error)
		Unmarshal(json.RawMessage) error
	}

	// DeliveryLocation hold information about where the items should be delivered
	DeliveryLocation struct {
		// Type is the type of the delivery - use some of the constant defined in the package like DeliverylocationTypeAddress
		Type string
		// Address contains the address of the delivery location, maybe not relevant if the type is not address
		Address *Address
		// UseBillingAddress if the address should be taken from billing (only relevant for type address)
		UseBillingAddress bool
		// Code is an optional identifier of this location/destination
		Code string
	}
)

const (
	// DeliveryWorkflowPickup constant for pickup delivery workflow
	DeliveryWorkflowPickup = "pickup"
	// DeliveryWorkflowDelivery constant for delivery delivery workflow
	DeliveryWorkflowDelivery = "delivery"
	// DeliveryWorkflowUnspecified constant for an unspecified delivery workflow
	DeliveryWorkflowUnspecified = "unspecified"

	// DeliverylocationTypeUnspecified constant for an unspecified delivery location type
	DeliverylocationTypeUnspecified = "unspecified"
	// DeliverylocationTypeCollectionpoint constant for collection points
	DeliverylocationTypeCollectionpoint = "collection-point"
	// DeliverylocationTypeStore constant for store delivery
	DeliverylocationTypeStore = "store"
	// DeliverylocationTypeAddress constant for deliveries to an address
	DeliverylocationTypeAddress = "address"
	// DeliverylocationTypeFreightstation constant for deliveries to an freight station
	DeliverylocationTypeFreightstation = "freight-station"
)

// SumRowTaxes returns all taxes applied to items of this delivery
func (d Delivery) SumRowTaxes() Taxes {
	var taxes Taxes
	for _, item := range d.Cartitems {
		for _, tax := range item.RowTaxes {
			taxes = taxes.AddTaxWithMerge(tax)
		}
	}
	return taxes
}

// SumTotalTaxAmount returns the sum of all applied item taxes
func (d Delivery) SumTotalTaxAmount() priceDomain.Price {
	prices := make([]priceDomain.Price, 0, len(d.Cartitems)+1)

	prices = append(prices, d.ShippingItem.TaxAmount)

	for _, item := range d.Cartitems {
		prices = append(prices, item.TotalTaxAmount())
	}
	result, _ := priceDomain.SumAll(prices...)

	return result
}

// HasItems returns true if there are items under the delivery
func (d Delivery) HasItems() bool {
	return len(d.Cartitems) > 0
}

// Tax is the Tax of the shipping
func (s ShippingItem) Tax() Tax {
	return Tax{
		Type:   "tax",
		Amount: s.TaxAmount,
	}
}

// LoadAdditionalInfo returns the additional Data
func (di *DeliveryInfo) LoadAdditionalInfo(key string, info AdditionalDeliverInfo) error {
	if di.AdditionalDeliveryInfos == nil {
		return ErrAdditionalInfosNotFound
	}
	if val, ok := di.AdditionalDeliveryInfos[key]; ok {
		return info.Unmarshal(val)
	}
	return ErrAdditionalInfosNotFound
}

// GetAdditionalData returns additional data
func (di DeliveryInfo) GetAdditionalData(key string) string {
	attribute := di.AdditionalData[key]
	return attribute
}

// AdditionalDataKeys lists all available keys
func (di DeliveryInfo) AdditionalDataKeys() []string {
	res := make([]string, len(di.AdditionalData))
	i := 0
	for k := range di.AdditionalData {
		res[i] = k
		i++
	}
	return res
}

// GetAdditionalDeliveryInfo returns additional delivery info
func (di DeliveryInfo) GetAdditionalDeliveryInfo(key string) json.RawMessage {
	attribute := di.AdditionalDeliveryInfos[key]
	return attribute
}

// AdditionalDeliveryInfoKeys lists all available keys
func (di DeliveryInfo) AdditionalDeliveryInfoKeys() []string {
	res := make([]string, len(di.AdditionalDeliveryInfos))
	i := 0
	for k := range di.AdditionalDeliveryInfos {
		res[i] = k
		i++
	}
	return res
}
