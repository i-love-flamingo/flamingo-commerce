//go:build integration
// +build integration

package restapi_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/application/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	"flamingo.me/flamingo-commerce/v3/test/integrationtest"
)

func Test_Checkout_SimplePlaceOrderProcess(t *testing.T) {

	e := integrationtest.NewHTTPExpect(t, "http://"+FlamingoURL)
	// add something to the cart
	response := e.POST("/api/v1/cart/delivery/delivery/item").WithQuery("deliveryCode", "delivery").WithQuery("marketplaceCode", "fake_simple").Expect()
	response.Status(200).JSON().Object().Value("Success").Boolean().IsEqual(true)

	// add billing
	response = e.PUT("/api/v1/cart/billing").WithFormField("firstname", "Max").WithFormField("lastname", "Mustermann").WithFormField("email", "test@test.de").Expect()
	response.Status(200).JSON().Object().Value("Success").Boolean().IsEqual(true)

	// add shipping
	response = e.PUT("/api/v1/cart/delivery/delivery/").WithFormField("deliveryAddress.firstname", "Max").WithFormField("deliveryAddress.lastname", "Mustermann").WithFormField("deliveryAddress.email", "test@test.de").Expect()
	response.Status(200).JSON().Object().Value("Success").Boolean().IsEqual(true)

	// add payment selection
	response = e.PUT("/api/v1/cart/payment-selection").WithQuery("gateway", "fake_payment_gateway").WithQuery("method", "payment_waiting_for_customer").Expect()
	response.Status(200).JSON().Object().Value("Success").Boolean().IsEqual(true)

	// start place order
	response = e.PUT("/api/v1/checkout/placeorder").WithQuery("returnURL", "http://www.example.org").Expect()
	response.Status(201).JSON().Object().Value("UUID").String().NotEmpty()
	uuid := response.Status(201).JSON().Object().Value("UUID").String().Raw()

	// get last place order context
	response = e.GET("/api/v1/checkout/placeorder").WithQuery("returnURL", "http://www.example.org").Expect()
	response.Status(200).JSON().Object().Value("UUID").String().IsEqual(uuid)

	// cancel place order
	response = e.POST("/api/v1/checkout/placeorder/cancel").Expect()
	response.Status(200).Body().IsEqual("true\n")

	// get last place order context
	response = e.GET("/api/v1/checkout/placeorder").WithQuery("returnURL", "http://www.example.org").Expect()
	response.Status(200).JSON().Object().Value("State").String().IsEqual("Failed")

	// clear last place order context
	response = e.DELETE("/api/v1/checkout/placeorder").WithQuery("returnURL", "http://www.example.org").Expect()
	response.Status(200).Body().IsEqual("true\n")

	// get last place order context
	response = e.GET("/api/v1/checkout/placeorder").WithQuery("returnURL", "http://www.example.org").Expect()
	response.Status(500).JSON().Object().Value("Message").String().IsEqual(placeorder.ErrNoPlaceOrderProcess.Error())

	// set payment selection to a working one
	response = e.PUT("/api/v1/cart/payment-selection").WithQuery("gateway", "fake_payment_gateway").WithQuery("method", "payment_completed").Expect()
	response.Status(200).JSON().Object().Value("Success").Boolean().IsEqual(true)

	// start place order again
	response = e.PUT("/api/v1/checkout/placeorder").WithQuery("returnURL", "http://www.example.org").Expect()
	response.Status(201).JSON().Object().Value("UUID").String().NotEmpty()

	// refresh place order
	response = e.POST("/api/v1/checkout/placeorder/refresh").Expect()
	response.Status(200).JSON().Object().Value("State").String().NotEmpty()
	response = e.POST("/api/v1/checkout/placeorder/refresh-blocking").Expect()
	response.Status(200).JSON().Object().Value("State").String().NotEmpty()

	// get last place order context
	response = e.GET("/api/v1/checkout/placeorder").Expect()
	response.Status(200).JSON().Object().Value("FailedReason").String().IsEqual("")
	response.Status(200).JSON().Object().Value("State").String().IsEqual(states.Success{}.Name())
}
