package placeorder

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	price "flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// Service  interface - Secondary PORT
	Service interface {
		PlaceGuestCart(ctx context.Context, cart *cart.Cart, payment *Payment) (PlacedOrderInfos, error)
		PlaceCustomerCart(ctx context.Context, identity auth.Identity, cart *cart.Cart, payment *Payment) (PlacedOrderInfos, error)
		ReserveOrderID(ctx context.Context, cart *cart.Cart) (string, error)
		// CancelGuestOrder cancels a previously placed guest order and returns the used cart
		CancelGuestOrder(ctx context.Context, orderInfos PlacedOrderInfos) error
		// CancelCustomerOrder cancels a previously placed guest order and returns the used cart
		CancelCustomerOrder(ctx context.Context, orderInfos PlacedOrderInfos, identity auth.Identity) error
	}
	// Payment represents all payments done for the cart and which items have been purchased by what method
	Payment struct {
		// The name of the Gateway that has returned the Payment for the cart
		Gateway string
		// Transactions is the list of individual transactions -  most cases only one Transaction might be part of the payment
		Transactions []Transaction
		// RawTransactionData can be used to store any additional stuff (specific for Gateway)
		RawTransactionData interface{}
		// PaymentID is a optional reference of the Payment (that contains the Transactions)
		PaymentID string
	}

	// Transaction representing the transaction
	Transaction struct {
		// PaymentProvider - optional - the underling processor of this transaction (e.g. "paymark")
		PaymentProvider string
		// Method like "paymark_cc" , "paypal",
		Method string
		// Status - Method specific status e.g. Auth, Captured, Open, ...
		Status string
		// TransactionID - The main reference of the payment that was done
		TransactionID string
		// AdditionalData - room for AdditionalData - specific to the payment
		AdditionalData map[string]string
		// AmountPayed the amount that have been paid - eventually in a different currency
		AmountPayed price.Price
		// ValuedPayed the value of the AmountPayed in the cart default currency
		ValuedAmountPayed price.Price
		// CreditCardInfo Optional
		CreditCardInfo *CreditCardInfo
		// Title - speaking title - optional may describe the payment and may be shown to the customer
		Title string
		// RawTransactionData - place to store any additional stuff (specific for Gateway)
		RawTransactionData interface{}
		// ChargeAssignments - optional the assignment of this transaction to charges - this might be required for payments that are really only done for a certain item
		ChargeByItem *ChargeByItem
	}

	// ChargeByItem - the Charge that is paid for the individual items
	ChargeByItem struct {
		cartItems     map[string]price.Charge
		shippingItems map[string]price.Charge
		totalItems    map[string]price.Charge
	}

	// CreditCardInfo contains the necessary data
	CreditCardInfo struct {
		AnonymizedCardNumber string
		Type                 string
		CardHolder           string
		Expire               string
	}

	// PlacedOrderInfos represents a slice of PlacedOrderInfo
	PlacedOrderInfos []PlacedOrderInfo

	// PlacedOrderInfo defines the additional info struct for placed orders
	PlacedOrderInfo struct {
		OrderNumber  string
		DeliveryCode string
	}
)

const (
	// PaymentStatusCaptured a payment which has been captured
	PaymentStatusCaptured = "CAPTURED"
	// PaymentStatusAuthorized a payment which has been AUTHORIZED
	PaymentStatusAuthorized = "AUTHORIZED"
	// PaymentStatusOpen payment is still open
	PaymentStatusOpen = "OPEN"
)

// AddTransaction for a paymentInfo with items
func (cp *Payment) AddTransaction(transaction Transaction) {
	cp.Transactions = append(cp.Transactions, transaction)
}

// TotalValue returns the Total Valued Price
func (cp *Payment) TotalValue() (price.Price, error) {
	var prices []price.Price

	for _, transaction := range cp.Transactions {
		prices = append(prices, transaction.ValuedAmountPayed)
	}

	return price.SumAll(prices...)
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

// CartItems return CartItems
func (c ChargeByItem) CartItems() map[string]price.Charge {
	return c.cartItems
}

// ChargeForCartItem returns Charge for a cart item by id
func (c ChargeByItem) ChargeForCartItem(itemid string) (*price.Charge, bool) {
	if charge, ok := c.cartItems[itemid]; ok {
		return &charge, true
	}
	return nil, false
}

// ChargeForDeliveryCode returns Charge for a shipping item by delivery code
func (c ChargeByItem) ChargeForDeliveryCode(itemid string) (*price.Charge, bool) {
	if charge, ok := c.shippingItems[itemid]; ok {
		return &charge, true
	}
	return nil, false
}

// ChargeForTotal returns Charge for a total item by code
func (c ChargeByItem) ChargeForTotal(itemid string) (*price.Charge, bool) {
	if charge, ok := c.totalItems[itemid]; ok {
		return &charge, true
	}
	return nil, false
}

// TotalItems returns totalItems
func (c ChargeByItem) TotalItems() map[string]price.Charge {
	return c.totalItems
}

// ShippingItems returns ShippingItems
func (c ChargeByItem) ShippingItems() map[string]price.Charge {
	return c.shippingItems
}

// AddCartItem modifies the current instance and adds a charge for a cart item
func (c ChargeByItem) AddCartItem(id string, charge price.Charge) ChargeByItem {
	if c.cartItems == nil {
		c.cartItems = make(map[string]price.Charge)
	}
	c.cartItems[id] = charge
	return c
}

// AddTotalItem modifies the current instance and adds a charge for a total item
func (c ChargeByItem) AddTotalItem(id string, charge price.Charge) ChargeByItem {
	if c.totalItems == nil {
		c.totalItems = make(map[string]price.Charge)
	}
	c.totalItems[id] = charge
	return c
}

// AddShippingItems modifies the current instance and adds a charge for a shipping item
func (c ChargeByItem) AddShippingItems(id string, charge price.Charge) ChargeByItem {
	if c.shippingItems == nil {
		c.shippingItems = make(map[string]price.Charge)
	}
	c.shippingItems[id] = charge
	return c
}
