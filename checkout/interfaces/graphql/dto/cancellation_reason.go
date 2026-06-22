package dto

import (
	paymentdomain "flamingo.me/flamingo-commerce/v3/payment/domain"
)

// PaymentCancellationReason is the GraphQL-facing string enum for the optional
// cancel-mutation reason argument. Values mirror the proto enum.
type PaymentCancellationReason string

const (
	// PaymentCancellationReasonUnspecified is the default / legacy value.
	PaymentCancellationReasonUnspecified PaymentCancellationReason = "UNSPECIFIED"
	// PaymentCancellationReasonAbortedByCustomer marks an explicit customer abort.
	PaymentCancellationReasonAbortedByCustomer PaymentCancellationReason = "ABORTED_BY_CUSTOMER"
)

// MapPaymentCancellationReason maps the optional GraphQL reason to the domain
// enum; nil / UNSPECIFIED / unknown all resolve to Unspecified.
func MapPaymentCancellationReason(r *PaymentCancellationReason) paymentdomain.CancellationReason {
	if r == nil {
		return paymentdomain.CancellationReasonUnspecified
	}

	switch *r {
	case PaymentCancellationReasonUnspecified:
		return paymentdomain.CancellationReasonUnspecified
	case PaymentCancellationReasonAbortedByCustomer:
		return paymentdomain.CancellationReasonAbortedByCustomer
	default:
		return paymentdomain.CancellationReasonUnspecified
	}
}
