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

	// FlowResult contains an url and a type to use to start a flow
	FlowResult struct {
		URL  *url.URL
		Type string
		Data interface{}
	}

	// Error should be used by PaymentGateway to indicate that payment failed (so that the customer can see a speaking message)
	Error struct {
		ErrorMessage string
		ErrorCode    string
	}
)

const (
	// PaymentErrorCodeCancelled cancelled
	PaymentErrorCodeCancelled = "payment_cancelled"
	// PaymentErrorCodeAuthorizeFailed error
	PaymentErrorCodeAuthorizeFailed = "authorization_failed"
	// PaymentErrorCodeCaptureFailed error
	PaymentErrorCodeCaptureFailed = "capture_failed"
	// ErrorPaymentAbortedByCustomer error
	ErrorPaymentAbortedByCustomer = "payment-aborted-by-customer"
)

// Error getter
func (pe *Error) Error() string {
	return pe.ErrorMessage
}
