package interfaces

import (
	"errors"
)

var (
	ErrCategoryNotFound = errors.New("category_not_found")
	ErrCategoryGeneral  = errors.New("category_general_error")
	// ErrCategoryPageOutOfRange is returned when the requested page is bigger than the number of result pages
	ErrCategoryPageOutOfRange = errors.New("category_page_out_of_range")
)
