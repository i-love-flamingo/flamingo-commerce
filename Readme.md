# Flamingo Commerce

Contains modules that helps building powerful and flexible ecommerce websites.

Read more under [go.flamingo.me](https://docs.flamingo.me/4.%20Flamingo%20Commerce/1.%20Introduction/About%20Flamingo%20Commerce.html)

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
