package dto

type (
	// State representation for graphql
	State interface{}

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
	// ShowIFrame state
	ShowIFrame struct {
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
