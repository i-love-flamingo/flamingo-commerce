package controller

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/category/application"
	"flamingo.me/flamingo-commerce/v3/category/domain"
	productApplication "flamingo.me/flamingo-commerce/v3/product/application"
	searchDomain "flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
)

type (
	// ViewController provides web-specific actions for category single view
	ViewController struct {
		commandHandler    QueryHandler
		breadcrumbService *application.BreadcrumbService
		responder         *web.Responder
		router            *web.Router
		template          string
		teaserTemplate    string
	}

	// ViewData for rendering context
	ViewData struct {
		ProductSearchResult *productApplication.SearchResult
		Category            domain.Category
		CategoryTree        domain.Tree
		SearchMeta          searchDomain.SearchMeta
		PaginationInfo      utils.PaginationInfo
	}
)

// Inject the ViewController controller required dependencies
func (vc *ViewController) Inject(
	queryCommandHandler QueryHandler,
	breadcrumbService *application.BreadcrumbService,
	responder *web.Responder,
	router *web.Router,
	config *struct {
		Template       string `inject:"config:commerce.category.view.template"`
		TeaserTemplate string `inject:"config:commerce.category.view.teaserTemplate"`
	},
) *ViewController {
	vc.commandHandler = queryCommandHandler
	vc.responder = responder
	vc.router = router
	vc.breadcrumbService = breadcrumbService

	if config != nil {
		vc.template = config.Template
		vc.teaserTemplate = config.TeaserTemplate
	}

	return vc
}

// Get Action to display a category page for the web
func (vc *ViewController) Get(c context.Context, request *web.Request) web.Result {

	result, redirect, err := vc.commandHandler.Execute(c, Request{
		Code:     request.Params["code"],
		Name:     request.Params["name"],
		URL:      *request.Request().URL,
		QueryAll: request.QueryAll(),
	})

	if err == domain.ErrNotFound || err == searchDomain.ErrNotFound {
		return vc.responder.NotFound(err)
	}
	if err != nil {
		return vc.responder.ServerError(err)
	}

	if redirect != nil {
		redirectParams := map[string]string{
			"code": redirect.Code,
			"name": redirect.Name,
		}
		u, _ := vc.router.Relative("category.view", redirectParams)
		u.RawQuery = request.QueryAll().Encode()
		return vc.responder.URLRedirect(u).Permanent()
	}

	// Deprecated
	vc.breadcrumbService.AddBreadcrumb(c, result.CategoryTree)

	var template string
	switch result.Category.CategoryType() {
	case domain.TypeTeaser:
		template = vc.teaserTemplate
	default:
		template = vc.template
	}

	return vc.responder.Render(template, ViewData{
		ProductSearchResult: result.ProductSearchResult,
		Category:            result.Category,
		CategoryTree:        result.CategoryTree,
		SearchMeta:          result.SearchMeta,
		PaginationInfo:      result.PaginationInfo,
	})
}
