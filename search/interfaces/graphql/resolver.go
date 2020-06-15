package graphql

import (
	"context"
	"flamingo.me/flamingo-commerce/v3/search/domain"
	"flamingo.me/flamingo-commerce/v3/search/interfaces/graphql/dto"
)

// CommerceSearchQueryResolver is a commerce search query resolver
type CommerceSearchQueryResolver struct{}

// SortOptions remaps search meta options to graphql structure
func (r *CommerceSearchQueryResolver) SortOptions(ctx context.Context, searchMeta *domain.SearchMeta) ([]*dto.CommerceSearchSortOption, error) {
	var options = make([]*dto.CommerceSearchSortOption, 0)
	for _, option := range searchMeta.SortOptions {
		if option.Asc != "" {
			options = append(options, &dto.CommerceSearchSortOption{
				Label:    option.Label + " (asc)",
				Field:    option.Asc,
				Selected: option.SelectedAsc,
			})
		}
		if option.Desc != "" {
			options = append(options, &dto.CommerceSearchSortOption{
				Label:    option.Label + " (desc)",
				Field:    option.Desc,
				Selected: option.SelectedDesc,
			})
		}
	}

	return options, nil
}
