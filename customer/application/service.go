package application

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
)

var (
	// ErrNoIdentity user is considered to be not logged in
	ErrNoIdentity = errors.New("no identity")
)

// Service for customer management
type Service struct {
	customerIdentityService domain.CustomerIdentityService

	webIdentityService *auth.WebIdentityService
}

// Inject dependencies
func (s *Service) Inject(
	webIdentityService *auth.WebIdentityService,
	customerIdentityService domain.CustomerIdentityService,
) *Service {
	s.webIdentityService = webIdentityService
	s.customerIdentityService = customerIdentityService

	return s
}

// GetForIdentity returns the authenticated user if logged in
func (s *Service) GetForIdentity(ctx context.Context, request *web.Request) (domain.Customer, error) {
	identity := s.webIdentityService.Identify(ctx, request)
	if identity == nil {
		return nil, ErrNoIdentity
	}

	return s.customerIdentityService.GetByIdentity(ctx, identity)
}

// GetUserID returns the current user ID if logged in
//
// Can be used to check if user is logged in. Returns ErrNoIdentity if user is considered to be not logged in.
func (s *Service) GetUserID(ctx context.Context, request *web.Request) (string, error) {
	identity := s.webIdentityService.Identify(ctx, request)
	if identity != nil {
		return identity.Subject(), nil
	}

	return "", ErrNoIdentity
}
