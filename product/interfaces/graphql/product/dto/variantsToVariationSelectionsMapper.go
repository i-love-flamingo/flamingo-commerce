package graphqlproductdto

import (
	"sort"

	"flamingo.me/flamingo-commerce/v3/product/domain"
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

	// variantSortingComparer compares the sorting of two variants
	variantSortingComparer struct {
		// Relevant attributes for comparison
		attributeCodes []string
		// A map of ordered labels for an attribute code
		attributesSorting map[string][]string
		// First variant for comparison
		variantA domain.Variant
		// Second variant for comparison
		variantB domain.Variant
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
				attributeCodes:    m.variationAttributes,
				attributesSorting: m.variationAttributesSorting,
				variantA:          m.variantsWithMatchingAttributes[i],
				variantB:          m.variantsWithMatchingAttributes[j],
			}
			return comparer.compare()
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
		option := m.createOptionWithoutActiveVariant(*attribute)

		if m.hasActiveVariant() {
			option = m.createOptionWithActiveVariant(*attribute)
		}

		if option != nil {
			option.UnitCode = attribute.UnitCode
			options = append(options, *option)
		}
	})

	return options
}

func (m *variantsToVariationSelectionsMapper) createOptionWithActiveVariant(attribute domain.Attribute) *VariationSelectionOption {
	mergedAttributes := m.getActiveAttributesWithOverwrite(attribute)
	exactMatchingOption := m.createOption(mergedAttributes, VariationSelectionOption{
		Label: attribute.Label,
		State: VariationSelectionOptionStateMatch,
	})

	if exactMatchingOption != nil {
		if exactMatchingOption.Variant.MarketPlaceCode() == m.activeVariant.MarketPlaceCode {
			exactMatchingOption.State = VariationSelectionOptionStateActive
		}

		return exactMatchingOption
	}

	return m.createOption([]domain.Attribute{attribute}, VariationSelectionOption{
		Label: attribute.Label,
		State: VariationSelectionOptionStateNoMatch,
	})
}

func (m *variantsToVariationSelectionsMapper) createOptionWithoutActiveVariant(attribute domain.Attribute) *VariationSelectionOption {
	return m.createOption([]domain.Attribute{attribute}, VariationSelectionOption{
		Label: attribute.Label,
		State: VariationSelectionOptionStateMatch,
	})
}

func (m *variantsToVariationSelectionsMapper) createOption(attributes []domain.Attribute, props VariationSelectionOption) *VariationSelectionOption {
	fallbackVariant := m.findMatchingVariant(attributes)

	if fallbackVariant != nil {
		return &VariationSelectionOption{
			Variant: VariationSelectionOptionVariant{
				variant: *fallbackVariant,
			},
			Label: props.Label,
			State: props.State,
		}
	}

	return nil
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
	if value, ok := attribute.RawValue.(string); ok {
		ag.Attributes[value] = attribute
	}
}

func (ag *attributeGroup) eachAttributeInOrder(callback func(*domain.Attribute)) {
	for _, attributeValue := range ag.AttributesSorting {
		attribute := ag.getAttributeByValue(attributeValue)
		if attribute != nil {
			callback(attribute)
		}
	}
}

func (ag *attributeGroup) hasAttribute(attribute domain.Attribute) bool {
	value, ok := attribute.RawValue.(string)
	if !ok {
		return false
	}

	_, found := ag.Attributes[value]
	return found
}

func (ag *attributeGroup) addAttributeIfNotExisting(attribute domain.Attribute) {
	if !ag.hasAttribute(attribute) {
		ag.addAttribute(attribute)
	}
}

func (ag *attributeGroup) getAttributeByValue(value string) *domain.Attribute {
	if attribute, ok := ag.Attributes[value]; ok {
		return &attribute
	}
	return nil
}

func (c *variantSortingComparer) compare() bool {
	for _, attributeCode := range c.attributeCodes {
		indexA := c.getSortingIndex(attributeCode, c.variantA)
		indexB := c.getSortingIndex(attributeCode, c.variantB)

		if indexA == indexB {
			continue
		}
		return indexA < indexB
	}
	return false
}

func (c *variantSortingComparer) getSortingIndex(code string, variant domain.Variant) int {
	sortedValues := c.attributesSorting[code]
	for index, value := range sortedValues {
		if variant.Attribute(code).RawValue == value {
			return index
		}
	}
	return -1 // we should not get here, our variants have all required attributes
}
