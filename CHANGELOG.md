# Changelog

## v3.10.0 [upcoming]

**sourcing**
* **Breaking:** Drop deprecated sourcing service `SourcingServiceDetail` and all things related.

**price**
* **Breaking:** Introduce currency library http://github.com/Rhymond/go-money for more flexible rounding. All currency codes should comply to ISO4217 from now on.

## v3.9.0

**search**
* Introduce an Action to search result to advise the frontend to perform an action for that search different from displaying the result page.
  The action model is completely open to anything, but the most common use case is to redirect certain searches directly to product detail pages or special brand landing pages.

**cart**
* Fix: Cart merge of bundle products was broken, bundle configuration is now correctly added to the new customer cart
* Fix: Add to cart of bundle products now only fails for missing required bundle choices. As prior passing all choices with qty of zero have been required. Now optional choices with qty of zero can be omitted.
* Fix: Correctly evaluate min / max quantities of bundle choices. Previously it was possible to add bundle choices with qty 0 even if min qty was 1.
* **Breaking**(In case you have implemented a custom cart service): Extend the cart service interface with `UpdateItemBundleConfig` to allow updating bundles that have already been placed inside the cart.
* GraphQL:
  * Add new mutation `Commerce_Cart_UpdateItemBundleConfig` to update bundle configs for existing cart items
  * **Breaking** Make the qty in the `Commerce_Cart_ChoiceConfigurationInput` type mandatory, previously 1 was taken as a default

**checkout**
* initialize place order metrics with 0 on application start to follow prometheus best practices

**product**
* GraphQL:
  * Expose Active Option (product+qty) for bundle products
  * **Breaking** Make the qty in the `Commerce_Product_ChoiceConfigurationInput` type mandatory, previously it was 0 which lead to taking the minimum required qty of that choice

## v3.8.0

**sourcing**
* Introduce constant `MaxSourceQty` to indicate unlimited stock.

**cart**
* DefaultBehavior: consider discounts when updating cart item row prices
* Make the cart merge strategy configurable, this strategy handles how a guest cart is transformed into a customer cart.
  The following strategies can be set via the config `commerce.cart.mergeStrategy`:
  * `merge` (default): Merge the content of the guest and customer cart
  * `replace`: Replace the customer cart with the guest cart content
  * `none`: Don't do anything, guest cart is lost during customer sign-in.
* Telephone in Address is now deprecated. To distinguish phone number parts new fields were introduced: TelephoneCountryCode, 
  TelephoneAreaCode, TelephoneNumber. Changes are expected to be non-breaking.

## v3.7.0

**cart**
* Updated the `DefaultCartBehaviour` to be able to support other `CartStorage` implementations
    * Dropped the fake implementations of `GiftCardHandler` and `VoucherHandler` to not magically add some discounts

**general**
* Add .golangci.yml and github workflow
* Add .mockery.yaml config

**sourcing**
* **Breaking**
    * SourcingService changed types for return values for bundle product support.
      Every sourcing value is now returned per product id.

## v3.6.0

