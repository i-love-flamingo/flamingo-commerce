package cart

import (
	"encoding/json"
	"time"

	"flamingo.me/flamingo-commerce/v3/price/domain"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	// Provider should be used to create the cart Value objects
	Provider func() *Cart

	// Cart Value Object (immutable data - because the cartservice is responsible to return a cart).
	Cart struct {
		//ID is the main identifier of the cart
		ID string
		//EntityID is a second identifier that may be used by some backends
		EntityID string

		// ReservedOrderID is an ID already known by the Cart of the future order ID
		ReservedOrderID string

		//CartTotals - the cart totals (contain summary costs and discounts etc)
		CartTotals Totals

		//BillingAdress - the main billing address (relevant for all payments/invoices)
		BillingAdress Address

		//Purchaser - additional infos for the legal contact person in this order
		Purchaser Person

		//Deliveries - list of desired Deliverys (or Shippments) involved in this cart
		Deliveries []Delivery

		//AdditionalData   can be used for Custom attributes
		AdditionalData AdditionalData

		//BelongsToAuthenticatedUser - false = Guest Cart true = cart from the authenticated user
		BelongsToAuthenticatedUser bool
		AuthenticatedUserID        string

		AppliedCouponCodes []CouponCode

		DefaultCurrency string
	}

	// Teaser - represents some teaser infos for cart
	Teaser struct {
		ProductCount  int
		ItemCount     int
		DeliveryCodes []string
	}

	// CouponCode value object
	CouponCode struct {
		Code string
	}

	// Person value object
	Person struct {
		Address         *Address
		PersonalDetails PersonalDetails
		//ExistingCustomerData if the current purchaser is an existing customer - this contains infos about existing customer
		ExistingCustomerData *ExistingCustomerData
	}

	// ExistingCustomerData value object
	ExistingCustomerData struct {
		//ID of the customer
		ID string
	}

	// PersonalDetails value object
	PersonalDetails struct {
		DateOfBirth     string
		PassportCountry string
		PassportNumber  string
		Nationality     string
	}

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

	//AdditionalDeliverInfo is an interface that allows to store "any" additional objects on the cart
	// see DeliveryInfoUpdateCommand
	AdditionalDeliverInfo interface {
		Marshal() (json.RawMessage, error)
		Unmarshal(json.RawMessage) error
	}

	// DeliveryLocation value object
	DeliveryLocation struct {
		Type string
		//Address - only set for type adress
		Address *Address
		//Code - optional idendifier of this location/destination - is used in special destination Types
		Code string
	}

	// Totals value object
	Totals struct {
		//Additional non taxable totals
		Totalitems        []Totalitem
		TotalShippingItem ShippingItem
		//Final sum that need to be payed: GrandTotal = SubTotal + TaxAmount - DiscountAmount + SOME of Totalitems = (Sum of Items RowTotalWithDiscountInclTax) + SOME of Totalitems
		GrandTotal domain.Price
		//TaxAmount = Sum of Item TaxAmount
		TaxAmount domain.Price
	}

	// DeliveryTotals value object
	DeliveryTotals struct {
		//SubTotal = SUM of Item RowTotal
		SubTotal domain.Price
		//SubTotalInclTax = SUM of Item RowTotalInclTax
		SubTotalInclTax domain.Price
		//SubTotalWithDiscounts = SubTotal - Sum of Item ItemRelatedDiscountAmount
		SubTotalWithDiscounts domain.Price
		//SubTotalWithDiscountsAndTax= Sum of RowTotalWithItemRelatedDiscountInclTax
		SubTotalWithDiscountsAndTax domain.Price

		//TotalDiscountAmount = SUM of Item TotalDiscountAmount
		TotalDiscountAmount domain.Price
		//TotalNonItemRelatedDiscountAmount= SUM of Item NonItemRelatedDiscountAmount
		TotalNonItemRelatedDiscountAmount domain.Price
	}

	// Item for Cart
	Item struct {
		//ID of the item - need to be unique under a delivery
		ID string
		//
		UniqueID        string
		MarketplaceCode string
		//VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceID string

		Qty int

		AdditionalData map[string]string

		//brutto (gross) for single item
		SinglePrice domain.Price

		//AppliedDiscounts contains the details about the discounts applied to this item - they can be "itemrelated" or not
		AppliedDiscounts []ItemDiscount

		//TaxAmount = SinglePriceInclTax-SinglePrice - Tax is normaly applied after discounts
		TaxAmount domain.Price
	}

	// ItemCartReference - value object that can be used to reference a Item in a Cart
	//@todo - Use in ServicePort methods...
	ItemCartReference struct {
		ItemID       string
		DeliveryCode string
	}

	// ItemDiscount value object
	ItemDiscount struct {
		Code  string
		Title string
		Price domain.Price
		//IsItemRelated is a flag indicating if the discount should be displayed in the item or if it the result of a cart discount
		IsItemRelated bool
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price domain.Price
		Type  string
	}

	// ShippingItem value object
	ShippingItem struct {
		Title          string
		Price          domain.Price
		TaxAmount      domain.Price
		DiscountAmount domain.Price
	}

	// InvalidateCartEvent value object
	InvalidateCartEvent struct {
		Session *web.Session
	}

	// AdditionalData defines the supplementary cart data
	AdditionalData struct {
		CustomAttributes map[string]string
		SelectedPayment  SelectedPayment
	}

	// SelectedPayment value object
	SelectedPayment struct {
		Provider string
		Method   string
	}

	// PlacedOrderInfos represents a slice of PlacedOrderInfo
	PlacedOrderInfos []PlacedOrderInfo

	// PlacedOrderInfo defines the additional info struct for placed orders
	PlacedOrderInfo struct {
		OrderNumber  string
		DeliveryCode string
	}
)

