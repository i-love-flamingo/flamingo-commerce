package restrictors

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	productDomain "flamingo.me/flamingo-commerce/v3/product/domain"
	"flamingo.me/flamingo-commerce/v3/sourcing/application"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"go.opencensus.io/trace"
)

type (
	// Restrictor restricts qty based on available stock
	Restrictor struct {
		logger          flamingo.Logger
		sourcingService application.SourcingApplication
	}
)

var _ validation.MaxQuantityRestrictor = new(Restrictor)

// Inject dependencies
func (r *Restrictor) Inject(
	l flamingo.Logger,
	sourcingService application.SourcingApplication,
) *Restrictor {
	r.logger = l.WithField(flamingo.LogKeyCategory, "SourceAvailableRestrictor")
	r.sourcingService = sourcingService

	return r
}

// Name returns the code of the restrictor
func (r *Restrictor) Name() string {
	return "SourceAvailableRestrictor"
}

// Restrict qty based on product data
func (r *Restrictor) Restrict(ctx context.Context, session *web.Session, product productDomain.BasicProduct, _ *cart.Cart, deliveryCode string) *validation.RestrictionResult {
	ctx, span := trace.StartSpan(ctx, "sourcing/restrictors/SourceAvailableRestrictor")
	defer span.End()

	unrestricted := &validation.RestrictionResult{
		IsRestricted:        false,
		MaxAllowed:          0,
		RemainingDifference: 0,
		RestrictorName:      r.Name(),
	}

	availableSources, err := r.sourcingService.GetAvailableSources(ctx, session, product, deliveryCode)
	if err != nil {
		return unrestricted
	}

	availableSourcesDeducted, err := r.sourcingService.GetAvailableSourcesDeductedByCurrentCart(ctx, session, product, deliveryCode)
	if err != nil {
		r.logger.Error(err)

		return &validation.RestrictionResult{
			IsRestricted:        true,
			MaxAllowed:          availableSources.QtySum(),
			RemainingDifference: availableSources.QtySum(),
			RestrictorName:      r.Name(),
		}
	}

	return &validation.RestrictionResult{
		IsRestricted:        true,
		MaxAllowed:          availableSources.QtySum(),
		RemainingDifference: availableSourcesDeducted.QtySum(),
		RestrictorName:      r.Name(),
	}
}
