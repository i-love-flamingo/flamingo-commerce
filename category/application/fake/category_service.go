package fake

import (
	"context"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/category/domain"
)

//go:generate go run github.com/go-bindata/go-bindata/v3/go-bindata -nometadata -pkg fake -prefix mock/ mock/

type CategoryService struct {
	testDataFiles map[string]string
	logger        flamingo.Logger
}

var _ domain.CategoryService = new(CategoryService)

func (f *CategoryService) Inject(logger flamingo.Logger, config *struct {
	TestDataFolder string `inject:"config:commerce.category.testDataFolder,optional"`
}) {
	f.logger = logger.WithField(flamingo.LogKeyModule, "category")
	if config != nil {
		if len(config.TestDataFolder) > 0 {
			f.testDataFiles = RegisterTestData(config.TestDataFolder, f.logger)
		}
	}
}

func (f CategoryService) Tree(_ context.Context, activeCategoryCode string) (domain.Tree, error) {
	categoryTree := LoadCategoryTree(f.testDataFiles, f.logger)
	index := findTreeRootIndex(activeCategoryCode, categoryTree)
	if index == -1 {
		return nil, domain.ErrNotFound
	}
	return categoryTree[index], nil
}

func (f CategoryService) Get(_ context.Context, categoryCode string) (domain.Category, error) {
	panic("Implement me")
}

func findTreeRootIndex(categoryCode string, trees []*domain.TreeData) int {
	if len(trees) == 1 {
		return 0
	}
	for i, currentTree := range trees {
		if currentTree.Code() == categoryCode || findTreeRootIndexBySubTree(categoryCode, currentTree.SubTrees(), i) != -1 {
			currentTree.IsActive = true
			return i
		}
	}
	return -1
}

func findTreeRootIndexBySubTree(categoryCode string, trees []domain.Tree, index int) int {
	for _, currentTree := range trees {
		if currentTree.Code() == categoryCode || findTreeRootIndexBySubTree(categoryCode, currentTree.SubTrees(), index) != -1 {
			if data, ok := currentTree.(*domain.TreeData); ok {
				data.IsActive = true
			}
			return index
		}
	}
	return -1
}
