package infrastructure

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/pkg/errors"

	domaincart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/events"
	priceDomain "flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name GiftCardHandler --case snake
//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name VoucherHandler --case snake

type (
	// DefaultCartBehaviour defines the default cart order behaviour
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

	// GiftCardHandler enables the projects to have specific GiftCard handling
	GiftCardHandler interface {
		ApplyGiftCard(ctx context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, error)
		RemoveGiftCard(ctx context.Context, cart *domaincart.Cart, giftCardCode string) (*domaincart.Cart, error)
	}

	// VoucherHandler enables the projects to have specific Voucher handling
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

const (
	logCategory = "DefaultCartBehaviour"
)

// Inject dependencies
func (cob *DefaultCartBehaviour) Inject(
	cartStorage CartStorage,
	productService domain.ProductService,
	logger flamingo.Logger,
	voucherHandler VoucherHandler,
	giftCardHandler GiftCardHandler,
	config *struct {
		DefaultTaxRate  float64 `inject:"config:commerce.cart.defaultCartAdapter.defaultTaxRate,optional"`
		ProductPricing  string  `inject:"config:commerce.cart.defaultCartAdapter.productPrices"`
		DefaultCurrency string  `inject:"config:commerce.cart.defaultCartAdapter.defaultCurrency"`
	},
) {
	cob.cartStorage = cartStorage
	cob.productService = productService
	cob.logger = logger
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
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error removing cart: %w", err)
	}

	return cart, nil, nil
}

