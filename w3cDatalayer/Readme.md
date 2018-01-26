# Datalayer Module

Simple modul that helps implementing tracking pixels and tagmanagers.
The datastructure is oriented at:
https://www.w3.org/community/custexpdata/

It:
* provides easy access to common informations relevant for tracking and datalayer
* listens to various events and adds informations to the datalayer


Usage:

If you want to populate the w3c datalayer (digitalData object)
```
script(type="text/javascript")
  - var w3cDatalayerData = data("w3cDatalayer")
  | var digitalData = !{w3cDatalayerData}

```

