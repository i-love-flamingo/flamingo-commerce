package order_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/order"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(config.Map{"commerce.order.useFakeAdapter": true}, new(order.Module)); err != nil {
		t.Error(err)
	}
}
