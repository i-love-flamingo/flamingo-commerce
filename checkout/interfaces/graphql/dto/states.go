package dto

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// State representation for graphql
	State interface {
		MapFrom(process.Context)
	}

	// Failed state
	Failed struct {
		Name   string
		Reason string
	}
	// Success state
	Success struct {
		Name string
	}
	// Wait state
	Wait struct {
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
)

var (
	_ State = new(Failed)
	_ State = new(Success)
	_ State = new(Wait)
	_ State = new(ShowIframe)
	_ State = new(ShowHTML)
	_ State = new(Redirect)
	_ State = new(PostRedirect)
)

// MapFrom the internal process state to the graphQL state fields
func (s *Failed) MapFrom(pctx process.Context) {
	s.Name = pctx.State
	s.Reason = pctx.FailedReason.Reason()
}

// MapFrom the internal process state to the graphQL state fields
func (s *Success) MapFrom(pctx process.Context) {
	s.Name = pctx.State
}

// MapFrom the internal process state to the graphQL state fields
func (s *Wait) MapFrom(pctx process.Context) {
	s.Name = pctx.State
}

// MapFrom the internal process state to the graphQL state fields
func (s *ShowIframe) MapFrom(pctx process.Context) {
	s.Name = pctx.State
	s.URL = pctx.URL.String()
}

// MapFrom the internal process state to the graphQL state fields
func (s *ShowHTML) MapFrom(pctx process.Context) {
	s.Name = pctx.State
	s.HTML = pctx.DisplayData
}

// MapFrom the internal process state to the graphQL state fields
func (s *Redirect) MapFrom(pctx process.Context) {
	s.Name = pctx.State
	s.URL = pctx.URL.String()
}

// MapFrom the internal process state to the graphQL state fields
func (s *PostRedirect) MapFrom(pctx process.Context) {
	s.Name = pctx.State
	s.URL = pctx.URL.String()
	parameters := make([]FormParameter, 0, len(pctx.FormParameter))
	for key, p := range pctx.FormParameter {
		parameters = append(parameters, FormParameter{
			Key:   key,
			Value: p.Value,
		})
	}

	s.Parameters = parameters
}
