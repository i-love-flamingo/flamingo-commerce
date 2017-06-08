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

	CartController struct {
		*responder.RenderAware `inject:""`
	}

	CartApiController struct {
		*responder.JSONAware `inject:""`
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
func (cc *CartApiController) Get(c web.Context) web.Response {
	//Cart := cc.Cartservice.GetSessionCart()
	//JsonCart, _ := json.Marshal(Cart)
	fmt.Println("Cart API: Get Cart")
	return &web.ContentResponse{
		Status: http.StatusOK,
		//Body:        bytes.NewReader([]byte(JsonCart)),
		ContentType: "text/html; charset=utf-8",
	}
}

// Update Cart via JSON
func (cc *CartApiController) Post(c web.Context) web.Response {
	//Cart := cc.Cartservice.GetSessionCart()
	cartJsonData := c.MustForm1("cart")
	//json.Unmarshal([]byte(cartJsonData), Cart)

	fmt.Println("Cart API: Update Cart")
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte("Update:" + cartJsonData)),
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
