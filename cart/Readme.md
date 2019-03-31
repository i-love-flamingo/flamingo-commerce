# Cart Module

The cart package offers a domain model for carts and its items - and the application services that should be used to get and modify cart(s).

It requires several ports to be implemented in order to have fully working cart functionality.
(for example the *flamingo-commerce-adapter-magento* package implements Adapters for the Ports)

The cart model supports multidelivery and multipayment.

Also the cart module and its services will be used by the checkout module.

## Principles / Summary

* The "Cart" aggregate in the Domain Model is a complex object that should be used as a pure **immutable Value object**:
    * Never change it directly!
    * Only read from it
* The Cart is only **modified by Commands** send to a CartBehaviour Object
* If you want to retrieve or change a  cart - **ONLY work with the application services**
    * This will ensure that the correct cache is used
   

## Usage

### Configurations

```yml
  cart:
    # To register the in memory cart service adapter (e.g. for development mode)
    useInMemoryCartServiceAdapters: true
    enableCartCache: true
    defaultDeliveryCode: "delivery"
    useEmailPlaceOrderAdapter: true
```

## Domain Model Details

### Cart Aggregate

Represents the Cart with PaymentInfos, DeliveryInfos and its Items:

![Cart Model](cart-model.png)


### About Delivery

In order to support Multidelivery the cart cannot directly have Items attached, instead the Items belong to a Delivery.

That also means when adding Items to the cart you need to specify the delivery with a "code".

In cases where you only need one Delivery this can be configured as default and will be added on the fly for you.

#### DeliveryInfo

DeliveryInfo represents the information about which deliverymethod should be used and what deliverylocation should be used.
A DeliveryInfo has:
* a "code" that should idendify a Delivery unique under the cart. Its up to you what code you want. You may want to follow the conventions used by the Default "DeliveryInfoBuilder"
* a "workflow" - that is used to be able to differenciate between different fullfillment workflows (e.g. pickup or delivery)
* a "method" - used to specify details for the delivery. Its up for the project what you want to use. E.g. use it to differenciate between "standard" and "express"
* a "deliverylocation": A deliverylocation can be a address, but also a location defined by a code (e.g. such as a collection point).

The DeliveryInfo object is normaly completed with all required infos during the checkout using the DeliveryInfoUpdateCommand

##### Optional Port: DeliveryInfoBuilder
The DeliveryInfoBuilder interface defines an interface that builds initial DeliveryInfo for a cart.

The "DefaultDeliveryInfoBuilder" that is part of the package should be ok for most cases, it simply takes the passed deliverycode and builds an initial DeliveryInfo object.
The code used by the "DefaultDeliveryInfoBuilder" should be speaking for that reason and is used to initialy create the DeliveryInfo:

The convention used by this default builder is as follow:
`WORKFLOW_LOCATIONTYPE_LOCATIONCODE_METHOD_anythingelse`

Valid codes are:
* delivery (default)
    * DeliveryInfo to have the item (home) delivered
* pickup_store_LOCATIONCODE
    * DeliveryInfo to pickup the item in a (in)store pickup location
* pickup_collection_LOCATIONCODE
    * DeliveryInfo to pickup the item in a a special pickup location (central collectionpoint)


### CartItem details

There are special properties that require some explainations:

* SourceId: Optional represents a location that should be used to fillfill this Item. This can be the code of a certain warehouse or even the code of a retailstore (if item should be picked(sourced) from that location)
    * There is a SourcingService interface - that allows you to register the logic of how to decide on the sourceId

### Decorated Cart

If you need all the product informations at hand - use the Decorated Cart - its decorating the cart with references to the product (dependency product package)

```graphviz
digraph hierarchy {
size="5,5"
node[shape=record,style=filled,fillcolor=gray95]
edge[dir=both, arrowtail=diamond, arrowhead=open]

decoratedCart[label = "{decoratedCart||...}"]
decoratedItem[label = "{decoratedItem||...}"]

product[label = "{BasicProduct|+ ID\n|...}"]

cart[label = "{Cart|+ ID\n|...}"]
item[label = "{Item|ID\nPrice\nMarketplaceCode\nVariantMarketPlaceCode\nPrice\nQty|...}"]

decoratedCart->cart[arrowtail=none]
decoratedItem->item[arrowtail=odiamond]

decoratedItem->product[arrowtail="none"]
decoratedCart->decoratedItem

cart->item[]

}

```

