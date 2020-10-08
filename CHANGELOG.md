# Changelog
## v3.4.0 [upcoming]
**cart**
* Added desired time to DeliveryForm
* GraphQL
  * Updated schema and resolver regarding desired time

**category**
* Added cue config to module
* Updated documentation of the module regarding the fake service
* FakeService
    * The category fake service was added which can return a project specific category tree and categories
    * Added configuration options are `fakeService.enabled` and `fakeService.testDataFolder` to enable the fake category service and to use json files as fake categories and tree. You can find examples in the documentation of the module

**checkout**
* Checkout Controller, update handling of aborted/canceled payments:
  * Cancel the order / restore the cart before generating the new idempotency key of the payment selection

**product**
* GraphQL
  * **Breaking** New schema for products:
    * `Commerce_Product` has been restructured and now has three subtypes: `Commerce_Product_SimpleProduct`, `Commerce_Product_ConfigurableProduct`, `Commerce_Product_ActiveVariantProduct`
    * Product variant data, that has previously been buried in `Commerce_ConfigurableProduct.variants`, has been mapped to the toplevel of each product and can be accessed directly.
    * Both `ActiveVariantProduct` and `ConfigurableProduct` provide a new property named `variationSelections` which exposes a list of possible attribute combinations for the configurable.
* FakeService
    * The product fake search service is now able to return products with an active variant via `fake_configurable_with_active_variant`. Variation attributes have been changed to only include `color` and `size`.
    * Added configuration option `jsonTestDataFolder` to use json files as fake products. You can find an example product under: `test/integrationtest/projecttest/tests/graphql/testdata/products/json_simple.json`
    * Added fakservice documentation to the product module.  
* Expose `VariantVariationAttributesSorting` on `domain.ConfigurableProductWithActiveVariant`

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
