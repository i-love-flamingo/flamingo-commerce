package customer_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/customer"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(config.Map{"commerce.customer.useNilCustomerAdapter": true}, new(customer.Module)); err != nil {
		t.Error(err)
	}
}