## Details about Price fields

Make sure you read the product package details about prices.

The cart need to show prices with its taxes and additional cart discounts and all different subtotals etc.
What you want to show depends on the project, the type of shop (B2B or B2C), the discount logic and the implemented calculation details of the underlying cartservice implementation.

Some of the prices you may want to show are "calculateable" on the fly (in this cases they are offered as methods) - but some highly depends on the tax and discount logic and they need to have the correct values set by the underlying Cartservice implementation.

### Cart invariants

* While the flamingo Price model (which is used) can calculate exact prices internal, we need "payable prices" to be set in the Cart. This is to allow consistent adding and substracting of prices and still always get a price that is payable.

* All sums - also the Cart GrantTotal is calculated and can be "tracked" back to the item row prices.

In order to get consistent SubTotals in the Cart, the Cart model need certain invariants to be matched:
* Itemlevel: RowPriceNet + TotalTaxAmount = RowPriceGross
* Itemlevel: TotalDiscount <= RowPriceGross
* All main prices need to have same currency (that may be extended later but currently is a constrain)

### About Tax calculation in general

It makes sense to understand the details of tax calculations - so please take the following example of a cart with 2 items:
* item1 netprice 14,71
* item2 netprice 10,18
* and a 19% VAT (G.S.T)

There are two principal ways of tax calculations (called vertical and horizontal)
* vertical: the final tax is calculated from the sum of all rounded item taxes: 
    * item1 +19% gross price rounded: 17,50 (tax 2,79)
    * item2 +19% gross price rounded: 12,11 (tax 1,93)
    * = GrantTotal= 29,61 / =totalTax: 4,72
    * => SubTotalNet = 24,89
    * Often prefered for B2C, when the item prices are shown as gross prices first.
    * Pro: Easier to calculate 
    * Con: rounding errors may sum up in big orders
    
* horizontal: The final tax is calculated from the sum of all net prices of the item and then rounded
    * => SubTotalNet: 24,89  => +19% rounded GrandTotal: 29,62 / included Tax: 4,73
    * Often prefered for B2B or when the prices shown to the customer are as net prices at first.

* In both cases the tax might be calculated from a given net price or a given gross price (see product package config "commerce.product.priceIsGross").

* discounts are normaly subtracted before tax calculation

* At least in germany the law don't force one algorithm over the over. In this cart module it is up to the CartService implementation what algorithm it uses to fill the tax

* Flamingo cart model calculates sums of prices and taxes using the "vertical" way - its doing this by basically adding the prices in the items up on delivery and cart levels. So if you want to use "horizontal" tax calculations the cartservice implementation need to make sure that the item prices are set correct (splitted correct with cent corrections - outgoing from the horizontal sums).



### Cartitems - price fields and method 

The Key with "()" in the list are methods and it is assumed as an invariant, that all prices in an item have the same currency.

