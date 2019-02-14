package formDto_test

import (
	"context"
	"net/http"
	"testing"

	"flamingo.me/flamingo-commerce/v3/checkout/interfaces/controller/formDto"
	"flamingo.me/flamingo-commerce/v3/customer/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/go-playground/form"
)

type (
	FakeCustomer struct{}
)

func (f *FakeCustomer) GetPersonalData() domain.PersonData {
	return domain.PersonData{}
}

func (f *FakeCustomer) GetAddresses() []domain.Address {
	return nil
}
func (f *FakeCustomer) GetDefaultShippingAddress() *domain.Address {
	return &domain.Address{
		Firstname: "first",
		Email:     "example@mail.dmn",
	}
}
func (f *FakeCustomer) GetDefaultBillingAddress() *domain.Address {
	return &domain.Address{
		Firstname: "first",
		Email:     "example@mail.dmn",
	}
}

func (f *FakeCustomer) GetId() string {
	return "customerID32929"
}

var (
	_ domain.Customer = &FakeCustomer{}
)

func TestFormService(t *testing.T) {
	service := formDto.CheckoutFormService{}
	service.Inject(form.NewDecoder(), flamingo.NullLogger{}, nil)
	service.SetCustomer(&FakeCustomer{})

	r := web.CreateRequest(&http.Request{}, nil)

	urlValues := make(map[string][]string)
	form, err := service.ParseFormData(context.Background(), r, urlValues)
	form = service.GetDefaultFormData(form)
	if checkoutForm, ok := form.(formDto.CheckoutFormData); ok {
		if checkoutForm.BillingAddress.Email != "example@mail.dmn" {
			t.Errorf("Wrong mail in data - expected to be initialized")
		}
	} else {
		t.Error(err)
	}
}
