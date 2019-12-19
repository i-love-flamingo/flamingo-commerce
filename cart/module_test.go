package cart_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(config.Map{
		"flamingo.session.backend":               "memory",
		"core.oauth.useFake":                     true,
		"core.oauth.preventSimultaneousSessions": true,
		"core.oauth.secret":                      "secret",
		"core.oauth.server":                      "server",
		"core.oauth.clientid":                    "clientid",
	}, new(cart.Module)); err != nil {
		t.Error(err)
	}
}
