package graphqlproductdto

import (
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"sort"
)

type (
	variantsToVariationSelectionsMapper struct {
		// All variants of a product
		variants []domain.Variant
		// The attributes that are configurable
		variationAttributes []string
		// The preferred sorting of attributes
		variationAttributesSorting map[string][]string
		// The currently active variant
		activeVariant *domain.Variant
		// Variants that have all required variation attributes
		variantsWithMatchingAttributes []domain.Variant
		// A group of attributes that have the same code
		attributeGroups map[string]*attributeGroup
	}

	attributeGroup struct {
		// Code of the group
		Code string
		// The label of the group
		Label string
		// unique Attributes matching the group code
		Attributes map[string]domain.Attribute
		// The order of the attributes
		AttributesSorting []string
	}

	variantSortingComparer struct {
		attributes        []string
		attributesSorting map[string][]string
		a                 domain.Variant
		b                 domain.Variant
	}
)

// NewVariantsToVariationSelections Converts a product to variation selections
func NewVariantsToVariationSelections(p domain.BasicProduct) []VariationSelection {
	if p.Type() == domain.TypeConfigurableWithActiveVariant {
		configurableWithActiveVariant := p.(domain.ConfigurableProductWithActiveVariant)

		mapper := variantsToVariationSelectionsMapper{
			variants:                   configurableWithActiveVariant.Variants,
			variationAttributes:        configurableWithActiveVariant.VariantVariationAttributes,
			variationAttributesSorting: configurableWithActiveVariant.VariantVariationAttributesSorting,
			activeVariant:              &configurableWithActiveVariant.ActiveVariant,
		}
		return mapper.Map()
	}

	if p.Type() == domain.TypeConfigurable {
		configurable := p.(domain.ConfigurableProduct)

		mapper := variantsToVariationSelectionsMapper{
			variants:                   configurable.Variants,
			variationAttributes:        configurable.VariantVariationAttributes,
			variationAttributesSorting: configurable.VariantVariationAttributesSorting,
			activeVariant:              nil,
		}

		return mapper.Map()
	}

	return []VariationSelection{}
}

func (m *variantsToVariationSelectionsMapper) Map() []VariationSelection {
	m.pickVariantsWithMatchingAttributes()
	m.unsetActiveVariantIfInvalid()
	m.sortVariants()
	m.groupAttributes()
	return m.buildVariationSelections()
}

func (m *variantsToVariationSelectionsMapper) pickVariantsWithMatchingAttributes() {
	for _, variant := range m.variants {
		if variant.HasAllAttributes(m.variationAttributes) {
			m.variantsWithMatchingAttributes = append(m.variantsWithMatchingAttributes, variant)
		}
	}
}

func (m *variantsToVariationSelectionsMapper) sortVariants() {
	if m.hasVariantsWithMatchingAttributes() {
		sort.Slice(m.variantsWithMatchingAttributes, func(i, j int) bool {
			comparer := variantSortingComparer{
				attributes:        m.variationAttributes,
				attributesSorting: m.variationAttributesSorting,
				a:                 m.variantsWithMatchingAttributes[i],
				b:                 m.variantsWithMatchingAttributes[j],
			}
			return comparer.compareSortingIndex()
		})
	}
}

func (m *variantsToVariationSelectionsMapper) groupAttributes() {
	if m.hasVariantsWithMatchingAttributes() {
		m.attributeGroups = make(map[string]*attributeGroup)

		for _, variant := range m.variantsWithMatchingAttributes {
			for _, attributeCode := range m.variationAttributes {
				attribute := variant.Attribute(attributeCode)
				group := m.createGroupIfNotExisting(attribute)
				group.addAttributeIfNotExisting(attribute)
			}
		}
	}
}

func (m *variantsToVariationSelectionsMapper) createGroupIfNotExisting(attribute domain.Attribute) *attributeGroup {
	if _, ok := m.attributeGroups[attribute.Code]; !ok {
		m.attributeGroups[attribute.Code] = newAttributeGroup(attribute, m.variationAttributesSorting[attribute.Code])
	}

	return m.attributeGroups[attribute.Code]
}

func (m *variantsToVariationSelectionsMapper) unsetActiveVariantIfInvalid() {
	if m.activeVariant != nil {
		for _, variant := range m.variantsWithMatchingAttributes {
			if variant.MarketPlaceCode == m.activeVariant.MarketPlaceCode {
				return
			}
		}
		m.activeVariant = nil
	}
}

func (m *variantsToVariationSelectionsMapper) hasActiveVariant() bool {
	return m.activeVariant != nil
}

func (m *variantsToVariationSelectionsMapper) buildVariationSelections() []VariationSelection {
	var variationSelections []VariationSelection

	if m.hasVariantsWithMatchingAttributes() {
		for _, attributeCode := range m.variationAttributes {
			attributeGroup := m.attributeGroups[attributeCode]

			variationSelections = append(variationSelections, VariationSelection{
				Code:    attributeGroup.Code,
				Label:   attributeGroup.Label,
				Options: m.buildVariationSelectionOptions(attributeGroup),
			})
		}
	}

	return variationSelections
}

