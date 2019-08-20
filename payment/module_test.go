package payment_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/payment"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(payment.Module)); err != nil {
		t.Error(err)
	}
}