| Key                                    | Desc                                                                                                                                                                                                                                                                      | Math Invariants                                                                                                             |
|----------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------|
| SinglePriceGross                       | Single price of product (brutto) (was SinglePrice)                                                                                                                                                                                                                        |                                                                                                                             |
| SinglePriceNet                         | (netto)                                                                                                                                                                                                                                                                   |                                                                                                                             |
| Qty                                    | Qty                                                                                                                                                                                                                                                                       |                                                                                                                             |
| RowPriceGross                          | (was RowTotal )                                                                                                                                                                                                                                                           | RowPriceGross ~ SinglePriceGross * Qty                                                                                      |
| RowPriceNet                            |                                                                                                                                                                                                                                                                           | RowPriceNet ~ SinglePriceNet * Qty                                                                                          |
| RowTaxes                                | Collection of (summed up) Taxes for that item row.                                                                                                                                                                                               |                                                                                                                             |
| TotalTaxAmount()                            | Sum of all Taxes for this Row                                                                                                                                                                                                                                             | = RowPriceGross-RowPriceNet                                                                                                 |
| AppliedDiscounts                       | List with the applied Discounts for this Item  (There are ItemRelated Discounts and Discounts that are not ItemRelated (Cart Related).  However it is important to know that at the end all DiscountAmounts are applied to an item (to make refunding logic easier later) |                                                                                                                             |
| TotalDiscountAmount()                  | Complete Discount for the Row. If the Discounts have no tax/duty (they can be considered as Gross). If they are applied from RowPriceGross or RowPriceNet depends on the calculations done in the cartservice implementation.                                             | TotalDiscountAmount = Sum of AppliedDiscounts TotalDiscountAmount = ItemRelatedDiscountAmount +NonItemRelatedDiscountAmount |
| NonItemRelatedDiscountAmount()         |                                                                                                                                                                                                                                                                           | NonItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = false                                          |
| ItemRelatedDiscountAmount()            |                                                                                                                                                                                                                                                                           | ItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = false                                             |
| RowPriceGrossWithDiscount()            |                                                                                                                                                                                                                                                                           | RowPriceGross-TotalDiscountAmount()                                                                                         |
| RowPriceNetWithDiscount()              |                                                                                                                                                                                                                                                                           |                                                                                                                             |
| RowPriceGrossWithItemRelatedDiscount() |                                                                                                                                                                                                                                                                           | RowPriceGross-ItemRelatedDiscountAmount()                                                                                   |
| RowPriceNetWithItemRelatedDiscount()   |                                                                                                                                                                                                                                                                           |                                                                                                                             |
|                                        |                                                                                                                                                                                                                                                                           |                                                                                                                             |

% use https://www.tablesgenerator.com/markdown_tables to update the table %

### Delivery - price fields and method 
| Key                               | Desc | Math                                       |
|-----------------------------------|------|--------------------------------------------|
| SubTotalGross()                   |      |  Sum of items RowPriceGross                |
| SumRowTaxes()                     |      | List of the sum of the different RowFees   |
| SumTotalTaxAmount()               |      | List of the sum of the TotalTaxAmount   |
| SubTotalNet()                     |      | Sum of items RowPriceNet                   |
| SubTotalGrossWithDiscounts()      |      | SubTotalGross() - SumTotalDiscountAmount() |
| SubTotalNetWithDiscounts          |      | SubTotalNet() - SumTotalDiscountAmount()   |
| SumTotalDiscountAmount()          |      | Sum of items TotalDiscountAmount           |
| SumNonItemRelatedDiscountAmount() |      | Sum of items NonItemRelatedDiscountAmount  |
| SumItemRelatedDiscountAmount()    |      | Sum of items ItemRelatedDiscountAmount     |


### Cart - price fields and method 

| Key                               | Desc                                                                                                                                                | Math                                                                                                                                                     |
|-----------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| GrandTotal()                        | The final amount that need to be payed by the customer (Gross)                                                                                      | = SubTotalGross()  - SumTotalDiscountAmount() + Totalitems |
| Totalitems                        | List of (additional) Totalitems. Each have a certain type - you may want to show some of them in the frontend.      |                                                                                                                                                          |
| SumShipping()                     |           Sum of all shipping costs as Price                                                                                                                                           | The sum of the deliveries shipping                                                                                                                       |
| SubTotalGross()                   |                                                                                                                                                     | Sum of deliveries SubTotalGross()                                                                                                                        |
| SumTaxes()                        | The total taxes of the cart -   as list of Tax                                                                   |                                                                                                                                                          |
| SumTotalTaxAmount()               |            The overall Tax of cart as Price                                                                                       |
| SubTotalNet()                     |                                                                                                                                                     | Sum of deliveries SubTotalNet()                                                                                                                          |
| SubTotalGrossWithDiscounts()      |                                                                                                                                                     | SubTotalGross() - SumTotalDiscountAmount()                                                                                                               |
| SubTotalNetWithDiscounts()          |                                                                                                                                                     | SubTotalNet() - SumTotalDiscountAmount()                                                                                                                 |
| SumTotalDiscountAmount()          | The overall applied discount                                                                                                                        | Sum of deliveries TotalDiscountAmount                                                                                                                    |
| SumNonItemRelatedDiscountAmount() |                                                                                                                                                     | Sum of deliveries NonItemRelatedDiscountAmount                                                                                                           |
| SumItemRelatedDiscountAmount()    |                                                                                                                                                     | Sum of deliveries ItemRelatedDiscountAmount                                                                                                              |
| GetVoucherSavings()               | Returns the sum of Totals from type voucher
| HasShippingCosts()                |  True if cart has in sum any shipping costs

### Typical B2C vs B2B usecases

B2C use cases:
* The item price is with all fees (gros). The discount will be reduced from the price and all the fees, duty and net price will be calculated from this.
* Typically you want to show per row:
    * SinglePriceGross
    * Qty
    * RowPriceGross
    * RowPriceGrossWithDiscount

* The cart normaly shows:
   * SubTotalGross
   * Carttotals (non taxable extra lines on cart level)
   * Shipping
   * The included Total Tax in Cart (SumTaxes) 
   * GrandTotal (is what the customer need to pay at the end inkl all fees)

B2B use cases:
* The item price is without fees (net). The discount will be reduced and then the fees will be added to get the gross price. Do propably you want to show per row:
    *  SinglePriceNet
    *  RowPriceNet
    *  RowPriceNetWithDiscount
    
* The cart then normaly shows:
   * SubTotalNet
   * Carttotals (non taxable extra lines on cart level)
   * Shipping
   * The included Total Tax in Cart (SumTaxes) 
   * GrandTotal (is what the customer need to pay at the end inkl all fees)
 
### Building a Cart with its Deliveries and Items
The domain concept is that the cart is returned and updated by the "Cartservice" - which is the main secondary port of the package that need to be implemented (see below).
The implementations should use the provided "Builder" objects to build up Items, Deliveries and Carts. The Builder items ensure that the invariants are met and help with calculations of missing values.

e.g. to build an item:
```go
builder := t.itemBuilderProvider()
item, err := builder.SetSinglePriceNet(priceDomain.NewFromInt(100, 100, "EUR")).SetQty(10).SetID("22").SetUniqueID("kkk").CalculatePricesAndTaxAmountsFromSinglePriceNet().Build()

``` 
 
### About charges
If you have read the sections above you know about the different prices that are available at item, delivery and cart level and how they are calculated.

There is something else that this cart modell supports - we call it "charges". All the price amounts mentioned is the previous chapters represents the value of the items 
in the cart default currency.

However this value need to be payed - when paying the value it can be that:
- customer wants to pay with different payment methods (e.g. 50% of the value with paypal and the rest with creditcart)
- also the value can be payed in a different currency


The desired split of charges is saved on the cart with the "UpdatePaymentSelection" command.
If you dont need the full flexibility of the charges, than you will simply alway pay one charge that matches the grandtotal of your cart.
Use the factory `NewSimplePaymentSelection` for this.


If you want to use the feature it is important to know how the cart charge split should be generated:

1. the product that is in the cart might require that his price is payed in certain charges. An example for this is products that need to be payed in miles.
2. the customer might want to select a split by himself

You can use the factory on the decorated cart to get a valid PaymentSelection based on the two facts



## Domain - Secondary Ports

### Must Have Secondary Ports

**GuestCartService, CustomerCartService (and ModifyBehavior)**

GuestCartService and CustomerCartService are the two interfaces that act as secondary ports.
They need to be implemented and registered:

```
injector.Bind((*cart.GuestCartService)(nil)).To(infrastructure.YourAdapter{})
injector.Bind((*cart.CustomerCartService)(nil)).To(infrastructure.YourAdapter{})
```

Most of the cart modification methods are part of the `ModifyBehaviour` interface - if you look at the seconary ports you will see, that they need to return an (initialized) implementation of the
`ModifyBehaviour` interface - so in fact this interface need to be implemented when writing an adapter as well.

There is a "InMemoryAdapter" implementation as part of the package.

**PlaceOrderService**

There is also a `PlaceOrderService` interface as secondary port.
Implement an Adapter for it to define what should happen in case the cart is placed.

There is a `EmailAdapter` implementation as part of the package, that sends out the content of the cart as mail.

#### Optional Port: CartValidator

The CartValidator interface defines an interface to validate the cart.

If you want to register an implementation, it will be used to pass the validation results to the web view.
Also the cart validator will be used by the checkout - to make sure only valid carts can be placed as order.

#### Optional Port: ItemValidator

ItemValidator defines an interface to validate an item **BEFORE** it is added to the cart.

If an Item is not valid according to the result of the registered *ItemValidator* it will **not** be added to the cart.

### Store "any" data on the cart

This package offers a flexible way to store any additional objects on the cart:

See this example:

```go
type (
  // FlightData value object
  FlightData struct {
    Direction          string
    FlightNumber       string
  }
)

var (
  //need to implement the cart interface AdditionalDeliverInfo
  _ cart.AdditionalDeliverInfo = new(FlightData)
)

func (f *FlightData) Marshal() ([]byte, error) {
  return json.Marshal(f)
}

func (f *FlightData) Unmarshal(data []byte) error {
  return json.Unmarshal(data, f)
}


//Helper for storing additional data
func StoreFlightData(duc *cart.DeliveryInfoUpdateCommand, flight *FlightData) ( error) {
  if flight == nil {
    return nil
  }
  return duc.SetAdditional("flight",flight)
}

//Helper for getting stored data:
func GetStoredFlightData(d cart.DeliveryInfo) (*FlightData, error) {
  flight := new(FlightData)
  err := d.LoadAdditionalInfo("flight",flight)
  if err != nil {
    return nil,err
  }
  return flight, nil
}
```

## Application Layer

Offers the following services:

* CartReceiverService:
    * Responsible to get the current users cart. This is either a GuestCart or a CustomerCart
    * Interacts with the local CartCache (if enabled)
* CartService
    * All manipulation actions should go over this service (!)
    * Interacts with the local CartCache (if enabled)

Example Sequence for AddToCart Application Services to

![Cart Flow](cart-flow.png)

## A typical Checkout "Flow"

A checkout package would use the cart package for adding informations to the cart, typically that would involve:

* Checkout might want to update Items:
    * Set sourceId (Sourcing logic) on an Item and then call CartBehaviour->UpdateItem(item,itemId)

* Updating DeliveryInformations by Calling CartBehaviour->UpdateDeliveryInfo())
    * (for updating ShippingAddress, WhishDate, ...)

