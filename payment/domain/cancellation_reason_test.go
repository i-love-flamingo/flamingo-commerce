package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/payment/domain"
)

func TestCancellationReason_Values(t *testing.T) {
	t.Parallel()

	// The zero value MUST be Unspecified — it is the gob default for old
	// process bytes and the proto3 default-0 on the wire (spec §8).
	assert.Equal(t, 0, int(domain.CancellationReasonUnspecified))
	assert.Equal(t, 1, int(domain.CancellationReasonAbortedByCustomer))
}
