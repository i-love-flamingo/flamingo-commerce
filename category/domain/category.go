package domain

type (
	// Category domain model
	Category interface {
		Code() string
		Name() string
		Categories() []Category
	}
)
