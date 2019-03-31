# Search Module

Defines a common search domain and provides controllers to render search result and autosuggest.

The module also offers a pagination util, that can be used to pass pagination data to the interfaces.

## Domain Layer

* A Search can return Results of different types. So a Document in a search result may be a "Product" or anything else
    * SearchService->Search(Filter)  returns Map of Results (by type)
    * Document

### Secondary Ports
* The SearchService needs to be implemented
* Please note that a `Document` is defined as an interface and can be "anything". This way the search can be used very generic and can return documents of any type (e.g. products, categories, content, brands etc).
