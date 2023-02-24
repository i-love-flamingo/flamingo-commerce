package graphqlproductdto

import (
	"sort"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

func MapVariantSelections(configurable domain.ConfigurableProduct) VariantSelection {
	return mapVariations(configurable.VariantVariationAttributes, configurable.VariantVariationAttributesSorting, configurable.Variants)
}

// mapVariations ranges over all variants and inserts every one existing selection into the structure
func mapVariations(variantVariation []string, variantVariationSorting map[string][]string, variants []domain.Variant) VariantSelection {
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

	return sortSelection(variantVariation, variantVariationSorting, selection)
}

func addToVariationSelection(v VariantSelection, variant domain.Variant, variantVariationValues map[string]domain.Attribute) VariantSelection {
	variantSelectionVariant := VariantSelectionMatch{
		Variant: VariantSelectionMatchVariant{MarketplaceCode: variant.MarketPlaceCode, VariantData: variant},
	}

	for variantVariation, value := range variantVariationValues {
		variantSelectionVariant.Attributes = append(variantSelectionVariant.Attributes, VariantSelectionMatchAttributes{
			Key:   variantVariation,
			Value: value.Label,
		})
		// we need positions because we cannot mutate variables belonging to loop
		// if position is -1 then it simply appends to a slice
		attribute, attributePosition := findOrCreateVariantSelectionAttribute(variantVariation, value, v.Attributes)
		attributeOption, attributeOptionPosition := findOrCreateVariantSelectionAttributeOption(value, attribute.Options)

		for restriction, restrictionValue := range variantVariationValues {
			// skip because we do not want to insert to attribute color with option red possible
			// attribute color with option red
			if variantVariation == restriction {
				continue
			}

			otherAttributeRestrictions, attributeRestrictionPosition := findOrCreateOtherAttributeRestriction(restriction, attributeOption.OtherAttributesRestrictions)

			// insert available options for current attribute (available options: m,l current attribute: size)
			otherAttributeRestrictions.AvailableOptions = append(otherAttributeRestrictions.AvailableOptions, restrictionValue.Label)

			// insert other attributes available for current attribute (current attribute: color-blue, available attribute: size-m-l)
			attributeOption.OtherAttributesRestrictions = appendOtherAttributeRestrictions(attributeOption.OtherAttributesRestrictions,
				attributeRestrictionPosition, otherAttributeRestrictions)
		}

		// insert options available for current attribute (current attribute: color available options: red, blue)
		attribute.Options = appendOptions(attribute.Options, attributeOptionPosition, attributeOption)

		// insert all possible selections
		v.Attributes = appendSelectionAttributes(v.Attributes, attributePosition, attribute)
	}

	// insert variant with its attributes
	v.Variants = append(v.Variants, variantSelectionVariant)

	return v
}

func sortSelection(variantVariation []string, variantVariationSorting map[string][]string, selection VariantSelection) VariantSelection {
	for attributeIndex, attribute := range selection.Attributes {
		for optionIndex, option := range attribute.Options {
			for restrictionIndex := range option.OtherAttributesRestrictions {
				// sort available options for current attribute (available options: m,l current attribute: size)
				sort.Slice(selection.Attributes[attributeIndex].Options[optionIndex].OtherAttributesRestrictions[restrictionIndex].AvailableOptions, func(i, j int) bool {
					return indexOf(variantVariationSorting[attribute.Code], selection.Attributes[attributeIndex].Options[optionIndex].OtherAttributesRestrictions[restrictionIndex].AvailableOptions[i]) <
						indexOf(variantVariationSorting[attribute.Code], selection.Attributes[attributeIndex].Options[optionIndex].OtherAttributesRestrictions[restrictionIndex].AvailableOptions[j])
				})
			}

			// sort other attributes available for current attribute (current attribute: color-blue, available attribute: size-m-l)
			sort.Slice(selection.Attributes[attributeIndex].Options[optionIndex].OtherAttributesRestrictions, func(i, j int) bool {
				return indexOf(variantVariationSorting[attribute.Code], selection.Attributes[attributeIndex].Options[optionIndex].OtherAttributesRestrictions[i].Code) <
					indexOf(variantVariationSorting[attribute.Code], selection.Attributes[attributeIndex].Options[optionIndex].OtherAttributesRestrictions[j].Code)
			})
		}

		// sort options available for current attribute (current attribute: color available options: red, blue)
		sort.Slice(attribute.Options, func(i, j int) bool {
			return indexOf(variantVariationSorting[attribute.Code], attribute.Options[i].Label) <
				indexOf(variantVariationSorting[attribute.Code], attribute.Options[j].Label)
		})
	}

	// sort selections
	sort.Slice(selection.Attributes, func(i, j int) bool {
		return indexOf(variantVariation, selection.Attributes[i].Code) <
			indexOf(variantVariation, selection.Attributes[j].Code)
	})

	for variantIndex := range selection.Variants {
		// sort all possible variants
		sort.Slice(selection.Variants[variantIndex].Attributes, func(i, j int) bool {
			return indexOf(variantVariation, selection.Variants[variantIndex].Attributes[i].Key) <
				indexOf(variantVariation, selection.Variants[variantIndex].Attributes[j].Key)
		})
	}

	return selection
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

func indexOf(slice []string, element string) int {
	for i, v := range slice {
		if v == element {
			return i
		}
	}
	return -1
}
