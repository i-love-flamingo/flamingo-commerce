package cart

import (
	"context"
	"errors"

	domainCart "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/infrastructure"
	"flamingo.me/flamingo-commerce/v3/price/domain"
)

type (
	// FakeVoucherHandler used vouchers in integration tests
	FakeVoucherHandler struct{}
)

var (
	_ infrastructure.VoucherHandler = &FakeVoucherHandler{}
)

// ApplyVoucher fake implementation
func (f FakeVoucherHandler) ApplyVoucher(ctx context.Context, cart *domainCart.Cart, couponCode string) (*domainCart.Cart, error) {
	if couponCode != "100-percent-off" {
		err := errors.New("voucher code invalid")
		return nil, err
	}

	coupon := domainCart.CouponCode{
		Code: couponCode,
	}

	if !f.voucherAlreadyApplied(cart, couponCode) {
		cart.AppliedCouponCodes = append(cart.AppliedCouponCodes, coupon)
	}

	if couponCode == "100-percent-off" {
		for delKey, delivery := range cart.Deliveries {
			for itemKey, item := range delivery.Cartitems {
				cart.Deliveries[delKey].Cartitems[itemKey].AppliedDiscounts = []domainCart.AppliedDiscount{{
					CampaignCode:  "100-percent-off",
					CouponCode:    "100-percent-off",
					Label:         "100% Off",
					Applied:       item.RowPriceGross.Inverse(),
					Type:          "coupon",
					IsItemRelated: false,
					SortOrder:     0,
				}}
				cart.Deliveries[delKey].Cartitems[itemKey].RowPriceGrossWithDiscount = domain.NewZero(item.RowPriceGross.Currency())
				cart.Deliveries[delKey].Cartitems[itemKey].RowPriceNetWithDiscount = domain.NewZero(item.RowPriceGross.Currency())
				cart.Deliveries[delKey].Cartitems[itemKey].NonItemRelatedDiscountAmount = item.RowPriceGross.Inverse()
				cart.Deliveries[delKey].Cartitems[itemKey].TotalDiscountAmount = item.RowPriceGross.Inverse()
			}
		}
	}

	return cart, nil
}

// RemoveVoucher fake implementation
func (f FakeVoucherHandler) RemoveVoucher(ctx context.Context, cart *domainCart.Cart, couponCode string) (*domainCart.Cart, error) {
	for i, coupon := range cart.AppliedCouponCodes {
		if coupon.Code == couponCode {
			cart.AppliedCouponCodes[i] = cart.AppliedCouponCodes[len(cart.AppliedCouponCodes)-1]
			cart.AppliedCouponCodes[len(cart.AppliedCouponCodes)-1] = domainCart.CouponCode{}
			cart.AppliedCouponCodes = cart.AppliedCouponCodes[:len(cart.AppliedCouponCodes)-1]
			return cart, nil
		}
	}

	return nil, errors.New("couldn't remove supplied voucher since it wasn't applied before")
}

func (FakeVoucherHandler) voucherAlreadyApplied(cart *domainCart.Cart, code string) bool {
	for _, coupon := range cart.AppliedCouponCodes {
		if coupon.Code == code {
			return true
		}
	}
	return false
}
