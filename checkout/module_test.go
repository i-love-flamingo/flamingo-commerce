package checkout_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/checkout"
	"flamingo.me/flamingo/framework/dingo"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(checkout.Module)); err != nil {
		t.Error(err)
	}
}
