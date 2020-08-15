# Order Module

The order module offers a domain model for orders to be used to list orders of a customer.

## Usage

### Show orders of a customer
The module offers a data controller "customerorders" to get all orders of the current authenticated user/customer:

`orders = data("customerorders")`

## Ports
The module offers a port that needs to be implemented to fetch customer orders `CustomerIdentityOrderService`.

The module comes with an adapter for the port:
* FakeAdapter: Just returns some dummy orders - useful for local testing

### possible configurations

```yaml
  order:
    # use fake adapter for order fetching
    useFakeAdapter: true
```
