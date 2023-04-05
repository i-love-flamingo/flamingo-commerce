package rediscart_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/rediscart"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	t.Parallel()

	if err := config.TryModules(nil, new(rediscart.Module)); err != nil {
		t.Error(err)
	}
}
