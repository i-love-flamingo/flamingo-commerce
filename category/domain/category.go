package domain

import (
	"fmt"
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
	Attributes map[string]Attribute //@name CategoryAttributes

	// Attribute instance representation
	Attribute struct {
		Code   string
		Label  string
		Values []AttributeValue
	} //@name CategoryAttribute

	//AttributeValue represents the value that a Attribute can have
	AttributeValue struct {
		Label    string
		RawValue interface{}
	}

	// Promotion defines promotion for a category
	Promotion struct {
		LinkType   string
		LinkTarget string
		Media      Medias
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

// Get by key
func (a Attributes) Get(code string) *Attribute {
	if att, ok := a[code]; ok {
		return &att
	}
	return nil
}

// Has by key
func (a Attributes) Has(code string) bool {
	if _, ok := a[code]; ok {
		return true
	}
	return false
}

// All returns all Attributes as slice
func (a Attributes) All() []Attribute {
	var att []Attribute
	for _, v := range a {
		att = append(att, v)
	}
	return att
}

// ToString returns a concatenated string of all the values of an attribute
func (a Attribute) ToString() string {
	var attValue []string
	for _, v := range a.Values {
		attValue = append(attValue, v.Value())
	}
	return strings.Join(attValue, ",")
}

// Value returns string representation of the RawValue
func (av AttributeValue) Value() string {
	if stringer, ok := av.RawValue.(fmt.Stringer); ok {
		return stringer.String()
	}
	string, _ := av.RawValue.(string)
	return string
}
