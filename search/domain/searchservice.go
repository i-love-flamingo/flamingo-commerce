package domain

import (
	"net/url"

	"go.aoe.com/flamingo/framework/web"
)

// LegacySearchService interface
type LegacySearchService interface {
	Search(ctx web.Context, query url.Values) (*SearchResult, error)
}
