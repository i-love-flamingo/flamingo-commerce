package events_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

func TestDefaultEventPublisher_PublishOrderPlacedEvent(t *testing.T) {
	type fields struct {
		logger flamingo.Logger
	}
	type args struct {
		ctx              context.Context
		cart             *cart.Cart
		placedOrderInfos placeorder.PlacedOrderInfos
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				ctx:  context.Background(),
				cart: &cart.Cart{},
				placedOrderInfos: placeorder.PlacedOrderInfos{
					placeorder.PlacedOrderInfo{
						OrderNumber:  "124",
						DeliveryCode: "test_delivery",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//// prepare the wantEvent for the mocked event router
			//wantEvent := &application.OrderPlacedEvent{
			//	Cart:             tt.args.cart,
			//	PlacedOrderInfos: tt.args.placedOrderInfos,
			//}
			//
			//// prepare the event router
			//eventRouter := new(mocks.Router)
			//eventRouter.On(
			//	"Dispatch",
			//	context.Background(),
			//	mock.MatchedBy(
			//		func(e flamingo.Event) bool {
			//			if diff := deep.Equal(e, wantEvent); diff != nil {
			//				t.Logf("PublishOrderPlacedEvent got!=want, diff: %#v", diff)
			//				return false
			//			}
			//
			//			return true
			//		},
			//	),
			//).Return(nil)
			//
			//// prepare the event publisher
			//dep := &application.DefaultEventPublisher{}
			//dep.Inject(
			//	tt.fields.logger,
			//	eventRouter,
			//)
			//
			//dep.PublishOrderPlacedEvent(tt.args.ctx, tt.args.cart, tt.args.placedOrderInfos)
			//eventRouter.AssertExpectations(t)
			//eventRouter.AssertNumberOfCalls(t, "Dispatch", 1)
		})
	}
}
