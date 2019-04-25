package infrastructure

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	domaincart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
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
		defaultTaxRate          float64
	}

	//CartStorage Interface - might be implemented by other persistence types later as well
	CartStorage interface {
		GetCart(id string) (*domaincart.Cart, error)
		HasCart(id string) bool
		StoreCart(cart *domaincart.Cart) error
	}
)

var (
	_ domaincart.ModifyBehaviour = (*InMemoryBehaviour)(nil)
)

// Inject dependencies
func (cob *InMemoryBehaviour) Inject(
	CartStorage CartStorage,
	ProductService domain.ProductService,
	Logger flamingo.Logger,
	itemBuilderProvider domaincart.ItemBuilderProvider,
	deliveryBuilderProvider domaincart.DeliveryBuilderProvider,
	cartBuilderProvider domaincart.BuilderProvider,
	eventPublisher application.EventPublisher,
	config *struct {
		DefaultTaxRate float64 `inject:"config:commerce.cart.inMemoryCartServiceAdapter.defaultTaxRate,optional"`
	},
) {
	cob.cartStorage = CartStorage
	cob.productService = ProductService
	cob.logger = Logger
	cob.itemBuilderProvider = itemBuilderProvider
	cob.deliveryBuilderProvider = deliveryBuilderProvider
	cob.cartBuilderProvider = cartBuilderProvider
	if config != nil {
		cob.defaultTaxRate = config.DefaultTaxRate
	}
}

// DeleteItem removes an item from the cart
func (cob *InMemoryBehaviour) DeleteItem(ctx context.Context, cart *domaincart.Cart, itemID string, deliveryCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot delete - Guestcart with id %v not existent", cart.ID)
	}

	if newDelivery, ok := cart.GetDeliveryByCode(deliveryCode); ok {
		cob.logger.WithField(flamingo.LogKeyCategory, "inmemorybehaviour").Info("Inmemory Service Delete %v in %#v", itemID, newDelivery.Cartitems)
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
func (cob *InMemoryBehaviour) UpdateItem(ctx context.Context, cart *domaincart.Cart, itemID string, deliveryCode string, itemUpdateCommand domaincart.ItemUpdateCommand) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot add - Guestcart with id %v not existend", cart.ID)
	}

	itemBuilder := cob.itemBuilderProvider()
	if delivery, ok := cart.GetDeliveryByCode(deliveryCode); ok {
		cob.logger.WithField(flamingo.LogKeyCategory, "inmemorybehaviour").Info("Inmemory Service Update %v in %#v", itemID, delivery.Cartitems)
		for _, item := range delivery.Cartitems {
			if itemID == item.ID {
				itemBuilder.SetFromItem(item).SetQty(*itemUpdateCommand.Qty).AddTaxInfo("default", big.NewFloat(cob.defaultTaxRate), nil).CalculatePricesAndTax()
				newItem, err := itemBuilder.Build()
				if err != nil {
					return nil, nil, err
				}
				for k, currentItem := range delivery.Cartitems {
					if currentItem.ID == itemID {
						delivery.Cartitems[k] = *newItem
					}
				}
			}
		}

		// update the delivery with the new info
		for j, delivery := range cart.Deliveries {
			if deliveryCode == delivery.DeliveryInfo.Code {
				cart.Deliveries[j] = delivery
			}
		}
	}

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

// AddToCart add an item to the cart
func (cob *InMemoryBehaviour) AddToCart(ctx context.Context, cart *domaincart.Cart, deliveryCode string, addRequest domaincart.AddRequest) (*domaincart.Cart, domaincart.DeferEvents, error) {

	if cart != nil && !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot add - Guestcart with id %v not existend", cart.ID)
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
	itemBuilder.SetQty(addRequest.Qty).AddTaxInfo("default", big.NewFloat(cob.defaultTaxRate), nil).SetByProduct(product).SetID(strconv.Itoa(rand.Int())).SetExternalReference(strconv.Itoa(rand.Int()))

	return itemBuilder.Build()
}

