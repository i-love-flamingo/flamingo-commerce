package application

import (
	"fmt"

	"github.com/pkg/errors"
	cartDomain "go.aoe.com/flamingo/core/cart/domain/cart"
	productDomain "go.aoe.com/flamingo/core/product/domain"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
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
func (cs CartService) ValidateCart(ctx web.Context, decoratedCart *cartDomain.DecoratedCart) cartDomain.CartValidationResult {

	if cs.CartValidator != nil {
		result := cs.CartValidator.Validate(ctx, decoratedCart)
		return result
	}

	return cartDomain.CartValidationResult{}
}

// ValidateCurrentCart validates the current active cart
func (cs CartService) ValidateCurrentCart(ctx web.Context) (cartDomain.CartValidationResult, error) {
	decoratedCart, err := cs.CartReceiverService.ViewDecoratedCart(ctx)
	if err != nil {
		return cartDomain.CartValidationResult{}, err
	}

	return cs.ValidateCart(ctx, decoratedCart), nil
}

func (cs CartService) UpdateDeliveryInfosAndBilling(ctx web.Context, billingAddress *cartDomain.Address, updateCommands []cartDomain.DeliveryInfoUpdateCommand) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
	if err != nil {
		return err
	}
	cart, err = behaviour.UpdateDeliveryInfosAndBilling(ctx, cart, billingAddress, updateCommands)
	if err != nil {
		cs.handleCartNotFound(ctx, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, cart)
	return nil
}
func (cs CartService) UpdatePurchaser(ctx web.Context, purchaser *cartDomain.Person, additionalData map[string]string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
	if err != nil {
		return err
	}
	cart, err = behaviour.UpdatePurchaser(ctx, cart, purchaser, additionalData)
	if err != nil {
		cs.handleCartNotFound(ctx, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, cart)
	return nil
}

// UpdateItemQty
func (cs CartService) UpdateItemQty(ctx web.Context, itemId string, qty int) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
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
		return cs.DeleteItem(ctx, itemId)
	}

	cs.EventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, qty, cart.ID)
	itemUpdate := cartDomain.ItemUpdateCommand{
		Qty: &qty,
	}
	cart, err = behaviour.UpdateItem(ctx, cart, itemId, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(ctx, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, cart)
	return nil
}

func (cs CartService) UpdateItemSourceId(ctx web.Context, itemId string, sourceId string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
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
		cs.handleCartNotFound(ctx, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemSourceId").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, cart)
	return nil
}

// DeleteItem in current cart
func (cs CartService) DeleteItem(ctx web.Context, itemId string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
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
		cs.handleCartNotFound(ctx, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteItem").Error(err)
		return err
	}
	cs.updateCartInCache(ctx, cart)
	return nil
}

// PlaceOrder
func (cs *CartService) PlaceOrder(ctx web.Context, payment *cartDomain.CartPayment) (string, error) {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
	if err != nil {
		return "", err
	}

	orderNumber, err := behaviour.PlaceOrder(ctx, cart, payment)
	if err != nil {
		cs.handleCartNotFound(ctx, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "PlaceOrder").Error(err)
		return "", err
	}
	cs.EventPublisher.PublishOrderPlacedEvent(ctx, cart, orderNumber)
	cs.DeleteSessionGuestCart(ctx)
	cs.deleteCartInCache(ctx, cart)
	return orderNumber, err
}

// BuildAddRequest Helper to build
func (cs *CartService) BuildAddRequest(ctx web.Context, marketplaceCode string, variantMarketplaceCode string, qty int, deliveryIntentStringRepresentation string) cartDomain.AddRequest {
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
func (cs *CartService) AddProduct(ctx web.Context, addRequest cartDomain.AddRequest) error {
	addRequest, product, err := cs.checkProductForAddRequest(ctx, addRequest)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)
		return err
	}

	cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Debugf("AddRequest received %#v  / %v", addRequest, addRequest.DeliveryIntent.String())

	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)
		return err
	}

	//Check if we can autodetect empty location code for pickup
	if addRequest.DeliveryIntent.Method == cartDomain.DELIVERY_METHOD_PICKUP && addRequest.DeliveryIntent.DeliveryLocationCode == "" && addRequest.DeliveryIntent.AutodetectDeliveryLocation {
		if cs.PickUpDetectionService != nil {
			locationCode, locationType, err := cs.PickUpDetectionService.Detect(product, addRequest)
			if err == nil {
				cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Debugf("Detected pickup location %v / %v", locationCode, locationType)
				addRequest.DeliveryIntent.DeliveryLocationCode = locationCode
				addRequest.DeliveryIntent.DeliveryLocationType = locationType
			}
		}
	}

	cart, err = behaviour.AddToCart(ctx, cart, addRequest)
	if err != nil {
		cs.handleCartNotFound(ctx, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)
		return err
	}
	cs.publishAddtoCartEvent(ctx, *cart, addRequest)
	cs.updateCartInCache(ctx, cart)
	return nil
}

func (cs *CartService) ApplyVoucher(ctx web.Context, couponCode string) (*cartDomain.Cart, error) {

	oldCart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "ApplyVoucher").Error(err)
		return nil, err
	}

	cart, err := behaviour.ApplyVoucher(ctx, oldCart, couponCode)
	cs.updateCartInCache(ctx, cart)
	return cart, err
}

func (cs *CartService) handleCartNotFound(ctx web.Context, err error) {
	if err == cartDomain.CartNotFoundError {
		cs.DeleteSessionGuestCart(ctx)
	}
}

// checkProductForAddRequest existence and validate with productService
func (cs *CartService) checkProductForAddRequest(ctx web.Context, addRequest cartDomain.AddRequest) (cartDomain.AddRequest, productDomain.BasicProduct, error) {
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
		return addRequest, product, cs.ItemValidator.Validate(ctx, addRequest, product)
	}

	return addRequest, product, nil
}

func (cs *CartService) publishAddtoCartEvent(ctx web.Context, currentCart cartDomain.Cart, addRequest cartDomain.AddRequest) {
	if cs.EventPublisher != nil {
		cs.EventPublisher.PublishAddToCartEvent(ctx, addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty)
	}
}

func (cs *CartService) updateCartInCache(ctx web.Context, cart *cartDomain.Cart) {
	if cs.CartCache != nil {
		id, err := BuildIdentifierFromCart(cart)
		if err != nil {
			return
		}
		err = cs.CartCache.CacheCart(ctx, *id, cart)
		if err != nil {
			cs.Logger.Errorf("Error while caching cart: %v", err)
		}
	}
}

func (cs *CartService) deleteCartInCache(ctx web.Context, cart *cartDomain.Cart) {
	if cs.CartCache != nil {
		id, err := BuildIdentifierFromCart(cart)
		if err != nil {
			return
		}
		err = cs.CartCache.Delete(ctx, *id)
		if err != nil {
			cs.Logger.Errorf("Error while deleting cart in cache: %v", err)
		}
	}
}
