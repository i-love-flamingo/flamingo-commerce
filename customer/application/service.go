package application

import (
	"context"

	"flamingo.me/flamingo-commerce/customer/domain"
	"flamingo.me/flamingo/core/auth/application"
	"github.com/gorilla/sessions"
)

type (
	Service struct {
		AuthManager     *application.AuthManager `inject:""`
		CustomerService domain.CustomerService   `inject:""`
	}
)

func (s *Service) GetForAuthenticatedUser(ctx context.Context, session *sessions.Session) (domain.Customer, error) {
	auth, err := s.AuthManager.Auth(ctx, session)
	if err != nil {
		return nil, err
	}
	return s.CustomerService.GetByAuth(ctx, auth)
}
