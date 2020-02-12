# Checkout Package

This package provides a one page standard checkout with the following features:

* Concept for PaymentProviders, that can be used to implement specific Payments
* An "offline payment" provider is part of the module

## Standard checkout flow

This module implements controller and services for the following checkout flow (the checkout process for the end customer):

1. StartAction (optional)
    * check if user is logged in
        * yes: next action
        * no: show start template (where the user can login)
1. Checkout Action
    * this is the main step, and the template should render the big checkout form (or at least the parts that are interesting). 
    * on submit and if everything was valid:
        * Action will update the cart - specifically the following information:
            * Update billing (on cart level)
            * Update deliveryinfos on (each) delivery in the cart (e.g. used to set shipping address)
            * Select Payment Gateway and preferred Payment Method
            * Optional Save the wished payment split for each item in the cart
            * Optional add Vouchers/GiftCards (may already happened before)            
        * If Review Action is skipped:
            * Start payment and place order if needed (EarlyPlaceOrder)
            * Redirect to Payment Action
        * If Review Step is not skipped:
            * Redirect to Review Action
1. Review Action (Optional)
    * Renders "review" template that can be used to review the cart
    * After confirming start the payment and place order if needed (EarlyPlaceOrder)
1. Payment Action
    * Ask Payment Gateway about FlowStatus and handle it
    * FlowStatus:
        * Error / Abort by customer: Regenerate Idempotency Key of PaymentSelection, redirect to checkout and reopen cart if needed
        * Success / Approved: Redirect to PalceOrderAction
        * Unapproved: Render payment template and let frontend decide how to continue in flow (e.g. redirect to payment provider)
1. Place Order Action
    * Check if order already placed (EarlyPlaceOrder)
    * If order not already placed check FlowStatus and place order
    * Put order infos in flash message and redirect to Success Action
1. Success Action:
    * Renders order success template

## Configurations

If your template does not want to ask for all the information required you can also set default values for the checkoutform (strings)

```yaml
commerce:
  checkout:
    # use a faked sourcing service
    useFakeSourcingService: false
    # to enable the offline payment provider
    enableOfflinePaymentProvider: true

    # checkout flow control flags:
    skipStartAction: false
    skipReviewAction: false
    showReviewStepAfterPaymentError: false
    showEmptyCartPageIfNoItems: false
    redirectToCartOnInvalideCart: false

    # checkout form settings:
    useDeliveryForms:                 true
	usePersonalDataForm:              false
	privacyPolicyRequired:           true
```


## Sourcing Service Secondary Ports
There is the an optional secondary port provided, that we call "Sourcing Service".
The Sourcing service is responsible for assigning an Item in the cart the correct source location. The source location is the location where the item should be fulfilled from. Typically a warehouse.

By providing an adapter for this port you can control the source locations for the items in your cart.
