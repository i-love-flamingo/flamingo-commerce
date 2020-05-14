# 14. June 2018

* Price Fields in Cartitems and Carttotals have been changed:
  * Cartitem:
    * Deleted (Dont use anymore): Price / DiscountAmount / PriceInclTax
    * Now Existing: SinglePrice / SinglePriceInclTax / RowTotal / TaxAmount/ RowTotalInclTax / TotalDiscountAmount / ItemRelatedDiscountAmount / NonItemRelatedDiscountAmount / RowTotalWithItemRelatedDiscount / RowTotalWithItemRelatedDiscountInclTax / RowTotalWithDiscountInclTax
    
  * Carttotal:
    * Deleted: DiscountAmount
    * Now Existing: SubTotal / SubTotalInclTax / SubTotalInclTax /SubTotalWithDiscounts / SubTotalWithDiscountsAndTax / TotalDiscountAmount / TotalNonItemRelatedDiscountAmount 

# 17. April 2019

* Cart Item `UniqueID` is removed
  * `Item.ID` is now supposed to be unique
  * The combination `ID` + `DeliveryCode` is no longer required to identify a cart item
  * For non-unique references of certain backend implementations the new field `Item.ExternalReference` can be used
  
# 23. July 2019
* Add general gift card support
  * `cart.AppliedGiftCards` contains a list of applied gift cards
  * Add convenience functions for gift card like `SumGrandTotalWithGiftCards()` and `HasAppliedGiftCards()`   
  
* Add support for gift cards in default payment selection handling
  * Adds new public function `NewDefaultPaymentSelection` which will generate a basic payment selection
  * Changed visibility of `NewSimplePaymentSelection` to private, please use `NewDefaultPaymentSelection` instead
  * Update ChargeQualifier, add additional Reference string field
  * Add support for multiple charges of the same type (unique Reference needed)
  
# 8. August 2019
* Renamed 'cart.BillingAdress' to 'cart.BillingAddress'

# 15. August 2019  
* Removed `ShippingItem.DiscountAmount` 
  * Added `ShippingItem.AppliedDiscounts`
    * ShippingItem now implements interface `WithDiscount`
    
# 9. October 2019
* Add `PlaceOrderWithCart` to `CartService` to be able to place an already fetched cart instead of triggering an additional call to the `CartReceiverService`

# 18. December 2019
* Introduce `UpdateItems` to `ModifyBehaviour`interface to reduce calls for updating items one by one
* Add helper function `GetDeliveryByItemID`
* Remove `itemID` and `deliveryCode` as parameters for `UpdateItem` as this information is part of the update command, respectively from the new helper

# 7. January 2020
* Add [Idempotency Key pattern](https://stripe.com/blog/idempotency) to the `PaymentSelection`
  * PaymentSelection interface now offers new functions for receiving (`IdempotencyKey()`) / generating (`GenerateNewIdempotencyKey()`) a new Idempotency-Key

# 10. January 2020
* Changes AppliedCouponCodes in the cart to an own struct to be able to add some functions to it
* Quantity item adjustments know also contain a bool that indicates if the respective adjustment caused a change to the AppliedCouponCodes slice of the cart
  * New template function to get if any of the currently stored adjustments caused a coupon code to be removed

# 17. January 2020
* Add `AppliedGiftCard` convenience function `Total()`

# 10. February 2020
* Add `additionalData` to `AddRequest` 
  * Breaking change: Update helper/builder function `BuildAddRequest`
* Breaking Change to `EventPublisher` interface, `PublishChangedQtyInCartEvent` and `PublishAddToCartEvent` now
include a cart as a parameter
* Breaking Change to behaviour of `AddToCartEvent` and `ChangedQtyInCartEvent`, they are now thrown after
the cart has been adjusted and written back to cache
* Events deferred from `ModifyBehaviour` are dispatched before `AddToCartEvent` and `ChangedQtyInCartEvent`
* The `AddToCartEvent` includes the current cart (with added product)
* The `ChangedQtyInCartEvent` includes the current cart (with updated quantities)
* Add Whitebox Test `TestCartService_CartInEvent` to check `AddToCartEvent`

# 20. February 2020

* Mark `CartReceiverService.RestoreCart()` as deprecated, use `CartService.RestoreCart()` instead,
  the cart adapter therefore needs to implement the `CompleteBehaviour` interface.
* Add `CartReceiverService.ModifyBehaviour()` to easily receive the current behaviour (guest/customer)

* Add `CompleteBehaviour` interface which ensures that the cart adapter offers Complete / Restore functionality
* Add `CartService.CompleteCurrentCart()` and `CartService.RestoreCart()` which rely on the new `CompleteBehaviour` interface
* **Breaking**: Update `CartService.CancelOrder()` to use `CartService.RestoreCart()` instead of `CartReceiverService.RestoreCart()`,
  if your cart supports completing/restoring please implement `CompleteBehaviour` interface
* Add `CartService.CancelOrderWithoutRestore()` to allow order cancellation without restoring the cart

* Mark `GuestCartService.RestoreCart` as deprecated, will be replaced by `CompleteBehaviour`
* Mark `CustomerCartService.RestoreCart` as deprecated, will be replaced by `CompleteBehaviour`

* Add mocks for all behaviours, you can use a specific one e.g. `&mocks.CompleteBehaviour{}` or the all in one `&mocks.AllBehaviour{}`

* Update `InMemoryBehaviour` to fulfill the `CompleteBehaviour` interface (adds `Complete()`/`Restore()`)
* Update `InMemoryCartStorage`, add Mutex to be thread safe

* Update `SimplePaymentFormService` to allow gift cards in the `PaymentSelection`, please use the
  config `commerce.cart.simplePaymentForm.giftCardPaymentMethod`to specify the default payment method for gift cards

* Add missing `product` module dependency to cart module

# 13. May 2020

* Update pagination module configuration. Use "commerce.pagination" namespace for configuration now.
* search and product module configuration is using CueConfig - and therefore config options can be looked up in the commandline