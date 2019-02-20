package cart

import (
	"context"

	"flamingo.me/flamingo/v3/core/auth/domain"
	"github.com/pkg/errors"
)

type (
	// GuestCartService interface - Secondary PORT
	GuestCartService interface {
		// GetBehaviour gets the behaviour for the guest cart service
		GetModifyBehaviour(context.Context) (ModifyBehaviour, error)
		GetCart(ctx context.Context, cartID string) (*Cart, error)

		//GetGuestCart - should return a new guest cart (including the id of the cart)
		GetNewCart(ctx context.Context) (*Cart, error)
	}

	// CustomerCartService  interface - Secondary PORT
	CustomerCartService interface {
		// GetBehaviour gets the behaviour for the customer cart service
		GetModifyBehaviour(context.Context, domain.Auth) (ModifyBehaviour, error)
		GetCart(ctx context.Context, auth domain.Auth, cartID string) (*Cart, error)
	}


	// ModifyBehaviour is a interface that can be implemented by other packages to provide cart actions
	// This port can not be registered directly but is provided by the registered "GuestCartService"
	ModifyBehaviour interface {
		DeleteItem(ctx context.Context, cart *Cart, itemID string, deliveryCode string) (*Cart, error)
		UpdateItem(ctx context.Context, cart *Cart, itemID string, deliveryCode string, itemUpdateCommand ItemUpdateCommand) (*Cart, error)
		AddToCart(ctx context.Context, cart *Cart, deliveryCode string, addRequest AddRequest) (*Cart, error)
		CleanCart(ctx context.Context, cart *Cart) (*Cart, error)
		CleanDelivery(ctx context.Context, cart *Cart, deliveryCode string) (*Cart, error)
		UpdatePurchaser(ctx context.Context, cart *Cart, purchaser *Person, additionalData *AdditionalData) (*Cart, error)
		UpdateAdditionalData(ctx context.Context, cart *Cart, additionalData *AdditionalData) (*Cart, error)
		//UpdateDeliveryInfosAndBilling(ctx context.Context, cart *Cart, billingAddress *Address, deliveryInfoUpdates []DeliveryInfoUpdateCommand) (*Cart, error)
		UpdateDeliveryInfo(ctx context.Context, cart *Cart, deliveryCode string, deliveryInfo DeliveryInfo) (*Cart, error)
		UpdateBillingAddress(ctx context.Context, cart *Cart, billingAddress *Address) (*Cart, error)
		UpdateDeliveryInfoAdditionalData(ctx context.Context, cart *Cart, deliveryCode string, additionalData *AdditionalData) (*Cart, error)
		ApplyVoucher(ctx context.Context, cart *Cart, couponCode string) (*Cart, error)
	}

	// AddRequest defines add to cart requeset
	AddRequest struct {
		MarketplaceCode        string
		Qty                    int
		VariantMarketplaceCode string
	}

	// ItemUpdateCommand defines the update item command
	ItemUpdateCommand struct {
		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceID       *string
		Qty            *int
		AdditionalData map[string]string
	}


	// PlaceOrderService  interface - Secondary PORT
	PlaceOrderService interface {
		PlaceGuestCart(ctx context.Context, cart *Cart, payment *Payment) (PlacedOrderInfos, error)
		PlaceCustomerCart(ctx context.Context, auth domain.Auth, cart *Cart, payment *Payment) (PlacedOrderInfos, error)
	}
)

var (
	// ErrCartNotFound is used if a cart was not found
	ErrCartNotFound = errors.New("Cart not found")
	// ErrDeliveryCodeNotFound is used if a delivery was not found
	ErrDeliveryCodeNotFound = errors.New("Delivery not found")
)
