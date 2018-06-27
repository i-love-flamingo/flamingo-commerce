package controller

import (
	"errors"

	"net/url"

	"flamingo.me/flamingo-commerce/breadcrumbs"
	"flamingo.me/flamingo-commerce/category/domain"
	productdomain "flamingo.me/flamingo-commerce/product/domain"
	productInterfaceViewData "flamingo.me/flamingo-commerce/product/interfaces/viewData"
	searchdomain "flamingo.me/flamingo-commerce/search/domain"
	"flamingo.me/flamingo-commerce/search/utils"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder.ErrorAware                                         `inject:""`
		responder.RenderAware                                        `inject:""`
		responder.RedirectAware                                      `inject:""`
		domain.CategoryService                                       `inject:""`
		productdomain.SearchService                                  `inject:""`
		*productInterfaceViewData.ProductSearchResultViewDataFactory `inject:""`
		Router                                                       *router.Router               `inject:""`
		Template                                                     string                       `inject:"config:core.category.view.template"`
		Logger                                                       flamingo.Logger              `inject:""`
		PaginationInfoFactory                                        *utils.PaginationInfoFactory `inject:""`
	}

	// ViewData for rendering context
	ViewData struct {
		ProductSearchResult productInterfaceViewData.ProductSearchResultViewData
		Category            domain.Category
		CategoryTree        domain.Category
		SearchMeta          searchdomain.SearchMeta
		PaginationInfo      utils.PaginationInfo
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
	if name, _ := c.Param1("name"); expectedName != name {

		redirectParams := router.P{
			"code": category.Code(),
			"name": expectedName,
		}

		u := vc.Router.URL("category.view", redirectParams)
		u.RawQuery = url.Values(c.QueryAll()).Encode()
		return vc.RedirectPermanentURL(u.String())
	}

	filter := make([]searchdomain.Filter, len(c.QueryAll())+1)
	filter[0] = domain.NewCategoryFacet(category)
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
	result := vc.ProductSearchResultViewDataFactory.NewProductSearchResultViewDataFromResult(c, products)

	paginationInfo := vc.PaginationInfoFactory.Build(result.SearchMeta.Page, result.SearchMeta.NumResults, 30, result.SearchMeta.NumPages, c.Request().URL)

	return vc.Render(c, vc.Template, ViewData{
		Category:            category,
		CategoryTree:        categoryRoot,
		ProductSearchResult: result,
		SearchMeta:          result.SearchMeta,
		PaginationInfo:      paginationInfo,
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
