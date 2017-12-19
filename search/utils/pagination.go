package utils

import (
	"net/url"
	"sort"
	"strconv"
)

type (
	PaginationInfoFactory struct {
		ShowFirstPage bool `inject:"config:pagination.showFirstPage"`
		ShowLastPage  bool `inject:"config:pagination.showLastPage"`
		//ShowAroundActivePageAmount - amount of pages to show before and after the current page (so a value of2 would show 2 pages before and 2 pages after)
		ShowAroundActivePageAmount float64 `inject:"config:pagination.showAroundActivePageAmount"`
	}
	PaginationInfo struct {
		NextPage       Page
		PreviousPage   Page
		TotalHits      int
		PageNavigation []Page
	}
	Page struct {
		Page     int
		Url      string
		IsActive bool
		IsSpacer bool
	}
)

//Build Pagination
func (f *PaginationInfoFactory) Build(activePage int, totalHits int, pageSize int, lastPage int, urlBase *url.URL) PaginationInfo {
	if pageSize < 1 {
		pageSize = 1
	}
	if activePage < 1 {
		activePage = 1
	}
	paginationInfo := PaginationInfo{
		TotalHits: totalHits,
	}
	var pagesToAdd []int
	if f.ShowFirstPage {
		pagesToAdd = append(pagesToAdd, 1)
	}
	if f.ShowLastPage {
		pagesToAdd = append(pagesToAdd, lastPage)
	}
	if activePage > 1 {
		paginationInfo.PreviousPage = Page{
			Page: activePage - 1,
			Url:  makeUrl(urlBase, activePage-1),
		}
	}
	if activePage < lastPage {
		paginationInfo.NextPage = Page{
			Page: activePage + 1,
			Url:  makeUrl(urlBase, activePage+1),
		}
	}

	pagesToAdd = append(pagesToAdd, activePage)
	showAroundActivePageAmount := int(f.ShowAroundActivePageAmount)
	for i := activePage - showAroundActivePageAmount; i <= activePage+showAroundActivePageAmount; i++ {
		if i > 0 && i < lastPage {
			pagesToAdd = append(pagesToAdd, i)
		}
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
			IsActive: pageNr == activePage,
			IsSpacer: false,
			Url:      makeUrl(urlBase, pageNr),
		}
		paginationInfo.PageNavigation = append(paginationInfo.PageNavigation, page)
		previousPageNr = pageNr
	}
	return paginationInfo
}

func makeUrl(base *url.URL, page int) string {
	q := base.Query()
	q.Set("page", strconv.Itoa(page))
	return (&url.URL{RawQuery: q.Encode()}).String()
}
