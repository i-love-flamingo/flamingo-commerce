package cart

import (
	"fmt"
	"math/big"

	"flamingo.me/flamingo-commerce/v3/price/domain"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

type (
	// Provider should be used to create the cart Value objects
	provider func() *Cart

	// Cart Value Object (immutable data - because the cartservice is responsible to return a cart).
	Cart struct {
		//ID is the main identifier of the cart
		ID string
		//EntityID is a second identifier that may be used by some backends
		EntityID string

		//BillingAdress - the main billing address (relevant for all payments/invoices)
		BillingAdress *Address

		//Purchaser - additional infos for the legal contact person in this order
		Purchaser *Person

		//Deliveries - list of desired Deliverys (or Shippments) involved in this cart
		Deliveries []Delivery

		//AdditionalData   can be used for Custom attributes
		AdditionalData AdditionalData

		//PaymentSelection - the saved PaymentSelection (saves "how" the customer want to pay)
		PaymentSelection *PaymentSelection

		//BelongsToAuthenticatedUser - false = Guest Cart true = cart from the authenticated user
		BelongsToAuthenticatedUser bool
		AuthenticatedUserID        string

		AppliedCouponCodes []CouponCode

		DefaultCurrency string

		//Additional non taxable totals
		Totalitems []Totalitem
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

	//Taxes is a list of Tax
	Taxes []Tax

	//Tax - it the Tax represented by an Amount and Optional the Rate.
	Tax struct {
		Amount domain.Price
		Type   string
		Rate   *big.Float
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price domain.Price
		Type  string
	}

	// InvalidateCartEvent value object
	InvalidateCartEvent struct {
		Session *web.Session
	}

	// AdditionalData defines the supplementary cart data
	AdditionalData struct {
		//CustomAttributes list of key values
		CustomAttributes map[string]string
		// ReservedOrderID is an ID already known by the Cart of the future order ID
		ReservedOrderID string
	}

	// PlacedOrderInfos represents a slice of PlacedOrderInfo
	PlacedOrderInfos []PlacedOrderInfo

	// PlacedOrderInfo defines the additional info struct for placed orders
	PlacedOrderInfo struct {
		OrderNumber  string
		DeliveryCode string
	}

	//Builder - the main builder for a cart
	Builder struct {
		cartInBuilding *Cart
	}
	// BuilderProvider should be used to create the cart by using the Builder
	BuilderProvider func() *Builder
)

var (
	// ErrAdditionalInfosNotFound is returned if the additional infos are not set
	ErrAdditionalInfosNotFound = errors.New("additional infos not found")
)

// Key constants
const (
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

// IsEmpty - returns true if cart is empty
func (cart Cart) IsEmpty() bool {
	return cart.GetCartTeaser().ItemCount == 0
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
func (cart Cart) GetByItemID(itemID string) (*Item, error) {
	for _, delivery := range cart.Deliveries {
		for _, currentItem := range delivery.Cartitems {
			if currentItem.ID == itemID {
				return &currentItem, nil
			}
		}
	}

	return nil, errors.Errorf("itemId %q in cart does not exist", itemID)
}

// GetByExternalReference gets an item by its external reference
func (cart Cart) GetByExternalReference(ref string) (*Item, error) {
	for _, delivery := range cart.Deliveries {
		for _, currentItem := range delivery.Cartitems {
			if currentItem.ExternalReference == ref {
				return &currentItem, nil
			}
		}
	}

	return nil, errors.Errorf("uitemID %v in cart not existing", ref)
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

// IsPaymentSelected - returns true if a valid payment is selected
func (cart Cart) IsPaymentSelected() bool {
	return cart.PaymentSelection != nil && cart.PaymentSelection.IsSelected()
}

// GetVoucherSavings returns the savings of all vouchers
func (cart Cart) GetVoucherSavings() domain.Price {
	price := domain.Price{}
	var err error
	for _, item := range cart.Totalitems {
		if item.Type == TotalsTypeVoucher {
			price, err = price.Add(item.Price)
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

// GrandTotal - Final sum (Valued price) that need to be payed: GrandTotal = SubTotal + TaxAmount - DiscountAmount + SOME of Totalitems = (Sum of Items RowTotalWithDiscountInclTax) + SOME of Totalitems
func (cart Cart) GrandTotal() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.GrandTotal())
	}
	for _, total := range cart.Totalitems {
		prices = append(prices, total.Price)
	}
	price, _ := domain.SumAll(prices...)
	return price
}

// SumShipping - returns sum price of deliveries ShippingItems
func (cart Cart) SumShippingNet() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.ShippingItem.PriceNet)
	}
	price, _ := domain.SumAll(prices...)
	return price
}

// HasShippingCosts - returns true if cart HasShippingCosts
func (cart Cart) HasShippingCosts() bool {
	return cart.SumShippingNet().IsPositive()
}

// AllShippingTitles - returns all ShippingItem titles
func (cart Cart) AllShippingTitles() []string {
	var label []string
	for _, del := range cart.Deliveries {
		label = append(label, del.ShippingItem.Title)
	}
	return label
}

// SubTotalGross - returns sum price of deliveries SubTotalGross
func (cart Cart) SubTotalGross() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.SubTotalGross())
	}
	price, _ := domain.SumAll(prices...)
	return price
}

