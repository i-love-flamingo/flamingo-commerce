package controller

import (
	"go.aoe.com/flamingo/core/category/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// Tree controller for category tree retrieval
	Tree struct {
		CategoryService domain.CategoryService `inject:""`
	}

	// Entity controller for category entity retrieval
	Entity struct {
		CategoryService domain.CategoryService `inject:""`
	}
)

// Data controller for category trees
func (controller *Tree) Data(c web.Context) interface{} {
	code, _ := c.Param1("code") // no err check, empty code is fine if not set

	categoryRoot, err := controller.CategoryService.Get(c, code)
	_ = err

	return categoryRoot
}

// Data controller for category entities
func (controller *Entity) Data(c web.Context) interface{} {
	code, _ := c.Param1("code") // no err check, empty code is fine if not set

	category, err := controller.CategoryService.Get(c, code)
	_ = err

	return category
}
