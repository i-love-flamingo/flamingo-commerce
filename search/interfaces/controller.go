package interfaces

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/search/application"
	"flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/utils"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// ViewController demonstrates a search view controller
	ViewController struct {
		Responder *web.Responder    `inject:""`
		Responder *web.Responder   `inject:""`
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
func (vc *ViewController) Get(c context.Context, r *web.Request) web.Result {
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

	if typ, ok := r.Params["type"]; ok {
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
		return vc.Responder.Render( "search/"+typ, vd)
	}

	searchResult, err := vc.SearchService.Find(c, searchRequest)
	if err != nil {
		if re, ok := err.(*domain.RedirectError); ok {
			return vc.RedirectPermanentURL(re.To)
		}

		return vc.Error(c, err)
	}
	vd.SearchResult = searchResult
	return vc.Responder.Render( "search/search", vd)
}
