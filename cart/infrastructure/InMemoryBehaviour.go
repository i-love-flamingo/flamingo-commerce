package infrastructure

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	domaincart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/pkg/errors"
)

type (
	// InMemoryBehaviour defines the in memory cart order behaviour
	InMemoryBehaviour struct {
		cartStorage             CartStorage
		productService          domain.ProductService
		logger                  flamingo.Logger
		itemBuilderProvider     domaincart.ItemBuilderProvider
		deliveryBuilderProvider domaincart.DeliveryBuilderProvider
		cartBuilderProvider     domaincart.BuilderProvider
		giftCardHandler         GiftCardHandler
		voucherHandler          VoucherHandler
		defaultTaxRate          float64
	}

	// CartStorage Interface - might be implemented by other persistence types later as well
	CartStorage interface {
		GetCart(id string) (*domaincart.Cart, error)
		HasCart(id string) bool
		StoreCart(cart *domaincart.Cart) error
		RemoveCart(cart *domaincart.Cart) error
	}

	// GiftCardHandler enables the projects to have specific GiftCard handling within the in-memory cart
	GiftCardHandler interface {
		ApplyGiftCard(ctx context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, error)
		RemoveGiftCard(ctx context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, error)
	}

	// VoucherHandler enables the projects to have specific Voucher handling within the in-memory cart
	VoucherHandler interface {
		ApplyVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, error)
		RemoveVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, error)
	}

	// DefaultGiftCardHandler implements a basic gift card handler
	DefaultGiftCardHandler struct{}

	// DefaultVoucherHandler implements a basic voucher handler
	DefaultVoucherHandler struct{}
)

var (
	_ domaincart.ModifyBehaviour             = (*InMemoryBehaviour)(nil)
	_ domaincart.GiftCardAndVoucherBehaviour = (*InMemoryBehaviour)(nil)
	_ domaincart.CompleteBehaviour           = (*InMemoryBehaviour)(nil)
	_ GiftCardHandler                        = (*DefaultGiftCardHandler)(nil)
	_ VoucherHandler                         = (*DefaultVoucherHandler)(nil)
)

