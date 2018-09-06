# Product Package

## About

* Provides product domain models and the related secondary ports
* Provides Controller for product detail view, including variant logic
* Provides templatefuncs for:
    * get a (filtered) list of products (e.g. for product teasers)
    * get product urls
    * get a specific product

## Domain Layer
* A product is idendified by a "MarketplaceCode"
* There are different product Types, each of them need to implement the "BasicProduct" interface.

### Secondary Ports
The module defines two seconfary ports:

* ProductService interface to receive products
* SearchService interface, to search for product by any passed filter

### Product Types:
#### Simple Products
Represent a simple product that can be purchased directly

#### ConfigurableProduct and ConfigurableProductWithActiveVariant
Represent a product, that has several Variants. The configurable product cannot be sold directly.

But from a "Configurable Product" you can get the Saleable Variants - this is a product with Type "ConfigurableProductWithActiveVariant"

Here is an example illustrating this:

```go
   ...
   //productService is the injected imeplementation of interface "ProductService"
   product, err := c.productService.Get(ctx, "id_of_a_configurable_product")
   if product.Type() == TYPECONFIGURABLE {
      //type assert ConfigurableProduct
      if configurableProduct, ok := product.(ConfigurableProduct); ok {
        variantProduct, err := configurableProduct.GetConfigurableWithActiveVariant("id_of_an_variant")
      }
   }


```

## Product Detail View

The view gets the following Data passed:

```
    productViewData struct {
        Product          domain.BasicProduct
        VariantSelected  bool
        VariantSelection variantSelection
    }
``` 

## Dependencies:
* search package: the product.SearchService uses the search Result and Filter objects
