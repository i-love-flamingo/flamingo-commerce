# 20. August 2019
* Add new PaymentAction which processes the payment flow status
* Add option to place an order as early as the payment is started

# 18. December 2019
* Reduce calls for updating items in `SetSourcesForCartItems`

# 7. January 2020
* Generate a new Idempotency Key in the PaymentSelection if an payment error occurs (canceled / aborted by customer) to allow the customer to retry

# 12. February 2020
* Move config to commerce namespace, from `checkout` to `commerce.checkout`
* Add cue based config

# 20. February 2020
* Add `OrderService.CancelOrderWithoutRestore()` which uses the new `CartService` function
* Add `OrderService.CartPlaceOrder()` to place a provided cart instead of fetching it from the `CartService`

* Add new GraphQL Place Order process which relies on a new state machine please referer to the module readme for more details
    * Transition all actions of the checkout controller to separate states
    * Add new `ContextStore` port to provide a storage for the place order process
        * Provide InMemory and Redis Adapter
    * Add new `TryLocker` port to provide an easy way to sync multiple order processes across different nodes
        * Provide InMemory and Redis Adapter
    * Breaking: Add new GraphQL mutations / queries to start / stop / refresh the place order process

# 22. April 2020
* Add coniguration to switch off cart validation on place order
