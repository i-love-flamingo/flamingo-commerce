## Checkout Package

This package provides a One Page standard checkout.

### Configurations:

If your template does not want to ask for all the informations required you can also set default values for the checkoutform (strings)

```
checkout:
  checkoutForm:
    defaultValues:
      billingAddress_phoneAreaCode: "03722"
    overrideValues:
      billingAsShipping: true
      billingAddress_countryCode: DE
      billingAddress_street: NoStreet
      billingAddress_streetNr: "0"
      billingAddress_city: NoCity
      billingAddress_postCode: "99999"
```
