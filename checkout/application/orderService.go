package application

import (
	"errors"

	"flamingo.me/flamingo-commerce/cart/application"
	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/checkout/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
)

type (

	// PaymentService
	OrderService struct {
		SourcingEngine      *domain.SourcingEngine           `inject:""`
		PaymentService      *PaymentService                  `inject:""`
		Logger              flamingo.Logger                  `inject:""`
		CartService         *application.CartService         `inject:""`
		CartReceiverService *application.CartReceiverService `inject:""`
		DeliveryInfoBuilder cart.DeliveryInfoBuilder         `inject:""`
	}
)

func (os *OrderService) SetSources(ctx web.Context) error {
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return err
	}
	err = os.SourcingEngine.SetSourcesForCartItems(ctx, decoratedCart)
	if err != nil {
		os.Logger.WithField("category", "checkout.orderService").Error("Error while getting sources: %v", err)
		return errors.New("Error while setting sources.")
	}
	return nil
}

func (os *OrderService) PlaceOrder(ctx web.Context, decoratedCart *cart.DecoratedCart, payment *cart.CartPayment) (orderid string, orderError error) {
	validationResult := os.CartService.ValidateCart(ctx, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return "", errors.New("Cart is Invalid.")
	}
	return os.CartService.PlaceOrder(ctx, payment)
}

func (os *OrderService) CurrentCartSaveInfos(ctx web.Context, billingAddress *cart.Address, shippingAddress *cart.Address, purchaser *cart.Person) error {
	os.Logger.Debug("CurrentCartSaveInfos call billingAddress:%v shippingAddress:%v payment:%v", billingAddress, shippingAddress)

	if billingAddress == nil {
		os.Logger.Warn("CurrentCartSaveInfos called without billing address")
		return errors.New("Billing Address is missing")
	}
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx)
	if err != nil {
		os.Logger.Error("CurrentCartSaveInfos GetDecoratedCart Error %v", err)
		return err
	}

	//update Billing
	err = os.CartService.UpdateBillingAddress(ctx, billingAddress)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder UpdateBillingAddress Error %v", err)
		return err
	}

	//Update ShippingAddress on ALL Deliveries in the Cart if given
	// Maybe later we need to support different shipping addresses in the Checkout
	if shippingAddress != nil {
		for _, d := range decoratedCart.Cart.Deliveries {
			newDeliveryInfo := d.DeliveryInfo
			newDeliveryInfo.DeliveryLocation.Address = shippingAddress
			err = os.CartService.UpdateDeliveryInfo(ctx, d.DeliveryInfo.Code, newDeliveryInfo)
			if err != nil {
				os.Logger.Error("OnStepCurrentCartPlaceOrder UpdateDeliveryInfosAndBilling Error %v", err)
				return err
			}
		}

	}

	//Update Purchaser
	err = os.CartService.UpdatePurchaser(ctx, purchaser, nil)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder UpdatePurchaser Error %v", err)
		return err
	}

	//After setting DeliveryInfos - call SourcingEnginge (this will reload the cart and update all items!)
	err = os.SetSources(ctx)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder SetSources Error %v", err)
		return err
	}
	return nil
}

//OnStepCurrentCartPlaceOrder - probably the best choice for a simple checkout
// Assumptions: Only one BuildDeliveryInfo is used on the cart!
func (os *OrderService) CurrentCartPlaceOrder(ctx web.Context, payment cart.CartPayment) (orderid string, orderError error) {
	decoratedCart, err := os.CartReceiverService.ViewDecoratedCart(ctx)
	if err != nil {
		os.Logger.Error("OnStepCurrentCartPlaceOrder GetDecoratedCart Error %v", err)
		return "", err
	}

	validationResult := os.CartService.ValidateCart(ctx, decoratedCart)
	if !validationResult.IsValid() {
		os.Logger.Warn("Try to place an invalid cart")
		return "", errors.New("Cart is Invalid.")
	}

	orderid, orderError = os.PlaceOrder(ctx, decoratedCart, &payment)

	if orderError != nil {
		os.Logger.WithField("category", "checkout.orderService").Error("Error during place Order: %v", err)
		return "", errors.New("Error while placing the order. Please contact customer support.")
	}
	return orderid, nil
}
