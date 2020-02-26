package domain

import "net/url"

type (
	// Method contains information about a general payment method
	Method struct {
		//A speaking title
		Title string
		//A unique Code
		Code string
	}

	// FlowResult contains information about a newly started flow
	FlowResult struct {
		// EarlyPlaceOrder indicates if the order should be placed with the beginning of the payment flow
		EarlyPlaceOrder bool
		// Status contains the current payment status
		Status FlowStatus
	}

	// FlowStatus contains information about the current payment status
	FlowStatus struct {
		// Status of the payment flow
		Status string
		// Action to perform to proceed in the payment flow
		Action string
		// Data contains additional information related to the action / flow
		Data       interface{}
		ActionData FlowActionData
		// Error contains additional information in case of an error (e.g. payment failed)
		Error *Error
	}

	// FlowActionData contains additional data for the current action
	FlowActionData struct {
		// URL is used to pass URL data to the user if the current state needs some
		URL *url.URL
		// DisplayData holds data, normally HTML to be displayed to the user
		DisplayData   string
		FormParameter map[string]FormField
	}

	// FormField contains form fields
	FormField struct {
		Value []string
	}

	// Error should be used by PaymentGateway to indicate that payment failed (so that the customer can see a speaking message)
	Error struct {
		ErrorMessage string
		ErrorCode    string
	}
)

const (
	// PaymentErrorCodeFailed error
	PaymentErrorCodeFailed = "failed"
	// PaymentErrorCodeAuthorizeFailed error
	PaymentErrorCodeAuthorizeFailed = "authorization_failed"
	// PaymentErrorCodeCaptureFailed error
	PaymentErrorCodeCaptureFailed = "capture_failed"
	// PaymentErrorAbortedByCustomer error
	PaymentErrorAbortedByCustomer = "aborted_by_customer"
	// PaymentErrorCodeCancelled error
	PaymentErrorCodeCancelled = "cancelled"

	// PaymentFlowStatusUnapproved payment started
	PaymentFlowStatusUnapproved = "payment_unapproved"
	// PaymentFlowStatusFailed payment failed
	PaymentFlowStatusFailed = "payment_failed"
	// PaymentFlowStatusAborted payment aborted by user
	PaymentFlowStatusAborted = "payment_aborted"
	// PaymentFlowStatusApproved payment approved by payment adapter
	PaymentFlowStatusApproved = "payment_approved"
	// PaymentFlowStatusCompleted payment approved and confirmed by customer
	PaymentFlowStatusCompleted = "payment_completed"
	// PaymentFlowWaitingForCustomer payment waiting for customer
	PaymentFlowWaitingForCustomer = "payment_waiting_for_customer"
	// PaymentFlowStatusCancelled payment cancelled by provider
	PaymentFlowStatusCancelled = "payment_cancelled"

	// PaymentFlowActionShowIframe signals the frontend to show an iframe
	PaymentFlowActionShowIframe = "show_iframe"
	// PaymentFlowActionShowHTML signals the frontend to show HTML
	PaymentFlowActionShowHTML = "show_html"
	// PaymentFlowActionRedirect signals the frontend to do a redirect to a hosted payment page
	PaymentFlowActionRedirect = "redirect"
	// PaymentFlowActionPostRedirect signals the frontend to do a post redirect to a hosted payment page
	PaymentFlowActionPostRedirect = "post_redirect"
)

// Error getter
func (pe *Error) Error() string {
	return pe.ErrorMessage
}
