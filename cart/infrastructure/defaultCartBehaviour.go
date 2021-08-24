package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// DefaultCartBehaviour defines the in memory cart order behaviour
	DefaultCartBehaviour struct {
		cartStorage             CartStorage
		productService          domain.ProductService
		logger                  flamingo.Logger
		itemBuilderProvider     cartDomain.ItemBuilderProvider
		deliveryBuilderProvider cartDomain.DeliveryBuilderProvider
		cartBuilderProvider     cartDomain.BuilderProvider
		giftCardHandler         GiftCardHandler
		voucherHandler          VoucherHandler
		defaultTaxRate          float64
	}

	// CartStorage Interface - might be implemented by other persistence types later as well
	CartStorage interface {
		GetCart(ctx context.Context, id string) (*cartDomain.Cart, error)
		HasCart(ctx context.Context, id string) bool
		StoreCart(ctx context.Context, cart *cartDomain.Cart) error
		RemoveCart(ctx context.Context, cart *cartDomain.Cart) error
	}

	// GiftCardHandler enables the projects to have specific GiftCard handling within the in-memory cart
	GiftCardHandler interface {
		ApplyGiftCard(ctx context.Context, cart *cartDomain.Cart, giftCardCode string) (*cartDomain.Cart, error)
		RemoveGiftCard(ctx context.Context, cart *cartDomain.Cart, giftCardCode string) (*cartDomain.Cart, error)
	}

	// VoucherHandler enables the projects to have specific Voucher handling within the in-memory cart
	VoucherHandler interface {
		ApplyVoucher(ctx context.Context, cart *cartDomain.Cart, couponCode string) (*cartDomain.Cart, error)
		RemoveVoucher(ctx context.Context, cart *cartDomain.Cart, couponCode string) (*cartDomain.Cart, error)
	}

	// DefaultGiftCardHandler implements a basic gift card handler
	DefaultGiftCardHandler struct{}

	// DefaultVoucherHandler implements a basic voucher handler
	DefaultVoucherHandler struct{}
)

var (
	_ cartDomain.ModifyBehaviour             = (*DefaultCartBehaviour)(nil)
	_ cartDomain.GiftCardAndVoucherBehaviour = (*DefaultCartBehaviour)(nil)
	_ cartDomain.CompleteBehaviour           = (*DefaultCartBehaviour)(nil)
	_ GiftCardHandler                        = (*DefaultGiftCardHandler)(nil)
	_ VoucherHandler                         = (*DefaultVoucherHandler)(nil)
)

// Inject dependencies
func (cob *DefaultCartBehaviour) Inject(
	CartStorage CartStorage,
	ProductService domain.ProductService,
	Logger flamingo.Logger,
	itemBuilderProvider cartDomain.ItemBuilderProvider,
	deliveryBuilderProvider cartDomain.DeliveryBuilderProvider,
	cartBuilderProvider cartDomain.BuilderProvider,
	voucherHandler VoucherHandler,
	giftCardHandler GiftCardHandler,
	config *struct {
		DefaultTaxRate float64 `inject:"config:commerce.cart.defaultCartAdapter.defaultTaxRate,optional"`
	},
) {
	cob.cartStorage = CartStorage
	cob.productService = ProductService
	cob.logger = Logger.WithField(flamingo.LogKeyCategory, "DefaultCartBehaviour")
	cob.itemBuilderProvider = itemBuilderProvider
	cob.deliveryBuilderProvider = deliveryBuilderProvider
	cob.cartBuilderProvider = cartBuilderProvider
	cob.voucherHandler = voucherHandler
	cob.giftCardHandler = giftCardHandler
	if config != nil {
		cob.defaultTaxRate = config.DefaultTaxRate
	}
}

// Complete a cart and remove from storage
func (cob *DefaultCartBehaviour) Complete(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	err := cob.cartStorage.RemoveCart(ctx, cart)
	if err != nil {
		return nil, nil, err
	}
	return cart, nil, nil
}

// Restore supplied cart (implements CompleteBehaviour)
func (cob *DefaultCartBehaviour) Restore(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, err
	}
	return &newCart, nil, nil
}

