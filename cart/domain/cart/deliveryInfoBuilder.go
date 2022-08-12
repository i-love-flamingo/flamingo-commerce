package cart

import (
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (

	// DeliveryInfoBuilder can be used to set delivery infos depending on supplied delivery code
	DeliveryInfoBuilder interface {
		BuildByDeliveryCode(deliveryCode string) (*DeliveryInfo, error)
	}

	// DefaultDeliveryInfoBuilder defines the default delivery info builder used
	DefaultDeliveryInfoBuilder struct {
		logger                   flamingo.Logger
		defaultUseBillingAddress bool
	}
)

// Inject dependencies
func (b *DefaultDeliveryInfoBuilder) Inject(
	logger flamingo.Logger,
	config *struct {
		DefaultUseBillingAddress bool `inject:"config:commerce.cart.defaultUseBillingAddress,optional"`
	},
) {
	b.logger = logger.WithField("category", "cart").WithField("subcategory", "DefaultDeliveryInfoBuilder")
	if config != nil {
		b.defaultUseBillingAddress = config.DefaultUseBillingAddress
	}
}

// BuildByDeliveryCode builds a (initial) DeliveryInfo by deliveryCode
// Convention that is used in this factory is to split infos in the build deliveryinfo by "_" like this:
//   - workflow_locationtype_locationdetail_method_anythingelse
//   - not all parts are required
//   - to "skip" parts in between use "-"
func (b *DefaultDeliveryInfoBuilder) BuildByDeliveryCode(deliverycode string) (*DeliveryInfo, error) {
	if deliverycode == "" {
		b.logger.Warn("Empty deliverycode")
	}

	intentParts := strings.SplitN(deliverycode, "_", 5)

	deliveryInfo := DeliveryInfo{
		Code: deliverycode,
	}
	if len(intentParts) > 0 && intentParts[0] != "" {
		deliveryInfo.Workflow = intentParts[0]
	} else {
		deliveryInfo.Workflow = DeliveryWorkflowUnspecified
	}

	if len(intentParts) > 1 && intentParts[1] != "" {
		deliveryInfo.DeliveryLocation.Type = intentParts[1]
	} else {
		deliveryInfo.DeliveryLocation.Type = DeliverylocationTypeUnspecified
	}
	if len(intentParts) > 2 && intentParts[2] != "" {
		deliveryInfo.DeliveryLocation.Code = intentParts[2]
	}
	if len(intentParts) > 3 && intentParts[3] != "" {
		deliveryInfo.Method = intentParts[3]
	}

	deliveryInfo.DeliveryLocation.UseBillingAddress = b.defaultUseBillingAddress

	return &deliveryInfo, nil
}
