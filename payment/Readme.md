# Payment Package

Uff Payment is a tough topic and this package offers a generic concept to implement payment processing.

Before we start we should clarify some namings:
* PaymentGateway: Is a digital tool that allows for online (credit card) payment request processing. It accepts PaymentRequest
* PaymentMethod: Represents the used Payment -Typical payment methods include cash, checks, credit or debit cards, money orders, bank transfers and online payment services such as PayPal.

The main thing this package offers is a *WebCartPaymentGateway* interface.

## WebCartPaymentGateway
This is an interface that supports the processing of cart payments in any possible flow.

The checkout will use this to start a payment of a cart. 
Since a payment may involve redirects to one or more external hosted payment pages - or the requirement to show some iframe this interface uses a very generic abstraction.

Basically the checkout will at some point call
```go
 result, err := theSelectedWebCartPaymentGateway.StartFlow(ctx, cart, selectedMethod, correlationID, returnURL) 
 if err != nil {
   return result
 }
```
This basically forwards the control of what should happen in the browser of the user 100% to the Payment Flow.
The only thing it need to make sure is, that at the end the user should be returned to the given "returnUL". Which will be normally the payment processing page of the checkout.

There the checkout will ask the Gateway implementation again to get the current status of the payment flow and can decide if the payment if successful / failed or if further actions need to take place to proceed in the payment status (e.g. rendering an iframe of the payment provider)
```go
    flowStatus, err := gateway.FlowStatus(ctx, &decoratedCart.Cart, application.PaymentFlowStandardCorrelationID)
 	if err != nil {
 		return err
 	}
 
 	switch flowStatus.Status {
 	case paymentDomain.PaymentFlowStatusUnapproved:
 		// payment just started render payment page which handles actions
 		return cc.responder.Render("checkout/payment", viewData).SetNoCache()
 	case paymentDomain.PaymentFlowStatusApproved:
 		// payment is done but not confirmed by customer, confirm it and place order 
 		return cc.responder.RouteRedirect("checkout.placeorder", nil)
 	case paymentDomain.PaymentFlowStatusCompleted:
 		// payment is done and confirmed, place order
 		return cc.responder.RouteRedirect("checkout.placeorder", nil)
 	case paymentDomain.PaymentFlowStatusAborted:
 		// payment was aborted by user, redirect to checkout so a new payment can be started
 		return cc.responder.RouteRedirect("checkout", nil)
 	case paymentDomain.PaymentFlowStatusFailed, paymentDomain.PaymentFlowStatusCancelled:
 		// payment failed, redirect back to checkout
 		return cc.responder.RouteRedirect("checkout", nil)
 	case paymentDomain.PaymentFlowWaitingForCustomer:
 		// payment pending, waiting for customer
 		return cc.responder.Render("checkout/payment", viewData).SetNoCache()
 	default:
 		// show payment page which can react to unknown payment status
 		return cc.responder.Render("checkout/payment", viewData).SetNoCache()
 	}
 ..
```

## Offline Payment

This module also offers a simple implementation of an OfflineWebCartPaymentGateway - that can be used to process cart payments that are not done online but offline.
This implementation can be activate with this setting in `config.yml`:
```
commerce.payment.enableOfflinePaymentGateway: true
```


## Registering own Payment Providers

You need to implement the secondary port "WebCartPaymentGateway" and register your Gateway implementation in your `module.go` using Dingo.
Since in a project multiple implementations can be active, we use BindMap and use the **key (or code) of the gateway** as key in the map.
For example:

```go
injector.BindMap((*interfaces.WebCartPaymentGateway)(nil), "offlinepayment").To(interfaces.OfflineWebCartPaymentGateway{})
```
