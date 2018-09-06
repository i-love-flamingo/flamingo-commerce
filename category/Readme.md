# Module category

* Provides domain model for (product)category and category tree, and the related categoryService interface (=Secondary Port)
* Provides Controller for rendering category pages, supporting different templates baed on the category type.
* Since its expected that products should be shown, there is a dependency to the product.Searchservice from the product package
  * (also the domain defines a special "CategoryFilter" - the used implementation of the product searchservice should understand this special filter)

## Configurations

```
category.view.template = "category/category"
category.view.teaserTemplate = "category/teaser"
```

## Dependencies:
* product package: (for product searchservice - expected to understand the CategoryFilter) 
* search package: (for pagination)
