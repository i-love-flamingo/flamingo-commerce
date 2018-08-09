package viewData

import (
	"net/url"

	"flamingo.me/flamingo-commerce/product/domain"
	searchdomain "flamingo.me/flamingo-commerce/search/domain"
	"flamingo.me/flamingo-commerce/search/utils"
)

type (
	ProductSearchResultViewDataFactory struct {
		PaginationInfoFactory *utils.PaginationInfoFactory `inject:""`
		PageSize              float64                      `inject:"pagination.defaultPageSize"`
	}

	//ProductSearchResultViewData - struct with common values typical for views that show product search results
	ProductSearchResultViewData struct {
		Products       []domain.BasicProduct
		SearchMeta     searchdomain.SearchMeta
		Facets         searchdomain.FacetCollection
		PaginationInfo utils.PaginationInfo
	}
)

func (f *ProductSearchResultViewDataFactory) NewProductSearchResultViewDataFromResult(url *url.URL, products domain.SearchResult) ProductSearchResultViewData {
	return ProductSearchResultViewData{
		Products:       products.Hits,
		SearchMeta:     products.SearchMeta,
		Facets:         products.Facets,
		PaginationInfo: f.PaginationInfoFactory.Build(products.SearchMeta.Page, products.SearchMeta.NumResults, int(f.PageSize), products.SearchMeta.NumPages, url),
	}
}
