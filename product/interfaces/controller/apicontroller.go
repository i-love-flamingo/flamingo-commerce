package controller

import (
	"context"

	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"

	"flamingo.me/flamingo-commerce/v3/product/application"
	"flamingo.me/flamingo-commerce/v3/product/domain"
)

type (
	// APIController for products
	APIController struct {
		responder      *web.Responder
		productService domain.ProductService
		uRLService     *application.URLService
	}

	// APIResult view data
	APIResult struct {
		Error   *resultError
		Success bool
		Product domain.BasicProduct
	}

	resultError struct {
		Message string
		Code    string
	} //@name productResultError
)

// Inject dependencies
func (c *APIController) Inject(responder *web.Responder,
	productService domain.ProductService,
	uRLService *application.URLService) *APIController {
	c.responder = responder
	c.productService = productService
	c.uRLService = uRLService
	return c
}

// Get Response for Product matching marketplacecode param
// @Summary Gets the requested product
// @Tags  Product
// @Produce json
// @Success 200 {object} APIResult{product=domain.SimpleProduct}
// @Failure 500 {object} APIResult
// @Failure 404 {object} APIResult
// @Param marketplacecode path string true "the marketplace code (idendifier) for the product"
// @Router /api/v1/products/{marketplacecode} [get]
func (c *APIController) Get(ctx context.Context, r *web.Request) web.Result {
	product, err := c.productService.Get(ctx, r.Params["marketplacecode"])
	if err != nil {
		switch errors.Cause(err).(type) {
		case domain.ProductNotFound:
			return c.responder.Data(APIResult{
				Success: false,
				Error:   &resultError{Code: "404", Message: err.Error()},
			})

		default:
			return c.responder.Data(APIResult{
				Success: false,
				Error:   &resultError{Code: "500", Message: err.Error()},
			})
		}
	}

	return c.responder.Data(APIResult{
		Success: true,
		Product: product,
	})
}
