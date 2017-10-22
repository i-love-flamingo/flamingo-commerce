## Product Module

### About

* Provides product domain models and the related Port (ProductService).
* Provides Controller for product detail view, including variant logic

### Product Detail View

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