// CleanCart removes all deliveries and their items from the cart
func (cob *InMemoryBehaviour) CleanCart(ctx context.Context, cart *domaincart.Cart) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if !cob.cartStorage.HasCart(cart.ID) {
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot delete - Guestcart with id %v not existend", cart.ID)
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
		return nil, nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot delete - Guestcart with id %v not existend", cart.ID)
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

// UpdatePurchaser @todo implement when needed
func (cob *InMemoryBehaviour) UpdatePurchaser(ctx context.Context, cart *domaincart.Cart, purchaser *domaincart.Person, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	return cart, nil, nil
}

// UpdateBillingAddress - updates address
func (cob *InMemoryBehaviour) UpdateBillingAddress(ctx context.Context, cart *domaincart.Cart, billingAddress domaincart.Address) (*domaincart.Cart, domaincart.DeferEvents, error) {

	cart.BillingAdress = &billingAddress

	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cart.infrastructure.InMemoryBehaviour: error on saving cart")
	}

	return cart, nil, nil
}

// UpdateAdditionalData @todo implement when needed
func (cob *InMemoryBehaviour) UpdateAdditionalData(ctx context.Context, cart *domaincart.Cart, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	return cart, nil, nil
}

//UpdatePaymentSelection updates payment on cart
func (cob *InMemoryBehaviour) UpdatePaymentSelection(ctx context.Context, cart *domaincart.Cart, paymentSelection *domaincart.PaymentSelection) (*domaincart.Cart, domaincart.DeferEvents, error) {
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
func (cob *InMemoryBehaviour) UpdateDeliveryInfoAdditionalData(ctx context.Context, cart *domaincart.Cart, deliveryCode string, additionalData *domaincart.AdditionalData) (*domaincart.Cart, domaincart.DeferEvents, error) {
	return cart, nil, nil
}

// GetCart returns the current cart from storage
func (cob *InMemoryBehaviour) GetCart(ctx context.Context, cartID string) (*domaincart.Cart, error) {
	if cob.cartStorage.HasCart(cartID) {
		// if cart exists, there is no error ;)
		cart, err := cob.cartStorage.GetCart(cartID)
		if err == nil {
			return cart, nil
		}
	}

	return nil, fmt.Errorf("cart.infrastructure.InMemoryBehaviour: Cannot get - Guestcart with id %v not existent", cartID)
}

// StoreCart in the memory
func (cob *InMemoryBehaviour) StoreCart(cart *domaincart.Cart) error {
	return cob.cartStorage.StoreCart(cart)
}

// ApplyVoucher applies a voucher to the cart
func (cob *InMemoryBehaviour) ApplyVoucher(ctx context.Context, cart *domaincart.Cart, couponCode string) (*domaincart.Cart, domaincart.DeferEvents, error) {
	if couponCode != "valid" {
		err := errors.New("Code invalid")
		return nil, nil, err
	}

	coupon := domaincart.CouponCode{
		Code: couponCode,
	}
	cart.AppliedCouponCodes = append(cart.AppliedCouponCodes, coupon)
	err := cob.cartStorage.StoreCart(cart)
	if err != nil {
		return nil, nil, err
	}

	return cob.resetPaymentSelectionIfInvalid(ctx, cart)
}

func (cob *InMemoryBehaviour) isCurrentPaymentSelectionValid(ctx context.Context, cart *domaincart.Cart) bool {
	return cob.checkPaymentSelection(ctx, cart, cart.PaymentSelection) == nil
}

// isPaymentSelectionValid checks if the grand total of the cart matches the total of the supplied payment selection
func (cob *InMemoryBehaviour) checkPaymentSelection(ctx context.Context, cart *domaincart.Cart, paymentSelection *domaincart.PaymentSelection) error {
	if paymentSelection == nil {
		return nil
	}
	paymentSelectionTotal, err := paymentSelection.TotalValue()
	if err != nil {
		return err
	}
	if !cart.GrandTotal().Equal(paymentSelectionTotal) {
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
		defers = append(defers, &application.PaymentSelectionHasBeenResetEvent{Cart: cart})
		return cart, defers, err
	}

	return cart, nil, nil
}