**cart**
* **Breaking**: Move all calculations to cart behaviour implementation
  * By moving calculation responsibility, we enable different implementation possibilities for calculations like tax before or after discounts, tax on single item or sum and different tax rounding modes instead of having it hard-coded in the flamingo cart.
  * All calculation functions on cart item, shipping item, delivery and cart are now public fields for which the values must be set by the cart behaviour implementation
  * The `DefaultCartBehaviour` calculates all new fields accordingly
  * Removed `ItemBuilder`, `DeliveryBuilder` and `Builder` since they didn't provide any meaningful functionality after removing the calculations. Please create structs directly.
  * Changed the GraphQL cart model accordingly.
  * To help with the migration there are sed commands for the following fields in `cart/migration.sed`: run `find . -type f -iname '*.go' -exec gsed -i -f migration.sed "{}" +;`
  * Cart items
      * | Old Function                           | New Field                            |
        |----------------------------------------|--------------------------------------|
        | RowPriceGrossWithDiscount()            | RowPriceGrossWithDiscount            |
        | RowPriceGrossWithItemRelatedDiscount() | RowPriceGrossWithItemRelatedDiscount |
        | RowPriceNetWithDiscount()              | RowPriceNetWithDiscount              |
        | RowPriceNetWithItemRelatedDiscount()   | RowPriceNetWithItemRelatedDiscount   |
        | TotalDiscountAmount()                  | TotalDiscountAmount                  |
        | ItemRelatedDiscountAmount()            | ItemRelatedDiscountAmount            |
        | NonItemRelatedDiscountAmount()         | NonItemRelatedDiscountAmount         |
  * Shipping items
      * | Old Function               | New Field               |
        |----------------------------|-------------------------|
        | TotalWithDiscountInclTax() | PriceGrossWithDiscounts |
        | -                          | PriceNetWithDiscounts   |
  * Deliveries
      * | Old Function                      | New Field                       |
        |-----------------------------------|---------------------------------|
        | SubTotalGross()                   | SubTotalGross                   |
        | SubTotalNet()                     | SubTotalNet                     |
        | SumTotalDiscountAmount()          | TotalDiscountAmount             |
        | SumSubTotalDiscountAmount()       | SubTotalDiscountAmount          |
        | SumNonItemRelatedDiscountAmount() | NonItemRelatedDiscountAmount    |
        | SumItemRelatedDiscountAmount()    | ItemRelatedDiscountAmount       |
        | SubTotalGrossWithDiscounts()      | SubTotalGrossWithDiscounts      |
        | SubTotalNetWithDiscounts()        | SubTotalNetWithDiscounts        |
        | GrandTotal()                      | GrandTotal                      |
  * Cart
      * | Old Function                      | New Field                       |
        |-----------------------------------|---------------------------------|
        | GrandTotal()                      | GrandTotal                      |
        | -                                 | GrandTotalNet                   |
        | SumShippingNet()                  | ShippingNet                     |
        | SumShippingNetWithDiscounts()     | ShippingNetWithDiscounts        |
        | SumShippingGross()                | ShippingGross                   |
        | SumShippingGrossWithDiscounts()   | ShippingGrossWithDiscounts      |
        | SubTotalGross()                   | SubTotalGross                   |
        | SubTotalNet()                     | SubTotalNet                     |
        | SubTotalGrossWithDiscounts()      | SubTotalGrossWithDiscounts      |
        | SubTotalNetWithDiscounts()        | SubTotalNetWithDiscounts        |
        | SumTotalDiscountAmount()          | TotalDiscountAmount             |
        | SumNonItemRelatedDiscountAmount() | NonItemRelatedDiscountAmount    |
        | SumItemRelatedDiscountAmount()    | ItemRelatedDiscountAmount       |
        | SumAppliedGiftCards()             | TotalGiftCardAmount             |
        | SumGrandTotalWithGiftCards()      | GrandTotalWithGiftCards         |
        | -                                 | GrandTotalNetWithGiftCards      |
* Dispatch a `PreCartMergeEvent` before and a `PostCartMergeEvent` after merging a guest and customer cart when logging in
* Update payment selection on cart merge with the guest carts' payment selection after applying coupons and giftcards and only if customer cart has no items
* Add possibility to have additional data in `AddToCartNotAllowed` error
* Deprecate `BuildAddRequest` in cart service (build your own add request)
* Provide `BaseCartReceiver` for fetching carts without `DecoratedCartFactory` dependency
* Provide `Receiver` interface to be able to mock receiving a cart
* GraphQL:
  * Aggregate many input values to `Commerce_Cart_AddToCart` mutation in `Commerce_Cart_AddToCartInput`
  * Add support for bundle and configurable products in `Commerce_Cart_AddToCart` mutation
* Added support for AdditionalData on CartItem Level for the defaultCartBehaviour.

**checkout**
* Add possibility to have additional data in `PaymentFlowActionTriggerClientSDK`
* Breaking: Disable auto canceling of orders during place order rollback, to restore old behaviour set `commerce.checkout.placeorder.states.placeorder.cancelOrdersDuringRollback: true`

**product**
* Introduce `Labels()` function on `Attribute` to handle translations for attributes with multiple values, will fallback to `Values()` function if not translated.
* Introduce `Stock` slice in `BasicProductData` to store more accurate information about availability for each delivery code
* Introduce new `BundleProduct` and `BundleProductWithActiveChoices` types
* GraphQL:
  * Add `unitCode` to`Commerce_Product_VariationSelection_Option` and `Commerce_Product_ActiveVariationSelection`
  * Fix mapping of VariationSelections
  * Introduce `Labels` for attributes here as well
  * Extend `Commerce_Product` query with non-mandatory bundle configuration argument
  * Add type `Commerce_Product_BundleProduct` that implements `Commerce_Product` interface and is used as bundle product graphql representation.
  * Change `variantionSelections` field in configurable products with `variantSelection` that handles all possible combinations of multi axis configurable.
