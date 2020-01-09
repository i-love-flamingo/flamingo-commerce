// +build integration

package moduletest

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/gavv/httpexpect"
	"net/http"

	"testing"
)

func Test_ProductPage2(t *testing.T) {
	info := integrationtest.Bootup(
		[]dingo.Module{
			new(product.Module),
		},
		"",
		config.Map{
			"commerce.product.fakeservice.enabled": true,
		},
	)
	defer info.ShutdownFunc()

	e := httpexpect.New(t, "http://"+info.BaseURL)
	e.GET("/product/fake_configurable/typeconfigurable-product.html").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("RenderContext").Equal("configurable")
}
