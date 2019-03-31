package domain

type (
	// Method contains information about a general payment method
	Method struct {
		//A speaking title
		Title string
		//A unique Code
		Code string
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
)

// Error getter
func (pe *Error) Error() string {
	return pe.ErrorMessage
}
