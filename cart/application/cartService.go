package application

import (
	"context"
	"encoding/gob"
	"fmt"

	"flamingo.me/flamingo/v3/core/auth"
	"github.com/pkg/errors"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
)

// CartService application struct
type (
	// CartService provides methods to modify the cart
	CartService struct {
		cartReceiverService *CartReceiverService
		webIdentityService  *auth.WebIdentityService
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
		CartItemID        string
		RestrictionResult validation.RestrictionResult
	}

	// QtyAdjustmentResult restriction result enriched with the respective item
	QtyAdjustmentResult struct {
		OriginalItem          cartDomain.Item
		DeliveryCode          string
		WasDeleted            bool
		RestrictionResult     *validation.RestrictionResult
		NewQty                int
		HasRemovedCouponCodes bool
	}

	// QtyAdjustmentResults slice of QtyAdjustmentResult
	QtyAdjustmentResults []QtyAdjustmentResult

	// PromotionFunction type takes ctx, cart, couponCode and applies the promotion
	promotionFunc func(context.Context, *cartDomain.Cart, string) (*cartDomain.Cart, cartDomain.DeferEvents, error)

	contextKeyType string
)

const (
	itemIDKey contextKeyType = "item_id"
)

func ItemIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(itemIDKey).(string)

	return id
}

func ContextWithItemID(ctx context.Context, itemID string) context.Context {
	return context.WithValue(ctx, itemIDKey, itemID)
}

func init() {
	gob.Register(RestrictionError{})
	gob.Register(QtyAdjustmentResults{})
}

// Error fetch error message
func (e *RestrictionError) Error() string {
	return e.message
}

var (
	_                          Service = &CartService{}
	ErrBundleConfigNotProvided         = errors.New("error no bundle config provided")
	ErrProductNotTypeBundle            = errors.New("product not of type bundle")
)

// Inject dependencies
func (cs *CartService) Inject(
	cartReceiverService *CartReceiverService,
	productService productDomain.ProductService,
	eventPublisher events.EventPublisher,
	eventRouter flamingo.EventRouter,
	deliveryInfoBuilder cartDomain.DeliveryInfoBuilder,
	restrictionService *validation.RestrictionService,
	webIdentityService *auth.WebIdentityService,
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
	cs.webIdentityService = webIdentityService
	cs.logger = logger.WithField(flamingo.LogKeyModule, "cart").WithField(flamingo.LogKeyCategory, "application.cartService")
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

		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdatePaymentSelection").Error(err)
		}

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

		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateBillingAddress").Error(err)
		}

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

		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateDeliveryInfo").Error(err)
		}

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

		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdatePurchaser").Error(err)
		}

		return err
	}

	return nil
}

// UpdateItemQty updates a single cart item qty
func (cs *CartService) UpdateItemQty(ctx context.Context, session *web.Session, itemID string, deliveryCode string, qty int) error {
	if qty < 1 {
		// item needs to be removed, let DeleteItem handle caching/events
		return cs.DeleteItem(ctx, session, itemID, deliveryCode)
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

	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	item, err := cart.GetByItemID(itemID)
	if err != nil {

		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemQty").Error(err)
		}

		return err
	}

	qtyBefore := item.Qty

	product, err := cs.productService.Get(ctx, item.MarketplaceCode)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemQty").Error(err)
		}
		return err
	}

	product, err = cs.getSpecificProductType(ctx, product, item.VariantMarketPlaceCode, item.BundleConfig)
	if err != nil {
		return err
	}

	err = cs.checkProductQtyRestrictions(ctx, session, product, cart, qty-qtyBefore, deliveryCode, itemID)
	if err != nil {
		cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemQty").Info(err)

		return err
	}

	itemUpdate := cartDomain.ItemUpdateCommand{
		Qty:            &qty,
		ItemID:         itemID,
		AdditionalData: nil,
	}

	cart, defers, err = behaviour.UpdateItem(ctx, cart, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemQty").Error(err)
		}
		return err
	}

	// append deferred events of behaviour with changed qty event
	updateEvent := &events.ChangedQtyInCartEvent{
		Cart:                   cart,
		CartID:                 cart.ID,
		MarketplaceCode:        item.MarketplaceCode,
		VariantMarketplaceCode: item.VariantMarketPlaceCode,
		ProductName:            product.TeaserData().ShortTitle,
		QtyBefore:              qtyBefore,
		QtyAfter:               qty,
	}
	defers = append(defers, updateEvent)

	return nil
}

