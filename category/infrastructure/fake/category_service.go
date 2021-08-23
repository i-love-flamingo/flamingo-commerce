package fake

import (
	"context"

	"flamingo.me/flamingo/v3/framework/flamingo"

	"flamingo.me/flamingo-commerce/v3/category/domain"
)

// CategoryService returns category test data
type CategoryService struct {
	testDataFiles map[string]string
	logger        flamingo.Logger
}

var _ domain.CategoryService = new(CategoryService)

// Inject dependencies
func (f *CategoryService) Inject(logger flamingo.Logger, config *struct {
	TestDataFolder string `inject:"config:commerce.category.fakeService.testDataFolder,optional"`
}) {
	f.logger = logger.WithField(flamingo.LogKeyModule, "category").WithField(flamingo.LogKeyCategory, "fakeService")
	if config != nil {
		if len(config.TestDataFolder) > 0 {
			f.testDataFiles = RegisterTestData(config.TestDataFolder, f.logger)
		}
	}
}

// Tree returns the tree the given category belongs to
func (f CategoryService) Tree(_ context.Context, activeCategoryCode string) (domain.Tree, error) {
	categoryTree := LoadCategoryTree(f.testDataFiles, f.logger)
	if len(activeCategoryCode) == 0 {
		return &domain.TreeData{
			IsActive:     true,
			SubTreesData: categoryTree,
		}, nil
	}
	index := findTreeRootIndex(activeCategoryCode, categoryTree)
	if index == -1 {
		return nil, domain.ErrNotFound
	}
	return categoryTree[index], nil
}

// Get the category of the given category code
func (f CategoryService) Get(_ context.Context, categoryCode string) (domain.Category, error) {
	category := LoadCategory(categoryCode, f.testDataFiles, f.logger)
	if category == nil {
		treeCategory := f.getCategoryByTree(categoryCode)
		if treeCategory == nil {
			return nil, domain.ErrNotFound
		}
		return &domain.CategoryData{
			CategoryCode: treeCategory.Code(),
			CategoryName: treeCategory.Name(),
			CategoryPath: treeCategory.Path(),
			IsActive:     treeCategory.Active(),
		}, nil
	}
	return category, nil
}

func (f CategoryService) getCategoryByTree(categoryCode string) domain.Tree {
	categoryTree := LoadCategoryTree(f.testDataFiles, f.logger)
	tree := make([]domain.Tree, len(categoryTree))
	for i, data := range categoryTree {
		tree[i] = data
	}
	return findCategoryInTreeByCode(categoryCode, tree)
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

func findCategoryInTreeByCode(categoryCode string, trees []domain.Tree) domain.Tree {
	for _, currentTree := range trees {
		if currentTree.Code() == categoryCode {
			return currentTree
		}
		subTree := findCategoryInTreeByCode(categoryCode, currentTree.SubTrees())
		if subTree != nil {
			return subTree
		}
	}
	return nil
}