var (
	// ErrAdditionalInfosNotFound is returned if the additional infos are not set
	ErrAdditionalInfosNotFound = errors.New("additional infos not found")
)

// Key constants
const (
	DeliveryWorkflowPickup      = "pickup"
	DeliveryWorkflowDelivery    = "delivery"
	DeliveryWorkflowUnspecified = "unspecified"

	DeliverylocationTypeUnspecified     = "unspecified"
	DeliverylocationTypeCollectionpoint = "collection-point"
	DeliverylocationTypeStore           = "store"
	DeliverylocationTypeAddress         = "address"
	DeliverylocationTypeFreightstation  = "freight-station"

	TotalsTypeDiscount      = "totals_type_discount"
	TotalsTypeVoucher       = "totals_type_voucher"
	TotalsTypeTax           = "totals_type_tax"
	TotalsTypeLoyaltypoints = "totals_loyaltypoints"
	TotalsTypeShipping      = "totals_type_shipping"
)

// GetMainShippingEMail returns the main shipping address email, empty string if not available
func (cart Cart) GetMainShippingEMail() string {
	for _, deliveries := range cart.Deliveries {
		if deliveries.DeliveryInfo.DeliveryLocation.Address != nil {
			if deliveries.DeliveryInfo.DeliveryLocation.Address.Email != "" {
				return deliveries.DeliveryInfo.DeliveryLocation.Address.Email
			}
		}
	}

	return ""
}

// GetDeliveryByCode gets a delivery by code
func (cart Cart) GetDeliveryByCode(deliveryCode string) (*Delivery, bool) {
	for _, delivery := range cart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			return &delivery, true
		}
	}

	return nil, false
}

// HasDeliveryForCode checks if a delivery with the given code exists in the cart
func (cart Cart) HasDeliveryForCode(deliveryCode string) bool {
	_, found := cart.GetDeliveryByCode(deliveryCode)

	return found == true
}

// GetDeliveryCodes returns a slice of all delivery codes in cart that have at least one cart item
func (cart Cart) GetDeliveryCodes() []string {
	var deliveryCodes []string

	for _, delivery := range cart.Deliveries {
		if len(delivery.Cartitems) > 0 {
			deliveryCodes = append(deliveryCodes, delivery.DeliveryInfo.Code)
		}
	}

	return deliveryCodes
}