// DeleteItem removes an item from the cart
func (cob *DefaultCartBehaviour) DeleteItem(ctx context.Context, cart *cartDomain.Cart, itemID string, deliveryCode string) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("newCart.infrastructure.DefaultCartBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	if newDelivery, ok := newCart.GetDeliveryByCode(deliveryCode); ok {
		cob.logger.WithContext(ctx).Info("Inmemory Service Delete %v in %#v", itemID, newDelivery.Cartitems)
		for k, item := range newDelivery.Cartitems {
			if item.ID == itemID {
				if len(newDelivery.Cartitems) > k {
					newDelivery.Cartitems = append(newDelivery.Cartitems[:k], newDelivery.Cartitems[k+1:]...)
				} else {
					newDelivery.Cartitems = newDelivery.Cartitems[:k]
				}

				// update the delivery with the new info
				for j, delivery := range newCart.Deliveries {
					if deliveryCode == delivery.DeliveryInfo.Code {
						newCart.Deliveries[j] = *newDelivery
					}
				}
			}
		}
	}

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving newCart: %w", err)
	}
	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdateItem updates a cart item
func (cob *DefaultCartBehaviour) UpdateItem(ctx context.Context, cart *cartDomain.Cart, itemUpdateCommand cartDomain.ItemUpdateCommand) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("guestcart with id %q not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	err = cob.updateItem(ctx, &newCart, itemUpdateCommand)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdateItems updates multiple cart items
func (cob *DefaultCartBehaviour) UpdateItems(ctx context.Context, cart *cartDomain.Cart, itemUpdateCommands []cartDomain.ItemUpdateCommand) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot update - Guestcart with id %v not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	for _, itemUpdateCommand := range itemUpdateCommands {
		err := cob.updateItem(ctx, &newCart, itemUpdateCommand)
		if err != nil {
			return nil, nil, err
		}
	}

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

func (cob *DefaultCartBehaviour) updateItem(ctx context.Context, cart *cartDomain.Cart, itemUpdateCommand cartDomain.ItemUpdateCommand) error {
	itemBuilder := cob.itemBuilderProvider()
	itemDelivery, err := cart.GetDeliveryByItemID(itemUpdateCommand.ItemID)
	if err != nil {
		return err
	}

	cob.logger.WithContext(ctx).Info("Inmemory Service Update %v in %#v", itemUpdateCommand.ItemID, itemDelivery.Cartitems)
	for _, item := range itemDelivery.Cartitems {
		if itemUpdateCommand.ItemID == item.ID {
			itemBuilder.SetFromItem(item)
			if itemUpdateCommand.Qty != nil {
				itemBuilder.SetQty(*itemUpdateCommand.Qty)
			}

			if itemUpdateCommand.SourceID != nil {
				itemBuilder.SetSourceID(*itemUpdateCommand.SourceID)
			}
			itemBuilder.AddTaxInfo("default", big.NewFloat(cob.defaultTaxRate), nil).CalculatePricesAndTax()
			newItem, err := itemBuilder.Build()
			if err != nil {
				return err
			}
			for k, currentItem := range itemDelivery.Cartitems {
				if currentItem.ID == itemUpdateCommand.ItemID {
					itemDelivery.Cartitems[k] = *newItem
				}
			}
		}
	}

	// update the delivery with the new info
	for j, delivery := range cart.Deliveries {
		if itemDelivery.DeliveryInfo.Code == delivery.DeliveryInfo.Code {
			cart.Deliveries[j] = delivery
		}
	}

	return nil
}

// AddToCart add an item to the cart
func (cob *DefaultCartBehaviour) AddToCart(ctx context.Context, cart *cartDomain.Cart, deliveryCode string, addRequest cartDomain.AddRequest) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	if cart != nil && !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot add - Guestcart with id %v not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	// create delivery if it does not yet exist
	if !newCart.HasDeliveryForCode(deliveryCode) {
		// create delivery and add item
		delivery := new(cartDomain.Delivery)
		delivery.DeliveryInfo.Code = deliveryCode
		newCart.Deliveries = append(newCart.Deliveries, *delivery)
	}

	// has cart current delivery, check if there is an item present for this delivery
	delivery, _ := newCart.GetDeliveryByCode(deliveryCode)

	// create and add new item
	cartItem, err := cob.buildItemForCart(ctx, addRequest)
	if err != nil {
		return nil, nil, err
	}

	// does the item already exist?
	itemFound := false

	for i, item := range delivery.Cartitems {
		if item.MarketplaceCode == addRequest.MarketplaceCode {
			delivery.Cartitems[i] = *cartItem
			itemFound = true
		}
	}

	if !itemFound {
		delivery.Cartitems = append(delivery.Cartitems, *cartItem)
	}

	for k, del := range newCart.Deliveries {
		if del.DeliveryInfo.Code == delivery.DeliveryInfo.Code {
			newCart.Deliveries[k] = *delivery
		}
	}

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

func (cob *DefaultCartBehaviour) buildItemForCart(ctx context.Context, addRequest cartDomain.AddRequest) (*cartDomain.Item, error) {
	itemBuilder := cob.itemBuilderProvider()

	// create and add new item
	product, err := cob.productService.Get(ctx, addRequest.MarketplaceCode)
	if err != nil {
		return nil, err
	}

	// Get variant of configurable product
	if configurableProduct, ok := product.(domain.ConfigurableProduct); ok && addRequest.VariantMarketplaceCode != "" {
		productWithActiveVariant, err := configurableProduct.GetConfigurableWithActiveVariant(addRequest.VariantMarketplaceCode)
		if err != nil {
			return nil, err
		}
		product = productWithActiveVariant
	}

	itemBuilder.
		SetQty(addRequest.Qty).
		AddTaxInfo("default", big.NewFloat(cob.defaultTaxRate), nil).
		SetByProduct(product).
		SetID(strconv.Itoa(rand.Int())).
		SetExternalReference(strconv.Itoa(rand.Int())).
		SetAdditionalData(addRequest.AdditionalData)

	return itemBuilder.Build()
}

// CleanCart removes everything from the cart, e.g. deliveries, billing address, etc
func (cob *DefaultCartBehaviour) CleanCart(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCart.Deliveries = []cartDomain.Delivery{}
	newCart.AppliedCouponCodes = nil
	newCart.AppliedGiftCards = nil
	newCart.PaymentSelection = nil
	newCart.AdditionalData = cartDomain.AdditionalData{}
	newCart.Purchaser = nil
	newCart.BillingAddress = nil
	newCart.Totalitems = nil

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// CleanDelivery removes a complete delivery with its items from the cart
func (cob *DefaultCartBehaviour) CleanDelivery(ctx context.Context, cart *cartDomain.Cart, deliveryCode string) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	// create delivery if it does not yet exist
	if !newCart.HasDeliveryForCode(deliveryCode) {
		return nil, nil, fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: delivery %s not found", deliveryCode)
	}

	var position int
	for i, delivery := range newCart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			position = i
			break
		}
	}

	newLength := len(newCart.Deliveries) - 1
	newCart.Deliveries[position] = newCart.Deliveries[newLength]
	newCart.Deliveries[newLength] = cartDomain.Delivery{}
	newCart.Deliveries = newCart.Deliveries[:newLength]

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdatePurchaser @todo implement when needed
func (cob *DefaultCartBehaviour) UpdatePurchaser(ctx context.Context, cart *cartDomain.Cart, purchaser *cartDomain.Person, additionalData *cartDomain.AdditionalData) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCart.Purchaser = purchaser

	if additionalData != nil {
		newCart.AdditionalData.CustomAttributes = additionalData.CustomAttributes
	}

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// UpdateBillingAddress - updates address
func (cob *DefaultCartBehaviour) UpdateBillingAddress(ctx context.Context, cart *cartDomain.Cart, billingAddress cartDomain.Address) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCart.BillingAddress = &billingAddress

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// UpdateAdditionalData updates additional data
func (cob *DefaultCartBehaviour) UpdateAdditionalData(ctx context.Context, cart *cartDomain.Cart, additionalData *cartDomain.AdditionalData) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCart.AdditionalData = *additionalData
	err = cob.cartStorage.StoreCart(ctx, &newCart)

	if err != nil {
		return nil, nil, fmt.Errorf("error on updating additional data: %w", err)
	}

	return &newCart, nil, nil
}

