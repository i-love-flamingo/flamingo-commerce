package interfaces

import (
	"flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"net/url"
)

type (
	// ViewController demonstrates a product view controller
	ViewController struct {
		*responder.ErrorAware    `inject:""`
		*responder.RenderAware   `inject:""`
		*responder.RedirectAware `inject:""`
		domain.ProductService    `inject:""`
	}

	// ViewData is used for product rendering
	ViewData struct {
		Product *domain.Product
	}
)

// Get Response for Product matching sku param
func (vc *ViewController) Get(c web.Context) web.Response {
	product, err := vc.ProductService.Get(c, c.MustParam1("uid"))

	// catch error
	if err != nil {
		return vc.Error(c, err)
	}

	// normalize URL
	if url.QueryEscape(product.InternalName) != c.MustParam1("name") {
		return vc.Redirect("product.view", router.P{"uid": c.MustParam1("uid"), "name": url.QueryEscape(product.InternalName)})
	}

	// render page
	return vc.Render(c, "pages/product/configurable", ViewData{Product: product})
}
