package application

import (
	"context"
	"fmt"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

// CartService application struct
type (
	//CartService provides methods to modify the cart
	CartService struct {
		cartReceiverService *CartReceiverService
		productService      productDomain.ProductService
		eventPublisher      EventPublisher
		deliveryInfoBuilder cartDomain.DeliveryInfoBuilder
		logger              flamingo.Logger
		defaultDeliveryCode string
		// optionals - these may be nil
		cartValidator cartDomain.CartValidator
		itemValidator cartDomain.ItemValidator
		cartCache     CartCache
	}
)

// Inject dependencies
func (cs *CartService) Inject(
	cartReceiverService *CartReceiverService,
	productService productDomain.ProductService,
	eventPublisher EventPublisher,
	deliveryInfoBuilder cartDomain.DeliveryInfoBuilder,
	logger flamingo.Logger,
	config *struct {
		DefaultDeliveryCode string `inject:"config:cart.defaultDeliveryCode,optional"`
	},
	cartValidator cartDomain.CartValidator,
	itemValidator cartDomain.ItemValidator,
	cartCache CartCache,
) {
	cs.cartReceiverService = cartReceiverService
	cs.productService = productService
	cs.eventPublisher = eventPublisher
	cs.deliveryInfoBuilder = deliveryInfoBuilder
	cs.logger = logger
	if config != nil {
		cs.defaultDeliveryCode = config.DefaultDeliveryCode
	}
	cs.cartValidator = cartValidator
	cs.itemValidator = itemValidator
	cs.cartCache = cartCache
}

// GetCartReceiverService returns the injected cart receiver service
func (cs *CartService) GetCartReceiverService() *CartReceiverService {
	return cs.cartReceiverService
}

// ValidateCart validates a carts content
func (cs *CartService) ValidateCart(ctx context.Context, session *web.Session, decoratedCart *cartDomain.DecoratedCart) cartDomain.CartValidationResult {

	if cs.cartValidator != nil {
		result := cs.cartValidator.Validate(ctx, session, decoratedCart)

		return result
	}

	return cartDomain.CartValidationResult{}
}

// ValidateCurrentCart validates the current active cart
func (cs *CartService) ValidateCurrentCart(ctx context.Context, session *web.Session) (cartDomain.CartValidationResult, error) {
	decoratedCart, err := cs.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		return cartDomain.CartValidationResult{}, err
	}

	return cs.ValidateCart(ctx, session, decoratedCart), nil
}

//UpdateBillingAddress updates the billing address on the cart
func (cs *CartService) UpdateBillingAddress(ctx context.Context, session *web.Session, billingAddress *cartDomain.Address) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
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
		cs.logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	return nil
}

// UpdateDeliveryInfo updates the delivery info on the cart
func (cs *CartService) UpdateDeliveryInfo(ctx context.Context, session *web.Session, deliveryCode string, deliveryInfo cartDomain.DeliveryInfo) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	cart, err = behaviour.UpdateDeliveryInfo(ctx, cart, deliveryCode, deliveryInfo)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	return nil
}

// UpdatePurchaser updates the purchaser on the cart
func (cs *CartService) UpdatePurchaser(ctx context.Context, session *web.Session, purchaser *cartDomain.Person, additionalData *cartDomain.AdditionalData) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
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
		cs.logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	return nil
}

// UpdateItemQty updates a single cart item qty
func (cs *CartService) UpdateItemQty(ctx context.Context, session *web.Session, itemID string, deliveryCode string, qty int) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	item, err := cart.GetByItemId(itemID, deliveryCode)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	qtyBefore := item.Qty
	if qty < 1 {
		return cs.DeleteItem(ctx, session, itemID, deliveryCode)
	}

	cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, qty, cart.ID)
	itemUpdate := cartDomain.ItemUpdateCommand{
		Qty: &qty,
	}

	cart, err = behaviour.UpdateItem(ctx, cart, itemID, deliveryCode, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	return nil
}

// UpdateItemSourceID updates an item source id
func (cs *CartService) UpdateItemSourceID(ctx context.Context, session *web.Session, itemID string, deliveryCode string, sourceID string) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}
	_, err = cart.GetByItemId(itemID, deliveryCode)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemSourceId").Error(err)

		return err
	}

	itemUpdate := cartDomain.ItemUpdateCommand{
		SourceID: &sourceID,
	}

	cart, err = behaviour.UpdateItem(ctx, cart, itemID, deliveryCode, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemSourceId").Error(err)

		return err
	}

	return nil
}

// DeleteItem in current cart
func (cs *CartService) DeleteItem(ctx context.Context, session *web.Session, itemID string, deliveryCode string) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	item, err := cart.GetByItemId(itemID, deliveryCode)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField("subCategory", "DeleteItem").Error(err)

		return err
	}

	qtyBefore := item.Qty
	cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, 0, cart.ID)

	cart, err = behaviour.DeleteItem(ctx, cart, itemID, deliveryCode)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithField("category", "cartService").WithField("subCategory", "DeleteItem").Error(err)

		return err
	}

	return nil
}

// DeleteAllItems in current cart
func (cs *CartService) DeleteAllItems(ctx context.Context, session *web.Session) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
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
			cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, &item, qtyBefore, 0, cart.ID)

			cart, err = behaviour.DeleteItem(ctx, cart, item.ID, delivery.DeliveryInfo.Code)
			if err != nil {
				cs.handleCartNotFound(session, err)
				cs.logger.WithField("category", "cartService").WithField("subCategory", "DeleteAllItems").Error(err)

				return err
			}
		}
	}

	return nil
}

