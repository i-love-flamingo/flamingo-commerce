# Checkout Package

This package provides a One Page standard checkout with the following features:

* Concept for PaymentProviders, that can be used to implement specific Payments
* A "Offline payment" Provider is part of the module

## Standard Checkout Flow

This module implements controller and services for the following checkout flow (the checkout process for the end customer):

1. check if user is logged in
    * yes: next action
    * no: show start template (where the user can login)
2. Show checkout for user or for guest
    * this is the main step, and the template should render the big checkout form (or at least the parts that are interesting)
    * on submit and if everything was valid and update of the cart was sucessfull - the next action is called
3. Show Review Page
4. Check and Process Payment: If the payment requires a redirect to an external PSP (hosted payment page) the redirect is done and the results are proccessed.
5. Submit Order and show success page

## Configurations

If your template does not want to ask for all the informations required you can also set default values for the checkoutform (strings)

```yml
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

  # Also the checkout form can be configured:
  checkoutForm:
    defaultValues:
      billingAddress_phoneAreaCode: "03722"
    additionalFormValues:
      - lp_membership_id
      - newsletter_opt_in
    overrideValues:
      billingAsShipping: true
      billingAddress_countryCode: DE
      billingAddress_street: NoStreet
      billingAddress_streetNr: "0"
      billingAddress_city: NoCity
      billingAddress_postCode: "99999"
```

## Registering own Payment Providers

You need to implement the secondary port "PaymentProvider" and register the Payment provider in your module.go:

```go
func (pp *PaymarkProvider) Configure(injector *dingo.Injector) {
  injector.BindMap((*payment.PaymentProvider)(nil), pp.ProviderCode).To(infrastructure.PaymarkAdapter{})
}
```
