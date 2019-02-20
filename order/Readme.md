# Order Module

The order module offers a domain model for orders to be used to list orders of a customer.

## Usage

### Show orders of a customer
The module offers a data controller "customerorders" to get all orders of the current authenticated user/customer:

`orders = data("customerorders")`

## Ports
The module offers a port that need to be implemented to fetch the orders.

The module comes with 2 possible Adapters for the port:
* FakeAdapter: Just returns some dummy orders - useful for local testing
* EmailAdapters: Used to send mails

### possible configurations

!WIP

```yaml
  order:
    # use fake adapters for order fetching
    useFakeAdapters: true
```
