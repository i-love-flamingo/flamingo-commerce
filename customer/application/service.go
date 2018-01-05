package application

import (
	"go.aoe.com/flamingo/core/auth/application"
	"go.aoe.com/flamingo/core/customer/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	Service struct {
		AuthManager            application.AuthManager `inject:""`
		domain.CustomerService `inject:""`
	}
)

func (s *Service) GetForAuthenticatedUser(ctx web.Context) (domain.Customer, error) {
	auth := s.AuthManager.Auth(ctx)
	return s.CustomerService.GetByAuth(auth)
}
