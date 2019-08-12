package domain

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
		// Status of the payment flow
		Status string
		// EarlyPlaceOrder indicates if the order should be placed with the beginning of the payment flow
		EarlyPlaceOrder bool
		// Action to perform to proceed in the payment flow
		Action string
		// Data contains additional information related to the action / flow
		Data interface{}
	}

	// FlowStatus contains information about the current payment status
	FlowStatus struct {
		// Status of the payment flow
		Status string
		// Action to perform to proceed in the payment flow
		Action string
		// Data contains additional information related to the action / flow
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
)

// Error getter
func (pe *Error) Error() string {
	return pe.ErrorMessage
}
