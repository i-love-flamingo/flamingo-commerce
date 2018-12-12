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

	ItemUpdateCommand struct {
		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceId       *string
		Qty            *int
		AdditionalData map[string]string
	}

	// CustomerCartService  interface
	CustomerCartService interface {
		GetCartOrderBehaviour(context.Context, Auth) (CartBehaviour, error)
		GetCart(ctx context.Context, auth Auth, cartId string) (*Cart, error)
	}
	PlacedOrderInfos []PlacedOrderInfo

	PlacedOrderInfo struct {
		OrderNumber  string
		DeliveryCode string
	}
	// CartBehaviour is a Port that can be implemented by other packages to implement  cart actions required for Ordering a Cart
	CartBehaviour interface {
		PlaceOrder(ctx context.Context, cart *Cart, payment *CartPayment) (PlacedOrderInfos, error)
		DeleteItem(ctx context.Context, cart *Cart, itemId string, deliveryCode string) (*Cart, error)
		UpdateItem(ctx context.Context, cart *Cart, itemId string, deliveryCode string, itemUpdateCommand ItemUpdateCommand) (*Cart, error)
		AddToCart(ctx context.Context, cart *Cart, deliveryCode string, addRequest AddRequest) (*Cart, error)
		CleanCart(ctx context.Context, cart *Cart) (*Cart, error)
		CleanDelivery(ctx context.Context, cart *Cart, deliveryCode string) (*Cart, error)
		UpdatePurchaser(ctx context.Context, cart *Cart, purchaser *Person, additionalData map[string]string) (*Cart, error)
		UpdateAdditionalData(ctx context.Context, cart *Cart, additionalData map[string]string) (*Cart, error)
		//UpdateDeliveryInfosAndBilling(ctx context.Context, cart *Cart, billingAddress *Address, deliveryInfoUpdates []DeliveryInfoUpdateCommand) (*Cart, error)
		UpdateDeliveryInfo(ctx context.Context, cart *Cart, deliveryCode string, deliveryInfo DeliveryInfo) (*Cart, error)
		UpdateBillingAddress(ctx context.Context, cart *Cart, billingAddress *Address) (*Cart, error)
		UpdateDeliveryInfoAdditionalData(ctx context.Context, cart *Cart, deliveryCode string, additionalData map[string]string) (*Cart, error)
		ApplyVoucher(ctx context.Context, cart *Cart, couponCode string) (*Cart, error)
	}

	AddRequest struct {
		MarketplaceCode        string
		Qty                    int
		VariantMarketplaceCode string
	}

	// Auth defines cart authentication information
	Auth struct {
		TokenSource oauth2.TokenSource
		IDToken     *oidc.IDToken
	}
)

var (
	CartNotFoundError    = errors.New("Cart not found")
	DeliveryCodeNotFound = errors.New("Delivery not found")
)

func (sv PlacedOrderInfos) GetOrderNumberForDeliverCode(deliveryCode string) string {
	for _, v := range sv {
		if v.DeliveryCode == deliveryCode {
			return v.OrderNumber
		}
	}
	return ""
}
