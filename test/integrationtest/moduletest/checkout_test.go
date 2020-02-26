// +build integration

package moduletest

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo-commerce/v3/checkout"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"flamingo.me/flamingo/v3/core/oauth"
	"flamingo.me/flamingo/v3/framework/config"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

func Test_CheckoutStartPage(t *testing.T) {
	info := integrationtest.Bootup(
		[]dingo.Module{
			new(oauth.Module),
			new(checkout.Module),
		},
		"",
		config.Map{
			"commerce.product.fakeservice.enabled":    true,
			"commerce.cart.emailAdapter.emailAddress": "test@test.de",
			"commerce.customer.useNilCustomerAdapter": true,
			//Waiting for refactor of auth module to avoid this dependency
			"core.oauth.secret":   "test",
			"core.oauth.server":   "https://accounts.google.com",
			"core.oauth.clientid": "test",
		},
	)
	defer info.ShutdownFunc()

	e := httpexpect.New(t, "http://"+info.BaseUrl)
	e.GET("/checkout/start").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("CustomerLoggedIn").Equal(false)

}
