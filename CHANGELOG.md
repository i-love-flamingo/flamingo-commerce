# Changelog

## v3

- general cleanups and linting fixes - that includes several renames of packagenames and types.
- price object introduced:
    - cart and product model don't use float64 anymore but a Price type
    - commercePriceFormat
- cart module:
    - Has a new secondary port: PlaceOrderService
    - The meaning of DeliveryInfo.Method has changed! The former meaning is now represented in the property DeliveryInfo.Workflow. See Readme of cart ackage for details
- checkout: 
    - removed depricated viewdata (CartTotals)