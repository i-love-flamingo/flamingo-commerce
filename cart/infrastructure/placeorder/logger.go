package logger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"golang.org/x/mod/modfile"
)

type (
	// PlaceOrderLoggerAdapter provides an implementation of the Service as email adapter
	PlaceOrderLoggerAdapter struct {
		useFlamingoLog bool
		logAsFile      bool
		logDirectory   string
		logger         flamingo.Logger
	}
)

var (
	_ placeorder.Service = new(PlaceOrderLoggerAdapter)
)

// Inject dependencies
func (e *PlaceOrderLoggerAdapter) Inject(logger flamingo.Logger,
	config *struct {
		UseFlamingoLog bool   `inject:"config:commerce.cart.placeOrderLogger.useFlamingoLog,optional"`
		LogAsFile      bool   `inject:"config:commerce.cart.placeOrderLogger.logAsFile,optional"`
		LogDirectory   string `inject:"config:commerce.cart.placeOrderLogger.logDirectory,optional"`
	}) {
	e.logger = logger.WithField("module", "cart").WithField("category", "emailAdapter")
	if config != nil {
		e.useFlamingoLog = config.UseFlamingoLog
		e.logAsFile = config.LogAsFile
		e.logDirectory = config.LogDirectory
	}
}

// PlaceGuestCart places a guest cart as order email
func (e *PlaceOrderLoggerAdapter) PlaceGuestCart(ctx context.Context, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	return e.placeCart(cart, payment)
}

// PlaceCustomerCart places a customer cart as order email
func (e *PlaceOrderLoggerAdapter) PlaceCustomerCart(ctx context.Context, auth auth.Identity, cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	return e.placeCart(cart, payment)
}

// placeCart
func (e *PlaceOrderLoggerAdapter) placeCart(cart *cartDomain.Cart, payment *placeorder.Payment) (placeorder.PlacedOrderInfos, error) {
	err := e.checkPayment(cart, payment)
	if err != nil {
		return nil, err
	}
	err = e.logOrder(cart, payment)
	if err != nil {
		return nil, err
	}
	var placedOrders placeorder.PlacedOrderInfos
	for _, del := range cart.Deliveries {
		placedOrders = append(placedOrders, placeorder.PlacedOrderInfo{
			OrderNumber:  cart.ID,
			DeliveryCode: del.DeliveryInfo.Code,
		})
	}

	return placedOrders, nil
}

// checkPayment
func (e *PlaceOrderLoggerAdapter) checkPayment(cart *cartDomain.Cart, payment *placeorder.Payment) error {
	if payment == nil && cart.GrandTotal.IsPositive() {
		return errors.New("no valid payment given")
	}
	if cart.GrandTotal.IsPositive() {
		totalPrice, err := payment.TotalValue()
		if err != nil {
			return err
		}
		if !totalPrice.Equal(cart.GrandTotal) {
			return errors.New("payment total does not match with grandtotal")
		}
	}
	return nil
}

// logOrder
func (e *PlaceOrderLoggerAdapter) logOrder(cart *cartDomain.Cart, payment *placeorder.Payment) error {
	if e.useFlamingoLog {
		e.logger.WithField("placeorder", cart.ID).WithField("cart", cart).Info("Order placed and logged")
	}
	if e.logAsFile && e.logDirectory != "" {
		if !modfile.IsDirectoryPath(e.logDirectory) {
			return fmt.Errorf("%v is not a valid directory path", e.logDirectory)
		}
		// Create folder if not exist
		if _, err := os.Stat(e.logDirectory); os.IsNotExist(err) {
			err = os.MkdirAll(e.logDirectory, os.ModePerm)
			if err != nil {
				e.logger.Error(err)
				return err
			}
		}
		type order struct {
			Cart    cartDomain.Cart
			Payment placeorder.Payment
		}
		content, err := json.Marshal(order{
			Cart:    *cart,
			Payment: *payment,
		})
		if err != nil {
			e.logger.Error(err)
			return err
		}
		fileName := fmt.Sprintf("order-%v-%v.json", time.Now().Format(time.RFC3339), cart.ID)
		err = os.WriteFile(path.Join(e.logDirectory, fileName), []byte(content), os.ModePerm)
		if err != nil {
			e.logger.WithField("placeorder", cart.ID).Error(err)
		}
		e.logger.WithField("placeorder", cart.ID).WithField("cart", cart).Info("Order placed")
	}
	return nil
}

// ReserveOrderID returns the reserved order id
func (e *PlaceOrderLoggerAdapter) ReserveOrderID(ctx context.Context, cart *cartDomain.Cart) (string, error) {
	return cart.ID, nil
}

// CancelGuestOrder cancels a guest order
func (e *PlaceOrderLoggerAdapter) CancelGuestOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos) error {
	// since we don't actual place orders we just return nil here
	return nil
}

// CancelCustomerOrder cancels a customer order
func (e *PlaceOrderLoggerAdapter) CancelCustomerOrder(ctx context.Context, orderInfos placeorder.PlacedOrderInfos, auth auth.Identity) error {
	// since we don't actual place orders we just return nil here
	return nil
}
