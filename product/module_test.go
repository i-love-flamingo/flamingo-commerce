package product_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/product"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(product.Module)); err != nil {
		t.Error(err)
	}
}
