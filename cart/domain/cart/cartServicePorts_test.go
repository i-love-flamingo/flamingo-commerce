package cart_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/product/domain"

	"github.com/stretchr/testify/assert"
)

func TestBundleConfiguration_Equals(t *testing.T) {
	t.Parallel()

	t.Run("empty bundle configuration equals nil", func(t *testing.T) {
		t.Parallel()

		bc := domain.BundleConfiguration{}

		assert.True(t, bc.Equals(nil))
	})

	t.Run("nil bundle configuration equals empty", func(t *testing.T) {
		t.Parallel()

		var bc domain.BundleConfiguration

		assert.True(t, bc.Equals(domain.BundleConfiguration{}))
	})

	t.Run("bundle configuration equals itself", func(t *testing.T) {
		t.Parallel()

		bc := domain.BundleConfiguration{
			"choice-1": {
				MarketplaceCode:        "code-1",
				VariantMarketplaceCode: "variant-1",
				Qty:                    1,
			},
			"choice-2": {
				MarketplaceCode:        "code-2",
				VariantMarketplaceCode: "variant-2",
				Qty:                    2,
			},
		}

		other := bc

		assert.True(t, bc.Equals(other))
	})

	t.Run("bundle configuration does not equal different choice", func(t *testing.T) {
		t.Parallel()

		bc := domain.BundleConfiguration{
			"choice-1": {
				MarketplaceCode:        "code-1",
				VariantMarketplaceCode: "variant-1",
				Qty:                    1,
			},
			"choice-2": {
				MarketplaceCode:        "code-2",
				VariantMarketplaceCode: "variant-2a",
				Qty:                    2,
			},
		}

		other := domain.BundleConfiguration{
			"choice-1": {
				MarketplaceCode:        "code-1",
				VariantMarketplaceCode: "variant-1",
				Qty:                    1,
			},
			"choice-2": {
				MarketplaceCode:        "code-2",
				VariantMarketplaceCode: "variant-2b",
				Qty:                    2,
			},
		}

		assert.False(t, bc.Equals(other))
	})

	t.Run("bundle configuration does not equal extra choice", func(t *testing.T) {
		t.Parallel()

		bc := domain.BundleConfiguration{
			"choice-1": {
				MarketplaceCode:        "code-1",
				VariantMarketplaceCode: "variant-1",
				Qty:                    1,
			},
			"choice-2": {
				MarketplaceCode:        "code-2",
				VariantMarketplaceCode: "variant-2",
				Qty:                    2,
			},
		}

		other := domain.BundleConfiguration{
			"choice-1": {
				MarketplaceCode:        "code-1",
				VariantMarketplaceCode: "variant-1",
				Qty:                    1,
			},
			"choice-2": {
				MarketplaceCode:        "code-2",
				VariantMarketplaceCode: "variant-2",
				Qty:                    2,
			},
			"choice-3": {
				MarketplaceCode:        "code-3",
				VariantMarketplaceCode: "variant-3",
				Qty:                    3,
			},
		}

		assert.False(t, bc.Equals(other))
	})
}
