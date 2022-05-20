package templatefunctions

import (
	"net/http"
	"net/url"
	"testing"

	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

// TestFilterProcessingIsAllowed - Test different combinations of filter processing
func TestFilterProcessingIsAllowed(t *testing.T) {

	// Case 1: only blacklist filled, everything not on the blacklist should be fine
	filterConstrains := make(map[string]string)
	filterConstrains["blackList"] = "blacklisted"
	filterConstrains["whiteList"] = ""
	filterProcessing := buildFilterProcessing(nil, filterConstrains)

	assert.True(t, filterProcessing.isAllowed("allowedKey"))
	assert.False(t, filterProcessing.isAllowed("blacklisted"))

	// Case 2: only whitelist filled, everything noy on the whitelist should be bad
	filterConstrains["blackList"] = ""
	filterConstrains["whiteList"] = "allowedKey"
	filterProcessing = buildFilterProcessing(nil, filterConstrains)

	assert.True(t, filterProcessing.isAllowed("allowedKey"))
	assert.False(t, filterProcessing.isAllowed("blacklisted"))

	// Case 3: both lists are empty, everything is allowed
	filterConstrains["blackList"] = ""
	filterConstrains["whiteList"] = ""
	filterProcessing = buildFilterProcessing(nil, filterConstrains)

	assert.True(t, filterProcessing.isAllowed("allowedKey"))
	assert.True(t, filterProcessing.isAllowed("blacklisted"))

	// Case 4: both lists are filled, blacklist is ignored, because whitelisting has a higher prio
	filterConstrains["blackList"] = "blacklisted,notAllowed, notAllowedWithSpace"
	filterConstrains["whiteList"] = "allowedKey, allowedWithSpace,allowedNoSpace,blacklisted"
	filterProcessing = buildFilterProcessing(nil, filterConstrains)

	assert.True(t, filterProcessing.isAllowed("allowedKey"))
	assert.True(t, filterProcessing.isAllowed("allowedWithSpace"))
	assert.True(t, filterProcessing.isAllowed("allowedNoSpace"))
	assert.True(t, filterProcessing.isAllowed("blacklisted"))
	assert.False(t, filterProcessing.isAllowed("notAllowedWithSpace"))
	assert.False(t, filterProcessing.isAllowed("notOnAnyList"))
}

// TestFilterProcessingModifyResultRemoveBlacklisted - Test modification of the searchResult
func TestFilterProcessingModifyResultRemoveBlacklisted(t *testing.T) {
	searchResult := buildSearchResult()
	assert.Len(t, searchResult.Facets, 2)
	filterConstrains := make(map[string]string)
	filterConstrains["blackList"] = "disallowed"

	filterProcessing := buildFilterProcessing(nil, filterConstrains)

	newResult := filterProcessing.modifyResult(searchResult)
	assert.Len(t, newResult.Facets, 1)
	assert.Len(t, searchResult.SearchMeta.SelectedFacets, 1)
	assert.Equal(t, newResult.Facets["allowed"].Name, "allowed")
}

// TestFilterProcessingModifyResultOnlyKeepWhitelisted - Test modification of the searchResult
func TestFilterProcessingModifyResultOnlyKeepWhitelisted(t *testing.T) {
	searchResult := buildSearchResult()
	assert.Len(t, searchResult.Facets, 2)
	assert.Len(t, searchResult.SearchMeta.SelectedFacets, 2)
	filterConstrains := make(map[string]string)
	filterConstrains["whiteList"] = "allowed"

	filterProcessing := buildFilterProcessing(nil, filterConstrains)

	newResult := filterProcessing.modifyResult(searchResult)
	assert.Len(t, newResult.Facets, 1)
	assert.Len(t, searchResult.SearchMeta.SelectedFacets, 1)
	assert.Equal(t, newResult.Facets["allowed"].Name, "allowed")
}

// helper function for test cases
func buildSearchResult() *application.SearchResult {
	searchResult := application.SearchResult{}

	// adding a facet, one allowed, on to remove
	facetAllowed := domain.Facet{Name: "allowed"}
	facetDisallowed := domain.Facet{Name: "disallowed"}
	facetCollection := domain.FacetCollection{}
	facetCollection["allowed"] = facetAllowed
	facetCollection["disallowed"] = facetDisallowed

	searchResult.Facets = facetCollection

	var selectedFacets []domain.Facet
	selectedFacets = append(selectedFacets, facetAllowed)
	selectedFacets = append(selectedFacets, facetDisallowed)
	searchResult.SearchMeta.SelectedFacets = selectedFacets

	return &searchResult
}

// helper function for test cases
func buildFilterProcessing(keyValueFilter map[string][]string, filterConstrains map[string]string) filterProcessing {
	requestURL := url.URL{}
	webRequest := web.CreateRequest(&http.Request{
		URL: &requestURL,
	}, nil)

	// we do not yet support multiple namespaces
	namespace := ""

	filterProcessing := newFilterProcessing(webRequest, namespace, nil, keyValueFilter, filterConstrains, nil)
	return filterProcessing
}
