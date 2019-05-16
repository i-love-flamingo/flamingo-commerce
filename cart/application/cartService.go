package application

import (
	"context"
	"fmt"

	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo/v3/core/oauth/application"

	"github.com/pkg/errors"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// CartService application struct
type (
	// CartService provides methods to modify the cart
	CartService struct {
		cartReceiverService *CartReceiverService
		authManager         *application.AuthManager
		productService      productDomain.ProductService
		eventPublisher      events.EventPublisher
		eventRouter         flamingo.EventRouter
		deliveryInfoBuilder cartDomain.DeliveryInfoBuilder
		logger              flamingo.Logger
		defaultDeliveryCode string
		restrictionService  *validation.RestrictionService
		deleteEmptyDelivery bool
		// optionals - these may be nil
		cartValidator     validation.Validator
		itemValidator     validation.ItemValidator
		cartCache         CartCache
		placeOrderService placeorder.Service
	}

	// RestrictionError error enriched with result of restrictions
	RestrictionError struct {
		message           string
		RestrictionResult validation.RestrictionResult
	}
)

// Error fetch error message
func (e *RestrictionError) Error() string {
	return e.message
}

// Inject dependencies
func (cs *CartService) Inject(
	cartReceiverService *CartReceiverService,
	productService productDomain.ProductService,
	eventPublisher events.EventPublisher,
	eventRouter flamingo.EventRouter,
	deliveryInfoBuilder cartDomain.DeliveryInfoBuilder,
	restrictionService *validation.RestrictionService,
	authManager *application.AuthManager,
	logger flamingo.Logger,
	config *struct {
		DefaultDeliveryCode string `inject:"config:commerce.cart.defaultDeliveryCode,optional"`
		DeleteEmptyDelivery bool   `inject:"config:commerce.cart.deleteEmptyDelivery,optional"`
	},
	optionals *struct {
		CartValidator     validation.Validator     `inject:",optional"`
		ItemValidator     validation.ItemValidator `inject:",optional"`
		CartCache         CartCache                `inject:",optional"`
		PlaceOrderService placeorder.Service       `inject:",optional"`
	},
) {
	cs.cartReceiverService = cartReceiverService
	cs.productService = productService
	cs.eventPublisher = eventPublisher
	cs.eventRouter = eventRouter
	cs.deliveryInfoBuilder = deliveryInfoBuilder
	cs.restrictionService = restrictionService
	cs.authManager = authManager
	cs.logger = logger.WithField("module", "cart").WithField("category", "application.cartService")
	if config != nil {
		cs.defaultDeliveryCode = config.DefaultDeliveryCode
		cs.deleteEmptyDelivery = config.DeleteEmptyDelivery
	}
	if optionals != nil {
		cs.cartValidator = optionals.CartValidator
		cs.itemValidator = optionals.ItemValidator
		cs.cartCache = optionals.CartCache
		cs.placeOrderService = optionals.PlaceOrderService
	}
}

// GetCartReceiverService returns the injected cart receiver service
func (cs *CartService) GetCartReceiverService() *CartReceiverService {
	return cs.cartReceiverService
}

// ValidateCart validates a carts content
func (cs *CartService) ValidateCart(ctx context.Context, session *web.Session, decoratedCart *decorator.DecoratedCart) validation.Result {

	if cs.cartValidator != nil {
		result := cs.cartValidator.Validate(ctx, session, decoratedCart)

		return result
	}

	return validation.Result{}
}

// ValidateCurrentCart validates the current active cart
func (cs *CartService) ValidateCurrentCart(ctx context.Context, session *web.Session) (validation.Result, error) {
	decoratedCart, err := cs.cartReceiverService.ViewDecoratedCart(ctx, session)
	if err != nil {
		return validation.Result{}, err
	}

	return cs.ValidateCart(ctx, session, decoratedCart), nil
}

// UpdatePaymentSelection updates the paymentselection in the cart
func (cs *CartService) UpdatePaymentSelection(ctx context.Context, session *web.Session, paymentSelection cartDomain.PaymentSelection) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	cart, defers, err = behaviour.UpdatePaymentSelection(ctx, cart, paymentSelection)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdatePaymentSelection").Error(err)

		return err
	}

	return nil
}

