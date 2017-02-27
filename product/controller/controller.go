package controller

import (
	"flamingo/core/flamingo/web"
	"flamingo/core/flamingo/web/responder"
	"flamingo/core/product/interfaces"
	"flamingo/core/product/models"
)

type (
	// ViewController demonstrates a product view controller
	ViewController struct {
		*responder.RenderAware    `inject:""`
		interfaces.ProductService `inject:""`
	}

	// ViewData is used for product rendering
	ViewData struct {
		Product models.Product
	}
)

// Get Response for Product matching sku param
func (vc *ViewController) Get(c web.Context) web.Response {
	product := vc.ProductService.Get(c, c.Param1("sku"))

	return vc.Render(c, "pages/product/view", ViewData{Product: product})
}
