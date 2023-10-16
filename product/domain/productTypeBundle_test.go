package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/product/domain"
)

func TestGetBundleProductWithActiveChoices(t *testing.T) {
	t.Parallel()

	t.Run("returns bundle product matching bundle configuration", func(t *testing.T) {
		t.Parallel()

		b := domain.BundleProduct{
			Choices: []domain.Choice{
				{
					Identifier: "A",
					Options: []domain.Option{{
						Product: domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "A"}},
						MinQty:  1,
						MaxQty:  2,
					}},
					Required: true,
				},
				{
					Identifier: "B",
					Options: []domain.Option{{
						Product: domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "B"}},
						MinQty:  1,
						MaxQty:  2,
					}},
					Required: false,
				},
			},
		}

		bc := domain.BundleConfiguration{
			"A": {MarketplaceCode: "A", Qty: 1, VariantMarketplaceCode: ""},
			"B": {MarketplaceCode: "B", Qty: 2, VariantMarketplaceCode: ""},
		}

		bpac, err := b.GetBundleProductWithActiveChoices(bc)

		assert.Nil(t, err)
		assert.Equal(t, domain.BundleProductWithActiveChoices{
			BundleProduct: b,
			ActiveChoices: map[domain.Identifier]domain.ActiveChoice{
				"A": {
					Product:    domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "A"}},
					Qty:        1,
					Label:      "",
					Required:   true,
					Identifier: "A",
				},
				"B": {
					Product:    domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "B"}},
					Qty:        2,
					Label:      "",
					Required:   false,
					Identifier: "B",
				},
			},
		}, bpac)
	})

	t.Run("returns error selected qty out of range", func(t *testing.T) {
		t.Parallel()

		b := domain.BundleProduct{
			Choices: []domain.Choice{
				{
					Identifier: "A",
					Options: []domain.Option{{
						Product: domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "A"}},
						MinQty:  2,
						MaxQty:  2,
					}},
					Required: true,
				},
				{
					Identifier: "B",
					Options: []domain.Option{{
						Product: domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "B"}},
						MinQty:  1,
						MaxQty:  1,
					}},
					Required: false,
				},
			},
		}

		bc := domain.BundleConfiguration{
			"A": {MarketplaceCode: "A", Qty: 2, VariantMarketplaceCode: ""},
			"B": {MarketplaceCode: "B", Qty: 2, VariantMarketplaceCode: ""},
		}

		_, err := b.GetBundleProductWithActiveChoices(bc)

		assert.ErrorIs(t, err, domain.ErrSelectedQuantityOutOfRange)
	})

	t.Run("returns error selected qty 0 out of range", func(t *testing.T) {
		t.Parallel()

		b := domain.BundleProduct{
			Choices: []domain.Choice{
				{
					Identifier: "A",
					Options: []domain.Option{{
						Product: domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "A"}},
						MinQty:  2,
						MaxQty:  2,
					}},
					Required: true,
				},
				{
					Identifier: "B",
					Options: []domain.Option{{
						Product: domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "B"}},
						MinQty:  1,
						MaxQty:  1,
					}},
					Required: false,
				},
			},
		}

		bc := domain.BundleConfiguration{
			"A": {MarketplaceCode: "A", Qty: 2, VariantMarketplaceCode: ""},
			"B": {MarketplaceCode: "B", Qty: 0, VariantMarketplaceCode: ""},
		}

		_, err := b.GetBundleProductWithActiveChoices(bc)

		assert.ErrorIs(t, err, domain.ErrSelectedQuantityOutOfRange)
	})

	t.Run("returns error required choices not selected when required choices are not in the bundle config", func(t *testing.T) {
		t.Parallel()

		b := domain.BundleProduct{
			Choices: []domain.Choice{{Identifier: "A", Required: true}},
		}

		bc := domain.BundleConfiguration{}

		_, err := b.GetBundleProductWithActiveChoices(bc)

		assert.Equal(t, domain.ErrRequiredChoicesAreNotSelected, err)
	})

	t.Run("returns no error for missing none required bundle configuration choices", func(t *testing.T) {
		t.Parallel()

		b := domain.BundleProduct{
			Choices: []domain.Choice{
				{
					Identifier: "A",
					Options: []domain.Option{{
						Product: domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "A"}},
						MinQty:  1,
						MaxQty:  2,
					}},
					Required: true,
				},
				{
					Identifier: "B",
					Options: []domain.Option{{
						Product: domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "B"}},
						MinQty:  0,
						MaxQty:  1,
					}},
					Required: false,
				},
			},
		}

		bc := domain.BundleConfiguration{
			"A": {MarketplaceCode: "A", Qty: 2, VariantMarketplaceCode: ""},
		}

		_, err := b.GetBundleProductWithActiveChoices(bc)

		assert.NoError(t, err)
	})

	t.Run("error when variant not found", func(t *testing.T) {
		t.Parallel()

		bundleProduct := domain.BundleProduct{
			Choices: []domain.Choice{
				{
					Identifier: "choice_1",
					Required:   true,
					Options: []domain.Option{
						{
							Product: domain.ConfigurableProduct{
								BasicProductData: domain.BasicProductData{
									MarketPlaceCode: "configurable_1",
								},
								Variants: []domain.Variant{
									{
										BasicProductData: domain.BasicProductData{
											MarketPlaceCode: "variant_1",
										},
									},
									{
										BasicProductData: domain.BasicProductData{
											MarketPlaceCode: "variant_2",
										},
									},
								},
							},
							MinQty: 1,
							MaxQty: 5,
						},
					},
				},
			},
		}
		bundleConfiguration := domain.BundleConfiguration{
			"choice_1": {
				MarketplaceCode:        "configurable_1",
				VariantMarketplaceCode: "non_existent_variant",
				Qty:                    1,
			},
		}

		result, err := bundleProduct.GetBundleProductWithActiveChoices(bundleConfiguration)

		assert.Error(t, err, "Expected error but got nil")
		assert.Equal(t, 0, len(result.ActiveChoices), "Expected 0 active choices but got %d", len(result.ActiveChoices))
	})

	t.Run("one option is a configurable product", func(t *testing.T) {
		t.Parallel()

		b := domain.BundleProduct{
			Choices: []domain.Choice{
				{
					Identifier: "choice1",
					Label:      "Choice 1",
					Required:   true,
					Options: []domain.Option{
						{
							Product: domain.SimpleProduct{
								BasicProductData: domain.BasicProductData{
									MarketPlaceCode: "mpc1",
								},
							},
							MinQty: 1,
							MaxQty: 2,
						},
						{
							Product: domain.ConfigurableProduct{
								BasicProductData: domain.BasicProductData{
									MarketPlaceCode: "mpc2",
								},
								Variants: []domain.Variant{
									{
										BasicProductData: domain.BasicProductData{
											MarketPlaceCode: "vmc1",
										},
									},
									{
										BasicProductData: domain.BasicProductData{
											MarketPlaceCode: "vmc2",
										},
									},
								},
							},
							MinQty: 1,
							MaxQty: 2,
						},
					},
				},
			},
		}

		bundleConfiguration := domain.BundleConfiguration{
			"choice1": domain.ChoiceConfiguration{
				MarketplaceCode:        "mpc2",
				Qty:                    1,
				VariantMarketplaceCode: "vmc1",
			},
		}

		bundleProductWithActiveChoices, err := b.GetBundleProductWithActiveChoices(bundleConfiguration)
		assert.NoError(t, err)
		assert.Equal(t, domain.BundleProduct{Choices: b.Choices}, bundleProductWithActiveChoices.BundleProduct)
	})
}

