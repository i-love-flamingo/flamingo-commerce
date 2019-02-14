package domain_test

import (
	"testing"

	"flamingo.me/flamingo-commerce/v3/order/domain"
)

func TestPlacedOrderInfos_GetOrderNumberForDeliveryCode(t *testing.T) {
	type args struct {
		deliveryCode string
	}
	tests := []struct {
		name string
		poi  domain.PlacedOrderInfos
		args args
		want string
	}{
		{
			name: "empty order infos, empty delivery code",
		},
		{
			name: "empty order infos, with code",
			args: args{
				deliveryCode: "delivery",
			},
			want: "",
		},
		{
			name: "delivery code not in placed orders",
			poi: domain.PlacedOrderInfos{
				domain.PlacedOrderInfo{
					OrderNumber:  "1",
					DeliveryCode: "delivery_1",
				},
				domain.PlacedOrderInfo{
					OrderNumber:  "2",
					DeliveryCode: "delivery_2",
				},
			},
			args: args{
				deliveryCode: "delivery",
			},
			want: "",
		},
		{
			name: "delivery code in placed orders",
			poi: domain.PlacedOrderInfos{
				domain.PlacedOrderInfo{
					OrderNumber:  "1",
					DeliveryCode: "delivery_1",
				},
				domain.PlacedOrderInfo{
					OrderNumber:  "2",
					DeliveryCode: "delivery_2",
				},
			},
			args: args{
				deliveryCode: "delivery_1",
			},
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.poi.GetOrderNumberForDeliveryCode(tt.args.deliveryCode); got != tt.want {
				t.Errorf("PlacedOrderInfos.GetOrderNumberForDeliveryCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
