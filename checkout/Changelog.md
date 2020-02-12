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
