package interfaces

import (
	"errors"
)

var (
	ErrProductNotFound = errors.New("product_not_found")
	ErrGeneralProduct  = errors.New("product_general_error")
)