// Restore supplied cart (implements CompleteBehaviour)
func (cob *DefaultCartBehaviour) Restore(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	cob.collectTotals(cart)

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// DeleteItem removes an item from the cart
func (cob *DefaultCartBehaviour) DeleteItem(ctx context.Context, cart *domaincart.Cart, itemID string, deliveryCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: %w for cart id %q during delete", domaincart.ErrCartNotFound, cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	if newDelivery, ok := newCart.GetDeliveryByCode(deliveryCode); ok {
		cob.logger.WithContext(ctx).WithField(flamingo.LogKeyCategory, logCategory).Info("DefaultCartBehaviour Delete %v in %#v", itemID, newDelivery.Cartitems)

		for index, item := range newDelivery.Cartitems {
			if item.ID == itemID {
				newDelivery.Cartitems = append(newDelivery.Cartitems[:index], newDelivery.Cartitems[index+1:]...)
				break
			}
		}

		// update the delivery with the new info
		for index, delivery := range newCart.Deliveries {
			if deliveryCode == delivery.DeliveryInfo.Code {
				newCart.Deliveries[index] = *newDelivery
			}
		}
	}

	cob.collectTotals(&newCart)

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdateItem updates a cart item
func (cob *DefaultCartBehaviour) UpdateItem(ctx context.Context, cart *domaincart.Cart, itemUpdateCommand domaincart.ItemUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
	return cob.UpdateItems(ctx, cart, []domaincart.ItemUpdateCommand{itemUpdateCommand})
}

// UpdateItems updates multiple cart items
func (cob *DefaultCartBehaviour) UpdateItems(ctx context.Context, cart *domaincart.Cart, itemUpdateCommands []domaincart.ItemUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: %w for cart id %q during update", domaincart.ErrCartNotFound, cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
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
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

func (cob *DefaultCartBehaviour) updateItem(ctx context.Context, cart *domaincart.Cart, itemUpdateCommand domaincart.ItemUpdateCommand) error {
	itemDelivery, err := cart.GetDeliveryByItemID(itemUpdateCommand.ItemID)
	if err != nil {
		return fmt.Errorf("DefaultCartBehaviour: error finding delivery of item: %w", err)
	}

	cob.logger.WithContext(ctx).WithField(flamingo.LogKeyCategory, logCategory).Info("DefaultCartBehaviour Update %v in %#v", itemUpdateCommand.ItemID, itemDelivery.Cartitems)

	for index, item := range itemDelivery.Cartitems {
		if itemUpdateCommand.ItemID != item.ID {
			continue
		}

		if itemUpdateCommand.SourceID != nil {
			itemDelivery.Cartitems[index].SourceID = *itemUpdateCommand.SourceID
		}

		if itemUpdateCommand.AdditionalData != nil {
			itemDelivery.Cartitems[index].AdditionalData = itemUpdateCommand.AdditionalData
		}

		if itemUpdateCommand.BundleConfiguration != nil {
			itemDelivery.Cartitems[index].BundleConfig = itemUpdateCommand.BundleConfiguration
		}

		if itemUpdateCommand.Qty == nil {
			break
		}

		// in case of qty 0 remove the item from the delivery
		if *itemUpdateCommand.Qty == 0 {
			itemDelivery.Cartitems = append(itemDelivery.Cartitems[:index], itemDelivery.Cartitems[index+1:]...)
			break
		}

		itemDelivery.Cartitems[index].Qty = *itemUpdateCommand.Qty

		gross := item.SinglePriceGross.Clone().Amount().Mul(item.SinglePriceGross.Amount(), big.NewFloat(float64(*itemUpdateCommand.Qty)))
		itemDelivery.Cartitems[index].RowPriceGross = priceDomain.NewFromBigFloat(*gross, item.SinglePriceGross.Currency())

		net := item.SinglePriceNet.Clone().Amount().Mul(item.SinglePriceNet.Amount(), big.NewFloat(float64(*itemUpdateCommand.Qty)))
		itemDelivery.Cartitems[index].RowPriceNet = priceDomain.NewFromBigFloat(*net, item.SinglePriceNet.Currency())

		itemDelivery.Cartitems[index].RowPriceGrossWithDiscount = itemDelivery.Cartitems[index].RowPriceGross
		if rowPriceGrossWithDiscount, err := itemDelivery.Cartitems[index].RowPriceGross.Sub(itemDelivery.Cartitems[index].TotalDiscountAmount); err == nil {
			itemDelivery.Cartitems[index].RowPriceGrossWithDiscount = rowPriceGrossWithDiscount
		}

		itemDelivery.Cartitems[index].RowPriceNetWithDiscount = itemDelivery.Cartitems[index].RowPriceNet
		if rowPriceNetWithDiscount, err := itemDelivery.Cartitems[index].RowPriceNet.Sub(itemDelivery.Cartitems[index].TotalDiscountAmount); err == nil {
			itemDelivery.Cartitems[index].RowPriceNetWithDiscount = rowPriceNetWithDiscount
		}

		itemDelivery.Cartitems[index].RowPriceGrossWithItemRelatedDiscount = itemDelivery.Cartitems[index].RowPriceGross
		if rowPriceGrossWithItemRelatedDiscount, err := itemDelivery.Cartitems[index].RowPriceGross.Sub(itemDelivery.Cartitems[index].ItemRelatedDiscountAmount); err == nil {
			itemDelivery.Cartitems[index].RowPriceGrossWithItemRelatedDiscount = rowPriceGrossWithItemRelatedDiscount
		}

		itemDelivery.Cartitems[index].RowPriceNetWithItemRelatedDiscount = itemDelivery.Cartitems[index].RowPriceNet
		if rowPriceNetWithItemRelatedDiscount, err := itemDelivery.Cartitems[index].RowPriceNet.Sub(itemDelivery.Cartitems[index].ItemRelatedDiscountAmount); err == nil {
			itemDelivery.Cartitems[index].RowPriceNetWithItemRelatedDiscount = rowPriceNetWithItemRelatedDiscount
		}

		if cob.defaultTaxRate > 0.0 {
			taxAmount, err := itemDelivery.Cartitems[index].RowPriceGross.Sub(itemDelivery.Cartitems[index].RowPriceNet)
			if err != nil {
				return fmt.Errorf("DefaultCartBehaviour: error calculating tax amount: %w", err)
			}

			itemDelivery.Cartitems[index].RowTaxes[0].Amount = taxAmount
		}
	}

	// update the delivery with the new info
	for index, delivery := range cart.Deliveries {
		if itemDelivery.DeliveryInfo.Code == delivery.DeliveryInfo.Code {
			cart.Deliveries[index] = *itemDelivery
		}
	}

	return nil
}

// AddToCart add an item to the cart
func (cob *DefaultCartBehaviour) AddToCart(ctx context.Context, cart *domaincart.Cart, deliveryCode string, addRequest domaincart.AddRequest) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if cart != nil && !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: %w for cart id %q during add", domaincart.ErrCartNotFound, cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	// create delivery if it does not yet exist
	if !newCart.HasDeliveryForCode(deliveryCode) {
		// create delivery and add item
		delivery := new(domaincart.Delivery)
		delivery.DeliveryInfo.Code = deliveryCode
		newCart.Deliveries = append(newCart.Deliveries, *delivery)
	}

	delivery, err := cob.addToDelivery(ctx, newCart.GetDeliveryByCodeWithoutBool(deliveryCode), addRequest)
	if err != nil {
		return nil, nil, err
	}

	for k, del := range newCart.Deliveries {
		if del.DeliveryInfo.Code == delivery.DeliveryInfo.Code {
			newCart.Deliveries[k] = *delivery
		}
	}

	cob.collectTotals(&newCart)

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// has cart current delivery, check if there is an item present for this delivery
func (cob *DefaultCartBehaviour) addToDelivery(ctx context.Context, delivery *domaincart.Delivery, addRequest domaincart.AddRequest) (*domaincart.Delivery, error) {
	for index, item := range delivery.Cartitems {
		if item.MarketplaceCode != addRequest.MarketplaceCode {
			continue
		}

		if item.VariantMarketPlaceCode != addRequest.VariantMarketplaceCode {
			continue
		}

		if !item.BundleConfig.Equals(addRequest.BundleConfiguration) {
			continue
		}

		addRequest.Qty += item.Qty

		if addRequest.AdditionalData == nil && len(item.AdditionalData) > 0 {
			addRequest.AdditionalData = make(map[string]string)
		}

		// copy additional data
		for key, val := range item.AdditionalData {
			addRequest.AdditionalData[key] = val
		}

		// create and add new item
		cartItem, err := cob.buildItemForCart(ctx, addRequest)
		if err != nil {
			return nil, err
		}

		delivery.Cartitems[index] = *cartItem

		return delivery, nil
	}

	// create and add new item
	cartItem, err := cob.buildItemForCart(ctx, addRequest)
	if err != nil {
		return nil, err
	}

	delivery.Cartitems = append(delivery.Cartitems, *cartItem)

	return delivery, nil
}

func (cob *DefaultCartBehaviour) buildItemForCart(ctx context.Context, addRequest domaincart.AddRequest) (*domaincart.Item, error) {
	// create and add new item
	product, err := cob.productService.Get(ctx, addRequest.MarketplaceCode)
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	// Get variant of configurable product
	if configurableProduct, ok := product.(domain.ConfigurableProduct); ok && addRequest.VariantMarketplaceCode != "" {
		productWithActiveVariant, err := configurableProduct.GetConfigurableWithActiveVariant(addRequest.VariantMarketplaceCode)
		if err != nil {
			return nil, fmt.Errorf("error getting configurable with active variant: %w", err)
		}

		product = productWithActiveVariant
	}

	if bundleProduct, ok := product.(domain.BundleProduct); ok && len(addRequest.BundleConfiguration) != 0 {
		bundleConfig := addRequest.BundleConfiguration

		bundleProductWithActiveChoices, err := bundleProduct.GetBundleProductWithActiveChoices(bundleConfig)
		if err != nil {
			return nil, fmt.Errorf("error getting bundle with active choices: %w", err)
		}

		product = bundleProductWithActiveChoices
	}

	return cob.createCartItemFromProduct(addRequest.Qty, addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.AdditionalData, addRequest.BundleConfiguration, product)
}
func (cob *DefaultCartBehaviour) createCartItemFromProduct(qty int, marketplaceCode string, variantMarketPlaceCode string,
	additonalData map[string]string, bundleConfig domain.BundleConfiguration, product domain.BasicProduct) (*domaincart.Item, error) {
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
			return nil, fmt.Errorf("DefaultCartBehaviour: error calculating tax amount: %w", err)
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

	item.BundleConfig = bundleConfig

	return item, nil
}

// CleanCart removes everything from the cart, e.g. deliveries, billing address, etc
func (cob *DefaultCartBehaviour) CleanCart(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: %w for cart id %q during clean cart", domaincart.ErrCartNotFound, cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
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
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// CleanDelivery removes a complete delivery with its items from the cart
func (cob *DefaultCartBehaviour) CleanDelivery(ctx context.Context, cart *domaincart.Cart, deliveryCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(ctx, cart.ID) {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: %w for cart id %q during clean delivery", domaincart.ErrCartNotFound, cart.ID)
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	// create delivery if it does not yet exist
	if !newCart.HasDeliveryForCode(deliveryCode) {
		return nil, nil, errors.Errorf("DefaultCartBehaviour: delivery %s not found", deliveryCode)
	}

	var position int

	for index, delivery := range newCart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			position = index
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
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, &newCart)
}

// UpdatePurchaser - updates purchaser
func (cob *DefaultCartBehaviour) UpdatePurchaser(ctx context.Context, cart *domaincart.Cart, purchaser *domaincart.Person, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	newCart.Purchaser = purchaser

	if additionalData != nil {
		if newCart.AdditionalData.CustomAttributes == nil {
			newCart.AdditionalData.CustomAttributes = make(map[string]string)
		}

		for key, val := range additionalData.CustomAttributes {
			newCart.AdditionalData.CustomAttributes[key] = val
		}
	}

	cob.collectTotals(&newCart)

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// UpdateBillingAddress - updates address
func (cob *DefaultCartBehaviour) UpdateBillingAddress(ctx context.Context, cart *domaincart.Cart, billingAddress domaincart.Address) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	newCart.BillingAddress = &billingAddress

	cob.collectTotals(&newCart)

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// UpdateAdditionalData updates additional data
func (cob *DefaultCartBehaviour) UpdateAdditionalData(ctx context.Context, cart *domaincart.Cart, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	newCart.AdditionalData = *additionalData

	cob.collectTotals(&newCart)

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error updating additional data: %w", err)
	}

	return &newCart, nil, nil
}

// UpdatePaymentSelection updates payment on cart
func (cob *DefaultCartBehaviour) UpdatePaymentSelection(ctx context.Context, cart *domaincart.Cart, paymentSelection domaincart.PaymentSelection) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
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
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return &newCart, nil, nil
}

// UpdateDeliveryInfo updates a delivery info
func (cob *DefaultCartBehaviour) UpdateDeliveryInfo(ctx context.Context, cart *domaincart.Cart, deliveryCode string, deliveryInfoUpdateCommand domaincart.DeliveryInfoUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	deliveryInfo := deliveryInfoUpdateCommand.DeliveryInfo
	deliveryInfo.AdditionalDeliveryInfos = deliveryInfoUpdateCommand.Additional()

	for key, delivery := range newCart.Deliveries {
		if delivery.DeliveryInfo.Code == deliveryCode {
			newCart.Deliveries[key].DeliveryInfo = deliveryInfo

			cob.collectTotals(&newCart)

			err := cob.cartStorage.StoreCart(ctx, &newCart)
			if err != nil {
				return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
			}

			return &newCart, nil, nil
		}
	}

	newCart.Deliveries = append(newCart.Deliveries, domaincart.Delivery{DeliveryInfo: deliveryInfo})

	cob.collectTotals(&newCart)

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
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
		err := fmt.Errorf("DefaultCartBehaviour: %w for cart id %q during get", domaincart.ErrCartNotFound, cartID)
		cob.logger.WithField(flamingo.LogKeyCategory, logCategory).Info(err)

		return nil, err
	}

	cart, err := cob.cartStorage.GetCart(ctx, cartID)
	if err != nil {
		cob.logger.WithField(flamingo.LogKeyCategory, logCategory).Info(fmt.Errorf("DefaultCartBehaviour: get cart from storage: %w ", err))
		return nil, domaincart.ErrCartNotFound
	}

	newCart, err := cart.Clone()
	if err != nil {
		cob.logger.WithField(flamingo.LogKeyCategory, logCategory).Info(fmt.Errorf("DefaultCartBehaviour: cart clone failed: %w ", err))
		return nil, domaincart.ErrCartNotFound
	}

	return &newCart, nil
}

// StoreNewCart created and stores a new cart.
func (cob *DefaultCartBehaviour) StoreNewCart(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, error) {
	if cart.ID == "" {
		return nil, errors.New("no id given")
	}

	newCart, err := cart.Clone()
	if err != nil {
		return nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
	}

	newCart.DefaultCurrency = cob.defaultCurrency

	cob.collectTotals(&newCart)

	err = cob.cartStorage.StoreCart(ctx, &newCart)
	if err != nil {
		return nil, fmt.Errorf("DefaultCartBehaviour: error saving cart: %w", err)
	}

	return &newCart, nil
}

// ApplyVoucher applies a voucher to the cart
func (cob *DefaultCartBehaviour) ApplyVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	newCart, err := cart.Clone()
	if err != nil {
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
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
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
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
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
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
		return nil, nil, fmt.Errorf("DefaultCartBehaviour: error cloning cart: %w", err)
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
func (cob *DefaultCartBehaviour) checkPaymentSelection(_ context.Context, cart *domaincart.Cart, paymentSelection domaincart.PaymentSelection) error {
	if paymentSelection == nil {
		return nil
	}

	paymentSelectionTotal := paymentSelection.TotalValue()

	if !cart.GrandTotal.LikelyEqual(paymentSelectionTotal) {
		return errors.New("Payment Total does not match with Grandtotal")
	}

	return nil
}

// resetPaymentSelectionIfInvalid checks for valid paymentselection on given cart and deletes in case it is invalid
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

// ApplyVoucher is called when applying a voucher
func (DefaultVoucherHandler) ApplyVoucher(_ context.Context, cart *domaincart.Cart, _ string) (*domaincart.Cart, error) {
	return cart, nil
}

// RemoveVoucher is called when removing a voucher
func (DefaultVoucherHandler) RemoveVoucher(_ context.Context, cart *domaincart.Cart, _ string) (*domaincart.Cart, error) {
	return cart, nil
}

// ApplyGiftCard is called when applying a gift card
func (DefaultGiftCardHandler) ApplyGiftCard(_ context.Context, cart *domaincart.Cart, _ string) (*domaincart.Cart, error) {
	return cart, nil
}

// RemoveGiftCard is called when removing a gift card
func (DefaultGiftCardHandler) RemoveGiftCard(_ context.Context, cart *domaincart.Cart, _ string) (*domaincart.Cart, error) {
	return cart, nil
}
