package models

type (
	ProductAttribute struct {
		ID    string
		Name  string
		Value string
	}

	Product struct {
		ID          string
		Name        string
		Description string
		Price       float64
		Images      []string
	}
)