func Test_ExtractBundleConfig(t *testing.T) {
	t.Parallel()

	t.Run("get config with simples and configurables correctly", func(t *testing.T) {
		t.Parallel()

		b := domain.BundleProductWithActiveChoices{
			ActiveChoices: map[domain.Identifier]domain.ActiveChoice{
				domain.Identifier("identifier1"): {
					Product:  domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "marketplace1"}},
					Qty:      1,
					Required: true,
				},
				domain.Identifier("identifier2"): {
					Product:  domain.SimpleProduct{BasicProductData: domain.BasicProductData{MarketPlaceCode: "marketplace2"}},
					Qty:      1,
					Required: true,
				},
				domain.Identifier("identifier3"): {
					Product: domain.ConfigurableProductWithActiveVariant{
						BasicProductData: domain.BasicProductData{
							MarketPlaceCode: "marketplace3",
						},
						ActiveVariant: domain.Variant{
							BasicProductData: domain.BasicProductData{
								MarketPlaceCode: "variantMarketplace3",
							},
						},
					},
					Qty:      1,
					Required: true,
				},
			},
		}

		want := domain.BundleConfiguration{
			domain.Identifier("identifier1"): {
				MarketplaceCode: "marketplace1",
				Qty:             1,
			},
			domain.Identifier("identifier2"): {
				MarketplaceCode: "marketplace2",
				Qty:             1,
			},
			domain.Identifier("identifier3"): {
				MarketplaceCode:        "marketplace3",
				VariantMarketplaceCode: "variantMarketplace3",
				Qty:                    1,
			},
		}

		config := b.ExtractBundleConfig()

		assert.Equal(t, want, config)
	})

	t.Run("return nil when active choices are empty", func(t *testing.T) {
		t.Parallel()

		b := domain.BundleProductWithActiveChoices{}

		config := b.ExtractBundleConfig()

		assert.Nil(t, config)
	})
}
