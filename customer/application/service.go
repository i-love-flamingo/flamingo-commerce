package application

import (
	"go.aoe.com/flamingo/core/auth/application"
	"go.aoe.com/flamingo/core/customer/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	Service struct {
		AuthManager     *application.AuthManager `inject:""`
		CustomerService domain.CustomerService   `inject:""`
	}
)

func (s *Service) GetForAuthenticatedUser(ctx web.Context) (domain.Customer, error) {
	auth, err := s.AuthManager.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return s.CustomerService.GetByAuth(auth)
}
