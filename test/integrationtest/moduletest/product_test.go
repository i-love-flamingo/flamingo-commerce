// +build integration

package moduletest

import (
	"net/http"

	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/gavv/httpexpect"

	"testing"
)

func Test_ProductPage(t *testing.T) {
	info := integrationtest.Bootup(
		[]dingo.Module{
			new(product.Module),
		},
		"",
		config.Map{
			"flamingo.systemendpoint.serviceAddr":  ":0",
			"commerce.product.fakeservice.enabled": true,
		},
	)
	defer info.ShutdownFunc()

	e := httpexpect.New(t, "http://"+info.BaseURL)
	e.GET("/product/fake_configurable/typeconfigurable-product.html").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("RenderContext").Equal("configurable")
}