// SumTaxes - returns sum of deliveries SumRowTaxes
func (cart Cart) SumTaxes() Taxes {
	newTaxes := Taxes{}
	for _, del := range cart.Deliveries {
		newTaxes = newTaxes.AddTaxesWithMerge(del.SumRowTaxes())
	}
	return newTaxes
}

// SumTotalTaxAmount - returns sum price of deliveries Taxes as total amount
func (cart Cart) SumTotalTaxAmount() domain.Price {
	return cart.SumTaxes().TotalAmount()
}

// SubTotalNet - returns sum price of deliveries SubTotalNet
func (cart Cart) SubTotalNet() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.SubTotalNet())
	}
	price, _ := domain.SumAll(prices...)
	return price
}

// SubTotalGrossWithDiscounts - returns sum price of deliveries SubTotalGrossWithDiscounts
func (cart Cart) SubTotalGrossWithDiscounts() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.SubTotalGrossWithDiscounts())
	}
	price, _ := domain.SumAll(prices...)
	return price
}

// SubTotalNetWithDiscounts - returns sum price of deliveries SubTotalNetWithDiscounts
func (cart Cart) SubTotalNetWithDiscounts() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.SubTotalNetWithDiscounts())
	}
	price, _ := domain.SumAll(prices...)
	return price
}

// SumTotalDiscountAmount - returns sum price of deliveries SumTotalDiscountAmount
func (cart Cart) SumTotalDiscountAmount() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.SumTotalDiscountAmount())
	}
	price, _ := domain.SumAll(prices...)
	return price
}

// SumNonItemRelatedDiscountAmount - returns sum price of deliveries SumNonItemRelatedDiscountAmount
func (cart Cart) SumNonItemRelatedDiscountAmount() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.SumNonItemRelatedDiscountAmount())
	}
	price, _ := domain.SumAll(prices...)
	return price
}