// UpdateItemSourceID updates an item source id
func (cs *CartService) UpdateItemSourceID(ctx context.Context, session *web.Session, itemID string, sourceID string) error {
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

	_, err = cart.GetByItemID(itemID)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemSourceId").Error(err)
		}
		return err
	}

	itemUpdate := cartDomain.ItemUpdateCommand{
		SourceID: &sourceID,
		ItemID:   itemID,
	}

	cart, defers, err = behaviour.UpdateItem(ctx, cart, itemUpdate)
	if err != nil {
		cs.handleCartNotFound(session, err)

		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemSourceId").Error(err)
		}

		return err
	}

	return nil
}

// UpdateItems updates multiple items
func (cs *CartService) UpdateItems(ctx context.Context, session *web.Session, updateCommands []cartDomain.ItemUpdateCommand) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}

	for _, command := range updateCommands {
		_, err := cart.GetByItemID(command.ItemID)
		if err != nil {
			return err
		}
	}

	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	cart, defers, err = behaviour.UpdateItems(ctx, cart, updateCommands)
	if err != nil {
		cs.handleCartNotFound(session, err)
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemSourceId").Error(err)
		}
		return err
	}

	return nil
}

// UpdateItemBundleConfig updates multiple item
func (cs *CartService) UpdateItemBundleConfig(ctx context.Context, session *web.Session, updateCommand cartDomain.ItemUpdateCommand) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}

	if updateCommand.BundleConfiguration == nil {
		return ErrBundleConfigNotProvided
	}

	item, err := cart.GetByItemID(updateCommand.ItemID)
	if err != nil {
		return err
	}

	product, err := cs.productService.Get(ctx, item.MarketplaceCode)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemBundleConfig").Error(err)
		}
		return err
	}

	product, err = cs.getBundleWithActiveChoices(ctx, product, updateCommand.BundleConfiguration)
	if err != nil {
		return fmt.Errorf("error converting product to bundle: %w", err)
	}

	if cs.itemValidator != nil {
		decoratedCart, _ := cs.cartReceiverService.DecorateCart(ctx, cart)
		delivery, err := cart.GetDeliveryByItemID(updateCommand.ItemID)
		if err != nil {
			return fmt.Errorf("delivery code not found by item, while updating bundle: %w", err)
		}

		if err := cs.itemValidator.Validate(ctx, session, decoratedCart, delivery.DeliveryInfo.Code, cartDomain.AddRequest{}, product); err != nil {
			return fmt.Errorf("error validating bundle update: %w", err)
		}
	}

	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	cart, defers, err = behaviour.UpdateItem(ctx, cart, updateCommand)
	if err != nil {
		cs.handleCartNotFound(session, err)
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "UpdateItemSourceId").Error(err)
		}
		return err
	}

	return nil
}

//nolint:wrapcheck // error wrapped in the code above, no need to wrap twice
func (cs *CartService) getBundleWithActiveChoices(_ context.Context, product productDomain.BasicProduct, bundleConfig productDomain.BundleConfiguration) (productDomain.BasicProduct, error) {
	var err error

	bundleProduct, ok := product.(productDomain.BundleProduct)
	if !ok {
		return nil, ErrProductNotTypeBundle
	}

	product, err = bundleProduct.GetBundleProductWithActiveChoices(bundleConfig)
	if err != nil {
		return nil, err
	}

	return product, nil
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
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "DeleteItem").Error(err)
		}
		return err
	}

	qtyBefore := item.Qty
	cart, defers, err = behaviour.DeleteItem(ctx, cart, itemID, deliveryCode)
	if err != nil {
		cs.handleCartNotFound(session, err)
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "DeleteItem").Error(errors.Wrap(err, "Trying to delete SKU :"+item.MarketplaceCode))
		}
		return err
	}

	// append deferred events of behaviour with changed qty event
	updateEvent := &events.ChangedQtyInCartEvent{
		Cart:                   cart,
		CartID:                 cart.ID,
		MarketplaceCode:        item.MarketplaceCode,
		VariantMarketplaceCode: item.VariantMarketPlaceCode,
		ProductName:            item.ProductName,
		QtyBefore:              qtyBefore,
		QtyAfter:               0,
	}
	defers = append(defers, updateEvent)

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
	var deleteItemEvents cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	// we throw a qty changed event for every item in delivery
	for _, delivery := range cart.Deliveries {
		for _, item := range delivery.Cartitems {
			updateEvent := &events.ChangedQtyInCartEvent{
				Cart:                   cart,
				CartID:                 cart.ID,
				MarketplaceCode:        item.MarketplaceCode,
				VariantMarketplaceCode: item.VariantMarketPlaceCode,
				ProductName:            item.ProductName,
				QtyBefore:              item.Qty,
				QtyAfter:               0,
			}
			deleteItemEvents = append(deleteItemEvents, updateEvent)

			cart, defers, err = behaviour.DeleteItem(ctx, cart, item.ID, delivery.DeliveryInfo.Code)
			if err != nil {
				cs.handleCartNotFound(session, err)
				if !errors.Is(err, context.Canceled) {
					cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "DeleteAllItems").Error(err)
				}
				return err
			}
		}
	}
	// append deferred events of behaviour with changed qty events
	defers = append(defers, deleteItemEvents)

	return nil
}

