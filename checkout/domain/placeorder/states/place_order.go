package states

import (
	"context"
	"fmt"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// PlaceOrder state
	PlaceOrder struct {
		orderService *application.OrderService
		cartService  *cartApplication.CartService
	}

	// PlaceOrderRollbackData needed for rollbacks
	PlaceOrderRollbackData struct {
		order application.PlaceOrderInfo
	}
)

var _ process.State = PlaceOrder{}

// Inject dependencies
func (po *PlaceOrder) Inject(
	orderService *application.OrderService,
	cartService *cartApplication.CartService,
) *PlaceOrder {
	po.orderService = orderService
	po.cartService = cartService

	return po
}

// Name get state name
func (PlaceOrder) Name() string {
	return "PlaceOrder"
}

// Run the state operations
func (po PlaceOrder) Run(ctx context.Context, p *process.Process) process.RunResult {
	cart := p.Context().Cart

	paymentGateway, err := po.orderService.GetPaymentGateway(ctx, interfaces.OfflineWebCartPaymentGatewayCode)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	payment, err := paymentGateway.OrderPaymentFromFlow(ctx, &cart, p.Context().UUID)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	// Todo: need new function in orderService to place a provided cart similar to:
	_, _ = po.orderService.CurrentCartPlaceOrder(ctx, web.SessionFromContext(ctx), *payment)

	// todo: next state depending on early place.. success / validate payment
	return process.RunResult{}
}

// Rollback the state operations
func (po PlaceOrder) Rollback(data process.RollbackData) error {
	rollbackData, ok := data.(PlaceOrderRollbackData)
	if !ok {
		return fmt.Errorf("rollback data not of expected type 'PlaceOrderRollbackData', but %T", rollbackData)
	}

	// todo: check if ctx/session needed.. cart restore needs also be done or?
	_, err := po.orderService.CancelOrder(context.Background(), web.SessionFromContext(context.Background()), &rollbackData.order)
	if err != nil {
		return err
	}

	return nil
}

// IsFinal if state is a final state
func (po PlaceOrder) IsFinal() bool {
	return false
}
