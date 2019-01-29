package application_test

import (
	"context"
	"errors"
	"reflect"
	"testing"


	"flamingo.me/flamingo-commerce/order/application"
	"flamingo.me/flamingo-commerce/order/domain"
	domainMocks "flamingo.me/flamingo-commerce/order/domain/mocks"
	authApplication "flamingo.me/flamingo/core/auth/application"
	authApplicationMocks "flamingo.me/flamingo/core/auth/application/mocks"
	"flamingo.me/flamingo/framework/flamingo"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/mock"
)

func TestOrderReceiverService_GetBehaviour(t *testing.T) {
	type fields struct {
		logger flamingo.Logger
	}
	type args struct {
		ctx                         context.Context
		session                     *sessions.Session
		isLoggedIn                  bool
		hasCustomerServiceBehaviour bool
		hasGuestServiceBehaviour    bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.Behaviour
		wantErr bool
	}{
		{
			name: "not logged in, error on get guest behaviour",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				ctx:                         context.Background(),
				session:                     &sessions.Session{},
				isLoggedIn:                  false,
				hasCustomerServiceBehaviour: false,
				hasGuestServiceBehaviour:    false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not logged in, can get guest behaviour",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				ctx:                         context.Background(),
				session:                     &sessions.Session{},
				isLoggedIn:                  false,
				hasCustomerServiceBehaviour: false,
				hasGuestServiceBehaviour:    true,
			},
			want:    new(domainMocks.Behaviour),
			wantErr: false,
		},
		{
			name: "logged in, no behaviour",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				ctx:                         context.Background(),
				session:                     &sessions.Session{},
				isLoggedIn:                  true,
				hasCustomerServiceBehaviour: false,
				hasGuestServiceBehaviour:    false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "logged in, with behaviour",
			fields: fields{
				logger: flamingo.NullLogger{},
			},
			args: args{
				ctx:                         context.Background(),
				session:                     &sessions.Session{},
				isLoggedIn:                  true,
				hasCustomerServiceBehaviour: true,
				hasGuestServiceBehaviour:    false,
			},
			want:    new(domainMocks.Behaviour),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userServiceMock := new(authApplicationMocks.UserServiceInterface)
			userServiceMock.On("IsLoggedIn", mock.Anything, mock.Anything).Return(tt.args.isLoggedIn)

			guestOrderServiceMock := new(domainMocks.GuestOrderService)
			if tt.args.hasGuestServiceBehaviour {
				guestOrderServiceMock.On("GetBehaviour", mock.Anything).Return(tt.want, nil)
			} else {
				guestOrderServiceMock.On("GetBehaviour", mock.Anything).Return(tt.want, errors.New("nope"))
			}

			customerOrderServiceMock := new(domainMocks.CustomerOrderService)
			if tt.args.hasCustomerServiceBehaviour {
				customerOrderServiceMock.On("GetBehaviour", mock.Anything, mock.Anything).Return(tt.want, nil)
			} else {
				customerOrderServiceMock.On("GetBehaviour", mock.Anything, mock.Anything).Return(tt.want, errors.New("nope"))
			}

			authManager := &authApplication.AuthManager{}

			ors := &application.OrderReceiverService{}
			ors.Inject(
				guestOrderServiceMock,
				customerOrderServiceMock,
				userServiceMock,
				authManager,
				tt.fields.logger,
			)
			got, err := ors.GetBehaviour(tt.args.ctx, tt.args.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderReceiverService.GetBehaviour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderReceiverService.GetBehaviour() = %v, want %v", got, tt.want)
			}

			userServiceMock.AssertExpectations(t)
			if !tt.args.isLoggedIn {
				guestOrderServiceMock.AssertExpectations(t)
			} else {
				customerOrderServiceMock.AssertExpectations(t)
			}
		})
	}
}
