# Sourcing Package

## What do we mean by sourcing

Sourcing is about finding the "best" possible location,
where a product or an ordered item should be fulfilled from.
A sourcing logic is therefore also used to allocate your ordered items to the best source locations.

Different things can play a role while figuring out the correct source location(s):

- available stock (or replenished stock)
- cost of delivery (e.g. pick warehouses close to delivery location)
- synergies regarding the complete order  
  (e.g. source the order from a warehouse where most of the items are available)
- delivery time (if time is more important than cost -
  then picking the source location that can deliver the fastest)
...

## Typical Use cases in e-commerce

For your shop it is helpful to have access to the Sourcing logic for advanced use cases like:

- On PDP:
  - you might restrict the allowed qty based on available source qty (e.g. QtyRestrictor that access Sourcing logic)
  - you might want to indicate delivery times based on source locations

- During Checkout or Place Order you can access the item allocation for:
  - make sure that a cart can always be sourced (e.g. as part of your CartValidator)
  - you might want to show potential packages and delivery time
  - you want to make sure that only carts can be placed that can be Sourced completly
  - you might want to attach the source locations for every item to your backend system (e.g. access the Sourcng logic in your PlaceOrder Adapter)

## About this package

Provides Port for Sourcing logic, that can be implemented according to your project needs.

The main Port is the "SourcingService" interface that you can provide a custom adapter and have all possible freedom to design your sourcing logic.

### Configurations

```yaml
  commerce:
    sourcing:
      # use the DefaultSourcingService (default: true)
      useDefaultSourcingService: true
```

### DefaultSourcingService

The package also offers a "DefaultSourcingService" that does sourcing based on two inputs:

1. The theoretical available or possible sourcelocations for a given delivery
2. The available stock for a specific sourcelocations

For this two inputs the DefaultSourcingService offers also Ports where you can provide individual adapters.
Based on this the DefaultSourcingService fetches the possible sourcelocations and will source items based on the available stock on that locations (starting from the first sourcelocation retrieved).
