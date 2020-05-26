package price_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/price"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(nil, new(price.Module)); err != nil {
		t.Error(err)
	}
}
