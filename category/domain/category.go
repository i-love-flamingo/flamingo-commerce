package domain

import (
	"strings"
)

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
		Promotion          Promotion
	}

	// Attributes define additional category attributes
	Attributes map[string]interface{}

	// Promotion defines promotion for a category
	Promotion struct {
		LinkType   string
		LinkTarget string
		Media      Medias
	}

	// AdditionalAttributes - concrete additional category attributes (see searchperience category)
	AdditionalAttributes struct {
		Title            string
		MarketingTitle   string
		ShortDescription string
		Content          string
	}
)

// Category Types
const (
	TypeProduct   = "product"
	TypeTeaser    = "teaser"
	TypePromotion = "promotion"
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

// Active indicator
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

// GetAdditionalAttributes - returns additional attributes
func (c CategoryData) GetAdditionalAttributes() AdditionalAttributes {
	if c.CategoryAttributes != nil {
		return c.CategoryAttributes.mapToAdditionalAttributes()
	}
	return AdditionalAttributes{}
}

// attributeKeys - lists all available keys
func (a Attributes) attributeKeys() []string {
	res := make([]string, len(a))
	i := 0
	for k := range a {
		res[i] = k
		i++
	}
	return res
}

// mapToAdditionalAttributes - maps attributes to AdditionalAttributes struct
func (a Attributes) mapToAdditionalAttributes() AdditionalAttributes {
	additionalAttributes := AdditionalAttributes{}
	attributeKeys := a.attributeKeys()

	for _, key := range attributeKeys {
		switch strings.ToLower(key) {
		case "title":
			if value, ok := a[key].(string); ok {
				additionalAttributes.Title = value
			}
		case "marketingtitle":
			if value, ok := a[key].(string); ok {
				additionalAttributes.MarketingTitle = value
			}
		case "shortdescription":
			if value, ok := a[key].(string); ok {
				additionalAttributes.ShortDescription = value
			}
		case "content":
			if value, ok := a[key].(string); ok {
				additionalAttributes.Content = value
			}
		}
	}

	return additionalAttributes
}
