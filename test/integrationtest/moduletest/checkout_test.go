//go:build integration
// +build integration

package moduletest

import (
	"net/http"
	"testing"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/auth/fake"

	"flamingo.me/flamingo/v3/framework/config"

	"flamingo.me/flamingo-commerce/v3/checkout"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
)

func Test_CheckoutStartPage(t *testing.T) {
	info := integrationtest.Bootup(
		[]dingo.Module{
			new(auth.WebModule),
			new(fake.Module),
			new(checkout.Module),
		},
		"",
		config.Map{
			"flamingo.systemendpoint.serviceAddr":     ":0",
			"commerce.product.fakeservice.enabled":    true,
			"commerce.cart.emailAdapter.emailAddress": "test@test.de",
			"commerce.customer.useNilCustomerAdapter": true,
			"core.auth.web": config.Map{
				"debugController": true,
				"broker": config.Slice{
					config.Map{
						"broker":                       "fake",
						"typ":                          "fake",
						"userConfig.username.password": "password",
						"validatePassword":             true,
						"usernameFieldId":              "username",
						"passwordFieldId":              "password",
					},
				},
			},
		},
	)
	defer info.ShutdownFunc()

	e := integrationtest.NewHTTPExpect(t, "http://"+info.BaseURL)
	e.GET("/checkout/start").
		Expect().
		Status(http.StatusOK).JSON().Object().Value("CustomerLoggedIn").Equal(false)

}
