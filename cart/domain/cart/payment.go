package cart

import "flamingo.me/flamingo-commerce/v3/price/domain"

type (
	// Payment represents all payments done for the cart and which items have been purchased by what method
	Payment struct {
		//The name of the Gateway that has returned the Payment for the cart
		Gateway string
		//Transactions - the list of individual transactions -  most cases only one Transaction might be part of the payment
		Transactions []Transaction
		//RawTransactionData - place to store any additional stuff (specific for Gateway)
		RawTransactionData interface{}
	}

	// Transaction representing the transaction
	Transaction struct {
		//PaymentProvider - optional - the underling processor of this transaction (e.g. "paymark")
		PaymentProvider string
		//Method like "paymark_cc" , "paypal",
		Method string
		//Status - Method specific status e.g. Auth, Captured, Open, ...
		Status string
		//TransactionID - The main reference of the payment that was done
		TransactionID string
		//AdditionalData - room for AdditionalData - specific to the payment
		AdditionalData map[string]string
		//AmountPayed the amount that have been payed - eventually in a different currency
		AmountPayed domain.Price
		//ValuedPayed the value of the AmountPayed in the cart default currency
		ValuedAmountPayed domain.Price
		//CreditCardInfo Optional
		CreditCardInfo *CreditCardInfo
		//Title - speaking title - optional may describe the payment and may be shown to the customer
		Title string
		//RawTransactionData - place to store any additional stuff (specific for Gateway)
		RawTransactionData interface{}
		//ItemChargeAssignments - optional the assignment of this transaction to item charges- this might be required for payments that are really only done for an item
		// List of items.UniqueID
		ItemChargeAssignments []ItemChargeAssignment
	}

	//ItemChargeAssignment holds the information what amount was assigned to a specific chargetype of a specific item in the cart
	ItemChargeAssignment struct {
		UniqueItemID      string
		ChargeType        string
		AmountPayed       domain.Price
		ValuedAmountPayed domain.Price
	}

	// CreditCardInfo contains the necessary data
	CreditCardInfo struct {
		AnonymizedCardNumber string
		Type                 string
		CardHolder           string
		Expire               string
	}

	// PaymentSelection value object - that represents the payment selection on the cart
	PaymentSelection struct {
		//Gateway - the selected Gateway
		Gateway string
		//ChargeSplits - the selected split
		ChargeSplits []ChargeSplit
	}

	//ChargeSplit - amount by Method and chargetype
	ChargeSplit struct {
		//The type of the charge that is supposed to be payed
		ChargeType string
		//The selected payment method (code) that should be used
		Method string
		//The amount that is supposed to be payed
		Amount domain.Price
		//The value in cart currency of the amount
		ValuedAmount domain.Price
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

//NewSimplePaymentSelection - returns a PaymentSelection that can be used to update the cart.
// 	multiple charges to pay the cart are not used here: The complete grandtotal is selected to be payed in one charge with the given paymentgateway and paymentmethod
func NewSimplePaymentSelection(gateway string, method string, grandTotal domain.Price) PaymentSelection {
	return PaymentSelection{
		Gateway: gateway,
		ChargeSplits: []ChargeSplit{
			ChargeSplit{
				Amount:     grandTotal,
				ChargeType: domain.ChargeTypeMain,
				Method:     method,
			},
		},
	}
}
