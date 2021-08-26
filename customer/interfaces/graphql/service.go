package graphql

import (
	// embed schema.graphql
	_ "embed"

	"flamingo.me/graphql"

	"flamingo.me/flamingo-commerce/v3/customer/domain"
	"flamingo.me/flamingo-commerce/v3/customer/interfaces/graphql/dtocustomer"
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
	types.Map("Commerce_Customer_Status_Result", dtocustomer.CustomerStatusResult{})
	types.Map("Commerce_Customer_Result", dtocustomer.CustomerResult{})
	types.Map("Commerce_Customer_PersonData", domain.PersonData{})
	types.Map("Commerce_Customer_Address", domain.Address{})
	types.GoField("Commerce_Customer_Address", "streetNumber", "StreetNr")
	types.Resolve("Query", "Commerce_Customer_Status", CustomerResolver{}, "CommerceCustomerStatus")
	types.Resolve("Query", "Commerce_Customer", CustomerResolver{}, "CommerceCustomer")
}
