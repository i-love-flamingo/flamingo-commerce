package interfaces

import (
	"context"

	"flamingo.me/flamingo-commerce/search/application"
	"flamingo.me/flamingo-commerce/search/domain"
	"flamingo.me/flamingo-commerce/search/utils"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	// ViewController demonstrates a search view controller
	ViewController struct {
		responder.ErrorAware    `inject:""`
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		SearchService           *application.SearchService   `inject:""`
		PaginationInfoFactory   *utils.PaginationInfoFactory `inject:""`
	}

	viewData struct {
		SearchMeta     domain.SearchMeta
		SearchResult   map[string]*application.SearchResult
		PaginationInfo utils.PaginationInfo
	}
)

// Get Response for search
func (vc *ViewController) Get(c context.Context, r *web.Request) web.Response {
	query, _ := r.Query1("q")

	vd := viewData{
		SearchMeta: domain.SearchMeta{
			Query: query,
		},
	}

	queryAll := r.QueryAll()
	filter := make(map[string]interface{})
	for k, v := range queryAll {
		filter[k] = v
	}
	searchRequest := application.SearchRequest{
		FilterBy: filter,
		Query:    query,
	}

	if typ, ok := r.Param1("type"); ok {
		searchResult, err := vc.SearchService.FindBy(c, typ, searchRequest)
		if err != nil {
			if re, ok := err.(*domain.RedirectError); ok {
				return vc.RedirectPermanentURL(re.To)
			}

			return vc.Error(c, err)
		}
		vd.SearchMeta = searchResult.SearchMeta
		vd.SearchMeta.Query = query
		vd.SearchResult = map[string]*application.SearchResult{typ: searchResult}
		vd.PaginationInfo = vc.PaginationInfoFactory.Build(
			searchResult.SearchMeta.Page,
			searchResult.SearchMeta.NumResults,
			searchRequest.PageSize,
			searchResult.SearchMeta.NumPages,
			r.Request().URL,
		)
		return vc.Render(c, "search/"+typ, vd)
	}

	searchResult, err := vc.SearchService.Find(c, searchRequest)
	if err != nil {
		if re, ok := err.(*domain.RedirectError); ok {
			return vc.RedirectPermanentURL(re.To)
		}

		return vc.Error(c, err)
	}
	vd.SearchResult = searchResult
	return vc.Render(c, "search/search", vd)
}