// CompleteCurrentCart and remove from cache
func (cs *CartService) CompleteCurrentCart(ctx context.Context) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, web.SessionFromContext(ctx))
	if err != nil {
		return nil, err
	}

	// dispatch potential events
	var defers cartDomain.DeferEvents
	defer func() {
		cs.dispatchAllEvents(ctx, defers)
	}()

	completeBehaviour, ok := behaviour.(cartDomain.CompleteBehaviour)
	if !ok {
		return nil, fmt.Errorf("not supported by used cart behaviour: %T", behaviour)
	}

	var completedCart *cartDomain.Cart
	completedCart, defers, err = completeBehaviour.Complete(ctx, cart)
	if err != nil {
		cs.handleCartNotFound(web.SessionFromContext(ctx), err)
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "CloseCurrentCart").Error(err)
		}
		return nil, err
	}
	cs.DeleteCartInCache(ctx, web.SessionFromContext(ctx), nil)

	session := web.SessionFromContext(ctx)
	if !cart.BelongsToAuthenticatedUser {
		session.Delete(GuestCartSessionKey)
	}

	return completedCart, nil
}

// RestoreCart and cache
func (cs *CartService) RestoreCart(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.Cart, error) {
	session := web.SessionFromContext(ctx)
	behaviour, err := cs.cartReceiverService.ModifyBehaviour(ctx)
	if err != nil {
		return nil, err
	}

	completeBehaviour, ok := behaviour.(cartDomain.CompleteBehaviour)
	if !ok {
		return nil, fmt.Errorf("not supported by used cart behaviour: %T", behaviour)
	}

	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	var restoredCart *cartDomain.Cart
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, restoredCart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	restoredCart, defers, err = completeBehaviour.Restore(ctx, cart)

	if err != nil {
		cs.handleCartNotFound(web.SessionFromContext(ctx), err)
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "RestoreCart").Error(err)
		}
		return nil, err
	}

	if !restoredCart.BelongsToAuthenticatedUser {
		session.Store(GuestCartSessionKey, restoredCart.ID)
	}

	return restoredCart, nil
}

// Clean current cart
func (cs *CartService) Clean(ctx context.Context, session *web.Session) error {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return err
	}
	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	var deleteItemEvents cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	// we throw a qty changed event for every item of each delivery
	for _, delivery := range cart.Deliveries {
		for _, item := range delivery.Cartitems {
			updateEvent := &events.ChangedQtyInCartEvent{
				Cart:                   cart,
				CartID:                 cart.ID,
				MarketplaceCode:        item.MarketplaceCode,
				VariantMarketplaceCode: item.VariantMarketPlaceCode,
				ProductName:            item.ProductName,
				QtyBefore:              item.Qty,
				QtyAfter:               0,
			}
			deleteItemEvents = append(deleteItemEvents, updateEvent)
		}
	}

	cart, defers, err = behaviour.CleanCart(ctx, cart)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "DeleteAllItems").Error(err)
		}
		return err
	}
	// append deferred events of behaviour with changed qty events for every deleted item
	defers = append(defers, deleteItemEvents)

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
	var deleteItemEvents cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	delivery, found := cart.GetDeliveryByCode(deliveryCode)
	if !found {
		return nil, errors.New("delivery not found: " + deliveryCode)
	}
	// we throw a qty changed event for every item in delivery// we throw a qty changed event for every item in delivery
	for _, item := range delivery.Cartitems {
		updateEvent := &events.ChangedQtyInCartEvent{
			Cart:                   cart,
			CartID:                 cart.ID,
			MarketplaceCode:        item.MarketplaceCode,
			VariantMarketplaceCode: item.VariantMarketPlaceCode,
			ProductName:            item.ProductName,
			QtyBefore:              item.Qty,
			QtyAfter:               0,
		}
		deleteItemEvents = append(deleteItemEvents, updateEvent)
	}

	cart, defers, err = behaviour.CleanDelivery(ctx, cart, deliveryCode)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "DeleteAllItems").Error(err)
		}
		return nil, err
	}

	// append deferred events of behaviour with changed qty events for every deleted item
	defers = append(defers, deleteItemEvents)

	return cart, nil
}

