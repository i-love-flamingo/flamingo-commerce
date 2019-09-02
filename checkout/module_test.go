package checkout_test

import (
"testing"

"flamingo.me/flamingo-commerce/v3/checkout"
"flamingo.me/dingo"
"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	cfgModule := &config.Module{
		Map: new(checkout.Module).DefaultConfig(),
	}

	cfgModule.Map["session.backend"] = ""
	cfgModule.Map["oauth.useFake"] = true
	cfgModule.Map["oauth.preventSimultaneousSessions"] = true

	if err := dingo.TryModule(cfgModule, new(checkout.Module)); err != nil {
		t.Error(err)
	}
}