// UpdatePaymentSelection updates payment on cart
func (cob *DefaultCartBehaviour) UpdatePaymentSelection(ctx context.Context, cart *cartDomain.Cart, paymentSelection cartDomain.PaymentSelection) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	if paymentSelection != nil {
		err := cob.checkPaymentSelection(ctx, &newCart, paymentSelection)
		if err != nil {
			return nil, nil, err
		}
	}
	newCart.PaymentSelection = paymentSelection

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// UpdateDeliveryInfo updates a delivery info
func (cob *DefaultCartBehaviour) UpdateDeliveryInfo(ctx context.Context, cart *cartDomain.Cart, deliveryCode string, deliveryInfoUpdateCommand cartDomain.DeliveryInfoUpdateCommand) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	deliveryInfo := deliveryInfoUpdateCommand.DeliveryInfo
	deliveryInfo.AdditionalDeliveryInfos = deliveryInfoUpdateCommand.Additional()

	for key, delivery := range newCart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			newCart.Deliveries[key].DeliveryInfo = deliveryInfo
			err := cob.cartStorage.StoreCart(ctx, &newCart)
			if err != nil {
				return nil, nil, fmt.Errorf("error on saving cart: %w", err)
			}
			return &newCart, nil, nil
		}
	}
	newCart.Deliveries = append(newCart.Deliveries, cartDomain.Delivery{DeliveryInfo: deliveryInfo})

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("error on saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdateDeliveryInfoAdditionalData @todo implement when needed
func (cob *DefaultCartBehaviour) UpdateDeliveryInfoAdditionalData(_ context.Context, cart *cartDomain.Cart, _ string, _ *cartDomain.AdditionalData) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	return cart, nil, nil
}

