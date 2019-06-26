package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

type Service struct{}

func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Product": graphql.ModelMapEntry{
			Type: new(domain.BasicProduct),
			Fields: map[string]string{
				"specifications": "GetSpecifications",
			},
		},
		"Commerce_SimpleProduct": graphql.ModelMapEntry{
			Type: domain.SimpleProduct{},
			Fields: map[string]string{
				"specifications": "GetSpecifications",
			},
		},
		"Commerce_BasicProductData":          domain.BasicProductData{},
		"Commerce_ProductTeaserData":         domain.TeaserData{},
		"Commerce_ProductSpecifications":     domain.Specifications{},
		"Commerce_ProductSpecificationGroup": domain.SpecificationGroup{},
		"Commerce_ProductSpecificationEntry": domain.SpecificationEntry{},
		"Commerce_ProductSaleable":           domain.Saleable{},
		"Commerce_ProductMedia":              domain.Media{},
		"Commerce_ProductAttributes":         domain.Attributes{},
		"Commerce_ProductAttribute":          domain.Attribute{},
		"Commerce_CategoryTeaser":            domain.CategoryTeaser{},
		"Commerce_ProductPriceInfo":          domain.PriceInfo{},
	}.Models()
}

type CommerceProductQueryResolver struct {
	productService domain.ProductService
}

func (r *CommerceProductQueryResolver) Inject(productService domain.ProductService) {
	r.productService = productService
}

func (r *CommerceProductQueryResolver) CommerceProduct(ctx context.Context, marketplaceCode string) (domain.BasicProduct, error) {
	return r.productService.Get(ctx, marketplaceCode)
}