// GetByItemID gets an item by its id
func (cart Cart) GetByItemID(itemID string, deliveryCode string) (*Item, error) {
	delivery, found := cart.GetDeliveryByCode(deliveryCode)
	if found != true {
		return nil, errors.Errorf("Delivery for code %v not found", deliveryCode)
	}
	for _, currentItem := range delivery.Cartitems {
		if currentItem.ID == itemID {
			return &currentItem, nil
		}
	}

	return nil, errors.Errorf("itemId %v in cart not existing", itemID)
}

func inStruct(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}

	return false
}

// ItemCount - returns amount of Cartitems
func (cart Cart) ItemCount() int {
	count := 0
	for _, delivery := range cart.Deliveries {
		for _, item := range delivery.Cartitems {
			count += item.Qty
		}
	}

	return count
}

// ProductCount - returns amount of different products
func (cart Cart) ProductCount() int {
	count := 0
	for _, delivery := range cart.Deliveries {
		count += len(delivery.Cartitems)
	}

	return count
}

// GetItemCartReferences returns a slice of all ItemCartReferences
func (cart Cart) GetItemCartReferences() []ItemCartReference {
	var ids []ItemCartReference
	for _, delivery := range cart.Deliveries {
		for _, item := range delivery.Cartitems {
			ids = append(ids, ItemCartReference{
				ItemID:       item.ID,
				DeliveryCode: delivery.DeliveryInfo.Code,
			})
		}
	}

	return ids
}

// GetVoucherSavings returns the savings of all vouchers
func (cart Cart) GetVoucherSavings() domain.Price {
	price := domain.Price{}
	for _, item := range cart.CartTotals.Totalitems {
		if item.Type == TotalsTypeVoucher {
			price, err := price.Add(item.Price)
			if err != nil {
				return price
			}
		}
	}
	if price.IsNegative() {
		return domain.Price{}
	}
	return price
}

// GetSavings retuns the total of all discount totals
func (cart Cart) GetSavings() domain.Price {
	price := domain.Price{}
	for _, item := range cart.CartTotals.Totalitems {
		if item.Type == TotalsTypeDiscount {
			price, err := price.Add(item.Price)
			if err != nil {
				return price
			}
		}
	}

	if price.IsNegative() {
		return domain.Price{}
	}
	return price
}

//SubTotal = SUM of Item RowTotal
func (cart Cart) SubTotal() domain.Price {
	return cart.deliverySum(func(d *Delivery) domain.Price {
		return d.DeliveryTotals().SubTotal
	})
}

//SubTotalInclTax = SUM of Item RowTotalInclTax
func (cart Cart) SubTotalInclTax() domain.Price {
	return cart.deliverySum(func(d *Delivery) domain.Price {
		return d.DeliveryTotals().SubTotalInclTax
	})
}

//deliverySum - private function that returns the sum of given deliveryPriceFunc
func (cart Cart) deliverySum(deliveryPriceFunc func(delivery *Delivery) domain.Price) domain.Price {
	if len(cart.Deliveries) == 0 {
		return domain.NewZero(cart.DefaultCurrency)
	}
	price := deliveryPriceFunc(&cart.Deliveries[0])
	for _, delivery := range cart.Deliveries[1:] {
		price, _ = price.Add(deliveryPriceFunc(&delivery))
	}
	return price
}

//SubTotalWithDiscounts = SubTotal - Sum of Item ItemRelatedDiscountAmount
func (cart Cart) SubTotalWithDiscounts() domain.Price {
	return cart.deliverySum(func(d *Delivery) domain.Price {
		return d.DeliveryTotals().SubTotalWithDiscounts
	})
}

//SubTotalWithDiscountsAndTax = Sum of RowTotalWithItemRelatedDiscountInclTax
func (cart Cart) SubTotalWithDiscountsAndTax() domain.Price {
	return cart.deliverySum(func(d *Delivery) domain.Price {
		return d.DeliveryTotals().SubTotalWithDiscountsAndTax
	})
}

