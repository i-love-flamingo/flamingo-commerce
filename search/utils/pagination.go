package utils

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

type (
	//PaginationConfig - represents configuration Options used by the PaginationInfo Build method
	PaginationConfig struct {
		ShowFirstPage bool `inject:"config:commerce.pagination.showFirstPage"`
		ShowLastPage  bool `inject:"config:commerce.pagination.showLastPage"`
		//ShowAroundActivePageAmount - amount of pages to show before and after the current page (so a value of2 would show 2 pages before and 2 pages after)
		ShowAroundActivePageAmount float64 `inject:"config:commerce.pagination.showAroundActivePageAmount"`
		NameSpace                  string
	}

	// CurrentResultInfos page information
	CurrentResultInfos struct {
		ActivePage int
		TotalHits  int
		PageSize   int
		LastPage   int
	}

	// PaginationInfo meta information
	PaginationInfo struct {
		NextPage       *Page
		PreviousPage   *Page
		TotalHits      int
		PageNavigation []Page
	}

	// Page page data
	Page struct {
		Page     int
		URL      string
		IsActive bool
		IsSpacer bool
	}

	// PaginationInfoFactory - used to build a configuration based on configured defaults
	PaginationInfoFactory struct {
		DefaultConfig *PaginationConfig `inject:""`
	}
)

// BuildWith builds a paginationInfo based on the given infos and config
func BuildWith(currentResult CurrentResultInfos, paginationConfig PaginationConfig, urlBase *url.URL) PaginationInfo {
	if currentResult.PageSize < 1 {
		currentResult.PageSize = 1
	}
	if currentResult.ActivePage < 1 {
		currentResult.ActivePage = 1
	}
	paginationInfo := PaginationInfo{
		TotalHits: currentResult.TotalHits,
	}
	var pagesToAdd []int
	if paginationConfig.ShowFirstPage {
		pagesToAdd = append(pagesToAdd, 1)
	}
	if paginationConfig.ShowLastPage {
		pagesToAdd = append(pagesToAdd, currentResult.LastPage)
	}
	if currentResult.ActivePage > 1 {
		paginationInfo.PreviousPage = &Page{
			Page: currentResult.ActivePage - 1,
			URL:  makeURL(urlBase, currentResult.ActivePage-1, paginationConfig.NameSpace),
		}
	}
	if currentResult.ActivePage < currentResult.LastPage {
		paginationInfo.NextPage = &Page{
			Page: currentResult.ActivePage + 1,
			URL:  makeURL(urlBase, currentResult.ActivePage+1, paginationConfig.NameSpace),
		}
	}

	pagesToAdd = append(pagesToAdd, currentResult.ActivePage)
	showAroundActivePageAmount := int(paginationConfig.ShowAroundActivePageAmount)
	for i := currentResult.ActivePage - showAroundActivePageAmount; i <= currentResult.ActivePage+showAroundActivePageAmount; i++ {
		if i > 0 && i < currentResult.LastPage {
			pagesToAdd = append(pagesToAdd, i)
		}
	}

	if currentResult.ActivePage == 1 && currentResult.LastPage > currentResult.ActivePage+2 {
		pagesToAdd = append(pagesToAdd, currentResult.ActivePage+2)
	}

	if currentResult.ActivePage == currentResult.LastPage && currentResult.LastPage > 2 {
		pagesToAdd = append(pagesToAdd, currentResult.ActivePage-2)
	}

	sort.Ints(pagesToAdd)

	previousPageNr := 0
	for _, pageNr := range pagesToAdd {
		//guard same pages / deduplication
		if previousPageNr == pageNr {
			continue
		}
		// add spacer if not subsequent pages
		if pageNr > previousPageNr+1 {
			paginationInfo.PageNavigation = append(paginationInfo.PageNavigation, Page{
				IsSpacer: true,
			})
		}
		page := Page{
			Page:     pageNr,
			IsActive: pageNr == currentResult.ActivePage,
			IsSpacer: false,
			URL:      makeURL(urlBase, pageNr, paginationConfig.NameSpace),
		}
		paginationInfo.PageNavigation = append(paginationInfo.PageNavigation, page)
		previousPageNr = pageNr
	}
	return paginationInfo
}

// Build Pagination with the default configuration
func (f *PaginationInfoFactory) Build(activePage int, totalHits int, pageSize int, lastPage int, urlBase *url.URL) PaginationInfo {
	return BuildWith(CurrentResultInfos{
		ActivePage: activePage,
		TotalHits:  totalHits,
		PageSize:   pageSize,
		LastPage:   lastPage,
	}, *f.DefaultConfig, urlBase)
}

func makeURL(base *url.URL, page int, namespace string) string {
	q := base.Query()
	parameterName := "page"
	if namespace != "" {
		parameterName = fmt.Sprintf("%v.%v", namespace, "page")
	}
	q.Set(parameterName, strconv.Itoa(page))
	return (&url.URL{RawQuery: q.Encode()}).String()
}
