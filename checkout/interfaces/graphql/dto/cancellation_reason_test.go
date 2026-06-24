package dto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	paymentdomain "flamingo.me/flamingo-commerce/v3/payment/domain"
)

func TestMapPaymentCancellationReason(t *testing.T) {
	t.Parallel()

	aborted := dto.PaymentCancellationReasonAbortedByCustomer
	unspecified := dto.PaymentCancellationReasonUnspecified
	unknown := dto.PaymentCancellationReason("SOMETHING_ELSE")

	assert.Equal(t, paymentdomain.CancellationReasonUnspecified, dto.MapPaymentCancellationReason(nil))
	assert.Equal(t, paymentdomain.CancellationReasonUnspecified, dto.MapPaymentCancellationReason(&unspecified))
	assert.Equal(t, paymentdomain.CancellationReasonAbortedByCustomer, dto.MapPaymentCancellationReason(&aborted))
	assert.Equal(t, paymentdomain.CancellationReasonUnspecified, dto.MapPaymentCancellationReason(&unknown))
}
