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
	// DefaultCartBehaviour defines the in memory cart order behaviour
	DefaultCartBehaviour struct {
		cartStorage     CartStorage
		productService  domain.ProductService
		logger          flamingo.Logger
		giftCardHandler GiftCardHandler
		voucherHandler  VoucherHandler
		defaultTaxRate  float64
		grossPricing    bool
		defaultCurrency string
	}

	// CartStorage Interface - might be implemented by other persistence types later as well
	CartStorage interface {
		GetCart(ctx context.Context, id string) (*domaincart.Cart, error)
		HasCart(ctx context.Context, id string) bool
		StoreCart(ctx context.Context, cart *domaincart.Cart) error
		RemoveCart(ctx context.Context, cart *domaincart.Cart) error
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
	_ domaincart.ModifyBehaviour             = (*DefaultCartBehaviour)(nil)
	_ domaincart.GiftCardAndVoucherBehaviour = (*DefaultCartBehaviour)(nil)
	_ domaincart.CompleteBehaviour           = (*DefaultCartBehaviour)(nil)
	_ GiftCardHandler                        = (*DefaultGiftCardHandler)(nil)
	_ VoucherHandler                         = (*DefaultVoucherHandler)(nil)
)

// Inject dependencies
func (cob *DefaultCartBehaviour) Inject(
	CartStorage CartStorage,
	ProductService domain.ProductService,
	Logger flamingo.Logger,
	voucherHandler VoucherHandler,
	giftCardHandler GiftCardHandler,
	config *struct {
		DefaultTaxRate  float64 `inject:"config:commerce.cart.defaultCartAdapter.defaultTaxRate,optional"`
		ProductPricing  string  `inject:"config:commerce.cart.defaultCartAdapter.productPrices"`
		DefaultCurrency string  `inject:"config:commerce.cart.defaultCartAdapter.defaultCurrency"`
	},
) {
	cob.cartStorage = CartStorage
	cob.productService = ProductService
	cob.logger = Logger.WithField(flamingo.LogKeyCategory, "inmemorybehaviour")
	cob.voucherHandler = voucherHandler
	cob.giftCardHandler = giftCardHandler
	if config != nil {
		cob.defaultTaxRate = config.DefaultTaxRate
		cob.defaultCurrency = config.DefaultCurrency
		if config.ProductPricing == "gross" {
			cob.grossPricing = true
		}
	}
}

// Complete a cart and remove from storage
func (cob *DefaultCartBehaviour) Complete(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	err := cob.cartStorage.RemoveCart(ctx, cart)
	if err != nil {
		return nil, nil, err
	}
	return cart, nil, nil
}

// Restore supplied cart (implements CompleteBehaviour)
func (cob *DefaultCartBehaviour) Restore(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}
	cob.collectTotals(cart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, err
	}
	return &newCart, nil, nil
}

// DeleteItem removes an item from the cart
func (cob *DefaultCartBehaviour) DeleteItem(ctx context.Context, cart *domaincart.Cart, itemID string, deliveryCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
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

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "newCart.infrastructure.DefaultCartBehaviour: error on saving newCart")
	}
	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdateItem updates a cart item
