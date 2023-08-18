//go:build integration
// +build integration

package frontend_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

type (
	// CartItems list of CartItem
	CartItems []CartItem
	// CartItem used as simple CartItem representation during test
	CartItem struct {
		ProductName     string
		MarketplaceCode string
		Qty             int
	}
)

const (
	routeCheckoutSubmit     = "/en/checkout"
	routeCheckoutReview     = "/en/checkout/review"
	routeCheckoutPlaceOrder = "/en/checkout/placeorder"
	routeCheckoutSuccess    = "/en/checkout/success"
)

// CartAddProduct helper
func CartAddProduct(t *testing.T, e *httpexpect.Expect, marketplaceCode string, qty int, variantMarketplaceCode string, deliveryCode string) {
	t.Helper()
	request := e.POST("/en/cart/add/"+marketplaceCode).WithQuery("qty", qty)
	if deliveryCode != "" {
		request = request.WithQuery("deliveryCode", deliveryCode)
	}
	if variantMarketplaceCode != "" {
		request = request.WithQuery("variantMarketplaceCode", variantMarketplaceCode)
	}
	request.Expect().
		Status(http.StatusOK)
}

// CartApplyVoucher applies a voucher via api
func CartApplyVoucher(t *testing.T, e *httpexpect.Expect, code string) {
	t.Helper()
	request := e.POST("/en/api/cart/applyvoucher").WithQuery("couponCode", code)
	request.Expect().Status(http.StatusOK)
}

// CartGetItems testhelper
func CartGetItems(t *testing.T, e *httpexpect.Expect) CartItems {
	t.Helper()
	var items CartItems

	cartItems := e.GET("/en/cart/").Expect().Status(http.StatusOK).JSON().Object().
		Value("DecoratedCart").Object().
		Value("Cart").Object().
		Value("Deliveries").Array().Value(0).Object().
		Value("Cartitems").Array()

	for _, v := range cartItems.Iter() {
		items = append(items, CartItem{
			ProductName:     v.Object().Value("ProductName").String().Raw(),
			Qty:             int(v.Object().Value("Qty").Number().Raw()),
			MarketplaceCode: v.Object().Value("MarketplaceCode").String().Raw(),
		})

	}
	return items
}

// MustContain checks and returns CartItem by marketplaceCode
func (c CartItems) MustContain(t *testing.T, marketplaceCode string) *CartItem {
	t.Helper()
	for _, v := range c {
		if v.MarketplaceCode == marketplaceCode {
			return &v
		}
	}
	t.Fatal(fmt.Sprintf("No CartItem with marketplaceCode: %v Only: %#v", marketplaceCode, c))
	return nil
}

// SubmitCheckoutForm sends a POST request to the checkout route with provided form params
func SubmitCheckoutForm(t *testing.T, e *httpexpect.Expect, form map[string]interface{}) *httpexpect.Response {
	t.Helper()
	return e.POST(routeCheckoutSubmit).WithForm(form).Expect().Status(http.StatusOK)
}

// SubmitReviewForm sends a POST request to the review route with provided form params
func SubmitReviewForm(t *testing.T, e *httpexpect.Expect, form map[string]interface{}) *httpexpect.Response {
	t.Helper()
	return e.POST(routeCheckoutReview).WithForm(form).Expect().Status(http.StatusOK)
}
