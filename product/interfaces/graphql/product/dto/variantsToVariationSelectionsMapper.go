package graphqlProductDto

import "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	VariantsToVariationSelectionsMapper struct {
		// All variants of a product
		variants []domain.Variant
		// The attributes that are configurable
		variationAttributes []string
		// The currently active variant
		activeVariant                  *domain.Variant
		// Variants that have all required variation attributes
		variantsWithMatchingAttributes []domain.Variant
		// A group of attributes that have the same code
		attributeGroups map[string]*AttributeGroup
	}

	AttributeGroup struct {
		// Code of the group
		Code            string
		// The label of the group
		Label           string
		// Attributes order containing the Attribute.Label
		AttributesOrder []string
		// A map of Attribute.Label to the attribute
		Attributes      map[string]domain.Attribute
	}
)

// Create a new mapper based on product type
func NewVariantsToVariationSelectionsMapper(p domain.BasicProduct) VariantsToVariationSelectionsMapper {
	if p.Type() == domain.TypeConfigurableWithActiveVariant {
		configurableWithActiveVariant := p.(domain.ConfigurableProductWithActiveVariant)

		return VariantsToVariationSelectionsMapper{
			variants:            configurableWithActiveVariant.Variants,
			variationAttributes: configurableWithActiveVariant.VariantVariationAttributes,
			activeVariant:       &configurableWithActiveVariant.ActiveVariant,
		}
	}

	if p.Type() == domain.TypeConfigurable {
		configurable := p.(domain.ConfigurableProduct)

		return VariantsToVariationSelectionsMapper{
			variants:            configurable.Variants,
			variationAttributes: configurable.VariantVariationAttributes,
			activeVariant:       nil,
		}
	}

	return VariantsToVariationSelectionsMapper{
		variants:            []domain.Variant{},
		variationAttributes: nil,
		activeVariant:       nil,
	}
}

func (m *VariantsToVariationSelectionsMapper) Map() []VariationSelection {
	m.pickVariantsWithMatchingAttributes()
	m.unsetActiveVariantIfInvalid()
	m.groupAttributes()
	return m.buildVariationSelections()
}

func (m *VariantsToVariationSelectionsMapper) pickVariantsWithMatchingAttributes() {
	for _, variant := range m.variants {
		if variant.HasAllAttributes(m.variationAttributes) {
			m.variantsWithMatchingAttributes = append(m.variantsWithMatchingAttributes, variant)
		}
	}
}

func (m *VariantsToVariationSelectionsMapper) groupAttributes() {
	if m.hasVariantsWithMatchingAttributes() {
		m.attributeGroups = make(map[string]*AttributeGroup)

		for _, variant := range m.variantsWithMatchingAttributes {
			for _, attributeCode := range m.variationAttributes {
				attribute := variant.Attribute(attributeCode)
				group := m.ensureGroupExists(attribute)
				group.addAttribute(attribute)
			}
		}
	}
}

func (m *VariantsToVariationSelectionsMapper) ensureGroupExists(attribute domain.Attribute) *AttributeGroup {
	if _, ok := m.attributeGroups[attribute.Code]; !ok {
		m.attributeGroups[attribute.Code] = NewAttributeGroup(attribute)
	}

	return m.attributeGroups[attribute.Code]
}

func (m *VariantsToVariationSelectionsMapper) unsetActiveVariantIfInvalid() {
	if m.activeVariant != nil {
		for _, variant := range m.variantsWithMatchingAttributes {
			if variant.MarketPlaceCode == m.activeVariant.MarketPlaceCode {
				return
			}
		}
		m.activeVariant = nil
	}

}

func (m *VariantsToVariationSelectionsMapper) hasActiveVariant() bool {
	return m.activeVariant != nil
}

func (m *VariantsToVariationSelectionsMapper) buildVariationSelections() []VariationSelection {
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

func (m *VariantsToVariationSelectionsMapper) hasVariantsWithMatchingAttributes() bool {
	return m.variantsWithMatchingAttributes != nil
}

func (m *VariantsToVariationSelectionsMapper) buildVariationSelectionOptions(attributeGroup *AttributeGroup) []VariationSelectionOption {
	var options []VariationSelectionOption
	for _, key := range attributeGroup.AttributesOrder {
		var state VariationSelectionOptionState
		var marketPlaceCode string

		attribute := domain.Attribute{
			Code:  attributeGroup.Code,
			Label: attributeGroup.Attributes[key].Label,
		}

		if m.hasActiveVariant() {
			mergedAttributes := m.getActiveAttributesWithOverwrite(attribute)
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
				fallbackVariant := m.findMatchingVariant([]domain.Attribute{attribute})

				if fallbackVariant != nil {
					state = VariationSelectionOptionStateNoMatch
					marketPlaceCode = fallbackVariant.MarketPlaceCode
				}
			}
		} else {
			fallbackVariant := m.findMatchingVariant([]domain.Attribute{attribute})

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
	}

	return options
}

func (m *VariantsToVariationSelectionsMapper) findMatchingVariant(attributes []domain.Attribute) *domain.Variant {
	for _, variant := range m.variantsWithMatchingAttributes {
		if variant.MatchesAllAttributes(attributes) {
			return &variant
		}
	}
	return nil
}

func (m *VariantsToVariationSelectionsMapper) getActiveAttributesWithOverwrite(overwrite domain.Attribute) []domain.Attribute {
	var attributes []domain.Attribute

	for _, attributeCode := range m.variationAttributes {
		attribute := m.activeVariant.Attribute(attributeCode)
		resultingAttribute := domain.Attribute{
			Code:  attribute.Code,
			Label: attribute.Label,
		}

		if overwrite.Code == attribute.Code {
			resultingAttribute.Label = overwrite.Label
		}

		attributes = append(attributes, resultingAttribute)
	}

	return attributes
}

func NewAttributeGroup(a domain.Attribute) *AttributeGroup {
	return &AttributeGroup{
		Code:            a.Code,
		Label:           a.CodeLabel,
		AttributesOrder: nil,
		Attributes:      make(map[string]domain.Attribute),
	}
}

func (ag *AttributeGroup) addAttribute(attribute domain.Attribute) {
	if _, ok := ag.Attributes[attribute.Label]; !ok {
		ag.AttributesOrder = append(ag.AttributesOrder, attribute.Label)
		ag.Attributes[attribute.Label] = attribute
	}
}


