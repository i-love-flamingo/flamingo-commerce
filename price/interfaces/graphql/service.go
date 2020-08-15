package graphql

import (
	"flamingo.me/flamingo-commerce/v3/price/domain"
	"flamingo.me/graphql"
)

//go:generate go run github.com/go-bindata/go-bindata/v3/go-bindata -nometadata -o schema.go -pkg graphql schema.graphql

// Service is the Graphql-Service of this module
type Service struct{}

var _ graphql.Service = new(Service)

// Schema returns graphql schema of this module
func (*Service) Schema() []byte {
	return MustAsset("schema.graphql")
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
