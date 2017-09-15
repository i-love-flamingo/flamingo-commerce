package controller

import (
	"flamingo/core/category/domain"
	productdomain "flamingo/core/product/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder.ErrorAware    `inject:""`
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		domain.CategoryService  `inject:""`

		Template string         `inject:"config:core.category.view.template"`
		Router   *router.Router `inject:""`
	}

	// ViewData for rendering context
	ViewData struct {
		Category domain.Category
		Products []productdomain.BasicProduct
	}
)

// Get Response for Product matching sku param
func (vc *View) Get(c web.Context) web.Response {
	category, err := vc.CategoryService.Get(c, c.MustParam1("categorycode"))

	// catch error
	if err == domain.NotFound {
		return vc.ErrorNotFound(c, err)
	} else if err != nil {
		return vc.Error(c, err)
	}

	expectedName := web.UrlTitle(category.Name)
	if expectedName != c.MustParam1("name") {
		return vc.Redirect("category.view", router.P{
			"categorycode": category.Code,
			"name":         expectedName,
		})
	}

	products, err := vc.CategoryService.GetProducts(c, c.MustParam1("categorycode"))

	// catch error
	if err != nil {
		return vc.Error(c, err)
	}

	return vc.Render(c, vc.Template, ViewData{
		Category: category,
		Products: products,
	})
}