* FakeService
  * Add configuration option `commerce.product.fakeservice.defaultProducts` which toggles the delivery of default test products. 
  * Add category facet functionality to the fake `SearchService` with default category facet items.
  * Add configuration option `commerce.product.fakeservice.jsonTestDataCategoryFacetItems` which can be used to provide your own category facet items. 
  * Add new `Stock` field to returned `SimpleProduct` from service
  * Add configuration option for delivery codes `commerce.product.fakeservice.deliveryCodes` which can be used to provide different delivery codes for stock.

## v3.5.0
**general**
* Switch to MIT License

**cart**
* Add convenience function to clone carts
* DefaultCartBehaviour now returns real cart clones to prevent data races on cart fields 
* API
  * **Breaking**: Update `DELETE /api/v1/cart` to actually clean the whole cart not only removing the cart items (introduces new route for the previous behaviour, see below)
  * Add new endpoint `DELETE /api/v1/cart/deliveries/items` to be able to remove all cart items from all deliveries but keeping delivery info and other cart data untouched
* Add new method `SumShippingGrossWithDiscounts` to the cart domain which returns gross shipping costs for the cart 
* When using the `ItemSplitter` to split items in items with single qty (`SplitInSingleQtyItems`) the split discounts are reversed to make splitting the row total stable.
* **Breaking**: `SumTotalTaxAmount` now takes taxes on shipping costs into account
* **Breaking**: Delivery discount sum calculations `SumTotalDiscountAmount`, `SumNonItemRelatedDiscountAmount`, `SumItemRelatedDiscountAmount` now take discount on shipping costs into account
  * Old calculation is now in `SumSubTotalDiscountAmount`.
* `CartService`
  * Add `UpdateAdditionalData` to be able to set additional data to cart
  * Add `UpdateDeliveryAdditionalData` to be able to set additional data to the delivery info
  * Introduce new [interface](cart/application/service.go) to be able to easier mock the whole `CartService`
  * Add auto generated mockery mock for the `CartService`
  * Add new field `PriceGross` of `shippingItem` to directly get the shipping cost incl tax (must be filled by cart adapter)
* GraphQL: 
  * Add new method `sumShippingGrossWithDiscounts` to the `Commerce_DecoratedCart` type
  * Add new field `sumShippingGross` to the `Commerce_DecoratedCart` type
  * Add new field `priceGross` to the `Commerce_Cart_ShippingItem` type
  * Add new mutation `Commerce_Cart_UpdateAdditionalData`
  * Add new mutation `Commerce_Cart_UpdateDeliveriesAdditionalData`
  * Add new field `customAttributes` to the `Commerce_CartAdditionalData` type
  * Add new field `additionalData` to the `Commerce_CartDeliveryInfo` type
  * Add new type `Commerce_Cart_CustomAttributes` with method for getting key/value pairs
  * **Breaking**: Make naming convention consistent in graphql schema `Commerce_Cart_*`
  * **Breaking**: Remove the fields `getAdditionalData, additionalDataKeys, additionalDeliveryInfoKeys` from the `Commerce_CartDeliveryInfo` type
  * **Breaking**: `Commerce_Cart_UpdateDeliveryShippingOptions` mutation responded with slice of `Commerce_Cart_DeliveryAddressForm` which was incorrect as we don't process any form data within the mutation. It responds now rightly only with `processed` state.
* **Breaking**: Upgrade github.com/go-playground/form to v4, all types are fully compatible, but import paths have to be changed


**checkout**
* Introducing Flamingo events on final states of the place order process
* Introduce a max ttl for the checkout state machine to avoid polluting the redis with stale checkout processes, defaults to 2h
* Checkout controller: force new order id reservation if an early place happened and there was a payment issue
* API
  * In case of an invalid cart during place order process we now expose the cart validation result, affected endpoints:
    ```
    GET /api/v1/checkout/placeorder
    POST /api/v1/checkout/placeorder/refresh
    POST /api/v1/checkout/placeorder/refresh-blocking
    ```
