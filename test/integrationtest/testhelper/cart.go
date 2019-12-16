package testhelper

import (
	"fmt"
	"github.com/gavv/httpexpect"
	"net/http"
	"testing"
)

type (
	CartItems []CartItem
	CartItem  struct {
		ProductName     string
		MarketplaceCode string
		Qty             int
	}
)

func CartAddProduct(e *httpexpect.Expect, marketplaceCode string, qty int, variantMarketplaceCode string, deliveryCode string) {
	request := e.GET("/en/cart/add/"+marketplaceCode).WithQuery("qty", qty)
	if deliveryCode != "" {
		request = request.WithQuery("deliveryCode", deliveryCode)
	}
	if variantMarketplaceCode != "" {
		request = request.WithQuery("variantMarketplaceCode", variantMarketplaceCode)
	}
	request.Expect().
		Status(http.StatusOK)
}

func CartGetItems(e *httpexpect.Expect) CartItems {
	var items CartItems

	cartItems := e.GET("/en/cart/").Expect().Status(http.StatusOK).JSON().Object().
		Value("DecoratedCart").Object().
		Value("Cart").Object().
		Value("Deliveries").Array().Element(0).Object().
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

func (c CartItems) MustContain(t *testing.T, marketplaceCode string) *CartItem {
	for _, v := range c {
		if v.MarketplaceCode == marketplaceCode {
			return &v
		}
	}
	t.Fatal(fmt.Sprintf("No CartItem with marketplaceCode: %v Only: %#v", marketplaceCode, c))
	return nil
}
