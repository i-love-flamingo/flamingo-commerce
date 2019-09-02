package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/category/domain"
)

// CommerceCategoryQueryResolver resolves graphql category queries
type CommerceCategoryQueryResolver struct {
	categoryService domain.CategoryService
}

// Inject dependencies
func (r *CommerceCategoryQueryResolver) Inject(service domain.CategoryService) {
	r.categoryService = service
}

// CommerceCategoryTree returns a Tree with the given activeCategoryCode from categoryService
func (r *CommerceCategoryQueryResolver) CommerceCategoryTree(ctx context.Context, activeCategoryCode string) (domain.Tree, error) {
	return r.categoryService.Tree(ctx, activeCategoryCode)
}
