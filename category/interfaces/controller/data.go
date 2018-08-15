package controller

import (
	"context"

	"flamingo.me/flamingo-commerce/category/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	// Tree controller for category tree retrieval
	Tree struct {
		categoryService domain.CategoryService
	}

	// Entity controller for category entity retrieval
	Entity struct {
		categoryService domain.CategoryService
	}
)

func (controller *Tree) Inject(service domain.CategoryService) {
	controller.categoryService = service
}

// Data controller for category trees
func (controller *Tree) Data(c context.Context, r *web.Request) interface{} {
	code, _ := r.Param1("code") // no err check, empty code is fine if not set

	categoryRoot, err := controller.categoryService.Tree(c, code)
	_ = err

	return categoryRoot
}

func (controller *Entity) Inject(service domain.CategoryService) {
	controller.categoryService = service
}

// Data controller for category entities
func (controller *Entity) Data(c context.Context, r *web.Request) interface{} {
	code, _ := r.Param1("code") // no err check, empty code is fine if not set

	category, err := controller.categoryService.Get(c, code)
	_ = err

	return category
}
