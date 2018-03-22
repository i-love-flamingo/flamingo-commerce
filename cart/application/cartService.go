package application

import (
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

		//DefaultDeliveryMethodForValidation - used for calling the CartValidator (this is something that might get obsolete if the Cart and the CartItems have theire Deliverymethod "saved")
		DefaultDeliveryMethodForValidation string `inject:"config:cart.validation.defaultDeliveryMethod,optional"`
	}
)

// ValidateCart validates a carts content
func (cs CartService) ValidateCart(ctx web.Context, decoratedCart cartDomain.DecoratedCart) cartDomain.CartValidationResult {
	if cs.CartValidator != nil {
		// TODO pass delivery Method from CART - once cart supports this!
		result := cs.CartValidator.Validate(ctx, decoratedCart, cs.DefaultDeliveryMethodForValidation)

		return result
	}

	return cartDomain.CartValidationResult{}
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

	item.Qty = qty
	err = behaviour.UpdateItem(ctx, cart, itemId, *item)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "UpdateItemQty").Error(err)
		return err
	}
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
	err = behaviour.DeleteItem(ctx, cart, itemId)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "DeleteItem").Error(err)
		return err
	}
	cs.EventPublisher.PublishChangedQtyInCartEvent(ctx, item, qtyBefore, 0, cart.ID)
	return nil
}

// SetShippingInformation
func (cs *CartService) SetShippingInformation(ctx web.Context, shippingAddress *cartDomain.Address, billingAddress *cartDomain.Address, shippingCarrierCode string, shippingMethodCode string) error {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "SetShippingInformation").Error(err)
		return err
	}

	err = behaviour.SetShippingInformation(ctx, cart, shippingAddress, billingAddress, shippingCarrierCode, shippingMethodCode)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "SetShippingInformation").Error(err)
	}
	return nil
}

// PlaceOrder
func (cs *CartService) PlaceOrder(ctx web.Context, payment *cartDomain.PaymentInfo) (string, error) {
	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
	if err != nil {
		return "", err
	}

	orderNumber, err := behaviour.PlaceOrder(ctx, cart, payment)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "PlaceOrder").Error(err)
		return "", err
	}
	cs.EventPublisher.PublishOrderPlacedEvent(ctx, cart, orderNumber)
	return orderNumber, err
}

// AddProduct Add a product
func (cs *CartService) AddProduct(ctx web.Context, addRequest cartDomain.AddRequest) error {
	addRequest, err := cs.checkProductAndEnrichAddRequest(ctx, addRequest)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)
		return err
	}

	cart, behaviour, err := cs.CartReceiverService.GetCart(ctx)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)
		return err
	}

	err = behaviour.AddToCart(ctx, cart, addRequest)
	if err != nil {
		cs.Logger.WithField("category", "cartService").WithField("subCategory", "AddProduct").Error(err)
		return err
	}
	cs.publishAddtoCartEvent(ctx, *cart, addRequest)

	return nil
}

// checkProductAndEnrichAddRequest existence and validate with productService
func (cs *CartService) checkProductAndEnrichAddRequest(ctx web.Context, addRequest cartDomain.AddRequest) (cartDomain.AddRequest, error) {
	product, err := cs.ProductService.Get(ctx, addRequest.MarketplaceCode)
	if err != nil {
		return addRequest, errors.New("cart.application.cartservice - AddProduct:Product not found")
	}

	if product.Type() == productDomain.TYPECONFIGURABLE {
		if addRequest.VariantMarketplaceCode == "" {
			return addRequest, errors.New("cart.application.cartservice - AddProduct:No Variant given for configurable product")
		}

		configurableProduct := product.(productDomain.ConfigurableProduct)
		_, err := configurableProduct.Variant(addRequest.VariantMarketplaceCode)
		if err != nil {
			return addRequest, errors.New("cart.application.cartservice - AddProduct:Product has not the given variant")
		}
	}

	//Now Validate the Item with the optional registered ItemValidator
	if cs.ItemValidator != nil {
		return addRequest, cs.ItemValidator.Validate(ctx, addRequest, product)
	}

	return addRequest, nil
}

func (cs *CartService) publishAddtoCartEvent(ctx web.Context, currentCart cartDomain.Cart, addRequest cartDomain.AddRequest) {
	if cs.EventPublisher != nil {
		cs.EventPublisher.PublishAddToCartEvent(ctx, addRequest.MarketplaceCode, addRequest.VariantMarketplaceCode, addRequest.Qty)
	}
}
