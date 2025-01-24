package interfaces

import (
	"errors"
)

var (
	ErrCategoryNotFound = errors.New("category_not_found")
	ErrGeneralCategory  = errors.New("category_general_error")
)
