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
    