func (cob *DefaultCartBehaviour) UpdateItem(ctx context.Context, cart *domaincart.Cart, itemUpdateCommand domaincart.ItemUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot add - Guestcart with id %v not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	err = cob.updateItem(ctx, &newCart, itemUpdateCommand)
	if err != nil {
		return nil, nil, err
	}

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdateItems updates multiple cart items
func (cob *DefaultCartBehaviour) UpdateItems(ctx context.Context, cart *domaincart.Cart, itemUpdateCommands []domaincart.ItemUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
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

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

func (cob *DefaultCartBehaviour) updateItem(ctx context.Context, cart *domaincart.Cart, itemUpdateCommand domaincart.ItemUpdateCommand) error {
	itemDelivery, err := cart.GetDeliveryByItemID(itemUpdateCommand.ItemID)
	if err != nil {
		return err
	}

	cob.logger.WithContext(ctx).Info("Inmemory Service Update %v in %#v", itemUpdateCommand.ItemID, itemDelivery.Cartitems)
	for k, item := range itemDelivery.Cartitems {
		if itemUpdateCommand.ItemID == item.ID {
			if itemUpdateCommand.Qty != nil {
				itemDelivery.Cartitems[k].Qty = *itemUpdateCommand.Qty

				gross := item.SinglePriceGross.Clone().Amount().Mul(item.SinglePriceGross.Amount(), big.NewFloat(float64(*itemUpdateCommand.Qty)))
				itemDelivery.Cartitems[k].RowPriceGross = priceDomain.NewFromBigFloat(*gross, item.SinglePriceGross.Currency())

				net := item.SinglePriceNet.Clone().Amount().Mul(item.SinglePriceNet.Amount(), big.NewFloat(float64(*itemUpdateCommand.Qty)))
				itemDelivery.Cartitems[k].RowPriceNet = priceDomain.NewFromBigFloat(*net, item.SinglePriceNet.Currency())

				itemDelivery.Cartitems[k].RowPriceGrossWithDiscount = itemDelivery.Cartitems[k].RowPriceGross
				itemDelivery.Cartitems[k].RowPriceNetWithDiscount = itemDelivery.Cartitems[k].RowPriceNet

				itemDelivery.Cartitems[k].RowPriceGrossWithItemRelatedDiscount = itemDelivery.Cartitems[k].RowPriceGross
				itemDelivery.Cartitems[k].RowPriceNetWithItemRelatedDiscount = itemDelivery.Cartitems[k].RowPriceNet

				if cob.defaultTaxRate > 0.0 {
					taxAmount, err := itemDelivery.Cartitems[k].RowPriceGross.Sub(itemDelivery.Cartitems[k].RowPriceNet)
					if err != nil {
						return err
					}
					itemDelivery.Cartitems[k].RowTaxes[0].Amount = taxAmount
				}

			}

			if itemUpdateCommand.SourceID != nil {
				itemDelivery.Cartitems[k].SourceID = *itemUpdateCommand.SourceID
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
func (cob *DefaultCartBehaviour) AddToCart(ctx context.Context, cart *domaincart.Cart, deliveryCode string, addRequest domaincart.AddRequest) (*domaincart.Cart, domaincart.DeferEvents, error) {
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
		delivery := new(domaincart.Delivery)
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

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

func (cob *DefaultCartBehaviour) buildItemForCart(ctx context.Context, addRequest domaincart.AddRequest) (*domaincart.Item, error) {
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

	return cob.createCartItemFromProduct(addRequest.Qty, addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.AdditionalData, product)
}
func (cob *DefaultCartBehaviour) createCartItemFromProduct(qty int, marketplaceCode string, variantMarketPlaceCode string, additonalData map[string]string, product domain.BasicProduct) (*domaincart.Item, error) {
	item := &domaincart.Item{
		ID:                     strconv.Itoa(rand.Int()),
		ExternalReference:      strconv.Itoa(rand.Int()),
		MarketplaceCode:        marketplaceCode,
		VariantMarketPlaceCode: variantMarketPlaceCode,
		ProductName:            product.BaseData().Title,
		Qty:                    qty,
		AdditionalData:         additonalData,
	}

	currency := product.SaleableData().ActivePrice.GetFinalPrice().Currency()

	if cob.grossPricing {
		item.SinglePriceGross = product.SaleableData().ActivePrice.GetFinalPrice().GetPayable()
		net := item.SinglePriceGross.Clone().Amount().Quo(item.SinglePriceGross.Amount(), big.NewFloat(1+(cob.defaultTaxRate/100)))
		item.SinglePriceNet = priceDomain.NewFromBigFloat(*net, currency).GetPayable()
	} else {
		item.SinglePriceNet = product.SaleableData().ActivePrice.GetFinalPrice().GetPayable()
		gross := item.SinglePriceGross.Clone().Amount().Mul(item.SinglePriceNet.Amount(), big.NewFloat(1+(cob.defaultTaxRate/100)))
		item.SinglePriceGross = priceDomain.NewFromBigFloat(*gross, currency).GetPayable()
	}

	gross := item.SinglePriceGross.Clone().Amount().Mul(item.SinglePriceGross.Amount(), big.NewFloat(float64(qty)))
	item.RowPriceGross = priceDomain.NewFromBigFloat(*gross, currency)
	_ = item.RowPriceGross.FloatAmount()
	net := item.SinglePriceNet.Clone().Amount().Mul(item.SinglePriceNet.Amount(), big.NewFloat(float64(qty)))
	item.RowPriceNet = priceDomain.NewFromBigFloat(*net, currency)
	_ = item.RowPriceNet.FloatAmount()

	item.RowPriceGrossWithDiscount, item.RowPriceNetWithDiscount = item.RowPriceGross, item.RowPriceNet
	item.RowPriceGrossWithItemRelatedDiscount, item.RowPriceNetWithItemRelatedDiscount = item.RowPriceGross, item.RowPriceNet
	if cob.defaultTaxRate > 0.0 {
		taxAmount, err := item.RowPriceGross.Sub(item.RowPriceNet)
		if err != nil {
			return nil, err
		}
		item.RowTaxes = []domaincart.Tax{{
			Amount: taxAmount,
			Type:   "default",
			Rate:   big.NewFloat(cob.defaultTaxRate),
		}}
	}

	item.TotalDiscountAmount = priceDomain.NewZero(currency)
	item.ItemRelatedDiscountAmount = priceDomain.NewZero(currency)
	item.NonItemRelatedDiscountAmount = priceDomain.NewZero(currency)

	return item, nil
}

// CleanCart removes everything from the cart, e.g. deliveries, billing address, etc
func (cob *DefaultCartBehaviour) CleanCart(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCart.Deliveries = []domaincart.Delivery{}
	newCart.AppliedCouponCodes = nil
	newCart.AppliedGiftCards = nil
	newCart.PaymentSelection = nil
	newCart.AdditionalData = domaincart.AdditionalData{}
	newCart.Purchaser = nil
	newCart.BillingAddress = nil
	newCart.Totalitems = nil

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return &newCart, nil, nil
}

// CleanDelivery removes a complete delivery with its items from the cart
func (cob *DefaultCartBehaviour) CleanDelivery(ctx context.Context, cart *domaincart.Cart, deliveryCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	// create delivery if it does not yet exist
	if !newCart.HasDeliveryForCode(deliveryCode) {
		return nil, nil, errors.Errorf("cart.infrastructure.DefaultCartBehaviour: delivery %s not found", deliveryCode)
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
	newCart.Deliveries[newLength] = domaincart.Delivery{}
	newCart.Deliveries = newCart.Deliveries[:newLength]

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdatePurchaser @todo implement when needed
func (cob *DefaultCartBehaviour) UpdatePurchaser(ctx context.Context, cart *domaincart.Cart, purchaser *domaincart.Person, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCart.Purchaser = purchaser

	if additionalData != nil {
		newCart.AdditionalData.CustomAttributes = additionalData.CustomAttributes
	}

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return &newCart, nil, nil
}

// UpdateBillingAddress - updates address
func (cob *DefaultCartBehaviour) UpdateBillingAddress(ctx context.Context, cart *domaincart.Cart, billingAddress domaincart.Address) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCart.BillingAddress = &billingAddress

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return &newCart, nil, nil
}

// UpdateAdditionalData updates additional data
func (cob *DefaultCartBehaviour) UpdateAdditionalData(ctx context.Context, cart *domaincart.Cart, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCart.AdditionalData = *additionalData
	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)

	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on updating additional data")
	}

	return &newCart, nil, nil
}

// UpdatePaymentSelection updates payment on cart
func (cob *DefaultCartBehaviour) UpdatePaymentSelection(ctx context.Context, cart *domaincart.Cart, paymentSelection domaincart.PaymentSelection) (*domaincart.Cart, domaincart.DeferEvents, error) {
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

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return &newCart, nil, nil
}

// UpdateDeliveryInfo updates a delivery info
func (cob *DefaultCartBehaviour) UpdateDeliveryInfo(ctx context.Context, cart *domaincart.Cart, deliveryCode string, deliveryInfoUpdateCommand domaincart.DeliveryInfoUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	deliveryInfo := deliveryInfoUpdateCommand.DeliveryInfo
	deliveryInfo.AdditionalDeliveryInfos = deliveryInfoUpdateCommand.Additional()

	for key, delivery := range newCart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			newCart.Deliveries[key].DeliveryInfo = deliveryInfo
			cob.collectTotals(&newCart)
			err := cob.cartStorage.StoreCart(ctx, &newCart)
			if err != nil {
				return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
			}
			return &newCart, nil, nil
		}
	}
	newCart.Deliveries = append(newCart.Deliveries, domaincart.Delivery{DeliveryInfo: deliveryInfo})

	cob.collectTotals(&newCart)
	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.DefaultCartBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdateDeliveryInfoAdditionalData @todo implement when needed
func (cob *DefaultCartBehaviour) UpdateDeliveryInfoAdditionalData(ctx context.Context, cart *domaincart.Cart, deliveryCode string, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	return cart, nil, nil
}

// GetCart returns the current cart from storage
func (cob *DefaultCartBehaviour) GetCart(ctx context.Context, cartID string) (*domaincart.Cart, error) {
	if !cob.cartStorage.HasCart(ctx, cartID) {
		cob.logger.Info(fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: Cannot get - cart with id %v not existent", cartID))
		return nil, domaincart.ErrCartNotFound
	}

	cart, err := cob.cartStorage.GetCart(ctx, cartID)
	if err != nil {
		cob.logger.Info(fmt.Errorf("cart.infrastructure.DefaultCartBehaviour: get cart from storage: %w ", err))
		return nil, domaincart.ErrCartNotFound
	}

	newCart, err := cart.Clone()
	if err != nil {
		cob.logger.Info(fmt.Errorf("cart clone failed: %w ", err))
		return nil, domaincart.ErrCartNotFound
	}

	return &newCart, nil
}

// StoreNewCart created and stores a new cart.
func (cob *DefaultCartBehaviour) StoreNewCart(ctx context.Context, newCart *domaincart.Cart) (*domaincart.Cart, error) {
	if newCart.ID == "" {
		return nil, errors.New("no id given")
	}
	newCart.DefaultCurrency = cob.defaultCurrency
	cob.collectTotals(newCart)
	return newCart, cob.cartStorage.StoreCart(ctx, newCart)
}

// ApplyVoucher applies a voucher to the cart
func (cob *DefaultCartBehaviour) ApplyVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCartWithVoucher, err := cob.voucherHandler.ApplyVoucher(ctx, &newCart, couponCode)
	if err != nil {
		return nil, nil, err
	}

	cob.collectTotals(newCartWithVoucher)
	err = cob.cartStorage.StoreCart(ctx, newCartWithVoucher)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, newCartWithVoucher)
}

// ApplyAny applies a voucher or giftcard to the cart
func (cob *DefaultCartBehaviour) ApplyAny(ctx context.Context, cart *domaincart.Cart, anyCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	currentCart, deferFunc, err := cob.ApplyVoucher(ctx, cart, anyCode)
	if err == nil {
		// successfully applied as voucher
		return currentCart, deferFunc, nil
	}

	// some error occurred, retry as giftcard
	return cob.ApplyGiftCard(ctx, cart, anyCode)
}

// RemoveVoucher removes a voucher from the cart
func (cob *DefaultCartBehaviour) RemoveVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCartWithoutVoucher, err := cob.voucherHandler.RemoveVoucher(ctx, &newCart, couponCode)
	if err != nil {
		return nil, nil, err
	}

	cob.collectTotals(newCartWithoutVoucher)
	err = cob.cartStorage.StoreCart(ctx, newCartWithoutVoucher)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, newCartWithoutVoucher)
}

// ApplyGiftCard applies a gift card to the cart
// if a GiftCard is applied, it will be added to the array AppliedGiftCards on the cart
func (cob *DefaultCartBehaviour) ApplyGiftCard(ctx context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCartWithGiftCard, err := cob.giftCardHandler.ApplyGiftCard(ctx, &newCart, giftCardCode)
	if err != nil {
		return nil, nil, err
	}

	cob.collectTotals(newCartWithGiftCard)
	err = cob.cartStorage.StoreCart(ctx, newCartWithGiftCard)
	if err != nil {
		return nil, nil, err
	}
	return cob.resetPaymentSelectionIfInvalid(ctx, newCartWithGiftCard)
}

// RemoveGiftCard removes a gift card from the cart
// if a GiftCard is removed, it will be removed from the array AppliedGiftCards on the cart
func (cob *DefaultCartBehaviour) RemoveGiftCard(ctx context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, err
	}

	newCartWithOutGiftCard, err := cob.giftCardHandler.RemoveGiftCard(ctx, &newCart, giftCardCode)
	if err != nil {
		return nil, nil, err
	}

	cob.collectTotals(newCartWithOutGiftCard)
	err = cob.cartStorage.StoreCart(ctx, newCartWithOutGiftCard)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, newCartWithOutGiftCard)
}