* Add new Flow Action `PaymentFlowActionTriggerClientSDK` to the checkout
* **Breaking**: Upgrade github.com/go-playground/form to v4, all types are fully compatible, but import paths have to be changed

**customer**
* Add mockery mocks for both `Customer` / `CustomerIdentityService` for easier testing
* Add `State` field to customer address to be closer to cart address type, expose via GraphQL

**price**
* When marshalling `domain.Price` to JSON the amount is rounded.
* Fix various rounding issues with negative prices, add all rounding modes and examples to moduel readme.

**product**
* Enhance the `PriceContext` to allow potential delivery specific pricing
* GraphQL:
  * **Breaking**: Change `activeBase` of `Commerce_Product_PriceInfo` from `Float` to `Commerce_Price`
  * Add `availablePrices` to the `Commerce_Product` interface to display potential pricing options in the frontend
  * Add `context` to the `Commerce_Product_PriceInfo` model to be able to differ between prices

## v3.4.0
**cart**
* Added desired time to DeliveryForm
* InMemoryCartStorage: initialize lock and storage already in Inject() to avoid potential race conditions
* API
  * Add endpoints for deleting / updating a item in the cart (DELETE/PUT: /api/v1/cart/delivery/{deliveryCode}/item)
  * **Breaking**: Affects v1 prefixed routes, switched to a more RESTful naming and use of the correct HTTP verbs to mark idempotent operations
    * | old HTTP verb | old route                                          | new HTTP verb | new route                                  |
      |--------------:|----------------------------------------------------|---------------|--------------------------------------------|
      |          POST | /api/v1/cart/delivery/{deliveryCode}/additem       | POST          | /api/v1/cart/delivery/{deliveryCode}/item  |
      |      POST/PUT | /api/v1/cart/applyvoucher                          | POST          | /api/v1/cart/voucher                       |
      |   POST/DELETE | /api/v1/cart/removevoucher                         | DELETE        | /api/v1/cart/voucher                       |
      |      POST/PUT | /api/v1/cart/applygiftcard                         | POST          | /api/v1/cart/gift-card                     |
      |          POST | /api/v1/cart/applycombinedvouchergift              | POST          | /api/v1/cart/voucher-gift-card             |
      |          POST | /api/v1/cart/removegiftcard                        | DELETE        | /api/v1/cart/gift-card                     |
      |          POST | /api/v1/cart/billing                               | PUT           | /api/v1/cart/billing                       |
      |          POST | /api/v1/cart/delivery/{deliveryCode}/deliveryinfo  | PUT           | /api/v1/cart/delivery/{deliveryCode}       |
      |           PUT | /api/v1/cart/updatepaymentselection                | PUT           | /api/v1/cart/payment-selection             |
  
* GraphQL
  * Update schema and resolver regarding desired time

**category**
* Added cue config to module
* Update fake service documentation
* FakeService
    * The category fake service was added which can return a project specific category tree and categories
    * Added configuration options are `fakeService.enabled` and `fakeService.testDataFolder` to enable the fake category service and to use json files as fake categories and tree. You can find examples in the documentation of the module

**checkout**
* Checkout Controller, update handling of aborted/canceled payments:
  * Cancel the order / restore the cart before generating the new idempotency key of the payment selection
* Resolve goroutine leak in redis locker
* **Breaking** Change StartPlaceOrder behaviour to always start a new one.
* **Breaking** Change ClearPlaceOrderProcess behaviour to be always possible (no matter on which state)
* API
  * **Breaking**: Affects v1 prefixed routes, switched to a more RESTful naming and use of the correct HTTP verbs to mark idempotent operations
    * | old HTTP verb | old route                                     | new HTTP verb | new route                                     |
      |--------------:|-----------------------------------------------|---------------|-----------------------------------------------|
      |          POST | /api/v1/checkout/placeorder/refreshblocking   | POST          | /api/v1/checkout/placeorder/refresh-blocking  |

**customer**
* GraphQL
  * Extend `Commerce_Customer_Address` with some useful fields
  * Extend `Commerce_Customer_Result` with a field for querying a specific address
  * **Breaking**:
    * `Commerce_Customer_Address`: rename field `StreetNr` to `StreetNumber`, `lastname` to `lastName` and `firstname` to `firstName` 
    * `Commerce_Customer_Result`: `defaultShippingAddress` and `defaultBillingAddress` now can return null if there is no default address
    * `Commerce_Customer_PersonData`: field `birthday` is now nullable and of type `Date`. 

