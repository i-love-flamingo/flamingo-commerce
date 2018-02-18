# Package category
* Provides domain model for (product)category and category tree, and the related categoryService interface (=Secondary Port)
* Provides Controller for rendering category pages. Since its expected that products should be shown, there is a dependency to the product.Searchservice from the product package

## Dependencies:
* product package: (for product searchservice)
* search package: (for pagination)
