package application

import (
	"context"
	"fmt"

	cartDomain "flamingo.me/flamingo-commerce/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/product/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

// CartService application struct
type (
	//CartService provides methods to modify the cart
	CartService struct {
		CartReceiverService *CartReceiverService         `inject:""`
		ProductService      productDomain.ProductService `inject:""`
		Logger              flamingo.Logger              `inject:""`
		CartValidator       cartDomain.CartValidator     `inject:",optional"`

		ItemValidator  cartDomain.ItemValidator `inject:",optional"`
		EventPublisher EventPublisher           `inject:""`

		PickUpDetectionService cartDomain.PickUpDetectionService `inject:",optional"`

		DeliveryIntentBuilder *cartDomain.DeliveryIntentBuilder `inject:""`

		DefaultDeliveryIntent string `inject:"config:cart.defaultDeliveryIntent,optional"`

		CartCache CartCache `inject:",optional"`
	}
)

// ValidateCart validates a carts content
func (cs CartService) ValidateCart(ctx context.Context, session *sessions.Session, decoratedCart *cartDomain.DecoratedCart) cartDomain.CartValidationResult {

	if cs.CartValidator != nil {
		result := cs.CartValidator.Validate(ctx, session, decoratedCart)
		return result
	}

	return cartDomain.CartValidationResult{}
}

// ValidateCurrentCart validates the current active cart
func (cs CartService) ValidateCurrentCart(ctx context.Context, session *sessions.Session) (cartDomain.CartValidationResult, error) {
	decoratedCart, err := cs.CartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		return cartDomain.CartValidationResult{}, err
	}

	return cs.ValidateCart(ctx, session, decoratedCart), nil
}

func (cs CartService) UpdateDeliveryInfosAndBilling(ctx context.Context, session *sessions.Session, billingAddress *cartDomain.Address, updateCommands []cartDomain.DeliveryInfoUpdateCommand) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	cart, err = behaviour.UpdateDeliveryInfosAndBilling(ctx, cart, billingAddress, updateCommands)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, session, cart)
	return nil
}
func (cs CartService) UpdatePurchaser(ctx context.Context, session *sessions.Session, purchaser *cartDomain.Person, additionalData map[string]string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	cart, err = behaviour.UpdatePurchaser(ctx, cart, purchaser, additionalData)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, session, cart)
	return nil
}

// UpdateItemQty
func (cs CartService) UpdateItemQty(ctx context.Context, session *sessions.Session, itemId string, qty int) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	item, err := cart.GetByItemId(itemId)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)
		return err
	}
	qtyBefore := item.Qty
	if qty < 1 {
		return cs.DeleteItem(ctx, session, itemId)
	}

	cs.EventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, qty, cart.ID)
	itemUpdate := cartDomain.ItemUpdateCommand{
		Qty: &qty,
	}
	cart, err = behaviour.UpdateItem(ctx, cart, itemId, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, session, cart)
	return nil
}

func (cs CartService) UpdateItemSourceId(ctx context.Context, session *sessions.Session, itemId string, sourceId string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	_, err = cart.GetByItemId(itemId)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemSourceId").Error(err)
		return err
	}
	itemUpdate := cartDomain.ItemUpdateCommand{
		SourceId: &sourceId,
	}
	cart, err = behaviour.UpdateItem(ctx, cart, itemId, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemSourceId").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, session, cart)
	return nil
}

// DeleteItem in current cart
func (cs CartService) DeleteItem(ctx context.Context, session *sessions.Session, itemId string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	item, err := cart.GetByItemId(itemId)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteItem").Error(err)
		return err
	}
	qtyBefore := item.Qty
	cs.EventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, 0, cart.ID)
	cart, err = behaviour.DeleteItem(ctx, cart, itemId)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteItem").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, session, cart)
	return nil
}

// DeleteAllItems in current cart
func (cs CartService) DeleteAllItems(ctx context.Context, session *sessions.Session) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}

	for _, itemId := range cart.GetItemIds() {
		item, err := cart.GetByItemId(itemId)
		if err != nil {
			cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteAllItems").Error(err)
			return err
		}

		qtyBefore := item.Qty
		cs.EventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, 0, cart.ID)
		cart, err = behaviour.DeleteItem(ctx, cart, itemId)
		if err != nil {
			cs.handleCartNotFound(session, err)
			cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteAllItems").Error(err)
			return err
		}
	}

	cs.updateCartInCache(ctx, session, cart)
	return nil
}

// PlaceOrder
func (cs *CartService) PlaceOrder(ctx context.Context, session *sessions.Session, payment *cartDomain.CartPayment) ([]string, error) {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	orderNumbers, err := behaviour.PlaceOrder(ctx, cart, payment)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "PlaceOrder").Error(err)
		return nil, err
	}
	cs.EventPublisher.PublishOrderPlacedEvent(ctx, cart, orderNumbers)
	cs.DeleteSavedSessionGuestCartId(session)
	cs.deleteCartInCache(ctx, session, cart)
	return orderNumbers, err
}