**payment**
* Introduced wallet payment method 

**product**
* Add support for product badges
* GraphQL
  * **Breaking** New schema for products:
    * `Commerce_Product` has been restructured and now has three subtypes: `Commerce_Product_SimpleProduct`, `Commerce_Product_ConfigurableProduct`, `Commerce_Product_ActiveVariantProduct`
    * Product variant data, that has previously been buried in `Commerce_ConfigurableProduct.variants`, has been mapped to the toplevel of each product and can be accessed directly.
    * Both `ActiveVariantProduct` and `ConfigurableProduct` provide a new property named `variationSelections` which exposes a list of possible attribute combinations for the configurable.
* FakeService
    * The product fake search service is now able to return products with an active variant via `fake_configurable_with_active_variant`. Variation attributes have been changed to only include `color` and `size`.
    * Added configuration option `jsonTestDataFolder` to use json files as fake products. You can find an example product under: `test/integrationtest/projecttest/tests/graphql/testdata/products/json_simple.json`
    * Added fakservice documentation to the product module.
    * The product fake search service is now able to return a specific product if the given query matches the marketplace code / name of the json file of the product  
    * The product fake search service returns no products if it is queried with `no-results`
* Expose `VariantVariationAttributesSorting` on `domain.ConfigurableProductWithActiveVariant`
* **Breaking**: Update stock handling: Remove magic alwaysInStock product attribute. Just Rely BasicProductData.StockLevel field instead.

**sourcing**
* **Breaking**
  * Optional pointer `DeliveryInfo` added as parameter to `StockProvider.GetStock`
  
**docs**
* Embed swagger.json via go-bindata, so it can be used from the outside


## v3.3.0
**product**
* Switch module config to CUE
* Extended product model with loyalty earnings
* Added Rest API route to get products
* GraphQL
  * Added values field of Attribute to schema
  * exposed loyalty earnings
  * Added facets to fake search service
* Expose unit of product variant attributes
* fake: add loyalty pricing for `fake_simple` product, introduced `fake_fixed_simple_without_discounts` product

**cart**
* **Breaking**: Moved to new flamingo auth module, breaks several interfaces which previously relied on the oauth module
* **Breaking**
  * Cart item validation now requires the decorated cart to be passed to assure that validators don't rely on a cart from any other source (e.g. session)
  * Session added as parameter to interface method `MaxQuantityRestrictor.Restrict`
  * Session added as parameter to `RestrictionService.RestrictQty`
  * Changed no cache entry found error for cartCache `Invalidate`, `Delete` and `DeleteAll` to `ErrNoCacheEntry`
* Switch module config to CUE
* GraphQL
  * Add new mutation to set / update one or multiple delivery addresses `Commerce_Cart_UpdateDeliveryAddresses`
  * Add new mutation to update the shipping options (carrier / method) of an existing delivery `Commerce_Cart_UpdateDeliveryShippingOptions`
  * Add new mutation to clean current users cart `Commerce_Cart_Clean`
  * Add new query to check if a product is restricted in terms of the allowed quantity `Commerce_Cart_QtyRestriction`
  * Add new field `sumPaymentSelectionCartSplitValueAmountByMethods` to the `Commerce_Cart_Summary` which sums up cart split amounts of the payment selection by the provided methods.   
  * Expose PaymentSelection.CartSplit() via GraphQL, add new types `Commerce_Cart_PaymentSelection_Split` and `Commerce_Cart_PaymentSelection_SplitQualifier`
  * **Breaking**: renamed the following GraphQL types
    * type `Commerce_Cart_BillingAddressFormData` is now `Commerce_Cart_AddressForm`
    * input `Commerce_BillingAddressFormInput` is now `Commerce_Cart_AddressFormInput`
    * type `Commerce_Charge` is now `Commerce_Price_Charge`
    * type `Commerce_ChargeQualifier` is now `Commerce_Price_ChargeQualifier`
    * input `Commerce_ChargeQualifierInput` is now `Commerce_Price_ChargeQualifierInput`
* Adjusted log level for cache entry not found error when trying to delete the cached cart   

