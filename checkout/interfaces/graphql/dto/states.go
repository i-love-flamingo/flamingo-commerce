package dto

import (
	"fmt"
	"net/url"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/states"
)

type (
	// StateProvider returns a state instance
	StateProvider func() map[string]State

	// State representation for graphql
	State interface {
		MapFrom(process.Context)
	}

	// StateMapper to create dto states from context states
	StateMapper struct {
		stateProvider StateProvider
	}

	// Failed state
	Failed struct {
		Name   string
		Reason process.FailedReason
	}
	// Success state
	Success struct {
		Name string
	}
	// Wait state
	Wait struct {
		Name string
	}
	// WaitForCustomer state
	WaitForCustomer struct {
		Name string
	}
	// ShowIframe state
	ShowIframe struct {
		Name string
		URL  string
	}
	// ShowHTML state
	ShowHTML struct {
		Name string
		HTML string
	}
	// Redirect state
	Redirect struct {
		Name string
		URL  string
	}
	// PostRedirect state
	PostRedirect struct {
		Name       string
		URL        string
		Parameters []FormParameter
	}
	// FormParameter holds redirect related form data
	FormParameter struct {
		Key   string
		Value []string
	}

	// ShowWalletPayment state
	ShowWalletPayment struct {
		Name              string
		PaymentMethod     string
		PaymentRequestAPI PaymentRequestAPI
	}

	// TriggerClientSDK state
	TriggerClientSDK struct {
		Name string
		URL  string
		Data string
	}

	// PaymentRequestAPI holds all data needed to create a PaymentRequest
	PaymentRequestAPI struct {
		MethodData            string
		Details               string
		Options               string
		MerchantValidationURL *string
		CompleteURL           string
	}
)

var (
	_ State = new(Failed)
	_ State = new(Success)
	_ State = new(Wait)
	_ State = new(WaitForCustomer)
	_ State = new(ShowIframe)
	_ State = new(ShowHTML)
	_ State = new(Redirect)
	_ State = new(PostRedirect)
	_ State = new(ShowWalletPayment)
	_ State = new(TriggerClientSDK)
)

// MapFrom the internal process state to the graphQL state fields
func (s *Failed) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
	s.Reason = pctx.FailedReason
}

// MapFrom the internal process state to the graphQL state fields
func (s *ShowWalletPayment) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
	if stateData, ok := pctx.CurrentStateData.(states.ShowWalletPaymentData); ok {
		s.PaymentMethod = stateData.UsedPaymentMethod
		s.PaymentRequestAPI = PaymentRequestAPI{
			CompleteURL: func() string {
				if stateData.PaymentRequestAPI.CompleteURL == nil {
					return ""
				}
				return stateData.PaymentRequestAPI.CompleteURL.String()
			}(),
			MerchantValidationURL: func() *string {
				if stateData.PaymentRequestAPI.MerchantValidationURL == nil {
					return nil
				}
				result := stateData.PaymentRequestAPI.MerchantValidationURL.String()
				return &result
			}(),
			Details:    stateData.PaymentRequestAPI.Details,
			Options:    stateData.PaymentRequestAPI.Options,
			MethodData: stateData.PaymentRequestAPI.Methods,
		}
	}
}

// MapFrom the internal process state to the graphQL state fields
func (s *Success) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
}

// MapFrom the internal process state to the graphQL state fields
func (s *Wait) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
}

// MapFrom the internal process state to the graphQL state fields
func (s *WaitForCustomer) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
}

// MapFrom the internal process state to the graphQL state fields
func (s *ShowIframe) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
	if stateData, ok := pctx.CurrentStateData.(*url.URL); ok {
		s.URL = stateData.String()
	}
}

// MapFrom the internal process state to the graphQL state fields
func (s *ShowHTML) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
	if stateData, ok := pctx.CurrentStateData.(string); ok {
		s.HTML = stateData
	}
}

// MapFrom the internal process state to the graphQL state fields
func (s *Redirect) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
	if stateData, ok := pctx.CurrentStateData.(*url.URL); ok {
		s.URL = stateData.String()
	}
}

// MapFrom the internal process state to the graphQL state fields
func (s *PostRedirect) MapFrom(pctx process.Context) {
	s.Name = pctx.CurrentStateName
	if stateData, ok := pctx.CurrentStateData.(states.PostRedirectData); ok {
		s.URL = stateData.URL.String()
		parameters := make([]FormParameter, 0, len(stateData.FormFields))
		for key, p := range stateData.FormFields {
			parameters = append(parameters, FormParameter{
				Key:   key,
				Value: p.Value,
			})
		}
		s.Parameters = parameters
	}
}

// MapFrom the internal process state to the graphQL state fields
func (t *TriggerClientSDK) MapFrom(pctx process.Context) {
	t.Name = pctx.CurrentStateName
	if stateData, ok := pctx.CurrentStateData.(states.TriggerClientSDKData); ok {
		t.URL = stateData.URL.String()
		t.Data = stateData.Data
	}
}

// Inject dependencies
func (sm *StateMapper) Inject(stateProvider StateProvider) *StateMapper {
	sm.stateProvider = stateProvider

	return sm
}

// Map a context into a state
func (sm *StateMapper) Map(pctx process.Context) (State, error) {
	resultState, found := sm.stateProvider()[pctx.CurrentStateName]
	if !found {
		return nil, fmt.Errorf("couldn't map the internal process state %q to a GraphQL state", pctx.CurrentStateName)
	}

	resultState.MapFrom(pctx)
	return resultState, nil
}