// UpdateBillingAddress updates the billing address on the cart
func (cs *CartService) UpdateBillingAddress(ctx context.Context, session *web.Session, billingAddress *cartDomain.Address) error {
	if billingAddress == nil {
		return nil
	}
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	cart, defers, err = behaviour.UpdateBillingAddress(ctx, cart, *billingAddress)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdateBillingAddress").Error(err)

		return err
	}

	return nil
}

// UpdateDeliveryInfo updates the delivery info on the cart
func (cs *CartService) UpdateDeliveryInfo(ctx context.Context, session *web.Session, deliveryCode string, deliveryInfo cartDomain.DeliveryInfoUpdateCommand) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	cart, defers, err = behaviour.UpdateDeliveryInfo(ctx, cart, deliveryCode, deliveryInfo)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdateDeliveryInfo").Error(err)

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
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	cart, defers, err = behaviour.UpdatePurchaser(ctx, cart, purchaser, additionalData)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdatePurchaser").Error(err)

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
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	item, err := cart.GetByItemID(itemID)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	qtyBefore := item.Qty
	if qty < 1 {
		return cs.DeleteItem(ctx, session, itemID, deliveryCode)
	}

	product, err := cs.productService.Get(ctx, item.MarketplaceCode)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	err = cs.checkProductQtyRestrictions(ctx, product, cart, qty-qtyBefore)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdateItemQty").Error(err)

		return err
	}

	cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, qty, cart.ID)
	itemUpdate := cartDomain.ItemUpdateCommand{
		Qty: &qty,
	}

	cart, defers, err = behaviour.UpdateItem(ctx, cart, itemID, deliveryCode, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdateItemQty").Error(err)

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
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}
	_, err = cart.GetByItemID(itemID)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdateItemSourceId").Error(err)

		return err
	}

	itemUpdate := cartDomain.ItemUpdateCommand{
		SourceID: &sourceID,
	}

	cart, defers, err = behaviour.UpdateItem(ctx, cart, itemID, deliveryCode, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).WithField("subCategory", "UpdateItemSourceId").Error(err)

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
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)

		cs.handleEmptyDelivery(ctx, session, cart, deliveryCode)
		cs.dispatchAllEvents(ctx, defers)
	}()

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	item, err := cart.GetByItemID(itemID)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "DeleteItem").Error(err)

		return err
	}

	qtyBefore := item.Qty
	cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, 0, cart.ID)

	cart, defers, err = behaviour.DeleteItem(ctx, cart, itemID, deliveryCode)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).WithField("subCategory", "DeleteItem").Error(err)

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
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	for _, delivery := range cart.Deliveries {
		for _, item := range delivery.Cartitems {
			qtyBefore := item.Qty
			cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, &item, qtyBefore, 0, cart.ID)

			cart, defers, err = behaviour.DeleteItem(ctx, cart, item.ID, delivery.DeliveryInfo.Code)
			if err != nil {
				cs.handleCartNotFound(session, err)
				cs.logger.WithContext(ctx).WithField("subCategory", "DeleteAllItems").Error(err)

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
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	for _, delivery := range cart.Deliveries {
		for _, item := range delivery.Cartitems {
			qtyBefore := item.Qty
			cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, &item, qtyBefore, 0, cart.ID)
		}
	}

	_, defers, err = behaviour.CleanCart(ctx, cart)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "DeleteAllItems").Error(err)
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
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	delivery, found := cart.GetDeliveryByCode(deliveryCode)
	if !found {
		return nil, errors.New("delivery not found: " + deliveryCode)
	}
	for _, item := range delivery.Cartitems {
		qtyBefore := item.Qty
		cs.eventPublisher.PublishChangedQtyInCartEvent(ctx, &item, qtyBefore, 0, cart.ID)
	}

	cart, defers, err = behaviour.CleanDelivery(ctx, cart, deliveryCode)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "DeleteAllItems").Error(err)
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

	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return nil, err
	}
	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)

		cs.handleEmptyDelivery(ctx, session, cart, deliveryCode)
		cs.dispatchAllEvents(ctx, defers)
	}()

	addRequest, product, err := cs.checkProductForAddRequest(ctx, session, deliveryCode, addRequest)
	if err != nil {
		cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return nil, err
	}

	cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Debug(fmt.Sprintf("AddRequest received %#v  / %v", addRequest, deliveryCode))

	cart, err = cs.CreateInitialDeliveryIfNotPresent(ctx, session, deliveryCode)
	if err != nil {
		cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return nil, err
	}

	err = cs.checkProductQtyRestrictions(ctx, product, cart, addRequest.Qty)
	if err != nil {
		cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)

		return nil, err
	}

	cart, defers, err = behaviour.AddToCart(ctx, cart, deliveryCode, addRequest)
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).WithField("subCategory", "AddProduct").Error(err)

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

	updateCommand := cartDomain.DeliveryInfoUpdateCommand{
		DeliveryInfo: *delInfo,
	}

	info, defers, err := behaviour.UpdateDeliveryInfo(ctx, cart, deliveryCode, updateCommand)
	defer func() {
		cs.dispatchAllEvents(ctx, defers)
	}()

	return info, err
}

