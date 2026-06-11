# Sourcing Package

## What do we mean by sourcing?

Sourcing is about finding the "best" possible location
from which a product or an ordered item should be fulfilled.
Sourcing logic is therefore also used to allocate ordered items to the best source locations.

Different factors can play a role when determining the correct source location(s):

- available stock (or replenished stock)
- cost of delivery (e.g. picking warehouses close to the delivery location)
- synergies regarding the complete order  
  (e.g. sourcing the order from a warehouse where most of the items are available)
- delivery time (if time is more important than cost,
  pick the source location that can deliver the fastest)
- ...

## Typical use cases in e-commerce

For your shop, it is helpful to have access to the sourcing logic for advanced use cases such as:

- On PDP:
  - you might restrict the allowed qty based on the available source qty (e.g. a QtyRestrictor that accesses the sourcing logic)
  - you might want to indicate delivery times based on source locations

- During checkout or place order you can access the item allocation to:
  - make sure that a cart can always be sourced (e.g. as part of your CartValidator)
  - show potential packages and delivery times
  - make sure that only carts which can be sourced completely are allowed to be placed
  - attach the source locations for every item to your backend system (e.g. by accessing the sourcing logic in your PlaceOrder adapter)

## About this package

Provides a port for sourcing logic that can be implemented according to your project's needs.

The main port is the `SourcingService` interface, for which you can provide a custom adapter, giving you full freedom to design your sourcing logic.

### Configurations

```yaml
  commerce:
    sourcing:
      # use the DefaultSourcingService (default: true)
      useDefaultSourcingService: true
```

### DefaultSourcingService

The package also offers a `DefaultSourcingService` that performs sourcing based on two inputs:

1. The theoretically available or possible source locations for a given delivery
2. The available stock for specific source locations

For these two inputs, the `DefaultSourcingService` also offers ports where you can provide individual adapters.
Based on this, the `DefaultSourcingService` fetches the possible source locations and sources items based on the available stock at those locations (starting from the first source location retrieved).

### Fake SourcingService

Enabled by adding:
```yaml
commerce:
  sourcing:
    fake:
      enable: true
```

When enabled, it overrides other fake services, and the user has to provide a fake data JSON file via:
```yaml
commerce:
  sourcing:
    fake:
      jsonPath: <your_json_path_here>
```

JSON structure example:
```json
{
  "deliveryCodes": {
    "inflight": 5
  },
  "products": {
    "0f0asdf-0asd0a9sd-askdlj123rw": 10,
    "0f0asdf-0asd0a9sd-askdlj123rx": 15
  }
}
```