func (m *variantsToVariationSelectionsMapper) hasVariantsWithMatchingAttributes() bool {
	return m.variantsWithMatchingAttributes != nil
}

func (m *variantsToVariationSelectionsMapper) buildVariationSelectionOptions(attributeGroup *attributeGroup) []VariationSelectionOption {
	var options []VariationSelectionOption

	attributeGroup.eachAttributeInOrder(func(attribute *domain.Attribute) {
		var state VariationSelectionOptionState
		var marketPlaceCode string

		if m.hasActiveVariant() {
			mergedAttributes := m.getActiveAttributesWithOverwrite(*attribute)
			exactMatchingVariant := m.findMatchingVariant(mergedAttributes)

			if exactMatchingVariant != nil {
				if exactMatchingVariant.MarketPlaceCode == m.activeVariant.MarketPlaceCode {
					state = VariationSelectionOptionStateActive
					marketPlaceCode = exactMatchingVariant.MarketPlaceCode
				} else {
					state = VariationSelectionOptionStateMatch
					marketPlaceCode = exactMatchingVariant.MarketPlaceCode
				}
			} else {
				fallbackVariant := m.findMatchingVariant([]domain.Attribute{*attribute})

				if fallbackVariant != nil {
					state = VariationSelectionOptionStateNoMatch
					marketPlaceCode = fallbackVariant.MarketPlaceCode
				}
			}
		} else {
			fallbackVariant := m.findMatchingVariant([]domain.Attribute{*attribute})

			if fallbackVariant != nil {
				state = VariationSelectionOptionStateMatch
				marketPlaceCode = fallbackVariant.MarketPlaceCode
			}
		}

		options = append(options, VariationSelectionOption{
			Label:                  attribute.Label,
			State:                  state,
			VariantMarketPlaceCode: marketPlaceCode,
		})
	})

	return options
}

func (m *variantsToVariationSelectionsMapper) findMatchingVariant(attributes []domain.Attribute) *domain.Variant {
	for _, variant := range m.variantsWithMatchingAttributes {
		if m.VariantMatchesAllAttributes(variant, attributes) {
			return &variant
		}
	}
	return nil
}

// VariantMatchesAllAttributes returns true, if a variant has all given attributes
func (m *variantsToVariationSelectionsMapper) VariantMatchesAllAttributes(variant domain.Variant, attributes []domain.Attribute) bool {
	for _, attribute := range attributes {
		currentAttribute := variant.Attribute(attribute.Code)
		if currentAttribute.Label != attribute.Label {
			return false
		}
	}
	return true
}

func (m *variantsToVariationSelectionsMapper) getActiveAttributesWithOverwrite(attributeOverwrite domain.Attribute) []domain.Attribute {
	var attributes []domain.Attribute

	for _, attributeCode := range m.variationAttributes {
		attribute := m.activeVariant.Attribute(attributeCode)
		resultingAttribute := domain.Attribute{
			Code:  attribute.Code,
			Label: attribute.Label,
		}

		if attribute.Code == attributeOverwrite.Code {
			resultingAttribute.Label = attributeOverwrite.Label
		}

		attributes = append(attributes, resultingAttribute)
	}

	return attributes
}

func newAttributeGroup(a domain.Attribute, attributesSorting []string) *attributeGroup {
	return &attributeGroup{
		Code:              a.Code,
		Label:             a.CodeLabel,
		Attributes:        map[string]domain.Attribute{},
		AttributesSorting: attributesSorting,
	}
}

func (ag *attributeGroup) addAttribute(attribute domain.Attribute) {
	ag.Attributes[attribute.Label] = attribute
}

func (ag *attributeGroup) eachAttributeInOrder(callback func(*domain.Attribute)) {
	for _, attributeLabel := range ag.AttributesSorting {
		attribute := ag.getAttributeByLabel(attributeLabel)
		if attribute != nil {
			callback(attribute)
		}
	}
}

func (ag *attributeGroup) hasAttribute(attribute domain.Attribute) bool {
	_, ok := ag.Attributes[attribute.Label]
	return ok
}

func (ag *attributeGroup) addAttributeIfNotExisting(attribute domain.Attribute) {
	if !ag.hasAttribute(attribute) {
		ag.addAttribute(attribute)
	}
}

func (ag *attributeGroup) getAttributeByLabel(label string) *domain.Attribute {
	if attribute, ok := ag.Attributes[label]; ok {
		return &attribute
	}
	return nil
}

func (c *variantSortingComparer) compareSortingIndex() bool {
	for _, attributeCode := range c.attributes {
		indexA := c.getSortingIndex(attributeCode, c.a)
		indexB := c.getSortingIndex(attributeCode, c.b)

		if indexA == indexB {
			continue
		}
		return indexA < indexB
	}
	return false
}

func (c *variantSortingComparer) getSortingIndex(code string, variant domain.Variant) int {
	sorting := c.attributesSorting[code]
	for index, label := range sorting {
		if variant.Attribute(code).Label == label {
			return index
		}
	}
	return -1 // we should not get here, our variants have all required attributes
}
