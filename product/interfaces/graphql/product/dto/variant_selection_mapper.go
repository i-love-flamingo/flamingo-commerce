package graphqlproductdto

import (
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

func MapVariantSelections(product domain.BasicProduct) VariantSelection {
	if product.Type() == domain.TypeConfigurable {
		configurable, ok := product.(domain.ConfigurableProduct)
		if ok {
			return mapVariations(configurable.VariantVariationAttributes, configurable.Variants)
		}
	}

	if product.Type() == domain.TypeConfigurableWithActiveVariant {
		configurableWithActiveVariants, ok := product.(domain.ConfigurableProductWithActiveVariant)
		if ok {
			return mapVariations(configurableWithActiveVariants.VariantVariationAttributes, configurableWithActiveVariants.Variants)
		}
	}

	return VariantSelection{}
}

// mapVariations ranges over all variants and inserts every one existing selection into the structure
func mapVariations(variantVariation []string, variants []domain.Variant) VariantSelection {
	selection := VariantSelection{}

	for _, variant := range variants {
		if !variant.HasAllAttributes(variantVariation) {
			continue
		}

		variantValues := map[string]domain.Attribute{}

		for _, variantVariation := range variantVariation {
			attribute := variant.Attribute(variantVariation)
			variantValues[variantVariation] = attribute
		}

		selection = addToVariationSelection(selection, variant, variantValues)
	}

	return selection
}

func addToVariationSelection(v VariantSelection, variant domain.Variant, variantVariationValues map[string]domain.Attribute) VariantSelection {
	variantSelectionVariant := VariantSelectionVariant{
		Variant: VariantSelectionVariantMatchingVariant{MarketplaceCode: variant.MarketPlaceCode},
	}

	for variantVariation, value := range variantVariationValues {
		variantSelectionVariant.MatchingAttributes = append(variantSelectionVariant.MatchingAttributes, MatchingVariantSelection{
			Key:   variantVariation,
			Value: value.Label,
		})

		attribute, attributePosition := findOrCreateVariantSelectionAttribute(variantVariation, value, v.Attributes)
		attributeOption, attributeOptionPosition := findOrCreateVariantSelectionAttributeOption(value, attribute.Options)

		for restriction, restrictionValue := range variantVariationValues {
			if variantVariation == restriction {
				continue
			}

			otherAttributeRestrictions, attributeRestrictionPosition := findOrCreateOtherAttributeRestriction(restriction, attributeOption.OtherAttributesRestrictions)

			otherAttributeRestrictions.AvailableOptions = append(otherAttributeRestrictions.AvailableOptions, restrictionValue.Label)
			attributeOption.OtherAttributesRestrictions = appendOtherAttributeRestrictions(attributeOption.OtherAttributesRestrictions,
				attributeRestrictionPosition, otherAttributeRestrictions)
		}

		attribute.Options = appendOptions(attribute.Options, attributeOptionPosition, attributeOption)

		v.Attributes = appendSelectionAttributes(v.Attributes, attributePosition, attribute)
	}

	v.Variants = append(v.Variants, variantSelectionVariant)

	return v
}

func findOrCreateVariantSelectionAttribute(key string, domainAttribute domain.Attribute, attributes []VariantSelectionAttribute) (VariantSelectionAttribute, int) {
	for i, attribute := range attributes {
		if attribute.Code == key {
			return attribute, i
		}
	}

	return VariantSelectionAttribute{
		Label: domainAttribute.CodeLabel,
		Code:  key,
	}, -1
}

func findOrCreateVariantSelectionAttributeOption(attribute domain.Attribute, options []VariantSelectionAttributeOption) (VariantSelectionAttributeOption, int) {
	for i, option := range options {
		if option.Label == attribute.Label {
			return option, i
		}
	}

	return VariantSelectionAttributeOption{
		Label:    attribute.Label,
		UnitCode: attribute.UnitCode,
	}, -1
}

func findOrCreateOtherAttributeRestriction(key string, restrictions []OtherAttributesRestriction) (OtherAttributesRestriction, int) {
	for i, restriction := range restrictions {
		if restriction.Code == key {
			return restriction, i
		}
	}

	return OtherAttributesRestriction{
		Code: key,
	}, -1
}

func appendOptions(options []VariantSelectionAttributeOption, pos int, option VariantSelectionAttributeOption) []VariantSelectionAttributeOption {
	if pos == -1 {
		options = append(options, option)
		return options
	}

	options[pos] = option
	return options
}

func appendOtherAttributeRestrictions(restrictions []OtherAttributesRestriction, pos int, restriction OtherAttributesRestriction) []OtherAttributesRestriction {
	if pos == -1 {
		restrictions = append(restrictions, restriction)
		return restrictions
	}

	restrictions[pos] = restriction
	return restrictions
}

func appendSelectionAttributes(attributes []VariantSelectionAttribute, pos int, attribute VariantSelectionAttribute) []VariantSelectionAttribute {
	if pos == -1 {
		attributes = append(attributes, attribute)
		return attributes
	}

	attributes[pos] = attribute
	return attributes
}
