package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

func TestURLService_GetURLParams(t *testing.T) {

	type fields struct {
		config *struct {
			GenerateSlug      bool   `inject:"config:commerce.product.generateSlug,optional"`
			SlugAttributecode string `inject:"config:commerce.product.slugAttributeCode,optional"`
		}
	}
	type args struct {
		product     domain.BasicProduct
		variantCode string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "nil product",
			want: make(map[string]string),
		},
		{
			name: "simple, generate slug",
			fields: fields{
				config: getConfig(true, "slug"),
			},
			args: args{
				product: domain.SimpleProduct{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
						Attributes: domain.Attributes{
							"slug": domain.Attribute{
								Code:     "slug",
								RawValue: "test-slug",
							},
						},
					},
				},
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "test-name",
			},
		},
		{
			name: "simple, use slug, attribute missing",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.SimpleProduct{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
						Attributes: domain.Attributes{
							"slug-invalid": domain.Attribute{
								Code:     "slug",
								RawValue: "test-slug",
							},
						},
					},
				},
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "test-name",
			},
		},
		{
			name: "simple, use slug, attribute empty",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.SimpleProduct{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
						Attributes: domain.Attributes{
							"slug": domain.Attribute{
								Code: "slug",
							},
						},
					},
				},
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "test-name",
			},
		},
		{
			name: "simple, use slug",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.SimpleProduct{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
						Attributes: domain.Attributes{
							"slug": domain.Attribute{
								Code:     "slug",
								RawValue: "slug-test-name",
							},
						},
					},
				},
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "slug-test-name",
			},
		},
		{
			name: "configurable, use slug",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.ConfigurableProduct{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
						Attributes: domain.Attributes{
							"slug": domain.Attribute{
								Code:     "slug",
								RawValue: "slug-test-name",
							},
						},
					},
				},
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "slug-test-name",
			},
		},
		{
			name: "configurable, active variant, use slug",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.ConfigurableProductWithActiveVariant{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
						Attributes: domain.Attributes{
							"slug": domain.Attribute{
								Code:     "slug",
								RawValue: "slug-test-name",
							},
						},
					},
					ActiveVariant: domain.Variant{
						BasicProductData: domain.BasicProductData{
							MarketPlaceCode: "variant-test-code",
							Title:           "Variant Name",
							Attributes: domain.Attributes{
								"slug": domain.Attribute{
									Code:     "slug",
									RawValue: "slug-variant-test-name",
								},
							},
						},
					},
				},
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "slug-variant-test-name",
				"variantcode":     "variant-test-code",
			},
		},
		{
			name: "configurable, active variant, withvariant code, use slug",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.ConfigurableProductWithActiveVariant{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
						Attributes: domain.Attributes{
							"slug": domain.Attribute{
								Code:     "slug",
								RawValue: "slug-test-name",
							},
						},
					},
					Variants: []domain.Variant{
						{
							BasicProductData: domain.BasicProductData{
								MarketPlaceCode: "selected-variant",
								Title:           "Select Variant Name",
								Attributes: domain.Attributes{
									"slug": domain.Attribute{
										Code:     "slug",
										RawValue: "slug-selected-variant-test-name",
									},
								},
							},
						},
					},
					ActiveVariant: domain.Variant{
						BasicProductData: domain.BasicProductData{
							MarketPlaceCode: "variant-test-code",
							Title:           "Variant Name",
							Attributes: domain.Attributes{
								"slug": domain.Attribute{
									Code:     "slug",
									RawValue: "slug-variant-test-name",
								},
							},
						},
					},
				},
				variantCode: "selected-variant",
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "slug-selected-variant-test-name",
				"variantcode":     "selected-variant",
			},
		},
		{
			name: "configurable, with teaser, use slug, no slug on teaser set",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.ConfigurableProduct{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
					},
					Variants: []domain.Variant{
						{
							BasicProductData: domain.BasicProductData{
								MarketPlaceCode: "selected-variant",
								Title:           "Select Variant Name",
								Attributes: domain.Attributes{
									"slug": domain.Attribute{
										Code:     "slug",
										RawValue: "slug-selected-variant-test-name",
									},
								},
							},
						},
					},
					Teaser: domain.TeaserData{
						PreSelectedVariantSku: "teaser-preselected-variant",
						ShortTitle:            "teaser-short-title",
					},
				},
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "teaser-short-title",
				"variantcode":     "teaser-preselected-variant",
			},
		},
		{
			name: "configurable, with teaser, use slug, no slug on teaser set",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.ConfigurableProduct{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
					},
					Variants: []domain.Variant{
						{
							BasicProductData: domain.BasicProductData{
								MarketPlaceCode: "selected-variant",
								Title:           "Select Variant Name",
								Attributes: domain.Attributes{
									"slug": domain.Attribute{
										Code:     "slug",
										RawValue: "slug-selected-variant-test-name",
									},
								},
							},
						},
					},
					Teaser: domain.TeaserData{
						PreSelectedVariantSku: "teaser-preselected-variant",
						ShortTitle:            "teaser-short-title",
						URLSlug:               "teaser-url-slug",
					},
				},
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "teaser-url-slug",
				"variantcode":     "teaser-preselected-variant",
			},
		},
		{
			name: "configurable, variant set, use slug",
			fields: fields{
				config: getConfig(false, "slug"),
			},
			args: args{
				product: domain.ConfigurableProduct{
					BasicProductData: domain.BasicProductData{
						MarketPlaceCode: "test-code",
						Title:           "Test Name",
					},
					Variants: []domain.Variant{
						{
							BasicProductData: domain.BasicProductData{
								MarketPlaceCode: "selected-variant",
								Title:           "Select Variant Name",
								Attributes: domain.Attributes{
									"slug": domain.Attribute{
										Code:     "slug",
										RawValue: "slug-selected-variant-test-name",
									},
								},
							},
						},
					},
				},
				variantCode: "selected-variant",
			},
			want: map[string]string{
				"marketplacecode": "test-code",
				"name":            "slug-selected-variant-test-name",
				"variantcode":     "selected-variant",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &application.URLService{}
			s.Inject(
				nil,
				tt.fields.config,
			)
			got := s.GetURLParams(tt.args.product, tt.args.variantCode)
			assert.Equal(t, tt.want, got, "url params")
		})
	}
}

func getConfig(generate bool, code string) *struct {
	GenerateSlug      bool   `inject:"config:commerce.product.generateSlug,optional"`
	SlugAttributecode string `inject:"config:commerce.product.slugAttributeCode,optional"`
} {

	return &struct {
		GenerateSlug      bool   `inject:"config:commerce.product.generateSlug,optional"`
		SlugAttributecode string `inject:"config:commerce.product.slugAttributeCode,optional"`
	}{
		GenerateSlug:      generate,
		SlugAttributecode: code,
	}
}
