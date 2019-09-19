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
		Data interface{}
		// Error contains additional information in case of an error (e.g. payment failed)
		Error *Error
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
)

// Error getter
func (pe *Error) Error() string {
	return pe.ErrorMessage
}
