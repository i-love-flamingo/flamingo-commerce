package interfaces

import (
	"errors"
)

var (
	ErrProductNotFound = errors.New("product_not_found")
	ErrProductGeneral  = errors.New("product_general_error")
)