// BuildAddRequest Helper to build
// Deprecated: build your own add request
func (cs *CartService) BuildAddRequest(_ context.Context, marketplaceCode string, variantMarketplaceCode string, qty int, additionalData map[string]string) cartDomain.AddRequest {
	if qty < 0 {
		qty = 0
	}

	return cartDomain.AddRequest{
		MarketplaceCode:        marketplaceCode,
		Qty:                    qty,
		VariantMarketplaceCode: variantMarketplaceCode,
		AdditionalData:         additionalData,
	}
}

// AddProduct adds a product to the cart
func (cs *CartService) AddProduct(ctx context.Context, session *web.Session, deliveryCode string, addRequest cartDomain.AddRequest) (productDomain.BasicProduct, error) {
	if deliveryCode == "" {
		deliveryCode = cs.defaultDeliveryCode
	}

	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)
		}
		return nil, err
	}
	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)

		cs.handleEmptyDelivery(ctx, session, cart, deliveryCode)
		cs.dispatchAllEvents(ctx, defers)
	}()

	product, err := cs.checkProductForAddRequest(ctx, session, cart, deliveryCode, addRequest)

	switch err.(type) {
	case nil:
	case *validation.AddToCartNotAllowed:
		cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Info(err)
		return nil, err
	default:
		cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Warn(err)
		return nil, err
	}

	cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Debug(fmt.Sprintf("AddRequest received %#v  / %v", addRequest, deliveryCode))

	cart, err = cs.CreateInitialDeliveryIfNotPresent(ctx, session, deliveryCode)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)
		}
		return nil, err
	}

	err = cs.checkProductQtyRestrictions(ctx, session, product, cart, addRequest.Qty, deliveryCode, "")
	if err != nil {
		cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Info(err)

		return nil, err
	}

	cart, defers, err = behaviour.AddToCart(ctx, cart, deliveryCode, addRequest)
	if err != nil {
		cs.handleCartNotFound(session, err)
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "AddProduct").Error(err)
		}

		return nil, err
	}

	// append deferred events of behaviour with add to cart event
	addToCart := &events.AddToCartEvent{
		Cart:                   cart,
		MarketplaceCode:        addRequest.MarketplaceCode,
		VariantMarketplaceCode: addRequest.VariantMarketplaceCode,
		ProductName:            product.TeaserData().ShortTitle,
		Qty:                    addRequest.Qty,
	}
	defers = append(defers, addToCart)

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

	updatedCart, defers, err := behaviour.UpdateDeliveryInfo(ctx, cart, deliveryCode, updateCommand)
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, updatedCart)
		cs.dispatchAllEvents(ctx, defers)
	}()

	return updatedCart, err
}

// GetInitialDelivery - calls the registered deliveryInfoBuilder to get the initial values for a Delivery based on the given code
func (cs *CartService) GetInitialDelivery(deliveryCode string) (*cartDomain.DeliveryInfo, error) {
	return cs.deliveryInfoBuilder.BuildByDeliveryCode(deliveryCode)
}

// ApplyVoucher applies a voucher to the cart
func (cs *CartService) ApplyVoucher(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.getCartAndBehaviour(ctx, session, "ApplyVoucher")
	if err != nil {
		return nil, err
	}
	return cs.executeVoucherBehaviour(ctx, session, cart, couponCode, behaviour.ApplyVoucher)
}

// ApplyAny applies a voucher or giftcard to the cart
func (cs *CartService) ApplyAny(ctx context.Context, session *web.Session, anyCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.getCartAndBehaviour(ctx, session, "ApplyAny")
	if err != nil {
		return nil, err
	}
	if giftCardAndVoucherBehaviour, ok := behaviour.(cartDomain.GiftCardAndVoucherBehaviour); ok {
		return cs.executeVoucherBehaviour(ctx, session, cart, anyCode, giftCardAndVoucherBehaviour.ApplyAny)
	}
	return nil, errors.New("ApplyAny not supported")
}

// RemoveVoucher removes a voucher from the cart
func (cs *CartService) RemoveVoucher(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.getCartAndBehaviour(ctx, session, "RemoveVoucher")
	if err != nil {
		return nil, err
	}
	return cs.executeVoucherBehaviour(ctx, session, cart, couponCode, behaviour.RemoveVoucher)
}

