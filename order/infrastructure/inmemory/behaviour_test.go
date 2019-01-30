package inmemory_test

import (
	"context"
	"reflect"
	"testing"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/order/domain"
	"flamingo.me/flamingo-commerce/order/infrastructure/inmemory"
)

func TestBehaviour_PlaceOrder(t *testing.T) {
	type fields struct {
		storage inmemory.Storager
	}
	type args struct {
		ctx     context.Context
		cart    *cart.Cart
		payment *cart.CartPayment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.PlacedOrderInfos
		wantErr bool
	}{
		{
			name: "empty nil cart, nil payment",
			fields: fields{
				storage: &inmemory.Storage{},
			},
			args: args{
				ctx:     context.Background(),
				cart:    nil,
				payment: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cart, nil payment",
			fields: fields{
				storage: &inmemory.Storage{},
			},
			args: args{
				ctx:     context.Background(),
				cart:    &cart.Cart{},
				payment: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil cart, payment",
			fields: fields{
				storage: &inmemory.Storage{},
			},
			args: args{
				ctx:     context.Background(),
				cart:    nil,
				payment: &cart.CartPayment{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cart empty",
			fields: fields{
				storage: &inmemory.Storage{},
			},
			args: args{
				ctx:     context.Background(),
				cart:    &cart.Cart{},
				payment: &cart.CartPayment{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "single delivery cart",
			fields: fields{
				storage: &inmemory.Storage{},
			},
			args: args{
				ctx: context.Background(),
				cart: &cart.Cart{
					ID:            "test_cart",
					EntityID:      "1",
					CartTotals:    cart.CartTotals{},
					BillingAdress: cart.Address{},
					Purchaser:     cart.Person{},
					Deliveries: []cart.Delivery{
						cart.Delivery{
							DeliveryInfo:   cart.DeliveryInfo{
								Code: "delivery",
								Method: "method",
								Carrier: "carrier",
							},
							Cartitems:      []cart.Item{
								cart.Item{
									ID: "item_1",
									UniqueId: "item_1",
									MarketplaceCode: "item_1",
									VariantMarketPlaceCode: "",
									ProductName: "Item 1",
									SourceId: "somewhere",
									Qty: 3,
									CurrencyCode: "EUR",
									SinglePrice: 1.99,
									SinglePriceInclTax: 1.99,
									RowTotal: 3 * 1.99,
									TaxAmount: 0,
									RowTotalInclTax: 3 * 1.99,
									TotalDiscountAmount: 0,
									ItemRelatedDiscountAmount: 0,
									NonItemRelatedDiscountAmount: 0,
									RowTotalWithItemRelatedDiscount: 3 * 1.99,
									RowTotalWithItemRelatedDiscountInclTax: 3 * 1.99,
									RowTotalWithDiscountInclTax: 3 * 1.99,
								},
							},
							DeliveryTotals: cart.DeliveryTotals{},
						},
					},
					BelongsToAuthenticatedUser: false,
				},
				payment: &cart.CartPayment{},
			},
			want:    domain.PlacedOrderInfos{
				domain.PlacedOrderInfo{
					OrderNumber: "1",
					DeliveryCode: "delivery",
				},
			},
			wantErr: false,
		},
		{
			name: "multi delivery cart",
			fields: fields{
				storage: &inmemory.Storage{},
			},
			args: args{
				ctx: context.Background(),
				cart: &cart.Cart{
					ID:            "test_cart",
					EntityID:      "1",
					CartTotals:    cart.CartTotals{},
					BillingAdress: cart.Address{},
					Purchaser:     cart.Person{},
					Deliveries: []cart.Delivery{
						cart.Delivery{
							DeliveryInfo:   cart.DeliveryInfo{
								Code: "delivery1",
								Method: "method",
								Carrier: "carrier",
							},
							Cartitems:      []cart.Item{
								cart.Item{
									ID: "item_1",
									UniqueId: "item_1",
									MarketplaceCode: "item_1",
									VariantMarketPlaceCode: "",
									ProductName: "Item 1",
									SourceId: "somewhere",
									Qty: 3,
									CurrencyCode: "EUR",
									SinglePrice: 1.99,
									SinglePriceInclTax: 1.99,
									RowTotal: 3 * 1.99,
									TaxAmount: 0,
									RowTotalInclTax: 3 * 1.99,
									TotalDiscountAmount: 0,
									ItemRelatedDiscountAmount: 0,
									NonItemRelatedDiscountAmount: 0,
									RowTotalWithItemRelatedDiscount: 3 * 1.99,
									RowTotalWithItemRelatedDiscountInclTax: 3 * 1.99,
									RowTotalWithDiscountInclTax: 3 * 1.99,
								},
							},
							DeliveryTotals: cart.DeliveryTotals{},
						},
						cart.Delivery{
							DeliveryInfo:   cart.DeliveryInfo{
								Code: "delivery2",
								Method: "method",
								Carrier: "carrier",
							},
							Cartitems:      []cart.Item{
								cart.Item{
									ID: "item_1",
									UniqueId: "item_1",
									MarketplaceCode: "item_1",
									VariantMarketPlaceCode: "",
									ProductName: "Item 1",
									SourceId: "somewhere",
									Qty: 3,
									CurrencyCode: "EUR",
									SinglePrice: 1.99,
									SinglePriceInclTax: 1.99,
									RowTotal: 3 * 1.99,
									TaxAmount: 0,
									RowTotalInclTax: 3 * 1.99,
									TotalDiscountAmount: 0,
									ItemRelatedDiscountAmount: 0,
									NonItemRelatedDiscountAmount: 0,
									RowTotalWithItemRelatedDiscount: 3 * 1.99,
									RowTotalWithItemRelatedDiscountInclTax: 3 * 1.99,
									RowTotalWithDiscountInclTax: 3 * 1.99,
								},
							},
							DeliveryTotals: cart.DeliveryTotals{},
						},
					},
					BelongsToAuthenticatedUser: false,
				},
				payment: &cart.CartPayment{},
			},
			want:    domain.PlacedOrderInfos{
				domain.PlacedOrderInfo{
					OrderNumber: "1",
					DeliveryCode: "delivery1",
				},
				domain.PlacedOrderInfo{
					OrderNumber: "2",
					DeliveryCode: "delivery2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &inmemory.Behaviour{}
			b.Inject(tt.fields.storage)
			got, err := b.PlaceOrder(tt.args.ctx, tt.args.cart, tt.args.payment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Behaviour.PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Behaviour.PlaceOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_HasOrder(t *testing.T) {
	type fields struct {
		orders []*domain.Order
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "no orders in storge",
			fields: fields{
				orders: nil,
			},
			args: args{
				id: "1234",
			},
			want: false,
		},
		{
			name: "order not in storage",
			fields: fields{
				orders: []*domain.Order{
					&domain.Order{
						ID: "0",
					},
				},
			},
			args: args{
				id: "1234",
			},
			want: false,
		},
		{
			name: "order in storage",
			fields: fields{
				orders: []*domain.Order{
					&domain.Order{
						ID: "1234",
					},
				},
			},
			args: args{
				id: "1234",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &inmemory.Storage{}
			for _, order := range tt.fields.orders {
				os.StoreOrder(order)
			}
			if got := os.HasOrder(tt.args.id); got != tt.want {
				t.Errorf("Storage.HasOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetOrder(t *testing.T) {
	type fields struct {
		orders []*domain.Order
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Order
		wantErr bool
	}{
		{
			name: "no orders in storage",
			fields: fields{
				orders: nil,
			},
			args: args{
				id: "1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "orders not in storage",
			fields: fields{
				orders: []*domain.Order{
					&domain.Order{
						ID: "0",
					},
				},
			},
			args: args{
				id: "1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "orders in storage",
			fields: fields{
				orders: []*domain.Order{
					&domain.Order{
						ID: "1234",
					},
				},
			},
			args: args{
				id: "1234",
			},
			want: &domain.Order{
				ID: "1234",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &inmemory.Storage{}
			for _, order := range tt.fields.orders {
				os.StoreOrder(order)
			}
			got, err := os.GetOrder(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.GetOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
