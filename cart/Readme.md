## Cart Package ##

The cart package offers a domain model for carts and its items.
It requires several ports to be implemented in order to use it.

### Ports ###

* Domain/GuestCartService: Implement for Guest Carts. This service is reponsible to get Guest Carts by Id
  * Domain/CartBehaviour: The Servie nneed to make sure that the a correct implementation for "Domain/CartBehaviour" is passed to the Cart:
    * This interface is reponsible to manage updates of items in the cart
* Domain/CustomerCartService: Implement to support Customer Carts: This service is reponsible to get Carts by Customer Token
  *  Domain/CartBehaviour: The CartBehaviour for Customer Carts
* Domain/CartValidator: Optional Implement a Validator for your Cart. The Validator returns a struct with validation messages and results.

### Configurations###

```
  
  cart:
    # To register the in memory cart service adapter (e.g. for development mode)
    useInMemoryCartServiceAdapters: true
    
```

### Cart API ###

### Get Cart Content:
* http://localhost:3210/en/api/cart

### Adding products:

* Simple product: http://localhost:3210/en/api/cart/add/fake_simple
* With qty: http://localhost:3210/en/api/cart/add/fake_simple?qty=10
* Adding configurables: http://localhost:3210/en/api/cart/add/fake_configurable?variantMarketplaceCode=shirt-white-s



