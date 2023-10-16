package graphql

import (
	"context"
	// embed schema.graphql
	_ "embed"

	formDomain "flamingo.me/form/domain"
	"flamingo.me/graphql"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
)

// Service describes the Commerce/Cart GraphQL Service
type Service struct{}

var _ graphql.Service = new(Service)

//go:embed schema.graphql
var schema []byte

// Schema for cart, delivery and addresses
func (*Service) Schema() []byte {
	return schema
}

// Types configures the GraphQL to Go resolvers
func (*Service) Types(types *graphql.Types) {
	types.Map("Commerce_Cart_CustomAttributes", dto.CustomAttributes{})
	types.Map("Commerce_Cart_KeyValue", dto.KeyValue{})
	types.Map("Commerce_Cart_KeyValueInput", dto.KeyValue{})
	types.Map("Commerce_Cart_DeliveryAdditionalDataInput", dto.DeliveryAdditionalData{})
	types.Map("Commerce_Cart_DecoratedCart", dto.DecoratedCart{})
	types.Map("Commerce_Cart_Cart", cart.Cart{})
	types.Resolve("Commerce_Cart_Cart", "getDeliveryByCode", Resolver{}, "GetDeliveryByCodeWithoutBool")
	types.Map("Commerce_Cart_Summary", dto.CartSummary{})
	types.Map("Commerce_Cart_DecoratedDelivery", dto.DecoratedDelivery{})
	types.Map("Commerce_Cart_Delivery", cart.Delivery{})
	types.Map("Commerce_Cart_DeliveryInfo", cart.DeliveryInfo{})
	types.Resolve("Commerce_Cart_DeliveryInfo", "additionalData", CommerceCartDeliveryInfoResolver{}, "AdditionalData")
	types.Map("Commerce_Cart_DeliveryLocation", cart.DeliveryLocation{})
	types.Map("Commerce_Cart_Totalitem", cart.Totalitem{})
	types.Map("Commerce_Cart_Tax", cart.Tax{})
	types.Map("Commerce_Cart_Taxes", dto.Taxes{})
	types.Map("Commerce_Cart_Teaser", cart.Teaser{})
	types.Map("Commerce_Cart_CouponCode", cart.CouponCode{})
	types.Map("Commerce_Cart_AdditionalData", cart.AdditionalData{})
	types.Resolve("Commerce_Cart_AdditionalData", "customAttributes", CommerceCartAdditionalDataResolver{}, "CustomAttributes")
	types.Map("Commerce_Cart_ShippingItem", cart.ShippingItem{})
	types.Resolve("Commerce_Cart_ShippingItem", "appliedDiscounts", dto.CartAppliedDiscountsResolver{}, "ForShippingItem")
	types.Map("Commerce_Cart_DecoratedItem", dto.DecoratedCartItem{})
	types.Map("Commerce_Cart_Item", cart.Item{})
	types.Resolve("Commerce_Cart_Item", "appliedDiscounts", dto.CartAppliedDiscountsResolver{}, "ForItem")
	types.Map("Commerce_Cart_Address", cart.Address{})
	types.Map("Commerce_Cart_Person", cart.Person{})
	types.Map("Commerce_Cart_ExistingCustomerData", cart.ExistingCustomerData{})
	types.Map("Commerce_Cart_PersonalDetails", cart.PersonalDetails{})
	types.Map("Commerce_Cart_AppliedDiscounts", dto.CartAppliedDiscounts{})
	types.Map("Commerce_Cart_AppliedDiscount", cart.AppliedDiscount{})
	types.Map("Commerce_Cart_AppliedGiftCard", cart.AppliedGiftCard{})
	types.Map("Commerce_Cart_PricedItems", dto.PricedItems{})
	types.Map("Commerce_Cart_PricedCartItem", dto.PricedCartItem{})
	types.Map("Commerce_Cart_PricedShippingItem", dto.PricedShippingItem{})
	types.Map("Commerce_Cart_PricedTotalItem", dto.PricedTotalItem{})
	types.Map("Commerce_Cart_BillingAddressForm", dto.BillingAddressForm{})
	types.Map("Commerce_Cart_AddressForm", forms.AddressForm{})
	types.Map("Commerce_Cart_AddressFormInput", forms.AddressForm{})
	types.Map("Commerce_Cart_Form_ValidationInfo", dto.ValidationInfo{})
	types.Map("Commerce_Cart_Form_Error", formDomain.Error{})
	types.Map("Commerce_Cart_Form_FieldError", dto.FieldError{})
	types.Map("Commerce_Cart_ValidationResult", validation.Result{})
	types.Map("Commerce_Cart_ItemValidationError", validation.ItemValidationError{})
	types.Map("Commerce_Cart_PlacedOrderInfo", placeorder.PlacedOrderInfo{})
	types.Map("Commerce_Cart_SelectedPaymentResult", dto.SelectedPaymentResult{})
	types.Map("Commerce_Cart_PaymentSelection", new(cart.PaymentSelection))
	types.Map("Commerce_Cart_DefaultPaymentSelection", cart.DefaultPaymentSelection{})
	types.Resolve("Commerce_Cart_DefaultPaymentSelection", "cartSplit", CommerceCartQueryResolver{}, "CartSplit")
	types.Map("Commerce_Cart_DeliveryAddressForm", dto.DeliveryAddressForm{})
	types.Map("Commerce_Cart_DeliveryAddressInput", forms.DeliveryForm{})
	types.Map("Commerce_Cart_UpdateDeliveryShippingOptions_Result", dto.UpdateShippingOptionsResult{})
	types.Map("Commerce_Cart_AddToCartInput", dto.AddToCart{})
	types.Map("Commerce_Cart_ChoiceConfigurationInput", dto.ChoiceConfiguration{})

	types.GoField("Commerce_Cart_DeliveryAddressInput", "deliveryCode", "LocationCode")
	types.GoField("Commerce_Cart_DeliveryAddressInput", "carrier", "ShippingCarrier")
	types.GoField("Commerce_Cart_DeliveryAddressInput", "method", "ShippingMethod")
	types.Map("Commerce_Cart_DeliveryShippingOptionInput", dto.DeliveryShippingOption{})
	types.Map("Commerce_Cart_QtyRestrictionResult", validation.RestrictionResult{})
	types.Map("Commerce_Cart_PaymentSelection_Split", dto.PaymentSelectionSplit{})
	types.Map("Commerce_Cart_PaymentSelection_SplitQualifier", cart.SplitQualifier{})
	types.GoField("Commerce_Cart_PaymentSelection_SplitQualifier", "type", "ChargeType")
	types.GoField("Commerce_Cart_PaymentSelection_SplitQualifier", "reference", "ChargeReference")

	types.Resolve("Query", "Commerce_Cart_DecoratedCart", CommerceCartQueryResolver{}, "CommerceCart")
	types.Resolve("Query", "Commerce_Cart_Validator", CommerceCartQueryResolver{}, "CommerceCartValidator")
	types.Resolve("Query", "Commerce_Cart_QtyRestriction", CommerceCartQueryResolver{}, "CommerceCartQtyRestriction")

	types.Resolve("Mutation", "Commerce_Cart_AddToCart", CommerceCartMutationResolver{}, "CommerceAddToCart")
	types.Resolve("Mutation", "Commerce_Cart_DeleteCartDelivery", CommerceCartMutationResolver{}, "CommerceDeleteCartDelivery")
	types.Resolve("Mutation", "Commerce_Cart_DeleteCartDelivery", CommerceCartMutationResolver{}, "CommerceDeleteCartDelivery")
	types.Resolve("Mutation", "Commerce_Cart_DeleteItem", CommerceCartMutationResolver{}, "CommerceDeleteItem")
	types.Resolve("Mutation", "Commerce_Cart_UpdateItemQty", CommerceCartMutationResolver{}, "CommerceUpdateItemQty")
	types.Resolve("Mutation", "Commerce_Cart_UpdateItemBundleConfig", CommerceCartMutationResolver{}, "CommerceUpdateItemBundleConfig")
	types.Resolve("Mutation", "Commerce_Cart_UpdateBillingAddress", CommerceCartMutationResolver{}, "CommerceCartUpdateBillingAddress")
	types.Resolve("Mutation", "Commerce_Cart_UpdateSelectedPayment", CommerceCartMutationResolver{}, "CommerceCartUpdateSelectedPayment")
	types.Resolve("Mutation", "Commerce_Cart_ApplyCouponCodeOrGiftCard", CommerceCartMutationResolver{}, "CommerceCartApplyCouponCodeOrGiftCard")
	types.Resolve("Mutation", "Commerce_Cart_RemoveGiftCard", CommerceCartMutationResolver{}, "CommerceCartRemoveGiftCard")
	types.Resolve("Mutation", "Commerce_Cart_RemoveCouponCode", CommerceCartMutationResolver{}, "CommerceCartRemoveCouponCode")
	types.Resolve("Mutation", "Commerce_Cart_UpdateDeliveryAddresses", CommerceCartMutationResolver{}, "CommerceCartUpdateDeliveryAddresses")
	types.Resolve("Mutation", "Commerce_Cart_UpdateDeliveryShippingOptions", CommerceCartMutationResolver{}, "CommerceCartUpdateDeliveryShippingOptions")
	types.Resolve("Mutation", "Commerce_Cart_Clean", CommerceCartMutationResolver{}, "CartClean")
	types.Resolve("Mutation", "Commerce_Cart_UpdateAdditionalData", CommerceCartMutationResolver{}, "UpdateAdditionalData")
	types.Resolve("Mutation", "Commerce_Cart_UpdateDeliveriesAdditionalData", CommerceCartMutationResolver{}, "UpdateDeliveriesAdditionalData")
}

// Resolver helper
type Resolver struct{}

// GetDeliveryByCodeWithoutBool helper
func (*Resolver) GetDeliveryByCodeWithoutBool(_ context.Context, cart *cart.Cart, code string) (*cart.Delivery, error) {
	return cart.GetDeliveryByCodeWithoutBool(code), nil
}
