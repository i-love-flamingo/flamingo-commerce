package controller

import (
	"context"
	"net/url"
	"strconv"

	breadcrumb "flamingo.me/flamingo-commerce/v3/category/application"

	"flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/product/application"
	searchApplication "flamingo.me/flamingo-commerce/v3/search/application"
	searchdomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// ViewController providing actions for category single view
	ViewController struct {
		responder *web.Responder
		domain.CategoryService
		SearchService         *application.ProductSearchService
		router                *web.Router
		template              string
		teaserTemplate        string
		paginationInfoFactory *utils.PaginationInfoFactory
		breadcrumbService     *breadcrumb.BreadcrumbService
	}

	// ViewData for rendering context
	ViewData struct {
		ProductSearchResult *application.SearchResult
		Category            domain.Category
		CategoryTree        domain.Tree
		SearchMeta          searchdomain.SearchMeta
		PaginationInfo      utils.PaginationInfo
	}
)

// Inject the ViewController controller required dependencies
func (vc *ViewController) Inject(
	responder *web.Responder,
	categoryService domain.CategoryService,
	searchService *application.ProductSearchService,
	router *web.Router,
	paginationInfoFactory *utils.PaginationInfoFactory,
	breadcrumbService *breadcrumb.BreadcrumbService,

	config *struct {
		Template       string `inject:"config:commerce.category.view.template"`
		TeaserTemplate string `inject:"config:commerce.category.view.teaserTemplate"`
	},
) *ViewController {
	vc.responder = responder
	vc.CategoryService = categoryService
	vc.SearchService = searchService
	vc.router = router
	vc.paginationInfoFactory = paginationInfoFactory
	if config != nil {
		vc.template = config.Template
		vc.teaserTemplate = config.TeaserTemplate
	}
	vc.breadcrumbService = breadcrumbService

	return vc
}

// Get Action to display a category page
func (vc *ViewController) Get(c context.Context, request *web.Request) web.Result {
	treeRoot, err := vc.CategoryService.Tree(c, request.Params["code"])
	if err == domain.ErrNotFound {
		return vc.responder.NotFound(err)
	} else if err != nil {
		return vc.responder.ServerError(err)
	}
	currentCategory, err := vc.CategoryService.Get(c, request.Params["code"])
	if err == domain.ErrNotFound {
		return vc.responder.NotFound(err)
	} else if err != nil {
		return vc.responder.ServerError(err)
	}

	//Normalize url if required:
	expectedName := web.URLTitle(currentCategory.Name())
	if name, _ := request.Params["name"]; expectedName != name {
		redirectParams := map[string]string{
			"code": currentCategory.Code(),
			"name": expectedName,
		}
		u, _ := vc.router.Relative("category.view", redirectParams)
		u.RawQuery = url.Values(request.QueryAll()).Encode()
		return vc.responder.URLRedirect(u).Permanent()
	}

	searchRequest := &searchApplication.SearchRequest{}
	for k,v := range request.QueryAll() {
		switch k {
		case "page":
			page,_ := strconv.ParseInt(v[0],10,64)
			searchRequest.SetAdditionalFilter(searchdomain.NewPaginationPageFilter(int(page)))
			break
		default:
			searchRequest.SetAdditionalFilter(searchdomain.NewKeyValueFilter(k,v))
		}
	}
	searchRequest.SetAdditionalFilter(domain.NewCategoryFacet(currentCategory.Code()))

	products, err := vc.SearchService.Find(c, searchRequest)
	if err != nil {
		return vc.responder.ServerError(err)
	}

	vc.breadcrumbService.AddBreadcrumb(c, treeRoot)

	paginationInfo := vc.paginationInfoFactory.Build(products.SearchMeta.Page, products.SearchMeta.NumResults, 30, products.SearchMeta.NumPages, request.Request().URL)

	var template string
	switch currentCategory.CategoryType() {
	case domain.TypeTeaser:
		template = vc.teaserTemplate
	default:
		template = vc.template
	}

	return vc.responder.Render(template, ViewData{
		Category:            currentCategory,
		CategoryTree:        treeRoot,
		ProductSearchResult: products,
		SearchMeta:          products.SearchMeta,
		PaginationInfo:      paginationInfo,
	})
}
