package cart_test

import (
	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"github.com/stretchr/testify/assert"

	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

func TestDefaultDeliveryInfoBuilder_BuildByDeliveryCode(t *testing.T) {
	type fields struct {
		logger flamingo.NullLogger
		config *struct {
			DefaultUseBillingAddress bool `inject:"config:commerce.cart.defaultUseBillingAddress,optional"`
		}
	}
	type args struct {
		deliverycode string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *cart.DeliveryInfo
		wantErr bool
	}{
		{
			args: args{
				deliverycode: "delivery",
			},
			wantErr: false,
			want: &cart.DeliveryInfo{
				Code:     "delivery",
				Workflow: "delivery",
				DeliveryLocation: cart.DeliveryLocation{
					Type: "unspecified",
				},
			},
			name: "test for delivery",
		},
		{
			args: args{
				deliverycode: "pickup_store",
			},
			wantErr: false,
			want: &cart.DeliveryInfo{
				Code:     "pickup_store",
				Workflow: "pickup",
				DeliveryLocation: cart.DeliveryLocation{
					Type: "store",
				},
			},
			name: "test for pickup_store",
		},
		{
			args: args{
				deliverycode: "workflow___method",
			},
			wantErr: false,
			want: &cart.DeliveryInfo{
				Code:     "workflow___method",
				Workflow: "workflow",
				Method:   "method",
				DeliveryLocation: cart.DeliveryLocation{
					Type: "unspecified",
				},
			},
			name: "test for empty type and locationdetail",
		},
		{
			args: args{
				deliverycode: "workflow____ignoreme",
			},
			wantErr: false,
			want: &cart.DeliveryInfo{
				Code:     "workflow____ignoreme",
				Workflow: "workflow",
				DeliveryLocation: cart.DeliveryLocation{
					Type: "unspecified",
				},
			},
			name: "test for empty type and locationdetail and method",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &cart.DefaultDeliveryInfoBuilder{}
			b.Inject(tt.fields.logger, tt.fields.config)
			got, err := b.BuildByDeliveryCode(tt.args.deliverycode)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultDeliveryInfoBuilder.BuildByDeliveryCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