// ApplyGiftCard adds a giftcard to the cart
func (cs *CartService) ApplyGiftCard(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.getCartAndBehaviour(ctx, session, "ApplyGiftCard")
	if err != nil {
		return nil, err
	}
	if giftCartBehaviour, ok := behaviour.(cartDomain.GiftCardBehaviour); ok {
		return cs.executeVoucherBehaviour(ctx, session, cart, couponCode, giftCartBehaviour.ApplyGiftCard)
	}
	return nil, errors.New("ApplyGiftCard not supported")
}

// RemoveGiftCard removes a giftcard from the cart
func (cs *CartService) RemoveGiftCard(ctx context.Context, session *web.Session, couponCode string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.getCartAndBehaviour(ctx, session, "RemoveGiftCard")
	if err != nil {
		return nil, err
	}
	if giftCartBehaviour, ok := behaviour.(cartDomain.GiftCardBehaviour); ok {
		return cs.executeVoucherBehaviour(ctx, session, cart, couponCode, giftCartBehaviour.RemoveGiftCard)
	}
	return nil, errors.New("RemoveGiftCard not supported")
}

// Get current cart from session and corresponding behaviour
func (cs *CartService) getCartAndBehaviour(ctx context.Context, session *web.Session, logKey string) (*cartDomain.Cart, cartDomain.ModifyBehaviour, error) {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, logKey).Error(err)
		}

		return nil, nil, err
	}
	return cart, behaviour, nil
}

// Executes provided behaviour regarding vouchers, this function serves to reduce duplicated code
// for voucher / giftcard behaviour as their internal logic is basically the same
func (cs *CartService) executeVoucherBehaviour(ctx context.Context, session *web.Session, cart *cartDomain.Cart, couponCode string, fn promotionFunc) (*cartDomain.Cart, error) {
	// cart cache must be updated - with the current value of cart
	var defers cartDomain.DeferEvents
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()
	cart, defers, err := fn(ctx, cart, couponCode)
	return cart, err
}

func (cs *CartService) handleCartNotFound(session *web.Session, err error) {
	if errors.Is(err, cartDomain.ErrCartNotFound) {
		_ = cs.DeleteSavedSessionGuestCartID(session)
	}
}

// checkProductForAddRequest existence and validate with productService
func (cs *CartService) checkProductForAddRequest(ctx context.Context, session *web.Session, cart *cartDomain.Cart, deliveryCode string, addRequest cartDomain.AddRequest) (productDomain.BasicProduct, error) {
	product, err := cs.productService.Get(ctx, addRequest.MarketplaceCode)
	if err != nil {
		return nil, err
	}

	if product.Type() == productDomain.TypeConfigurable {
		if addRequest.VariantMarketplaceCode == "" {
			return nil, ErrNoVariantForConfigurable
		}

		configurableProduct := product.(productDomain.ConfigurableProduct)
		if !configurableProduct.HasVariant(addRequest.VariantMarketplaceCode) {
			return nil, ErrVariantDoNotExist
		}

		product, err = configurableProduct.GetConfigurableWithActiveVariant(addRequest.VariantMarketplaceCode)
		if err != nil {
			return nil, err
		}
	}

	if product.Type() == productDomain.TypeBundle {
		if len(addRequest.BundleConfiguration) == 0 {
			return nil, ErrNoBundleConfigurationGiven
		}

		bundleProduct := product.(productDomain.BundleProduct)
		domainBundleConfig := addRequest.BundleConfiguration

		product, err = bundleProduct.GetBundleProductWithActiveChoices(domainBundleConfig)
		if err != nil {
			return nil, fmt.Errorf("error getting bundle with active choices: %w", err)
		}
	}

	// Now Validate the Item with the optional registered ItemValidator
	if cs.itemValidator != nil {
		decoratedCart, _ := cs.cartReceiverService.DecorateCart(ctx, cart)
		return product, cs.itemValidator.Validate(ctx, session, decoratedCart, deliveryCode, addRequest, product)
	}

	return product, nil
}

func (cs *CartService) checkProductQtyRestrictions(ctx context.Context, sess *web.Session, product productDomain.BasicProduct, cart *cartDomain.Cart, qtyToCheck int, deliveryCode string, itemID string) error {
	restrictionResult := cs.restrictionService.RestrictQty(ctx, sess, product, cart, deliveryCode)

	if restrictionResult.IsRestricted {
		if qtyToCheck > restrictionResult.RemainingDifference {
			return &RestrictionError{
				message:           fmt.Sprintf("Can't update item quantity, product max quantity of %d would be exceeded. Restrictor: %v", restrictionResult.MaxAllowed, restrictionResult.RestrictorName),
				CartItemID:        itemID,
				RestrictionResult: *restrictionResult,
			}
		}
	}

	return nil
}

