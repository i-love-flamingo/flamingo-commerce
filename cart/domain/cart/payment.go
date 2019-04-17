package cart

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

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
		// ChargeAssignments - optional the assignment of this transaction to charges - this might be required for payments that are really only done
		ChargeAssignments *ChargeAssignments
	}

	// ChargeAssignments collects all assignments for item, shipment and total charges
	ChargeAssignments struct {
		ItemChargeAssignments     []ItemChargeAssignment
		TotalChargeAssignments    []TotalChargeAssignment
		ShipmentChargeAssignments []ShipmentChargeAssignment
	}

	// ItemChargeAssignment holds the information what amount was assigned to a specific chargetype of a specific item in the cart
	ItemChargeAssignment struct {
		ItemID string
		Charge domain.Charge
	}

	// TotalChargeAssignment holds the information what amount was assigned to a specific chargetype of a specific total in the cart
	TotalChargeAssignment struct {
		Type   string
		Charge domain.Charge
	}

	// ShipmentChargeAssignment holds the information what amount was assigned to a specific chargetype of a specific shipment in the cart
	ShipmentChargeAssignment struct {
		DeliveryCode string
		Charge       domain.Charge
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
		//The Charge that need to be payed (with Amount, Value and Type in it)
		Charge domain.Charge

		//The selected payment method (code) that should be used
		Method string
    
		// ChargeAssignments - optional the assignment of this transaction to charges - this might be required for payments that are really only done
		ChargeAssignments *ChargeAssignments
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
			{
				Charge: domain.Charge{
					Type:  domain.ChargeTypeMain,
					Price: grandTotal,
					Value: grandTotal,
				},
				Method: method,
			},
		},
	}
}

//IsSelected - returns true if a Gateway  is selected
func (s PaymentSelection) IsSelected() bool {
	return s.Gateway != ""
}

//GetCharges - sum per chargetype
func (s PaymentSelection) GetCharges() domain.Charges {
	result := domain.Charges{}
	for _, cs := range s.ChargeSplits {
		result = result.AddCharge(cs.Charge)
	}
	return result
}