// GetCart returns the current cart from storage
func (cob *DefaultCartBehaviour) GetCart(ctx context.Context, cartID string) (*cartDomain.Cart, error) {
	if !cob.cartStorage.HasCart(ctx, cartID) {
		cob.logger.Info(fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot get - cart with id %v not existent", cartID))
		return nil, cartDomain.ErrCartNotFound
	}

	cart, err := cob.cartStorage.GetCart(ctx, cartID)
	if err != nil {
		cob.logger.Info(fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: get cart from storage: %w ", err))
		return nil, cartDomain.ErrCartNotFound
	}

	newCart, err := cart.Clone()
	if err != nil {
		cob.logger.Info(fmt.Errorf("cart clone failed: %w ", err))
		return nil, cartDomain.ErrCartNotFound
	}

	return &newCart, nil
}

// StoreNewCart created and stores a new cart.
func (cob *DefaultCartBehaviour) StoreNewCart(ctx context.Context, newCart *cartDomain.Cart) (*cartDomain.Cart, error) {
	if newCart.ID == "" {
		return nil, errors.New("no id given")
	}
	return newCart, cob.cartStorage.StoreCart(ctx, newCart)
}

// ApplyVoucher applies a voucher to the cart
func (cob *DefaultCartBehaviour) ApplyVoucher(ctx context.Context, cart *cartDomain.Cart, couponCode string) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCartWithVoucher, err := cob.voucherHandler.ApplyVoucher(ctx, &newCart, couponCode)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(ctx, newCartWithVoucher)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, newCartWithVoucher)
}

// ApplyAny applies a voucher or gift card to the cart
func (cob *DefaultCartBehaviour) ApplyAny(ctx context.Context, cart *cartDomain.Cart, anyCode string) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	currentCart, deferFunc, err := cob.ApplyVoucher(ctx, cart, anyCode)
	if err == nil {
		// successfully applied as voucher
		return currentCart, deferFunc, nil
	}

	// some error occurred, retry as gift card
	return cob.ApplyGiftCard(ctx, cart, anyCode)
}

// RemoveVoucher removes a voucher from the cart
func (cob *DefaultCartBehaviour) RemoveVoucher(ctx context.Context, cart *cartDomain.Cart, couponCode string) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCartWithoutVoucher, err := cob.voucherHandler.RemoveVoucher(ctx, &newCart, couponCode)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(ctx, newCartWithoutVoucher)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, newCartWithoutVoucher)
}

// ApplyGiftCard applies a gift card to the cart
// if a GiftCard is applied, it will be added to the array AppliedGiftCards on the cart
func (cob *DefaultCartBehaviour) ApplyGiftCard(ctx context.Context, cart *cartDomain.Cart, giftCardCode string) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCartWithGiftCard, err := cob.giftCardHandler.ApplyGiftCard(ctx, &newCart, giftCardCode)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(ctx, newCartWithGiftCard)
	if err != nil {
		return nil, nil, err
	}
	return cob.resetPaymentSelectionIfInvalid(ctx, newCartWithGiftCard)
}

// RemoveGiftCard removes a gift card from the cart
// if a GiftCard is removed, it will be removed from the array AppliedGiftCards on the cart
func (cob *DefaultCartBehaviour) RemoveGiftCard(ctx context.Context, cart *cartDomain.Cart, giftCardCode string) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCartWithOutGiftCard, err := cob.giftCardHandler.RemoveGiftCard(ctx, &newCart, giftCardCode)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(ctx, newCartWithOutGiftCard)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, newCartWithOutGiftCard)
}

