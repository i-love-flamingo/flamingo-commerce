
# Flamingo Commerce
[![Go Report Card](https://goreportcard.com/badge/github.com/i-love-flamingo/flamingo-commerce)](https://goreportcard.com/report/github.com/i-love-flamingo/flamingo-commerce) [![Build Status](https://travis-ci.org/i-love-flamingo/flamingo-commerce.svg)](https://travis-ci.org/i-love-flamingo/flamingo-commerce)

Contains modules that helps building powerful and flexible ecommerce websites.

Read more under [docs.flamingo.me](https://docs.flamingo.me/4.%20Flamingo%20Commerce/1.%20Introduction/About%20Flamingo%20Commerce.html)

## Commerce Modules:

* **price**: 
    * Offers value objects for prices and charges - supporting calculations, rounding and splitting
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/price?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/price)
    * [Readme](`docs/2. Flamingo Commerce Modules/10. Price.md`)
* **product**: 
    * Offers domain models and interface logic for handling different product types
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/product?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/product) 
    * [Readme](`docs/2. Flamingo Commerce Modules/01. Product Module.md`)
* **category**: 
    * Offers domain models and interface logic for category tree and category views
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/category?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/category) 
    * [Readme](`docs/2. Flamingo Commerce Modules/02. Category Module.md`)
* **cart**: 
    * The cart module is one of the main modules in Flamingo Commerce. It offers domain models and logic for multi delivery, multi payment carts.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/cart/domain/cart?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/domain/cart) 
    * [Readme](`docs/2. Flamingo Commerce Modules/03. Cart Module.md`)
* **payment**: 
    * Offers a generic payment value objects as well as a generic web payment interface and comes with the "offlinepayment" gateway.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/payment/domain?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/payment/domain) 
    * [Readme](`docs/2. Flamingo Commerce Modules/11. Payment.md`)
* **search**: 
    * Offers domain models and interface logic for generic search and search filters.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/search/domain?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/search/domain) 
    * [Readme](`docs/2. Flamingo Commerce Modules/07. Search Module.md`)
* **checkout**: 
    * Offers a default checkout implementation that can be used.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/checkout?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/checkout) 
    * [Readme](`docs/2. Flamingo Commerce Modules/04. Checkout Module.md`)
* **customer**: 
    * Offers domain models for customer
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/customer/domain?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/customer/domain) 
    * [Readme](`docs/2. Flamingo Commerce Modules/06. Customer Module.md`)
* **order**: 
    * Offers domain models for orders. For example to use it on a "My Orders" page.
    * [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/order/domain?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo-commerce/order/domain) 
    * [Readme](`docs/2. Flamingo Commerce Modules/05. Order Module`)

* **w3cdatalayer**: 
    * Offers interface logic to render a Datalayer that can be used for e-commerce tracking
    * [Readme](`docs/2. Flamingo Commerce Modules/08. W3C Datalayer.md`)
    
# Flamingo Commerce in Beta

Flamingo Commerce is Beta because we will still change the API (models and methods).
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
