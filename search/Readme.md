# Search package

Defines a common search domain and provides controllers to render search result and autosuggest.

The module also offers a Pagination Util, that can be used to pass pagination data to the interfaces.

## Domain Layer

* A Search can return Results of different types. So a Document in a search result may be a "Product" or anything else
    * SearchService->Search(Filter)  returns Map of Results (per type)
    * Document

### Secondary Ports
* The SearchService need to be implemented
