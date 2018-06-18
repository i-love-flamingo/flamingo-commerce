package viewData

import (
	"flamingo.me/flamingo-commerce/product/domain"
	searchdomain "flamingo.me/flamingo-commerce/search/domain"
	"flamingo.me/flamingo-commerce/search/utils"
	"flamingo.me/flamingo/framework/web"
)

type (
	ProductSearchResultLiveViewDataFactory struct {
		PaginationInfoFactory *utils.PaginationInfoFactory `inject:""`
		PageSize              float64                      `inject:"pagination.defaultPageSize,optional"`
	}

	//ProductSearchResultViewData - struct with common values typical for views that show product search results
	ProductSearchResultLiveViewData struct {
		Products       []domain.BasicProduct
		SearchMeta     searchdomain.SearchMeta
		Suggestions    []searchdomain.Suggestion
		PaginationInfo utils.PaginationInfo
	}
)

func (f *ProductSearchResultLiveViewDataFactory) NewProductSearchResultLiveViewDataFromResult(c web.Context, products domain.SearchResult) ProductSearchResultLiveViewData {
	if f.PageSize == 0 {
		f.PageSize = 36
	}
	return ProductSearchResultLiveViewData{
		Products:       products.Hits,
		SearchMeta:     products.SearchMeta,
		Suggestions:    products.Suggestion,
		PaginationInfo: f.PaginationInfoFactory.Build(products.SearchMeta.Page, products.SearchMeta.NumResults, int(f.PageSize), products.SearchMeta.NumPages, c.Request().URL),
	}
}