* Optional Updating Purchaser Infos  by Calling CartBehaviour->UpdatePurchaser()

* Finish by calling CartService->PlaceOrder(CartPayment)
    * CartPayment is an object, which holds the informations which Payment is used for which item

## Interface Layer

### Cart Controller

The main Cart Controller expects the following templates by default:

* checkout/cart
* checkout/carterror

The templates get the following variables passed:

* DecoratedCart
* CartValidationResult

### Cart template function
Use the `getCart` templatefunction to get the decorated cart:

```
-
  var decoratedCart = getCart()
  var currentCount = decoratedCart.cart.getCartTeaser.itemCount
```


### Cart Ajax API

There are also of course ajax endpoints, that can be used to interact with the cart directly from your browser and the javascript functionality of your template.

* Get Cart Content: http://localhost:3210/en/api/cart
* Add Simple product: http://localhost:3210/en/api/cart/add/fake_simple
* Add With qty: http://localhost:3210/en/api/cart/add/fake_simple?qty=10
* Adding configurables: http://localhost:3210/en/api/cart/add/fake_configurable?variantMarketplaceCode=shirt-white-s
* Adding configurables with a given delivery: http://localhost:3210/en/api/cart/add/fake_configurable?variantMarketplaceCode=shirt-white-s&deliveryCode=pickup_store