// Inject dependencies
func (cob *InMemoryBehaviour) Inject(
	CartStorage CartStorage,
	ProductService domain.ProductService,
	Logger flamingo.Logger,
	itemBuilderProvider domaincart.ItemBuilderProvider,
	deliveryBuilderProvider domaincart.DeliveryBuilderProvider,
	cartBuilderProvider domaincart.BuilderProvider,
	voucherHandler VoucherHandler,
	giftCardHandler GiftCardHandler,
	config *struct {
		DefaultTaxRate float64 `inject:"config:commerce.cart.inMemoryCartServiceAdapter.defaultTaxRate,optional"`
	},
) {
	cob.cartStorage = CartStorage
	cob.productService = ProductService
	cob.logger = Logger.WithField(flamingo.LogKeyCategory, "inmemorybehaviour")
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
func (cob *InMemoryBehaviour) Complete(_ context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	err := cob.cartStorage.RemoveCart(cart)
	if err != nil {
		return nil, nil, err
	}
	return cart, nil, nil
}

// Restore supplied cart
func (cob *InMemoryBehaviour) Restore(_ context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart := cart
	err := cob.cartStorage.StoreCart(newCart)
	if err != nil {
		return nil, nil, err
	}
	return newCart, nil, nil
}

// DeleteItem removes an item from the cart
func (cob *InMemoryBehaviour) DeleteItem(ctx context.Context, cart *domaincart.Cart, itemID string, deliveryCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	if newDelivery, ok := cart.GetDeliveryByCode(deliveryCode); ok {
		cob.logger.WithContext(ctx).Info("Inmemory Service Delete %v in %#v", itemID, newDelivery.Cartitems)
		for k, item := range newDelivery.Cartitems {
			if item.ID == itemID {
				if len(newDelivery.Cartitems) > k {
					newDelivery.Cartitems = append(newDelivery.Cartitems[:k], newDelivery.Cartitems[k+1:]...)
				} else {
					newDelivery.Cartitems = newDelivery.Cartitems[:k]
				}

				// update the delivery with the new info
				for j, delivery := range cart.Deliveries {
					if deliveryCode == delivery.DeliveryInfo.Code {
						cart.Deliveries[j] = *newDelivery
					}
				}

			}
		}
	}

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}
	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

// UpdateItem updates a cart item
func (cob *InMemoryBehaviour) UpdateItem(ctx context.Context, cart *domaincart.Cart, itemUpdateCommand domaincart.ItemUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot add - Guestcart with id %v not existent", cart.ID)
	}

	err := cob.updateItem(ctx, cart, itemUpdateCommand)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

// UpdateItems updates multiple cart items
func (cob *InMemoryBehaviour) UpdateItems(ctx context.Context, cart *domaincart.Cart, itemUpdateCommands []domaincart.ItemUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot update - Guestcart with id %v not existent", cart.ID)
	}

	for _, itemUpdateCommand := range itemUpdateCommands {
		err := cob.updateItem(ctx, cart, itemUpdateCommand)
		if err != nil {
			return nil, nil, err
		}
	}

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

func (cob *InMemoryBehaviour) updateItem(ctx context.Context, cart *domaincart.Cart, itemUpdateCommand domaincart.ItemUpdateCommand) error {
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
func (cob *InMemoryBehaviour) AddToCart(ctx context.Context, cart *domaincart.Cart, deliveryCode string, addRequest domaincart.AddRequest) (*domaincart.Cart, domaincart.DeferEvents, error) {

	if cart != nil && !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot add - Guestcart with id %v not existent", cart.ID)
	}

	// create delivery if it does not yet exist
	if !cart.HasDeliveryForCode(deliveryCode) {
		// create delivery and add item
		delivery := new(domaincart.Delivery)
		delivery.DeliveryInfo.Code = deliveryCode
		cart.Deliveries = append(cart.Deliveries, *delivery)
	}

	// has cart current delivery, check if there is an item present for this delivery
	delivery, _ := cart.GetDeliveryByCode(deliveryCode)

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

	for k, del := range cart.Deliveries {
		if del.DeliveryInfo.Code == delivery.DeliveryInfo.Code {
			cart.Deliveries[k] = *delivery
		}
	}

	err = cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

func (cob *InMemoryBehaviour) buildItemForCart(ctx context.Context, addRequest domaincart.AddRequest) (*domaincart.Item, error) {
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

// CleanCart removes all deliveries and their items from the cart
func (cob *InMemoryBehaviour) CleanCart(_ context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	cart.Deliveries = []domaincart.Delivery{}

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cart, nil, nil
}

// CleanDelivery removes a complete delivery with its items from the cart
func (cob *InMemoryBehaviour) CleanDelivery(ctx context.Context, cart *domaincart.Cart, deliveryCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	// create delivery if it does not yet exist
	if !cart.HasDeliveryForCode(deliveryCode) {
		return nil, nil, errors.Errorf("cart.infrastructure.InMemoryBehaviour: delivery %s not found", deliveryCode)
	}

	var position int
	for i, delivery := range cart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			position = i
			break
		}
	}

	newLength := len(cart.Deliveries) - 1
	cart.Deliveries[position] = cart.Deliveries[newLength]
	cart.Deliveries[newLength] = domaincart.Delivery{}
	cart.Deliveries = cart.Deliveries[:newLength]

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

// UpdatePurchaser sets the purchaser data and the additional data on the cart
func (cob *InMemoryBehaviour) UpdatePurchaser(_ context.Context, cart *domaincart.Cart, purchaser *domaincart.Person, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	cart.Purchaser = purchaser

	if additionalData != nil {
		cart.AdditionalData.CustomAttributes = additionalData.CustomAttributes
	}

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cart, nil, nil
}

// UpdateBillingAddress updates the billing address
func (cob *InMemoryBehaviour) UpdateBillingAddress(_ context.Context, cart *domaincart.Cart, billingAddress domaincart.Address) (*domaincart.Cart, domaincart.DeferEvents, error) {

	cart.BillingAddress = &billingAddress

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cart, nil, nil
}

// UpdateAdditionalData updates additional data
func (cob *InMemoryBehaviour) UpdateAdditionalData(_ context.Context, cart *domaincart.Cart, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	cart.AdditionalData = *additionalData
	err := cob.cartStorage.StoreCart(cart)

	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on updating additional data")
	}

	return cart, nil, nil
}

//UpdatePaymentSelection updates payment on cart
func (cob *InMemoryBehaviour) UpdatePaymentSelection(ctx context.Context, cart *domaincart.Cart, paymentSelection domaincart.PaymentSelection) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if paymentSelection != nil {
		err := cob.checkPaymentSelection(ctx, cart, paymentSelection)
		if err != nil {
			return nil, nil, err
		}
	}
	cart.PaymentSelection = paymentSelection

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cart, nil, nil
}

