package dto

import (
	cartDomain "flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/interfaces/controller/forms"

	"time"

	formDomain "flamingo.me/form/domain"
)

type (
	// BillingAddressForm is the GraphQL representation of the billing form
	BillingAddressForm struct {
		FormData       forms.AddressForm
		Processed      bool
		ValidationInfo ValidationInfo
	}

	// DeliveryAddressForm is the GraphQL representation of the delivery form
	DeliveryAddressForm struct {
		FormData          forms.AddressForm
		Processed         bool
		ValidationInfo    ValidationInfo
		UseBillingAddress bool
		DeliveryCode      string
		Method            string
		Carrier           string
		DesiredTime       time.Time
	}

	// DeliveryShippingOption used to update shipping method/carrier for a specific delivery
	DeliveryShippingOption struct {
		DeliveryCode string
		Method       string
		Carrier      string
	}

	// ValidationInfo contains form related validation information
	ValidationInfo struct {
		GeneralErrors []formDomain.Error
		FieldErrors   []FieldError
	}

	// FieldError contains field related errors
	FieldError struct {
		// MessageKey - a key of the error message. Often used to pass to translation func in the template
		MessageKey string
		// DefaultLabel - a speaking error label. OFten used to show to end user - in case no translation exists
		DefaultLabel string
		//FieldName
		FieldName string
	}

	// SelectedPaymentResult represents the selected payment
	SelectedPaymentResult struct {
		Processed      bool
		ValidationInfo ValidationInfo
	}

	// UpdateShippingOptionsResult definition
	UpdateShippingOptionsResult struct {
		Processed bool
	}

	AddToCart struct {
		MarketplaceCode        string
		Qty                    int
		DeliveryCode           string
		VariantMarketplaceCode string
		BundleConfiguration    []ChoiceConfiguration
	}

	ChoiceConfiguration struct {
		Identifier             string
		MarketplaceCode        string
		VariantMarketplaceCode *string
		Qty                    *int
	}
)

func MapBundleConfigToDomain(graphqlBundleConfig []ChoiceConfiguration) cartDomain.BundleConfiguration {
	cartBundleConfiguration := make(map[cartDomain.ChoiceID]cartDomain.ChoiceConfiguration)

	for _, configuration := range graphqlBundleConfig {
		variantMarketplaceCode := ""
		quantity := 0

		if configuration.VariantMarketplaceCode != nil {
			variantMarketplaceCode = *configuration.VariantMarketplaceCode
		}
		if configuration.Qty != nil {
			quantity = *configuration.Qty
		}

		cartBundleConfiguration[cartDomain.ChoiceID(configuration.Identifier)] = cartDomain.ChoiceConfiguration{
			MarketplaceCode:        configuration.MarketplaceCode,
			VariantMarketplaceCode: variantMarketplaceCode,
			Qty:                    quantity,
		}
	}

	return cartBundleConfiguration
}
