package interfaces

import (
	"go.aoe.com/flamingo/core/search/domain"
	"go.aoe.com/flamingo/core/search/utils"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	// ViewController demonstrates a search view controller
	ViewController struct {
		responder.ErrorAware    `inject:""`
		responder.RenderAware   `inject:""`
		responder.RedirectAware `inject:""`
		domain.SearchService    `inject:""`
		PaginationInfoFactory   *utils.PaginationInfoFactory `inject:""`
	}

	viewData struct {
		SearchMeta     domain.SearchMeta
		SearchResult   map[string]domain.Result
		PaginationInfo utils.PaginationInfo
	}
)

// Get Response for search
func (vc *ViewController) Get(c web.Context) web.Response {
	filter := make([]domain.Filter, len(c.QueryAll()))
	i := 0
	for k, v := range c.QueryAll() {
		filter[i] = domain.NewKeyValueFilter(k, v)
		i++
	}

	query, _ := c.Query1("q")

	vd := viewData{
		SearchMeta: domain.SearchMeta{
			Query: query,
		},
	}

	//if err != nil {
	//	return vc.Render(c, "search/search", vd)
	//}

	if typ, err := c.Param1("type"); err == nil {
		searchResult, err := vc.SearchService.SearchFor(c, typ, filter...)
		if err != nil {
			if re, ok := err.(*domain.RedirectError); ok {
				return vc.RedirectPermanentURL(re.To)
			}

			return vc.Error(c, err)
		}
		vd.SearchMeta = searchResult.SearchMeta
		vd.SearchMeta.Query = query
		vd.SearchResult = map[string]domain.Result{typ: searchResult}
		vd.PaginationInfo = vc.PaginationInfoFactory.Build(searchResult.SearchMeta.Page, searchResult.SearchMeta.NumResults, 30, searchResult.SearchMeta.NumPages, c.Request().URL)
		return vc.Render(c, "search/"+typ, vd)
	}

	searchResult, err := vc.SearchService.Search(c, filter...)
	if err != nil {
		if re, ok := err.(*domain.RedirectError); ok {
			return vc.RedirectPermanentURL(re.To)
		}

		return vc.Error(c, err)
	}
	vd.SearchResult = searchResult
	//vd.PaginationInfo = vc.PaginationInfoFactory.Build(1, 0, 50, c.Request().URL)
	return vc.Render(c, "search/search", vd)
}
