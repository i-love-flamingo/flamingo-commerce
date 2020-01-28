package states_test

import (
	"context"
	"testing"

	"flamingo.me/flamingo-commerce/v3/cart/application"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
	authApplication "flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/core/oauth/domain"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
)

func TestCompleteCart_IsFinal(t *testing.T) {
	assert.False(t, states.CompleteCart{}.IsFinal())
}

func TestCompleteCart_Name(t *testing.T) {
	assert.Equal(t, "CompleteCart", states.CompleteCart{}.Name())
}

func TestCompleteCart_Run(t *testing.T) {
	//factory := provideProcessFactory()
	//p, _ := factory.New(&url.URL{}, cartDomain.Cart{})

	cartReceiverService := &application.CartReceiverService{}

	//cartReceiverService.Inject()
	state := states.CompleteCart{}
	state.Inject(&application.CartService{}, cartReceiverService)

	//state.Run(context.Background(), p)

}

func TestCompleteCart_Rollback(t *testing.T) {

}

type (
	MockUserService struct {
		LoggedIn bool
	}
)

var _ authApplication.UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) GetUser(_ context.Context, _ *web.Session) *domain.User {
	return &domain.User{
		Name: "Test",
	}
}

func (m *MockUserService) IsLoggedIn(_ context.Context, _ *web.Session) bool {
	return m.LoggedIn
}
