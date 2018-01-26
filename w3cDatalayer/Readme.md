# Datalayer Module

Simple modul that helps implementing tracking pixels and tagmanagers.
The datastructure is oriented at:
https://www.w3.org/community/custexpdata/

It:
* provides easy access to common informations relevant for tracking and datalayer
* listens to various events and adds informations to the datalayer


Usage:

The templatefunc provides you access to the current requests datalayer.
You can get the datalayer and you can modify it:

For some values in the datalayer the template knows better than the backend what to put in:
```
  - var result = w3cDatalayerService().setPageCategories("masterdata","brand","detail")
  - var result = w3cDatalayerService().setBreadCrumb("Home/Checkout/Step1")
```


If you want to populate the w3c datalayer to your page (digitalData object)
```
- var w3cDatalayerData = w3cDatalayerService().get()
script(type="text/javascript").
  var digitalData = !{w3cDatalayerData}
  digitalData.page.pageInfo.referringUrl = document.referrer
  digitalData.siteInfo.domain = document.location.hostname
```

