// +build integration

package moduletest

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/prefixrouter"
	"github.com/gavv/httpexpect"
	"net/http"

	"testing"
)

func Test_ProductPage2(t *testing.T) {
	deferfunc, _ := integrationtest.Bootup(
		[]dingo.Module{
			new(product.Module),
			new(prefixrouter.Module),
		},
		"",
		config.Map{
			"commerce.product.fakeservice.enabled": true,
			"flamingo.router.path":                 "/en",
		},
	)
	defer deferfunc()

	e := httpexpect.New(t, "http://localhost:3210")
	e.GET("/en/product/fake_configurable/typeconfigurable-product.html").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("RenderContext").Equal("configurable")
}
