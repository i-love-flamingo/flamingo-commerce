package graphqlProductDto

import "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	VariantsToVariationSelectionsMapper struct {
		// All variants of a product
		variants []domain.Variant
		// The attributes that are configurable
		variationAttributes []string
		// The currently active variant
		activeVariant *domain.Variant
		// Variants that have all required variation attributes
		variantsWithMatchingAttributes []domain.Variant
		// A group of attributes that have the same code
		attributeGroups map[string]*AttributeGroup
	}

	AttributeGroup struct {
		// Code of the group
		Code string
		// The label of the group
		Label string
		// unique Attributes matching the group code
		Attributes []domain.Attribute
	}
)

// Create a new mapper based on product type
func NewVariantsToVariationSelections(p domain.BasicProduct) []VariationSelection {
	if p.Type() == domain.TypeConfigurableWithActiveVariant {
		configurableWithActiveVariant := p.(domain.ConfigurableProductWithActiveVariant)

		mapper := VariantsToVariationSelectionsMapper{
			variants:            configurableWithActiveVariant.Variants,
			variationAttributes: configurableWithActiveVariant.VariantVariationAttributes,
			activeVariant:       &configurableWithActiveVariant.ActiveVariant,
		}
		return mapper.Map()
	}

	if p.Type() == domain.TypeConfigurable {
		configurable := p.(domain.ConfigurableProduct)

		mapper := VariantsToVariationSelectionsMapper{
			variants:            configurable.Variants,
			variationAttributes: configurable.VariantVariationAttributes,
			activeVariant:       nil,
		}

		return mapper.Map()
	}

	mapper := VariantsToVariationSelectionsMapper{
		variants:            []domain.Variant{},
		variationAttributes: nil,
		activeVariant:       nil,
	}
	return mapper.Map()
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
				group := m.createGroupIfNotExisting(attribute)
				group.addAttributeIfNotExisting(attribute)
			}
		}
	}
}

func (m *VariantsToVariationSelectionsMapper) createGroupIfNotExisting(attribute domain.Attribute) *AttributeGroup {
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
	for _, attribute := range attributeGroup.Attributes {
		var state VariationSelectionOptionState
		var marketPlaceCode string

		attribute := domain.Attribute{
			Code:  attributeGroup.Code,
			Label: attribute.Label,
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
		if m.VariantMatchesAllAttributes(variant, attributes) {
			return &variant
		}
	}
	return nil
}

// VariantMatchesAllAttributes returns true, if a variant has all given attributes
func (m *VariantsToVariationSelectionsMapper) VariantMatchesAllAttributes(variant domain.Variant, attributes []domain.Attribute) bool {
	for _, attribute := range attributes {
		currentAttribute := variant.Attribute(attribute.Code)
		if currentAttribute.Label != attribute.Label {
			return false
		}
	}
	return true
}

func (m *VariantsToVariationSelectionsMapper) getActiveAttributesWithOverwrite(attributeOverwrite domain.Attribute) []domain.Attribute {
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

func NewAttributeGroup(a domain.Attribute) *AttributeGroup {
	return &AttributeGroup{
		Code:            a.Code,
		Label:           a.CodeLabel,
		Attributes:      nil,
	}
}

func (ag *AttributeGroup) addAttribute(attribute domain.Attribute) {
	ag.Attributes = append(ag.Attributes, attribute)
}

func (ag *AttributeGroup) hasAttribute(attribute domain.Attribute) bool {
	for _, currentAttribute := range ag.Attributes {
		if currentAttribute.Label == attribute.Label {
			return true
		}
	}
	return false
}


func (ag *AttributeGroup) addAttributeIfNotExisting(attribute domain.Attribute) {
	if !ag.hasAttribute(attribute) {
		ag.addAttribute(attribute)
	}
}

