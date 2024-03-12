package price_test

import (
	"testing"

	"flamingo.me/flamingo/v3/framework/config"

	"flamingo.me/flamingo-commerce/v3/price"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(price.Module)); err != nil {
		t.Error(err)
	}
}
