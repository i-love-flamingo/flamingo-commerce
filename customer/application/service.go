package application

import (
	"context"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/web"
)

// Service for customer management
type Service struct {
	AuthManager     *application.AuthManager `inject:""`
	CustomerService domain.CustomerService   `inject:""`
}

// GetForAuthenticatedUser returns the authenticated user
func (s *Service) GetForAuthenticatedUser(ctx context.Context, session *web.Session) (domain.Customer, error) {
	auth, err := s.AuthManager.Auth(ctx, session)
	if err != nil {
		return nil, err
	}
	return s.CustomerService.GetByAuth(ctx, auth)
}
