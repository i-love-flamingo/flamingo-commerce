package domain

type (

	// Category domain model
	Category interface {
		Code() string
		Name() string
		Path() string
		Promoted() bool
		Active() bool
		CategoryType() string
		Media() Medias
		Attributes() Attributes
		Attribute(string) interface{}
	}

	// CategoryData defines the default domain category data model
	CategoryData struct {
		CategoryCode       string
		CategoryName       string
		CategoryPath       string
		IsPromoted         bool
		IsActive           bool
		CategoryMedia      Medias
		CategoryTypeCode   string
		CategoryAttributes Attributes
	}

	// Attributes define additional category attributes
	Attributes map[string]interface{}
)

// Category Types
const (
	TypeProduct = "product"
	TypeTeaser  = "teaser"
)

var _ Category = (*CategoryData)(nil)

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

// Promoted gets the category promoted state
func (c CategoryData) Promoted() bool {
	return c.IsPromoted
}

func (c CategoryData) Active() bool {
	return c.IsActive
}

// CategoryType gets the category type code
func (c CategoryData) CategoryType() string {
	return c.CategoryTypeCode
}

// Attributes gets the additional category attributes
func (c CategoryData) Attributes() Attributes {
	return c.CategoryAttributes
}

// Attribute gets an additional category attribute, returns nil if not available
func (c CategoryData) Attribute(code string) interface{} {
	if v, ok := c.CategoryAttributes[code]; ok {
		return v
	}

	return nil
}
