package graphql

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/graphql/dto"
	formDomain "flamingo.me/form/domain"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

//go:generate go run github.com/go-bindata/go-bindata/go-bindata -nometadata -o schema.go -pkg graphql schema.graphql

// Service describes the Commerce/Cart GraphQL Service
type Service struct{}

// Schema for cart, delivery and addresses
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
}

// Models mapping for Commerce_Cart types
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_DecoratedCart": dto.DecoratedCart{},
		"Commerce_Cart": graphql.ModelMapEntry{
			Type: cart.Cart{},
			Fields: map[string]string{
				"getDeliveryByCode": "GetDeliveryByCodeWithoutBool",
			},
		},
		"Commerce_Cart_Summary":                 dto.CartSummary{},
		"Commerce_CartDecoratedDelivery":        decorator.DecoratedDelivery{},
		"Commerce_CartDelivery":                 cart.Delivery{},
		"Commerce_CartDeliveryInfo":             cart.DeliveryInfo{},
		"Commerce_CartDeliveryLocation":         cart.DeliveryLocation{},
		"Commerce_CartTotalitem":                cart.Totalitem{},
		"Commerce_Cart_Tax":                     cart.Tax{},
		"Commerce_Cart_Taxes":                   dto.Taxes{},
		"Commerce_Cart_Teaser":                  cart.Teaser{},
		"Commerce_CartCouponCode":               cart.CouponCode{},
		"Commerce_CartAdditionalData":           cart.AdditionalData{},
		"Commerce_CartShippingItem":             cart.ShippingItem{},
		"Commerce_CartDecoratedItem":            decorator.DecoratedCartItem{},
		"Commerce_CartItem":                     cart.Item{},
		"Commerce_CartAddress":                  cart.Address{},
		"Commerce_CartPerson":                   cart.Person{},
		"Commerce_CartExistingCustomerData":     cart.ExistingCustomerData{},
		"Commerce_CartPersonalDetails":          cart.PersonalDetails{},
		"Commerce_CartAppliedDiscounts":         cart.AppliedDiscounts{},
		"Commerce_CartAppliedDiscount":          cart.AppliedDiscount{},
		"Commerce_CartAppliedGiftCard":          cart.AppliedGiftCard{},
		"Commerce_Cart_PricedItems":             dto.PricedItems{},
		"Commerce_Cart_PricedCartItem":          dto.PricedCartItem{},
		"Commerce_Cart_PricedShippingItem":      dto.PricedShippingItem{},
		"Commerce_Cart_PricedTotalItem":         dto.PricedTotalItem{},
		"Commerce_Cart_BillingAddressForm":      dto.BillingAddressForm{},
		"Commerce_Cart_AddressForm":             forms.AddressForm{},
		"Commerce_Cart_AddressFormInput":        forms.AddressForm{},
		"Commerce_Cart_Form_ValidationInfo":     dto.ValidationInfo{},
		"Commerce_Cart_Form_Error":              formDomain.Error{},
		"Commerce_Cart_Form_FieldError":         dto.FieldError{},
		"Commerce_Cart_ValidationResult":        validation.Result{},
		"Commerce_Cart_ItemValidationError":     validation.ItemValidationError{},
		"Commerce_Cart_PlacedOrderInfo":         placeorder.PlacedOrderInfo{},
		"Commerce_Cart_SelectedPaymentResult":   dto.SelectedPaymentResult{},
		"Commerce_Cart_PaymentSelection":        new(cart.PaymentSelection),
		"Commerce_Cart_DefaultPaymentSelection": cart.DefaultPaymentSelection{},
		"Commerce_Cart_DeliveryAddressForm":     dto.DeliveryAddressForm{},
		"Commerce_Cart_DeliveryAddressInput": graphql.ModelMapEntry{
			Type: forms.DeliveryForm{},
			Fields: map[string]string{
				"deliveryCode": "LocationCode",
				"carrier":      "ShippingCarrier",
				"method":       "ShippingMethod",
			},
		},
		"Commerce_Cart_DeliveryShippingOption": dto.DeliveryShippingOption{},
		"Commerce_Cart_QtyRestrictionResult":   validation.RestrictionResult{},
		"Commerce_Cart_PaymentSelection_Split": dto.PaymentSelectionSplit{},
		"Commerce_Cart_PaymentSelection_SplitQualifier": graphql.ModelMapEntry{
			Type: cart.SplitQualifier{},
			Fields: map[string]string{
				"type":      "ChargeType",
				"reference": "ChargeReference",
			},
		},
	}.Models()
}
