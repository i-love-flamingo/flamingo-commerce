package graphql

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/graphql"
	"github.com/99designs/gqlgen/codegen/config"
)

// Service describes the Commerce/Cart GraphQL Service
type Service struct{}

// Schema for cart, delivery and addresses
func (*Service) Schema() []byte {
	// language=graphql
	return []byte(`
type Commerce_DecoratedCart {
	cart: Commerce_Cart!
	decorated_deliveries: [Commerce_CartDecoratedDelivery!]
}

type Commerce_Cart {
	id: ID!
	entityID: String!
	billingAdress: Commerce_CartAddress
	purchaser: Commerce_CartPerson
	deliveries: [Commerce_CartDelivery!]
#	additionalData: Commerce_CartAdditionalData!
#	paymentSelection: Commerce_CartPaymentSelection!
	belongsToAuthenticatedUser: Boolean!
	authenticatedUserID: String!
#	appliedCouponCodes: [Commerce_CartCouponCode!]
	defaultCurrency: String!
#	totalitems: [Commerce_CartTotalitem!]
#	appliedGiftCards: [Commerce_CartAppliedGiftCard!]
}

type Commerce_CartDecoratedDelivery {
	delivery: Commerce_CartDelivery!
	decoratedItems: [Commerce_CartDecoratedItem!]
}

type Commerce_CartDelivery {
	cartitems: [Commerce_CartItem!]
}

type Commerce_CartDecoratedItem {
	item: Commerce_CartItem
	product: Commerce_Product
}

type Commerce_CartItem {
	id: ID!
	qty: Int
}

type Commerce_CartAddress {
	vat:                    String!
	firstname:              String!
	lastname:               String!
	middleName:             String!
	title:                  String!
	salutation:             String!
	street:                 String!
	streetNr:               String!
	additionalAddressLines: [String!]
	company:                String!
	city:                   String!
	postCode:               String!
	state:                  String!
	regionCode:             String!
	country:                String!
	countryCode:            String!
	telephone:              String!
	email:                  String!
}

type Commerce_CartPerson {
	address: Commerce_CartAddress
	personalDetails: Commerce_CartPersonalDetails!
	existingCustomerData: Commerce_CartExistingCustomerData
}

type Commerce_CartExistingCustomerData {
	id: ID!
}

type Commerce_CartPersonalDetails {
	dateOfBirth: String!
	passportCountry: String!
	passportNumber: String!
	nationality: String!
}

extend type Query {
	commerce_Cart: Commerce_DecoratedCart!
}

extend type Mutation {
	commerce_AddToCart(id: ID!, qty: Int, deliveryCode: String!): Commerce_DecoratedCart!
}
`)
}

// Models mapping for Commerce_Cart types
func (*Service) Models() map[string]config.TypeMapEntry {
	return graphql.ModelMap{
		"Commerce_DecoratedCart":            decorator.DecoratedCart{},
		"Commerce_Cart":                     cart.Cart{},
		"Commerce_CartDecoratedDelivery":    decorator.DecoratedDelivery{},
		"Commerce_CartDelivery":             cart.Delivery{},
		"Commerce_CartDecoratedItem":        decorator.DecoratedCartItem{},
		"Commerce_CartItem":                 cart.Item{},
		"Commerce_CartAddress":              cart.Address{},
		"Commerce_CartPerson":               cart.Person{},
		"Commerce_CartExistingCustomerData": cart.ExistingCustomerData{},
		"Commerce_CartPersonalDetails":      cart.PersonalDetails{},
	}.Models()
}

// CommerceCartQueryResolver resolver for carts
type CommerceCartQueryResolver struct {
	applicationCartReceiverService *application.CartReceiverService
}

// Inject dependencies
func (r *CommerceCartQueryResolver) Inject(applicationCartReceiverService *application.CartReceiverService) {
	r.applicationCartReceiverService = applicationCartReceiverService
}

// CommerceCart getter for queries
func (r *CommerceCartQueryResolver) CommerceCart(ctx context.Context) (*decorator.DecoratedCart, error) {
	req := web.RequestFromContext(ctx)

	return r.applicationCartReceiverService.ViewDecoratedCart(ctx, req.Session())
}

// CommerceCartMutationResolver resolves cart mutations
type CommerceCartMutationResolver struct {
	q                      *CommerceCartQueryResolver
	applicationCartService *application.CartService
}

// Inject dependencies
func (r *CommerceCartMutationResolver) Inject(q *CommerceCartQueryResolver, applicationCartService *application.CartService) {
	r.q = q
	r.applicationCartService = applicationCartService
}

// CommerceAddToCart mutation for adding products to the current users cart
func (r *CommerceCartMutationResolver) CommerceAddToCart(ctx context.Context, id string, qty *int, deliveryCode string) (*decorator.DecoratedCart, error) {
	if qty == nil {
		one := 1
		qty = &one
	}

	req := web.RequestFromContext(ctx)

	addRequest := r.applicationCartService.BuildAddRequest(ctx, id, "", *qty)

	_, err := r.applicationCartService.AddProduct(ctx, req.Session(), deliveryCode, addRequest)
	if err != nil {
		return nil, err
	}

	return r.q.CommerceCart(ctx)
}
