package graphql

import (
	// embed schema.graphql
	_ "embed"

	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	"flamingo.me/graphql"
)

// Service is the Graphql-Service of this module
type Service struct{}

var _ graphql.Service = new(Service)

//go:embed schema.graphql
var schema []byte

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return schema
}

// Types configures the GraphQL to Go resolvers
func (*Service) Types(types *graphql.Types) {
	types.Map("Commerce_Checkout_PlaceOrderContext", dto.PlaceOrderContext{})
	types.Map("Commerce_Checkout_StartPlaceOrder_Result", dto.StartPlaceOrderResult{})
	types.Map("Commerce_Checkout_PlacedOrderInfos", dto.PlacedOrderInfos{})
	types.Map("Commerce_Checkout_PlaceOrderPaymentInfo", application.PlaceOrderPaymentInfo{})
	types.Map("Commerce_Checkout_PlaceOrderState_State", new(dto.State))
	types.Map("Commerce_Checkout_PlaceOrderState_State_Wait", dto.Wait{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_WaitForCustomer", dto.WaitForCustomer{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_Success", dto.Success{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_Failed", dto.Failed{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_ShowIframe", dto.ShowIframe{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_ShowHTML", dto.ShowHTML{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_ShowWalletPayment", dto.ShowWalletPayment{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_TriggerClientSDK", dto.TriggerClientSDK{})
	types.Map("Commerce_Checkout_PlaceOrderState_PaymentRequestAPI", dto.PaymentRequestAPI{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_Redirect", dto.Redirect{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_PostRedirect", dto.PostRedirect{})
	types.Map("Commerce_Checkout_PlaceOrderState_Form_Parameter", dto.FormParameter{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_FailedReason", new(process.FailedReason))
	types.Map("Commerce_Checkout_PlaceOrderState_State_FailedReason_Error", process.ErrorOccurredReason{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentError", process.PaymentErrorOccurredReason{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_FailedReason_CartValidationError", process.CartValidationErrorReason{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_FailedReason_CanceledByCustomer", process.CanceledByCustomerReason{})
	types.Map("Commerce_Checkout_PlaceOrderState_State_FailedReason_PaymentCanceledByCustomer", process.PaymentCanceledByCustomerReason{})

	types.Resolve("Query", "Commerce_Checkout_ActivePlaceOrder", CommerceCheckoutQueryResolver{}, "CommerceCheckoutActivePlaceOrder")
	types.Resolve("Query", "Commerce_Checkout_CurrentContext", CommerceCheckoutQueryResolver{}, "CommerceCheckoutCurrentContext")
	types.Resolve("Mutation", "Commerce_Checkout_StartPlaceOrder", CommerceCheckoutMutationResolver{}, "CommerceCheckoutStartPlaceOrder")
	types.Resolve("Mutation", "Commerce_Checkout_CancelPlaceOrder", CommerceCheckoutMutationResolver{}, "CommerceCheckoutCancelPlaceOrder")
	types.Resolve("Mutation", "Commerce_Checkout_ClearPlaceOrder", CommerceCheckoutMutationResolver{}, "CommerceCheckoutClearPlaceOrder")
	types.Resolve("Mutation", "Commerce_Checkout_RefreshPlaceOrder", CommerceCheckoutMutationResolver{}, "CommerceCheckoutRefreshPlaceOrder")
	types.Resolve("Mutation", "Commerce_Checkout_RefreshPlaceOrderBlocking", CommerceCheckoutMutationResolver{}, "CommerceCheckoutRefreshPlaceOrderBlocking")
}