// GetInitialDelivery - calls the registered deliveryInfoBuilder to get the initial values for a Delivery based on the given code
func (cs *CartService) GetInitialDelivery(deliveryCode string) (*cartDomain.DeliveryInfo, error) {
	return cs.deliveryInfoBuilder.BuildByDeliveryCode(deliveryCode)
}

// ApplyVoucher applies a voucher to the cart
func (cs *CartService) ApplyVoucher(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "ApplyVoucher").Error(err)

		return nil, err
	}
	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	cart, defers, err = behaviour.ApplyVoucher(ctx, cart, couponCode)

	return cart, err
}

func (cs *CartService) handleCartNotFound(session *web.Session, err error) {
	if err == cartDomain.ErrCartNotFound {
		cs.DeleteSavedSessionGuestCartID(session)
	}
}

// checkProductForAddRequest existence and validate with productService
func (cs *CartService) checkProductForAddRequest(ctx context.Context, session *web.Session, deliveryCode string, addRequest cartDomain.AddRequest) (cartDomain.AddRequest, productDomain.BasicProduct, error) {
	product, err := cs.productService.Get(ctx, addRequest.MarketplaceCode)
	if err != nil {
		return addRequest, nil, fmt.Errorf("cart.application.cartservice - AddProduct Error: %v", err)
	}

	if product.Type() == productDomain.TypeConfigurable {
		if addRequest.VariantMarketplaceCode == "" {
			return addRequest, nil, errors.New("cart.application.cartservice - AddProduct:No Variant given for configurable product")
		}

		configurableProduct := product.(productDomain.ConfigurableProduct)
		if !configurableProduct.HasVariant(addRequest.VariantMarketplaceCode) {
			return addRequest, nil, errors.New("cart.application.cartservice - AddProduct:Product has not the given variant")
		}
	}

	// Now Validate the Item with the optional registered ItemValidator
	if cs.itemValidator != nil {
		return addRequest, product, cs.itemValidator.Validate(ctx, session, deliveryCode, addRequest, product)
	}

	return addRequest, product, nil
}

func (cs *CartService) checkProductQtyRestrictions(ctx context.Context, product productDomain.BasicProduct, cart *cartDomain.Cart, qtyToCheck int) error {
	restrictionResult := cs.restrictionService.RestrictQty(ctx, product, cart)

	if restrictionResult.IsRestricted {
		if qtyToCheck > restrictionResult.RemainingDifference {
			return &RestrictionError{
				message:           fmt.Sprintf("Can't update item quantity, product max quantity of %d would be exceeded", restrictionResult.MaxAllowed),
				RestrictionResult: *restrictionResult,
			}
		}
	}

	return nil
}

func (cs *CartService) publishAddtoCartEvent(ctx context.Context, currentCart cartDomain.Cart, addRequest cartDomain.AddRequest) {
	if cs.eventPublisher != nil {
		cs.eventPublisher.PublishAddToCartEvent(ctx, addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty)
	}
}

