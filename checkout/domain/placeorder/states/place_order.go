package states

import (
	"context"
	"encoding/gob"
	"fmt"

	cartApplication "flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/cart/domain/decorator"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	paymentApplication "flamingo.me/flamingo-commerce/v3/payment/application"
	"flamingo.me/flamingo/v3/framework/web"
	"go.opencensus.io/trace"
)

type (
	// PlaceOrder state
	PlaceOrder struct {
		orderService               *application.OrderService
		cartService                *cartApplication.CartService
		cartDecoratorFactory       *decorator.DecoratedCartFactory
		paymentService             *paymentApplication.PaymentService
		cancelOrdersDuringRollback bool
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
	paymentService *paymentApplication.PaymentService,
	cfg *struct {
		CancelOrdersDuringRollback bool `inject:"config:commerce.checkout.placeorder.states.placeorder.cancelOrdersDuringRollback"`
	},
) *PlaceOrder {
	po.orderService = orderService
	po.cartService = cartService
	po.cartDecoratorFactory = cartDecoratorFactory
	po.paymentService = paymentService

	if cfg != nil {
		po.cancelOrdersDuringRollback = cfg.CancelOrdersDuringRollback
	}

	return po
}

// Name get state name
func (PlaceOrder) Name() string {
	return "PlaceOrder"
}

// Run the state operations
func (po PlaceOrder) Run(ctx context.Context, p *process.Process) process.RunResult {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/PlaceOrder/Run")
	defer span.End()

	cart := p.Context().Cart
	decoratedCart := po.cartDecoratorFactory.Create(ctx, cart)

	payment := &placeorder.Payment{}
	if !cart.GrandTotal.IsZero() {
		paymentGateway, err := po.paymentService.PaymentGatewayByCart(cart)
		if err != nil {
			return process.RunResult{
				Failed: process.PaymentErrorOccurredReason{Error: err.Error()},
			}
		}

		payment, err = paymentGateway.OrderPaymentFromFlow(ctx, &cart, p.Context().UUID)
		if err != nil {
			return process.RunResult{
				Failed: process.ErrorOccurredReason{Error: err.Error()},
			}
		}

		p.UpdateState(ValidatePayment{}.Name(), nil)
	} else {
		p.UpdateState(Success{}.Name(), nil)
	}

	infos, err := po.orderService.CartPlaceOrder(ctx, decoratedCart, *payment)
	if err != nil {
		return process.RunResult{
			Failed: process.ErrorOccurredReason{Error: err.Error()},
		}
	}

	p.UpdateOrderInfo(infos)
	return process.RunResult{
		RollbackData: PlaceOrderRollbackData{OrderInfos: *infos},
	}
}

// Rollback the state operations
func (po PlaceOrder) Rollback(ctx context.Context, data process.RollbackData) error {
	ctx, span := trace.StartSpan(ctx, "placeorder/state/PlaceOrder/Rollback")
	defer span.End()

	rollbackData, ok := data.(PlaceOrderRollbackData)
	if !ok {
		return fmt.Errorf("rollback data not of expected type 'PlaceOrderRollbackData', but %T", rollbackData)
	}

	if !po.cancelOrdersDuringRollback {
		return nil
	}

	err := po.orderService.CancelOrderWithoutRestore(ctx, web.SessionFromContext(ctx), &rollbackData.OrderInfos)
	if err != nil {
		return err
	}

	return nil
}

// IsFinal if state is a final state
func (po PlaceOrder) IsFinal() bool {
	return false
}
