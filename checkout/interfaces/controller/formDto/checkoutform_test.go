package formDto

import (
	"testing"

	form2 "github.com/go-playground/form"
	"go.aoe.com/flamingo/core/customer/domain"
	"go.aoe.com/flamingo/framework/flamingo"
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
	service := CheckoutFormService{
		Customer: &FakeCustomer{},
		Logger:   flamingo.NullLogger{},
		Decoder:  form2.NewDecoder(),
	}

	urlValues := make(map[string][]string)
	form, err := service.ParseFormData(nil, urlValues)
	form = service.GetDefaultFormData(form)
	if checkoutForm, ok := form.(CheckoutFormData); ok {
		if checkoutForm.BillingAddress.Email != "example@mail.dmn" {
			t.Errorf("Wrong mail in data - expected to be initialized")
		}
	} else {
		t.Error(err)
	}
}