**customer**
* **Breaking**: renamed `GetId` to `GetID` in `domain.Customer` interface
* introduced new `CustomerIdentityService` to retrieve authenticated customers by `auth.Identity` 
* **Breaking**: removed `CustomerService` please use `CustomerIdentityService`
* GraphQL: Add new customer queries:
  * `Commerce_Customer_Status` returns the customer's login status
  * `Commerce_Customer` returns the logged-in customer

**checkout**
* Deprecate Sourcing service port in checkout (activate if required with setting `commerce.checkout.activateDeprecatedSourcing`)
* Make cart validation before place order optional with configuration
* State Machine
  * Add additional metrics to monitor place order flow
    * flamingo_commerce_checkout_placeorder_starts
    * flamingo_commerce_checkout_placeorder_state_run_count
    * flamingo_commerce_checkout_placeorder_state_failed_count
   * Add a step to validate the payment selection if needed. The step provides a port to be implemented if needed.
* Expose placeorder endpoints also via rest
* Checkout controller, update to the error handling:
  * In case of a payment error the checkout controller will now redirect to the checkout/review action instead of just rendering the matching template on the current route.
  * Same applies in case of an error during place order, the checkout controller will now redirect to the checkout step.
  * In both cases the error will be stored as a flash message in the session before redirecting, the target action will then receive it and pass it to the template view data. 
 
**search**
* Switch module config to CUE
* Update `pagination` module configuration. Use `commerce.pagination` namespace for configuration now.
* GraphQL
  * **Breaking** simplified Commerce_Search_SortOption type
  * **Breaking** use GraphQL specific search result type
  * Added facets resolver
  
**category**
* GraphQL
  * **Breaking** moved GraphQL `dto` package to `categorydto`

