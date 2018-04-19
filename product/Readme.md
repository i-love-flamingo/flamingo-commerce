# Product Package

## About

* Provides product domain models and the related Port (ProductService).
* Provides Controller for product detail view, including variant logic

## Domain Layer
* There are 

## Product Detail View

The view gets the following Data passed:

```
	productViewData struct {
		// simple / configurable / configurable_with_variant
		RenderContext       string
		SimpleProduct       domain.SimpleProduct
		ConfigurableProduct domain.ConfigurableProduct
		ActiveVariant       domain.Variant
		VariantSelected     bool
		VariantSelection    variantSelection
	}
``` 

## Dependencies:
* search package: the product.SearchService uses the search Result and Filter objects
