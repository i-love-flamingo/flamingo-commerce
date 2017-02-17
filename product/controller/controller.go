package controller

import (
	"flamingo/core/flamingo/web"
	"flamingo/core/flamingo/web/responder"
	"flamingo/core/product/interfaces"
	"flamingo/core/product/models"
)

type (
	ViewController struct {
		*responder.RenderAware `inject:""`

		interfaces.ProductService `inject:""`
	}

	ViewData struct {
		Product models.Product
	}
)

func (vc *ViewController) Get(c web.Context) web.Response {
	product := vc.ProductService.Get(c.Param1("sku"))

	return vc.Render(c, "pages/product/view", ViewData{Product: product})
}
