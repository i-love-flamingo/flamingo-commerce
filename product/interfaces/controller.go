package interfaces

import (
	"flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

type (
	// ViewController demonstrates a product view controller
	ViewController struct {
		responder.ErrorAware    `inject:""`
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		domain.ProductService   `inject:""`
	}

	// ViewData is used for product rendering
	ViewData struct {
		Product domain.Product
	}
)

// Get Response for Product matching sku param
func (vc *ViewController) Get(c web.Context) web.Response {
	product, err := vc.ProductService.Get(c, c.MustParam1("uid"))

	// catch error
	if err != nil {
		switch errors.Cause(err).(type) {
		case domain.ProductNotFound:
			return vc.ErrorNotFound(c, err)

		default:
			return vc.Error(c, err)
		}
	}
	fmt.Println(product)

	// normalize URL
	if url.QueryEscape(product.InternalName) != c.MustParam1("name") {
		return vc.Redirect("product.view", router.P{"uid": c.MustParam1("uid"), "name": url.QueryEscape(product.InternalName)})
	}

	// render page
	if product.ProductType == "configurable" {
		return vc.Render(c, "product/configurable", ViewData{Product: *product})
	} else {
		return vc.Render(c, "product/simple", ViewData{Product: *product})
	}

}