//TotalDiscountAmount = SUM of Item TotalDiscountAmount
func (cart Cart) TotalDiscountAmount() domain.Price {
	return cart.deliverySum(func(d *Delivery) domain.Price {
		return d.DeliveryTotals().TotalDiscountAmount
	})
}

//TotalNonItemRelatedDiscountAmount = SUM of Item NonItemRelatedDiscountAmount
func (cart Cart) TotalNonItemRelatedDiscountAmount() domain.Price {
	return cart.deliverySum(func(d *Delivery) domain.Price {
		return d.DeliveryTotals().TotalNonItemRelatedDiscountAmount
	})
}

// HasAppliedCouponCode checks if a coupon code is applied to the cart
func (cart Cart) HasAppliedCouponCode() bool {
	return len(cart.AppliedCouponCodes) > 0
}

// GetCartTeaser returns the teaser
func (cart Cart) GetCartTeaser() *Teaser {
	return &Teaser{
		DeliveryCodes: cart.GetDeliveryCodes(),
		ItemCount:     cart.ItemCount(),
		ProductCount:  cart.ProductCount(),
	}
}

// GetTotalItemsByType gets a slice of all Totalitems by typeCode
func (ct Totals) GetTotalItemsByType(typeCode string) []Totalitem {
	var totalitems []Totalitem
	for _, item := range ct.Totalitems {
		if item.Type == typeCode {
			totalitems = append(totalitems, item)
		}
	}

	return totalitems
}

// GetSavingsByItem gets the savings by item
func (item Item) GetSavingsByItem() domain.Price {
	price, _ := item.ItemRelatedDiscountAmount().Add(item.NonItemRelatedDiscountAmount())
	return price.GetPayable()
}

// RowTotal = SinglePrice * Qty
func (item Item) RowTotal() domain.Price {
	return item.SinglePrice.Multiply(item.Qty).GetPayable()
}

//SinglePriceInclTax - netto (net) for single item. Normaly the Tax is applied after discounting - but that detail is up to the secondary adpater
func (item Item) SinglePriceInclTax() domain.Price {
	price, _ := item.SinglePrice.Add(item.TaxAmount)
	return price.GetPayable()
}

//TotalTaxAmount =Qty * (SinglePriceInclTax-SinglePrice)
func (item Item) TotalTaxAmount() domain.Price {
	price, _ := item.RowTotalInclTax().Sub(item.RowTotal())
	return price.GetPayable()
}

//RowTotalInclTax = RowTotal + TaxAmount
func (item Item) RowTotalInclTax() domain.Price {
	return item.SinglePriceInclTax().Multiply(item.Qty).GetPayable()
}

// TotalDiscountAmount = Sum of AppliedDiscounts = ItemRelatedDiscountAmount +NonItemRelatedDiscountAmount
func (item Item) TotalDiscountAmount() domain.Price {
	return item.GetSavingsByItem()
}

// ItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = True
func (item Item) ItemRelatedDiscountAmount() domain.Price {
	price := domain.Price{}
	for _, discount := range item.AppliedDiscounts {
		if !discount.IsItemRelated {
			continue
		}
		price, err := price.Add(discount.Price)
		if err != nil {
			return price
		}
	}

	if price.IsNegative() {
		return domain.Price{}
	}
	return price.GetPayable()
}

//NonItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = false
func (item Item) NonItemRelatedDiscountAmount() domain.Price {
	price := domain.Price{}
	for _, discount := range item.AppliedDiscounts {
		if discount.IsItemRelated {
			continue
		}
		price, err := price.Add(discount.Price)
		if err != nil {
			return price
		}
	}

	if price.IsNegative() {
		return domain.Price{}
	}
	return price.GetPayable()
}

