package application

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name Service --case snake --structname CartService

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/framework/web"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// Service that provides functionality regarding the cart
	Service interface {
		GetCartReceiverService() *CartReceiverService
		ValidateCart(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart) validation.Result
		ValidateCurrentCart(ctx context.Context, session *web.Session) (validation.Result, error)
		UpdatePaymentSelection(ctx context.Context, session *web.Session, paymentSelection cartDomain.PaymentSelection) error
		UpdateBillingAddress(ctx context.Context, session *web.Session, billingAddress *cartDomain.Address) error
		UpdateDeliveryInfo(ctx context.Context, session *web.Session, deliveryCode string, deliveryInfo cartDomain.DeliveryInfoUpdateCommand) error
		UpdatePurchaser(ctx context.Context, session *web.Session, purchaser *cartDomain.Person, additionalData *cartDomain.AdditionalData) error
		UpdateItemQty(ctx context.Context, session *web.Session, itemID string, deliveryCode string, qty int) error
		UpdateItemSourceID(ctx context.Context, session *web.Session, itemID string, sourceID string) error
		UpdateItems(ctx context.Context, session *web.Session, updateCommands []cartDomain.ItemUpdateCommand) error
		UpdateItemBundleConfig(ctx context.Context, session *web.Session, updateCommand cartDomain.ItemUpdateCommand) error
		DeleteItem(ctx context.Context, session *web.Session, itemID string, deliveryCode string) error
		DeleteAllItems(ctx context.Context, session *web.Session) error
		CompleteCurrentCart(ctx context.Context) (*cartDomain.Cart, error)
		RestoreCart(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.Cart, error)
		Clean(ctx context.Context, session *web.Session) error
		DeleteDelivery(ctx context.Context, session *web.Session, deliveryCode string) (*cartDomain.Cart, error)
		// Deprecated: build your own add request
		BuildAddRequest(ctx context.Context, marketplaceCode string, variantMarketplaceCode string, qty int, additionalData map[string]string) cartDomain.AddRequest
		AddProduct(ctx context.Context, session *web.Session, deliveryCode string, addRequest cartDomain.AddRequest) (productDomain.BasicProduct, error)
		CreateInitialDeliveryIfNotPresent(ctx context.Context, session *web.Session, deliveryCode string) (*cartDomain.Cart, error)
		GetInitialDelivery(deliveryCode string) (*cartDomain.DeliveryInfo, error)
		ApplyVoucher(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error)
		ApplyAny(ctx context.Context, session *web.Session, anyCode string) (*cartDomain.Cart, error)
		RemoveVoucher(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error)
		ApplyGiftCard(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error)
		RemoveGiftCard(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error)
		DeleteCartInCache(ctx context.Context, session *web.Session, cart *cartDomain.Cart)
		ReserveOrderIDAndSave(ctx context.Context, session *web.Session) (*cartDomain.Cart, error)
		ForceReserveOrderIDAndSave(ctx context.Context, session *web.Session) (*cartDomain.Cart, error)
		PlaceOrderWithCart(ctx context.Context, session *web.Session, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error)
		PlaceOrder(ctx context.Context, session *web.Session, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error)
		CancelOrder(ctx context.Context, session *web.Session, orderInfos placeorder.PlacedOrderInfos, cart cartDomain.Cart) (*cartDomain.Cart, error)
		CancelOrderWithoutRestore(ctx context.Context, session *web.Session, orderInfos placeorder.PlacedOrderInfos) error
		GetDefaultDeliveryCode() string
		AdjustItemsToRestrictedQty(ctx context.Context, session *web.Session) (QtyAdjustmentResults, error)
		UpdateAdditionalData(ctx context.Context, session *web.Session, additionalData map[string]string) (*cartDomain.Cart, error)
		UpdateDeliveryAdditionalData(ctx context.Context, session *web.Session, deliveryCode string, additionalData map[string]string) (*cartDomain.Cart, error)
	}
)

var (
	// ErrNoBundleConfigurationGiven returned when bundle configuration is not provided
	ErrNoBundleConfigurationGiven = errors.New("no bundle configuration given for configurable product")

	// ErrNoVariantForConfigurable returned when type configurable with active variant do not have variant selected
	ErrNoVariantForConfigurable = errors.New("no variant given for configurable product")

	// ErrVariantDoNotExist returned when selected variant do not exist
	ErrVariantDoNotExist = errors.New("product has not the given variant")
)
