package controller

import (
	"context"
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
		responder.ErrorAware
		responder.RenderAware
		responder.RedirectAware
		domain.CategoryService
		productdomain.SearchService
		*productInterfaceViewData.ProductSearchResultViewDataFactory
		router                *router.Router
		template              string
		teaserTemplate        string
		logger                flamingo.Logger
		paginationInfoFactory *utils.PaginationInfoFactory
	}

	// ViewData for rendering context
	ViewData struct {
		ProductSearchResult productInterfaceViewData.ProductSearchResultViewData
		Category            domain.Category
		CategoryTree        domain.Category
		CategoryChildren    []domain.Category
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

func (vc *View) Inject(
	errorAware responder.ErrorAware,
	renderAware responder.RenderAware,
	redirectAware responder.RedirectAware,
	categoryService domain.CategoryService,
	searchService productdomain.SearchService,
	productSearchResultViewDataFactory *productInterfaceViewData.ProductSearchResultViewDataFactory,
	router *router.Router,
	config *struct {
		Template       string `inject:"config:core.category.view.template"`
		TeaserTemplate string `inject:"config:core.category.view.teaserTemplate"`
	},
	logger flamingo.Logger,
	paginationInfoFactory *utils.PaginationInfoFactory,
) {
	vc.ErrorAware = errorAware
	vc.RenderAware = renderAware
	vc.RedirectAware = redirectAware
	vc.CategoryService = categoryService
	vc.SearchService = searchService
	vc.ProductSearchResultViewDataFactory = productSearchResultViewDataFactory
	vc.router = router
	vc.logger = logger
	vc.paginationInfoFactory = paginationInfoFactory
	vc.template = config.Template
	vc.teaserTemplate = config.TeaserTemplate
}

// Get Response for Product matching sku param
func (vc *View) Get(c context.Context, request *web.Request) web.Response {
	categoryRoot, err := vc.CategoryService.Tree(c, request.MustParam1("code"))
	if err == domain.ErrNotFound {
		return vc.ErrorNotFound(c, err)
	} else if err != nil {
		return vc.Error(c, err)
	}

	category := getActive(categoryRoot)
	if category == nil {
		return vc.ErrorNotFound(c, errors.New("Active Category not found"))
	}

	categoryChildren, err := vc.CategoryService.Children(c, category.Code())
	if err == domain.ErrNotFound {
		return vc.ErrorNotFound(c, err)
	} else if err != nil {
		return vc.Error(c, err)
	}

	expectedName := web.URLTitle(category.Name())
	if name, _ := request.Param1("name"); expectedName != name {

		redirectParams := router.P{
			"code": category.Code(),
			"name": expectedName,
		}

		u := vc.router.URL("category.view", redirectParams)
		u.RawQuery = url.Values(request.QueryAll()).Encode()
		return vc.RedirectPermanentURL(u.String())
	}

	filter := make([]searchdomain.Filter, len(request.QueryAll())+1)
	filter[0] = domain.NewCategoryFacet(category)
	i := 1
	for k, v := range request.QueryAll() {
		filter[i] = searchdomain.NewKeyValueFilter(k, v)
		i++
	}

	products, err := vc.SearchService.Search(c, filter...)
	if err != nil {
		return vc.Error(c, err)
	}

	vc.addBreadcrumb(c, categoryRoot)
	result := vc.ProductSearchResultViewDataFactory.NewProductSearchResultViewDataFromResult(request.Request().URL, products)

	paginationInfo := vc.PaginationInfoFactory.Build(result.SearchMeta.Page, result.SearchMeta.NumResults, 30, result.SearchMeta.NumPages, request.Request().URL)

	var template string
	switch category.CategoryType() {
	case domain.TypeTeaser:
		template = vc.teaserTemplate
	default:
		template = vc.template
	}

	return vc.Render(c, template, ViewData{
		Category:            category,
		CategoryTree:        categoryRoot,
		CategoryChildren:    categoryChildren,
		ProductSearchResult: result,
		SearchMeta:          result.SearchMeta,
		PaginationInfo:      paginationInfo,
	})
}

func (vc *View) addBreadcrumb(c context.Context, category domain.Category) {
	if !category.Active() {
		return
	}
	if category.Code() != "" {
		breadcrumbs.Add(c, breadcrumbs.Crumb{
			category.Name(),
			vc.router.URL(URLWithName(category.Code(), web.URLTitle(category.Name()))).String(),
		})
	}

	for _, subcat := range category.Categories() {
		vc.addBreadcrumb(c, subcat)
	}
}
