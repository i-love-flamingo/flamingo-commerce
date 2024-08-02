//go:build integration
// +build integration

package moduletest

import (
	"net/http"
	"testing"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/config"

	"flamingo.me/flamingo-commerce/v3/product"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
)

func Test_ProductPage(t *testing.T) {
	info := integrationtest.Bootup(
		[]dingo.Module{
			new(product.Module),
		},
		"",
		config.Map{
			"flamingo.systemendpoint.serviceAddr": ":0",
			"commerce.product.fakeservice": config.Map{
				"enabled": true,
				"sorting": config.Slice{
					config.Map{
						"desc":  "recommended",
						"key":   "recommended",
						"label": "Recommended",
					},
				},
			},
		},
	)
	defer info.ShutdownFunc()

	e := integrationtest.NewHTTPExpect(t, "http://"+info.BaseURL)
	e.GET("/product/fake_configurable/typeconfigurable-product.html").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("RenderContext").Equal("configurable")

	e.GET("/api/v1/products/fake_simple").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("Success").Boolean().True()

	e.GET("/api/v1/products/404").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("Success").Boolean().False()
}