func (cs *CartService) updateCartInCacheIfCacheIsEnabled(ctx context.Context, session *web.Session, cart *cartDomain.Cart) {
	if cs.cartCache != nil && cart != nil {
		id, err := cs.cartCache.BuildIdentifier(ctx, session)
		if err != nil {
			return
		}
		if cart.BelongsToAuthenticatedUser != id.IsCustomerCart {
			cs.logger.WithContext(ctx).Error(fmt.Sprintf("Request to cache a cart with wrong idendifier. Cart BelongsToAuthenticatedUser: %v / CacheIdendifier: %v", cart.BelongsToAuthenticatedUser, id.IsCustomerCart))
			return
		}

		err = cs.cartCache.CacheCart(ctx, session, id, cart)
		if err != nil {
			cs.logger.WithContext(ctx).Error("Error while caching cart: %v", err)
		}
	}
}

// DeleteCartInCache removes the cart from cache
func (cs *CartService) DeleteCartInCache(ctx context.Context, session *web.Session, _ *cartDomain.Cart) {
	if cs.cartCache != nil {
		id, err := cs.cartCache.BuildIdentifier(ctx, session)
		if err != nil {
			return
		}

		err = cs.cartCache.Delete(ctx, session, id)
		if err == ErrNoCacheEntry {
			cs.logger.WithContext(ctx).Info("Cart in cache not found: %v", err)
			return
		}
		if err != nil {
			cs.logger.WithContext(ctx).Error("Error while deleting cart in cache: %v", err)
		}
	}
}

// ReserveOrderIDAndSave reserves order id by using the PlaceOrder behaviour, sets and saves it on the cart.
// If the cart already holds a reserved order id no set/save is performed and the existing cart is returned.
// You may want to use this before proceeding with payment to ensure having a useful reference in the payment processing
func (cs *CartService) ReserveOrderIDAndSave(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	return cs.reserveOrderIDAndSave(ctx, session, false)
}

// ForceReserveOrderIDAndSave reserves order id by using the PlaceOrder behaviour, sets and saves it on the cart.
// Each call of this method reserves a new order ID, even if it is already set on the cart.
// You may want to use this before proceeding with payment to ensure having a useful reference in the payment processing
func (cs *CartService) ForceReserveOrderIDAndSave(ctx context.Context, session *web.Session) (*cartDomain.Cart, error) {
	return cs.reserveOrderIDAndSave(ctx, session, true)
}

func (cs *CartService) reserveOrderIDAndSave(ctx context.Context, session *web.Session, force bool) (*cartDomain.Cart, error) {
	if cs.placeOrderService == nil {
		return nil, errors.New("No placeOrderService registered")
	}
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	// the cart already has a reserved order id - no need to generate and update data, if force parameter is true
	// a new reserved order id will always be generated
	if !force && cart.AdditionalData.ReservedOrderID != "" {
		return cart, nil
	}

	reservedOrderID, err := cs.placeOrderService.ReserveOrderID(ctx, cart)
	if err != nil {
		cs.logger.WithContext(ctx).Debug("Reserve order id:", reservedOrderID)
		return nil, err
	}
	additionalData := cart.AdditionalData
	additionalData.ReservedOrderID = reservedOrderID
	cart, defers, err := behaviour.UpdateAdditionalData(ctx, cart, &additionalData)
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()
	return cart, err
}

// PlaceOrderWithCart converts the given cart with payments into orders by calling the Service
func (cs *CartService) PlaceOrderWithCart(ctx context.Context, session *web.Session, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	return cs.placeOrder(ctx, session, cart, payment)
}

// PlaceOrder converts the cart (possibly cached) with payments into orders by calling the Service
func (cs *CartService) PlaceOrder(ctx context.Context, session *web.Session, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	cart, _, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}
	return cs.placeOrder(ctx, session, cart, payment)
}

