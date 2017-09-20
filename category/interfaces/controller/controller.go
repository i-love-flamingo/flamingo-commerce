package controller

import (
	"flamingo/core/category/domain"
	productdomain "flamingo/core/product/domain"
	searchdomain "flamingo/core/search/domain"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"flamingo/framework/web/responder"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder.ErrorAware       `inject:""`
		responder.RenderAware      `inject:""`
		responder.RedirectAware    `inject:""`
		domain.CategoryService     `inject:""`
		searchdomain.SearchService `inject:""`

		Router   *router.Router `inject:""`
		Template string         `inject:"config:core.category.view.template"`
	}

	// ViewData for rendering context
	ViewData struct {
		Category domain.Category
		Products []productdomain.BasicProduct
	}
)

// Get Response for Product matching sku param
func (vc *View) Get(c web.Context) web.Response {
	category, err := vc.CategoryService.Get(c, c.MustParam1("code"))
	if err == domain.NotFound {
		return vc.ErrorNotFound(c, err)
	} else if err != nil {
		return vc.Error(c, err)
	}

	expectedName := web.UrlTitle(category.Name())
	if expectedName != c.MustParam1("name") {
		return vc.Redirect("category.view", router.P{
			"code": category.Code(),
			"name": expectedName,
		})
	}

	_, products, _, err := vc.SearchService.GetProducts(c, searchdomain.SearchMeta{}, domain.NewCategoryFacet(c.MustParam1("code")))
	if err != nil {
		return vc.Error(c, err)
	}

	return vc.Render(c, vc.Template, ViewData{
		Category: category,
		Products: products,
	})
}
