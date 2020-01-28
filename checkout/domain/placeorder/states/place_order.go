package states

import (
	"context"
	"encoding/gob"
	"fmt"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// PlaceOrder state
	PlaceOrder struct {
		orderService         *application.OrderService
		cartService          *cartApplication.CartService
		cartDecoratorFactory *decorator.DecoratedCartFactory
	}

	// PlaceOrderRollbackData needed for rollbacks
	PlaceOrderRollbackData struct {
		OrderInfos application.PlaceOrderInfo
	}
)

func init() {
	gob.Register(PlaceOrderRollbackData{})
}

var _ process.State = PlaceOrder{}

// Inject dependencies
func (po *PlaceOrder) Inject(
	orderService *application.OrderService,
	cartService *cartApplication.CartService,
	cartDecoratorFactory *decorator.DecoratedCartFactory,
) *PlaceOrder {
	po.orderService = orderService
	po.cartService = cartService
	po.cartDecoratorFactory = cartDecoratorFactory

	return po
}

// Name get state name
func (PlaceOrder) Name() string {
	return "PlaceOrder"
}

// Run the state operations
func (po PlaceOrder) Run(ctx context.Context, p *process.Process, stateData process.StateData) process.RunResult {
	cart := p.Context().Cart
	decoratedCart := po.cartDecoratorFactory.Create(ctx, cart)

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

	infos, err := po.orderService.CartPlaceOrder(ctx, decoratedCart, *payment)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	p.UpdateState(ValidatePayment{}.Name(), nil)
	return process.RunResult{
		RollbackData: PlaceOrderRollbackData{OrderInfos: *infos},
	}
}

// Rollback the state operations
func (po PlaceOrder) Rollback(ctx context.Context, data process.RollbackData) error {
	rollbackData, ok := data.(PlaceOrderRollbackData)
	if !ok {
		return fmt.Errorf("rollback data not of expected type 'PlaceOrderRollbackData', but %T", rollbackData)
	}

	_, err := po.orderService.CancelOrder(ctx, web.SessionFromContext(context.Background()), &rollbackData.OrderInfos)
	if err != nil {
		return err
	}

	return nil
}

// IsFinal if state is a final state
func (po PlaceOrder) IsFinal() bool {
	return false
}
