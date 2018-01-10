package controller

import (
	"errors"

	"go.aoe.com/flamingo/core/breadcrumbs"
	"go.aoe.com/flamingo/core/category/domain"
	productdomain "go.aoe.com/flamingo/core/product/domain"
	searchdomain "go.aoe.com/flamingo/core/search/domain"
	"go.aoe.com/flamingo/core/search/utils"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder.ErrorAware        `inject:""`
		responder.RenderAware       `inject:""`
		responder.RedirectAware     `inject:""`
		domain.CategoryService      `inject:""`
		productdomain.SearchService `inject:""`

		Router                *router.Router              `inject:""`
		Template              string                      `inject:"config:core.category.view.template"`
		PaginationInfoFactory utils.PaginationInfoFactory `inject:""`
	}

	// ViewData for rendering context
	ViewData struct {
		Category       domain.Category
		CategoryTree   domain.Category
		Products       []productdomain.BasicProduct
		SearchMeta     searchdomain.SearchMeta
		Facets         searchdomain.FacetCollection
		PaginationInfo utils.PaginationInfo
	}
)

// URL to category
func URL(code string) (string, map[string]string) {
	return "category.view", map[string]string{"code": code}
}

// URLWithName points to a category with a given name
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
	categoryRoot, err := vc.CategoryService.Tree(c, c.MustParam1("code"))
	if err == domain.ErrNotFound {
		return vc.ErrorNotFound(c, err)
	} else if err != nil {
		return vc.Error(c, err)
	}

	category := getActive(categoryRoot)

	if category == nil {
		return vc.ErrorNotFound(c, errors.New("Active Category not found"))
	}

	expectedName := web.URLTitle(category.Name())
	if expectedName != c.MustParam1("name") {
		return vc.Redirect("category.view", router.P{
			"code": category.Code(),
			"name": expectedName,
		})
	}

	filter := make([]searchdomain.Filter, len(c.QueryAll())+1)
	filter[0] = domain.NewCategoryFacet(c.MustParam1("code"))
	i := 1
	for k, v := range c.QueryAll() {
		filter[i] = searchdomain.NewKeyValueFilter(k, v)
		i++
	}

	products, err := vc.SearchService.Search(c, filter...)
	if err != nil {
		return vc.Error(c, err)
	}

	vc.addBreadcrumb(c, categoryRoot)

	return vc.Render(c, vc.Template, ViewData{
		Category:       category,
		CategoryTree:   categoryRoot,
		Products:       products.Hits,
		SearchMeta:     products.SearchMeta,
		Facets:         products.Facets,
		PaginationInfo: vc.PaginationInfoFactory.Build(products.SearchMeta.Page, products.SearchMeta.NumResults, 30, products.SearchMeta.NumPages, c.Request().URL),
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