// isPaymentSelectionValid checks if the grand total of the cart matches the total of the supplied payment selection
func (cob *DefaultCartBehaviour) checkPaymentSelection(ctx context.Context, cart *domaincart.Cart, paymentSelection domaincart.PaymentSelection) error {
	if paymentSelection == nil {
		return nil
	}
	paymentSelectionTotal := paymentSelection.TotalValue()

	if !cart.GrandTotal.LikelyEqual(paymentSelectionTotal) {
		return errors.New("Payment Total does not match with Grandtotal")
	}
	return nil
}

// resetPaymentSelectionIfInvalid checks for valid paymentselection on givencart and deletes in in case it is invalid
func (cob *DefaultCartBehaviour) resetPaymentSelectionIfInvalid(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
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

func (cob *DefaultCartBehaviour) collectTotals(cart *domaincart.Cart) {
	cart.TotalGiftCardAmount = priceDomain.NewZero(cart.DefaultCurrency)
	cart.GrandTotalWithGiftCards = priceDomain.NewZero(cart.DefaultCurrency)
	cart.GrandTotal = priceDomain.NewZero(cart.DefaultCurrency)
	cart.GrandTotalNet = priceDomain.NewZero(cart.DefaultCurrency)
	cart.GrandTotalNetWithGiftCards = priceDomain.NewZero(cart.DefaultCurrency)
	cart.ShippingNet = priceDomain.NewZero(cart.DefaultCurrency)
	cart.ShippingNetWithDiscounts = priceDomain.NewZero(cart.DefaultCurrency)
	cart.ShippingGross = priceDomain.NewZero(cart.DefaultCurrency)
	cart.ShippingGrossWithDiscounts = priceDomain.NewZero(cart.DefaultCurrency)
	cart.SubTotalGross = priceDomain.NewZero(cart.DefaultCurrency)
	cart.SubTotalNet = priceDomain.NewZero(cart.DefaultCurrency)
	cart.SubTotalGrossWithDiscounts = priceDomain.NewZero(cart.DefaultCurrency)
	cart.SubTotalNetWithDiscounts = priceDomain.NewZero(cart.DefaultCurrency)
	cart.TotalDiscountAmount = priceDomain.NewZero(cart.DefaultCurrency)
	cart.NonItemRelatedDiscountAmount = priceDomain.NewZero(cart.DefaultCurrency)
	cart.ItemRelatedDiscountAmount = priceDomain.NewZero(cart.DefaultCurrency)

	for i := 0; i < len(cart.Deliveries); i++ {
		delivery := &cart.Deliveries[i]
		delivery.SubTotalGross = priceDomain.NewZero(cart.DefaultCurrency)
		delivery.SubTotalNet = priceDomain.NewZero(cart.DefaultCurrency)
		delivery.TotalDiscountAmount = priceDomain.NewZero(cart.DefaultCurrency)
		delivery.SubTotalDiscountAmount = priceDomain.NewZero(cart.DefaultCurrency)
		delivery.NonItemRelatedDiscountAmount = priceDomain.NewZero(cart.DefaultCurrency)
		delivery.ItemRelatedDiscountAmount = priceDomain.NewZero(cart.DefaultCurrency)
		delivery.SubTotalGrossWithDiscounts = priceDomain.NewZero(cart.DefaultCurrency)
		delivery.SubTotalNetWithDiscounts = priceDomain.NewZero(cart.DefaultCurrency)
		delivery.GrandTotal = priceDomain.NewZero(cart.DefaultCurrency)

		if !delivery.ShippingItem.PriceGrossWithDiscounts.IsZero() {
			delivery.GrandTotal = delivery.GrandTotal.ForceAdd(delivery.ShippingItem.PriceGrossWithDiscounts)
			discounts, _ := delivery.ShippingItem.AppliedDiscounts.Sum()
			delivery.TotalDiscountAmount = delivery.TotalDiscountAmount.ForceAdd(discounts)
		}

		for _, cartitem := range delivery.Cartitems {
			delivery.SubTotalGross = delivery.SubTotalGross.ForceAdd(cartitem.RowPriceGross)
			delivery.SubTotalNet = delivery.SubTotalNet.ForceAdd(cartitem.RowPriceNet)
			delivery.TotalDiscountAmount = delivery.TotalDiscountAmount.ForceAdd(cartitem.TotalDiscountAmount)
			delivery.SubTotalDiscountAmount = delivery.SubTotalDiscountAmount.ForceAdd(cartitem.TotalDiscountAmount)
			delivery.NonItemRelatedDiscountAmount = delivery.NonItemRelatedDiscountAmount.ForceAdd(cartitem.NonItemRelatedDiscountAmount)
			delivery.ItemRelatedDiscountAmount = delivery.ItemRelatedDiscountAmount.ForceAdd(cartitem.ItemRelatedDiscountAmount)
			delivery.SubTotalGrossWithDiscounts = delivery.SubTotalGrossWithDiscounts.ForceAdd(cartitem.RowPriceGrossWithDiscount)
			delivery.SubTotalNetWithDiscounts = delivery.SubTotalNetWithDiscounts.ForceAdd(cartitem.RowPriceNetWithDiscount)
			delivery.GrandTotal = delivery.GrandTotal.ForceAdd(cartitem.RowPriceGrossWithDiscount)
		}

		cart.GrandTotal = cart.GrandTotal.ForceAdd(delivery.GrandTotal)
		cart.GrandTotalNet = cart.GrandTotalNet.ForceAdd(delivery.SubTotalNetWithDiscounts).ForceAdd(delivery.ShippingItem.PriceNetWithDiscounts)
		cart.ShippingNet = cart.ShippingNet.ForceAdd(delivery.ShippingItem.PriceNet)
		cart.ShippingNetWithDiscounts = cart.ShippingNetWithDiscounts.ForceAdd(delivery.ShippingItem.PriceNetWithDiscounts)
		cart.ShippingGross = cart.ShippingGross.ForceAdd(delivery.ShippingItem.PriceGross)
		cart.ShippingGrossWithDiscounts = cart.ShippingGrossWithDiscounts.ForceAdd(delivery.ShippingItem.PriceGrossWithDiscounts)
		cart.SubTotalGross = cart.SubTotalGross.ForceAdd(delivery.SubTotalGross)
		cart.SubTotalNet = cart.SubTotalNet.ForceAdd(delivery.SubTotalNet)
		cart.SubTotalGrossWithDiscounts = cart.SubTotalGrossWithDiscounts.ForceAdd(delivery.SubTotalGrossWithDiscounts)
		cart.SubTotalNetWithDiscounts = cart.SubTotalNetWithDiscounts.ForceAdd(delivery.SubTotalNetWithDiscounts)
		cart.TotalDiscountAmount = cart.TotalDiscountAmount.ForceAdd(delivery.TotalDiscountAmount)
		cart.NonItemRelatedDiscountAmount = cart.NonItemRelatedDiscountAmount.ForceAdd(delivery.NonItemRelatedDiscountAmount)
		cart.ItemRelatedDiscountAmount = cart.ItemRelatedDiscountAmount.ForceAdd(delivery.ItemRelatedDiscountAmount)
	}

	for _, totalitem := range cart.Totalitems {
		cart.GrandTotal = cart.GrandTotal.ForceAdd(totalitem.Price)
	}

	sumAppliedGiftCards := priceDomain.NewZero(cart.DefaultCurrency)
	for _, card := range cart.AppliedGiftCards {
		sumAppliedGiftCards = sumAppliedGiftCards.ForceAdd(card.Applied)
	}

	cart.TotalGiftCardAmount = sumAppliedGiftCards

	cart.GrandTotalWithGiftCards, _ = cart.GrandTotal.Sub(cart.TotalGiftCardAmount)
	if cart.GrandTotalWithGiftCards.IsNegative() {
		cart.GrandTotalWithGiftCards = priceDomain.NewZero(cart.DefaultCurrency)
	}

	cart.GrandTotalNetWithGiftCards, _ = cart.GrandTotalNet.Sub(cart.TotalGiftCardAmount)
	if cart.GrandTotalNetWithGiftCards.IsNegative() {
		cart.GrandTotalNetWithGiftCards = priceDomain.NewZero(cart.DefaultCurrency)
	}
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
		Applied:   priceDomain.NewFromInt(10, 100, cart.DefaultCurrency),
		Remaining: priceDomain.NewFromInt(0, 100, cart.DefaultCurrency),
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