// BuildAddRequest Helper to build
func (cs *CartService) BuildAddRequest(ctx context.Context, marketplaceCode string, variantMarketplaceCode string, qty int, deliveryIntentStringRepresentation string) cartDomain.AddRequest {
	if qty < 0 {
		qty = 0
	}
	if deliveryIntentStringRepresentation == "" {
		deliveryIntentStringRepresentation = cs.DefaultDeliveryIntent
	}
	return cartDomain.AddRequest{
		MarketplaceCode: marketplaceCode,
		Qty:             qty,
		VariantMarketplaceCode: variantMarketplaceCode,
		DeliveryIntent:         cs.DeliveryIntentBuilder.BuildDeliveryIntent(deliveryIntentStringRepresentation),
	}
}

// AddProduct Add a product
func (cs *CartService) AddProduct(ctx context.Context, session *sessions.Session, addRequest cartDomain.AddRequest) (error, productDomain.BasicProduct) {
	addRequest, product, err := cs.checkProductForAddRequest(ctx, session, addRequest)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)
		return err, nil
	}

	cs.Logger.WithField(flamingo.LogKeyCategory, "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Debug("AddRequest received %#v  / %v", addRequest, addRequest.DeliveryIntent.String())

	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)
		return err, nil
	}

	//Check if we can autodetect empty location code for pickup
	if addRequest.DeliveryIntent.Method == cartDomain.DELIVERY_METHOD_PICKUP && addRequest.DeliveryIntent.DeliveryLocationCode == "" && addRequest.DeliveryIntent.AutodetectDeliveryLocation {
		if cs.PickUpDetectionService != nil {
			locationCode, locationType, err := cs.PickUpDetectionService.Detect(product, addRequest)
			if err == nil {
				cs.Logger.WithField(flamingo.LogKeyCategory, "cartService").WithField("subCategory", "AddProduct").Debug("Detected pickup location %v / %v", locationCode, locationType)
				addRequest.DeliveryIntent.DeliveryLocationCode = locationCode
				addRequest.DeliveryIntent.DeliveryLocationType = locationType
			}
		}
	}

	cart, err = behaviour.AddToCart(ctx, cart, addRequest)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)
		return err, nil
	}
	cs.publishAddtoCartEvent(ctx, *cart, addRequest)
	cs.updateCartInCache(ctx, session, cart)
	return nil, product
}

func (cs *CartService) ApplyVoucher(ctx context.Context, session *sessions.Session, couponCode string) (*cartDomain.Cart, error) {

	oldCart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "ApplyVoucher").Error(err)
		return nil, err
	}

	cart, err := behaviour.ApplyVoucher(ctx, oldCart, couponCode)
	cs.updateCartInCache(ctx, session, cart)
	return cart, err
}

func (cs *CartService) handleCartNotFound(session *sessions.Session, err error) {
	if err == cartDomain.CartNotFoundError {
		cs.DeleteSavedSessionGuestCartId(session)
	}
}

// checkProductForAddRequest existence and validate with productService
func (cs *CartService) checkProductForAddRequest(ctx context.Context, session *sessions.Session, addRequest cartDomain.AddRequest) (cartDomain.AddRequest, productDomain.BasicProduct, error) {
	product, err := cs.ProductService.Get(ctx, addRequest.MarketplaceCode)
	if err != nil {
		return addRequest, nil, fmt.Errorf("cart.application.cartservice - AddProduct Error: %v", err)
	}

	if product.Type() == productDomain.TYPECONFIGURABLE {
		if addRequest.VariantMarketplaceCode == "" {
			return addRequest, nil, errors.New("cart.application.cartservice - AddProduct:No Variant given for configurable product")
		}

		configurableProduct := product.(productDomain.ConfigurableProduct)
		if !configurableProduct.HasVariant(addRequest.VariantMarketplaceCode) {
			return addRequest, nil, errors.New("cart.application.cartservice - AddProduct:Product has not the given variant")
		}
	}

	//Now Validate the Item with the optional registered ItemValidator
	if cs.ItemValidator != nil {
		return addRequest, product, cs.ItemValidator.Validate(ctx, session, addRequest, product)
	}

	return addRequest, product, nil
}

func (cs *CartService) publishAddtoCartEvent(ctx context.Context, currentCart cartDomain.Cart, addRequest cartDomain.AddRequest) {
	if cs.EventPublisher != nil {
		cs.EventPublisher.PublishAddToCartEvent(ctx, addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty)
	}
}

func (cs *CartService) updateCartInCache(ctx context.Context, session *sessions.Session, cart *cartDomain.Cart) {
	if cs.CartCache != nil {
		id, err := BuildIdentifierFromCart(cart)
		if err != nil {
			return
		}
		err = cs.CartCache.CacheCart(ctx, session, *id, cart)
		if err != nil {
			cs.Logger.Error("Error while caching cart: %v", err)
		}
	}
}

func (cs *CartService) deleteCartInCache(ctx context.Context, session *sessions.Session, cart *cartDomain.Cart) {
	if cs.CartCache != nil {
		id, err := BuildIdentifierFromCart(cart)
		if err != nil {
			return
		}
		err = cs.CartCache.Delete(ctx, session, *id)
		if err != nil {
			cs.Logger.Error("Error while deleting cart in cache: %v", err)
		}
	}
}
