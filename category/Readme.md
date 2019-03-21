# Category Module

* Domain Layer: Provides domain model for:
    * (product)category with potential data a category can have (like name, media, ..)
    * and a (category) tree
    * and the related categoryService interface (=Secondary Port)
* Interface Layer:
    * Provides controller for rendering category pages, supporting different templates baed on the category type.
    * Provides data controller to access category and tree from inside templates
    
* Product Search:
    * Since its expected that products should be shown, there is a dependency to the "product" module - more specific to the productSearchService.
    * the category module defines in its domain also a "CategoryFilter" (that implements the search filter interface): This filter is passed to the priductSearchService, so any implementation of the product searchservice should understand this special filter.

## Configurations

Set the templates for the category single view
```
//default template
category.view.template = "category/category"

//template used for category type "teaser"
category.view.teaserTemplate = "category/teaser"
```

## Usage in templates
This module provides two data controller that can be used to get category and tree objects:
```
- var rootCategoryTree = data('category.tree', {'code': ''})

 each category in rootCategoryTree.categories
 
- var category = data("category",{'code': 'category-code'})

```

## Dependencies:
* product package: (for product searchservice) 
* search package: (for pagination)