//RowTotalWithItemRelatedDiscount =RowTotal-ItemRelatedDiscountAmount
func (item Item) RowTotalWithItemRelatedDiscount() domain.Price {
	price, _ := item.RowTotal().Sub(item.ItemRelatedDiscountAmount())
	return price.GetPayable()
}

//RowTotalWithItemRelatedDiscountInclTax =RowTotalInclTax-ItemRelatedDiscountAmount
func (item Item) RowTotalWithItemRelatedDiscountInclTax() domain.Price {
	price, _ := item.RowTotalInclTax().Sub(item.ItemRelatedDiscountAmount())
	return price.GetPayable()
}

//RowTotalWithDiscountInclTax This is the price the customer finaly need to pay for this item:  RowTotalWithDiscountInclTax = RowTotalInclTax-TotalDiscountAmount
func (item Item) RowTotalWithDiscountInclTax() domain.Price {
	price, _ := item.RowTotalInclTax().Sub(item.TotalDiscountAmount())
	return price.GetPayable()
}

//TotalWithDiscountInclTax - the price the customer need to pay for the shipping
func (s ShippingItem) TotalWithDiscountInclTax() domain.Price {
	price, _ := s.Price.Add(s.TaxAmount)
	price, _ = price.Sub(s.DiscountAmount)
	return price.GetPayable()
}

// GetOrderNumberForDeliveryCode returns the order number for a delivery code
func (poi PlacedOrderInfos) GetOrderNumberForDeliveryCode(deliveryCode string) string {
	for _, v := range poi {
		if v.DeliveryCode == deliveryCode {
			return v.OrderNumber
		}
	}
	return ""
}

//LoadAdditionalInfo - returns the additional Data
func (d *DeliveryInfo) LoadAdditionalInfo(key string, info AdditionalDeliverInfo) error {
	if d.AdditionalDeliveryInfos == nil {
		return ErrAdditionalInfosNotFound
	}
	if val, ok := d.AdditionalDeliveryInfos[key]; ok {
		return info.Unmarshal(val)
	}
	return ErrAdditionalInfosNotFound
}

//DeliveryTotals - Totals with the intent to use them to display the customer summary costs for this delivery
func (d Delivery) DeliveryTotals() DeliveryTotals {
	if len(d.Cartitems) == 0 {
		return DeliveryTotals{}
	}
	firstItem := d.Cartitems[0]
	deliveryTotals := DeliveryTotals{
		SubTotal:                          firstItem.RowTotal(),
		SubTotalInclTax:                   firstItem.RowTotalInclTax(),
		TotalDiscountAmount:               firstItem.TotalDiscountAmount(),
		SubTotalWithDiscounts:             firstItem.RowTotalWithItemRelatedDiscount(),
		SubTotalWithDiscountsAndTax:       firstItem.RowTotalWithItemRelatedDiscountInclTax(),
		TotalNonItemRelatedDiscountAmount: firstItem.NonItemRelatedDiscountAmount(),
	}

	for _, cartItem := range d.Cartitems[1:] {
		deliveryTotals.SubTotal, _ = deliveryTotals.SubTotal.Add(cartItem.RowTotal())
		deliveryTotals.SubTotalInclTax, _ = deliveryTotals.SubTotalInclTax.Add(cartItem.RowTotalInclTax())
		deliveryTotals.TotalDiscountAmount, _ = deliveryTotals.TotalDiscountAmount.Add(cartItem.TotalDiscountAmount())
		deliveryTotals.TotalNonItemRelatedDiscountAmount, _ = deliveryTotals.TotalNonItemRelatedDiscountAmount.Add(cartItem.NonItemRelatedDiscountAmount())
		deliveryTotals.SubTotalWithDiscounts, _ = deliveryTotals.SubTotalWithDiscounts.Add(cartItem.RowTotalWithItemRelatedDiscount())
		deliveryTotals.SubTotalWithDiscountsAndTax, _ = deliveryTotals.SubTotalWithDiscountsAndTax.Add(cartItem.RowTotalWithItemRelatedDiscountInclTax())
	}
	return deliveryTotals
}
