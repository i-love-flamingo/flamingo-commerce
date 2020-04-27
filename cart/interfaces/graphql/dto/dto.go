package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	formDomain "flamingo.me/form/domain"
)

type (
	// BillingAddressForm is the GraphQL representation of the billing form
	BillingAddressForm struct {
		FormData       forms.AddressForm
		Processed      bool
		ValidationInfo ValidationInfo
	}

	// DeliveryAddressForm is the GraphQL representation of the delivery form
	DeliveryAddressForm struct {
		FormData          forms.AddressForm
		Processed         bool
		ValidationInfo    ValidationInfo
		UseBillingAddress bool
		DeliveryCode      string
		Method            string
		Carrier           string
	}

	// DeliveryShippingOption used to update shipping method/carrier for a specific delivery
	DeliveryShippingOption struct {
		DeliveryCode string
		Method       string
		Carrier      string
	}

	// ValidationInfo contains form related validation information
	ValidationInfo struct {
		GeneralErrors []formDomain.Error
		FieldErrors   []FieldError
	}

	// FieldError contains field related errors
	FieldError struct {
		// MessageKey - a key of the error message. Often used to pass to translation func in the template
		MessageKey string
		// DefaultLabel - a speaking error label. OFten used to show to end user - in case no translation exists
		DefaultLabel string
		//FieldName
		FieldName string
	}

	// SelectedPaymentResult represents the selected payment
	SelectedPaymentResult struct {
		//Processed
		Processed bool
		//ValidationInfo
		ValidationInfo ValidationInfo
	}
)
