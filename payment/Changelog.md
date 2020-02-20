# 30. September 2019
* Enhance WebCartPaymentGateway with CancelOrderPayment

# 20. August 2019
* Reworked WebCartPaymentGateway Interface
    * Removed StartWebFlow func
    * Renamed GetFlowResult to OrderPaymentFromFlow
    * Add function to get newly introduced FlowStatus
* Add new generic API endpoint to fetch the current FlowStatus of a payment
* Enhanced / updated domain model
    * Changed FlowResult to contain FlowStatus and allow flag to represent and early place order
    * Add FlowStatus which contains current status of the payment flow
    
# 20. February 2020
* Add `PaymentService` to easily work with bound PaymentGateway's
    * `PaymentService.AvailablePaymentGateways()` returns all bound gateways
    * `PaymentService.PaymentGateway()` gets the payment gateway by gateway code
    * `PaymentService.PaymentGatewayByCart()` gets the payment gateway of the cart payment selection

* Extend the `FlowStatus` struct with more standardized `FlowActionData`
* Add standardized Flow Actions `PaymentFlowActionShowIframe`, `PaymentFlowActionShowHTML`, `PaymentFlowActionRedirect`,
  `PaymentFlowActionPostRedirect` please use these in your payment adapter since the standard place order relies on them.