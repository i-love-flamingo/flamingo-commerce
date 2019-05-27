package category_test

import (
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/category"
)

func TestModule_Configure(t *testing.T) {
	if err := dingo.TryModule(new(category.Module)); err != nil {
		t.Error(err)
	}
}
