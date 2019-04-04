# About Flamingo Commerce

With "Flamingo Commerce" and "FlamingoCarotene" you get your toolkit for building fast and flexible commerce experience applications.
Flamingo Commerce contains:

* Flamingomodules for typical e-commerce domains: Each providing a seperated bounded context with its „domain“, „application“ and „interface“ logic.
* Using „ports and adapters“ to seperate domain from technical details, all these modules can be used with your own „Adapters“ to interact with any API or microservice you want.

**Flamingo Commerce is build on top of the FlamingoFramework so it makes sense that you read through the Flamingodocs first**

## A flamingo-commerce project

A typical Flamingo Commerce based e-commerce project would have:

* Its own templates (of course). The templates can be build with the pugtemplate engine and can use the FlamingoCarotene frontend build pipeline.
* Has the dependency to Flamingoand flamingo-commerce packages (using go.mod)
* Project specific Flamingomodules, that provide adapters (=implementation of secondary ports in flamingo-commerce modules) for the commerce modules. This can be adapters that call other microservices and external APIs for example.
* Instead of the specific implementations of adapters, you can also select between available adapters from the open source community:
    * *flamingo-commerce-adapter-standalone* Implementations that work complete without communications to any external service. It provides features to load product data from CSV and keeps cart and checkout in memory. It can be used as a quickstart
    * *flamingo-commerce-adapter-magento* Implementations to use Flamingowith Magento 2.

So a possible e-commerce project build with Flamingo Commerce may look like:

![Logo](./flamingo-commerce-overview.png)

## Possible Architectures

### Flamingo Commerce with Magento2

Flamingo can for example be used as „Head“ for a „Headless“ Magento 2 Setup.

![Logo](./flamingo-magento2.png)

* Products are loaded in memory for better performance on start up
* Cart and Checkout interacts with the Magento2 APIs
* Currently there is still the need to install additional Magento2 extension, in order to expose missing features in Magento2 standard Rest API, that we need for cart features.


### Flamingo Commerce with Magento2 and Elasticsearch

An improved setup will use a proper product-service (e.g. based on elasticsearch) to query and search product data:

![Logo](./flamingo-magento2-es.png)

* Products live in a product-service, that offers a blazing fast and feature rich API to access, search and filter for products. (E.g. based on Elasticsearch)
* Flamingo product and search features use this Elasticsearch service
* Cart and Checkout interacts still with the Magento2 APIs
* In this scenario magento also need an additional extension to load at least basic product data from the product service to use the cart features.


### Flamingo Commerce in a Microservice architecture

A typical commerce based microservice architecture could look like this:

![Logo](./flamingo-ms.png)

In this example we see two different Flamingoprojects:

* one for the core commerce experience - including product search, cart and checkout. This project interacts with microservices like productsearch, cartservice, stockservice and a cms service for example.
* a second one that includes the my account features, that interacts with a order management system for example.

In this szenario also a single sign on solution such as [keycloak](https://www.aoe.com/techradar/tools/keycloak.html) is suggested. Flamingoof course comes with modules supporting the openId Connect flow and OAuth2.0 out of the box.


### Flamingo Commerce standalone

By using the *flamingo-commerce-adapter-standalone* module its also possible to run a Flamingobased webshop like this.

![Logo](./flamingo-standalone.png)

* In this scenario also a SSO solution might be useful, in case you want to support login
* This scenario is not recommended for large scale shops, but might be a possible start.

## Demo

The deome Flamingo Commerce project shows some of the features, and is using  *flamingo-commerce-adapter-standalone*.