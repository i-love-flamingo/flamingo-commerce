package controller

import (
	"flamingo/framework/web/responder"
	"flamingo/core/cart/application"
	"flamingo/framework/web"
	"flamingo/core/cart/domain"
	"net/http"
	"bytes"
	"fmt"
	"encoding/json"
	"strconv"
)

type (
	BaseCartController struct {
		*responder.RenderAware `inject:""`

		Cartservice *application.Cartservice `inject:""`

		//pageservice interfaces.PageService
	}
	// ViewData is used for product rendering
	ViewData struct {
		Cart *domain.Cart
	}

	CartViewController struct {
		*BaseCartController `inject:""`
	}

	CartApiController struct {
		*BaseCartController `inject:""`
	}

	CartItemAddApiController struct {
		*BaseCartController `inject:""`
	}
)



// Get the Cart View ( / cart)
func (cc *CartViewController) Get(c web.Context) web.Response {
	Cart := cc.Cartservice.GetSessionCart()
	//return cc.Render(c, "pages/home", ViewData{ Cart: Cart})

	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte(fmt.Sprintf("Here is the Cart: %s",Cart))),
		ContentType: "text/html; charset=utf-8",
	}

}




// Get JSON Format of API
func (cc *CartApiController) Get(c web.Context) web.Response {
	Cart := cc.Cartservice.GetSessionCart()
	JsonCart,_ := json.Marshal(Cart)
	fmt.Println("Cart API: Get Cart")
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte(JsonCart)),
		ContentType: "text/html; charset=utf-8",
	}
}


// Update Cart via JSON
func (cc *CartApiController) Post(c web.Context) web.Response {
	Cart := cc.Cartservice.GetSessionCart()
	cartJsonData := c.Form1("cart")
	json.Unmarshal([]byte(cartJsonData),Cart)

	fmt.Println("Cart API: Update Cart")
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte("Update:"+cartJsonData)),
		ContentType: "text/html; charset=utf-8",
	}
}



// Add Item to cart
func (cc *CartItemAddApiController) Get(c web.Context) web.Response {
	productCode := c.QueryFirst("code")
	qty,_ := strconv.Atoi(c.QueryFirst("qty"))

	cc.Cartservice.AddItem(productCode, qty)

	fmt.Println("Cart Item API: Add Item in Cart")
	return &web.ContentResponse{
		Status:      http.StatusOK,
		Body:        bytes.NewReader([]byte(fmt.Sprintf("Added Item?: %s qty",productCode,qty))),
		ContentType: "text/html; charset=utf-8",
	}
}