func (cob *DefaultCartBehaviour) isCurrentPaymentSelectionValid(ctx context.Context, cart *cartDomain.Cart) bool {
	return cob.checkPaymentSelection(ctx, cart, cart.PaymentSelection) == nil
}

// isPaymentSelectionValid checks if the grand total of the cart matches the total of the supplied payment selection
func (cob *DefaultCartBehaviour) checkPaymentSelection(_ context.Context, cart *cartDomain.Cart, paymentSelection cartDomain.PaymentSelection) error {
	if paymentSelection == nil {
		return nil
	}
	paymentSelectionTotal := paymentSelection.TotalValue()

	if !cart.GrandTotal().LikelyEqual(paymentSelectionTotal) {
		return errors.New("payment total does not match with grand total")
	}
	return nil
}

// resetPaymentSelectionIfInvalid checks for valid payment selection on given cart and deletes it in case it is invalid
func (cob *DefaultCartBehaviour) resetPaymentSelectionIfInvalid(ctx context.Context, cart *cartDomain.Cart) (*cartDomain.Cart, cartDomain.DeferEvents, error) {
	if cart.PaymentSelection == nil {
		return cart, nil, nil
	}
	err := cob.checkPaymentSelection(ctx, cart, cart.PaymentSelection)
	if err != nil {
		cart, defers, err := cob.UpdatePaymentSelection(ctx, cart, nil)
		defers = append(defers, &events.PaymentSelectionHasBeenResetEvent{Cart: cart})
		return cart, defers, err
	}

	return cart, nil, nil
}

// ApplyVoucher checks the voucher and adds the voucher to the supplied cart if valid
func (DefaultVoucherHandler) ApplyVoucher(_ context.Context, cart *cartDomain.Cart, couponCode string) (*cartDomain.Cart, error) {
	if couponCode != "valid_voucher" && couponCode != "valid" {
		return nil, errors.New("Code invalid")
	}

	coupon := cartDomain.CouponCode{
		Code: couponCode,
	}

	cart.AppliedCouponCodes = append(cart.AppliedCouponCodes, coupon)
	return cart, nil
}

// RemoveVoucher removes the voucher from the cart if possible
func (DefaultVoucherHandler) RemoveVoucher(_ context.Context, cart *cartDomain.Cart, couponCode string) (*cartDomain.Cart, error) {
	for i, coupon := range cart.AppliedCouponCodes {
		if coupon.Code == couponCode {
			cart.AppliedCouponCodes[i] = cart.AppliedCouponCodes[len(cart.AppliedCouponCodes)-1]
			cart.AppliedCouponCodes[len(cart.AppliedCouponCodes)-1] = cartDomain.CouponCode{}
			cart.AppliedCouponCodes = cart.AppliedCouponCodes[:len(cart.AppliedCouponCodes)-1]
			return cart, nil
		}
	}

	return cart, nil
}

// ApplyGiftCard checks the gift card and adds it to the supplied cart if valid
func (DefaultGiftCardHandler) ApplyGiftCard(_ context.Context, cart *cartDomain.Cart, giftCardCode string) (*cartDomain.Cart, error) {
	if giftCardCode != "valid_giftcard" && giftCardCode != "valid" {
		return nil, errors.New("Code invalid")
	}

	giftCard := cartDomain.AppliedGiftCard{
		Code:      giftCardCode,
		Applied:   priceDomain.NewFromInt(10, 100, "$"),
		Remaining: priceDomain.NewFromInt(0, 100, "$"),
	}
	cart.AppliedGiftCards = append(cart.AppliedGiftCards, giftCard)

	return cart, nil
}

// RemoveGiftCard removes the gift card from the cart if possible
func (DefaultGiftCardHandler) RemoveGiftCard(_ context.Context, cart *cartDomain.Cart, giftCardCode string) (*cartDomain.Cart, error) {
	for i, giftCard := range cart.AppliedGiftCards {
		if giftCard.Code == giftCardCode {
			cart.AppliedGiftCards[i] = cart.AppliedGiftCards[len(cart.AppliedGiftCards)-1]
			cart.AppliedGiftCards[len(cart.AppliedGiftCards)-1] = cartDomain.AppliedGiftCard{}
			cart.AppliedGiftCards = cart.AppliedGiftCards[:len(cart.AppliedGiftCards)-1]
			return cart, nil
		}
	}

	return cart, nil
}
