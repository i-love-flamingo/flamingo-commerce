# W3C Datalayer Module

Modul that makes it easy to add common datalayer to your website - which makes it easier to connect to analytics and to implement tracking pixels and tagmanagers

The datastructure is oriented at:
https://www.w3.org/community/custexpdata/

The datalayer informations provides an easy access to common informations relevant for tracking and datalayer.
The datalayer module therefore listens to various events and adds informations to the datalayer


## Configurations:
```
w3cDatalayer:
  pageInstanceIDPrefix: "mywebsite"
  pageInstanceIDStage: "%%ENV:STAGE%%production"
  pageNamePrefix: My Shop
  siteName:  My Shop
  defaultCurrency: GBP
  version: 1.0
  //If you want sha512 hashes instead real user values
  hashUserValues: true
```

Also it reuse the configuration from locale package to extract the language:
```
locale:
  locale: en-gb
``` 


## Usage example:

### Templatefunc `w3cDatalayerService`

The templatefunc provides you access to the current requests datalayer.
You can get the datalayer and you can modify it:

For some values in the datalayer the template knows better than the backend what to put in, so please call the approriate setter like this:
```
  - var result = w3cDatalayerService().setPageCategories("masterdata","brand","detail")
  - var result = w3cDatalayerService().setBreadCrumb("Home/Checkout/Step1")
  - var result = w3cDatalayerService().setPageInfos("pageID","pageName")
  - var result = w3cDatalayerService().setCartData(decoratedCart)
  - var result = w3cDatalayerService().setTransaction(cartTotals, decoratedItems, orderid)
  - var result = w3cDatalayerService().addProduct(product)
  - var result = w3cDatalayerService().addEvent("eventName")
      
```

If you want to populate the w3c datalayer to your page you can render the final digitalData object like this:
```
- var w3cDatalayerData = w3cDatalayerService().get()
script(type="text/javascript").
  var digitalData = !{w3cDatalayerData}
  digitalData.page.pageInfo.referringUrl = document.referrer
  digitalData.siteInfo.domain = document.location.hostname
```

