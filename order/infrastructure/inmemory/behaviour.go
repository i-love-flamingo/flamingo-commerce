package inmemory

import (
	"context"
	"errors"
	"fmt"

	"flamingo.me/flamingo-commerce/cart/domain/cart"
	"flamingo.me/flamingo-commerce/order/domain"
)

type (
	// Behaviour defines the in memory order behaviour
	Behaviour struct {
		storage Storager
	}

	// Storager interface for in memory order storage
	Storager interface {
		GetNextOrderID() string
		GetOrder(id string) (*domain.Order, error)
		HasOrder(id string) bool
		StoreOrder(order *domain.Order) error
	}

	// Storage as default implementation of an order storage
	Storage struct {
		orders map[string]*domain.Order
	}
)

var (
	_ domain.Behaviour = (*Behaviour)(nil)
)

// Inject dependencies
func (b *Behaviour) Inject(
	OrderStorage Storager,
) {
	b.storage = OrderStorage
}

// PlaceOrder handles the in memory order service
func (b *Behaviour) PlaceOrder(ctx context.Context, cart *cart.Cart, payment *cart.CartPayment) (domain.PlacedOrderInfos, error) {
	if cart == nil {
		return nil, errors.New("cart is nil")
	}

	if payment == nil {
		return nil, errors.New("payment is nil")
	}

	if len(cart.Deliveries) == 0 {
		return nil, errors.New("no deliveries in cart")
	}

	result := domain.PlacedOrderInfos{}
	for _, delivery := range cart.Deliveries {
		info, err := b.placeOrder(ctx, cart, &delivery)
		if err != nil {
			return nil, err
		}

		result = append(result, *info)
	}

	return result, nil
}

func (b *Behaviour) placeOrder(ctx context.Context, c *cart.Cart, d *cart.Delivery) (*domain.PlacedOrderInfo, error) {
	if len(d.Cartitems) == 0 {
		return nil, errors.New("no items in delivery")
	}

	o := &domain.Order{
		ID: b.storage.GetNextOrderID(),
		OrderItems: func(d *cart.Delivery) []*domain.OrderItem {
			result := make([]*domain.OrderItem, len(d.Cartitems))
			for i, item := range d.Cartitems {
				result[i] = &domain.OrderItem{
					MarketplaceCode:        item.MarketplaceCode,
					VariantMarketplaceCode: item.VariantMarketPlaceCode,
					Qty:                    float64(item.Qty),
					CurrencyCode:           item.CurrencyCode,
					SinglePrice:            item.SinglePrice,
					SinglePriceInclTax:     item.SinglePriceInclTax,
					RowTotal:               item.RowTotal,
					TaxAmount:              item.TaxAmount,
					RowTotalInclTax:        item.RowTotalInclTax,
					Name:                   item.ProductName,
					Price:                  item.SinglePrice,
					PriceInclTax:           item.SinglePriceInclTax,
					SourceID:               item.SourceId,
				}
			}
			return result
		}(d),
		Status: "pending",
		Total: func(d *cart.Delivery) float64 {
			var total float64
			for _, item := range d.Cartitems {
				total += item.RowTotalWithDiscountInclTax
			}

			return total
		}(d),
		CurrencyCode: "EUR",
	}

	err := b.storage.StoreOrder(o)
	if err != nil {
		return nil, err
	}

	result := &domain.PlacedOrderInfo{
		OrderNumber:  o.ID,
		DeliveryCode: d.DeliveryInfo.Code,
	}

	return result, nil
}

var (
	_ Storager = (*Storage)(nil)
)

// init the in memory order storage if required
func (os *Storage) init() {
	if os.orders == nil {
		os.orders = make(map[string]*domain.Order)
	}
}

// GetNextOrderID retuns the next possible order id
func (os *Storage) GetNextOrderID() string {
	return fmt.Sprintf("%v", len(os.orders)+1)
}

// GetOrder gets an order from the storage
func (os *Storage) GetOrder(id string) (*domain.Order, error) {
	os.init()
	if !os.HasOrder(id) {
		return nil, errors.New("no such order")
	}

	result := os.orders[id]

	return result, nil
}

// HasOrder checks if an order with `id` is in the storage
func (os *Storage) HasOrder(id string) bool {
	os.init()
	_, result := os.orders[id]

	return result
}

// StoreOrder puts an order into the in memory order storage
func (os *Storage) StoreOrder(order *domain.Order) error {
	os.init()
	os.orders[order.ID] = order

	return nil
}
