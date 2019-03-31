# Customer Module

Provides the basic customer domain, which is required in other e-commerce related packages.

A customer represents the legal entity of a single person and is normally used in a "my account" area.

## Secondary Ports

An implementation needs to provide an adapter for the `CustomerService` as well as an implementation of the `Customer` interface.

A typical implementation would fetch the Token from the auth object and with that information call a microservice that will return the customer data (authorized by the token).

Your specific implementation of a customer can also include much more properties - as long as the two interfaces (ports) are implemented.

### No customer data needed?

You can enable the provided adapter for the customerService with:

```yaml
commerce.customer.useNilCustomerAdapter: true
```
