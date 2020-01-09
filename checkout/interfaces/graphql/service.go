package graphql

import (
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/graphql/dto"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -nometadata -o fs.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models return the 'Schema name' => 'Go model' mapping of this module
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_Checkout_PlaceOrderContext":                                        dto.PlaceOrderContext{},
		"Commerce_Checkout_PlacedOrderInfos":                                         dto.PlacedOrderInfos{},
		"Commerce_Checkout_PlaceOrderPaymentInfo":                                    application.PlaceOrderPaymentInfo{},
		"Commerce_Checkout_PlaceOrderState_State":                                    new(dto.State),
		"Commerce_Checkout_PlaceOrderState_State_Wait":                               dto.StateWait{},
		"Commerce_Checkout_PlaceOrderState_State_Success":                            dto.StateSuccess{},
		"Commerce_Checkout_PlaceOrderState_State_FatalError":                         dto.StateFatalError{},
		"Commerce_Checkout_PlaceOrderState_State_ShowIframe":                         dto.StateShowIframe{},
		"Commerce_Checkout_PlaceOrderState_State_ShowHtml":                           dto.StateShowHTML{},
		"Commerce_Checkout_PlaceOrderState_State_Redirect":                           dto.StateRedirect{},
		"Commerce_Checkout_PlaceOrderState_State_Cancelled":                          dto.StateCancelled{},
		"Commerce_Checkout_PlaceOrderState_State_CancellationReason":                 new(dto.CancellationReason),
		"Commerce_Checkout_PlaceOrderState_State_CancellationReason_PaymentError":    dto.CancellationReasonPaymentError{},
		"Commerce_Checkout_PlaceOrderState_State_CancellationReason_ValidationError": dto.CancellationReasonValidationError{},
	}.Models()
}
