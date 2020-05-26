package category_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/category"
	"flamingo.me/flamingo/v3/framework/config"
)

func TestModule_Configure(t *testing.T) {
	if err := config.TryModules(config.Map{"commerce.category.useCategoryFixedAdapter": true}, new(category.Module)); err != nil {
		t.Error(err)
	}
}
