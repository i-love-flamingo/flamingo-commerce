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
        * Action will update the cart - specifically the following informations:
            * Update billing (on cart level)
            * Update deliveryinfos on (each) delivery in the cart (e.g. used to set shipping address)
            * Select Payment Gateway and prefered Payment Method
            * Optional Save the wished payment split for each item in the cart
            * Optional add Vouchers (may already happend before)            
        * Forward to next Action
1. Review Action (Optional)
    * Renderes "review" template that can be used to review the cart
    * After confirming:
        * Action will give control to the selected PaymentFlow (see payment module)
1. Place Order Action
    * Get PaymentFlow Result
    * Place Order and forward to success
1. Success Action:
    * Renders order success template

## Configurations

If your template does not want to ask for all the information required you can also set default values for the checkoutform (strings)

```yaml
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

```
