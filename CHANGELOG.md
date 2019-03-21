# Changelog

## v3

- general cleanups and linting fixes - that includes several renames of packagenames and types.
- price object introduced:
    - cart and product model don't use float64 anymore but a Price type
    - use commercePriceFormat templatefunc instead (core) priceFormat where you want to render a price object. This will automatically render a "Payable" price.
- cart module:
    - Has a new secondary port: PlaceOrderService
    - The meaning of DeliveryInfo.Method has changed! The former meaning is now represented in the property DeliveryInfo.Workflow. See Readme of cart ackage for details
    - The complete pricefields are changed! Check readme for details on the new price fields and methods
- checkout: 
    - removed depricated viewdata (CartTotals)
- category:
    - Tree object uses a Tree Entity now which contains NOT all category properties. You have to fetch the category details seperate on demand:
        - search for usages of the data funcs - they may need changes in rendering the data: `data('category´´..`