func (cs *CartService) placeOrder(ctx context.Context, session *web.Session, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	if cs.placeOrderService == nil {
		return nil, errors.New("No placeOrderService registered")
	}
	var placeOrderInfos placeorder.PlacedOrderInfos
	var errPlaceOrder error

	identity := cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx))

	if identity != nil {
		placeOrderInfos, errPlaceOrder = cs.placeOrderService.PlaceCustomerCart(ctx, identity, cart, payment)
	} else {
		placeOrderInfos, errPlaceOrder = cs.placeOrderService.PlaceGuestCart(ctx, cart, payment)
	}

	if errPlaceOrder != nil {
		cs.handleCartNotFound(session, errPlaceOrder)
		return nil, errPlaceOrder
	}

	cs.eventPublisher.PublishOrderPlacedEvent(ctx, cart, placeOrderInfos)
	_ = cs.DeleteSavedSessionGuestCartID(session)
	cs.DeleteCartInCache(ctx, session, cart)

	return placeOrderInfos, nil
}

// CancelOrder cancels a previously placed order and restores the cart content
func (cs *CartService) CancelOrder(ctx context.Context, session *web.Session, orderInfos placeorder.PlacedOrderInfos, cart cartDomain.Cart) (*cartDomain.Cart, error) {
	err := cs.cancelOrder(ctx, session, orderInfos)
	if err != nil {
		return nil, err
	}

	restoredCart, err := cs.RestoreCart(ctx, &cart)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			cs.logger.Error(fmt.Sprintf("couldn't restore cart err: %v", err))
		}
		return nil, err
	}

	return restoredCart, nil
}

func (cs *CartService) cancelOrder(ctx context.Context, session *web.Session, orderInfos placeorder.PlacedOrderInfos) error {
	if cs.placeOrderService == nil {
		return errors.New("No placeOrderService registered")
	}

	var err error

	identity := cs.webIdentityService.Identify(ctx, web.RequestFromContext(ctx))
	if identity != nil {
		err = cs.placeOrderService.CancelCustomerOrder(ctx, orderInfos, identity)
	} else {
		err = cs.placeOrderService.CancelGuestOrder(ctx, orderInfos)
	}

	if err != nil {
		err = fmt.Errorf("couldn't cancel order %q, err: %w", orderInfos, err)

		if !errors.Is(err, context.Canceled) {
			cs.logger.Error(err)
		}

		return err
	}
	return nil
}

// CancelOrderWithoutRestore cancels a previously placed order
func (cs *CartService) CancelOrderWithoutRestore(ctx context.Context, session *web.Session, orderInfos placeorder.PlacedOrderInfos) error {
	return cs.cancelOrder(ctx, session, orderInfos)
}

// GetDefaultDeliveryCode returns the configured default deliverycode
func (cs *CartService) GetDefaultDeliveryCode() string {
	return cs.defaultDeliveryCode
}

// handleEmptyDelivery - delete an empty delivery when found and feature flag is set
func (cs *CartService) handleEmptyDelivery(ctx context.Context, session *web.Session, cart *cartDomain.Cart, deliveryCode string) {
	if !cs.deleteEmptyDelivery {
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
		if !errors.Is(err, context.Canceled) {
			cs.logger.WithContext(ctx).WithField(flamingo.LogKeySubCategory, "handleEmptyDelivery").Error(err)
		}
		return
	}
}

func (cs *CartService) dispatchAllEvents(ctx context.Context, events []flamingo.Event) {
	for _, e := range events {
		cs.eventRouter.Dispatch(ctx, e)
	}
}

// AdjustItemsToRestrictedQty checks the quantity restrictions for each item of the cart and returns what quantities have been adjusted
func (cs *CartService) AdjustItemsToRestrictedQty(ctx context.Context, session *web.Session) (QtyAdjustmentResults, error) {
	qtyAdjustmentResults, err := cs.generateRestrictedQtyAdjustments(ctx, session)
	if err != nil {
		return nil, err
	}

	cart, _, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	for index, qtyAdjustmentResult := range qtyAdjustmentResults {
		couponCodes := cart.AppliedCouponCodes
		if qtyAdjustmentResult.NewQty < 1 {
			err = cs.DeleteItem(ctx, session, qtyAdjustmentResult.OriginalItem.ID, qtyAdjustmentResult.DeliveryCode)
		} else {
			err = cs.UpdateItemQty(ctx, session, qtyAdjustmentResult.OriginalItem.ID, qtyAdjustmentResult.DeliveryCode, qtyAdjustmentResult.NewQty)
		}
		if err != nil {
			return nil, err
		}
		cart, _, err = cs.cartReceiverService.GetCart(ctx, session)
		if err != nil {
			return nil, err
		}

		qtyAdjustmentResult.HasRemovedCouponCodes = !cartDomain.AppliedCouponCodes(couponCodes).ContainedIn(cart.AppliedCouponCodes)
		qtyAdjustmentResults[index] = qtyAdjustmentResult
	}

	return qtyAdjustmentResults, nil
}

