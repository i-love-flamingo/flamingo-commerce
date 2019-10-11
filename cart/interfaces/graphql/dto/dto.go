package dto

import (
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	formDomain "flamingo.me/form/domain"
)

type (
	//BillingAddressForm - dto for graphql
	BillingAddressForm struct {
		FormData       forms.BillingAddressForm
		Processed      bool
		ValidationInfo ValidationInfo
	}

	//ValidationInfo - dto for graphql
	ValidationInfo struct {
		GeneralErrors []formDomain.Error
		FieldErrors   []FieldError
	}

	//FieldError - dto for graphql
	FieldError struct {
		// MessageKey - a key of the error message. Often used to pass to translation func in the template
		MessageKey string
		// DefaultLabel - a speaking error label. OFten used to show to end user - in case no translation exists
		DefaultLabel string
		//FieldName
		FieldName string
	}

	//SelectedPaymentResult represents the selected payment
	SelectedPaymentResult struct {
		//Processed
		Processed bool
		//ValidationInfo
		ValidationInfo ValidationInfo
	}
)
