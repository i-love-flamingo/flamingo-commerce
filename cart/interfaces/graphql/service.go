package graphql

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
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
		"Commerce_DecoratedCart": graphql.ModelMapEntry{
			Type: decorator.DecoratedCart{},
			Fields: map[string]string{
				"getDecoratedDeliveryByCode": "GetDecoratedDeliveryByCodeWithoutBool",
			},
		},
		"Commerce_Cart": graphql.ModelMapEntry{
			Type: cart.Cart{},
			Fields: map[string]string{
				"getDeliveryByCode": "GetDeliveryByCodeWithoutBool",
			},
		},
		"Commerce_CartDecoratedDelivery":       decorator.DecoratedDelivery{},
		"Commerce_CartDelivery":                cart.Delivery{},
		"Commerce_CartDeliveryInfo":            cart.DeliveryInfo{},
		"Commerce_CartDeliveryLocation":        cart.DeliveryLocation{},
		"Commerce_CartTotalitem":               cart.Totalitem{},
		"Commerce_CartCouponCode":              cart.CouponCode{},
		"Commerce_CartAdditionalData":          cart.AdditionalData{},
		"Commerce_CartShippingItem":            cart.ShippingItem{},
		"Commerce_CartDecoratedItem":           decorator.DecoratedCartItem{},
		"Commerce_CartItem":                    cart.Item{},
		"Commerce_CartAddress":                 cart.Address{},
		"Commerce_CartPerson":                  cart.Person{},
		"Commerce_CartExistingCustomerData":    cart.ExistingCustomerData{},
		"Commerce_CartPersonalDetails":         cart.PersonalDetails{},
		"Commerce_CartAppliedDiscounts":        cart.AppliedDiscounts{},
		"Commerce_CartAppliedDiscount":         cart.AppliedDiscount{},
		"Commerce_CartAppliedGiftCard":         cart.AppliedGiftCard{},
		"Commerce_Cart_BillingAddressForm":     dto.BillingAddressForm{},
		"Commerce_Cart_BillingAddressFormData": forms.BillingAddressForm{},
		"Commerce_BillingAddressFormInput":     forms.BillingAddressForm{},
		"Commerce_Cart_Form_ValidationInfo":    dto.ValidationInfo{},
		"Commerce_Cart_Form_Error":             formDomain.Error{},
		"Commerce_Cart_Form_FieldError":        dto.FieldError{},
	}.Models()
}