// Clean current cart
func (cs *CartService) Clean(ctx context.Context, session *web.Session) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
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
			cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, &item, qtyBefore, 0, cart.ID)
		}
	}

	_, err = behaviour.CleanCart(ctx, cart)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField("subCategory", "DeleteAllItems").Error(err)
		return err
	}

	return nil
}

// DeleteDelivery in current cart
func (cs *CartService) DeleteDelivery(ctx context.Context, session *web.Session, deliveryCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	delivery, found := cart.GetDeliveryByCode(deliveryCode)
	if !found {
		return nil, errors.New("delivery not found: " + deliveryCode)
	}
	for _, item := range delivery.Cartitems {
		qtyBefore := item.Qty
		cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, &item, qtyBefore, 0, cart.ID)
	}

	cart, err = behaviour.CleanDelivery(ctx, cart, deliveryCode)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField("subCategory", "DeleteAllItems").Error(err)
		return nil, err
	}

	return cart, nil
}

// BuildAddRequest Helper to build
func (cs *CartService) BuildAddRequest(ctx context.Context, marketplaceCode string, variantMarketplaceCode string, qty int) cartDomain.AddRequest {
	if qty < 0 {
		qty = 0
	}

	return cartDomain.AddRequest{
		MarketplaceCode:        marketplaceCode,
		Qty:                    qty,
		VariantMarketplaceCode: variantMarketplaceCode,
	}
}

// AddProduct adds a product to the cart
func (cs *CartService) AddProduct(ctx context.Context, session *web.Session, deliveryCode string, addRequest cartDomain.AddRequest) (productDomain.BasicProduct, error) {
	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	addRequest, product, err := cs.checkProductForAddRequest(ctx, session, deliveryCode, addRequest)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return nil, err
	}

	cs.logger.WithField(flamingo.LogKeyCategory, "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Debug("AddRequest received %#v  / %v", addRequest, deliveryCode)

	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return nil, err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	cart, err = cs.CreateInitialDeliveryIfNotPresent(ctx, session, deliveryCode)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return nil, err
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
		cs.logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)

		return nil, err
	}

	cs.publishAddtoCartEvent(ctx, *cart, addRequest)

	return product, nil
}

// CreateInitialDeliveryIfNotPresent creates the initial delivery
func (cs *CartService) CreateInitialDeliveryIfNotPresent(ctx context.Context, session *web.Session, deliveryCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	if cart.HasDeliveryForCode(deliveryCode) {
		return cart, nil
	}

	delInfo, err := cs.deliveryInfoBuilder.BuildByDeliveryCode(deliveryCode)
	if err != nil {
		return nil, err
	}

	return behaviour.UpdateDeliveryInfo(ctx, cart, deliveryCode, delInfo)
}

// ApplyVoucher applies a voucher to the cart
func (cs *CartService) ApplyVoucher(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		cs.logger.WithField("category", "cartService").WithField("subCategory", "ApplyVoucher").Error(err)

		return nil, err
	}
	// cart cache must be updated - with the current value of cart
	defer func() {
		cs.updateCartInCache(ctx, session, cart)
	}()

	cart, err = behaviour.ApplyVoucher(ctx, cart, couponCode)

	return cart, err
}

func (cs *CartService) handleCartNotFound(session *web.Session, err error) {
	if err == cartDomain.CartNotFoundError {
		cs.DeleteSavedSessionGuestCartID(session)
	}
}

// checkProductForAddRequest existence and validate with productService
func (cs *CartService) checkProductForAddRequest(ctx context.Context, session *web.Session, deliveryCode string, addRequest cartDomain.AddRequest) (cartDomain.AddRequest, productDomain.BasicProduct, error) {
	product, err := cs.productService.Get(ctx, addRequest.MarketplaceCode)
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
	if cs.itemValidator != nil {
		return addRequest, product, cs.itemValidator.Validate(ctx, session, deliveryCode, addRequest, product)
	}

	return addRequest, product, nil
}

func (cs *CartService) publishAddtoCartEvent(ctx context.Context, currentCart cartDomain.Cart, addRequest cartDomain.AddRequest) {
	if cs.eventPublisher != nil {
		cs.eventPublisher.PublishAddToCartEvent(ctx, addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty)
	}
}

func (cs *CartService) updateCartInCache(ctx context.Context, session *web.Session, cart *cartDomain.Cart) {
	if cs.cartCache != nil {
		id, err := BuildIdentifierFromCart(cart)
		if err != nil {
			return
		}

		err = cs.cartCache.CacheCart(ctx, session, *id, cart)
		if err != nil {
			cs.logger.Error("Error while caching cart: %v", err)
		}
	}
}

// DeleteCartInCache removes the cart from cache
func (cs *CartService) DeleteCartInCache(ctx context.Context, session *web.Session, cart *cartDomain.Cart) {
	if cs.cartCache != nil {
		id, err := BuildIdentifierFromCart(cart)
		if err != nil {
			return
		}

		err = cs.cartCache.Delete(ctx, session, *id)
		if err != nil {
			cs.logger.Error("Error while deleting cart in cache: %v", err)
		}
	}
}
