// +build integration

package integrationtest

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

func Test_ProductPage(t *testing.T) {
	t.Log("Booting Up Flamingo Commerce Test Project")
	bootup()
	e := httpexpect.New(t, "http://localhost:3210")
	e.GET("/en/product/fake_configurable/typeconfigurable-product.html").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("RenderContext").Equal("configurable")

}