// UpdateDeliveryInfo updates a delivery info
func (cob *InMemoryBehaviour) UpdateDeliveryInfo(ctx context.Context, cart *domaincart.Cart, deliveryCode string, deliveryInfoUpdateCommand domaincart.DeliveryInfoUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {

	deliveryInfo := deliveryInfoUpdateCommand.DeliveryInfo
	deliveryInfo.AdditionalDeliveryInfos = deliveryInfoUpdateCommand.Additional()

	for key, delivery := range cart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			cart.Deliveries[key].DeliveryInfo = deliveryInfo
			err := cob.cartStorage.StoreCart(cart)
			if err != nil {
				return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
			}
			return cart, nil, nil
		}
	}
	cart.Deliveries = append(cart.Deliveries, domaincart.Delivery{DeliveryInfo: deliveryInfo})

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

// UpdateDeliveryInfoAdditionalData @todo implement when needed
func (cob *InMemoryBehaviour) UpdateDeliveryInfoAdditionalData(_ context.Context, cart *domaincart.Cart, _ string, _ *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	return cart, nil, nil
}

// GetCart returns the current cart from storage
func (cob *InMemoryBehaviour) GetCart(_ context.Context, cartID string) (*domaincart.Cart, error) {
	if cob.cartStorage.HasCart(cartID) {
		// if cart exists, there is no error ;)
		cart, err := cob.cartStorage.GetCart(cartID)
		if err == nil {
			return cart, nil
		}
	}

	return nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot get - Guestcart with id %v not existent", cartID)
}

// storeCart in the memory
func (cob *InMemoryBehaviour) storeCart(cart *domaincart.Cart) error {
	return cob.cartStorage.StoreCart(cart)
}

// ApplyVoucher applies a voucher to the cart
func (cob *InMemoryBehaviour) ApplyVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	cart, err := cob.voucherHandler.ApplyVoucher(ctx, cart, couponCode)

	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

// ApplyAny applies a voucher or giftcard to the cart
func (cob *InMemoryBehaviour) ApplyAny(ctx context.Context, cart *domaincart.Cart, anyCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	currentCart, deferFunc, err := cob.ApplyVoucher(ctx, cart, anyCode)
	if err == nil {
		// successfully applied as voucher
		return currentCart, deferFunc, nil
	}

	// some error occurred, retry as giftcard
	return cob.ApplyGiftCard(ctx, cart, anyCode)
}

// RemoveVoucher removes a voucher from the cart
func (cob *InMemoryBehaviour) RemoveVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	cart, err := cob.voucherHandler.RemoveVoucher(ctx, cart, couponCode)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

// ApplyGiftCard applies a gift card to the cart
// if a GiftCard is applied, it will be added to the array AppliedGiftCards on the cart
func (cob *InMemoryBehaviour) ApplyGiftCard(ctx context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	cart, err := cob.giftCardHandler.ApplyGiftCard(ctx, cart, giftCardCode)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, err
	}
	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

