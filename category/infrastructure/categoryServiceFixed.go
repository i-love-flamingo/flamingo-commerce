package infrastructure

import (
	"context"
	"sort"

	"flamingo.me/flamingo-commerce/v3/category/domain"

	"flamingo.me/flamingo/v3/framework/config"
)

type (
	//CategoryServiceFixed - a secondary adapter for category service that returns simple categories based from configuration
	CategoryServiceFixed struct {
		mappedTree       *domain.TreeData
		mappedCategories map[string]*domain.CategoryData
	}

	configCat struct {
		Code   string
		Name   string
		Sort   int
		Childs map[string]configCat
	}
)

var (
	_ domain.CategoryService = new(CategoryServiceFixed)
)

// Inject - dingo injector
func (c *CategoryServiceFixed) Inject(config *struct {
	Tree config.Map `inject:"config:commerce.category.categoryServiceFixed.tree"`
}) {
	c.mapConfig(config.Tree)
}

func (c *CategoryServiceFixed) mapConfig(config config.Map) {
	structure := make(map[string]configCat)
	config.MapInto(&structure)

	c.mappedCategories = make(map[string]*domain.CategoryData)

	c.mappedTree = &domain.TreeData{
		CategoryCode: "root",
		SubTreesData: c.getSubTree(structure),
	}

}

func (c *CategoryServiceFixed) getSubTree(subs map[string]configCat) []*domain.TreeData {
	var trees []*domain.TreeData
	var sliceOfCategoryDto []configCat
	for _, cat := range subs {
		sliceOfCategoryDto = append(sliceOfCategoryDto, cat)
	}
	sort.Slice(sliceOfCategoryDto, func(i, j int) bool {
		return sliceOfCategoryDto[i].Sort < sliceOfCategoryDto[j].Sort
	})
	for _, cat := range sliceOfCategoryDto {
		trees = append(trees, &domain.TreeData{
			CategoryCode: cat.Code,
			CategoryName: cat.Name,
			SubTreesData: c.getSubTree(cat.Childs),
		})

		c.mappedCategories[cat.Code] = &domain.CategoryData{
			CategoryCode: cat.Code,
			CategoryName: cat.Name,
		}
	}

	return trees
}

// Tree a category
func (c *CategoryServiceFixed) Tree(ctx context.Context, activeCategoryCode string) (domain.Tree, error) {
	if c.mappedTree == nil {
		return nil, domain.ErrNotFound
	}
	return c.mappedTree, nil
}

// Get a category with more data
func (c *CategoryServiceFixed) Get(ctx context.Context, categoryCode string) (domain.Category, error) {
	if cat, ok := c.mappedCategories[categoryCode]; ok {
		return cat, nil
	}
	return nil, domain.ErrNotFound
}
