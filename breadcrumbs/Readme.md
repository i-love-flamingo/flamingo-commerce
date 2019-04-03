# Breadcrumbs

The breadcrumbs module contains a small helper that supports breadcrumb navigation data.

It can be used from two perspectives:
* other modules controller can add typical breadcrumb parts - for example the category module adds breadcrumb parts that shows the root category down to the current category. 
Or the product controller adds the correct breadcrumb to the product view.
* you can use the collected breadcrumb in your template

## Usage in templates

Call the data function. For example, in the pug template engine a usage can look like this:
    
```pug
var breadCrumbData = data('breadcrumbs')
ul
  each item, index in items    
    li
      if item.url === ""
        span.breadcrumbNoLink=item.title
      else
        a(href=item.url)=item.title
```
