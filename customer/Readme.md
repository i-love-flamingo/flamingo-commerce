# Customer Modul

Provides the basic customer domain, which is required in other ecommerce related packages.

A customer represents the legal entity of a single person and is normaly used in a "my account" area.

## Secondary Ports

A implementation need to provide an adpater for the `CustomerService`as well as a implementation of the `customer` interface.

A typical implementation would fetch the Token from the auth object and with that information call a microservice that will return the customer data (authorized by the token).

Your specific implementation of a customer can also include much more properties - as long as the two interfaces (ports) are implemented.

### No customer data needed?
You can enable the provided adapter for the customerService with:
```
commerce.customer.useNilCustomerAdapter: true
```