**docs**
* Add Swagger/OpenAPI 2.0 specification to project, using [swaggo/swag](https://github.com/swaggo/swag)
**sourcing**
* Add new "sourcing" module that can be used standalone. See sourcing/Readme.md for more details

**w3cdatalayer**
* **Breaking**: Switch from flamingo `oauth` module to the new `auth` module, to keep w3cdatalayer working please configure the new `auth` module accordingly

**order**
* **Breaking**: removed interface `CustomerOrderService` please use `CustomerIdentityOrderService`
* Update config path: `order.useFakeAdapters` to `commerce.order.useFakeAdapter`

## v3.2.0
**w3cdatalayer**
* Fixed a bug that causes the datalayer to panic if it failed to build an absolute url
* Introduced a configuration option to choose between base64url and hex encoding for the hashed values
* Move config to commerce namespace, from `w3cDatalayer` to `commerce.w3cDatalayer`
* Add legacy config mapping so old mappings can still be used
* Add cue based config to have config validation in place

**checkout**
* Controller
  * Allow checkout for fully discounted carts without payment processing. Previously all checkouts needed a valid payment to continue.
    In case there is nothing to pay this can be skipped.
    * Order ID will be reserved as soon as the user hits the checkout previously it was done before starting the payment
* GraphQL
  * Update place order process to also allow zero carts which don't need payment, this leads to a state flow that lacks the payment steps.
    See module readme for further details.
* Update source service to support external location codes.
  * Adds `ExternalLocationCode` to the `Source` struct.
  * Update `SetSourcesForCartItems()` to use the new `SourcingServiceDetail` functionality if the bound service implements the interface
* Update `OrderService` to expose more metrics regarding the place order process:
    ```
    flamingo-commerce/checkout/orders/cart_validation_failed
    flamingo-commerce/checkout/orders/no_payment_selection
    flamingo-commerce/checkout/orders/payment_gateway_not_found
    flamingo-commerce/checkout/orders/payment_flow_status_error
    flamingo-commerce/checkout/orders/order_payment_from_flow_error
    flamingo-commerce/checkout/orders/payment_flow_status_failed_canceled
    flamingo-commerce/checkout/orders/payment_flow_status_aborted
    flamingo-commerce/checkout/orders/place_order_failed
    flamingo-commerce/checkout/orders/place_order_successful
    ```
  
**cart**
* inMemoryBehaviour: Allow custom logic for GiftCard / Voucher handling
  * We introduced two new interfaces `GiftCardHandler` + `VoucherHandler`
  * This enables users of the in-memory cart to add project specific gift card and voucher handling 
* Fix `CreateInitialDeliveryIfNotPresent` so that cache gets updated now when an initial delivery is created
* GraphQL: Add new cart validation query `Commerce_Cart_Validator` to check if cart contains valid items

**price**
* IsZero() now uses LikelyEqual() instead of Equal() to avoid issues occurring due to floating-point arithmetic

**product**
* product attributes:
  * Added `AttributesByKey` domain method to filter attributes by key and exposed this method as `getAttributesByKey` in GraphQL
  * GraphQL: Exposing `codeLabel` property in the `Commerce_ProductAttribute` type
  
**payment**
* Introduced error message for already used idempotency key 

## v3.1.0
**tests**
* Added GraphQL integration tests for new Place Order Process, run manually with `make integrationtest`
* To run the GraphQL Demo project use `make run-integrationtest-demo-project`
* To regenerate the GraphQL files used by the integration tests / demo project use  `make generate-integrationtest-graphql`

**cart**
* Add `additionalData` to the `AddRequest` used during add to cart
    * Breaking: Update helper/builder function `BuildAddRequest`
* Breaking: Change to `EventPublisher` interface, `PublishChangedQtyInCartEvent` and `PublishAddToCartEvent` now
include a cart as a parameter
* Breaking: Change to behaviour of `AddToCartEvent` and `ChangedQtyInCartEvent`, they are now thrown after
the cart has been adjusted and written back to cache
* Events deferred from `ModifyBehaviour` are dispatched before `AddToCartEvent` and `ChangedQtyInCartEvent`
* The `AddToCartEvent` includes the current cart (with added product)
* The `ChangedQtyInCartEvent` includes the current cart (with updated quantities)

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

**checkout**
* Move config to commerce namespace, from `checkout` to `commerce.checkout`
* Add legacy config mapping so old mappings can still be used
* Add cue based config to have config validation in place

* Add `OrderService.CancelOrderWithoutRestore()` which uses the new `CartService` function
* Add `OrderService.CartPlaceOrder()` to place a provided cart instead of fetching it from the `CartService`

* Add new GraphQL Place Order process which relies on a new state machine please referer to the module readme for more details
    * Transition all actions of the checkout controller to separate states
    * Add new `ContextStore` port to provide a storage for the place order process
        * Provide InMemory and Redis Adapter
    * Add new `TryLocker` port to provide an easy way to sync multiple order processes across different nodes
        * Provide InMemory and Redis Adapter
    * Breaking: Add new GraphQL mutations / queries to start / stop / refresh the place order process
    
**payment**
* Add `PaymentService` to easily work with bound PaymentGateway's
    * `PaymentService.AvailablePaymentGateways()` returns all bound gateways
    * `PaymentService.PaymentGateway()` gets the payment gateway by gateway code
    * `PaymentService.PaymentGatewayByCart()` gets the payment gateway of the cart payment selection

* Extend the `FlowStatus` struct with more standardized `FlowActionData`
* Add standardized Flow Actions `PaymentFlowActionShowIframe`, `PaymentFlowActionShowHTML`, `PaymentFlowActionRedirect`,
  `PaymentFlowActionPostRedirect` please use these in your payment adapter since the standard place order relies on them.
 
**search**
* Extend `Suggestion` struct with `Type` and `AdditionalAttributes` to be able to distinguish between product/category suggestions

## v3.0.1
- Update dingo and form dependency to latest version

## v3.0.0
- general cleanups and linting fixes - that includes several renames of packagenames and types.
- price object introduced:
    - cart and product model don't use float64 anymore but a Price type
    - use commercePriceFormat templatefunc instead (core) priceFormat where you want to render a price object. This will automatically render a "Payable" price.
- cart module:
    - Has a new secondary port: PlaceOrderService
    - The meaning of DeliveryInfo.Method has changed! The former meaning is now represented in the property DeliveryInfo.Workflow. See Readme of cart ackage for details
    - The complete pricefields are changed! Check readme for details on the new price fields and methods
- checkout: 
    - removed depricated viewdata (CartTotals)
- products:
    - product category breadcrumb is not filled in controller - if you want a breadcrum you can use category data functions
    - product category fields are changed to use a categoryTeaser
- category:
    - Tree object uses a Tree Entity now which contains NOT all category properties. You have to fetch the category details separate on demand:
        - search for usages of the data funcs - they may need changes in rendering the data: `data('category´´..`
