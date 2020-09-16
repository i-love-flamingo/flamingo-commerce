package graphqlProductDto

import "flamingo.me/flamingo-commerce/v3/product/domain"

type (
	VariantsToVariationSelectionMapper struct {
		// All variants of a product
		Variants []domain.Variant
		// The attributes that are configurable
		VariationAttributes []string
		// The currently active variant
		ActiveVariantMarketPlaceCode *string
		matchingVariants             []domain.Variant
		activeVariant                *domain.Variant
		attributesByKey              map[string]*Attribute
	}

	Attribute struct {
		Code         string
		Label        string
		ValuesOrder  []string
		ValuesByCode map[string]*AttributeValue
	}

	AttributeValue struct {
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
			ActiveVariantMarketPlaceCode: &configurableWithActiveVariant.ActiveVariant.MarketPlaceCode,
		}
	}

	if p.Type() == domain.TypeConfigurable {
		configurable := p.(domain.ConfigurableProduct)

		return VariantsToVariationSelectionMapper{
			Variants:                     configurable.Variants,
			VariationAttributes:          configurable.VariantVariationAttributes,
			ActiveVariantMarketPlaceCode: nil,
		}
	}

	return VariantsToVariationSelectionMapper{
		Variants:                     []domain.Variant{},
		VariationAttributes:          nil,
		ActiveVariantMarketPlaceCode: nil,
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
		if m.variantHasAllAttributes(variant) {
			m.matchingVariants = append(m.matchingVariants, variant)
		}
	}
}

func (m *VariantsToVariationSelectionMapper) groupAttributes() {
	m.attributesByKey = make(map[string]*Attribute)

	for _, variant := range m.matchingVariants {
		for _, attributeKey := range m.VariationAttributes {
			attribute := variant.Attribute(attributeKey)

			if _, ok := m.attributesByKey[attributeKey]; !ok {
				m.attributesByKey[attributeKey] = &Attribute{
					Code:         attributeKey,
					Label:        attribute.CodeLabel,
					ValuesOrder:  nil,
					ValuesByCode: make(map[string]*AttributeValue),
				}
			}

			if _, ok := m.attributesByKey[attributeKey].ValuesByCode[attribute.Code]; !ok {
				m.attributesByKey[attributeKey].ValuesOrder = append(m.attributesByKey[attributeKey].ValuesOrder, attribute.Code)
				m.attributesByKey[attributeKey].ValuesByCode[attribute.Code] = &AttributeValue{
					Code:  attribute.Code,
					Label: attribute.Label,
				}
			}
		}
	}
}

func (m *VariantsToVariationSelectionMapper) findActiveVariant() {
	if m.ActiveVariantMarketPlaceCode != nil {
		for _, variant := range m.Variants {
			if &variant.MarketPlaceCode == m.ActiveVariantMarketPlaceCode {
				m.activeVariant = &variant
				return
			}
		}
	}
}

func (m *VariantsToVariationSelectionMapper) hasActiveVariant() bool {
	return m.activeVariant != nil
}

//func (m *VariantsToVariationSelectionMapper) hasActiveVariant() bool {
//	m.activeVariantAttributes = make(map[string]string)
//	m.matchingVariants = append(m.matchingVariants, variant)
//	for _, attributeKey := range m.VariationAttributes {
//		m.activeVariantAttributes[attributeKey] = variant.Attribute(attributeKey).Code
//	}
//}


func (m *VariantsToVariationSelectionMapper) variantHasAllAttributes(variant domain.Variant) bool {
	for _, attributeKey := range m.VariationAttributes {
		if !variant.HasAttribute(attributeKey) {
			return false
		}
	}
	return true
}

func (m *VariantsToVariationSelectionMapper) buildVariationSelections() []VariationSelection {
	var variationSelections []VariationSelection
	for _, attributeKey := range m.VariationAttributes {
		attribute := m.attributesByKey[attributeKey]

		variationSelections = append(variationSelections, VariationSelection{
			Code:    attribute.Code,
			Label:   attribute.Label,
			Options: m.buildVariationSelectionOptions(attribute),
		})

	}

	return variationSelections
}

func (m *VariantsToVariationSelectionMapper) buildVariationSelectionOptions(attribute *Attribute) []VariationSelectionOption {
	var options []VariationSelectionOption
	for _, valueCode := range attribute.ValuesOrder {
		value := attribute.ValuesByCode[valueCode]
		options = append(options, VariationSelectionOption{
			Code:                   value.Code,
			Label:                  value.Label,
			State:                  VariationSelectionOptionStateMatch,
			VariantMarketPlaceCode: m.getBestMatchingVariantMarketPlaceCode(attribute, value.Code),
		})
	}

	return options
}

func (m *VariantsToVariationSelectionMapper) getBestMatchingVariantMarketPlaceCode(attribute *Attribute, attributeCode string) string {

	if m.hasActiveVariant() {
		// TODO
	}

	return m.getFirstMatchingVariantMarketPlaceCode(attribute, attributeCode)
}


func (m *VariantsToVariationSelectionMapper) getFirstMatchingVariantMarketPlaceCode(attribute *Attribute, attributeCode string) string {
	for _, variant := range m.matchingVariants {
		attribute := variant.Attribute(attribute.Code)
		if attribute.Code == attributeCode {
			return variant.MarketPlaceCode
		}
	}
	return "" // We should never get here
}
