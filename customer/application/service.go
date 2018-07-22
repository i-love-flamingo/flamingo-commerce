package application

import (
	"flamingo.me/flamingo-commerce/customer/domain"
	"flamingo.me/flamingo/core/auth"
	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/framework/web"
)

type (
	Service struct {
		AuthManager     *application.AuthManager `inject:""`
		CustomerService domain.CustomerService   `inject:""`
	}
)

func (s *Service) GetForAuthenticatedUser(ctx web.Context) (domain.Customer, error) {
	auth, err := s.AuthManager.Auth(auth.CtxSession(ctx))
	if err != nil {
		return nil, err
	}
	return s.CustomerService.GetByAuth(ctx, auth)
}
