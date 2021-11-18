package logger_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/placeorder"
	logger "flamingo.me/flamingo-commerce/v3/cart/infrastructure/placeorder"
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"github.com/stretchr/testify/assert"
)

type (
	stubLogger struct {
		flamingo.NullLogger
		loggedcart interface{}
	}
)

func (l *stubLogger) WithField(key flamingo.LogKey, value interface{}) flamingo.Logger {
	if key == "cart" {
		l.loggedcart = value
	}
	return l
}

func TestPlaceOrderLoggerAdapter_PlaceGuestCart(t *testing.T) {
	stubLogger := &stubLogger{}
	placeOrderAdapter := &logger.PlaceOrderLoggerAdapter{}
	placeOrderAdapter.Inject(stubLogger, &struct {
		UseFlamingoLog bool   `inject:"config:commerce.cart.placeOrderLogger.useFlamingoLog,optional"`
		LogAsFile      bool   `inject:"config:commerce.cart.placeOrderLogger.logAsFile,optional"`
		LogDirectory   string `inject:"config:commerce.cart.placeOrderLogger.logDirectory,optional"`
	}{
		UseFlamingoLog: true,
		LogAsFile:      false,
		LogDirectory:   "",
	})
	exampleCart := &cart.Cart{
		ID:       "testid",
		EntityID: "",
		BillingAddress: &cart.Address{
			Vat:                    "",
			Firstname:              "Adrianna",
			Lastname:               "Mustermann",
			MiddleName:             "",
			Title:                  "",
			Salutation:             "Mr",
			Street:                 "Musterstraße",
			StreetNr:               "7",
			AdditionalAddressLines: nil,
			Company:                "AOE",
			City:                   "Wiesbaden",
			PostCode:               "65200",
			State:                  "",
			RegionCode:             "",
			Country:                "Germany",
			CountryCode:            "",
			Telephone:              "",
			Email:                  "adrianna@mail.de",
		},
		Purchaser: &cart.Person{
			Address: &cart.Address{
				Vat:                    "",
				Firstname:              "Max",
				Lastname:               "Mustermann",
				MiddleName:             "",
				Title:                  "",
				Salutation:             "Mr",
				Street:                 "Musterstraße",
				StreetNr:               "7",
				AdditionalAddressLines: nil,
				Company:                "AOE",
				City:                   "Wiesbaden",
				PostCode:               "65200",
				State:                  "",
				RegionCode:             "",
				Country:                "Germany",
				CountryCode:            "",
				Telephone:              "",
				Email:                  "mail@mail.de",
			},
			PersonalDetails:      cart.PersonalDetails{},
			ExistingCustomerData: nil,
		},
		Deliveries: []cart.Delivery{
			{
				DeliveryInfo: cart.DeliveryInfo{
					Code:     "delivery",
					Workflow: "",
					Method:   "",
					Carrier:  "",
					DeliveryLocation: cart.DeliveryLocation{
						Type: "",
						Address: &cart.Address{
							Vat:                    "",
							Firstname:              "Opa",
							Lastname:               "Mustermann",
							MiddleName:             "",
							Title:                  "",
							Salutation:             "Mr",
							Street:                 "Musterstraße",
							StreetNr:               "7",
							AdditionalAddressLines: nil,
							Company:                "AOE",
							City:                   "Wiesbaden",
							PostCode:               "65200",
							State:                  "",
							RegionCode:             "",
							Country:                "Germany",
							CountryCode:            "",
							Telephone:              "",
							Email:                  "mail@mail.de",
						},
						UseBillingAddress: false,
						Code:              "",
					},
					AdditionalData:          nil,
					AdditionalDeliveryInfos: nil,
				},
				Cartitems: []cart.Item{
					{
						ID:                     "1",
						ExternalReference:      "",
						MarketplaceCode:        "",
						VariantMarketPlaceCode: "",
						ProductName:            "ProductName",
						SourceID:               "",
						Qty:                    1,
						AdditionalData:         nil,
						SinglePriceGross:       domain.NewFromInt(1190, 100, "€"),
						SinglePriceNet:         domain.NewFromInt(1000, 100, "€"),
						RowPriceGross:          domain.NewFromInt(2380, 100, "€"),
						RowPriceNet:            domain.NewFromInt(2000, 100, "€"),
						RowTaxes:               nil,
						AppliedDiscounts:       nil,
					},
				},
				ShippingItem: cart.ShippingItem{
					Title:            "Express",
					PriceNet:         domain.NewFromInt(1000, 100, "€"),
					TaxAmount:        domain.NewFromInt(190, 100, "€"),
					AppliedDiscounts: nil,
				},
			},
			{
				DeliveryInfo: cart.DeliveryInfo{
					Code:     "pickup",
					Workflow: "pickup",
					Method:   "",
					Carrier:  "",
					DeliveryLocation: cart.DeliveryLocation{
						Type:              "pickup",
						UseBillingAddress: false,
						Code:              "location1",
					},
					AdditionalData:          nil,
					AdditionalDeliveryInfos: nil,
				},
				Cartitems: []cart.Item{
					{
						ID:                     "2",
						ExternalReference:      "",
						MarketplaceCode:        "",
						VariantMarketPlaceCode: "",
						ProductName:            "ProductName 2",
						SourceID:               "",
						Qty:                    1,
						AdditionalData:         nil,
						SinglePriceGross:       domain.NewFromInt(1190, 100, "€"),
						SinglePriceNet:         domain.NewFromInt(1000, 100, "€"),
						RowPriceGross:          domain.NewFromInt(2380, 100, "€"),
						RowPriceNet:            domain.NewFromInt(2000, 100, "€"),
						RowTaxes:               nil,
						AppliedDiscounts:       nil,
					},
					{
						ID:                     "3",
						ExternalReference:      "",
						MarketplaceCode:        "",
						VariantMarketPlaceCode: "",
						ProductName:            "ProductName 3",
						SourceID:               "",
						Qty:                    1,
						AdditionalData:         nil,
						SinglePriceGross:       domain.NewFromInt(1190, 100, "€"),
						SinglePriceNet:         domain.NewFromInt(1000, 100, "€"),
						RowPriceGross:          domain.NewFromInt(2380, 100, "€"),
						RowPriceNet:            domain.NewFromInt(2000, 100, "€"),
						RowTaxes:               nil,
						AppliedDiscounts:       nil,
					},
				},
			},
		},
		AdditionalData:             cart.AdditionalData{},
		PaymentSelection:           nil,
		BelongsToAuthenticatedUser: false,
		AuthenticatedUserID:        "",
		AppliedCouponCodes:         nil,
		DefaultCurrency:            "",
		Totalitems:                 nil,
		AppliedGiftCards:           nil,
	}

	payment := &placeorder.Payment{
		Gateway: "test",
		Transactions: []placeorder.Transaction{
			placeorder.Transaction{
				Method:            "testmethod",
				Status:            placeorder.PaymentStatusOpen,
				ValuedAmountPayed: exampleCart.GrandTotal,
				AmountPayed:       exampleCart.GrandTotal,
				TransactionID:     "t1",
			},
		},
		RawTransactionData: nil,
		PaymentID:          "p1",
	}
	poi, err := placeOrderAdapter.PlaceGuestCart(context.Background(), exampleCart, payment)
	assert.NoError(t, err)
	assert.Equal(t, poi.GetOrderNumberForDeliveryCode("delivery"), "testid")
	assert.NotNil(t, stubLogger.loggedcart)
	assert.IsType(t, stubLogger.loggedcart, &cart.Cart{})
}
