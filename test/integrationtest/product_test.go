// +build integration

package integrationtest

import (
	"github.com/gavv/httpexpect"
	"net/http"

	"testing"
)

func Tes_ProductPage(t *testing.T) {
	e := httpexpect.New(t, "http://localhost:3210")
	e.GET("/en/product/fake_configurable/typeconfigurable-product.html").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("RenderContext").Equal("configurable")
}
