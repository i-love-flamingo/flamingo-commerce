package domain

type (
	// Category domain model
	Category interface {
		Code() string
		Name() string
		Path() string
		Categories() []Category
		Active() bool
		Promoted() bool
		CategoryType() string
		Media() Medias
	}

	// CategoryData defines the default domain category data model
	CategoryData struct {
		CategoryCode     string
		CategoryName     string
		CategoryPath     string
		Children         []*CategoryData
		IsActive         bool
		IsPromoted       bool
		CategoryMedia    Medias
		CategoryTypeCode string
	}
)

// Category Types
const (
	TypeProduct = "product"
	TypeTeaser  = "teaser"
)

// GetActive returns the active category of the category tree described by a category with children
// retuns nil if no active category could be found
func GetActive(c Category) Category {
	if c == nil {
		return nil
	}
	for _, sub := range c.Categories() {
		if active := GetActive(sub); active != nil {
			return active
		}
	}
	if c.Active() {
		return c
	}
	return nil
}

var _ Category = (*CategoryData)(nil)

// Active gets the category active state
func (c CategoryData) Active() bool {
	return c.IsActive
}

// Code gets the category code
func (c CategoryData) Code() string {
	return c.CategoryCode
}

// Media gets the category media
func (c CategoryData) Media() Medias {
	return c.CategoryMedia
}

// Name gets the category name
func (c CategoryData) Name() string {
	return c.CategoryName
}

// Path gets the category path
func (c CategoryData) Path() string {
	return c.CategoryPath
}

// Categories gets the child categories
func (c CategoryData) Categories() []Category {
	result := make([]Category, len(c.Children))
	for i, child := range c.Children {
		result[i] = Category(child)
	}

	return result
}

// Promoted gets the category promoted state
func (c CategoryData) Promoted() bool {
	return c.IsPromoted
}

// CategoryType gets the category type code
func (c CategoryData) CategoryType() string {
	return c.CategoryTypeCode
}
