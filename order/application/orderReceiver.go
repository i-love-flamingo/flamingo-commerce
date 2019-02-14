package application

import (
	"context"
	"flamingo.me/flamingo/v3/framework/web"

	"flamingo.me/flamingo-commerce/v3/order/domain"
	authApplication "flamingo.me/flamingo/v3/core/auth/application"
	authDomain "flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// OrderReceiverService provides methods to place and get ortders
	OrderReceiverService struct {
		guestOrderService    domain.GuestOrderService
		customerOrderService domain.CustomerOrderService
		userService          authApplication.UserServiceInterface
		authManager          *authApplication.AuthManager
		logger               flamingo.Logger
	}
)

// Inject dependencies
func (ors *OrderReceiverService) Inject(
	GuestOrderService domain.GuestOrderService,
	CustomerOrderService domain.CustomerOrderService,
	UserService authApplication.UserServiceInterface,
	AuthManager *authApplication.AuthManager,
	Logger flamingo.Logger,
) {
	ors.guestOrderService = GuestOrderService
	ors.customerOrderService = CustomerOrderService
	ors.userService = UserService
	ors.authManager = AuthManager
	ors.logger = Logger
}

// GetBehaviour returns the order behaviour depending on the logged in state
func (ors *OrderReceiverService) GetBehaviour(ctx context.Context, session *web.Session) (domain.Behaviour, error) {
	if ors.userService.IsLoggedIn(ctx, session) {
		behaviour, err := ors.customerOrderService.GetBehaviour(ctx, ors.Auth(ctx, session))
		if err != nil {
			return nil, err
		}

		return behaviour, nil
	}

	behaviour, err := ors.guestOrderService.GetBehaviour(ctx)
	if err != nil {
		return nil, err
	}
	return behaviour, nil
}

// Auth tries to retrieve the authentication context for a active session
func (ors *OrderReceiverService) Auth(ctx context.Context, session *web.Session) authDomain.Auth {
	ts, _ := ors.authManager.TokenSource(ctx, session)
	idToken, _ := ors.authManager.IDToken(ctx, session)

	return authDomain.Auth{
		TokenSource: ts,
		IDToken:     idToken,
	}
}
