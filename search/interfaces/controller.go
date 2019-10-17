package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/search/application"
	"flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// ViewController demonstrates a search view controller
	ViewController struct {
		responder             *web.Responder
		searchService         *application.SearchService
		paginationInfoFactory *utils.PaginationInfoFactory
	}

	viewData struct {
		SearchMeta     domain.SearchMeta
		SearchResult   map[string]*application.SearchResult
		PaginationInfo utils.PaginationInfo
	}
)

// Inject dependencies
func (vc *ViewController) Inject(responder *web.Responder,
	paginationInfoFactory *utils.PaginationInfoFactory,
	searchService *application.SearchService,
) *ViewController {
	vc.responder = responder
	vc.paginationInfoFactory = paginationInfoFactory
	vc.searchService = searchService

	return vc
}

// Get Response for search
func (vc *ViewController) Get(c context.Context, r *web.Request) web.Result {
	query, _ := r.Query1("q")

	vd := viewData{
		SearchMeta: domain.SearchMeta{
			Query: query,
		},
	}

	searchRequest := application.SearchRequest{
		Query: query,
	}
	searchRequest.AddAdditionalFilters(domain.NewKeyValueFilters(r.QueryAll())...)

	if typ, ok := r.Params["type"]; ok {
		searchResult, err := vc.searchService.FindBy(c, typ, searchRequest)
		if err != nil {
			if re, ok := err.(*domain.RedirectError); ok {
				u, _ := url.Parse(re.To)
				return vc.responder.URLRedirect(u).Permanent()
			}

			return vc.responder.ServerError(err)
		}
		vd.SearchMeta = searchResult.SearchMeta
		vd.SearchMeta.Query = query
		vd.SearchResult = map[string]*application.SearchResult{typ: searchResult}
		vd.PaginationInfo = vc.paginationInfoFactory.Build(
			searchResult.SearchMeta.Page,
			searchResult.SearchMeta.NumResults,
			searchRequest.PageSize,
			searchResult.SearchMeta.NumPages,
			r.Request().URL,
		)
		return vc.responder.Render("search/"+typ, vd)
	}

	searchResult, err := vc.searchService.Find(c, searchRequest)
	if err != nil {
		if re, ok := err.(*domain.RedirectError); ok {
			u, _ := url.Parse(re.To)
			return vc.responder.URLRedirect(u).Permanent()
		}

		return vc.responder.ServerError(err)
	}
	vd.SearchResult = searchResult
	return vc.responder.Render("search/search", vd)
}
