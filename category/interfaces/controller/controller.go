package controller

import (
	"context"
	"errors"
	breadcrumb "flamingo.me/flamingo-commerce/v3/category/application"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// View demonstrates a product view controller
	View struct {
		responder *web.Responder
		domain.CategoryService
		SearchService         *application.ProductSearchService
		router                *web.Router
		template              string
		teaserTemplate        string
		logger                flamingo.Logger
		paginationInfoFactory *utils.PaginationInfoFactory
		breadcrumbService     *breadcrumb.BreadcrumbService
	}

	// ViewData for rendering context
	ViewData struct {
		ProductSearchResult *application.SearchResult
		Category            domain.Category
		CategoryTree        domain.Category
		SearchMeta          searchdomain.SearchMeta
		PaginationInfo      utils.PaginationInfo
	}
)

// Inject the View controller required dependencies
func (vc *View) Inject(
	responder *web.Responder,
	categoryService domain.CategoryService,
	searchService *application.ProductSearchService,
	router *web.Router,
	logger flamingo.Logger,
	paginationInfoFactory *utils.PaginationInfoFactory,
	breadcrumbService *breadcrumb.BreadcrumbService,

	config *struct {
		Template       string `inject:"config:category.view.template"`
		TeaserTemplate string `inject:"config:category.view.teaserTemplate"`
	},
) {
	vc.responder=responder
	vc.CategoryService = categoryService
	vc.SearchService = searchService
	vc.router = router
	vc.logger = logger
	vc.paginationInfoFactory = paginationInfoFactory
	vc.template = config.Template
	vc.teaserTemplate = config.TeaserTemplate
	vc.breadcrumbService = breadcrumbService
}

// Get Response for Product matching sku param
func (vc *View) Get(c context.Context, request *web.Request) web.Result {
	categoryRoot, err := vc.CategoryService.Tree(c, request.Params["code"])
	if err == domain.ErrNotFound {
		return vc.responder.NotFound(err)
	} else if err != nil {
		return vc.responder.ServerError(err)
	}

	category := domain.GetActive(categoryRoot)
	if category == nil {
		return vc.responder.NotFound(errors.New("Active Category not found"))
	}

	expectedName := web.URLTitle(category.Name())
	if name, _ := request.Params["name"]; expectedName != name {
		redirectParams := map[string]string{
			"code": category.Code(),
			"name": expectedName,
		}

		u, _ := vc.router.URL("category.view", redirectParams)
		u.RawQuery = url.Values(request.QueryAll()).Encode()
		return vc.responder.URLRedirect(u).Permanent()
	}

	queryAll := request.QueryAll()
	filter := make(map[string]interface{}, len(queryAll)+1)
	for k, v := range queryAll {
		filter[k] = v
	}
	filter[string(domain.CategoryKey)] = domain.NewCategoryFacet(category)

	searchRequest := &searchApplication.SearchRequest{
		FilterBy: filter,
	}

	products, err := vc.SearchService.Find(c, searchRequest)
	if err != nil {
		return vc.responder.ServerError(err)
	}

	vc.breadcrumbService.AddBreadcrumb(c, categoryRoot)

	paginationInfo := vc.paginationInfoFactory.Build(products.SearchMeta.Page, products.SearchMeta.NumResults, 30, products.SearchMeta.NumPages, request.Request().URL)

	var template string
	switch category.CategoryType() {
	case domain.TypeTeaser:
		template = vc.teaserTemplate
	default:
		template = vc.template
	}

	return vc.responder.Render( template, ViewData{
		Category:            category,
		CategoryTree:        categoryRoot,
		ProductSearchResult: products,
		SearchMeta:          products.SearchMeta,
		PaginationInfo:      paginationInfo,
	})
}