func (cs *CartService) updateCartInCacheIfCacheIsEnabled(ctx context.Context, session *web.Session, cart *cartDomain.Cart) {
	if cs.cartCache != nil && cart != nil {
		id, err := cs.cartCache.BuildIdentifier(ctx, session)
		if err != nil {
			return
		}
		if cart.BelongsToAuthenticatedUser != id.IsCustomerCart {
			cs.logger.WithContext(ctx).Error("Request to cache a cart with wrong idendifier. %v / %v", cart.BelongsToAuthenticatedUser, id.IsCustomerCart)
			return
		}

		err = cs.cartCache.CacheCart(ctx, session, id, cart)
		if err != nil {
			cs.logger.WithContext(ctx).Error("Error while caching cart: %v", err)
		}
	}
}

// DeleteCartInCache removes the cart from cache
func (cs *CartService) DeleteCartInCache(ctx context.Context, session *web.Session, cart *cartDomain.Cart) {
	if cs.cartCache != nil {
		id, err := cs.cartCache.BuildIdentifier(ctx, session)
		if err != nil {
			return
		}

		err = cs.cartCache.Delete(ctx, session, id)
		if err != nil {
			cs.logger.WithContext(ctx).Error("Error while deleting cart in cache: %v", err)
		}
	}
}

// ReserveOrderIDAndSave - reserves order id by using the PlaceOrder Behaviour and sets saves it on the cart. You may want to use this before proceeding with payment to ensure having a useful reference in the payment processing
func (cs *CartService) ReserveOrderIDAndSave(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	if cs.placeOrderService == nil {
		return nil, errors.New("No placeOrderService registered")
	}
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}
	reservedOrderID, err := cs.placeOrderService.ReserveOrderID(ctx, cart)
	if err != nil {
		cs.logger.WithContext(ctx).Debug("Reserve order id:", reservedOrderID)
		return nil, err
	}
	additionalData := cart.AdditionalData
	additionalData.ReservedOrderID = reservedOrderID
	data, defers, err := behaviour.UpdateAdditionalData(ctx, cart, &additionalData)
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()
	return data, err
}

// PlaceOrder converts the given cart with payments into orders by calling the Service
func (cs *CartService) PlaceOrder(ctx context.Context, session *web.Session, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	if cs.placeOrderService == nil {
		return nil, errors.New("No placeOrderService registered")
	}
	cart, _, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}
	var placeOrderInfos placeorder.PlacedOrderInfos
	if cs.cartReceiverService.IsLoggedIn(ctx, session) {
		auth, err := cs.authManager.Auth(ctx, session)
		if err != nil {
			return nil, err
		}
		placeOrderInfos, err = cs.placeOrderService.PlaceCustomerCart(ctx, auth, cart, payment)
	} else {
		placeOrderInfos, err = cs.placeOrderService.PlaceGuestCart(ctx, cart, payment)
	}
	if err != nil {
		cs.handleCartNotFound(session, err)
		cs.logger.WithContext(ctx).Error(err)
		return nil, err
	}

	cs.eventPublisher.PublishOrderPlacedEvent(ctx, cart, placeOrderInfos)
	cs.DeleteSavedSessionGuestCartID(session)
	cs.DeleteCartInCache(ctx, session, cart)

	return placeOrderInfos, err
}

// GetDefaultDeliveryCode returns the configured default deliverycode
func (cs *CartService) GetDefaultDeliveryCode() string {
	return cs.defaultDeliveryCode
}

// handleEmptyDelivery - delete an empty delivery when found and feature flag is set
func (cs *CartService) handleEmptyDelivery(ctx context.Context, session *web.Session, cart *cartDomain.Cart, deliveryCode string) {
	if cs.deleteEmptyDelivery != true {
		return
	}

	if cart == nil {
		return
	}

	delivery, found := cart.GetDeliveryByCode(deliveryCode)
	if !found {
		return
	}

	if len(delivery.Cartitems) > 0 {
		return
	}

	_, err := cs.DeleteDelivery(ctx, session, deliveryCode)
	if err != nil {
		cs.logger.WithContext(ctx).WithField("subCategory", "handleEmptyDelivery").Error(err)
		return
	}
}

func (cs *CartService) dispatchAllEvents(ctx context.Context, events []flamingo.Event) {
	for _, e := range events {
		cs.eventRouter.Dispatch(ctx, e)
	}
}
