package interfaces

import (
	"go.aoe.com/flamingo/core/search/domain"
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
	}

	viewData struct {
		SearchMeta   domain.SearchMeta
		SearchResult map[string]domain.Result
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

	vd := viewData{
		SearchMeta: domain.SearchMeta{
			Query: c.MustQuery1("q"),
		},
	}

	if typ, err := c.Param1("type"); err == nil {
		searchResult, err := vc.SearchService.SearchFor(c, typ, filter...)
		if err != nil {
			return vc.Error(c, err)
		}
		vd.SearchResult = map[string]domain.Result{typ: searchResult}
		return vc.Render(c, "search/"+typ, vd)
	}
	searchResult, err := vc.SearchService.Search(c, filter...)
	if err != nil {
		return vc.Error(c, err)
	}
	vd.SearchResult = searchResult
	return vc.Render(c, "search/search", vd)
}
