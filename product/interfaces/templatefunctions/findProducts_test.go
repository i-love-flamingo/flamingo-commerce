package templatefunctions

import (
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

// Test different combinations of filter processing
func TestFilterProcessingIsAllowed(t *testing.T) {

	// Case 1: only blacklist filled, everything not on the blacklist should be fine
	filterConstrains := make(map[string]string)
	filterConstrains["blackList"] = "blacklisted"
	filterConstrains["whiteList"] = ""
	filterProcessing := buildFilterProcessing(filterConstrains)

	assert.True(t, filterProcessing.isAllowed("allowedKey"))
	assert.False(t, filterProcessing.isAllowed("blacklisted"))

	// Case 2: only whitelist filled, everything noy on the whitelist should be bad
	filterConstrains["blackList"] = ""
	filterConstrains["whiteList"] = "allowedKey"
	filterProcessing = buildFilterProcessing(filterConstrains)

	assert.True(t, filterProcessing.isAllowed("allowedKey"))
	assert.False(t, filterProcessing.isAllowed("blacklisted"))

}

func buildFilterProcessing(filterConstrains map[string]string) filterProcessing {
	requestURL := url.URL{}
	webRequest := web.CreateRequest(&http.Request{
		URL: &requestURL,
	}, nil)
	namespace := ""

	filterProcessing := newFilterProcessing(webRequest, namespace, nil, nil, filterConstrains)
	return filterProcessing
}
