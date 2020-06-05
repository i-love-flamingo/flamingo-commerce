// +build integration

package restapi_test

import (
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
	"testing"
)

func Test_Checkout_SimplePlaceOrderProcess(t *testing.T) {

	e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
	// add something to the cart
	response := e.POST("/api/v1/cart/delivery/delivery/additem").WithQuery("deliveryCode", "delivery").WithQuery("marketplaceCode", "fake_simple").Expect()
	response.Status(200).JSON().Object().Value("Success").Boolean().Equal(true)

	// add billing
	response = e.POST("/api/v1/cart/billing").WithFormField("firstname", "Max").WithFormField("lastname", "Mustermann").WithFormField("email", "test@test.de").Expect()
	response.Status(200).JSON().Object().Value("Success").Boolean().Equal(true)

	// add shipping
	response = e.POST("/api/v1/cart/delivery/delivery/deliveryinfo").WithFormField("deliveryAddress.firstname", "Max").WithFormField("deliveryAddress.lastname", "Mustermann").WithFormField("deliveryAddress.email", "test@test.de").Expect()
	//b,err:=httputil.DumpRequest(response.Raw().Request, true)
	//t.Log(string(b),err)
	response.Status(200).JSON().Object().Value("Success").Boolean().Equal(true)

	// add shipping
	response = e.PUT("/api/v1/cart/updatepaymentselection").WithQuery("gateway", "offline").WithQuery("method", "offlinepayment_cashondelivery").Expect()
	response.Status(200).JSON().Object().Value("Success").Boolean().Equal(true)

	// start place order
	response = e.PUT("/api/v1/checkout/placeorder").WithQuery("returnURL", "http://www.example.org").Expect()
	response.Status(201).JSON().Object().Value("UUID").String().NotEmpty()

	// refresh place order
	response = e.POST("/api/v1/checkout/placeorder/refreshplaceorderblocking").Expect()
	response.Status(200).JSON().Object().Value("State").String().NotEmpty()

	// get place order
	response = e.GET("/api/v1/checkout/placeorder").Expect()
	response.Status(200).JSON().Object().Value("FailedReason").String().Equal("")
	response.Status(200).JSON().Object().Value("State").String().Equal("Success")

}
