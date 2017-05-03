package controller

import (
	"flamingo/core/product/interfaces"
	"flamingo/core/product/models"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

type (
	// ViewController demonstrates a product view controller
	ViewController struct {
		*responder.ErrorAware     `inject:""`
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
	product, errorData := vc.ProductService.Get(c, c.ParamAll()["Uid"])

	if errorData.HasError() {
		response := vc.RenderError(c, errorData)

		return response
	}

	return vc.Render(c, "pages/product/configurable", ViewData{Product: product})
}
