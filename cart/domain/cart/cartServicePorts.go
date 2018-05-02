package cart

import (
	"context"

	oidc "github.com/coreos/go-oidc"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type (
	// GuestCartService interface
	GuestCartService interface {
		GetCart(ctx context.Context, cartId string) (*Cart, error)

		//GetGuestCart - should return a new guest cart (including the id of the cart)
		GetNewCart(ctx context.Context) (*Cart, error)

		GetCartOrderBehaviour(context.Context) (CartBehaviour, error)
	}

	//DeliveryInfoUpdateCommand that is consumed by the CartBehaviour and is used to update the Cart with DeliveryInfos
	DeliveryInfoUpdateCommand struct {
		//DeliveryInfo - the deliveryinfo for update
		DeliveryInfo *DeliveryInfo
		//DeliveryInfoID - if set the ID that should be updated (if not given a new DeliveryInfo should be added to the cart)
		DeliveryInfoID  string
		AssignedItemIds []string
	}

	ItemUpdateCommand struct {
		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceId               *string
		Qty                    *int
		OriginalDeliveryIntent *DeliveryIntent
		AdditionalData         map[string]string
	}

	// CustomerCartService  interface
	CustomerCartService interface {
		GetCartOrderBehaviour(context.Context, Auth) (CartBehaviour, error)
		GetCart(ctx context.Context, auth Auth, cartId string) (*Cart, error)
	}

	// CartBehaviour is a Port that can be implemented by other packages to implement  cart actions required for Ordering a Cart
	CartBehaviour interface {
		PlaceOrder(ctx context.Context, cart *Cart, payment *CartPayment) (string, error)
		DeleteItem(ctx context.Context, cart *Cart, itemId string) (*Cart, error)
		UpdateItem(ctx context.Context, cart *Cart, itemId string, itemUpdateCommand ItemUpdateCommand) (*Cart, error)
		AddToCart(ctx context.Context, cart *Cart, addRequest AddRequest) (*Cart, error)
		UpdatePurchaser(ctx context.Context, cart *Cart, purchaser *Person, additionalData map[string]string) (*Cart, error)
		UpdateAdditionalData(ctx context.Context, cart *Cart, additionalData map[string]string) (*Cart, error)
		UpdateDeliveryInfosAndBilling(ctx context.Context, cart *Cart, billingAddress *Address, deliveryInfoUpdates []DeliveryInfoUpdateCommand) (*Cart, error)
		ApplyVoucher(ctx context.Context, cart *Cart, couponCode string) (*Cart, error)
	}

	AddRequest struct {
		MarketplaceCode        string
		Qty                    int
		VariantMarketplaceCode string
		DeliveryIntent         DeliveryIntent
	}

	// Auth defines cart authentication information
	Auth struct {
		TokenSource oauth2.TokenSource
		IDToken     *oidc.IDToken
	}
)

var (
	CartNotFoundError = errors.New("Cart not found")
)
