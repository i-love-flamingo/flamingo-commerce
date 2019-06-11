package cart_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(cart.Module).DefaultConfig(),
	}

	cfgModule.Map["session.backend"] = ""
	cfgModule.Map["oauth.useFake"] = true
	cfgModule.Map["oauth.preventSimultaneousSessions"] = true

	if err := dingo.TryModule(cfgModule, new(cart.Module)); err != nil {
		t.Error(err)
	}
}
