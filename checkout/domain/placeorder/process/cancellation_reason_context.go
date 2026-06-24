package process

import (
	"context"

	paymentdomain "flamingo.me/flamingo-commerce/v3/payment/domain"
)

type cancellationReasonContextKey struct{}

// ContextWithCancellationReason returns a context carrying the cancellation
// reason for the single rollback -> state.Rollback hop.
func ContextWithCancellationReason(ctx context.Context, reason paymentdomain.CancellationReason) context.Context {
	return context.WithValue(ctx, cancellationReasonContextKey{}, reason)
}

// CancellationReasonFromContext reads the reason placed by
// ContextWithCancellationReason; missing / wrong type defaults to Unspecified.
func CancellationReasonFromContext(ctx context.Context) paymentdomain.CancellationReason {
	if reason, ok := ctx.Value(cancellationReasonContextKey{}).(paymentdomain.CancellationReason); ok {
		return reason
	}

	return paymentdomain.CancellationReasonUnspecified
}
