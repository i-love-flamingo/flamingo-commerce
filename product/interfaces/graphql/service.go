package graphql

import (
	"context"
	categoryDomain "flamingo.me/flamingo-commerce/v3/category/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -nometadata -o fs.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models return the 'Schema name' => 'Go model' mapping of this module
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
		"Commerce_ProductAttributes": graphql.ModelMapEntry{
			Type: domain.Attributes{},
			Fields: map[string]string{
				"getAttribute":     "Attribute",
				"getAttributes":    "Attributes",
				"getAttributeKeys": "AttributeKeys",
			},
		},
		"Commerce_ProductAttribute":        domain.Attribute{},
		"Commerce_CategoryTeaser":          domain.CategoryTeaser{},
		"Commerce_ProductPriceInfo":        domain.PriceInfo{},
		"Commerce_ProductLoyaltyPriceInfo": domain.LoyaltyPriceInfo{},
		"Commerce_PriceContext":            domain.PriceContext{},

		"Commerce_CategoryProductList": CommerceCategoryProductList{},
		"Commerce_Category":            new(categoryDomain.Category),
		"Commerce_CategoryData":        categoryDomain.CategoryData{},
	}.Models()
}

// CommerceProductQueryResolver resolves graphql product queries
type CommerceProductQueryResolver struct {
	productService domain.ProductService
}

// CommerceProductListByCategoryQueryResolver resolves graphql product queries
type CommerceProductListByCategoryQueryResolver struct {
	categoryService categoryDomain.CategoryService
	searchService   domain.SearchService
}

// CommerceCategoryProductList holds structure for graphql products by category endpoint
type CommerceCategoryProductList struct {
	Category categoryDomain.Category
	Products []domain.BasicProduct
}

// Inject dependencies
func (r *CommerceProductQueryResolver) Inject(productService domain.ProductService) {
	r.productService = productService
}

// Inject dependencies
func (r *CommerceProductListByCategoryQueryResolver) Inject(categoryService categoryDomain.CategoryService, searchService domain.SearchService) {
	r.categoryService = categoryService
	r.searchService = searchService
}

// CommerceProduct returns a product with the given marketplaceCode from productService
func (r *CommerceProductQueryResolver) CommerceProduct(ctx context.Context, marketplaceCode string) (domain.BasicProduct, error) {
	return r.productService.Get(ctx, marketplaceCode)
}

// CommerceProductListByCategory returns products with the given categoryCode from searchService
func (r *CommerceProductListByCategoryQueryResolver) CommerceProductListByCategory(ctx context.Context, categoryCode string) (*CommerceCategoryProductList, error) {
	category, err := r.categoryService.Get(ctx, categoryCode)

	if err != nil {
		return &CommerceCategoryProductList{Category: category, Products: nil}, err
	}

	// - use categoryDomain.NewCategoryFacet as filter to use product/category endpoint from searchperience
	result, err := r.searchService.Search(ctx, categoryDomain.NewCategoryFacet(categoryCode))
	// result.SearchMeta includes pagination info, use implementation of usage of paginationInfoFactory in category controller as reference
	if err != nil {
		return &CommerceCategoryProductList{Category: category, Products: nil}, err
	}

	return &CommerceCategoryProductList{Category: category, Products: result.Hits}, err
}
