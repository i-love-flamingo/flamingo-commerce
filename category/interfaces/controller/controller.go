package controller

import (
	"flamingo/core/breadcrumbs"
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
		Category     domain.Category
		CategoryTree domain.Category
		Products     []productdomain.BasicProduct
	}
)

// URL to category
func URL(code string) (string, map[string]string) {
	return "category.view", map[string]string{"code": code}
}

// URL with name to category
func URLWithName(code, name string) (string, map[string]string) {
	return "category.view", map[string]string{"code": code, "name": name}
}

func getActive(category domain.Category) domain.Category {
	for _, sub := range category.Categories() {
		if active := getActive(sub); active != nil {
			return active
		}
	}
	if category.Active() {
		return category
	}
	return nil
}

// Get Response for Product matching sku param
func (vc *View) Get(c web.Context) web.Response {
	categoryRoot, err := vc.CategoryService.Get(c, c.MustParam1("code"))
	if err == domain.NotFound {
		return vc.ErrorNotFound(c, err)
	} else if err != nil {
		return vc.Error(c, err)
	}

	category := getActive(categoryRoot)

	expectedName := web.URLTitle(category.Name())
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

	vc.addBreadcrumb(c, categoryRoot)

	return vc.Render(c, vc.Template, ViewData{
		Category:     category,
		CategoryTree: categoryRoot,
		Products:     products,
	})
}

func (vc *View) addBreadcrumb(c web.Context, category domain.Category) {
	if !category.Active() {
		return
	}
	if category.Code() != "" {
		breadcrumbs.Add(c, breadcrumbs.Crumb{
			category.Name(),
			vc.Router.URL(URLWithName(category.Code(), web.URLTitle(category.Name()))).String(),
		})
	}

	for _, subcat := range category.Categories() {
		vc.addBreadcrumb(c, subcat)
	}
}