// RemoveGiftCard removes a gift card from the cart
// if a GiftCard is removed, it will be removed from the array AppliedGiftCards on the cart
func (cob *InMemoryBehaviour) RemoveGiftCard(ctx context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	cart, err := cob.giftCardHandler.RemoveGiftCard(ctx, cart, giftCardCode)
	if err != nil {
		return nil, nil, err
	}

	err = cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

func (cob *InMemoryBehaviour) isCurrentPaymentSelectionValid(ctx context.Context, cart *domaincart.Cart) bool {
	return cob.checkPaymentSelection(ctx, cart, cart.PaymentSelection) == nil
}

// isPaymentSelectionValid checks if the grand total of the cart matches the total of the supplied payment selection
func (cob *InMemoryBehaviour) checkPaymentSelection(_ context.Context, cart *domaincart.Cart, paymentSelection domaincart.PaymentSelection) error {
	if paymentSelection == nil {
		return nil
	}
	paymentSelectionTotal := paymentSelection.TotalValue()

	if !cart.GrandTotal().LikelyEqual(paymentSelectionTotal) {
		return errors.New("Payment Total does not match with Grandtotal")
	}
	return nil
}

// resetPaymentSelectionIfInvalid checks for valid paymentselection on givencart and deletes in in case it is invalid
func (cob *InMemoryBehaviour) resetPaymentSelectionIfInvalid(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
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
func (DefaultVoucherHandler) ApplyVoucher(_ context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, error) {
	if couponCode != "valid_voucher" && couponCode != "valid" {
		return nil, errors.New("Code invalid")
	}

	coupon := domaincart.CouponCode{
		Code: couponCode,
	}

	cart.AppliedCouponCodes = append(cart.AppliedCouponCodes, coupon)
	return cart, nil
}

// RemoveVoucher removes the voucher from the cart if possible
func (DefaultVoucherHandler) RemoveVoucher(_ context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, error) {
	for i, coupon := range cart.AppliedCouponCodes {
		if coupon.Code == couponCode {
			cart.AppliedCouponCodes[i] = cart.AppliedCouponCodes[len(cart.AppliedCouponCodes)-1]
			cart.AppliedCouponCodes[len(cart.AppliedCouponCodes)-1] = domaincart.CouponCode{}
			cart.AppliedCouponCodes = cart.AppliedCouponCodes[:len(cart.AppliedCouponCodes)-1]
			return cart, nil
		}
	}

	return cart, nil
}

// ApplyGiftCard checks the gift card and adds it to the supplied cart if valid
func (DefaultGiftCardHandler) ApplyGiftCard(_ context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, error) {
	if giftCardCode != "valid_giftcard" && giftCardCode != "valid" {
		return nil, errors.New("Code invalid")
	}

	giftCard := domaincart.AppliedGiftCard{
		Code:      giftCardCode,
		Applied:   priceDomain.NewFromInt(10, 100, "$"),
		Remaining: priceDomain.NewFromInt(0, 100, "$"),
	}
	cart.AppliedGiftCards = append(cart.AppliedGiftCards, giftCard)

	return cart, nil
}

// RemoveGiftCard removes the gift card from the cart if possible
func (DefaultGiftCardHandler) RemoveGiftCard(_ context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, error) {
	for i, giftcard := range cart.AppliedGiftCards {
		if giftcard.Code == giftCardCode {
			cart.AppliedGiftCards[i] = cart.AppliedGiftCards[len(cart.AppliedGiftCards)-1]
			cart.AppliedGiftCards[len(cart.AppliedGiftCards)-1] = domaincart.AppliedGiftCard{}
			cart.AppliedGiftCards = cart.AppliedGiftCards[:len(cart.AppliedGiftCards)-1]
			return cart, nil
		}
	}

	return cart, nil
}
