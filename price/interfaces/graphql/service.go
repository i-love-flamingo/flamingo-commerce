package graphql

import (
	// embed schema.graphql
	_ "embed"

	"flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/graphql"
)

// Service is the Graphql-Service of this module
type Service struct{}

var _ graphql.Service = new(Service)

//go:embed schema.graphql
var schema []byte

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return schema
}

// Types configures the GraphQL to Go resolvers
func (*Service) Types(types *graphql.Types) {
	types.Map("Commerce_Price", domain.Price{})
	types.GoField("Commerce_Price", "amount", "FloatAmount")
	types.Map("Commerce_Price_Charges", domain.Charges{})
	types.Map("Commerce_Price_Charge", domain.Charge{})
	types.Map("Commerce_Price_ChargeQualifier", domain.ChargeQualifier{})
	types.Map("Commerce_Price_ChargeQualifierInput", domain.ChargeQualifier{})
}
