package application_test

import (
	"context"
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/application"
	"flamingo.me/flamingo-commerce/v3/payment/interfaces"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

func TestOrderService_LastPlacedOrder(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *application.PlaceOrderInfo
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: contextWithPlaceOrderInfoInSession(application.PlaceOrderInfo{
					PaymentInfos: nil,
					PlacedOrders: nil,
					ContactEmail: "test@test.de",
				}),
			},
			want: &application.PlaceOrderInfo{
				PaymentInfos: nil,
				PlacedOrders: nil,
				ContactEmail: "test@test.de",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGatewayProvider := func() map[string]interfaces.WebCartPaymentGateway {
				return map[string]interfaces.WebCartPaymentGateway{}
			}

			os := &application.OrderService{}
			os.Inject(nil, flamingo.NullLogger{}, nil, nil, nil, fakeGatewayProvider, nil, nil)
			got, err := os.LastPlacedOrder(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("LastPlacedOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LastPlacedOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_ClearLastPlacedOrder(t *testing.T) {
	fakeGatewayProvider := func() map[string]interfaces.WebCartPaymentGateway {
		return map[string]interfaces.WebCartPaymentGateway{}
	}

	os := &application.OrderService{}
	os.Inject(nil, flamingo.NullLogger{}, nil, nil, nil, fakeGatewayProvider, nil, nil)

	want := &application.PlaceOrderInfo{
		PaymentInfos: nil,
		PlacedOrders: nil,
		ContactEmail: "test@test.de",
	}

	ctx := contextWithPlaceOrderInfoInSession(application.PlaceOrderInfo{
		PaymentInfos: nil,
		PlacedOrders: nil,
		ContactEmail: "test@test.de",
	})

	got, err := os.LastPlacedOrder(ctx)
	if err != nil {
		t.Errorf("LastPlacedOrder() shouldn't return an error: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("LastPlacedOrder() got = %v, want %v", got, want)
	}

	os.ClearLastPlacedOrder(ctx)

	got, err = os.LastPlacedOrder(ctx)
	if err != nil {
		t.Errorf("LastPlacedOrder() shouldn't return an error: %v", err)
	}

	if got != nil {
		t.Errorf("LastPlacedOrder() shouldn't return an order")
	}

}

func contextWithPlaceOrderInfoInSession(info application.PlaceOrderInfo) context.Context {
	session := web.EmptySession()
	session.Store(application.LastPlacedOrderSessionKey, info)
	return web.ContextWithSession(context.Background(), session)
}
