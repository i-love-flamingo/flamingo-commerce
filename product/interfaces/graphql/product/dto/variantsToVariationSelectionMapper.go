package graphqlProductDto

import "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	VariantsToVariationSelectionMapper struct {
		// All variants of a product
		Variants []domain.Variant
		// The attributes that are configurable
		VariationAttributes []string
		// The currently active variant
		ActiveVariantMarketPlaceCode string
		matchingVariants             []domain.Variant
		activeVariant                *domain.Variant
		attributeGroupsByCode        map[string]*AttributeGroup
	}

	AttributeGroup struct {
		Code            string
		Label           string
		AttributesOrder []string
		Attributes      map[string]*Attribute
	}

	Attribute struct {
		Label string
	}

	PopulatedAttribute struct {
		Code  string
		Label string
	}
)

func New(p domain.BasicProduct) VariantsToVariationSelectionMapper {
	if p.Type() == domain.TypeConfigurableWithActiveVariant {
		configurableWithActiveVariant := p.(domain.ConfigurableProductWithActiveVariant)

		return VariantsToVariationSelectionMapper{
			Variants:                     configurableWithActiveVariant.Variants,
			VariationAttributes:          configurableWithActiveVariant.VariantVariationAttributes,
			ActiveVariantMarketPlaceCode: configurableWithActiveVariant.ActiveVariant.MarketPlaceCode,
		}
	}

	if p.Type() == domain.TypeConfigurable {
		configurable := p.(domain.ConfigurableProduct)

		return VariantsToVariationSelectionMapper{
			Variants:                     configurable.Variants,
			VariationAttributes:          configurable.VariantVariationAttributes,
			ActiveVariantMarketPlaceCode: "",
		}
	}

	return VariantsToVariationSelectionMapper{
		Variants:                     []domain.Variant{},
		VariationAttributes:          nil,
		ActiveVariantMarketPlaceCode: "",
	}
}

func (m *VariantsToVariationSelectionMapper) Map() []VariationSelection {
	m.findMatchingVariants()
	m.findActiveVariant()
	m.groupAttributes()

	return m.buildVariationSelections()
}

func (m *VariantsToVariationSelectionMapper) findMatchingVariants() {
	for _, variant := range m.Variants {
		if m.variantHasAllRequiredAttributes(variant) {
			m.matchingVariants = append(m.matchingVariants, variant)
		}
	}
}

func (m *VariantsToVariationSelectionMapper) groupAttributes() {
	m.attributeGroupsByCode = make(map[string]*AttributeGroup)

	for _, variant := range m.matchingVariants {
		for _, attributeCode := range m.VariationAttributes {
			attribute := variant.Attribute(attributeCode)

			if _, ok := m.attributeGroupsByCode[attributeCode]; !ok {
				m.attributeGroupsByCode[attributeCode] = &AttributeGroup{
					Code:            attributeCode,
					Label:           attribute.CodeLabel,
					AttributesOrder: nil,
					Attributes:      make(map[string]*Attribute),
				}
			}

			if _, ok := m.attributeGroupsByCode[attributeCode].Attributes[attribute.Label]; !ok {
				m.attributeGroupsByCode[attributeCode].AttributesOrder = append(m.attributeGroupsByCode[attributeCode].AttributesOrder, attribute.Label)
				m.attributeGroupsByCode[attributeCode].Attributes[attribute.Label] = &Attribute{
					Label: attribute.Label,
				}
			}
		}
	}
}

func (m *VariantsToVariationSelectionMapper) findActiveVariant() {
	if m.ActiveVariantMarketPlaceCode != "" {
		for _, variant := range m.Variants {
			if variant.MarketPlaceCode == m.ActiveVariantMarketPlaceCode {
				m.activeVariant = &variant
				return
			}
		}
	}
}

func (m *VariantsToVariationSelectionMapper) hasActiveVariant() bool {
	return m.activeVariant != nil
}

func (m *VariantsToVariationSelectionMapper) variantHasAllRequiredAttributes(variant domain.Variant) bool {
	for _, attributeCode := range m.VariationAttributes {
		if !variant.HasAttribute(attributeCode) {
			return false
		}
	}
	return true
}

func (m *VariantsToVariationSelectionMapper) buildVariationSelections() []VariationSelection {
	var variationSelections []VariationSelection
	for _, attributeCode := range m.VariationAttributes {
		attribute := m.attributeGroupsByCode[attributeCode]

		variationSelections = append(variationSelections, VariationSelection{
			Code:    attribute.Code,
			Label:   attribute.Label,
			Options: m.buildVariationSelectionOptions(attribute),
		})
	}

	return variationSelections
}

func (m *VariantsToVariationSelectionMapper) buildVariationSelectionOptions(attributeGroup *AttributeGroup) []VariationSelectionOption {
	var options []VariationSelectionOption
	for _, key := range attributeGroup.AttributesOrder {
		attribute := attributeGroup.Attributes[key]
		populatedAttribute := PopulatedAttribute{
			Code:  attributeGroup.Code,
			Label: attribute.Label,
		}

		var state VariationSelectionOptionState
		var marketPlaceCode string

		if m.hasActiveVariant()  {
			mergedAttributes := m.getPopulatedAttributesFromActiveVariantWithOverwrite(populatedAttribute);
			exactMatchingVariant := m.findExactMatchingVariant(mergedAttributes)

			if exactMatchingVariant != nil {
				if exactMatchingVariant.MarketPlaceCode == m.activeVariant.MarketPlaceCode {
					state = VariationSelectionOptionStateActive
				} else {
					state = VariationSelectionOptionStateMatch
				}

				marketPlaceCode = exactMatchingVariant.MarketPlaceCode
			} else {
				state = VariationSelectionOptionStateNoMatch
				fallbackVariant := m.findExactMatchingVariant([]PopulatedAttribute{populatedAttribute});

				if fallbackVariant != nil {
					marketPlaceCode = fallbackVariant.MarketPlaceCode
				}
			}
		} else {
			state = VariationSelectionOptionStateMatch
			fallbackVariant := m.findExactMatchingVariant([]PopulatedAttribute{populatedAttribute});

			if fallbackVariant != nil {
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

func (m *VariantsToVariationSelectionMapper) findExactMatchingVariant(populatedAttributes []PopulatedAttribute) *domain.Variant {
	for _, variant := range m.matchingVariants {
		if m.variantHasAllAttributes(variant, populatedAttributes) {
			return &variant
		}
	}

	return nil
}


func (m *VariantsToVariationSelectionMapper) getPopulatedAttributesFromActiveVariantWithOverwrite(overwrite PopulatedAttribute) []PopulatedAttribute {
	var populatedAttributes []PopulatedAttribute

	for _, attributeCode := range m.VariationAttributes {
		attribute := m.activeVariant.Attribute(attributeCode)

		if overwrite.Code == attribute.Code {
			populatedAttributes = append(populatedAttributes, PopulatedAttribute{
				Code:  attribute.Code,
				Label: overwrite.Label,
			})
		} else {
			populatedAttributes = append(populatedAttributes, PopulatedAttribute{
				Code:  attribute.Code,
				Label: attribute.Label,
			})
		}
	}

	return populatedAttributes
}

func (m *VariantsToVariationSelectionMapper) variantHasAllAttributes(variant domain.Variant, populatedAttributes []PopulatedAttribute) bool {
	for _, populatedAttribute := range populatedAttributes {
		attribute := variant.Attribute(populatedAttribute.Code)
		if attribute.Label != populatedAttribute.Label {
			return false
		}
	}

	return true
}

