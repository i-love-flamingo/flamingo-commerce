package interfaces

import (
	"errors"
)

var (
	ErrCategoryNotFound = errors.New("category_not_found")
	ErrCategoryGeneral  = errors.New("category_general_error")
)