// generateRestrictedQtyAdjustments checks the quantity restrictions for each item of the cart and returns which items should be adjusted and how
func (cs *CartService) generateRestrictedQtyAdjustments(ctx context.Context, session *web.Session) (QtyAdjustmentResults, error) {
	cart, _, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	result := make([]QtyAdjustmentResult, 0)
	for _, delivery := range cart.Deliveries {
		for _, item := range delivery.Cartitems {
			product, err := cs.productService.Get(ctx, item.MarketplaceCode)
			if err != nil {
				return nil, err
			}

			product, err = cs.getSpecificProductType(ctx, product, item.VariantMarketPlaceCode, item.BundleConfig)
			if err != nil {
				return nil, err
			}

			itemContext := ContextWithItemID(ctx, item.ID)
			restrictionResult := cs.restrictionService.RestrictQty(itemContext, session, product, cart, delivery.DeliveryInfo.Code)

			if restrictionResult.RemainingDifference >= 0 {
				continue
			}

			newQty := item.Qty + restrictionResult.RemainingDifference

			result = append(result, QtyAdjustmentResult{
				item,
				delivery.DeliveryInfo.Code,
				newQty < 1,
				restrictionResult,
				newQty,
				false,
			})
		}
	}

	return result, nil
}

func (cs *CartService) getSpecificProductType(_ context.Context, product productDomain.BasicProduct, variantMarketplaceCode string, bundleConfig productDomain.BundleConfiguration) (productDomain.BasicProduct, error) {
	var err error

	if product.Type() != productDomain.TypeConfigurable && product.Type() != productDomain.TypeBundle {
		return product, nil
	}

	if variantMarketplaceCode == "" {
		return product, nil
	}

	if configurableProduct, ok := product.(productDomain.ConfigurableProduct); ok {
		product, err = configurableProduct.GetConfigurableWithActiveVariant(variantMarketplaceCode)

		if err != nil {
			return nil, err
		}
	}

	if bundleProduct, ok := product.(productDomain.BundleProduct); ok {
		product, err = bundleProduct.GetBundleProductWithActiveChoices(bundleConfig)
		if err != nil {
			return nil, err
		}
	}

	return product, nil
}

// HasRemovedCouponCodes returns if a QtyAdjustmentResults has any adjustment that removed a coupon code from the cart
func (qar QtyAdjustmentResults) HasRemovedCouponCodes() bool {
	for _, qtyAdjustmentResult := range qar {
		if qtyAdjustmentResult.HasRemovedCouponCodes {
			return true
		}
	}

	return false
}

// UpdateAdditionalData of cart
func (cs *CartService) UpdateAdditionalData(ctx context.Context, session *web.Session, additionalData map[string]string) (*cartDomain.Cart, error) {
	cart, behaviour, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	cartAdditionalData := cart.AdditionalData
	if cartAdditionalData.CustomAttributes == nil {
		cartAdditionalData.CustomAttributes = map[string]string{}
	}

	for key, value := range additionalData {
		cartAdditionalData.CustomAttributes[key] = value
	}

	cart, defers, err := behaviour.UpdateAdditionalData(ctx, cart, &cartAdditionalData)
	defer func() {
		cs.updateCartInCacheIfCacheIsEnabled(ctx, session, cart)
		cs.dispatchAllEvents(ctx, defers)
	}()
	return cart, err
}

// UpdateDeliveryAdditionalData of cart
func (cs *CartService) UpdateDeliveryAdditionalData(ctx context.Context, session *web.Session, deliveryCode string, additionalData map[string]string) (*cartDomain.Cart, error) {
	cart, _, err := cs.cartReceiverService.GetCart(ctx, session)
	if err != nil {
		return nil, err
	}

	delivery, found := cart.GetDeliveryByCode(deliveryCode)
	if !found {
		return cart, nil
	}

	if delivery.DeliveryInfo.AdditionalData == nil {
		delivery.DeliveryInfo.AdditionalData = map[string]string{}
	}

	for key, value := range additionalData {
		delivery.DeliveryInfo.AdditionalData[key] = value
	}
	newDeliveryInfoUpdateCommand := cartDomain.CreateDeliveryInfoUpdateCommand(delivery.DeliveryInfo)

	err = cs.UpdateDeliveryInfo(ctx, session, deliveryCode, newDeliveryInfoUpdateCommand)
	if err != nil {
		return nil, err
	}

	cart, _, err = cs.cartReceiverService.GetCart(ctx, session)
	return cart, err
}
