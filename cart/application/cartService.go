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

		DefaultDeliveryCode string `inject:"config:cart.defaultDeliveryCode,optional"`

		DeliveryInfoBuilder cartDomain.DeliveryInfoBuilder `inject:""`

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

//UpdateBillingAddress updates the billing address on the cart
func (cs CartService) UpdateBillingAddress(ctx context.Context, session *sessions.Session, billingAddress *cartDomain.Address) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	cart, err = behaviour.UpdateBillingAddress(ctx, cart, billingAddress)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	return nil
}

// UpdateDeliveryInfo updates the delivery info on the cart
func (cs CartService) UpdateDeliveryInfo(ctx context.Context, session *sessions.Session, deliveryCode string, deliveryInfo cartDomain.DeliveryInfo) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.DefaultDeliveryCode
	}

	cart, err = behaviour.UpdateDeliveryInfo(ctx, cart, deliveryCode, deliveryInfo)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	return nil
}

// UpdatePurchaser updates the purchaser on the cart
func (cs CartService) UpdatePurchaser(ctx context.Context, session *sessions.Session, purchaser *cartDomain.Person, additionalData map[string]string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	cart, err = behaviour.UpdatePurchaser(ctx, cart, purchaser, additionalData)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	return nil
}

// UpdateItemQty updates a single cart item qty
func (cs CartService) UpdateItemQty(ctx context.Context, session *sessions.Session, itemId string, deliveryCode string, qty int) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.DefaultDeliveryCode
	}

	item, err := cart.GetByItemId(itemId, deliveryCode)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	qtyBefore := item.Qty
	if qty < 1 {
		return cs.DeleteItem(ctx, session, itemId, deliveryCode)
	}

	cs.EventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, qty, cart.ID)
	itemUpdate := cartDomain.ItemUpdateCommand{
		Qty: &qty,
	}

	cart, err = behaviour.UpdateItem(ctx, cart, itemId, deliveryCode, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	return nil
}

// UpdateItemSourceId updates an item source id
func (cs CartService) UpdateItemSourceId(ctx context.Context, session *sessions.Session, itemId string, deliveryCode string, sourceId string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.DefaultDeliveryCode
	}
	_, err = cart.GetByItemId(itemId, deliveryCode)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemSourceId").Error(err)

		return err
	}

	itemUpdate := cartDomain.ItemUpdateCommand{
		SourceId: &sourceId,
	}

	cart, err = behaviour.UpdateItem(ctx, cart, itemId, deliveryCode, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemSourceId").Error(err)

		return err
	}

	return nil
}

// DeleteItem in current cart
func (cs CartService) DeleteItem(ctx context.Context, session *sessions.Session, itemId string, deliveryCode string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.DefaultDeliveryCode
	}

	item, err := cart.GetByItemId(itemId, deliveryCode)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteItem").Error(err)

		return err
	}

	qtyBefore := item.Qty
	cs.EventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, 0, cart.ID)

	cart, err = behaviour.DeleteItem(ctx, cart, itemId, deliveryCode)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteItem").Error(err)

		return err
	}

	return nil
}

// DeleteAllItems in current cart
func (cs CartService) DeleteAllItems(ctx context.Context, session *sessions.Session) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	for _, delivery := range cart.Deliveries {
		for _, item := range delivery.Cartitems {
			qtyBefore := item.Qty
			cs.EventPublisher.PublishChangedQtyInCartEvent(ctx, &item, qtyBefore, 0, cart.ID)

			cart, err = behaviour.DeleteItem(ctx, cart, item.ID, delivery.DeliveryInfo.Code)
			if err != nil {
				cs.handleCartNotFound(session, err)
				cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteAllItems").Error(err)

				return err
			}
		}
	}

	return nil
}

// PlaceOrder submits the order
func (cs *CartService) PlaceOrder(ctx context.Context, session *sessions.Session, payment *cartDomain.CartPayment) (cartDomain.PlacedOrderInfos, error) {
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
	cs.DeleteSavedSessionGuestCartID(session)
	cs.deleteCartInCache(ctx, session, cart)

	return orderNumbers, err
}

// BuildAddRequest Helper to build
func (cs *CartService) BuildAddRequest(ctx context.Context, marketplaceCode string, variantMarketplaceCode string, qty int) cartDomain.AddRequest {
	if qty < 0 {
		qty = 0
	}

	return cartDomain.AddRequest{
		MarketplaceCode: marketplaceCode,
		Qty:             qty,
		VariantMarketplaceCode: variantMarketplaceCode,
	}
}

// AddProduct adds a product to the cart
func (cs *CartService) AddProduct(ctx context.Context, session *sessions.Session, deliveryCode string, addRequest cartDomain.AddRequest) (error, productDomain.BasicProduct) {
	if deliveryCode == "" {
		deliveryCode = cs.DefaultDeliveryCode
	}

	addRequest, product, err := cs.checkProductForAddRequest(ctx, session, deliveryCode, addRequest)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return err, nil
	}

	cs.Logger.WithField(flamingo.LogKeyCategory, "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Debug("AddRequest received %#v  / %v", addRequest, deliveryCode)

	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return err, nil
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	cart, err = cs.CreateInitialDeliveryIfNotPresent(ctx, session, deliveryCode)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return err, nil
	}

	cart, err = behaviour.AddToCart(ctx, cart, deliveryCode, addRequest)
	if err == cartDomain.DeliveryCodeNotFound {
		//old Magento adapter will never return this
		//Edge case...
		// For later - todo:
		// call initialCreateDelivery
		// retry AddToCart
	}

	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)

		return err, nil
	}

	cs.publishAddtoCartEvent(ctx, *cart, addRequest)

	return nil, product
}

// CreateInitialDeliveryIfNotPresent creates the initial delivery
func (cs *CartService) CreateInitialDeliveryIfNotPresent(ctx context.Context, session *sessions.Session, deliveryCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	if cart.HasDeliveryForCode(deliveryCode) {
		return cart, nil
	}

	delInfo, err := cs.DeliveryInfoBuilder.BuildByDeliveryCode(deliveryCode)
	if err != nil {
		return nil, err
	}

	return behaviour.UpdateDeliveryInfo(ctx, cart, deliveryCode, delInfo)
}

// ApplyVoucher applies a voucher to the cart
func (cs *CartService) ApplyVoucher(ctx context.Context, session *sessions.Session, couponCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx, session)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "ApplyVoucher").Error(err)

		return nil, err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	cart, err = behaviour.ApplyVoucher(ctx, cart, couponCode)

	return cart, err
}

func (cs *CartService) handleCartNotFound(session *sessions.Session, err error) {
	if err == cartDomain.CartNotFoundError {
		cs.DeleteSavedSessionGuestCartID(session)
	}
}

// checkProductForAddRequest existence and validate with productService
func (cs *CartService) checkProductForAddRequest(ctx context.Context, session *sessions.Session, deliveryCode string, addRequest cartDomain.AddRequest) (cartDomain.AddRequest, productDomain.BasicProduct, error) {
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
		return addRequest, product, cs.ItemValidator.Validate(ctx, session, deliveryCode, addRequest, product)
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
