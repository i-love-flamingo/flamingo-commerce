package checkout_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/checkout"
	"flamingo.me/flamingo/v3/core/oauth"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	cfg := &config.Module{Map: new(oauth.Module).DefaultConfig()}
	cfg.Add(config.Map{
		"session.backend": "memory",
	})

	if err := dingo.TryModule(cfg, new(checkout.Module)); err != nil {
		t.Error(err)
	}
}
