package domain

type (
	Category struct {
		Code       string
		Name       string
		Categories []*Category
	}
)
