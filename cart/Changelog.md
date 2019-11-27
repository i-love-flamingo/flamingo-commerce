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