// SumItemRelatedDiscountAmount - returns sum price of deliveries SumItemRelatedDiscountAmount
func (cart Cart) SumItemRelatedDiscountAmount() domain.Price {
	var prices []domain.Price
	for _, del := range cart.Deliveries {
		prices = append(prices, del.SumItemRelatedDiscountAmount())
	}
	price, _ := domain.SumAll(prices...)
	return price
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

// GetPaymentReference - Returns a string that can be used as reference to pass to payment gateway. You may want to use it. It returns either the reserved Order id or the cart id/entityid
func (cart Cart) GetPaymentReference() string {
	if cart.AdditionalData.ReservedOrderID != "" {
		return cart.AdditionalData.ReservedOrderID
	}
	return fmt.Sprintf("%v-%v",cart.ID,cart.EntityID)
}

// GetTotalItemsByType gets a slice of all Totalitems by typeCode
func (cart Cart) GetTotalItemsByType(typeCode string) []Totalitem {
	var totalitems []Totalitem
	for _, item := range cart.Totalitems {
		if item.Type == typeCode {
			totalitems = append(totalitems, item)
		}
	}
	return totalitems
}

// GrandTotalCharges - Final sum that need to be payed - splitted by the charges that need to be payed
func (cart Cart) GrandTotalCharges() domain.Charges {
	// Check if a specific split was saved:
	if cart.PaymentSelection != nil {
		return cart.PaymentSelection.GetCharges()
	}
	// else return the grandtotal as main charge
	charges := domain.Charges{}
	mainCharge := domain.Charge{
		Value: cart.GrandTotal(),
		Price: cart.GrandTotal(),
		Type:  domain.ChargeTypeMain,
	}
	charges = charges.AddCharge(mainCharge)
	return charges
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

// AddTax returns new Tax with this Tax added
func (t Taxes) AddTax(tax Tax) Taxes {
	newTaxes := Taxes(t)
	newTaxes = append(newTaxes, tax)
	return newTaxes
}

// AddTaxWithMerge returns new Taxes with this Tax added
func (t Taxes) AddTaxWithMerge(taxToAddOrMerge Tax) Taxes {
	newTaxes := Taxes(t)
	for k, tax := range newTaxes {
		if tax.Type == taxToAddOrMerge.Type {
			if tax.Rate == nil && taxToAddOrMerge.Rate == nil {
				newTaxes[k].Amount, _ = tax.Amount.Add(taxToAddOrMerge.Amount)
				return newTaxes
			} else if tax.Rate != nil && taxToAddOrMerge.Rate != nil && (tax.Rate.Cmp(taxToAddOrMerge.Rate) == 0) {
				newTaxes[k].Amount, _ = tax.Amount.Add(taxToAddOrMerge.Amount)
				return newTaxes
			}
		}
	}
	newTaxes = newTaxes.AddTax(taxToAddOrMerge)

	return newTaxes
}

// AddTaxesWithMerge - returns new Taxes with the given Taxes all added or merged in
func (t Taxes) AddTaxesWithMerge(taxes Taxes) Taxes {
	newTaxes := Taxes(t)
	for _, tax := range taxes {
		newTaxes = newTaxes.AddTaxWithMerge(tax)
	}
	return newTaxes
}

// TotalAmount - returns the sum of all taxes as price
func (t Taxes) TotalAmount() domain.Price {
	var prices []domain.Price
	for _, tax := range t {
		prices = append(prices, tax.Amount)
	}

	result, _ := domain.SumAll(prices...)
	return result
}

// ###################

// Build - main factory method
func (b *Builder) Build() (*Cart, error) {
	return b.reset(nil)
}

// SetIds - sets the cart ids
func (b *Builder) SetIds(id string, entityID string) *Builder {
	b.init()
	b.cartInBuilding.ID = id
	b.cartInBuilding.EntityID = entityID
	return b
}

// SetReservedOrderID - optional
func (b *Builder) SetReservedOrderID(id string) *Builder {
	b.init()
	b.cartInBuilding.AdditionalData.ReservedOrderID = id
	return b
}

// SetBillingAdress - optional
func (b *Builder) SetBillingAdress(a Address) *Builder {
	b.init()
	b.cartInBuilding.BillingAdress = &a
	return b
}

// SetPurchaser - optional
func (b *Builder) SetPurchaser(p Person) *Builder {
	b.init()
	b.cartInBuilding.Purchaser = &p
	return b
}

// AddDelivery - add a delivery subobject - use the DeliveryBuilder
func (b *Builder) AddDelivery(d Delivery) *Builder {
	b.init()
	b.cartInBuilding.Deliveries = append(b.cartInBuilding.Deliveries, d)
	return b
}

// SetAdditionalData - to add additional data
func (b *Builder) SetAdditionalData(d AdditionalData) *Builder {
	b.init()
	b.cartInBuilding.AdditionalData = d
	return b
}

// SetPaymentSelection - to add additional data
func (b *Builder) SetPaymentSelection(d PaymentSelection) *Builder {
	b.init()
	b.cartInBuilding.PaymentSelection = &d
	return b
}

// SetAuthenticatedUserID - to mark the art as authenticated users cart
func (b *Builder) SetAuthenticatedUserID(id string) *Builder {
	b.init()
	b.cartInBuilding.AuthenticatedUserID = id
	b.cartInBuilding.BelongsToAuthenticatedUser = true
	return b
}

// SetBelongsToAuthenticatedUser -  mark the art as authenticated users cart
func (b *Builder) SetBelongsToAuthenticatedUser(v bool) *Builder {
	b.init()
	b.cartInBuilding.BelongsToAuthenticatedUser = v
	return b
}

// AddAppliedCouponCode - optional - add the coupon that is applied for the  cart
func (b *Builder) AddAppliedCouponCode(code CouponCode) *Builder {
	b.init()
	b.cartInBuilding.AppliedCouponCodes = append(b.cartInBuilding.AppliedCouponCodes, code)
	return b
}

// SetDefaultCurrency - sets the default currency
func (b *Builder) SetDefaultCurrency(d string) *Builder {
	b.init()
	b.cartInBuilding.DefaultCurrency = d
	return b
}

// AddTotalitem - adds nontaxable extra totals on cartlevel
func (b *Builder) AddTotalitem(totali Totalitem) *Builder {
	b.init()
	b.cartInBuilding.Totalitems = append(b.cartInBuilding.Totalitems, totali)
	return b
}

func (b *Builder) init() {
	if b.cartInBuilding == nil {
		b.cartInBuilding = &Cart{}
	}
}

func (b *Builder) reset(err error) (*Cart, error) {
	cart := b.cartInBuilding
	b.cartInBuilding = nil
	return cart, err
}
