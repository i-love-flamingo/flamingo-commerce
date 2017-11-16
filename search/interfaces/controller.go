package interfaces

import (
	"go.aoe.com/flamingo/core/search/domain"
	"go.aoe.com/flamingo/framework/router"
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

	// ViewData is used for search rendering
	ViewData struct {
		SearchResults map[string]domain.Result
		SearchHost    string
	}
)

func getSearchType(st string) string {
	switch st {
	case
		"retailer",
		"location",
		"brand":
		return st
	}
	return "product"
}

// Get Response for search
func (vc *ViewController) Get(c web.Context) web.Response {
	query, queryErr := c.Query1("q")
	_ = queryErr
	searchType := getSearchType(c.MustParam1("type"))

	if searchType != c.MustParam1("type") {
		return vc.Redirect("search.search?q="+query, router.P{"type": searchType})
	}

	//if query == "" || queryErr != nil {
	//	return vc.Render(c, "search/search", vd)
	//}

	//searchResult, err := vc.SearchService.Search(c, c.Request().URL.Query())
	searchResult, err := vc.SearchService.Search(c)
	if err != nil {
		return vc.Error(c, err)
	}

	vd := ViewData{
		SearchResults: searchResult,
		SearchHost:    c.Request().Host,
	}

	// render page
	return vc.Render(c, "search/search", vd)
}
