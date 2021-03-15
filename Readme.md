
# Flamingo Commerce
[![Go Report Card](https://goreportcard.com/badge/github.com/i-love-flamingo/flamingo-commerce)](https://goreportcard.com/report/github.com/i-love-flamingo/flamingo-commerce)
[![Tests](https://github.com/i-love-flamingo/flamingo-commerce/workflows/Tests/badge.svg?branch=master)](https://github.com/i-love-flamingo/flamingo-commerce/actions?query=branch%3Amaster+workflow%3ATests)
[![Release](https://img.shields.io/github/release/i-love-flamingo/flamingo-commerce?style=flat-square)](https://github.com/i-love-flamingo/flamingo-commerce/releases)


With "Flamingo Commerce" you get your toolkit for building fast and flexible commerce experience applications.

A demoshop using the standalone adapters is online here https://demoshop.flamingo.me - you can also try the [GraphQL](https://demoshop.flamingo.me/en/graphql-console) support
 
## What problems does Flamingo Commerce solve?

* Modern Architecture: Break monolithic e-commerce architeture to allow scaling and maintainability. 
* Modern Architecture: Use it to build commerce for headless commerce solutions
* Real time commerce: Build personalized experiences - without the need to cache rendered pages

## What are the main design goals of Flamingo Commerce?


* **Performance**: We do not want to rely on any frontend caching. Instead it is no problem to show every customer an individual experience.
* **Clean architecture**: We use "domaind driven design" and "ports and adapters" to build a maintainable and clean application. 
* **Suiteable for Microservice architectures**: Adapters concept and various resilience concepts makes it easy to connect to other (micro) services.
* **Decoupled and flexible frontend development**: Frontend development is decoupled from the "Backend for Frontend" - so that it is possible to use "any" frontend technology.
* **Testability**: By providing "Fake Adapters" that provide test data, it is possible to test your application without external dependencies.
* **Great Developer Experience**: :-)
* **Open Source**: Flamingo Commerce and Flamingo is Open Source and will be.


## Whats does Flamingo Commerce provide?

* Different e-commerce Flamingo Modules for typical e-commerce domains: Each providing a separated bounded context with its „domain“, „application“ and „interface“ logic.
* Using „ports and adapters“ to separate domain from technical details, all these modules can be used with your own „Adapters“ to interact with any API or microservice you want.
* Some of the major Flamingo Commerce modules (bounded contexts) are:
    * product: Offering domain models for different product types. Supporting multiple prices (including loyalty prices) etc..
    * cart: Powerful cart domain model. Supporting multi delivery, multiple payment transactions, and a lot more. 
    * search: Generic search service and features
    * checkout: Offering logic and interfaces for an example (default) checkout.
    
* Each of the modules provide graphql support that you can use.
    

**Flamingo Commerce is build on top of the Flamingo Framework so it makes sense that you read through the Flamingo docs also**

Read more under [docs.flamingo.me](https://docs.flamingo.me/1.%20Introduction/1.%20Getting%20Started.html)

## Commerce Modules:

* **price**: 
    * Offers value objects for prices and charges - supporting calculations, rounding and splitting
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/price?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/price)
    * [Readme](price/Readme.md)
* **product**: 
    * Offers domain models and interface logic for handling different product types
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/product?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/product) 
    * [Readme](product/Readme.md)
* **category**: 
    * Offers domain models and interface logic for category tree and category views
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/category?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/category) 
    * [Readme](category/Readme.md)
* **cart**: 
    * The cart module is one of the main modules in Flamingo Commerce. It offers domain models and logic for multi delivery, multi payment carts.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/cart/domain/cart?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/domain/cart) 
    * [Readme](cart/Readme.md)
* **payment**: 
    * Offers a generic payment value objects as well as a generic web payment interface and comes with the "offlinepayment" gateway.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/payment/domain?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/payment/domain) 
    * [Readme](payment/Readme.md)
* **search**: 
    * Offers domain models and interface logic for generic search and search filters.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/search/domain?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/search/domain) 
    * [Readme](search/Readme.md)
* **checkout**: 
    * Offers a default checkout implementation that can be used.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/checkout?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/checkout) 
    * [Readme](checkout/Readme.md)
* **customer**: 
    * Offers domain models for customer
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/customer/domain?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/customer/domain) 
    * [Readme](customer/Readme.md)
* **order**: 
    * Offers domain models for orders. For example to use it on a "My Orders" page.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/order/domain?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/order/domain) 
    * [Readme](order/Readme.md)

* **w3cdatalayer**: 
    * Offers interface logic to render a Datalayer that can be used for e-commerce tracking
    * [Readme](w3cdatalayer/Readme.md)
    
# Flamingo Commerce Release Status

Flamingo Commerce API is Beta because we will still change the API (models and methods).
You are encourages to use it but if you update you might need to adjust your code to the latest changes.


## Setup

We recommend to use go modules, so you just need to add Flamingo Commerce to your main go file as import:

e.g. to use the product module add

```go
import (
  "flamingo.me/flamingo-commerce/v3/product"
)
```

And then load the Module in your application bootstrap:

```go

// main is our entry point
func main() {

	flamingo.App([]dingo.Module{
	    ...
		//flamingo-commerce modules
		new(product.Module),
		
	}, nil)
}


```

To update the dependency in go.mod run

```
go get flamingo.me/flamingo-commerce/v3
```
## Demo 

There is a demo: https://demoshop.flamingo.me

And the code is also published: https://github.com/i-love-flamingo/commerce-demo-carotene
