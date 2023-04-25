package restrictors

import (
	"context"
	"errors"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/domain/cart"
	"flamingo.me/flamingo-commerce/v3/cart/domain/validation"
	"flamingo.me/flamingo-commerce/v3/product/domain"
	sourcingApplication "flamingo.me/flamingo-commerce/v3/sourcing/application"
	sourcingDomain "flamingo.me/flamingo-commerce/v3/sourcing/domain"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/stretchr/testify/assert"
)

type (
	sourcingServiceMock struct {
		AvailableSources              sourcingDomain.AvailableSourcesPerProduct
		AvailableSourcesError         error
		DeductedAvailableSources      sourcingDomain.AvailableSourcesPerProduct
		DeductedAvailableSourcesError error
	}
)

func (s *sourcingServiceMock) GetAvailableSourcesDeductedByCurrentCart(_ context.Context, _ *web.Session, _ domain.BasicProduct, _ string) (sourcingDomain.AvailableSourcesPerProduct, error) {
	return s.DeductedAvailableSources, s.DeductedAvailableSourcesError
}

func (s *sourcingServiceMock) GetAvailableSources(_ context.Context, _ *web.Session, _ domain.BasicProduct, _ string) (sourcingDomain.AvailableSourcesPerProduct, error) {
	return s.AvailableSources, s.AvailableSourcesError
}

var (
	_ sourcingApplication.SourcingApplication = new(sourcingServiceMock)
)

func TestRestrictor_Restrict(t *testing.T) {
	fixtureProduct := domain.SimpleProduct{Identifier: "productid"}

	fixtureCart := &cart.Cart{}

	t.Run("error handing on error fetching available sources", func(t *testing.T) {
		want := &validation.RestrictionResult{
			IsRestricted:        false,
			MaxAllowed:          0,
			RemainingDifference: 0,
			RestrictorName:      "SourceAvailableRestrictor",
		}

		restrictor := &Restrictor{}
		restrictor.Inject(
			flamingo.NullLogger{},
			&sourcingServiceMock{
				AvailableSources:         nil,
				DeductedAvailableSources: nil,
				AvailableSourcesError:    errors.New("mocked available sources error"),
			},
		)

		got := restrictor.Restrict(
			context.Background(),
			web.EmptySession(),
			fixtureProduct,
			fixtureCart,
			"test",
		)

		assert.Equal(t, got, want)
	})

	t.Run("available sources were fetched but deduction failed, returning full source stock", func(t *testing.T) {
		want := &validation.RestrictionResult{
			IsRestricted:        true,
			MaxAllowed:          3,
			RemainingDifference: 3,
			RestrictorName:      "SourceAvailableRestrictor",
		}

		restrictor := &Restrictor{}
		restrictor.Inject(
			flamingo.NullLogger{},
			&sourcingServiceMock{
				AvailableSources: sourcingDomain.AvailableSourcesPerProduct{
					"productid": sourcingDomain.AvailableSources{
						sourcingDomain.Source{
							LocationCode:         "testCode1",
							ExternalLocationCode: "testExternalLocation1",
						}: 3,
					},
				},
				DeductedAvailableSources:      nil,
				AvailableSourcesError:         nil,
				DeductedAvailableSourcesError: errors.New("mocked available sources error"),
			},
		)

		got := restrictor.Restrict(
			context.Background(),
			web.EmptySession(),
			fixtureProduct,
			fixtureCart,
			"test",
		)

		assert.Equal(t, got, want)
	})

	t.Run("returning deducted source normally", func(t *testing.T) {
		want := &validation.RestrictionResult{
			IsRestricted:        true,
			MaxAllowed:          5,
			RemainingDifference: 3,
			RestrictorName:      "SourceAvailableRestrictor",
		}

		restrictor := &Restrictor{}
		restrictor.Inject(
			flamingo.NullLogger{},
			&sourcingServiceMock{
				AvailableSources: sourcingDomain.AvailableSourcesPerProduct{
					"productid": sourcingDomain.AvailableSources{
						sourcingDomain.Source{
							LocationCode:         "testCode1",
							ExternalLocationCode: "testExternalLocation1",
						}: 3,
						sourcingDomain.Source{
							LocationCode:         "testCode2",
							ExternalLocationCode: "testExternalLocation1",
						}: 2,
					},
				},
				DeductedAvailableSources: sourcingDomain.AvailableSourcesPerProduct{
					"productid": sourcingDomain.AvailableSources{
						sourcingDomain.Source{
							LocationCode:         "testCode1",
							ExternalLocationCode: "testExternalLocation",
						}: 2,
						sourcingDomain.Source{
							LocationCode:         "testCode2",
							ExternalLocationCode: "testExternalLocation",
						}: 1,
					},
				},
				AvailableSourcesError:         nil,
				DeductedAvailableSourcesError: nil,
			},
		)

		got := restrictor.Restrict(
			context.Background(),
			web.EmptySession(),
			fixtureProduct,
			fixtureCart,
			"test",
		)

		assert.Equal(t, got, want)
	})
}
