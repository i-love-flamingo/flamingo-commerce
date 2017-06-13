package controller

import (
	"bytes"
	"flamingo/core/cart/domain"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"fmt"
	"net/http"
)

type (
	// ViewData is used for product rendering
	ViewData struct {
		Cart *domain.Cart
	}

	// CartController for carts
	CartController struct {
		responder.RenderAware `inject:""`
	}

	// CartAPIController for cart api
	CartAPIController struct {
		responder.JSONAware `inject:""`
	}
)

// Get the Cart View ( / cart)
func (cc *CartController) Get(c web.Context) web.Response {
	//Cart := cc.Cartservice.GetSessionCart()
	//return cc.Render(c, "pages/home", ViewData{ Cart: Cart})

	return &web.ContentResponse{
		Status: http.StatusOK,
		//Body:        bytes.NewReader([]byte(fmt.Sprintf("Here is the Cart: %s", Cart))),
		ContentType: "text/html; charset=utf-8",
	}

}

// Get JSON Format of API
func (cc *CartAPIController) Get(c web.Context) web.Response {
	//Cart := cc.Cartservice.GetSessionCart()
	//JsonCart, _ := json.Marshal(Cart)
	fmt.Println("Cart API: Get Cart")
	return &web.ContentResponse{
		Status: http.StatusOK,
		//Body:        bytes.NewReader([]byte(JsonCart)),
		ContentType: "text/html; charset=utf-8",
	}
}

// Post Update Cart via JSON
func (cc *CartAPIController) Post(c web.Context) web.Response {
	//Cart := cc.Cartservice.GetSessionCart()
	cartJSONData := c.MustForm1("cart")
	//json.Unmarshal([]byte(cartJSONData), Cart)

	fmt.Println("Cart API: Update Cart")
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte("Update:" + cartJSONData)),
		ContentType: "text/html; charset=utf-8",
	}
}

// Add Item to cart
//func (cc *CartItemAddApiController) AddToBasketAction(c web.Context) web.Response {
//	productCode := c.MustQuery1("code")
//	qty, _ := strconv.Atoi(c.MustQuery1("qty"))
//
//	cc.Cartservice.AddItem(productCode, qty)
//
//	fmt.Println("Cart Item API: Add Item in Cart")
//	return &web.ContentResponse{
//		Status:      http.StatusOK,
//		Body:        bytes.NewReader([]byte(fmt.Sprintf("Added Item?: %s qty", productCode, qty))),
//		ContentType: "text/html; charset=utf-8",
//	}
//}
