package application

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/core/auth"
	"flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
)

var (
	// ErrNoIdentity user is considered to be not logged in
	ErrNoIdentity = errors.New("no identity")
)

// Service for customer management
type Service struct {
	AuthManager             *application.AuthManager
	CustomerService         domain.CustomerService
	customerIdentityService domain.CustomerIdentityService

	webIdentityService *auth.WebIdentityService
}

// Inject dependencies
func (s *Service) Inject(
	authManager *application.AuthManager,
	customerService domain.CustomerService,
	webIdentityService *auth.WebIdentityService,
	customerIdentityService domain.CustomerIdentityService,
) *Service {
	s.AuthManager = authManager
	s.CustomerService = customerService
	s.webIdentityService = webIdentityService
	s.customerIdentityService = customerIdentityService

	return s
}

// GetForAuthenticatedUser returns the authenticated user
func (s *Service) GetForAuthenticatedUser(ctx context.Context, session *web.Session) (domain.Customer, error) {
	userAuth, err := s.AuthManager.Auth(ctx, session)
	if err != nil {
		return nil, err
	}
	return s.CustomerService.GetByAuth(ctx, userAuth)
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
