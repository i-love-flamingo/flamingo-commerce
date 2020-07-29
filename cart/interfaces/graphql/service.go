package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	formDomain "flamingo.me/form/domain"
	"flamingo.me/graphql"
)

//go:generate go run github.com/go-bindata/go-bindata/v3/go-bindata -nometadata -o schema.go -pkg graphql schema.graphql

// Service describes the Commerce/Cart GraphQL Service
type Service struct{}

var _ graphql.Service = new(Service)

// Schema for cart, delivery and addresses
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Types configures the GraphQL to Go resolvers
func (*Service) Types(types *graphql.Types) {
	types.Map("Commerce_DecoratedCart", dto.DecoratedCart{})
	types.Map("Commerce_Cart", cart.Cart{})
	types.Resolve("Commerce_Cart", "getDeliveryByCode", Resolver{}, "GetDeliveryByCodeWithoutBool")
	types.Map("Commerce_Cart_Summary", dto.CartSummary{})
	types.Map("Commerce_CartDecoratedDelivery", decorator.DecoratedDelivery{})
	types.Map("Commerce_CartDelivery", cart.Delivery{})
	types.Map("Commerce_CartDeliveryInfo", cart.DeliveryInfo{})
	types.Map("Commerce_CartDeliveryLocation", cart.DeliveryLocation{})
	types.Map("Commerce_CartTotalitem", cart.Totalitem{})
	types.Map("Commerce_Cart_Tax", cart.Tax{})
	types.Map("Commerce_Cart_Taxes", dto.Taxes{})
	types.Map("Commerce_Cart_Teaser", cart.Teaser{})
	types.Map("Commerce_CartCouponCode", cart.CouponCode{})
	types.Map("Commerce_CartAdditionalData", cart.AdditionalData{})
	types.Map("Commerce_CartShippingItem", cart.ShippingItem{})
	types.Resolve("Commerce_CartShippingItem", "appliedDiscounts", dto.CartAppliedDiscountsResolver{}, "ForShippingItem")
	types.Map("Commerce_CartDecoratedItem", decorator.DecoratedCartItem{})
	types.Map("Commerce_CartItem", cart.Item{})
	types.Resolve("Commerce_CartItem", "appliedDiscounts", dto.CartAppliedDiscountsResolver{}, "ForItem")
	types.Map("Commerce_CartAddress", cart.Address{})
	types.Map("Commerce_CartPerson", cart.Person{})
	types.Map("Commerce_CartExistingCustomerData", cart.ExistingCustomerData{})
	types.Map("Commerce_CartPersonalDetails", cart.PersonalDetails{})
	types.Map("Commerce_CartAppliedDiscounts", dto.CartAppliedDiscounts{})
	types.Map("Commerce_CartAppliedDiscount", cart.AppliedDiscount{})
	types.Map("Commerce_CartAppliedGiftCard", cart.AppliedGiftCard{})
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
	types.GoField("Commerce_Cart_DeliveryAddressInput", "deliveryCode", "LocationCode")
	types.GoField("Commerce_Cart_DeliveryAddressInput", "carrier", "ShippingCarrier")
	types.GoField("Commerce_Cart_DeliveryAddressInput", "method", "ShippingMethod")
	types.Map("Commerce_Cart_DeliveryShippingOption", dto.DeliveryShippingOption{})
	types.Map("Commerce_Cart_QtyRestrictionResult", validation.RestrictionResult{})
	types.Map("Commerce_Cart_PaymentSelection_Split", dto.PaymentSelectionSplit{})
	types.Map("Commerce_Cart_PaymentSelection_SplitQualifier", cart.SplitQualifier{})
	types.GoField("Commerce_Cart_PaymentSelection_SplitQualifier", "type", "ChargeType")
	types.GoField("Commerce_Cart_PaymentSelection_SplitQualifier", "reference", "ChargeReference")

	types.Resolve("Query", "Commerce_Cart", CommerceCartQueryResolver{}, "CommerceCart")
	types.Resolve("Query", "Commerce_Cart_Validator", CommerceCartQueryResolver{}, "CommerceCartValidator")
	types.Resolve("Query", "Commerce_Cart_QtyRestriction", CommerceCartQueryResolver{}, "CommerceCartQtyRestriction")

	types.Resolve("Mutation", "Commerce_AddToCart", CommerceCartMutationResolver{}, "CommerceAddToCart")
	types.Resolve("Mutation", "Commerce_DeleteCartDelivery", CommerceCartMutationResolver{}, "CommerceDeleteCartDelivery")
	types.Resolve("Mutation", "Commerce_DeleteCartDelivery", CommerceCartMutationResolver{}, "CommerceDeleteCartDelivery")
	types.Resolve("Mutation", "Commerce_DeleteItem", CommerceCartMutationResolver{}, "CommerceDeleteItem")
	types.Resolve("Mutation", "Commerce_UpdateItemQty", CommerceCartMutationResolver{}, "CommerceUpdateItemQty")
	types.Resolve("Mutation", "Commerce_Cart_UpdateBillingAddress", CommerceCartMutationResolver{}, "CommerceCartUpdateBillingAddress")
	types.Resolve("Mutation", "Commerce_Cart_UpdateSelectedPayment", CommerceCartMutationResolver{}, "CommerceCartUpdateSelectedPayment")
	types.Resolve("Mutation", "Commerce_Cart_ApplyCouponCodeOrGiftCard", CommerceCartMutationResolver{}, "CommerceCartApplyCouponCodeOrGiftCard")
	types.Resolve("Mutation", "Commerce_Cart_RemoveGiftCard", CommerceCartMutationResolver{}, "CommerceCartRemoveGiftCard")
	types.Resolve("Mutation", "Commerce_Cart_RemoveCouponCode", CommerceCartMutationResolver{}, "CommerceCartRemoveCouponCode")
	types.Resolve("Mutation", "Commerce_Cart_UpdateDeliveryAddresses", CommerceCartMutationResolver{}, "CommerceCartUpdateDeliveryAddresses")
	types.Resolve("Mutation", "Commerce_Cart_UpdateDeliveryShippingOptions", CommerceCartMutationResolver{}, "CommerceCartUpdateDeliveryShippingOptions")
}

// Resolver helper
type Resolver struct{}

// GetDeliveryByCodeWithoutBool helper
func (*Resolver) GetDeliveryByCodeWithoutBool(_ context.Context, cart *cart.Cart, code string) (*cart.Delivery, error) {
	return cart.GetDeliveryByCodeWithoutBool(code), nil
}
