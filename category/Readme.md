# Package category
* Provides domain model for (product)category and category tree, and the related categoryService interface (=Secondary Port)
  * (also the domain defines a special "CategoryFilter" - the used implementation of the product searchservice should understand this special filter)
* Provides Controller for rendering category pages. Since its expected that products should be shown, there is a dependency to the product.Searchservice from the product package

## Dependencies:
* product package: (for product searchservice - expected to understand the CategoryFilter) 
* search package: (for pagination)
