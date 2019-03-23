# Product Module

* Provides product domain models and the related secondary ports
* Provides Controller for product detail view, including variant logic
* Provides templatefuncs for:
    * get a (filtered) list of products (e.g. for product teasers)
    * get product urls
    * get a specific product

## Domain Layer
* A product is idendified by a "MarketplaceCode"
* There are different product types, each of them need to implement the "BasicProduct" interface.

## Details about Price fields

* The product has one main price (see PriceInfo Property) that stands for the value of that product.
* Product has maintained price in certain currency - eventually different by channel or locale.
* The price is either maintained as gross (B2C use cases) or net price (B2B use cases). How to interpret the price is up to the cart and view logic. 
    * Any logic depending on the gross/net knowledge should use the configpath `commerce.product.priceIsGross`

* the product might be currently discounted and has a discounted price (the discounted price is also either gros or net price like the normal price)

About Charges:
* A Charge is a price that need to be payed for that product. The is normaly the product price.
* But this concept allows to control "in what currency and type" a customer need to pay the price of the product (See loyalty below)

Loyaltyprices:
* The product might also have loyaltyprices that allows to set a price in points, but also a minimum amount of points that need to be spend
* In cases where the customer need to spend a certain amount of points, the Method "GetCharges()" will return the different charges. Again it is up for other modules to interpret if this is gross or net.

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
    // productViewData is used for product rendering
    productViewData struct {
        // simple / configurable / configurable_with_variant
        RenderContext    string
        Product          domain.BasicProduct
        VariantSelected  bool
        VariantSelection variantSelection
        BackURL          string
    }
``` 

## Templatefuncs

### getProduct

Returns the product or nil:
`- var product = getProduct("marketplacecode")`


### getProductUrl

Returns the correct url to the product:
`getProductUrl(product)`

### findProducts

findProducts is a template function that returns a searchresult to show products:
* it will dynamically adjust the search request based on url parameters that match the given namespace
* the templatefunctions requires 4 parameters:
    * namespace: A name of the findProduct call - it is used to url check parameters with that namespace
    * searchConfig: A map with the following keys supported:
        * query: optional - the searchstring that a "human" might have entered to filter the the search
        * pageSize, page: Optional - set the page and the pageSize (for pagination)
        * sortBy, sortDirection (A/D): Optional set the field that should be used to sort the search result
    * keyValueFilters: A map of key values that are used as additional keyValue Filters in the searchRequest
    * filterConstrains: Optional - A map that supports the following keys:
        * blackList or whiteList (if both given whiteList is prefered)
        * This is a comma seperated list of filter keys, that are evaluated during:
            * using url Parameters to set keyValueFilters - only allowed keys will be used
            * modifying the searchResult Facetcollection and remove disallowed facets

See the following example:
```
  - 
    var searchResult1 = findProducts("list1",
                                    {"query":"","pageSize": "10","sortDirection":"A","sortBy":'name'},
                                    {"brand":"Hello Kitty"},
                                    {"whiteList":"brand,color,sortBy"})

  ul
    each product in searchResult.products
      li Title: #{product.baseData.title} Retailer: #{product.baseData.retailerCode} Brand: #{product.baseData.attributes["brandCode"].value}
      
    - 
      var searchResult2 = findProducts("list2",
                                      {"query":"","pageSize": "10","sortDirection":"A","sortBy":'name'},
                                      {"brand":"Hello Kitty"},
                                      {"whiteList":"brand,color"})
  
    ul
      each product in searchResult.products
        li Title: #{product.baseData.title} Retailer: #{product.baseData.retailerCode} Brand: #{product.baseData.attributes["brandCode"].value}
```

It will show in the list1 10 products that match the brand "Hello Kitty".
If you call the page that contains the call to that template function with urlparameters, the response can be dynamically modified.
For example:

URL/?list1.brand=Apple&list2.color=black

Will modify the result of the two lists accordingly - while

URL/?list2.notallowedkey=black

will not.


## Dependencies:
* search package: the product.SearchService uses the search Result and Filter objects
