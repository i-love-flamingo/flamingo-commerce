package viewData

import (
	"net/url"

	"flamingo.me/flamingo-commerce/product/domain"
	searchdomain "flamingo.me/flamingo-commerce/search/domain"
	"flamingo.me/flamingo-commerce/search/utils"
)

type (
	ProductSearchResultLiveViewDataFactory struct {
		PaginationInfoFactory *utils.PaginationInfoFactory `inject:""`
		PageSize              float64                      `inject:"pagination.defaultPageSize"`
	}

	//ProductSearchResultViewData - struct with common values typical for views that show product search results
	ProductSearchResultLiveViewData struct {
		Products       []domain.BasicProduct
		SearchMeta     searchdomain.SearchMeta
		Suggestions    []searchdomain.Suggestion
		PaginationInfo utils.PaginationInfo
	}
)

func (f *ProductSearchResultLiveViewDataFactory) NewProductSearchResultLiveViewDataFromResult(url *url.URL, products domain.SearchResult) ProductSearchResultLiveViewData {
	return ProductSearchResultLiveViewData{
		Products:       products.Hits,
		SearchMeta:     products.SearchMeta,
		Suggestions:    products.Suggestion,
		PaginationInfo: f.PaginationInfoFactory.Build(products.SearchMeta.Page, products.SearchMeta.NumResults, int(f.PageSize), products.SearchMeta.NumPages, url),
	}
}
