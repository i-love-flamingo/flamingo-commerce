package fake_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/category/application/fake"
	"flamingo.me/flamingo-commerce/v3/category/domain"
)

func TestCategoryService_Tree(t *testing.T) {
	service := &fake.CategoryService{}
	service.Inject(flamingo.NullLogger{}, nil)

	t.Run("category not found", func(t *testing.T) {
		_, err := service.Tree(context.Background(), "not-found")
		if assert.NotNil(t, err) {
			assert.Equal(t, err, domain.ErrNotFound)
		}
	})
	t.Run("category electronics", func(t *testing.T) {
		tree, err := service.Tree(context.Background(), "electronics")
		if assert.Nil(t, err) {
			assert.Equal(t, "electronics", tree.Code())
			assert.Equal(t, true, tree.Active())
		}
	})
	t.Run("category headphones", func(t *testing.T) {
		tree, err := service.Tree(context.Background(), "headphones")
		if assert.Nil(t, err) {
			assert.Equal(t, "electronics", tree.Code())
			assert.Equal(t, true, tree.Active())

			if assert.Equal(t, len(tree.SubTrees()), 3) {
				headphonesTree := tree.SubTrees()[1]
				assert.Equal(t, true, headphonesTree.Active())
			}
		}
	})
	t.Run("category headphone_accessories", func(t *testing.T) {
		tree, err := service.Tree(context.Background(), "headphone_accessories")
		if assert.Nil(t, err) {
			assert.Equal(t, "electronics", tree.Code())
			assert.Equal(t, true, tree.Active())

			if assert.Equal(t, 3, len(tree.SubTrees())) {
				headphonesTree := tree.SubTrees()[1]
				assert.Equal(t, true, headphonesTree.Active())
				if assert.Equal(t, 1, len(headphonesTree.SubTrees())) {
					accessoriesTree := headphonesTree.SubTrees()[0]
					assert.Equal(t, true, accessoriesTree.Active())
				}
			}
		}
	})
	t.Run("category clothing", func(t *testing.T) {
		tree, err := service.Tree(context.Background(), "clothing")
		if assert.Nil(t, err) {
			assert.Equal(t, "clothing", tree.Code())
			assert.Equal(t, true, tree.Active())
		}
	})
	t.Run("category jumpsuits", func(t *testing.T) {
		tree, err := service.Tree(context.Background(), "jumpsuits")
		if assert.Nil(t, err) {
			assert.Equal(t, "clothing", tree.Code())
			assert.Equal(t, true, tree.Active())

			if assert.Equal(t, 1, len(tree.SubTrees())) {
				jumpsuitsTree := tree.SubTrees()[0]
				assert.Equal(t, true, jumpsuitsTree.Active())
			}
		}
	})
}
